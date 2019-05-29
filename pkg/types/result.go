package types

type ScanResult map[string]error

const (
	SetPassword = iota
	AvoidDuplicateUser
	AvoidDuplicateGroup
	AvoidRootDefault
	AvoidRootRun
	AvoidLargeImage
	UseHealthcheck
	AvoidUpdate
	AvoidEnvKeySecret
	UseContentTrust

	DeleteTmpFiles
	DeleteCacheFiles
	PHPini
	AvoidCredential
	InvalidHost
	FilePermission
	RunSingleProcess
	AvoidLatestTag
)

const (
	InfoLevel = iota
	WarnLevel
	FatalLevel
	_minLevel = InfoLevel
	_maxLevel = FatalLevel
)

var AlertLevels = map[int]int{
	SetPassword:         FatalLevel,
	AvoidRootDefault:    WarnLevel,
	AvoidDuplicateUser:  WarnLevel,
	AvoidDuplicateGroup: WarnLevel,
	AvoidRootRun:        WarnLevel,
	AvoidLargeImage:     InfoLevel,
	UseHealthcheck:      InfoLevel,
	AvoidUpdate:         InfoLevel,
	AvoidEnvKeySecret:   WarnLevel,
	UseContentTrust:     WarnLevel,
}

type Assessment struct {
	Type     int
	Filename string
	Desc     string
}
