package json

import (
	"fmt"
	"strconv"
	"unicode"
)

// ValidationError type
type ValidationError struct {
	msg string
}

func (e ValidationError) Error() string {
	return e.msg
}

// Validate a json string
func Validate(s string) error {
	iter := iterator{s: []byte(s)}
	err := validate(&iter)
	if err != nil {
		return err
	}
	iter.AdvancePastAllWhiteSpace()
	if iter.Cursor() != iter.Len() {
		return ValidationError{msg: "Extraneous characters at the end"}
	}
	return nil
}

func validate(iter *iterator) error {
	iter.AdvancePastAllWhiteSpace()
	switch {
	case iter.Current() == 'n':
		return validateKeyword(iter, "null")
	case iter.Current() == 't':
		return validateKeyword(iter, "true")
	case iter.Current() == 'f':
		return validateKeyword(iter, "false")
	case isNumber(iter):
		return validateNumber(iter)
	case iter.Current() == '"':
		return validateString(iter)
	case iter.Current() == '[':
		return validateArray(iter)
	case iter.Current() == '{':
		return validateObject(iter)
	default:
		return ValidationError{msg: "Unknown value"}
	}
}

func validateKeyword(iter *iterator, literal string) error {
	for _, char := range literal {
		if rune(iter.Current()) != char {
			return ValidationError{msg: fmt.Sprintf("Error when trying to unmarshall '%v'", literal)}
		}
	}
	return nil
}

func validateNumber(iter *iterator) error {
	// move this block to a func so it can be reused across validation - perhaps have a step before either
	// that computes certain things before?
	start := iter.Cursor()
	isFloat := false
	if (iter.Current() == '-') || (iter.Current() == '+') {
		iter.Next()
	}
	for unicode.IsDigit(rune(iter.Current())) {
		iter.Next()
	}
	if iter.Current() == '.' {
		isFloat = true
		iter.Next()
	}
	for unicode.IsDigit(rune(iter.Current())) {
		iter.Next()
	}
	if (iter.Current() == 'e') || (iter.Current() == 'E') {
		isFloat = true
		iter.Next()
	}
	if (iter.Current() == '-') || (iter.Current() == '+') {
		iter.Next()
	}
	for unicode.IsDigit(rune(iter.Current())) {
		iter.Next()
	}
	if isFloat {
		_, err := strconv.ParseFloat(string(iter.SliceTillCursor(start)), 64)
		if err != nil {
			return err
		}
		return nil
	}

	_, err := strconv.ParseInt(string(iter.SliceTillCursor(start)), 10, 64)
	if err != nil {
		return err
	}
	return nil
}

func validateString(iter *iterator) error {
	var err error
	start := iter.Cursor()
	err = iter.AdvancePast('"')
	if err != nil {
		return err
	}
	if iter.Current() == '"' {
		return nil
	}
	for iter.HasNext() && iter.Current() != '"' {
		if iter.Current() == '\\' {
			iter.Next()
			// skip over an escaped `"` and `\`
			if (iter.Current() == '"') || (iter.Current() == '\\') {
				iter.Next()
			}
		} else {
			iter.Next()
		}
	}
	err = iter.AdvancePast('"')
	if err != nil {
		return err
	}
	_, err = strconv.Unquote(string(iter.SliceTillCursor(start)))
	if err != nil {
		return err
	}
	return nil
}

func validateArray(iter *iterator) error {
	var err error
	err = iter.AdvancePast('[')
	if err != nil {
		return err
	}
	if iter.Current() == ']' {
		iter.Next()
		return nil
	}
	for iter.HasNext() {
		err = validate(iter)
		if err != nil {
			return err
		}

		iter.AdvancePastAllWhiteSpace()
		if iter.Current() == ']' {
			break
		}
		err = iter.AdvancePast(',')
		if err != nil {
			return err
		}
		iter.AdvancePastAllWhiteSpace()
	}
	err = iter.AdvancePast(']')
	if err != nil {
		return err
	}
	return nil
}

func validateObject(iter *iterator) error {
	iter.AdvancePast('{')
	if iter.Current() == '}' {
		iter.Next()
		return nil
	}
	var err error
	for iter.HasNext() {
		err = validate(iter)
		if err != nil {
			return err
		}
		err = iter.AdvancePast(':')
		if err != nil {
			return err
		}
		err = validate(iter)
		if err != nil {
			return err
		}

		iter.AdvancePastAllWhiteSpace()
		if iter.Current() == '}' {
			break
		}
		err = iter.AdvancePast(',')
		if err != nil {
			return err
		}
	}
	iter.AdvancePastAllWhiteSpace()
	// this should return err
	err = iter.AdvancePast('}')
	if err != nil {
		return err
	}
	return nil
}
