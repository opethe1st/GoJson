/*
Package json implements the json specification.
Check https://www.json.org for more details.
*/
package json

import (
	"strconv"
	"unicode"
)

// Load is used load an object from a string - should I have called it Unmarshal to be consistent?
func Load(s []byte) interface{} {
	iter := &iterator{s: s}
	return load(iter)
}

func load(iter *iterator) interface{} {
	iter.AdvancePastAllWhiteSpace()
	switch {
	case iter.Current() == 'n':
		return loadKeyword(iter, "null", nil)
	case iter.Current() == 't':
		return loadKeyword(iter, "true", true)
	case iter.Current() == 'f':
		return loadKeyword(iter, "false", false)
	case isNumber(iter):
		return loadNumber(iter)
	case iter.Current() == '"':
		return loadString(iter)
	case iter.Current() == '[':
		return loadArray(iter)
	case iter.Current() == '{':
		return loadObject(iter)
	default:
		// should I be using panic at all? or just return the error as a value?
		panic(errorMsg(iter, "Cannot detect the value here"))
	}
}

func loadKeyword(iter *iterator, keyword string, value interface{}) interface{} {
	for _, val := range keyword {
		if rune(iter.Current()) != val {
			panic(errorMsg(iter, "There was an error while reading in %s", keyword))
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

func loadNumber(iter *iterator) interface{} {
	start := iter.Offset
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
		floatValue, err := strconv.ParseFloat(string(iter.SliceTillOffset(start)), 64)
		if err != nil {
			panic(errorMsg(iter, "This error %s occurred while trying to parse a number", err))
		} else {
			return floatValue
		}
	}

	intValue, err := strconv.ParseInt(string(iter.SliceTillOffset(start)), 10, 64)
	if err != nil {
		panic(errorMsg(iter, "This error %s occurred while trying to parse a number", err))
	}
	return intValue
}

func loadString(iter *iterator) string {
	var str string
	start := iter.Offset
	iter.AdvancePast('"')
	if iter.Current() == '"' {
		return str
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
	str, err := strconv.Unquote(string(iter.Slice(start, iter.Offset)))
	if err != nil {
		panic(errorMsg(iter, "There was an error unquoting this %s", string(iter.SliceTillOffset(start))))
	}
	return str
}

func loadArray(iter *iterator) []interface{} {
	array := make([]interface{}, 0)
	iter.AdvancePast('[')
	if iter.Current() == ']' {
		iter.AdvancePast(']')
		return array
	}
	var item interface{}
	for iter.HasNext() {
		item = load(iter)
		array = append(array, item)
		iter.AdvancePastAllWhiteSpace()
		if iter.Current() == ']' {
			break
		}
		iter.AdvancePast(',')
		iter.AdvancePastAllWhiteSpace()
	}
	iter.AdvancePast(']')
	return array
}

func loadObject(iter *iterator) map[string]interface{} {
	object := make(map[string]interface{}, 0)
	iter.AdvancePast('{')
	if iter.Current() == '}' {
		iter.AdvancePast('}')
		return object
	}
	var key, value interface{}
	for iter.HasNext() {
		key = load(iter)
		iter.AdvancePast(':')
		value = load(iter)

		object[key.(string)] = value

		iter.AdvancePastAllWhiteSpace()
		if iter.Current() == '}' {
			break
		}
		iter.AdvancePast(',')
	}
	iter.AdvancePastAllWhiteSpace()
	iter.AdvancePast('}')
	return object
}
