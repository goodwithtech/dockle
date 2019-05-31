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
	MinTypeNumber = SetPassword
	MaxTypeNumber = UseContentTrust
)

const (
	InfoLevel = iota
	WarnLevel
	FatalLevel
	Pass
	_minLevel_ = InfoLevel
	_maxLevel_ = FatalLevel
)

type AlertDetail struct {
	DefaultLevel int
	Title        string
}

var AlertDetails = map[int]AlertDetail{
	SetPassword: {
		DefaultLevel: FatalLevel,
		Title:        "Check password",
	},
	AvoidRootDefault: {
		DefaultLevel: WarnLevel,
		Title:        "Running as root",
	},
	AvoidDuplicateUser: {
		DefaultLevel: WarnLevel,
		Title:        "Check users",
	},
	AvoidDuplicateGroup: {
		DefaultLevel: WarnLevel,
		Title:        "Check groups",
	},
	AvoidLargeImage: {
		DefaultLevel: InfoLevel,
		Title:        "Check image size",
	},
	UseHealthcheck: {
		DefaultLevel: InfoLevel,
		Title:        "Check healthcheck setting",
	},
	AvoidUpdate: {
		DefaultLevel: WarnLevel,
		Title:        "Check commands",
	},
	AvoidEnvKeySecret: {
		DefaultLevel: FatalLevel,
		Title:        "Check environment vars",
	},
	UseContentTrust: {
		DefaultLevel: WarnLevel,
		Title:        "Check DOCKER CONTENT TRUST setting",
	},
}

type Assessment struct {
	Type     int
	Filename string
	Desc     string
}
