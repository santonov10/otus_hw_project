package usecase

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	resizeOrigPath      = testResizeDir + "/" + "_gopher_original_1024x504.jpg"
	compareResizeResult = []resizeTestImages{
		{
			width:   50,
			height:  50,
			imgPath: testResizeDir + "/" + "gopher_50x50.jpg",
		},
		{
			width:   200,
			height:  700,
			imgPath: testResizeDir + "/" + "gopher_200x700.jpg",
		},
		{
			width:   256,
			height:  126,
			imgPath: testResizeDir + "/" + "gopher_256x126.jpg",
		},
		{
			width:   333,
			height:  666,
			imgPath: testResizeDir + "/" + "gopher_333x666.jpg",
		},
		{
			width:   500,
			height:  500,
			imgPath: testResizeDir + "/" + "gopher_500x500.jpg",
		},
		{
			width:   1024,
			height:  252,
			imgPath: testResizeDir + "/" + "gopher_1024x252.jpg",
		},
		{
			width:   2000,
			height:  1000,
			imgPath: testResizeDir + "/" + "gopher_2000x1000.jpg",
		},
	}
)

type resizeTestImages struct {
	width, height int
	imgPath       string
}

func TestResizeImage(t *testing.T) {
	t.Run("resizeImages", func(t *testing.T) {
		for _, tc := range compareResizeResult {
			tc := tc
			testName := fmt.Sprintf("resize to width - %d height - %d", tc.width, tc.height)
			t.Run(testName, func(t *testing.T) {
				fileData, err := ioutil.ReadFile(resizeOrigPath)
				require.NoError(t, err)
				origImageBytes := bytes.NewBuffer(fileData)

				resizedBytes, err := resizeImage(origImageBytes, tc.width, tc.height)
				require.NoError(t, err)
				// создать новые эталонные файлы для сравнения.
				// newFileName := testResizeDir + "/" + "gopher_" + strconv.Itoa(tc.width) + "x" + strconv.Itoa(tc.height) + ".jpg"
				// ioutil.WriteFile(newFileName, resizedBytes.Bytes(), 0755)

				compareWithReferenceFile, err := os.ReadFile(tc.imgPath)
				require.NoError(t, err)
				require.Equal(t, resizedBytes.Bytes(), compareWithReferenceFile)
			})
		}
	})
}
