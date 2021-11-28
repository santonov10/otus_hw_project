package lrucache

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var testDir = "./testcachedir"

func getTestLruCacheItem(capacity int) (*lruCacheImages, error) {
	return NewLruCacheImages(capacity, testDir)
}

func removeTestDir() {
	os.RemoveAll(testDir)
}

func TestNewLruCacheImages(t *testing.T) {
	t.Run("test CreateDir", func(t *testing.T) {
		removeTestDir()
		defer removeTestDir()
		_, err := getTestLruCacheItem(10)
		require.NoError(t, err)
	})

	t.Run("save cache to file", func(t *testing.T) {
		removeTestDir()
		defer removeTestDir()
		c, err := getTestLruCacheItem(10)
		require.NoError(t, err)

		a := []byte{'a'}
		filename, err := c.saveToNewCacheFile(bytes.NewBuffer(a))
		require.NoError(t, err)
		require.FileExists(t, c.getPathToCacheFile(filename))

		c.deleteCacheFile(filename)
		require.NoFileExists(t, c.getPathToCacheFile(filename))
	})

	t.Run("clearFiles", func(t *testing.T) {
		defer removeTestDir()
		c, err := getTestLruCacheItem(10)
		require.NoError(t, err)

		fileNameShouldClear := "notCleared.jpg"
		fileNameShouldStay := "randomFile.txt"
		f1, err := os.Create(c.getPathToCacheFile(fileNameShouldClear))
		f1.Close()
		require.NoError(t, err)
		f2, err := os.Create(c.getPathToCacheFile(fileNameShouldStay))
		f2.Close()
		require.NoError(t, err)

		err = c.ClearCacheDirImages()
		require.NoError(t, err)

		require.FileExists(t, c.getPathToCacheFile(fileNameShouldStay))
		require.NoFileExists(t, c.getPathToCacheFile(fileNameShouldClear))
	})

	t.Run("empty cache", func(t *testing.T) {
		defer removeTestDir()
		c, err := getTestLruCacheItem(10)
		require.NoError(t, err)

		var val *bytes.Buffer
		_, err = c.Set("aaa", val)
		require.Error(t, err, "должна быть ошибка записи nil value")

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		defer removeTestDir()
		c, err := getTestLruCacheItem(5)
		require.NoError(t, err)

		a := []byte{'a'}
		b := []byte{'b', 'b'}
		wasInCache, err := c.Set("aaa", bytes.NewBuffer(a))
		require.NoError(t, err)
		require.False(t, wasInCache)

		wasInCache, err = c.Set("bbb", bytes.NewBuffer(b))
		require.NoError(t, err)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, a, val.Bytes())

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, b, val.Bytes())

		newA := []byte{'n', 'e', 'w'}
		wasInCache, err = c.Set("aaa", bytes.NewBuffer(newA))
		require.NoError(t, err)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, newA, val.Bytes())

		files, _ := ioutil.ReadDir(c.imgCacheDirPath)
		require.Equal(t, c.queue.Len(), len(files), "количество картинок должно оставаться количеству записей в кеше")
	})

	t.Run("purge logic", func(t *testing.T) {
		defer removeTestDir()
		c, err := getTestLruCacheItem(2)
		require.NoError(t, err)
		c.Set("a", bytes.NewBuffer([]byte{'2'})) // a
		c.Set("a", bytes.NewBuffer([]byte{'a'})) // a
		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, val.Bytes(), []byte{'a'})

		c.Set("b", bytes.NewBuffer([]byte{'b'})) // b,a
		c.Get("a")                               // a,b
		c.Set("c", bytes.NewBuffer([]byte{'c'})) // c,a
		c.Set("c", bytes.NewBuffer([]byte{'c'})) // c,a
		val, ok = c.Get("b")                     // больше нет в кеше
		require.False(t, ok)
		require.Nil(t, val)
		_, ok = c.Get("a") // a,c
		require.True(t, ok)
		_, ok = c.Get("c") // c,a
		require.True(t, ok)

		files, _ := ioutil.ReadDir(c.imgCacheDirPath) // 2 файла
		require.Equal(t, c.queue.Len(), len(files), "количество картинок должно оставаться количеству записей в кеше")

		c.Clear()
		_, ok = c.Get("a") //
		require.False(t, ok)
		_, ok = c.Get("c") //
		require.False(t, ok)

		files, _ = ioutil.ReadDir(c.imgCacheDirPath) // 0 файлов
		require.Equal(t, c.queue.Len(), len(files), "количество картинок должно оставаться количеству записей в кеше")
	})
}

func TestCacheMultithreading(t *testing.T) {
	defer removeTestDir()
	c, err := getTestLruCacheItem(2)
	require.NoError(t, err)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_0; i++ {
			c.Set(strconv.Itoa(i), new(bytes.Buffer))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_0; i++ {
			c.Get(strconv.Itoa(rand.Intn(1_000_000))) // nolint
		}
	}()

	wg.Wait()
}
