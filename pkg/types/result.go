package types

type ScanResult map[string]error

const (
	SetPassword      = "SetPassword"
	AvoidRootDefault = "AvoidRootDefault"
	AvoidRootRun     = "AvoidRootRun"
	LargeImage       = "LargeImage"
	DeleteTmpFiles   = "DeleteTmpFiles"
	DeleteCacheFiles = "DeleteCacheFiles"
	PHPini           = "PHPini"
	EnvKeySuspition  = "EnvKeySuspition"
	EnvVarSuspition  = "EnvVarSuspition"
	AvoidCredential  = "AvoidCredential"
	InvalidHost      = "InvalidEtcHost"
	FilePermission   = "FilePermission"
	RunSingleProcess = "RunSingleProcess"
	AvoidLatestTag   = "AvoidLatestTag"
)

const (
	InfoLevel = iota
	WarnLevel
	FatalLevel
	_minLevel = InfoLevel
	_maxLevel = FatalLevel
)

type Assessment struct {
	Level    int
	Filename string
	Desc     string
}
