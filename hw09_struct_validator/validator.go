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

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for i, err := range v {
		errs[i] = fmt.Sprintf("%s: %s", err.Field, err.Err)
	}

	return strings.Join(errs, ", ")
}

var (
	ErrIsNotStruct = errors.New("отсутствует структура")

	ErrValidLength   = errors.New("длина строки не равна")
	ErrValidRegExp   = errors.New("регулярное выражение не совпадает")
	ErrValidStrNotIn = errors.New("строка должна входить в множество строк")
	ErrValidMin      = errors.New("число не может быть меньше")
	ErrValidMax      = errors.New("число не может быть больше")
	ErrValidIntNotIn = errors.New("число должно входить в множество чисел")
)

func isErrValid(err error) bool {
	if errors.Is(err, ErrValidLength) ||
		errors.Is(err, ErrValidRegExp) ||
		errors.Is(err, ErrValidStrNotIn) ||
		errors.Is(err, ErrValidMin) ||
		errors.Is(err, ErrValidMax) ||
		errors.Is(err, ErrValidIntNotIn) {
		return true
	}
	return false
}

func Validate(v interface{}) error {
	refValue := reflect.ValueOf(v)
	refType := reflect.TypeOf(v)

	if refValue.Kind() != reflect.Struct {
		return ErrIsNotStruct
	}

	var errs ValidationErrors

	for i := 0; i < refType.NumField(); i++ {
		field := refValue.Type().Field(i)
		fieldValue := refValue.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			err := validStruct(fieldValue, rule)
			if err != nil {
				if !isErrValid(err) {
					return err
				}
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validStruct(field reflect.Value, tag string) error {
	parts := strings.Split(tag, ":")
	if len(parts) != 2 {
		return nil
	}
	ruleType, ruleVal := parts[0], parts[1]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		return validString(field, ruleType, ruleVal)
	case reflect.Int:
		return validInt(field, ruleType, ruleVal)
	case reflect.Slice:
		return validSlice(field, tag)
	default:
		return nil
	}
}

func validString(field reflect.Value, ruleType, ruleVal string) error {
	value := field.String()
	switch ruleType {
	case "len":
		expectedLen, err := strconv.Atoi(ruleVal)
		if err != nil {
			return err
		}

		if len(value) != expectedLen {
			return ErrValidLength
		}
	case "regexp":
		regex, err := regexp.Compile(ruleVal)
		if err != nil {
			return err
		}

		if !regex.MatchString(value) {
			return ErrValidRegExp
		}
	case "in":
		vals := strings.Split(ruleVal, ",")
		for _, val := range vals {
			if value == val {
				return nil
			}
		}

		return ErrValidStrNotIn
	}

	return nil
}

func validInt(field reflect.Value, ruleType, ruleVal string) error {
	value := int(field.Int())

	switch ruleType {
	case "min":
		minVal, err := strconv.Atoi(ruleVal)
		if err != nil {
			return err
		}

		if value < minVal {
			return ErrValidMin
		}
	case "max":
		maxVal, err := strconv.Atoi(ruleVal)
		if err != nil {
			return err
		}

		if value > maxVal {
			return ErrValidMax
		}
	case "in":
		nums := strings.Split(ruleVal, ",")
		for _, num := range nums {
			num, err := strconv.Atoi(num)
			if err != nil {
				return err
			}

			if num == value {
				return nil
			}
		}

		return ErrValidIntNotIn
	}

	return nil
}

func validSlice(field reflect.Value, tag string) error {
	for i := 0; i < field.Len(); i++ {
		err := validStruct(field.Index(i), tag)
		if err != nil {
			return err
		}
	}

	return nil
}
