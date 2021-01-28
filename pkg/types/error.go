package types

import "errors"

var (
	ErrSetImageOrFile          = errors.New("image name or image file must be specified")
	ErrorCreateDockerExtractor = errors.New("error create docker extractor")
	ErrorAnalyze               = errors.New("error analyze")
)

type ScanError struct {
	ErrMsg  string
	ErrType ScanErrorType
}

type ScanErrorType string

const (
	CreateDockerExtractor ScanErrorType = "errorCreateDockerExtractor"
	Analyze               ScanErrorType = "errorAnalyze"
	SetImageOrFile        ScanErrorType = "errorSetImageOrFile"
	Unknown               ScanErrorType = "unknown"
)

func ConvertError(err error) *ScanError {
	if errors.Is(err, ErrorAnalyze) {
		return &ScanError{
			ErrMsg:  err.Error(),
			ErrType: Analyze,
		}
	} else if errors.Is(err, ErrorCreateDockerExtractor) {
		return &ScanError{
			ErrMsg:  err.Error(),
			ErrType: CreateDockerExtractor,
		}
	} else if errors.Is(err, ErrSetImageOrFile) {
		return &ScanError{
			ErrMsg:  err.Error(),
			ErrType: SetImageOrFile,
		}
	} else {
		return &ScanError{
			ErrMsg:  err.Error(),
			ErrType: Unknown,
		}
	}
}
