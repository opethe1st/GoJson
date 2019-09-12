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
func Validate(s []byte) error {
	iter := iterator{s: s}
	err := validate(&iter)
	if err != nil {
		return err
	}
	iter.AdvancePastAllWhiteSpace()
	if iter.Cursor() != iter.Len() {
		return ValidationError{msg: "Extra characters at the end of the json string"}
	}
	return nil
}

func validate(iter *iterator) error {
	iter.AdvancePastAllWhiteSpace()
	switch {
	case iter.Current() == '[':
		return validateArray(iter)
	case iter.Current() == '"':
		return validateString(iter)
	case iter.Current() == '{':
		return validateObject(iter)
	case iter.Current() == 'n':
		return validateLiteral(iter, "null")
	case iter.Current() == 't':
		return validateLiteral(iter, "true")
	case iter.Current() == 'f':
		return validateLiteral(iter, "false")
	case isNumber(iter):
		return validateNumber(iter)
	default:
		return ValidationError{msg: fmt.Sprintf("Unknown value at %d", iter.Cursor())}
	}
}

func validateLiteral(iter *iterator, literal string) error {
	for _, char := range literal {
		if rune(iter.Current()) != char {
			return ValidationError{msg: fmt.Sprintf("Error when trying to unmarshall '%v'", literal)}
		}
		iter.Next()
	}
	return nil
}

func validateNumber(iter *iterator) error {
	// maybe this should be an explicit state machine
	hasSign := false
	if (iter.Current() == '-') || (iter.Current() == '+') {
		iter.Next()
		hasSign = true
	}
	// there needs to be a digit after - or +
	if hasSign && !unicode.IsDigit(rune(iter.Current())) {
		return ValidationError{msg: "There needs to be a digit after - or +"}
	}
	for unicode.IsDigit(rune(iter.Current())) {
		iter.Next()
	}
	hasDot := false
	if iter.Current() == '.' {
		iter.Next()
		hasDot = true
	}
	if hasDot && !unicode.IsDigit(rune(iter.Current())) {
		return ValidationError{msg: "There needs to be a digit after . "}
	}
	for unicode.IsDigit(rune(iter.Current())) {
		iter.Next()
	}

	hasExponent := false
	if (iter.Current() == 'e') || (iter.Current() == 'E') {
		hasExponent = true
		iter.Next()
	}
	if (iter.Current() == '-') || (iter.Current() == '+') {
		iter.Next()
	}
	beforeExponentNumber := iter.Current()
	for unicode.IsDigit(rune(iter.Current())) {
		iter.Next()
	}
	// if we have encountered e/E then make sure there is at least one digit after e/E
	if hasExponent && (iter.Current()-beforeExponentNumber) == 0 {
		return ValidationError{msg: "There needs to be at least one digit after e/E when parsing a number"}
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

	// if empty string
	if iter.Current() == '"' {
		err = iter.AdvancePast('"')
		if err != nil {
			return err
		}
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
	}
	err = iter.AdvancePast(']')
	if err != nil {
		return err
	}
	return nil
}

func validateObject(iter *iterator) error {
	var err error

	err = iter.AdvancePast('{')
	if err != nil {
		return err
	}
	if iter.Current() == '}' {
		iter.Next()
		return nil
	}
	for iter.HasNext() {
		// key needs to be a string
		iter.AdvancePastAllWhiteSpace()
		err = validateString(iter)
		if err != nil {
			return ValidationError{msg: errorMsg(iter, "Key needs to be a valid string")}
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
	// this should return err
	err = iter.AdvancePast('}')
	if err != nil {
		return err
	}
	return nil
}
