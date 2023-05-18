package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

type Gif struct {
	Height int
	Width  int

	FrameIndex int
	Frames     [][]Pixel
}

//go:embed test1.gif
var imageBytes []byte
var frameStringBuffer strings.Builder
var asciiChars = []byte{' ', '.', ',', ':', ';', 'i', 't', 'L', 'C', 'o', 'Y', 'X', 'Z', 'N', 'M', 'W', 'Q', 'K', 'D', '#', '*'}

func New() (g *Gif) {
	g = new(Gif)
	g.Width = 80 // because cli

	file := bytes.NewReader(imageBytes)

	err := g.framing(file)
	if err != nil {
		panic(err)
	}

	return g
}

func (g *Gif) framing(r io.Reader) (err error) {
	gifDecoded, _ := gif.DecodeAll(r)

	g.Height = int(float64(gifDecoded.Image[0].Bounds().Max.Y) * float64(g.Width) / (float64(gifDecoded.Image[0].Bounds().Max.X) * 1.6))
	frame := image.NewRGBA(image.Rect(0, 0, gifDecoded.Config.Width, gifDecoded.Config.Height))
	for _, img := range gifDecoded.Image {
		draw.Draw(frame, frame.Bounds(), img, image.Point{}, draw.Over)
		g.Frames = append(g.Frames, g.convertToAscii(frame))
	}

	return
}

func (g *Gif) convertToAscii(img image.Image) []Pixel {
	img = resize.Resize(uint(g.Width), uint(g.Height), img, resize.Lanczos3)
	buffer := make([]Pixel, 0, g.Height*g.Width)

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			i := int((brightness / 65535) * float64(len(asciiChars)-1))

			p := Pixel{
				Symbol: asciiChars[i],
				R:      uint8(r >> 8),
				G:      uint8(g >> 8),
				B:      uint8(b >> 8),
				A:      uint8(a >> 8),
			}

			buffer = append(buffer, p)
		}

		buffer = append(buffer, Pixel{Symbol: '\n'})
	}

	return buffer
}

func (g *Gif) getNextFrame() string {
	frame := g.Frames[g.FrameIndex]
	g.FrameIndex = (g.FrameIndex + 1) % len(g.Frames)
	frameStringBuffer.Reset()
	for i := 0; i < len(frame); i++ {
		frameStringBuffer.WriteString(frame[i].GetColor())
	}

	return frameStringBuffer.String()
}

func (g *Gif) Print() {
	timeout := time.Duration(Timeout) * time.Second
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	endTime := time.Now().Add(timeout)

	for ; time.Now().Before(endTime); <-ticker.C {
		fmt.Print("\033[0;0H", g.getNextFrame())
	}
}
