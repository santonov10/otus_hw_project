package imagepreview

import (
	"bytes"
)

type ImageCacheRepository interface {
	Set(key string, value *bytes.Buffer) (bool, error)
	Get(key string) (*bytes.Buffer, bool)
	Clear()
	ClearCacheDirImages() error
}
