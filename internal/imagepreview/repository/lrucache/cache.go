package lrucache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type cacheItem struct {
	key      string
	fileName string
}

type lruCacheImages struct {
	capacity        int
	queue           List
	mu              sync.Mutex
	items           map[string]*ListItem
	imgCacheDirPath string
}

// id файла для кэша.
var cacheImageID uint32

func (c *lruCacheImages) Set(key string, value *bytes.Buffer) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fileName, err := c.saveToNewCacheFile(value)
	if err != nil {
		return false, err
	}

	if listItem, ok := c.items[key]; ok {
		oldFileName := listItem.Value.(*cacheItem).fileName
		c.deleteCacheFile(oldFileName)
		listItem.Value.(*cacheItem).fileName = fileName
		c.queue.MoveToFront(listItem)
		return true, nil
	}

	cItem := &cacheItem{key: key, fileName: fileName}
	insertItem := NewListItem(cItem)
	c.items[key] = insertItem
	c.queue.PushFront(insertItem)
	if c.queue.Len() > c.capacity {
		removeItem := c.queue.Back()
		c.removeItem(removeItem)
	}
	return false, nil
}

func (c *lruCacheImages) removeItem(removeItem *ListItem) error {
	removeKey := removeItem.Value.(*cacheItem).key
	removeFileName := removeItem.Value.(*cacheItem).fileName
	delete(c.items, removeKey)
	c.queue.Remove(removeItem)
	return c.deleteCacheFile(removeFileName)
}

func (c *lruCacheImages) Get(key string) (*bytes.Buffer, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if listItem, ok := c.items[key]; ok {
		data, err := c.getCacheFile(listItem.Value.(*cacheItem).fileName)
		if err != nil {
			return nil, false
		}
		c.queue.MoveToFront(listItem)
		return data, true
	}
	return nil, false
}

func (c *lruCacheImages) getCacheFile(fileName string) (*bytes.Buffer, error) {
	data, err := os.Open(c.getPathToCacheFile(fileName))
	if err != nil {
		return nil, err
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)

	var count int
	for {
		count, err = reader.Read(part)
		if err != nil {
			break
		}
		buffer.Write(part[:count])
	}

	if !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("Error Reading "+fileName+": %w", err)
	}

	return buffer, nil
}

func (c *lruCacheImages) saveToNewCacheFile(bytes *bytes.Buffer) (fileName string, err error) {
	if bytes == nil {
		return "", errors.New("значения для записи = nil")
	}

	fileName = c.generateNewFileName(".jpg")
	f, err := os.Create(c.getPathToCacheFile(fileName))
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(bytes.Bytes())
	if err != nil {
		return "", err
	}
	return fileName, nil
}

// Генерируем новое имя файла для кэша.
func (c *lruCacheImages) generateNewFileName(format string) string {
	return strconv.FormatUint(uint64(atomic.AddUint32(&cacheImageID, 1)), 10) + format
}

func (c *lruCacheImages) deleteCacheFile(fileName string) error {
	return os.Remove(c.getPathToCacheFile(fileName))
}

func (c *lruCacheImages) getPathToCacheFile(fileName string) string {
	return c.imgCacheDirPath + fileName
}

func (c *lruCacheImages) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, val := range c.items {
		c.removeItem(val)
	}
}

func NewLruCacheImages(capacity int, imgCacheDirPath string) (*lruCacheImages, error) { //nolint
	cache := &lruCacheImages{
		capacity:        capacity,
		imgCacheDirPath: imgCacheDirPath,
	}

	if err := cache.init(); err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *lruCacheImages) init() error {
	if c.imgCacheDirPath[len(c.imgCacheDirPath)-1:] != "/" {
		c.imgCacheDirPath += "/"
	}

	c.items = make(map[string]*ListItem, c.capacity)
	c.queue = NewList()
	err := c.createCacheDirIfNotExist()
	if err != nil {
		return fmt.Errorf("ошибка инициализации:%w", err)
	}
	if err != nil {
		return fmt.Errorf("ошибка инициализации:%w", err)
	}
	return nil
}

func (c *lruCacheImages) createCacheDirIfNotExist() error {
	_, err := os.Stat(c.imgCacheDirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(c.imgCacheDirPath, 0o755)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

// ClearCacheDirImages Удаляем все картинки в кэш категории.
func (c *lruCacheImages) ClearCacheDirImages() error {
	deleteExtensions := map[string]interface{}{
		".jpg": nil, ".png": nil,
	}
	files, err := ioutil.ReadDir(c.imgCacheDirPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		fileName := f.Name()
		fileExtension := fileName[len(fileName)-4:]
		if _, ok := deleteExtensions[fileExtension]; ok {
			_ = os.Remove(c.getPathToCacheFile(fileName))
		}
	}
	return nil
}
