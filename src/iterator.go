package json

type iterator struct {
	s      []byte
	Offset int
}

func (iter *iterator) Current() byte {
	if iter.Offset < len(iter.s){
		return iter.s[iter.Offset]
	}
	return 0
}

func (iter *iterator) Next() {
	iter.Offset++
}

func (iter *iterator) HasNext() bool {
	return iter.Offset < len(iter.s)
}

func (iter *iterator) Slice(start int, end int) []byte {
	return iter.s[start:end]
}

func (iter *iterator) Len() int {
	return len(iter.s)
}

func (iter *iterator) AdvancePassWhiteSpace(){
	for iter.HasNext() && isSpace(iter.Current()) {
		iter.Next()
	}
}

func isSpace(ch byte) bool {
	switch ch {
	case '\t', '\n', '\r', ' ':
		return true
	}
	return false
}
