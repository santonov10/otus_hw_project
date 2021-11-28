package imagepreview

import (
	"bytes"
	"net/http"
)

type UseCase interface {
	GetPreviewImageFromUrl(url string, header http.Header, width, height int) (data *bytes.Buffer, ok bool)
}
