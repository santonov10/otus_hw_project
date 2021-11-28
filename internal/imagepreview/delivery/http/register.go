package http

import (
	"github.com/gin-gonic/gin"
	"github.com/santonov10/otus_hw_project/internal/imagepreview"
)

func RegisterPreviewImageEndPoints(router *gin.RouterGroup, uc imagepreview.UseCase) {
	h := NewHandler(uc)
	router.GET("/fill/:width/:height/*imgUrl", h.GetCachedImagePreview)
}
