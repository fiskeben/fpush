package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"strings"
	"time"
)

func concatenateFiles(files []string) (string, error) {
	if len(files) == 1 {
		return files[0], nil
	}

	images := make([]image.Image, len(files))
	verboseLog("concatenate %d files", len(files))
	for i, f := range files {
		img, err := openImage(f)
		if err != nil {
			return "", err
		}
		images[i] = img
	}
	res := drawImages(images)
	filename := createTempFilename()
	if err := writeImage(res, filename); err != nil {
		return "", err
	}
	return filename, nil
}

func openImage(path string) (image.Image, error) {
	r, err := os.Open(path)
	if err != nil {
		return image.Opaque, err
	}
	i, _, err := image.Decode(r)
	if err != nil {
		return image.Opaque, err
	}
	return i, nil
}

func drawImages(images []image.Image) image.Image {
	if len(images) == 1 {
		return images[0]
	}

	width := 640 * 2
	height := 480 + 480*int(math.Floor(float64((len(images)-1)/2)))

	rect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rect)

	var r image.Rectangle
	for i, src := range images {
		x := (i % 2) * 640
		y := int(math.Floor(float64(i/2))) * 480
		r = image.Rect(x, y, x+640, y+480)
		draw.Draw(dst, r, src, image.ZP, draw.Over)
	}
	return dst
}

func createTempFilename() string {
	t := os.TempDir()
	if !strings.HasSuffix(t, "/") {
		t = t + "/"
	}
	return t + "cam_composite_" + time.Now().Format(time.RFC3339) + ".jpg"
}

func writeImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create temp file %s", err)
	}
	opts := jpeg.Options{Quality: 80}
	if err = jpeg.Encode(f, img, &opts); err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}
	return nil
}
