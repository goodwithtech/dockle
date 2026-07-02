package types

import (
	"archive/tar"
	"os"
)

type FilterFunc func(*tar.Header) (bool, error)

type FileMap map[string]FileData
type FileData struct {
	Body     []byte
	FileMode os.FileMode
}
