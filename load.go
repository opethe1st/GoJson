/*
Package json implements a subset of the json specification defined at https://www.json.org.
Json is a very simple data format that supports serializing primitive types like string and numbers as well as composite ones like
sequences and mappings. You can check https://www.json.org for more details.
This package doesn't implement the entire specification. The only types it supports are string (strings without escaped characters),
sequences (aka arrays) and mapping (aka dictionaries)
*/
package json

import (
	"fmt"
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
	case iter.isEnd():
		return nil
	case iter.getCurrent() == 'n':
		return loadKeyword(iter, "null", nil)
	case iter.getCurrent() == 't':
		return loadKeyword(iter, "true", true)
	case iter.getCurrent() == 'f':
		return loadKeyword(iter, "false", false)
	case isNumber(iter):
		return loadNumber(iter)
	case iter.getCurrent() == '"':
		return loadString(iter)
	case iter.getCurrent() == '[':
		return loadSequence(iter)
	case iter.getCurrent() == '{':
		return loadMapping(iter)
	default:
		end := iter.offset + 100
		if len(iter.s) < end {
			end = len(iter.s)
		}
		panic(fmt.Sprintf("There is an error around\n\n%s", iter.s[iter.offset:end]))
	}
}

func loadKeyword(iter *iterator, keyword string, value interface{}) interface{} {
	for _, val := range keyword {
		if rune(iter.getCurrent()) != val {
			panic(fmt.Sprintf("There was an error while reading in %s", keyword))
		}
		iter.advance()
	}
	return value
}

func isNumber(iter *iterator) bool {
	switch iter.getCurrent() {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

func loadNumber(iter *iterator) interface{} {
	num := 0
	for !iter.isEnd() && unicode.IsDigit(rune(iter.getCurrent())) {
		num *= 10
		val, _ := strconv.ParseInt(string(iter.getCurrent()), 10, 64)
		num += int(val)
		iter.advance()
	}
	return num
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
	for !iter.isEnd() && (iter.getCurrent() != '"') {
		if iter.getCurrent() == '\\' {
			iter.advance()
			current := iter.getCurrent()
			switch current {
			case '"', '\\', 'b', 'f', 'n', 'r', 't':
				s = append(s, mapping[rune(current)])
				//need to handle the default case and handle u and hex digits
			case 'u':
				var ans rune
				// I should make sure these are valid hex digits btw, but will leave it for error reporting
				for i := 0; i < 4; i++ {
					iter.advance() // move past the 'u'
					fmt.Println(i, ans, string(iter.getCurrent()))
					ans = ans * 16
					ans += convertToDecimal[rune(iter.getCurrent())]
				}
				s = append(s, ans)
			default:
				s = append(s, rune('\\'))
				s = append(s, rune(iter.getCurrent()))
			}
		} else {
			s = append(s, rune(iter.getCurrent()))
		}
		iter.advance()
	}
	consume(iter, '"')
	return string(s)
}

func loadSequence(iter *iterator) []interface{} {
	seq := make([]interface{}, 0)
	consume(iter, '[')
	var item interface{}
	for !iter.isEnd() && (iter.getCurrent() != ']') {
		item = load(iter)
		seq = append(seq, item)
		consumeWhiteSpace(iter)
		if iter.getCurrent() == ']' {
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
	for !iter.isEnd() && (iter.s[iter.offset] != '}') {
		key = load(iter)
		consumeWhiteSpace(iter)
		consume(iter, ':')
		consumeWhiteSpace(iter)
		value = load(iter)
		mapping[key.(string)] = value
		if iter.getCurrent() == '}' {
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
	for !iter.isEnd() && unicode.IsSpace(rune(iter.getCurrent())) {
		iter.advance()
	}
}

func consume(iter *iterator, char byte) {
	// actually should probably raise an error if char isn't consumed
	if iter.getCurrent() == char {
		iter.advance()
	}
}
