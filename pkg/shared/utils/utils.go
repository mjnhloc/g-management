package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"g-management/pkg/shared/wraperror"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rivo/uniseg"
	"golang.org/x/exp/utf8string"
)

const (
	queryDemiliter = "?"
)

// Paging struct
type Paging struct {
	Page    int
	PerPage int
}

type Paginate struct {
	Offset int
	Limit  int
}

// IsValid comment
// en: Check Paginate valid
func (p *Paginate) IsValid() bool {
	return p.Offset >= 0 && p.Limit > 0
}

const (
	IOS            = "iOS"
	Android        = "Android"
	WebBrowser     = "null"
	WebView        = "NFC/WebView"
	WebViewBrowser = "WebBrowser"
	MobileConstant = "mobile"
)

// GetValueFieldByName func
func GetValueFieldByName(v interface{}, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}

// EncryptPassword method SHA256 func
func EncryptPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(h[:])
}

// DerefString func convert string pointer to string
func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

// DerefInt func convert int pointer to int
func DerefInt(i *int) int {
	if i != nil {
		return *i
	}

	return 0
}

// ReplaceSpecialCharacters search like
func ReplaceSpecialCharacters(str string) string {
	r := strings.NewReplacer("_", "\\_", "%", "\\%", "'", "''", "\\", "\\\\")
	return r.Replace(str)
}

func ConvertToString(input interface{}) (string, error) {
	var output string
	var err error
	switch v := input.(type) {
	case int:
		output = strconv.Itoa(v)
	case bool:
		output = strconv.FormatBool(v)
	case float64:
		output = fmt.Sprintf("%f", v)
	case string:
		output = v
	default:
		err = fmt.Errorf("undefined type to convert %T", v)
	}

	return output, err
}

func IsUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}

func CheckInputIDListWithCurrent(inputIDList, currentIDList []int, errMsgSprintf string, err error) error {
	if len(currentIDList) == 0 {
		msg := make(map[string]interface{})
		for i := range inputIDList {
			fieldName := fmt.Sprintf(errMsgSprintf, i)
			msg[fieldName] = ErrorInputFail
		}

		return wraperror.NewValidationError(
			msg,
			err,
		)
	}
	exists := make(map[int]bool)
	for _, value := range currentIDList {
		exists[value] = true
	}
	outputMsg := make(map[string]interface{})
	for i, value := range inputIDList {
		if !exists[value] {
			fieldName := fmt.Sprintf(errMsgSprintf, i)
			outputMsg[fieldName] = ErrorInputFail
		}
	}

	if len(outputMsg) != 0 {
		return wraperror.NewValidationError(
			outputMsg,
			err,
		)
	}

	return nil
}

/**
 * CheckValidationPageAndPerPage check validation page and per page
 * @return paging include page and per page type integer
 * @return msgErr message error
 */
func CheckValidationPageAndPerPage(
	c *gin.Context,
	defaultPerPage int,
) (Paging, map[string]interface{}) {
	return CheckValidationPageAndPerPageWithMaximumPerPage(
		c,
		defaultPerPage,
		MaxPerPage,
	)
}

/**
 * CheckValidationPageAndPerPageWithMaximumPerPage check validation page and per page with maximum per page
 * @return paging include page and per page type integer
 * @return msgErr message error
 */
func CheckValidationPageAndPerPageWithMaximumPerPage(
	c *gin.Context,
	defaultPerPage int,
	maxPerPage int,
) (Paging, map[string]interface{}) {
	msgErr := map[string]interface{}{}
	paging := Paging{
		Page:    DefaultPage,
		PerPage: defaultPerPage,
	}

	if pageAsString, ok := c.GetQuery("page"); ok {
		if len(pageAsString) > 0 {
			page, err := strconv.Atoi(pageAsString)
			if err != nil {
				msgErr["page"] = ErrorInputFail
			}
			if page < MinPaging {
				msgErr["page"] = ErrorInputFail
			}
			paging.Page = page
		}
	}

	if perPageAsString, ok := c.GetQuery("per_page"); ok {
		if len(perPageAsString) > 0 {
			perPage, err := strconv.Atoi(perPageAsString)
			if err != nil {
				msgErr["per_page"] = ErrorInputFail
			}
			if perPage < MinPaging || perPage > maxPerPage {
				msgErr["per_page"] = ErrorInputFail
			}
			paging.PerPage = perPage
		}
	}

	return paging, msgErr
}

func CheckIdsIsDelete(idsNew []map[string]interface{}, idsOld []int) []int {
	idExists := map[int]bool{}
	var idDeletes []int

	for _, IDNew := range idsNew {
		if IDNew["id"] != nil {
			idExists[int(IDNew["id"].(float64))] = true
		}
	}

	for _, IDOld := range idsOld {
		_, ok := idExists[IDOld]
		if !ok {
			idDeletes = append(idDeletes, IDOld)
		}
	}

	return idDeletes
}

func CheckElementExistInSlice(slice, val interface{}) (int, bool) {
	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		s := reflect.ValueOf(slice)

		for i := 0; i < s.Len(); i++ {
			if val == s.Index(i).Interface() {
				return i, true
			}
		}
	}
	return -1, false
}

func CheckSortEnumFormat(sort string, fields []string) string {
	sortArr := []rune(sort)
	var orderBy string
	var orderField string
	if string(sortArr[0]) == "-" {
		orderBy = "desc"
		orderField = string(sortArr[1:])
	} else {
		orderBy = "asc"
		orderField = sort
	}
	mapFields := map[string]bool{}
	for _, field := range fields {
		if field != "" {
			mapFields[field] = true
		}
	}
	if mapFields[orderField] {
		return fmt.Sprintf("%v %v", orderField, orderBy)
	}

	return ""
}

// func GetKeycloakUserIDByHeader(c *gin.Context) (string, error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	var token string
// 	if authHeader != "" && strings.Index(authHeader, jwt.BearerAuthorizationPrefix) == 0 {
// 		token = authHeader[7:]
// 	} else {
// 		return "", nil
// 	}

// 	isAuth0Iss := c.GetString("auth0_user_id") != ""

// 	keyCloakUserID, err := jwt.GetKeycloakUserIDByToken(token, isAuth0Iss)
// 	if err != nil {
// 		return "", err
// 	}
// 	return keyCloakUserID, nil
// }

// func GetRealmNameByHeader(c *gin.Context) (string, error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	var token string
// 	if authHeader != "" && strings.Index(authHeader, jwt.BearerAuthorizationPrefix) == 0 {
// 		token = authHeader[7:]
// 	} else {
// 		return "", nil
// 	}
// 	realmName, err := jwt.GetRealmNameByToken(token)
// 	if err != nil {
// 		return "", err
// 	}
// 	return realmName, nil
// }

// GetUserIDByHeader comment
// en: Extract user_id (auth0_user_id, keycloak_user_id) information from Context
// en: Can only be used when the request has been verified by the middleware.
// func GetUserIDByHeader(c *gin.Context) (jwt.FanclubMemberUserID, error) {
// 	userID, err := jwt.ExtractUserIDFromContext(c)
// 	if err != nil {
// 		return jwt.FanclubMemberUserID{}, err
// 	}
// 	return userID, nil
// }

// GetFanclubGroupIDByHeader comment
// en: Get fanclub_group_id from header.access_token
// en: only Auth0 access token has nfc_group_id
// func GetFanclubGroupIDByHeader(c *gin.Context) (int, error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	var token string
// 	if authHeader != "" && strings.Index(authHeader, jwt.BearerAuthorizationPrefix) == 0 {
// 		token = authHeader[7:]
// 	} else {
// 		return 0, nil
// 	}
// 	fcGroupID, err := jwt.GetFanclubGroupIDByToken(token)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return fcGroupID, nil
// }

// IsAuth0IssByHeader comment
// en: Return true if the token is issued by auth0. Otherwise, return false
func IsAuth0IssByHeader(c *gin.Context) (bool, error) {
	// en: Access token auth0 will have auth0_user_id info in context
	isAuth0Iss := c.GetString("auth0_user_id") != ""

	return isAuth0Iss, nil
}

// GetAzpByHeader comment
// get azp information from header
// func GetAzpByHeader(c *gin.Context) (string, error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	var token string
// 	if authHeader != "" && strings.Index(authHeader, jwt.BearerAuthorizationPrefix) == 0 {
// 		token = authHeader[7:]
// 	} else {
// 		return "", nil
// 	}

// 	azp, err := jwt.GetAzpByToken(token)
// 	if err != nil {
// 		return "", err
// 	}

// 	return azp, nil
// }

// GetSubByHeader comment
// en: get access_token.sub
// func GetSubByHeader(c *gin.Context) (string, error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	var token string
// 	if authHeader != "" && strings.Index(authHeader, jwt.BearerAuthorizationPrefix) == 0 {
// 		token = authHeader[7:]
// 	} else {
// 		return "", nil
// 	}
// 	sub, err := jwt.GetUserIDByToken(token)
// 	if err != nil {
// 		return "", err
// 	}
// 	return sub, nil
// }

// GetEmailByHeader comment
// en: get email from access_token of auth0/keycloak
// func GetEmailByHeader(c *gin.Context) (string, error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	var token string
// 	if authHeader != "" && strings.Index(authHeader, jwt.BearerAuthorizationPrefix) == 0 {
// 		token = authHeader[7:]
// 	} else {
// 		return "", nil
// 	}
// 	keyCloakUserID, err := jwt.GetEmailByToken(token)
// 	if err != nil {
// 		return "", err
// 	}
// 	return keyCloakUserID, nil
// }

func CheckSliceDiff(inputList, currentList []int) []int {
	exists := make(map[int]bool)
	var output []int
	for _, value := range currentList {
		exists[value] = true
	}

	for _, value := range inputList {
		if !exists[value] {
			output = append(output, value)
		}
	}

	return output
}

// UnorderedEqualSliceInt This check only true if two slice is unique
func UnorderedEqualSliceInt(firstSlice, secondSlice []int) bool {
	if len(firstSlice) != len(secondSlice) {
		return false
	}
	exists := make(map[int]bool)
	for _, value := range firstSlice {
		exists[value] = true
	}
	for _, value := range secondSlice {
		if !exists[value] {
			return false
		}
	}
	return true
}

func GetSliceSameElements(inputList, currentList []int) []int {
	exists := make(map[int]bool)
	var output []int
	for _, value := range currentList {
		exists[value] = true
	}

	for _, value := range inputList {
		if exists[value] {
			output = append(output, value)
		}
	}

	return output
}

// DerefInt func convert int pointer to int
func DerefFloat64(i *float64) float64 {
	if i != nil {
		return *i
	}

	return 0
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ConvertPointerFloat64ToString(dataInput *float64, numberDigitAfterDecimal int) string {
	if dataInput != nil {
		return fmt.Sprintf("%."+strconv.Itoa(numberDigitAfterDecimal)+"f", *dataInput)
	}
	return ""
}

func CheckKeyMatch(key1, key2 string) (bool, error) {
	key2 = strings.ReplaceAll(key2, "/*", "/.*")

	re := regexp.MustCompile(`:[^/]+`)
	key2 = re.ReplaceAllString(key2, "$1[^/]+$2")

	return regexp.MatchString("^"+key2+"$", key1)
}

func GenerateArchivedPath(contentCode string) string {
	return JstreamSftpFolderPublicPath + "/" + contentCode + JstreamSftpFolderArchivedPath + "/"
}

func GenerateTranscodedPath(contentCode string) string {
	return JstreamSftpFolderPublicPath + "/" + contentCode + JstreamSftpFolderTranscodedPath + "/"
}

func MakeTimestampMillisecond(timeInput time.Time) int64 {
	return timeInput.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func GenerateRandomHash() (string, error) {
	b := [16]byte{}
	n, err := rand.Read(b[:])
	if err != nil {
		return "", err
	}

	if n != len(b) {
		return "", fmt.Errorf("invalid length : %v", n)
	}

	return fmt.Sprintf("%x", b), nil
}

func NormalizeRelativePath(path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	return path
}

// DerefBool func convert bool pointer to bool
func DerefBool(i *bool) bool {
	if i != nil {
		return *i
	}

	return false
}

func RemoveSpace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

func SplitStringByWhiteSpace(stringInput string) []string {
	return strings.Fields(stringInput)
}

func GetStringSliceHasUniqueValue(stringInputs []string) []string {
	mapStringExist := map[string]bool{}
	listUniqueString := []string{}
	for _, itemStringsInput := range stringInputs {
		if _, exist := mapStringExist[itemStringsInput]; !exist {
			mapStringExist[itemStringsInput] = true
			listUniqueString = append(listUniqueString, itemStringsInput)
		}
	}
	return listUniqueString
}

func NewBoolPointer(v bool) *bool {
	return &v
}

func NewIntPointer(v int) *int {
	return &v
}

func NewStringPointer(v string) *string {
	return &v
}

func GetFirstNotNilInt(values ...*int) *int {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil
}

// GetStringCount gets the number of characters in a person's view.
func GetStringCount(str string) int {
	return uniseg.GraphemeClusterCount(str)
}

// SliceUTF8 gets the characters from the beginning to the specified position,
// in UTF-8-based characters.
func SliceUTF8(str string, pos int, addString string) string {
	s := utf8string.NewString(str)
	length := GetStringCount(str)
	if pos >= length {
		return s.Slice(0, length)
	}
	return s.Slice(0, pos) + addString
}

// ReadCloserToString converts io.ReadClose to string.
func ReadCloserToString(r *io.ReadCloser) (string, error) {
	if r == nil {
		return "", fmt.Errorf("r argument equals nil")
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(*r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func CheckInputIDListWithCurrentArticleAuth(inputIDList, currentIDList []int, errMsgSprintf string, err error) error {
	if len(currentIDList) == 0 && len(inputIDList) > 0 {
		msg := make(map[string]interface{})
		for i := range inputIDList {
			fieldName := fmt.Sprintf(errMsgSprintf, i)
			msg[fieldName] = ErrorInputFail
		}

		return wraperror.NewValidationError(
			msg,
			err,
		)
	}
	exists := make(map[int]bool)
	for _, value := range currentIDList {
		exists[value] = true
	}
	outputMsg := make(map[string]interface{})
	for i, value := range inputIDList {
		if !exists[value] {
			fieldName := fmt.Sprintf(errMsgSprintf, i)
			outputMsg[fieldName] = ErrorInputFail
		}
	}

	if len(outputMsg) != 0 {
		return wraperror.NewValidationError(
			outputMsg,
			err,
		)
	}

	return nil
}

func ReplaceLineBreak(str string) string {
	regexpNewLine := regexp.MustCompile(`\r?\n`)
	return regexpNewLine.ReplaceAllString(str, "￥ｎ")
}

func RoundFloat64ToString(number *float64) string {
	return fmt.Sprintf("%d", int(math.Round(DerefFloat64(number))))
}

type UniqueIDs struct {
	IDs      []int
	InnerMap map[int]bool
}

func (uids *UniqueIDs) AppendIDIfNotExist(ids ...int) {
	for _, id := range ids {
		if !uids.InnerMap[id] {
			uids.InnerMap[id] = true
			uids.IDs = append(uids.IDs, id)
		}
	}
}

type UniqueStringIDs struct {
	IDs      []string
	InnerMap map[string]bool
}

func (uids *UniqueStringIDs) AppendStringIDIfNotExist(ids ...string) {
	for _, id := range ids {
		if !uids.InnerMap[id] {
			uids.InnerMap[id] = true
			uids.IDs = append(uids.IDs, id)
		}
	}
}

func CopyConditions(conditions map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range conditions {
		result[k] = v
	}
	return result
}

func GetPathAndParams(urlPath string) (string, string) {
	idx := strings.Index(urlPath, queryDemiliter)
	if idx >= 0 {
		return urlPath[:idx], urlPath[idx+1:]
	}
	return urlPath, ""
}

func GetIntSliceFromStringIgnoreOther(s string, delimiter string) ([]int, error) {
	stringArr := strings.Split(s, delimiter)
	var result []int
	for _, str := range stringArr {
		i, err := strconv.Atoi(str)
		if err == nil {
			result = append(result, i)
		}
	}
	return result, nil
}

func GetSliceInterfaceFromString(s string, delimiter string) []interface{} {
	stringArr := strings.Split(s, delimiter)
	var result []interface{}
	for _, str := range stringArr {
		i, err := strconv.Atoi(str)
		if err != nil {
			result = append(result, str)
		} else {
			result = append(result, i)
		}
	}
	return result
}

// SpliceInt func splice integer slice in multiple smaller slices with chunk size
// Example: []int{1, 2, 3, 4, 5}, chunkSize = 2
// [][]int{[]int{1, 2}, []int{3, 4}, []int{5}}
func SpliceInt(arr []int, chunkSize int) [][]int {
	if len(arr) == 0 {
		return [][]int{}
	}
	var divided [][]int
	if chunkSize < 1 {
		divided = append(divided, arr)
		return divided
	}
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize

		if end > len(arr) {
			end = len(arr)
		}

		divided = append(divided, arr[i:end])
	}
	return divided
}

// GetMaxLength return greater number between n and len(arr)
func GetMaxLength(n int, arr []int) int {
	if len(arr) > n {
		return len(arr)
	}
	return n
}

// FillArrayWithZero fill arr with zeros until arr length large than n
func FillArrayWithZero(maxLength int, arr []int) []int {
	if len(arr) >= maxLength {
		return arr
	}

	for i := 0; i < maxLength-len(arr); i++ {
		arr = append(arr, 0)
	}

	return arr
}

// ReplaceLineBreakWithEmpty return string without line break
// input: string
// output: string without line break
func ReplaceLineBreakWithEmpty(str string) string {
	r := strings.NewReplacer("\n", "")
	return r.Replace(str)
}

func CheckSliceIntEqualUnOrder(firstSlice, secondSlice []int) bool {
	if len(firstSlice) != len(secondSlice) {
		return false
	}
	count := make(map[int]int)
	for _, value := range firstSlice {
		count[value]++
	}
	for _, value := range secondSlice {
		if _, ok := count[value]; ok {
			count[value]--
			continue
		}
		return false
	}
	for _, v := range count {
		if v != 0 {
			return false
		}
	}
	return true
}

func GetPaginate(countWithOrder []int, paging Paging) map[int]Paginate {
	paginateMap := make(map[int]Paginate)
	incrementCount := make([]int, len(countWithOrder))
	for i := 0; i < len(incrementCount); i++ {
		if i == 0 {
			incrementCount[i] = countWithOrder[i]
		} else {
			incrementCount[i] = countWithOrder[i] + incrementCount[i-1]
		}
	}

	startOffSet, endOffSet := GetPaginateInfo(paging)
	var start int
	var end int
	var startValidIdx int
	endValidIdx := len(incrementCount) - 1

	for ; startValidIdx < len(incrementCount); startValidIdx++ {
		if incrementCount[startValidIdx] != 0 {
			break
		}
	}

	for ; endValidIdx >= 1; endValidIdx-- {
		if incrementCount[endValidIdx] != incrementCount[endValidIdx-1] {
			break
		}
	}

	if startValidIdx >= len(incrementCount) || startOffSet > incrementCount[endValidIdx] || startValidIdx > endValidIdx {
		return map[int]Paginate{}
	}

	for i := startValidIdx; i <= endValidIdx; i++ {
		if incrementCount[i] >= startOffSet {
			start = i
			break
		}
	}

	for i := startValidIdx; i <= endValidIdx; i++ {
		if incrementCount[i] >= endOffSet || i == endValidIdx {
			end = i
			break
		}
	}

	if start == end {
		startIndex := 1
		if start > 0 {
			startIndex = incrementCount[start-1] + 1
		}
		paginateMap[start] = Paginate{
			Offset: startOffSet - startIndex,
			Limit:  int(math.Min(float64(endOffSet-startOffSet+1), float64(incrementCount[start]-startOffSet+1))),
		}
	} else {
		for i := start; i <= end; i++ {
			switch {
			case i == start:
				tempPaginate := Paginate{
					Limit: incrementCount[i] - startOffSet + 1,
				}
				if i == 0 {
					tempPaginate.Offset = startOffSet - 1
				} else {
					tempPaginate.Offset = startOffSet - incrementCount[i-1] - 1
				}
				paginateMap[i] = tempPaginate
			case i == end:
				paginateMap[i] = Paginate{
					Offset: 0,
					Limit:  int(math.Min(float64(endOffSet-incrementCount[i-1]), float64(incrementCount[i]-incrementCount[i-1]))),
				}
			case incrementCount[i]-incrementCount[i-1] != 0:
				paginateMap[i] = Paginate{
					Offset: 0,
					Limit:  incrementCount[i] - incrementCount[i-1],
				}
			}
		}
	}

	return paginateMap
}

func GetPaginateInfo(paging Paging) (int, int) {
	page := 1
	if paging.Page > 0 {
		page = paging.Page
	}

	pageSize := DefaultPerPage
	if paging.PerPage > 0 {
		pageSize = paging.PerPage
	}

	offset := (page - 1) * pageSize
	return offset + 1, offset + pageSize
}

var ErrDeviceTypeUnknown = errors.New("device type is not in list [ios, android, web browser]")

func GetDeviceTypeByHeader(c *gin.Context) (string, error) {
	deviceType := c.Request.Header.Get("fc_use_device")
	_, ok := CheckElementExistInSlice([]string{
		IOS,
		Android,
		WebBrowser,
	}, deviceType)
	if !ok {
		return "", ErrDeviceTypeUnknown
	}

	return deviceType, nil
}

func BuildString(ids []int) string {
	buildString := strings.Builder{}
	if len(ids) == 1 {
		_, err := buildString.WriteString("(?)")
		if err != nil {
			return buildString.String()
		}
	} else {
		for i := range ids {
			switch i {
			case 0:
				_, err := buildString.WriteString("(?")
				if err != nil {
					return buildString.String()
				}
			case len(ids) - 1:
				_, err := buildString.WriteString(", ?)")
				if err != nil {
					return buildString.String()
				}
			default:
				_, err := buildString.WriteString(", ?")
				if err != nil {
					return buildString.String()
				}
			}
		}
	}
	return buildString.String()
}

// GetUserAgentTypeByHeader get user-agent from header
func GetUserAgentTypeByHeader(c *gin.Context) string {
	var userAgent string
	userAgentInput := c.Request.Header.Get("user-agent")
	if userAgentInput != "" {
		if strings.Contains(userAgentInput, WebView) {
			userAgent = WebView
		} else {
			userAgent = WebViewBrowser
		}
	} else {
		return userAgentInput
	}
	return userAgent
}

var ErrConflictUserAgentAndFcUseDevice = errors.New("user-agent and fc_use_device are conflict")

// ValidateUserAgentAndFcUseDevice validate user-agent and fc_use_device
func ValidateUserAgentAndFcUseDevice(userAgent string, fcUseDevice string) error {
	if (userAgent == WebView && fcUseDevice == WebBrowser) ||
		(userAgent == WebViewBrowser && (fcUseDevice == IOS || fcUseDevice == Android)) {
		return ErrConflictUserAgentAndFcUseDevice
	}
	return nil
}

func ReplaceLineBreakWithRegexp(re *regexp.Regexp, str string) string {
	return re.ReplaceAllString(str, "￥ｎ")
}

func ReplaceLineBreakWithEmptyWithRegexp(re *regexp.Regexp, str string) string {
	return re.ReplaceAllString(str, "")
}

// Range create int array with item values from start to end
func Range(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

// RangeN create int array with item values from 0 to end
func RangeN(end int) []int {
	return Range(0, end)
}

// StripHtmlTagsFromString comment
// en: use bluemonday package to strip all html tags in string
func StripHtmlTagsFromString(data string) string {
	p := bluemonday.StrictPolicy()
	data = p.Sanitize(data)
	return data
}

// GetHostAndAppNameFromRtmpUrl comment
// en: return host from rtmp url
func GetHostAndAppNameFromRtmpUrl(rtmpUrl string) (string, string, error) {
	u, err := url.Parse(rtmpUrl)
	if err != nil {
		return "", "", err
	}

	arr := strings.Split(rtmpUrl, "/")
	appName := arr[len(arr)-1]

	if u == nil || u.Host == "" {
		return "", "", errors.New("invalid url")
	}

	return u.Host, appName, nil
}

// GetDomainFromHttpUrl comment.
// en: extract and return domain from http url.
func GetDomainFromHttpUrl(urlString string) (string, error) {
	// extract domain if it has scheme.
	if strings.Contains(urlString, fmt.Sprintf("%s://", DefaultHttpScheme)) || strings.Contains(urlString, fmt.Sprintf("%s://", DefaultHttpsScheme)) {
		u, err := url.Parse(urlString)
		if err != nil {
			return "", err
		}
		return u.Hostname(), nil
	}

	return urlString, nil
}

// GetSubDomainForCookieFromDomain comment.
// en: extract subDomain for cookies.
func GetSubDomainForCookieFromDomain(domain string) string {
	hostPart := strings.Split(domain, ".")
	if len(hostPart) > 2 {
		return fmt.Sprintf(".%s", strings.Join(hostPart[1:], "."))
	}

	return domain
}

// GetUserIDConditions comment
// en: get user_id map query conditions
func GetUserIDConditions(auth0UserID, keycloakUserID string) map[string]interface{} {
	conditions := map[string]interface{}{}
	if auth0UserID != "" {
		conditions["auth0_user_id"] = auth0UserID
	}
	if keycloakUserID != "" {
		conditions["keycloak_user_id"] = keycloakUserID
	}

	return conditions
}

// GetGroupConditions comment
// en: get group map query conditions
// func GetGroupConditions(fanclubGroupId int, realmName string) map[string]interface{} {
// 	conditions := map[string]interface{}{}
// 	if fanclubGroupId != 0 {
// 		conditions["id"] = fanclubGroupId
// 	} else {
// 		newRealmName := realmName
// 		if strings.Contains(newRealmName, fanclubMembersUtils.CPRealmSuffix) {
// 			newRealmName = strings.TrimSuffix(newRealmName, fanclubMembersUtils.CPRealmSuffix)
// 		}

// 		conditions["realm_name"] = newRealmName
// 	}

// 	return conditions
// }

// FindIndexByColumnKey comment
// en: return index of element
func FindIndexByColumnKey(s []string, element string) int {
	for index, v := range s {
		if v == element {
			return index
		}
	}

	return -1
}

// RemoveIndex comment
// en: remove element by index
func RemoveIndex(s []string, index int) []string {
	if index < 0 || index >= len(s) {
		return s
	}

	return append(s[:index], s[index+1:]...)
}

// GetMailSignatureHtml comment
// en: Get mail signature html
func GetMailSignatureHtml(mailSignature string) string {
	if mailSignature != "" {
		regexpNewLine := regexp.MustCompile(`\r?\n`)
		mailSignature = regexpNewLine.ReplaceAllString(mailSignature, "<br>")
	}

	return mailSignature
}

// ReplaceSpecialCharactersOpenSearch comment
// en: replace special character for search OpenSearch
func ReplaceSpecialCharactersOpenSearch(str string) string {
	r := strings.NewReplacer(`"`, `\"`, `\`, `\\`)
	return r.Replace(str)
}

// GetUserAgentByHeader get user-agent from header
// en: func get user_agent from header
func GetUserAgentByHeader(c *gin.Context) string {
	return c.Request.Header.Get("user-agent")
}

// CheckMobileDevice comment
// en: check if user agent contain mobile
func CheckMobileDevice(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	return strings.Contains(userAgent, MobileConstant)
}

// ConvertJsonToArrayInt comment
// en: convert from json to array
// en: input : string array or nil (if input is not nil or string array return error)
// en: output : return []int (empty array if input is nil)
func ConvertJsonToArrayInt(jsonString *string) ([]int, error) {
	var result []int
	if jsonString != nil {
		err := json.Unmarshal([]byte(*jsonString), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ConvertArrayIntToJSON comment
// en: convert from json to array
// en: input : string array
// en: output : return *string
func ConvertArrayIntToJSON(numbers []int) (*string, error) {
	jsonData, err := json.Marshal(numbers)
	if err != nil {
		return nil, err
	}
	result := string(jsonData)
	return &result, nil
}

// CreatePlaceholdersAndArgs comment
// en: Generates a string of placeholders for SQL queries and a slice of arguments corresponding to those placeholders.
// en: It takes a slice of integers (ids) and returns a string of placeholders and a slice of interfaces containing the ids.
// en: Parameters:
// en: - ids: A slice of integers representing the IDs to be used in the query.
// en: Returns:
// en: - A comma-separated string.
func CreatePlaceholdersAndArgs(ids []int) string {
	if len(ids) == 0 {
		return NullStatusString
	}

	placeholders := make([]string, len(ids))
	for i, id := range ids {
		placeholders[i] = strconv.Itoa(id)
	}

	return strings.Join(placeholders, ",")
}

// GetIntFromAny comment
// en: Tries to convert a value of type any to an int.
// en: It supports float64 and *float64 types.
// en: If the input is nil or not a supported type, it returns 0.
func GetIntFromAny(val any) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case *float64:
		if v != nil {
			return int(*v)
		}
	}
	return 0
}
