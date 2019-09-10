/*
Package json implements the json specification.
Check https://www.json.org for more details.
*/
package json

import (
	"strconv"
	"unicode"
)

// Any is an alias for interface{}
type Any = interface{}

// Unmarshall is used load an object from a string
func Unmarshall(s []byte) Any {
	return unmarshall(&iterator{s: s})
}

func unmarshall(iter *iterator) Any {
	iter.AdvancePastAllWhiteSpace()
	switch {
	case iter.Current() == 'n':
		return unmarshallLiteral(iter, "null", nil)
	case iter.Current() == 't':
		return unmarshallLiteral(iter, "true", true)
	case iter.Current() == 'f':
		return unmarshallLiteral(iter, "false", false)
	case isNumber(iter):
		return unmarshallNumber(iter)
	case iter.Current() == '"':
		return unmarshallString(iter)
	case iter.Current() == '[':
		return unmarshallArray(iter)
	case iter.Current() == '{':
		return unmarshallObject(iter)
	default:
		// should I be using panic at all? or just return the error as a value?
		panic(errorMsg(iter, "Cannot detect the value here"))
	}
}

func unmarshallLiteral(iter *iterator, literal string, value Any) Any {
	for _, val := range literal {
		if rune(iter.Current()) != val {
			panic(errorMsg(iter, "There was an error while reading in %s", literal))
		}
		iter.Next()
	}
	return value
}

func isNumber(iter *iterator) bool {
	switch iter.Current() {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '0':
		return true
	}
	return false
}

func unmarshallNumber(iter *iterator) Any {
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
		floatValue, err := strconv.ParseFloat(string(iter.SliceTillCursor(start)), 64)
		if err != nil {
			panic(errorMsg(iter, "This error %s occurred while trying to parse a number", err))
		} else {
			return floatValue
		}
	}

	intValue, err := strconv.ParseInt(string(iter.SliceTillCursor(start)), 10, 64)
	if err != nil {
		panic(errorMsg(iter, "This error %s occurred while trying to parse a number", err))
	}
	return intValue
}

func unmarshallString(iter *iterator) (str string) {
	start := iter.Cursor()
	iter.AdvancePast('"')
	if iter.Current() == '"' {
		return
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
	iter.AdvancePast('"')
	str, err := strconv.Unquote(string(iter.SliceTillCursor(start)))
	if err != nil {
		panic(errorMsg(iter, "There was an error unquoting this %s", string(iter.SliceTillCursor(start))))
	}
	return
}

func unmarshallArray(iter *iterator) []Any {
	var err error
	array := make([]Any, 0)
	err = iter.AdvancePast('[')
	if err != nil {
		panic(err)
	}

	if iter.Current() == ']' {
		iter.Next()
		return array
	}
	var item Any
	for iter.HasNext() {
		item = unmarshall(iter)
		array = append(array, item)
		iter.AdvancePastAllWhiteSpace()
		if iter.Current() == ']' {
			break
		}
		iter.AdvancePast(',')
	}
	iter.AdvancePast(']')
	return array
}

func unmarshallObject(iter *iterator) map[string]Any {
	var err error

	object := make(map[string]Any, 0)
	err = iter.AdvancePast('{')
	if err != nil {
		panic(err)
	}
	if iter.Current() == '}' {
		iter.Next()
		return object
	}
	var (
		key   string
		value Any
	)
	for iter.HasNext() {
		iter.AdvancePastAllWhiteSpace()
		key = unmarshallString(iter)
		err = iter.AdvancePast(':')
		if err != nil {
			panic(err)
		}
		value = unmarshall(iter)

		object[key] = value

		iter.AdvancePastAllWhiteSpace()
		if iter.Current() == '}' {
			break
		}
		iter.AdvancePast(',')
		if err != nil {
			panic(err)
		}
	}
	iter.AdvancePastAllWhiteSpace()
	err = iter.AdvancePast('}')
	if err != nil {
		panic(err)
	}
	return object
}
