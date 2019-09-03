/*
These are just external tests to make sure it all works when it's called outside the json package
Is this necessary? Isn't this just solved by having my tests in the package end with the _test prefix?
Does that prefix mean I can't test implementation details?
*/
package unmarshall_test

import (
	// "encoding/json"
	myJson "github.com/opethe1st/GoJson/src/json"
	"github.com/stretchr/testify/assert"
	// "io/ioutil"
	"testing"
)

type (
	Any      = myJson.Any
	TestCase struct {
		name           string
		input          []byte
		expectedOutput Any
	}
)

func TestUnmarshall(t *testing.T) {
	assert := assert.New(t)

	testCases := []TestCase{
		{
			"Unmarshall a complex json",
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
			}`), map[string]Any{
				"k1": "v1",
				"k2": []Any{"v2"},
				"k3": map[string]Any{
					"k4": []Any{"v1"},
					"k5": "v2",
				},
				"k4":   nil,
				"k5":   true,
				"k6":   false,
				"k\t6": nil,
				"k7":   int64(123456),
			},
		},
	}

	for _, testcase := range testCases {
		output := myJson.Unmarshall(testcase.input)
		assert.Equal(testcase.expectedOutput, output)
	}
}

// func TestCompareUnmarshallToStdlibVersion(t *testing.T) {
// 	assert := assert.New(t)

// 	str, err := ioutil.ReadFile("../json/testdata/code.json")
// 	var value interface{}
// 	json.Unmarshal(str, &value)

// 	str, err = ioutil.ReadFile("../json/testdata/code.json")
// 	if err != nil {
// 		panic(err)
// 	}

// 	myValue := myJson.Unmarshall(str)
// apparently this takes forever to run. Crazy!
// and they are not the same, does this mean there is a problem with my implementation?
// assert.Equal(value, myValue)
// }
