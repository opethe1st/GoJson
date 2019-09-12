package json

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func Benchmark_Validate_NestedJson(b *testing.B) {
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/nested_array.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	var capture any
	for n := 0; n < b.N; n++ {
		capture = Validate(str)
	}
	res = capture
}

func Benchmark_Validate_NestedJsonStdlib(b *testing.B) {
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/nested_array.json")
	b.StartTimer()
	if err != nil {
		panic(err)
	}

	var capture any
	for n := 0; n < b.N; n++ {
		capture = json.Valid(str)
	}
	res = capture
}
