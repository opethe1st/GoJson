package json

// this is an iterator that keeps track of the last read position of s in Offset.
// s should never be changed since then the Offset and len won't make sense anymore
// after calling s.Len() - s.len should equal len(s)
type iterator struct {
	s      []byte
	Cursor int
	len    int
}

// Selectors

func (iter *iterator) Current() byte {
	if iter.Cursor < len(iter.s) {
		return iter.s[iter.Cursor]
	}
	return 0
}

func (iter *iterator) HasNext() bool {
	return iter.Cursor < len(iter.s)
}

func (iter *iterator) Slice(start int, end int) []byte {
	if end > len(iter.s) {
		end = len(iter.s)
	}
	return iter.s[start:end]
}

func (iter *iterator) SliceTillCursor(start int) []byte {
	return iter.s[start:iter.Cursor]
}

func (iter *iterator) Len() int {
	// is len(iter.s) cached?
	if iter.len != 0 {
		return iter.len
	}
	iter.len = len(iter.s)
	return len(iter.s)
}

// Mutators

func (iter *iterator) Next() {
	// I could have called this Advance to be consistent wih AdvancePast etc
	if iter.Cursor < iter.Len() {
		iter.Cursor++
	}
}

func (iter *iterator) AdvancePastAllWhiteSpace() {
	for isSpace(iter.Current()) {
		iter.Next()
	}
}

func (iter *iterator) AdvancePast(char byte) {
	iter.AdvancePastAllWhiteSpace()
	if iter.Current() == char {
		iter.Next()
	} else {
		panic(errorMsg(iter, "Was expecting %s but got %s instead", string(char), string(iter.Current())))
	}
}

func isSpace(ch byte) bool {
	switch ch {
	case '\t', '\n', '\r', ' ':
		return true
	}
	return false
}
