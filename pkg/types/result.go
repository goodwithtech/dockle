package types

type ScanResult map[string]error

const (
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
