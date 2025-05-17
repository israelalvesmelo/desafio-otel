package entity

import "errors"

var (
	ErrZipcodeNotFound = errors.New("can not found zipcode")
	ErrZipcodeNotValid = errors.New("invalid zipcode")
	ErrEmptyAPIkey     = errors.New("you should provide a not empty API key")
)
