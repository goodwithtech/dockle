package types

const (
	MinTypeNumber = AvoidRootDefault

	// CIS-DI
	AvoidRootDefault = iota
	UseContentTrust
	AddHealthcheck
	UseAptGetUpdateNoCache

	// TODO: change fanal FileMap structure
	RemoveSetuidSetgid
	UseCOPY
	AvoidEnvKeySecret
	AvoidCredentialFile

	// DG-DI
	AvoidSudo
	AvoidSensitiveDirectoryMounting
	AvoidDistUpgrade
	UseApkAddNoCache
	MinimizeAptGet
	AvoidLatestTag

	// DG-LI
	AvoidEmptyPassword
	AvoidDuplicateUser
	AvoidDuplicateGroup

	MaxTypeNumber = AvoidDuplicateGroup
)

const (
	PassLevel = iota + 1
	IgnoreLevel
	SkipLevel
	InfoLevel
	WarnLevel
	FatalLevel
)

type AlertDetail struct {
	DefaultLevel int
	Title        string
	Code         string
}

var AlertDetails = map[int]AlertDetail{
	AvoidRootDefault: {
		DefaultLevel: WarnLevel,
		Title:        "Create a user for the container",
		Code:         "CIS-DI-0001",
	},

	UseContentTrust: {
		DefaultLevel: InfoLevel,
		Title:        "Enable Content trust for Docker",
		Code:         "CIS-DI-0005",
	},

	AddHealthcheck: {
		DefaultLevel: WarnLevel,
		Title:        "Add HEALTHCHECK instruction to the container image",
		Code:         "CIS-DI-0006",
	},

	UseAptGetUpdateNoCache: {
		DefaultLevel: FatalLevel,
		Title:        "Do not use update instructions alone in the Dockerfile",
		Code:         "CIS-DI-0007",
	},

	RemoveSetuidSetgid: {
		DefaultLevel: InfoLevel,
		Title:        "Remove setuid and setgid permissions in the images",
		Code:         "CIS-DI-0008",
	},
	UseCOPY: {
		DefaultLevel: FatalLevel,
		Title:        "Use COPY instead of ADD in Dockerfile",
		Code:         "CIS-DI-0009",
	},

	AvoidEnvKeySecret: {
		DefaultLevel: FatalLevel,
		Title:        "Do not store secrets in ENVIRONMENT variables",
		Code:         "CIS-DI-0010",
	},
	AvoidCredentialFile: {
		DefaultLevel: FatalLevel,
		Title:        "Do not store secret files",
		Code:         "CIS-DI-0010",
	},

	// Docker Guard Checkpoints for Docker
	AvoidSudo: {
		DefaultLevel: FatalLevel,
		Title:        "Avoid sudo command",
		Code:         "DKL-DI-0001",
	},

	AvoidSensitiveDirectoryMounting: {
		DefaultLevel: FatalLevel,
		Title:        "Avoid sensitive directory mounting",
		Code:         "DKL-DI-0002",
	},
	AvoidDistUpgrade: {
		DefaultLevel: FatalLevel,
		Title:        "Avoid apt-get/apk/dist-upgrade",
		Code:         "DKL-DI-0003",
	},
	UseApkAddNoCache: {
		DefaultLevel: FatalLevel,
		Title:        "Use apk add with --no-cache",
		Code:         "DKL-DI-0004",
	},
	MinimizeAptGet: {
		DefaultLevel: FatalLevel,
		Title:        "Clear apt-get caches",
		Code:         "DKL-DI-0005",
	},
	AvoidLatestTag: {
		DefaultLevel: WarnLevel,
		Title:        "Avoid latest tag",
		Code:         "DKL-DI-0006",
	},

	// Docker Guard Checkpoints for Linux
	AvoidEmptyPassword: {
		DefaultLevel: FatalLevel,
		Title:        "Avoid empty password",
		Code:         "DKL-LI-0001",
	},
	AvoidDuplicateUser: {
		DefaultLevel: FatalLevel,
		Title:        "Be unique UID",
		Code:         "DKL-LI-0002",
	},
	AvoidDuplicateGroup: {
		DefaultLevel: FatalLevel,
		Title:        "Be unique GROUP",
		Code:         "DKL-LI-0002",
	},
}
