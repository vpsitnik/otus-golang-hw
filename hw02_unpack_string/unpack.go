package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var unpackedStr strings.Builder
	rStr := []rune(str)

	for index, char := range rStr {

		currC := string(char)
		var nextC rune
		res := ""

		if index < len(str)-1 {
			nextC = rStr[index+1]
		}

		if unicode.IsDigit(char) {

			if index == 0 {
				return "", ErrInvalidString
			}

			prevC := rStr[index-1]

			if unicode.IsDigit(prevC) {
				return "", ErrInvalidString
			} else if prevC == 92 {
				continue
			}

			repeat, _ := strconv.Atoi(currC)

			res = strings.Repeat(string(prevC), repeat)

		} else if !unicode.IsDigit(nextC) {
			res = currC
		}

		unpackedStr.WriteString(res)

	}

	return unpackedStr.String(), nil
}
