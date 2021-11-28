package usecase

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/santonov10/otus_hw_project/internal/imagepreview/repository/lrucache"
	"github.com/stretchr/testify/require"
)

var (
	testCacheDir  = "./testcachedir"
	testResizeDir = "./resizeTest"
	badURL        = "http://badURL"
	notImageURL   = "https://github.com/OtusGolang/final_project/tree/master/examples/image-previewer"
	imageURL      = "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg" //nolint
)

func getUseCaseWithLruCache() (*imagePreviewUseCase, error) {
	cacheRepo, err := lrucache.NewLruCacheImages(10, testCacheDir)
	if err != nil {
		return nil, err
	}

	return NewImagePreviewUseCase(cacheRepo), nil
}

func removeTestDir() {
	os.RemoveAll(testCacheDir)
}

func getRequestHeader() http.Header {
	request, _ := http.NewRequest("Get", "", nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36") //nolint
	return request.Header
}

func TestImagePreviewUseCase(t *testing.T) {
	t.Run("getBodyFromUrl", func(t *testing.T) {
		_, err := getBodyFromURL(context.TODO(), badURL, getRequestHeader()) // wrongurl
		require.Error(t, err)

		data, err := getBodyFromURL(context.TODO(), imageURL, getRequestHeader()) // okUrl
		require.NoError(t, err)
		require.NotEmpty(t, data)
	})

	t.Run("GetPreviewImageFromUrl", func(t *testing.T) {
		defer removeTestDir()

		imagePreviewUC, err := getUseCaseWithLruCache()
		require.NoError(t, err)

		// добавляет пустой файл в кеш.
		data, ok := imagePreviewUC.GetPreviewImageFromUrl(badURL, getRequestHeader(), 200, 200)
		require.False(t, ok)
		require.Empty(t, data.Bytes())

		// добавляет пустой файл в кеш.
		data, ok = imagePreviewUC.GetPreviewImageFromUrl(notImageURL, getRequestHeader(), 200, 200)
		require.False(t, ok)
		require.Empty(t, data.Bytes())

		// добавляет файл в кеш.
		data, ok = imagePreviewUC.GetPreviewImageFromUrl(imageURL, getRequestHeader(), 200, 200)
		require.True(t, ok)
		require.NotEmpty(t, data.Bytes())

		files, _ := ioutil.ReadDir(testCacheDir)
		require.Equal(t, 3, len(files), "3 файла = 3 записям в кеше")

		// из кеша.
		data, ok = imagePreviewUC.GetPreviewImageFromUrl(badURL, getRequestHeader(), 200, 200)
		require.False(t, ok)
		require.Empty(t, data.Bytes())

		// из кеша.
		data, ok = imagePreviewUC.GetPreviewImageFromUrl(notImageURL, getRequestHeader(), 200, 200)
		require.False(t, ok)
		require.Empty(t, data.Bytes())

		// из кеша.
		data, ok = imagePreviewUC.GetPreviewImageFromUrl(imageURL, getRequestHeader(), 200, 200)
		require.True(t, ok)
		require.NotEmpty(t, data.Bytes())

		files, _ = ioutil.ReadDir(testCacheDir)
		require.Equal(t, 3, len(files), "не должны создаваться лишние файлы т.к. берется из кэша")
	})
}
