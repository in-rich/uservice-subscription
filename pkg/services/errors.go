package services

import "errors"

var (
	ErrNoteEditsExhausted = errors.New("note edits exhausted")

	ErrInvalidRequest = errors.New("invalid request")
)
