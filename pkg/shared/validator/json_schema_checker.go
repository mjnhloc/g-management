package validator

import (
	"time"

	"g-management/pkg/shared/utils"

	"github.com/xeipuuv/gojsonschema"
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

type MaxLengthInvalidError struct {
	gojsonschema.ResultErrorFields
}

// NewMaxLengthError returns a custom gojsonschema error for validation max length
func NewMaxLengthError(context *gojsonschema.JsonContext, value interface{}, details gojsonschema.ErrorDetails) *MaxLengthInvalidError {
	err := MaxLengthInvalidError{}
	err.SetContext(context)
	err.SetType(maxLengthByte)
	// it is important to use SetDescriptionFormat() as this is used to call SetDescription() after it has been parsed
	// using the description of err will be overridden by this.
	err.SetDescriptionFormat(utils.ErrorInputByteLimit)
	err.SetValue(value)
	err.SetDetails(details)

	return &err
}
