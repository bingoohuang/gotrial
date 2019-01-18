package randimg

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/pbnjay/pixfont"
)

// RandUint64 generate a random uint64
func RandUint64() uint64 {
	buf := make([]byte, 8)
	mathrand.Read(buf) // Always succeeds, no need to check error
	return binary.LittleEndian.Uint64(buf)
}

// GenerateRandomImageFile generate image file.
func GenerateRandomImageFile(width, height int, randomText, fileName string, fixedSize int64) {
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	imgbytes, imgSize := GenerateRandomImage(width, height, 20, randomText, filepath.Ext(fileName))
	f.Write(imgbytes)
	if fixedSize > int64(imgSize) {
		// refer to https://stackoverflow.com/questions/16797380/how-to-create-a-10mb-file-filled-with-000000-data-in-golang
		// use f.Truncate to change size of the file
		// If you are using unix, then you can create a sparse file very quickly.
		// A sparse file is filled with zero (ascii NUL) and doesn't actually take up the disk space
		// until it is written to, but reads correctly.
		f.Truncate(fixedSize)
	}
}

// GenerateRandomImage generate a random image with imageFormat (jpg/png) .
// refer: https://onlinejpgtools.com/generate-random-jpg
func GenerateRandomImage(width, height, pixelSize int, randomText, imageFormat string) ([]byte, int) {
	var img draw.Image
	switch imageFormat {
	case "jpg":
		img = image.NewNRGBA(image.Rect(0, 0, width, height))
	default: // png
		img = image.NewRGBA(image.Rect(0, 0, width, height))
	}

	yp := height / pixelSize
	xp := width / pixelSize
	for yi := 0; yi < yp; yi++ {
		for xi := 0; xi < xp; xi++ {
			randomColor := GenerateRandomColor()
			DrawPixelWithRandomColor(img, yi, xi, pixelSize, randomColor)
		}
	}

	if randomText != "" {
		pixfont.DrawString(img, 10, 10, randomText, color.Black)
	}

	var buf bytes.Buffer
	switch imageFormat {
	case "jpg":
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}) // 图像质量值为100，是最好的图像显示
	default: // png
		png.Encode(&buf, img)
	}

	return buf.Bytes(), buf.Len()
}

// DrawPixelWithRandomColor draw pixels on img from yi, xi and randomColor with size of pixelSize x pixelSize
func DrawPixelWithRandomColor(img draw.Image, yi, xi, pixelSize int, randomColor color.Color) {
	ys := yi * pixelSize
	ym := ys + pixelSize
	xs := xi * pixelSize
	xm := xs + pixelSize

	for y := ys; y < ym; y++ {
		for x := xs; x < xm; x++ {
			img.Set(x, y, randomColor)
		}
	}
}

// GenerateRandomColor generate a random color
func GenerateRandomColor() color.Color {
	mathrand.Seed(time.Now().UnixNano())
	return color.RGBA{
		R: uint8(mathrand.Intn(255)),
		G: uint8(mathrand.Intn(255)),
		B: uint8(mathrand.Intn(255)),
		A: uint8(mathrand.Intn(255)),
	}
}
