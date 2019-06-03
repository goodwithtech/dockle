package types

const (
	// Security Checkpoints
	SetPassword = iota
	AvoidRootDefault
	AvoidSensitiveDirectoryMounting
	UseContentTrust
	AvoidEnvKeySecret
	AvoidCredentialFile
	AvoidDuplicateUser
	AvoidDuplicateGroup

	// Dockerfile Checkpoints
	AvoidUpgrade
	AvoidSudo
	UseNoCacheAPK
	MinimizeAptGet
	AvoidLatestTag

	MinTypeNumber = SetPassword
	MaxTypeNumber = AvoidLatestTag
)

const (
	InfoLevel = iota
	WarnLevel
	FatalLevel
	PassLevel
	SkipLevel
)

type AlertDetail struct {
	DefaultLevel int
	Title        string
	Code         string
}

var AlertDetails = map[int]AlertDetail{
	SetPassword: {
		DefaultLevel: FatalLevel,
		Title:        "Check password",
		Code:         "SC0001",
	},
	AvoidRootDefault: {
		DefaultLevel: WarnLevel,
		Title:        "Check running user isn't root",
		Code:         "SC0002",
	},
	AvoidSensitiveDirectoryMounting: {
		DefaultLevel: FatalLevel,
		Title:        "Check volumes",
		Code:         "SC0003",
	},
	UseContentTrust: {
		DefaultLevel: WarnLevel,
		Title:        "Check DOCKER CONTENT TRUST setting",
		Code:         "SC0004",
	},
	AvoidEnvKeySecret: {
		DefaultLevel: FatalLevel,
		Title:        "Check environment vars",
		Code:         "SC0005",
	},
	AvoidCredentialFile: {
		DefaultLevel: WarnLevel,
		Title:        "Check credential files",
		Code:         "SC0005",
	},
	AvoidDuplicateUser: {
		DefaultLevel: WarnLevel,
		Title:        "Check user names",
		Code:         "SC0006",
	},
	AvoidDuplicateGroup: {
		DefaultLevel: WarnLevel,
		Title:        "Check group names",
		Code:         "SC0006",
	},

	AvoidUpgrade: {
		DefaultLevel: WarnLevel,
		Title:        "Check upgrade commands",
		Code:         "DC0001",
	},
	AvoidSudo: {
		DefaultLevel: WarnLevel,
		Title:        "Check sudo commands",
		Code:         "DC0002",
	},

	UseNoCacheAPK: {
		DefaultLevel: InfoLevel,
		Title:        "Check apk add command",
		Code:         "DC0003",
	},

	MinimizeAptGet: {
		DefaultLevel: InfoLevel,
		Title:        "Check apt-get install command",
		Code:         "DC0004",
	},

	AvoidLatestTag: {
		DefaultLevel: WarnLevel,
		Title:        "Check image tag",
		Code:         "DC0005",
	},
}

type Assessment struct {
	Type     int
	Level    int
	Filename string
	Desc     string
}
