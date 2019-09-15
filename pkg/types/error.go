package types

import "errors"

var (
	ErrSetImageOrFile = errors.New("image name or image file must be specified")
)
