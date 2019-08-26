package json

type iterator struct {
	s      string
	offset int
}

func (iter *iterator) getCurrent() byte {
	return iter.s[iter.offset]
}

func (iter *iterator) advance() {
	iter.offset++
}

func (iter *iterator) isEnd() bool {
	return iter.offset >= len(iter.s)
}
