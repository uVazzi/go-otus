package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrConversion    = errors.New("error conversion")
)

func Unpack(inputText string) (string, error) {
	var result strings.Builder

	runes := []rune(inputText)
	if len(runes) > 0 && unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	for i := 0; i < len(runes); i++ {
		isNext := i != len(runes)-1
		if isNext && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+1]) {
			return "", ErrInvalidString
		}
		// write result
		if isNext && unicode.IsDigit(runes[i+1]) {
			nextParseInt, parseErr := strconv.Atoi(string(runes[i+1]))
			if parseErr != nil {
				return "", ErrConversion
			}
			result.WriteString(strings.Repeat(string(runes[i]), nextParseInt))
			continue
		}
		if !unicode.IsDigit(runes[i]) {
			result.WriteRune(runes[i])
		}
	}

	return result.String(), nil
}
