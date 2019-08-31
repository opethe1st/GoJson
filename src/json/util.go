package json

import (
	"fmt"
	"math"
)

func errorMsg(iter *iterator, msg string, msgArgs ...interface{}) string {
	startBefore := iter.Offset - 50
	if startBefore < 0 {
		startBefore = 0
	}
	endBefore := iter.Offset
	if endBefore < 0 {
		endBefore = 0
	}
	before := iter.Slice(startBefore, endBefore)

	startAfter := iter.Offset + 1
	if startAfter > iter.Len() {
		startAfter = iter.Len()
	}
	endAfter := iter.Offset + 50
	if endAfter > iter.Len() {
		endAfter = iter.Len()
	}
	after := iter.Slice(startAfter, endAfter)
	// this doesnt work well if the character to underline is a whitespace
	underlined := string([]byte{iter.Current(), 204, 179})
	return fmt.Sprintf(`There is an error around
	%s%s%s

	%s
	`, before, underlined, after, fmt.Sprintf(msg, msgArgs...))
}

func floatEquals(a, b float64) bool {
	if math.Abs(a-b) < 0.00000001 {
		return true
	}
	return false
}
