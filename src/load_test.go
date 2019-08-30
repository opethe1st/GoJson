package json

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	input          []byte
	expectedOutput interface{}
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)

	testCases := []TestCase{
		{
			[]byte(`{
				"k1": "v1",
				"k2": [
					"v2"
				],
				"k3": {
					"k4": ["v1"],
					"k5": "v2"
				},
				"k4": null,
				"k5": true,
				"k6": false,
				"k	6": null,
				"k7": 123456
			}`), map[string]interface{}{
				"k1": "v1",
				"k2": []interface{}{"v2"},
				"k3": map[string]interface{}{
					"k4": []interface{}{"v1"},
					"k5": "v2",
				},
				"k4":   nil,
				"k5":   true,
				"k6":   false,
				"k\t6": nil,
				"k7":   123456.0,
			},
		},
	}

	for _, testcase := range testCases {
		output := Load(testcase.input)
		assert.Equal(testcase.expectedOutput, output)
	}
}

func TestLoadKeyword(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		input   []byte
		keyword string
		value   interface{}
	}{
		{[]byte(`null`), `null`, nil},
		{[]byte(`true`), `true`, true},
		{[]byte(`false`), `false`, false},
	}
	for _, testcase := range testCases {
		iter := &iterator{s: testcase.input}
		output := loadKeyword(iter, testcase.keyword, testcase.value)
		assert.Equal(testcase.value, output, "Expected loadKeyword(%v, %v, %v) to be %v but got %v", iter, testcase.keyword, testcase.value, testcase.value, output)
	}

}

func TestLoadString(t *testing.T) {
	assert := assert.New(t) // redefinition here -- ugly!
	testCases := []TestCase{
		{[]byte(`"Key"`), "Key"},
		{[]byte(`"   Key"`), "   Key"},
		{[]byte(`"Key"`), "Key"},
		{[]byte(`"Key"   `), "Key"},
		{[]byte(`"abc\t123"   `), "abc\t123"},
		{[]byte(`"abc\n123"`), "abc\n123"},
		// {[]byte(`"she said \"a\""`), "she said \"a\""},
		{[]byte(`"\\"`), "\\"},
		// {[]byte(`"abc\123"`), "abc\\123"},
		{[]byte(`"\u1234"`), "ሴ"},
	}
	for _, testcase := range testCases {
		iter := &iterator{s: testcase.input}
		output := loadString(iter)
		assert.Equal(testcase.expectedOutput, output, "Expected loadNumber(%v) to be %v but got %v", iter, testcase.expectedOutput, output)

	}

}

func TestLoadNumber(t *testing.T) {
	testCases := []TestCase{
		{[]byte(`123`), 123.0},
		{[]byte(`-123`), -123.0},
		{[]byte(`-123.123`), -123.123},
		{[]byte(`0.234`), 0.234},
		{[]byte(`1.234e2`), 123.4},
		{[]byte(`-1.234e2`), -123.4},
		{[]byte(`-0.234e2`), -23.4},
	}
	for _, testcase := range testCases {
		iter := &iterator{s: testcase.input}
		output := loadNumber(iter)
		if !floatEquals(output.(float64), testcase.expectedOutput.(float64)) {
			t.Errorf("Expected loadNumber(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
		}
	}
}

func floatEquals(a, b float64) bool {
	if math.Abs(a-b) < 0.00000001 {
		return true
	}
	return false
}

func TestLoadSequence(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{[]byte(`[]`), make([]interface{}, 0)},
		{[]byte(`["value"]`), []interface{}{"value"}},
		{[]byte(`["v1", "v2", "v3"]`), []interface{}{"v1", "v2", "v3"}},
		{[]byte(`["v1", ["v2", "v3"]]`), []interface{}{"v1", []interface{}{"v2", "v3"}}},
		{[]byte(`["v1", ["v2", ["v3"]]]`), []interface{}{"v1", []interface{}{"v2", []interface{}{"v3"}}}},
		{[]byte(`["v1", {"v2": "v3"}]`), []interface{}{"v1", map[string]interface{}{"v2": "v3"}}},
	}
	for _, testcase := range testCases {
		iter := &iterator{s: testcase.input}
		output := loadSequence(iter)
		assert.Equal(testcase.expectedOutput, output, "Expected loadSequence(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
	}
}

func TestLoadMapping(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{[]byte(`{}`), make(map[string]interface{}, 0)},
		{[]byte(`{"key": "value"}`), map[string]interface{}{"key": "value"}},
		{[]byte(`{"k1": "v1", "k2":"v2"}`), map[string]interface{}{"k1": "v1", "k2": "v2"}},
		{[]byte(`{"v1": ["v2", "v3"]}`), map[string]interface{}{"v1": []interface{}{"v2", "v3"}}},
	}
	for _, testcase := range testCases {
		iter := &iterator{s: testcase.input}
		output := loadMapping(iter)
		assert.Equal(testcase.expectedOutput, output, "Expected loadMapping(%v) to be %v but got %v", iter, testcase.expectedOutput, output)
	}
}

func TestConsumeSpaces(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{[]byte(`"key"`), 0},
		{[]byte(` "key" `), 1},
		{[]byte(`   "key"`), 3},
	}
	for _, testcase := range testCases {
		iter := &iterator{s: testcase.input}
		iter.AdvancePassWhiteSpace()
		assert.Equal(testcase.expectedOutput, iter.Offset, "Expected consumeWhiteSpace(%v) to be %d but got %d", iter, testcase.expectedOutput, iter.Offset)
	}
}
