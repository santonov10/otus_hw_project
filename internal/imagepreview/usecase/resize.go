package usecase

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"math"

	"golang.org/x/image/draw"
)

func resizeImage(imgDataReader io.Reader, width, height int) (*bytes.Buffer, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("ширина и высота картинки должны быть > 0")
	}

	src, err := jpeg.Decode(imgDataReader)
	if err != nil {
		return nil, err
	}

	var fromX, fromY, toX, toY int

	fromX = 0
	fromY = 0
	toX = src.Bounds().Max.X
	toY = src.Bounds().Max.Y

	originalAspectRatio := float64(src.Bounds().Max.X) / float64(src.Bounds().Max.Y)
	newAspectRatio := float64(width) / float64(height)

	// вычисляем от каких координат производить ресайз картинки если соотношение сторон отличается
	if !isAlmostEqual(originalAspectRatio, newAspectRatio) {
		if originalAspectRatio > newAspectRatio {
			toX = src.Bounds().Max.X - int((float64(src.Bounds().Max.Y)*(originalAspectRatio-newAspectRatio))/2.0)
			fromX = src.Bounds().Max.X - toX
		} else {
			toY = src.Bounds().Max.Y - int((float64(src.Bounds().Max.Y)/(newAspectRatio-originalAspectRatio))/2.0)
			fromY = src.Bounds().Max.Y - toY
		}
	}

	fromSrcRect := image.Rect(fromX, fromY, toX, toY)

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, fromSrcRect, draw.Src, nil)
	writeTo := new(bytes.Buffer)
	err = jpeg.Encode(writeTo, dst, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}

	return writeTo, nil
}

const float64EqualityThreshold = 0.0001

func isAlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
