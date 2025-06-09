package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenSubKeycloakPrefix     string = "auth0|keycloak|"
	TokenSubPrefix             string = "auth0|"
	TokenSubOauth2Prefix       string = "oauth2|"
	TokenSubNiconicoPrefix     string = "oauth2|niconico|"
	TokenSubNiconicoOnlyPrefix string = "niconico|"
	BearerAuthorizationPrefix  string = "Bearer "
)

// FanclubMemberUserID comment
// en: Auth0UserID contains information auth0.user_id of fanclub_member
// en: KeycloakUserID contains information keycloak.user_id of fanclub_member
type FanclubMemberUserID struct {
	Auth0UserID    string
	KeycloakUserID string
}

// GetCommentUserID comment
// en: fanclub_member always exists at least auth0_user_id or key_cloak_user_id (there are no cases where both are Null).
// en: return user_id information (prefer keycloak_user_id)
// en: using for get comment user id
func (fm FanclubMemberUserID) GetCommentUserID() string {
	if fm.KeycloakUserID != "" {
		return fm.KeycloakUserID
	}
	return fm.Auth0UserID
}

// Decoding JWT to get payload, not verifying JWT
func Decode(jWTToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jWTToken, func(token *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return token, nil
	}
	return nil, err
}

// Generate HS256 JWT token
func GenerateHS256JWT(secret string, keyID string, payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for key, val := range payload {
		claims[key] = val
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = keyID
	signedToken, err := token.SignedString([]byte(secret))
	return signedToken, err
}

// Verify JWT token
func Verify(token, keyString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	publicKey := "-----BEGIN PUBLIC KEY-----\n" + keyString + "\n" + "-----END PUBLIC KEY-----"
	keyByte, keyErr := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if keyErr != nil {
		// if key is not public key, use HS256 for test
		if os.Getenv("DW_TESTS") == "true" && (keyErr == jwt.ErrNotRSAPublicKey || keyErr == jwt.ErrKeyMustBePEMEncoded) {
			_, parseErr := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(keyString), nil
			})
			return claims, parseErr
		}
		return claims, keyErr

	}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return keyByte, nil
	})
	return claims, err
}

// Do not get it directly from the context
// because it is not certain that the API has a `CheckAuthentication` middlerware
func GetKeycloakUserIDByToken(token string, isAuth0Iss bool) (string, error) {
	userToken, err := Decode(token)
	if err != nil {
		return "", err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", nil
	}

	if isAuth0Iss {
		if strings.HasPrefix(sub, TokenSubKeycloakPrefix) {
			return sub[len(TokenSubKeycloakPrefix):], nil
		}

		return "", nil
	}

	return sub, nil
}

// GetAzpByToken comment
// get azp information from token
func GetAzpByToken(token string) (string, error) {
	userToken, err := Decode(token)
	if err != nil {
		return "", err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	azp, ok := claims["azp"].(string)
	if !ok {
		return "", nil
	}

	return azp, nil
}

func GetUserIDByToken(token string) (string, error) {
	userToken, err := Decode(token)
	if err != nil {
		return "", err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}
	keycloakUserID, ok := claims["sub"].(string)
	if !ok {
		return "", nil
	}

	return keycloakUserID, nil
}

func GenerateCommentJWT(secret string, payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for key, val := range payload {
		claims[key] = val
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	signedToken, err := token.SignedString([]byte(secret))
	return signedToken, err
}

func GetRealmNameByToken(token string) (string, error) {
	userToken, err := Decode(token)
	if err != nil {
		return "", err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}
	authUrl, ok := claims["iss"].(string)
	if !ok {
		return "", nil
	}
	stringSplit := strings.Split(authUrl, "/")

	return stringSplit[len(stringSplit)-1], nil
}

func GetEmailByToken(token string) (string, error) {
	userToken, err := Decode(token)
	if err != nil {
		return "", err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}
	keycloakUserID, ok := claims["email"].(string)
	if !ok {
		return "", nil
	}

	return keycloakUserID, nil
}

func VerifyExpiredByToken(token string) (bool, error) {
	userToken, err := Decode(token)
	if err != nil {
		return false, err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	v := jwt.NewValidator(jwt.WithExpirationRequired())
	err = v.Validate(claims)
	return !errors.Is(err, jwt.ErrTokenExpired), nil
}

func GetTokenByHeader(c *gin.Context) (string, bool) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "" && strings.Index(authHeader, BearerAuthorizationPrefix) == 0 {
		return authHeader[7:], true
	}
	return "", false
}

// GenerateES256JWT generate jwt token with algo ES256
// secret is private key pem encode without begin and end type
// func return jwt token sign with private key generate from secret.
func GenerateES256JWT(secret string, keyID string, payload map[string]interface{}) (string, error) {
	// get private key from secret
	block, _ := pem.Decode([]byte("-----BEGIN PRIVATE KEY-----" + "\n" + secret + "\n" + "-----END PRIVATE KEY-----"))
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// create payload
	claims := jwt.MapClaims{}
	for key, val := range payload {
		claims[key] = val
	}

	// create token with algo ES256
	token := jwt.NewWithClaims(jwt.GetSigningMethod("ES256"), claims)
	token.Header["kid"] = keyID

	// sign with private key and return token
	return token.SignedString(privateKey)
}

// GetFanclubGroupIDByToken comment
// en: Get fanclub_group_id from access token
// en: only Auth0 access token has nfc_group_id
func GetFanclubGroupIDByToken(token string) (int, error) {
	userToken, err := Decode(token)
	if err != nil {
		return 0, err
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, nil
	}
	fcGroupID, ok := claims["nfc_group_id"].(float64)
	if !ok {
		return 0, nil
	}

	return int(fcGroupID), nil
}

// IsAuth0Iss comment
// en: verify that access_token had been pulbish by auth0
func IsAuth0Iss(iss string, auth0Domain string, auth0DefaultDomain string) (bool, error) {
	urlDomain, err := url.Parse(iss)
	if err != nil {
		return false, err
	}
	domain := urlDomain.Hostname()
	return domain == auth0Domain || domain == auth0DefaultDomain, nil
}

// ExtractUserIDFromContext comment
// en: Extract user_id (auth0_user_id, keycloak_user_id) information from context
// en: Can only be used when the request has been verified by the middleware
func ExtractUserIDFromContext(c *gin.Context) (FanclubMemberUserID, error) {
	return FanclubMemberUserID{
		Auth0UserID:    c.GetString("auth0_user_id"),
		KeycloakUserID: c.GetString("keycloak_user_id"),
	}, nil
}

// ExtractUserIDFromString comment
// en: Extract user_id (auth0_user_id, keycloak_user_id) information from string (Example: `sub` field from `access_token`)
func ExtractUserIDFromString(sub string) FanclubMemberUserID {
	if len(sub) == 0 {
		return FanclubMemberUserID{}
	}

	// en: Judgment is issued by Auth0
	// en: format "auth0|keycloak|USER-ID-UNIQUE-STRING" if Auth0 uses Keycloak information (migrate Keycloak user to Auth0)
	// en: format "auth0|USER-ID-UNIQUE-STRING" if Auth0 does not use Keycloak information
	// en: format "oauth2|niconico|NICONICO-ID" if user login via niconico
	if strings.HasPrefix(sub, TokenSubKeycloakPrefix) {
		return FanclubMemberUserID{
			Auth0UserID:    sub,
			KeycloakUserID: sub[len(TokenSubKeycloakPrefix):],
		}
	}
	if strings.HasPrefix(sub, TokenSubPrefix) || strings.HasPrefix(sub, TokenSubNiconicoPrefix) {
		return FanclubMemberUserID{
			Auth0UserID: sub,
		}
	}

	// en: Otherwise, judgment is issued by Keycloak
	return FanclubMemberUserID{
		KeycloakUserID: sub,
	}
}

// IsHasUserID comment
// en: Return `true` if exist user_id information
func IsHasUserID(fcMemberUserID FanclubMemberUserID) bool {
	return fcMemberUserID.Auth0UserID != "" || fcMemberUserID.KeycloakUserID != ""
}
