package firstservice

import (
	"encoding/json"
)


// Error is a struct for JSON answer
type Error struct {
	Err string `json:"error"`
}

//NewError Create new error from error interface
func NewError(err error) *Error {
	return &Error{Err: err.Error()}
}

//Marshall return a slice of bytes of JSON and error
func (err *Error) Marshall() ([]byte, error) {
	return json.Marshal(err)
}


