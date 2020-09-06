package domain

type FileMeta struct {
	Path        string
	ModifyTime  int64
	IsDirectory bool
	Size        int64
}

type SyncMode string

const (
	SyncModeCommon SyncMode = "Common"
	SyncModeClear  SyncMode = "Clear"
)
