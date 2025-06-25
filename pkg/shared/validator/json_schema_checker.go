package validator

import (
	"regexp"
	"time"

	"g-management/pkg/shared/utils"
)

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

const (
	maxLengthByte string = "max_length_byte"
)

type NotFullwidthOrHalfwidthFormatChecker struct{}

func (f NotFullwidthOrHalfwidthFormatChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return !(regexFullWidth.MatchString(inputAsString) || regexHalfWidth.MatchString(inputAsString))
	}

	return false
}

// This is only necessary because the API started to be developed
// with a different format than the specifications and the standard,
// so we want to keep backward-compatibility
type NonStandardDateTimeFormatChecker struct{}

func (f NonStandardDateTimeFormatChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		_, err := time.Parse(utils.FormatDateTimeDb, inputAsString)
		if err == nil {
			return true
		}
	}

	return false
}

type PasswordChecker struct{}

func (f PasswordChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return (asciiCharacterPassword.MatchString(inputAsString) && exceptSpecialCharacter.MatchString(inputAsString))
	}

	return false
}

type StrongPasswordChecker struct{}

func (f StrongPasswordChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return len(inputAsString) >= 12 &&
			asciiCharacterPassword.MatchString(inputAsString) &&
			exceptSpecialCharacter.MatchString(inputAsString) &&
			atLeastOneNumber.MatchString(inputAsString) &&
			atLeastOneLowerCase.MatchString(inputAsString) &&
			atLeastOneUpperCase.MatchString(inputAsString)
	}

	return false
}

type Auth0PasswordChecker struct{}

func (f Auth0PasswordChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return auth0Password.MatchString(inputAsString)
	}

	return false
}

type DomainChecker struct{}

func (f DomainChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return domain.MatchString(inputAsString)
	}

	return false
}

type HiraganaChecker struct{}

func (f HiraganaChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return regexHiragana.MatchString(inputAsString)
	}

	return false
}

type GaMeasurementIDChecker struct{}

func (f GaMeasurementIDChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return regexGaMeasurementID.MatchString(inputAsString)
	}

	return false
}

type GTMChecker struct{}

func (f GTMChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return regexGTM.MatchString(inputAsString)
	}

	return false
}

type MaxLengthChecker struct{}

func (f MaxLengthChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return len([]rune(inputAsString)) <= 50
	}

	return true
}

type UrlChecker struct{}

func (f UrlChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return regexUrl.MatchString(inputAsString)
	}

	return false
}

type IDSnSChecker struct{}

func (f IDSnSChecker) IsFormat(input interface{}) bool {
	if inputAsString, ok := input.(string); ok {
		return regexIDSnS.MatchString(inputAsString)
	}

	return false
}
