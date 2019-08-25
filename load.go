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
	"unicode"
)

// Load is used load an object from a string
func Load(s string) interface{} {
	_, value := load(s, 0)
	return value
}

func load(s string, current int) (int, interface{}) {
	current = consumeWhiteSpace(s, current)
	switch {
	case isString(s, current):
		return loadString(s, current)
	case isSequence(s, current):
		return loadSequence(s, current)
	case isMapping(s, current):
		return loadMapping(s, current)
	default:
		end := current + 100
		if len(s) < end {
			end = len(s)
		}
		panic(fmt.Sprintf("There is an error around\n\n%s", s[current:end]))
	}
}

// strings

func isString(s string, current int) bool {
	return s[current] == '"'
}

func loadString(s string, current int) (int, interface{}) {
	// actually should probably raise an error if '"' isn't consumed
	start := consume(s, current, '"')
	current = start
	for current < len(s) && s[current] != '"' {
		current++
	}
	// and current + 1 since the next not visited character in s is at current + 1
	return current + 1, s[start:current]
}

// sequences

func isSequence(s string, current int) bool {
	return s[current] == '['
}

func loadSequence(s string, current int) (int, []interface{}) {
	seq := make([]interface{}, 0)
	// actually should probably raise an error if '[' isn't consumed
	current = consume(s, current, '[')
	var item interface{}
	for (current < len(s)) && (s[current] != ']') {
		current, item = load(s, current)
		current = consumeWhiteSpace(s, current)
		// technically not allowed to have ["key",] but it is currently allowed
		current = consume(s, current, ',')
		current = consumeWhiteSpace(s, current)
		seq = append(seq, item)
		current = consumeWhiteSpace(s, current)
	}
	return current + 1, seq
}

// mappings

func isMapping(s string, current int) bool {
	return s[current] == '{'
}

func loadMapping(s string, current int) (int, map[string]interface{}) {
	mapping := make(map[string]interface{}, 0)
	// actually should probably raise an error if '{' isn't consumed
	current = consume(s, current, '{')
	var key, value interface{}
	for (current < len(s)) && (s[current] != '}') {
		current, key = load(s, current)
		current = consumeWhiteSpace(s, current)
		current = consume(s, current, ':')
		current = consumeWhiteSpace(s, current)
		current, value = load(s, current)
		mapping[key.(string)] = value
		// technically not allowed to have {"key":"value",} but it is currently allowed
		current = consume(s, current, ',')
		current = consumeWhiteSpace(s, current)
	}
	return current + 1, mapping
}

// utils - this could be in a separate file

// consumeWhiteSpace returns the next index after current such that s[index] is not whitespace
func consumeWhiteSpace(s string, current int) int {
	for current < len(s) && unicode.IsSpace(rune(s[current])) {
		current++
	}
	return current
}

func consume(s string, current int, char byte) int {
	if s[current] == char {
		return current + 1
	}
	return current
}
