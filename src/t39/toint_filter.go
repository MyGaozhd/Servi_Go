package t39

import (
	"errors"
	"strconv"
)

var ToIntFilterWrongFormatError = errors.New("input data should be []string")

type ToIntFilter struct {
}

func NewTointFilter() *ToIntFilter {
	return &ToIntFilter{}
}

func (ti *ToIntFilter) Process(data Request) (Response, error) {
	parts, ok := data.([]string)
	if !ok {
		return nil, ToIntFilterWrongFormatError
	}

	ret := []int{}

	for _, part := range parts {
		s, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		} else {
			ret = append(ret, s)
		}
	}

	return ret, nil
}
