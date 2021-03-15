package randimg

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
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
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// GenerateRandomImageFile generate image file.
// If fastMode is true, a sparse file is filled with zero (ascii NUL) and doesn't actually take up the disk space
// until it is written to, but reads correctly.
// $ ls -lh 424661641.png
// -rw-------  1 bingoobjca  staff   488K Mar 15 12:19 424661641.png
// $ du -hs 424661641.png
// 8.0K    424661641.png
// If fastMode is false, an actually sized file will generated.
// $ ls -lh 1563611881.png
// -rw-------  1 bingoobjca  staff   488K Mar 15 12:28 1563611881.png
// $ du -hs 1563611881.png
// 492K    1563611881.png

type RandImageConfig struct {
	Width      int
	Height     int
	RandomText string
	FileName   string
	FixedSize  int64
	FastMode   bool
	PixelSize  int
}

func (c *RandImageConfig) GenerateFile() {
	f, _ := os.OpenFile(c.FileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	if c.PixelSize == 0 {
		c.PixelSize = 50
	}

	imgbytes, imgSize := c.GenerateImage(filepath.Ext(c.FileName))
	f.Write(imgbytes)
	if c.FixedSize <= int64(imgSize) {
		return
	}

	if !c.FastMode {
		b, _ := GenerateRandomBytes(int(c.FixedSize) - imgSize)
		f.Write(b)
		return
	}

	// refer to https://stackoverflow.com/questions/16797380/how-to-create-a-10mb-file-filled-with-000000-data-in-golang
	// use f.Truncate to change size of the file
	// If you are using unix, then you can create a sparse file very quickly.
	// A sparse file is filled with zero (ascii NUL) and doesn't actually take up the disk space
	// until it is written to, but reads correctly.
	f.Truncate(c.FixedSize)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.RawURLEncoding.EncodeToString(b), err
}

// GenerateImage generate a random image with imageFormat (jpg/png) .
// refer: https://onlinejpgtools.com/generate-random-jpg
func (c *RandImageConfig) GenerateImage(imageFormat string) ([]byte, int) {
	var img draw.Image

	switch imageFormat {
	case "jpg":
		img = image.NewNRGBA(image.Rect(0, 0, c.Width, c.Height))
	default: // png
		img = image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	}

	yp := c.Height / c.PixelSize
	xp := c.Width / c.PixelSize
	for yi := 0; yi < yp; yi++ {
		for xi := 0; xi < xp; xi++ {
			randomColor := GenerateRandomColor()
			DrawPixelWithRandomColor(img, yi, xi, c.PixelSize, randomColor)
		}
	}

	if c.RandomText != "" {
		pixfont.DrawString(img, 10, 10, c.RandomText, color.Black)
	}

	var buf bytes.Buffer
	switch imageFormat {
	case "jpg", "jpeg":
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
