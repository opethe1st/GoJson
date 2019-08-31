package json


// Error is a json error
type Error struct{

}

func (e Error) Error() string{
	return ""
}

// Validate validates that json is well-formed json. It returns nil if well-form or an error object
// with details on what was wrong with the Json
func Validate(s []byte) error {
	return validate(&iterator{s: s})
}


func validate(iter *iterator) error {
	iter.AdvancePastAllWhiteSpace()
	switch {
	case iter.Current() == 'n':
		return validateKeyword(iter, "null")
	case iter.Current() == 't':
		return validateKeyword(iter, "true")
	case iter.Current() == 'f':
		return validateKeyword(iter, "false")
	case isNumber(iter):
		return validateNumber(iter)
	case iter.Current() == '"':
		return validateString(iter)
	case iter.Current() == '[':
		return validateArray(iter)
	case iter.Current() == '{':
		return validateObject(iter)
	default:
		return Error{}
	}
}

func validateKeyword(iter *iterator, keyword string) error {
	if keyword == string(iter.Slice(iter.Offset, iter.Offset+len(keyword))){
		return nil
	}
	return Error{}
}


func validateNumber(iter *iterator) error{
	return nil
}

func validateString(iter *iterator) error {
	return nil
}

func validateArray(iter *iterator) error {
	return nil
}

func validateObject(iter *iterator) error {
	return nil
}
