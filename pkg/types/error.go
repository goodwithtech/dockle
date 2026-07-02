package types

import "errors"

var (
	ErrSetImageOrFile = errors.New("image name or image file must be specified")
	InvalidURLPattern = errors.New("invalid url pattern")
	ErrNoRpmCmd       = errors.New("no rpm command")
)
