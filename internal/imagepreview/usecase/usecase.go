package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/santonov10/otus_hw_project/internal/imagepreview"
)

type imagePreviewUseCase struct {
	imageRepo imagepreview.ImageCacheRepository
}

func NewImagePreviewUseCase(cacheRepo imagepreview.ImageCacheRepository) *imagePreviewUseCase { //nolint
	return &imagePreviewUseCase{
		imageRepo: cacheRepo,
	}
}

func (i *imagePreviewUseCase) GetPreviewImageFromUrl(url string, header http.Header, width, height int) (data *bytes.Buffer, ok bool) { //nolint
	cacheKey := fmt.Sprintf("%d/%d/%s", width, height, url)
	img, ok := i.imageRepo.Get(cacheKey)
	if !ok {
		errorEmptyContent := bytes.NewBuffer([]byte{})
		imageBody, err := getBodyFromURL(context.TODO(), url, header)
		if err != nil {
			i.imageRepo.Set(cacheKey, errorEmptyContent)
			return errorEmptyContent, false
		}
		img, err = resizeImage(bytes.NewReader(imageBody), width, height)
		if err != nil {
			i.imageRepo.Set(cacheKey, errorEmptyContent)
			return errorEmptyContent, false
		}
		i.imageRepo.Set(cacheKey, img)
	}

	haveImg := false
	if img.Len() > 0 {
		haveImg = true
	}
	return img, haveImg
}

func getBodyFromURL(ctx context.Context, url string, header http.Header) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
