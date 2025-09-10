package requesto

import "errors"

var (
	ErrReadingBody         = errors.New("requesto: error reading response body")
	ErrUnmarshallingJSON   = errors.New("requesto: error unmarshalling JSON response")
	ErrUnmarshallingStruct = errors.New("requesto: error unmarshalling struct from JSON response")
)
