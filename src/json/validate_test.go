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
		{"empty array", []byte(`[]`), nil},
		{"empty dictionary", []byte(`{}`), nil},
		{"empty string", []byte(`""`), nil},
		{"standalone number", []byte(`12234`), nil},
		{"More values at the end", []byte(`"12234", 123`), ValidationError{msg: "Extraneous values at the end of the string"}},
		{"true", []byte(`true`), nil},
		{"false", []byte(`false`), nil},
		{"null", []byte(`null`), nil},
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
