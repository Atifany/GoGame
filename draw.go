package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// background tiles
func drawTiles(screen *ebiten.Image) {
	for _, tile := range tiles {
		tileT := (*tile).transform
		tileWidth := float64((*tile).sprite.Bounds().Dx())
		tileHeight := float64((*tile).sprite.Bounds().Dy())
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(math.Round((*tileT).position.x * tileWidth),
			math.Round((*tileT).position.y * tileHeight))

		screen.DrawImage((*tile).sprite, op)
	}
}

// cells
func drawCells(screen *ebiten.Image) {
	for _, cell := range cells {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate((*cell).direction)
		MoveAfterRotation(cell, op)

		screen.DrawImage((*cell).sprite, op)
	}
}

/*
	Summary

because op.GeoM.Rotate rotates an image the wrong way this function
is needed to adjust rotation result.
It also moves image to cell's coordinates
*/
func MoveAfterRotation(c *Cell, op *ebiten.DrawImageOptions) {
	cellWidth := float64((*c).sprite.Bounds().Dx())
	cellHeight := float64((*c).sprite.Bounds().Dy())
	
	cos := math.Cos((*c).direction)
	sin := math.Sin((*c).direction)

	rot := Point{cos - sin, sin + cos}
	if rot.x >= 0 {
		rot.x = 0.0
	}
	if rot.y >= 0 {
		rot.y = 0.0
	}

	cellT := (*c).transform
	op.GeoM.Translate(math.Round(((*cellT).position.x - rot.x) * cellWidth),
		math.Round(((*cellT).position.y - rot.y) * cellHeight))
}

// Called every frame to draw
func (g *Game) Draw(screen *ebiten.Image) {
	drawTiles(screen)
	drawCells(screen)
	
	// UI
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		(*pauseButton).position.x*float64((*pauseButton).sprite.Bounds().Dx()),
		(*pauseButton).position.y*float64((*pauseButton).sprite.Bounds().Dy()))
	
	screen.DrawImage((*pauseButton).sprite, op)
}

// whatever
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return int(float64(screenWidth) / SCALE), int(float64(screenHeight) / SCALE)
}
