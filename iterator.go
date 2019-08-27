package json

type iterator struct {
	s      string
	offset int
}

func (iter *iterator) Current() byte {
	return iter.s[iter.offset]
}

func (iter *iterator) Next() {
	iter.offset++
}

func (iter *iterator) HasNext() bool {
	return iter.offset < len(iter.s)
}
