/*
Package json implements the json specification.
Check https://www.json.org for more details.
*/
package json

import (
	"fmt"
	"math"
	"strconv"
	"unicode"
)

// Load is used load an object from a string
func Load(s string) interface{} {
	iter := &iterator{s: s}
	return load(iter)
}

func load(iter *iterator) interface{} {
	consumeWhiteSpace(iter)
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
		return loadSequence(iter)
	case iter.Current() == '{':
		return loadMapping(iter)
	default:
		panic(errorMessage(iter))
	}
}

func loadKeyword(iter *iterator, keyword string, value interface{}) interface{} {
	for _, val := range keyword {
		if rune(iter.Current()) != val {
			panic(fmt.Sprintf("There was an error while reading in %s", keyword))
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
	//negative numbers
	sign := 1.0
	if iter.Current() == '-' {
		sign = -1.0
		iter.Next()
	}
	num := 0.0
	for iter.HasNext() && unicode.IsDigit(rune(iter.Current())) {
		num *= 10
		val, _ := strconv.ParseInt(string(iter.Current()), 10, 64)
		num += float64(val)
		iter.Next()
	}

	// decimal
	// some of the code here is a duplicate of what is above, I should consolidate into one function.
	consume(iter, '.')
	frac := 0.0
	power := 0.1
	for iter.HasNext() && unicode.IsDigit(rune(iter.Current())) {
		val, _ := strconv.ParseInt(string(iter.Current()), 10, 64)
		frac += power * float64(val)
		power *= 0.1
		iter.Next()
	}

	exponent := 0.0
	//exponent
	if iter.HasNext() && ((iter.Current() == 'e') || (iter.Current() == 'E')) {
		if iter.Current() == 'e' {
			consume(iter, 'e')
		}
		if iter.Current() == 'E' {
			consume(iter, 'E')
		}
		exponentSign := 1.0
		//TODO(ope) this is subtly wrong since it allows +-123234, I will fix later
		if iter.HasNext() && iter.Current() == '+' {
			consume(iter, '+')
		}
		if iter.HasNext() && iter.Current() == '-' {
			consume(iter, '-')
			exponentSign = -1.0
		}
		// there needs to be at least one digit after an exponent
		exponent = 0.0
		for iter.HasNext() && unicode.IsDigit(rune(iter.Current())) {
			exponent *= 10
			val, _ := strconv.ParseInt(string(iter.Current()), 10, 64)
			exponent += float64(val)
			iter.Next()
		}
		exponent *= exponentSign
	}
	return (num + frac) * sign * math.Pow(10, exponent)
}

func loadString(iter *iterator) string {
	consume(iter, '"')
	s := make([]rune, 0)
	mapping := map[rune]rune{
		'"':  '"',
		'\\': '\\',
		'b':  '\b',
		'f':  '\f',
		'n':  '\n',
		'r':  '\r',
		't':  '\t',
	}
	// TODO(better as a function?)
	convertToDecimal := map[rune]rune{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'a': 10,
		'A': 10,
		'b': 11,
		'B': 11,
		'c': 12,
		'C': 12,
		'd': 13,
		'D': 13,
		'e': 14,
		'E': 14,
		'f': 15,
		'F': 15,
	}
	for iter.HasNext() && (iter.Current() != '"') {
		if iter.Current() == '\\' {
			iter.Next()
			current := iter.Current()
			switch current {
			case '"', '\\', 'b', 'f', 'n', 'r', 't':
				s = append(s, mapping[rune(current)])
				//need to handle the default case and handle u and hex digits
			case 'u':
				var ans rune
				// I should make sure these are valid hex digits btw, but will leave it for error reporting
				for i := 0; i < 4; i++ {
					iter.Next() // move past the 'u'
					ans = ans * 16
					ans += convertToDecimal[rune(iter.Current())]
				}
				s = append(s, ans)
			default:
				s = append(s, rune('\\'))
				s = append(s, rune(iter.Current()))
			}
		} else {
			s = append(s, rune(iter.Current()))
		}
		iter.Next()
	}
	consume(iter, '"')
	return string(s)
}

func loadSequence(iter *iterator) []interface{} {
	seq := make([]interface{}, 0)
	consume(iter, '[')
	var item interface{}
	for iter.HasNext() && (iter.Current() != ']') {
		item = load(iter)
		seq = append(seq, item)
		consumeWhiteSpace(iter)
		if iter.Current() == ']' {
			break
		}
		consume(iter, ',')
		consumeWhiteSpace(iter)
	}
	consume(iter, ']')
	return seq
}

func loadMapping(iter *iterator) map[string]interface{} {
	mapping := make(map[string]interface{}, 0)
	consume(iter, '{')
	var key, value interface{}
	for iter.HasNext() && (iter.Current() != '}') {
		key = load(iter)
		consumeWhiteSpace(iter)
		consume(iter, ':')
		consumeWhiteSpace(iter)
		value = load(iter)
		mapping[key.(string)] = value
		consumeWhiteSpace(iter)
		if iter.Current() == '}' {
			break
		}
		consume(iter, ',')
		consumeWhiteSpace(iter)
	}
	consume(iter, '}')
	return mapping
}

// utils - this could be in a separate file

func consumeWhiteSpace(iter *iterator) {
	for iter.HasNext() && unicode.IsSpace(rune(iter.Current())) {
		iter.Next()
	}
}

func consume(iter *iterator, char byte) {
	// actually should probably raise an error if char isn't consumed
	if iter.HasNext() && iter.Current() == char {
		iter.Next()
	}
	// add this else part later - this is part of better error handling
	// else {
	// 	panic(fmt.Sprintf("Expected %q but got %q", char, iter.Current())+errorMessage(iter))
	// }
}

func errorMessage(iter *iterator) string {
	startBefore := iter.offset - 50
	if startBefore < 0 {
		startBefore = 0
	}
	endBefore := iter.offset
	if endBefore < 0 {
		endBefore = 0
	}
	before := iter.s[startBefore:endBefore]

	startAfter := iter.offset + 1
	if startAfter > len(iter.s) {
		startAfter = len(iter.s)
	}
	endAfter := iter.offset + 50
	if endAfter > len(iter.s) {
		endAfter = len(iter.s)
	}
	after := iter.s[startAfter:endAfter]
	// this doesnt work well if the character to underline is a whitespace
	return fmt.Sprintf(`There is an error around
	%s%s%s
	`, before, string([]byte{iter.Current(), 204, 179}), after)
}
