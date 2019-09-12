package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	testcases := []TestCase{

		{"Unexpected end of string", []byte(`"k1`), ValidationError{msg: "Was expecting '\"' but we are at the end"}},
		{"empty string", []byte(`""`), nil},
		{"single quote in string", []byte(`"'"`), nil},
		{"double quote in string", []byte(`"\""`), nil},
		{"slash in string", []byte(`"\\"`), nil},

		{"standalone number", []byte(`12234`), nil},
		{"number with extras at the end", []byte(`1234tr`), ValidationError{msg: "Extra characters at the end of the json string"}},
		{"number with exponent", []byte(`1234e123`), nil},
		{"number with exponent with eE", []byte(`1234eE123`), ValidationError{msg: "There needs to be at least one digit after e/E when parsing a number"}},
		{"number with exponent without a digit", []byte(`1234e`), ValidationError{msg: "There needs to be at least one digit after e/E when parsing a number"}},
		{"fraction with exponent", []byte(`0.12e123`), nil},
		{"fraction with positive exponent", []byte(`0.12e+123`), nil},
		{"fraction with negative exponent", []byte(`0.12e-123`), nil},
		{". without number after", []byte(`0.`), ValidationError{msg: "There needs to be a digit after . "}},
		{"- on its own", []byte(`-`), ValidationError{msg: "There needs to be a digit after - or +"}},
		{"+ on its own", []byte(`+`), ValidationError{msg: "There needs to be a digit after - or +"}},
		{"zero with exponent", []byte(`0e10`), nil},
		{"one then fraction", []byte(`1.34`), nil},

		{"true", []byte(`true`), nil},
		{"false", []byte(`false`), nil},
		{"null", []byte(`null`), nil},

		{"Unexpected end of array", []byte(`["k1",`), ValidationError{msg: "Was expecting ']' but we are at the end"}},
		{"empty array", []byte(`[]`), nil},
		{"array with mixed types", []byte(`[1234, true]`), nil},

		{"empty object", []byte(`{}`), nil},
		{"Unexpected end of object", []byte(`{"k1":"v1"`), ValidationError{msg: "Was expecting ',' but we are at the end"}},
		{"object with key type that is not string", []byte(`{1234: true}`), ValidationError{msg: "There is an error around {1̳234: true}. Key needs to be a valid string"}},

		{"More character at the end", []byte(`"12234", 123`), ValidationError{msg: "Extra characters at the end of the json string"}},
		{"More spaces at the end", []byte(`"12234"  `), nil},
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
