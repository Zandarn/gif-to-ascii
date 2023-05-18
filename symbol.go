package main

import (
	"github.com/gookit/color"
)

type Pixel struct {
	R      uint8
	G      uint8
	B      uint8
	A      uint8
	Symbol byte
}

func (p *Pixel) GetColor() string {
	return color.RGB(p.R, p.G, p.B, true).Sprint(string(p.Symbol))
}
