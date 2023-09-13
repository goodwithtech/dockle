package types

import "os"

type FileMap map[string]FileData

type FileData struct {
	Body     []byte
	FileMode os.FileMode
}

type FilterFunc func(filePath string, fileMode os.FileMode) (bool, error)
