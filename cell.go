package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Cell struct {
	sprite *ebiten.Image
	x      float64
	y      float64
	scaleX float64
	scaleY float64
}
