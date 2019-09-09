package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCase struct {
	name           string
	input          []byte
	expectedOutput Any
}

func TestUnmarshallLiteral(t *testing.T) {
	assert := assert.New(t)
	testcases := []struct {
		name    string
		input   []byte
		literal string
		value   Any
	}{
		{"Null", []byte(`null`), `null`, nil},
		{"True", []byte(`true`), `true`, true},
		{"False", []byte(`false`), `false`, false},
	}
	for _, testcase := range testcases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				iter := &iterator{s: testcase.input}
				output := unmarshallLiteral(iter, testcase.literal, testcase.value)
				assert.Equal(testcase.value, output, "Expected UnmarshallKeyword(%v, %v, %v) to be %v but got %v", iter, testcase.literal, testcase.value, testcase.value, output)
			},
		)
	}

}

func TestUnmarshallString(t *testing.T) {
	assert := assert.New(t) // redefinition here -- ugly!
	testCases := []TestCase{
		{"Empty String", []byte(`""`), ""},
		{"Simple String", []byte(`"Key"`), "Key"},
		{"String with space at the beginning", []byte(`"   Key"`), "   Key"},
		{"String with space after it", []byte(`"Key"   `), "Key"},
		{"String with escaped character: tab", []byte(`"abc\t123"   `), "abc\t123"},
		{"String with escaped character: newline", []byte(`"abc\n123"`), "abc\n123"},
		{"String with escaped character: quote", []byte(`"\""`), "\""},
		{"String with more than one escaped character", []byte(`"she said \"a\""`), "she said \"a\""},
		{"String with only backslash", []byte(`"\\"`), "\\"},
		{"String with backslash", []byte(`"abc\\123"`), "abc\\123"},
		{"String with escaped unicode", []byte(`"\u1234"`), "ሴ"},
		{"String with just space", []byte(`"        "`), "        "},
	}
	for _, testcase := range testCases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				iter := &iterator{s: testcase.input}
				output := unmarshallString(iter)
				assert.Equal(testcase.expectedOutput, output, "Expected UnmarshallNumber(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
			},
		)

	}

}

func TestUnmarshallNumber(t *testing.T) {
	testCases := []TestCase{
		{"Float with 0 decimal", []byte(`123.0`), 123.0},
		{"Float with + sign", []byte(`+123.0`), 123.0},
		{"Float with - sign", []byte(`-123.0`), -123.0},
		{"Float with decimals", []byte(`-123.123`), -123.123},
		{"Float that's a decimal fraction", []byte(`0.234`), 0.234},
		{"Float with exponent", []byte(`1.234e2`), 123.4},
		{"Negative Float with exponent", []byte(`-1.234e2`), -123.4},
		{"Float with negative exponent", []byte(`-0.234e-2`), -0.00234},
	}
	for _, testcase := range testCases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				iter := &iterator{s: testcase.input}
				output := unmarshallNumber(iter)
				if !floatEquals(output.(float64), testcase.expectedOutput.(float64)) {
					t.Errorf("Expected UnmarshallNumber(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
				}
			},
		)
	}

	assert := assert.New(t)

	testCases = []TestCase{
		{"", []byte(`123`), int64(123)},
		{"", []byte(`-123`), int64(-123)},
	}
	for _, testcase := range testCases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				iter := &iterator{s: testcase.input}
				output := unmarshallNumber(iter)
				assert.Equal(testcase.expectedOutput, output)
			},
		)
	}
}

func TestUnmarshallArray(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{"Empty Array", []byte(`[]`), make([]Any, 0)},
		{"Array with a single value", []byte(`["value"]`), []Any{"value"}},
		{"Array with mor than one value", []byte(`["v1", "v2", "v3"]`), []Any{"v1", "v2", "v3"}},
		{"Nested array of depth 2", []byte(`["v1", ["v2", "v3"]]`), []Any{"v1", []Any{"v2", "v3"}}},
		{"Nested array of depth 3", []byte(`["v1", ["v2", ["v3"]]]`), []Any{"v1", []Any{"v2", []Any{"v3"}}}},
		{"Array that has an object", []byte(`["v1", {"v2": "v3"}]`), []Any{"v1", map[string]Any{"v2": "v3"}}},
	}
	for _, testcase := range testCases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				iter := &iterator{s: testcase.input}
				output := unmarshallArray(iter)
				assert.Equal(testcase.expectedOutput, output, "Expected UnmarshallArray(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
			},
		)
	}
}

func TestUnmarshallObject(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{"Empty object", []byte(`{}`), make(map[string]Any, 0)},
		{"Object with one item", []byte(`{"key": "value"}`), map[string]Any{"key": "value"}},
		{"Object with two items", []byte(`{"k1": "v1", "k2":"v2"}`), map[string]Any{"k1": "v1", "k2": "v2"}},
		{"Object with array value", []byte(`{"v1": ["v2", "v3"]}`), map[string]Any{"v1": []Any{"v2", "v3"}}},
	}
	for _, testcase := range testCases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				iter := &iterator{s: testcase.input}
				output := unmarshallObject(iter)
				assert.Equal(testcase.expectedOutput, output, "Expected UnmarshallObject(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
			},
		)
	}
}
