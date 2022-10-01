package hw09structvalidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNotStruct = errors.New("error: not structure received")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func NewValidationError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (v ValidationErrors) Error() string {
	var errmsg strings.Builder
	for _, e := range v {
		errmsg.WriteString("Field: " + e.Field + ", error: " + e.Err.Error() + "\n")
	}
	return errmsg.String()
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	valueType := reflect.TypeOf(v)

	if valueType.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	errList := make(ValidationErrors, 0)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		rules := strings.Split(field.Tag.Get("validate"), "|")

		switch field.Type.Kind() {
		case reflect.Int:
			errList = ruleInt(fieldValue.Int(), field.Name, rules)
		case reflect.String:
		case reflect.Slice:
		default:
			errList = append(errList, NewValidationError(
				field.Name,
				errors.New("field cannot be validated"),
			))
		}
	}
	return errList
}

// min:10 - число не может быть меньше 10;
// max:20 - число не может быть больше 20;
// in:256,1024 - число должно входить в множество чисел {256, 1024};
func ruleInt(fieldValue int64, fieldName string, rules []string) (errList ValidationErrors) {
	for _, rule := range rules {
		i := strings.Index(rule, ":")
		switch rule[:i] {
		case "min":
			n, e := strconv.Atoi(rule[i+1:])
			if e != nil {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("the number in 'min' rule is not a number"),
				))
			}
			if fieldValue < int64(n) {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("less then "+rule[i+1:]),
				))
			}
		case "max":
			n, e := strconv.Atoi(rule[i+1:])
			if e != nil {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("the number in 'max' rule is not a number"),
				))
			}
			if fieldValue > int64(n) {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("greater then "+rule[i+1:]),
				))
			}
		case "in":
			n := strings.Split(rule[i+1:], ",")
			n1, e := strconv.Atoi(n[0])
			if e != nil {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("the first number in 'in' rule is not a number"),
				))
			}
			n2, e := strconv.Atoi(n[1])
			if e != nil {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("the second number in 'in' rule is not a number"),
				))
			}
			if fieldValue < int64(n1) && fieldValue > int64(n2) {
				errList = append(errList, NewValidationError(
					fieldName,
					errors.New("not in "+rule[i+1:]),
				))
			}
		default:
		}
	}
	return errList
}
