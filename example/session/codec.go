package session

import "errors"

type codec struct {
}

func (c codec) Encode(pkg interface{}) ([]byte, error) {
	if pkg == nil {
		return nil, errors.New("pkg is illegal")
	}
	data, ok := pkg.(string)
	if !ok {
		return nil, errors.New("pkg type must be string")
	}

	if len(data) != 5 || data != "hello" {
		return nil, errors.New("pkg string must be \"hello\"")
	}

	return []byte(data), nil
}

func (c codec) Decode(bytes []byte) (interface{}, int, error) {
	//TODO implement me
	panic("implement me")
}
