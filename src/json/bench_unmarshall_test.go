package json

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

var res Any

func Benchmark_MapOfString(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/map_of_string.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		any = Unmarshall(str)
	}
	res = any
}

func Benchmark_MapOfString_Stdlib(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/map_of_string.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	var data interface{}
	for n := 0; n < b.N; n++ {
		any = json.Unmarshal(str, &data)
	}
	res = any
}

// these two are about the same
func Benchmark_ArrayOfInt(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/array_of_int.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		any = Unmarshall(str)
	}
	res = any
}
func Benchmark_ArrayOfInt_Stdlib(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/array_of_int.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	var data interface{}
	for n := 0; n < b.N; n++ {
		any = json.Unmarshal(str, &data)
	}
	res = any
}

// Array of Strigs
// There is a clear performance difference between these two - why?
// and mine is faster!
func Benchmark_ArrayOfString(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/array_of_string.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		any = Unmarshall(str)
	}
	res = any
}
func Benchmark_ArrayOfString_Stdlib(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/array_of_string.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	var data interface{}
	for n := 0; n < b.N; n++ {
		any = json.Unmarshal(str, &data)
	}
	res = any
}

// Crazy difference in performance and mine is faster
func Benchmark_BigString(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/big_string.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		any = Unmarshall(str)
	}
	res = any
}
func Benchmark_BigString_Stdlib(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/big_string.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	var data interface{}
	for n := 0; n < b.N; n++ {
		any = json.Unmarshal(str, &data)
	}
	res = any
}

// Interesting how my big gains in string processing evaporate here 🤔

func Benchmark_Code(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/code.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		any = Unmarshall(str)
	}
	res = any
}

func Benchmark_Code_Stdlib(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/code.json")
	if err != nil {
		panic(err)
	}
	b.StartTimer()

	var data interface{}
	for n := 0; n < b.N; n++ {
		any = json.Unmarshal(str, &data)
	}
	res = any
}

// similar performance
func Benchmark_NestedJson(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/nested_array.json")
	b.StartTimer()
	if err != nil {
		panic(err)
	}
	for n := 0; n < b.N; n++ {
		any = Unmarshall(str)
	}
	res = any
}
func Benchmark_NestedJson_Stdlib(b *testing.B) {
	var any Any
	b.StopTimer()
	str, err := ioutil.ReadFile("testdata/nested_array.json")
	b.StartTimer()
	if err != nil {
		panic(err)
	}
	var data interface{}
	for n := 0; n < b.N; n++ {
		any = json.Unmarshal(str, &data)
	}
	res = any
}
