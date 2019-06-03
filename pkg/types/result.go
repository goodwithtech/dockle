package types

type ScanResult map[string]error

const (
	SetPassword = iota
	AvoidDuplicateUser
	AvoidDuplicateGroup
	AvoidRootDefault
	UseHealthcheck
	AvoidUpdate
	AvoidUpgrade
	AvoidSudo
	AvoidEnvKeySecret
	AvoidMountSensitiveDir
	AvoidCredential
	UseContentTrust

	AvoidLargeImage
	AvoidRootRun
	DeleteTmpFiles
	DeleteCacheFiles
	PHPini
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
		Title:        "Check running user isn't root",
	},
	AvoidDuplicateUser: {
		DefaultLevel: WarnLevel,
		Title:        "Check user names",
	},
	AvoidDuplicateGroup: {
		DefaultLevel: WarnLevel,
		Title:        "Check group names",
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
		Title:        "Check update commands",
	},
	AvoidUpgrade: {
		DefaultLevel: WarnLevel,
		Title:        "Check upgrade commands",
	},
	AvoidSudo: {
		DefaultLevel: WarnLevel,
		Title:        "Check sudo commands",
	},
	AvoidEnvKeySecret: {
		DefaultLevel: FatalLevel,
		Title:        "Check environment vars",
	},
	AvoidMountSensitiveDir: {
		DefaultLevel: FatalLevel,
		Title:        "Check volumes",
	},
	AvoidCredential: {
		DefaultLevel: WarnLevel,
		Title:        "Check credential files",
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
