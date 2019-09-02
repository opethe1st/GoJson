package json

import (
	"fmt"
)

// ExampleLoad demonstrates how to use Load
func ExampleUnmarshall() {
	s := Unmarshall([]byte(`["v1", "v2"]`))
	fmt.Println(s)
	// Output:
	// [v1 v2]
}
