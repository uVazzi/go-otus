package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type rule struct {
	ruleKey   string
	ruleValue string
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrNotStruct            = errors.New("validation data is not struct")
	ErrIncorrectRule        = errors.New("incorrect rule")
	ErrIncorrectType        = errors.New("incorrect type")
	ErrValidateIntMin       = errors.New("int value less than min")
	ErrValidateIntMax       = errors.New("int value is greater than max")
	ErrValidateIntIn        = errors.New("invalid int value (in)")
	ErrValidateStringLen    = errors.New("invalid string length (len)")
	ErrValidateStringRegexp = errors.New("invalid string value (regexp)")
	ErrValidateStringIn     = errors.New("invalid string value (in)")
)

func (v ValidationErrors) Error() string {
	validationErrors := make([]string, 0, len(v))
	for _, err := range v {
		validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", err.Field, err.Err))
	}
	return strings.Join(validationErrors, "; ")
}

func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)
	reflectType := reflect.TypeOf(v)

	if reflectValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	validationErrors := ValidationErrors{}

	for i := 0; i < reflectValue.NumField(); i++ {
		dataRuleItem := reflectType.Field(i).Tag.Get("validate")
		if dataRuleItem == "" {
			continue
		}
		var rules []rule
		ruleData := strings.Split(dataRuleItem, "|")
		for _, ruleItem := range ruleData {
			dataValue := strings.Split(ruleItem, ":")
			if len(dataValue) != 2 {
				return ErrIncorrectRule
			}
			rules = append(rules, rule{dataValue[0], dataValue[1]})
		}

		err := validateItem(reflectType.Field(i).Name, reflectValue.Field(i), rules, &validationErrors)
		if err != nil {
			return err
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateItem(fieldName string, itemValue reflect.Value, rules []rule, validationErrors *ValidationErrors) error {
	//nolint:exhaustive
	switch itemValue.Kind() {
	case reflect.Int:
		return validateInt(fieldName, int(itemValue.Int()), rules, validationErrors)
	case reflect.String:
		return validateString(fieldName, itemValue.String(), rules, validationErrors)
	case reflect.Slice:

		switch itemValue.Type().Elem().Kind() {
		case reflect.Int:
			for i := 0; i < itemValue.Len(); i++ {
				subFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
				err := validateInt(subFieldName, int(itemValue.Index(i).Int()), rules, validationErrors)
				if err != nil {
					return err
				}
			}
		case reflect.String:
			for i := 0; i < itemValue.Len(); i++ {
				subFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
				err := validateString(subFieldName, itemValue.Index(i).String(), rules, validationErrors)
				if err != nil {
					return err
				}
			}
		default:
			return ErrIncorrectType
		}

	default:
		return ErrIncorrectType
	}

	return nil
}

func validateInt(fieldName string, fieldValue int, rules []rule, validationErrors *ValidationErrors) error {
	for _, ruleItem := range rules {
		switch ruleItem.ruleKey {
		case "min":
			ruleValue, err := strconv.Atoi(ruleItem.ruleValue)
			if err != nil {
				return ErrIncorrectRule
			}
			if fieldValue < ruleValue {
				*validationErrors = append(*validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrValidateIntMin,
				})
			}
		case "max":
			ruleValue, err := strconv.Atoi(ruleItem.ruleValue)
			if err != nil {
				return ErrIncorrectRule
			}
			if fieldValue > ruleValue {
				*validationErrors = append(*validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrValidateIntMax,
				})
			}
		case "in":
			ruleValue := strings.Split(ruleItem.ruleValue, ",")
			isFound := false
			for _, checkValue := range ruleValue {
				checkValueInt, err := strconv.Atoi(checkValue)
				if err != nil {
					return ErrIncorrectRule
				}
				if fieldValue == checkValueInt {
					isFound = true
					break
				}
			}
			if !isFound {
				*validationErrors = append(*validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrValidateIntIn,
				})
			}
		default:
			return ErrIncorrectRule
		}
	}

	return nil
}

func validateString(fieldName string, fieldValue string, rules []rule, validationErrors *ValidationErrors) error {
	for _, ruleItem := range rules {
		switch ruleItem.ruleKey {
		case "len":
			ruleValue, err := strconv.Atoi(ruleItem.ruleValue)
			if err != nil {
				return ErrIncorrectRule
			}
			if len(fieldValue) != ruleValue {
				*validationErrors = append(*validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrValidateStringLen,
				})
			}
		case "regexp":
			re, err := regexp.Compile(ruleItem.ruleValue)
			if err != nil {
				return ErrIncorrectRule
			}
			if !re.MatchString(fieldValue) {
				*validationErrors = append(*validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrValidateStringRegexp,
				})
			}
		case "in":
			ruleValue := strings.Split(ruleItem.ruleValue, ",")
			isFound := false
			for _, checkValue := range ruleValue {
				if fieldValue == checkValue {
					isFound = true
					break
				}
			}
			if !isFound {
				*validationErrors = append(*validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrValidateStringIn,
				})
			}
		default:
			return ErrIncorrectRule
		}
	}

	return nil
}
