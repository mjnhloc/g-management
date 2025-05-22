package validator

import "regexp"

// regular expressions for validating string
const (
	FullWidth              string = `^[ぁ-んァ-ン一-龥０-９ーａ-ｚ]+$`
	HalfWidth              string = `^[ｧ-ﾝﾞﾟ]+$`
	AllCharactersHalfWidth string = `^[\w\sｧ-ﾝﾞﾟ!@#$%^&*()-_=+{}|;:'",<.>/?]*$`
	AlphaNumericRegex      string = `^[a-zA-Z0-9]+$`
	Domain                 string = `^https:\/\/[a-z0-9-]+([\-\.]{1}[a-z0-9-]+)*\.[a-z]{2,5}(\/[\.a-z0-9_-]+)*$`
	AsciiCharacterPassword string = `^[\x00-\x7F]{8,32}$`
	ExceptSpecialCharacter string = `^[^¥\\ ]+$`
	AtLeastOneNumber       string = `[0-9]`
	AtLeastOneLowerCase    string = `[a-z]`
	AtLeastOneUpperCase    string = `[A-Z]`
	AtLeastOneSpecialChar  string = `[!@#~$%^&*()+|_]{1}`
	Hiragana               string = `^[ぁ-んー＝・ゔ]+$`
	GaMeasurementID        string = `^G-[a-zA-Z0-9_]{0,18}$`
	GTM                    string = `^GTM-[a-zA-Z0-9_]{0,16}$`
	Url                    string = `^https?:\/\/[\w/:%#$@&?()~.=+-]+$`
	IDSnS                  string = `^[A-Za-z0-9-.]*$`
	Auth0Password          string = `^[\w!"#$%&'()*+,–\-./:;<=>?@[\]^_` + "`" + `{|}~]{8,}$`
	HalfWidthOnlyNumber    string = `^[0-9]+$`
)

var (
	regexFullWidth         = regexp.MustCompile(FullWidth)
	regexHalfWidth         = regexp.MustCompile(HalfWidth)
	domain                 = regexp.MustCompile(Domain)
	asciiCharacterPassword = regexp.MustCompile(AsciiCharacterPassword)
	exceptSpecialCharacter = regexp.MustCompile(ExceptSpecialCharacter)
	atLeastOneNumber       = regexp.MustCompile(AtLeastOneNumber)
	atLeastOneLowerCase    = regexp.MustCompile(AtLeastOneLowerCase)
	atLeastOneUpperCase    = regexp.MustCompile(AtLeastOneUpperCase)
	regexHiragana          = regexp.MustCompile(Hiragana)
	regexGaMeasurementID   = regexp.MustCompile(GaMeasurementID)
	regexGTM               = regexp.MustCompile(GTM)
	regexUrl               = regexp.MustCompile(Url)
	regexIDSnS             = regexp.MustCompile(IDSnS)
	auth0Password          = regexp.MustCompile(Auth0Password)
)
