//nolint:exhaustive
package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("error: not structure received")

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
		errmsg.WriteString(e.Err.Error() + ", ")
	}
	return errmsg.String()
}

func Validate(v interface{}) error {
	valueType := reflect.TypeOf(v)
	if valueType.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	value := reflect.ValueOf(v)

	errList := make(ValidationErrors, 0)

	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		if _, ok := field.Tag.Lookup("validate"); !ok {
			continue
		}

		rules := strings.Split(field.Tag.Get("validate"), "|")

		switch field.Type.Kind() {
		case reflect.Int:
			ruleInt(fieldValue.Int(), field.Name, rules, &errList)
		case reflect.String:
			ruleString(fieldValue.String(), field.Name, rules, &errList)
		case reflect.Slice:
			ruleSlice(fieldValue.Interface(), field.Name, rules, &errList)
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
func ruleInt(fieldValue int64, fieldName string, rules []string, errList *ValidationErrors) {
	for _, rule := range rules {
		if !strings.Contains(rule, ":") {
			*errList = append(*errList, NewValidationError(
				fieldName,
				errors.New("rule not contain : "),
			))
			return
		}
		i := strings.Index(rule, ":")
		switch rule[:i] {
		case "min":
			n, e := strconv.Atoi(rule[i+1:])
			if e != nil {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("the number in 'min' rule is not a number"),
				))
			}
			if fieldValue < int64(n) {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("less then "+rule[i+1:]),
				))
			}
		case "max":
			n, e := strconv.Atoi(rule[i+1:])
			if e != nil {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("the number in 'max' rule is not a number"),
				))
			}
			if fieldValue > int64(n) {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("greater then "+rule[i+1:]),
				))
			}
		case "in":
			n := strings.Split(rule[i+1:], ",")
			for y := range n {
				nn, e := strconv.Atoi(n[y])
				if e != nil {
					*errList = append(*errList, NewValidationError(
						fieldName,
						errors.New("the number in 'in' rule is not a number"),
					))
				}
				if fieldValue == int64(nn) {
					return
				}
			}
			*errList = append(*errList, NewValidationError(
				fieldName,
				errors.New("not in these mass of int "),
			))
		default:
			*errList = append(*errList, NewValidationError(
				fieldName,
				errors.New("wrong int rule"),
			))
		}
	}
}

// len:32 - длина строки должна быть ровно 32 символа;
// regexp:\\d+ - согласно регулярному выражению строка должна состоять из цифр (\\ - экранирование слэша);
// in:foo,bar - строка должна входить в множество строк {"foo", "bar"}.
func ruleString(fieldValue string, fieldName string, rules []string, errList *ValidationErrors) {
	for _, rule := range rules {
		if !strings.Contains(rule, ":") {
			return
		}
		i := strings.Index(rule, ":")
		switch rule[:i] {
		case "len":
			n, e := strconv.Atoi(rule[i+1:])
			if e != nil {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("the number in 'len' rule is not a number "),
				))
			}
			if len(fieldValue) != n {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("len is not equal "+rule[i+1:]),
				))
			}
		case "regexp":
			reg, _ := regexp.Compile(rule[i+1:])
			if !reg.MatchString(fieldValue) {
				*errList = append(*errList, NewValidationError(
					fieldName,
					errors.New("string is not matched regexp expression "),
				))
			}
		case "in":
			n := strings.Split(rule[i+1:], ",")
			for _, y := range n {
				if fieldValue == y {
					return
				}
			}
			*errList = append(*errList, NewValidationError(
				fieldName,
				errors.New("not in these mass of strings "),
			))

		default:
			*errList = append(*errList, NewValidationError(
				fieldName,
				errors.New("wrong string rule"),
			))
		}
	}
}

// Для слайсов валидируется каждый элемент слайса.
func ruleSlice(fieldValue interface{}, fieldName string, rules []string, errList *ValidationErrors) {
	switch sliceValues := fieldValue.(type) {
	case []int:
		for _, i := range sliceValues {
			ruleInt(int64(i), fieldName, rules, errList)
		}
	case []string:
		for _, i := range sliceValues {
			ruleString(i, fieldName, rules, errList)
		}
	default:
		*errList = append(*errList, NewValidationError(
			fieldName,
			errors.New("wrong slice rule"),
		))
	}
}
