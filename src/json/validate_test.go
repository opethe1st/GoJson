package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	testcases := []TestCase{
		{"Unexpected end", []byte(`["k1","value"`), ValidationError{msg: "Was expecting ',' but we are at the end"}},
		{"Unexpected end of string", []byte(`"k1`), ValidationError{msg: "Was expecting '\"' but we are at the end"}},
		{"Unexpected end of array", []byte(`["k1",`), ValidationError{msg: "Was expecting ']' but we are at the end"}},
		{"Unexpected end of map", []byte(`{"k1":"v1"`), ValidationError{msg: "Was expecting ',' but we are at the end"}},
		{"empty array", []byte(`[]`), nil},
		{"empty dictionary", []byte(`{}`), nil},
		{"empty string", []byte(`""`), nil},
		{"standalone number", []byte(`12234`), nil},
		{"More character at the end", []byte(`"12234", 123`), ValidationError{msg: "Extra characters at the end of the json string"}},
		{"More spaces at the end", []byte(`"12234"  `), nil},
		{"true", []byte(`true`), nil},
		{"false", []byte(`false`), nil},
		{"null", []byte(`null`), nil},
		{"number with extras at the end", []byte(`1234tr`), ValidationError{msg: "Extra characters at the end of the json string"}},
		{"array with mixed types", []byte(`[1234, true]`), nil},
		{"dictionary with key type that is not string", []byte(`{1234: true}`), ValidationError{msg: "There is an error around {1̳234: true}. Key needs to be a valid string"}},
	}
	for _, testcase := range testcases {
		t.Run(
			testcase.name,
			func(t *testing.T) {
				assert.Equal(testcase.expectedOutput, Validate(testcase.input))
			},
		)
	}
}
