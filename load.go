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

func loadString(iter *iterator) interface{} {
	consume(iter, '"')
	start := iter.offset
	for !iter.isEnd() && (iter.getCurrent() != '"') {
		iter.advance()
	}
	end := iter.offset
	consume(iter, '"')
	return iter.s[start:end]
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
