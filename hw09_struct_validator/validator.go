package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrValidate       = errors.New("value is not struct")
	ErrValidateLen    = errors.New("length is not valid")
	ErrValidateMin    = errors.New("value is less than the MIN")
	ErrValidateMax    = errors.New("value is more than the MAX")
	ErrValidateIn     = errors.New("value is not in the range")
	ErrValidateRegexp = errors.New("values is not compile by regexp")
)

func (v ValidationErrors) Error() string {
	msgError := strings.Builder{}

	for _, entry := range v {
		msgError.WriteString(fmt.Sprintf("%v : %v", entry.Field, entry.Err.Error()))
	}

	return msgError.String()
}

func Validate(v interface{}) error {
	var validErrors ValidationErrors

	vValue := reflect.ValueOf(v)
	vType := vValue.Type()

	if kind := vType.Kind(); kind != reflect.Struct {
		return ErrValidate
	}

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)

		if tag := field.Tag.Get("validate"); tag != "" {
			fieldName := field.Name
			fieldValue := vValue.Field(i)

			rules := strings.Split(tag, "|")
			for _, rule := range rules {
				errs := ValidateRule(rule, fieldValue)
				for _, err := range errs {
					validErrors = append(validErrors, ValidationError{fieldName, err})
				}
			}
		}
	}

	return validErrors
}

func ValidateRule(rule string, fieldValue reflect.Value) []error {
	var validErr []error

	ruleSplit := strings.SplitN(rule, ":", 2)
	verb, value := ruleSplit[0], ruleSplit[1]

	switch fieldValue.Kind() { //nolint: exhaustive
	case reflect.Int:
		err := ValidateInt(fieldValue, verb, value)
		if err != nil {
			return []error{err}
		}
	case reflect.String:
		err := ValidateString(fieldValue, verb, value)
		if err != nil {
			return []error{err}
		}
	case reflect.Slice:
		fieldElem := fieldValue.Type().Elem()
		fieldLen := fieldValue.Len()

		switch fieldElem.Kind() { //nolint: exhaustive
		case reflect.Int:
			for i := 0; i < fieldLen; i++ {
				err := ValidateInt(fieldValue.Index(i), verb, value)
				if err != nil {
					validErr = append(validErr, err)
				}
			}
			return validErr
		case reflect.String:
			for i := 0; i < fieldLen; i++ {
				err := ValidateString(fieldValue.Index(i), verb, value)
				if err != nil {
					validErr = append(validErr, err)
				}
			}
			return validErr
		}
	}
	return nil
}

func ValidateInt(fieldValue reflect.Value, verb string, value string) error {
	actualValue := fieldValue.Int()

	switch verb {
	case "min":
		min, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if actualValue < int64(min) {
			return errors.Join(ErrValidateMin, fmt.Errorf("must be %v > %v", actualValue, min))
		}
	case "max":
		max, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if actualValue > int64(max) {
			return errors.Join(ErrValidateMax, fmt.Errorf("must be %v < %v", actualValue, max))
		}
	case "in":
		valRanges := strings.Split(value, ",")
		for _, valRange := range valRanges {
			intRange, err := strconv.Atoi(valRange)
			if err != nil {
				return err
			}

			if actualValue == int64(intRange) {
				return nil
			}
		}

		return errors.Join(ErrValidateIn, fmt.Errorf("%v must be in range %v", actualValue, valRanges))
	}
	return nil
}

func ValidateString(fieldValue reflect.Value, verb string, value string) error {
	actualValue := fieldValue.String()

	switch verb {
	case "len":
		maxLength, err := strconv.Atoi(value)
		actualLen := len(actualValue)
		if err != nil {
			return err
		}
		if actualLen > maxLength {
			return errors.Join(ErrValidateLen, fmt.Errorf("%v must be < %v", actualLen, maxLength))
		}
	case "regexp":
		re, err := regexp.Compile(value)
		if err != nil {
			return err
		}

		if !re.MatchString(actualValue) {
			return errors.Join(ErrValidateRegexp, fmt.Errorf("%v not compile by %v", actualValue, value))
		}
	case "in":
		valRanges := strings.Split(value, ",")

		for _, valRange := range valRanges {
			if actualValue == valRange {
				return nil
			}
		}

		return errors.Join(ErrValidateIn, fmt.Errorf("%v must be in range %v", actualValue, valRanges))
	}
	return nil
}
