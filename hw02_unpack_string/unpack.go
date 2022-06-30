package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	r := []rune(s)
	var res strings.Builder
	for i := 0; i < len(r); i++ {
		switch {
		case r[i] == '\\': // Если обратный слеш
			i++
			res.WriteRune(r[i])

		case !unicode.IsDigit(r[i]):
			if i > 1 && r[i-1] == '\\' {
				return "", ErrInvalidString
			}
			res.WriteRune(r[i])

		case unicode.IsDigit(r[i]):
			// Если первый символ цифра или две цифры подряд и перед ними нет экранирования
			if i == 0 || unicode.IsDigit(r[i-1]) && r[i-2] != '\\' {
				return "", ErrInvalidString
			} else if !unicode.IsDigit(r[i-1]) || unicode.IsDigit(r[i-1]) && r[i-2] == '\\' {
				// Если цифра, а перед ней буква или две цифры подряд и перед ними есть экранирование
				n, err := strconv.Atoi(string(r[i]))
				if err != nil {
					return "", err
				}
				if n != 0 {
					res.WriteString(strings.Repeat(string(r[i-1]), n-1))
				} else {
					temp := res.String()
					res.Reset()
					res.WriteString(temp[:len(temp)-1])
				}
			}
		}
	}

	return res.String(), nil
}
