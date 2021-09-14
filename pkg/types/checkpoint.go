package types

const (
	// CIS-DI
	AvoidRootDefault       = "CIS-DI-0001"
	UseContentTrust        = "CIS-DI-0005"
	AddHealthcheck         = "CIS-DI-0006"
	UseAptGetUpdateNoCache = "CIS-DI-0007"
	CheckSuidGuid          = "CIS-DI-0008"
	UseCOPY                = "CIS-DI-0009"
	AvoidCredential        = "CIS-DI-0010"

	// DG-DI
	AvoidSudo                       = "DKL-DI-0001"
	AvoidSensitiveDirectoryMounting = "DKL-DI-0002"
	AvoidDistUpgrade                = "DKL-DI-0003"
	UseApkAddNoCache                = "DKL-DI-0004"
	MinimizeAptGet                  = "DKL-DI-0005"
	AvoidLatestTag                  = "DKL-DI-0006"

	// DG-LI
	AvoidEmptyPassword      = "DKL-LI-0001"
	AvoidDuplicateUserGroup = "DKL-LI-0002"
	InfoDeletableFiles      = "DKL-LI-0003"
)

const (
	PassLevel int = iota + 1
	IgnoreLevel
	SkipLevel
	InfoLevel
	WarnLevel
	FatalLevel
)

// DefaultLevelMap save risk level each checkpoints
var DefaultLevelMap = map[string]int{
	AvoidRootDefault:       WarnLevel,
	UseContentTrust:        InfoLevel,
	AddHealthcheck:         InfoLevel,
	UseAptGetUpdateNoCache: FatalLevel,
	CheckSuidGuid:          InfoLevel,
	UseCOPY:                FatalLevel,
	AvoidCredential:        FatalLevel,

	AvoidSudo:                       FatalLevel,
	AvoidSensitiveDirectoryMounting: FatalLevel,
	AvoidDistUpgrade:                WarnLevel,
	UseApkAddNoCache:                FatalLevel,
	MinimizeAptGet:                  FatalLevel,
	AvoidLatestTag:                  WarnLevel,

	AvoidEmptyPassword:      FatalLevel,
	AvoidDuplicateUserGroup: FatalLevel,
	InfoDeletableFiles:      InfoLevel,
}

// TitleMap save title each checkpoints
var TitleMap = map[string]string{
	AvoidRootDefault:                "Create a user for the container",
	UseContentTrust:                 "Enable Content trust for Docker",
	AddHealthcheck:                  "Add HEALTHCHECK instruction to the container image",
	UseAptGetUpdateNoCache:          "Do not use update instructions alone in the Dockerfile",
	CheckSuidGuid:                   "Confirm safety of setuid/setgid files",
	UseCOPY:                         "Use COPY instead of ADD in Dockerfile",
	AvoidCredential:                 "Do not store credential in environment variables/files",
	AvoidSudo:                       "Avoid sudo command",
	AvoidSensitiveDirectoryMounting: "Avoid sensitive directory mounting",
	AvoidDistUpgrade:                `Avoid "apt-get dist-upgrade"`,
	UseApkAddNoCache:                `Use "apk add" with --no-cache`,
	MinimizeAptGet:                  `Clear apt-get caches`,
	AvoidLatestTag:                  "Avoid latest tag",
	AvoidEmptyPassword:              "Avoid empty password",
	AvoidDuplicateUserGroup:         "Be unique UID/GROUP",
	InfoDeletableFiles:              "Only put necessary files",
}
