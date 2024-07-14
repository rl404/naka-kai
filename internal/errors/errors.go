package errors

import "errors"

// Error list.
var (
	ErrInternalDB        = errors.New("internal database error")
	ErrInternalCache     = errors.New("internal cache error")
	ErrInternalServer    = errors.New("internal server error")
	ErrInvalidDBFormat   = errors.New("invalid db address")
	ErrInvalidYoutubeURL = errors.New("invalid youtube url")
	ErrInvalidYoutubeID  = errors.New("invalid youtube video id")
	ErrNotInVC           = errors.New("not in voice channel")
	ErrInvalidPrompt     = errors.New("invalid prompt response")
)
