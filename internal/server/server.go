package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/santonov10/otus_hw_project/internal/config"
	previewImageDelivery "github.com/santonov10/otus_hw_project/internal/imagepreview/delivery/http"
	"github.com/santonov10/otus_hw_project/internal/imagepreview/repository/lrucache"
	previewImageUC "github.com/santonov10/otus_hw_project/internal/imagepreview/usecase"
)

type Server struct { // TODO
	server *http.Server
}

func NewServer() *Server {
	config := config.Get()
	r := gin.Default()

	previewImageRepo, err := lrucache.NewLruCacheImages(config.CacheImagesLRU.Capacity, config.CacheImagesLRU.Dir)
	_ = previewImageRepo.ClearCacheDirImages()
	if err != nil {
		log.Fatalln(err)
	}
	mainGroup := r.Group("/")

	previewImageUC := previewImageUC.NewImagePreviewUseCase(previewImageRepo)
	previewImageDelivery.RegisterPreviewImageEndPoints(mainGroup, previewImageUC)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		server: s,
	}
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	fmt.Println("закрываем сервер")
	return err
}
