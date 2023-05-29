package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// background represents a gnome background gif based.
type background struct {
	// frames holds paths of gif frames.
	frames []string

	// delay indicates the delay during frame transitioning.
	delay int
}

// anime completes a gif frame cycle.
func (b *background) animate() {
	for _, uri := range b.frames {
		fmt.Println(uri)

		cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", uri)
		cmd.Run()

		time.Sleep(time.Duration(b.delay*30) * time.Millisecond)
	}
}

// newBackground creates a gif based background by given gif file path.
func newBackground(gifPath string) (*background, error) {
	d, err := os.ReadFile(gifPath)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(d)

	gif, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}

	var frames []string

	w, h := calcBackgroundDim(gif)
	canvas := image.NewRGBA(image.Rect(0, 0, w, h))

	for i, srcImg := range gif.Image {
		draw.Draw(canvas, canvas.Bounds(), srcImg, image.ZP, draw.Over)

		abs, err := filepath.Abs(".")
		if err != nil {
			return nil, err
		}

		frame := fmt.Sprintf("%s/frame%d", abs, i)

		file, err := os.Create(frame)
		if err != nil {
			return nil, err
		}

		err = png.Encode(file, canvas)
		if err != nil {
			return nil, err
		}

		file.Close()

		frames = append(frames, frame)
	}

	return &background{
		frames: frames,
		delay:  gif.Delay[0],
	}, nil
}

// calcBackgroundDim calculates the dimensions of background.
func calcBackgroundDim(gif *gif.GIF) (x, y int) {
	var (
		lowX  int
		lowY  int
		highX int
		highY int
	)

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowX {
			lowX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowY {
			lowY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highX {
			highX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highY {
			highY = img.Rect.Max.Y
		}
	}

	return highX - lowX, highY - lowY
}
