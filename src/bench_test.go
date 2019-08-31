package json

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func BenchmarkMyMapOfString(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/map_of_string.json")
	if err != nil {
		panic(err)
	}

	Load(str)
}

func BenchmarkMapOfString(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/map_of_string.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	json.Unmarshal(str, &data)
}

// these two are about the same
func BenchmarkMyArrayOfInt(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/array_of_int.json")
	if err != nil {
		panic(err)
	}

	Load(str)
}
func BenchmarkArrayOfInt(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/array_of_int.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	json.Unmarshal(str, &data)
}

// Array of Strigs
// There is a clear performance difference between these two - why?
// and mine is faster!
func BenchmarkMyArrayOfString(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/array_of_string.json")
	if err != nil {
		panic(err)
	}

	Load(str)
}
func BenchmarkArrayOfString(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/array_of_string.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	json.Unmarshal(str, &data)
}

// Crazy difference in performance and mine is faster
func BenchmarkMyBigString(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/big_string.json")
	if err != nil {
		panic(err)
	}

	Load(str)
}
func BenchmarkBigString(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/big_string.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	json.Unmarshal(str, &data)
}

// Interesting how my big gains in string processing evaporate here 🤔
func BenchmarkMyCodejson(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/code.json")
	if err != nil {
		panic(err)
	}

	Load(str)
}
func BenchmarkCodejson(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/code.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	json.Unmarshal(str, &data)
}

// similar performance
func BenchmarkMyNestedJson(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/nested_array.json")
	if err != nil {
		panic(err)
	}

	Load(str)
}
func BenchmarkNestedJson(b *testing.B) {
	b.ReportAllocs()
	str, err := ioutil.ReadFile("testdata/nested_array.json")
	if err != nil {
		panic(err)
	}
	var data interface{}
	json.Unmarshal(str, &data)
}
