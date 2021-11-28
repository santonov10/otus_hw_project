package http

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/santonov10/otus_hw_project/internal/imagepreview"
)

type Handler struct {
	useCase imagepreview.UseCase
}

func NewHandler(useCase imagepreview.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetCachedImagePreview(ctx *gin.Context) {
	width, err := strconv.Atoi(ctx.Param("width"))
	if err != nil {
		responseImageBadGateway(ctx)
		return
	}
	height, err := strconv.Atoi(ctx.Param("height"))
	if err != nil {
		responseImageBadGateway(ctx)
		return
	}
	imgURL := "http:/" + ctx.Param("imgUrl")

	content, ok := h.useCase.GetPreviewImageFromUrl(imgURL, ctx.Request.Header, width, height)

	if !ok {
		responseImageBadGateway(ctx)
		return
	}
	responseImageOk(ctx, content)
}

func responseImageBadGateway(ctx *gin.Context) {
	responseCode := http.StatusBadGateway
	ctx.Header("Content-Type", "image/jpeg")
	ctx.Header("Content-Length", "0")
	ctx.String(responseCode, "")
}

func responseImageOk(ctx *gin.Context, content *bytes.Buffer) {
	responseCode := http.StatusOK
	contentBytes := content.Bytes()
	ctx.Header("Content-Type", "image/jpeg")
	ctx.Header("Content-Length", strconv.Itoa(len(contentBytes)))

	ctx.String(responseCode, string(contentBytes))
}
