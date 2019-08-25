package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCase struct {
	input          string
	expectedOutput interface{}
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)

	testCases := []TestCase{
		{`{
			"k1": "v1",
			"k2": [
				"v2"
			],
			"k3": {
				"k4": ["v1"],
				"k5": "v2"
			}
		}`, map[string]interface{}{
			"k1": "v1",
			"k2": []interface{}{"v2"},
			"k3": map[string]interface{}{
				"k4": []interface{}{"v1"},
				"k5": "v2",
			},
		},
		},
	}

	for _, testcase := range testCases {
		output := Load(testcase.input)
		assert.Equal(output, testcase.expectedOutput, "These should be equal")
	}
}

func TestIsString(t *testing.T) {
	testCases := []TestCase{
		{`This is a string`, false},
		{`"Key"`, true},
	}
	for _, testcase := range testCases {
		if output := isString(testcase.input, 0); output != testcase.expectedOutput {
			t.Errorf("Expected isString (%v) to be %t but got %t", testcase.input, testcase.expectedOutput, output)
		}
	}
}

func TestLoadString(t *testing.T) {
	testCases := []TestCase{
		{`"Key"`, "Key"},
		{`"   Key"`, "   Key"},
		{`"Key"`, "Key"},
		{`"Key"   `, "Key"},
	}
	for _, testcase := range testCases {
		if _, output := loadString(testcase.input, 0); output != testcase.expectedOutput {
			t.Errorf("Expected loadString (%v) to be %v but got %v", testcase.input, testcase.expectedOutput, output)
		}
	}

	testCases = []TestCase{
		{`"Key"`, 5},
		{`"   Key"`, 8},
		{`"Key"   `, 5},
	}
	for _, testcase := range testCases {
		if output, _ := loadString(testcase.input, 0); output != testcase.expectedOutput {
			t.Errorf("Expected loadString (%v) to be %v but got %v", testcase.input, testcase.expectedOutput, output)
		}
	}
}

func TestIsSequence(t *testing.T) {
	testCases := []TestCase{
		{`[]`, true},
		{`"Key"`, false},
	}
	for _, testcase := range testCases {
		if output := isSequence(testcase.input, 0); output != testcase.expectedOutput {
			t.Errorf("Expected isSequence (%v) to be %t but got %t", testcase.input, testcase.expectedOutput, output)
		}
	}
}

func TestLoadSequence(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{`[]`, make([]interface{}, 0)},
		{`["value"]`, []interface{}{"value"}},
		{`["v1", "v2", "v3"]`, []interface{}{"v1", "v2", "v3"}},
		{`["v1", ["v2", "v3"]]`, []interface{}{"v1", []interface{}{"v2", "v3"}}},
		{`["v1", ["v2", ["v3"]]]`, []interface{}{"v1", []interface{}{"v2", []interface{}{"v3"}}}},
		{`["v1", {"v2": "v3"}]`, []interface{}{"v1", map[string]interface{}{"v2": "v3"}}},
	}
	for _, testcase := range testCases {
		_, output := loadSequence(testcase.input, 0)
		assert.Equal(output, testcase.expectedOutput, "Expected loadSequence(%v) to be %v but got %v", testcase.input, testcase.expectedOutput, output)
	}
}

func TestIsMapping(t *testing.T) {
	testCases := []TestCase{
		{`[]`, false},
		{`{`, true},
	}
	for _, testcase := range testCases {
		if output := isMapping(testcase.input, 0); output != testcase.expectedOutput {
			t.Errorf("Expected isSequence (%v) to be %t but got %t", testcase.input, testcase.expectedOutput, output)
		}
	}
}

func TestLoadMapping(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{`{}`, make(map[string]interface{}, 0)},
		{`{"key": "value"}`, map[string]interface{}{"key": "value"}},
		{`{"k1": "v1", "k2":"v2"}`, map[string]interface{}{"k1": "v1", "k2": "v2"}},
		{`{"v1": ["v2", "v3"]}`, map[string]interface{}{"v1": []interface{}{"v2", "v3"}}},
	}
	for _, testcase := range testCases {
		_, output := loadMapping(testcase.input, 0)
		assert.Equal(output, testcase.expectedOutput, "Expected loadMapping(%v) to be %v but got %v", testcase.input, testcase.expectedOutput, output)
	}
}

func TestConsumeSpaces(t *testing.T) {
	testCases := []TestCase{
		{`"key"`, 0},
		{` "key" `, 1},
		{`   "key"`, 3},
	}
	for _, testcase := range testCases {
		if output := consumeWhiteSpace(testcase.input, 0); output != testcase.expectedOutput {
			t.Errorf("Expected consumeWhiteSpace (%v) to be %d but got %d", testcase.input, testcase.expectedOutput, output)
		}
	}
}
