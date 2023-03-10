package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenWidth int = 1024
const screenHeight int = 768
const SCALE float64 = 2.0

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
		t := (*cell).transform
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate((*t).direction)
		MoveAfterRotation(cell, op)

		screen.DrawImage((*cell).sprite, op)
	}
}

func drawButtons(screen *ebiten.Image) {
	var buttons []*Button
	buttons = append(buttons, pauseButton)
	buttons = append(buttons, replayButton)
	buttons = append(buttons, turnByTurnButton)
	for _, button := range buttons {
		sprite := button.spriteChooser(*button)
		op := &ebiten.DrawImageOptions{}
		buttonT := (*button).transform
		op.GeoM.Translate(
		(*buttonT).position.x * float64(sprite.Bounds().Dx()),
		(*buttonT).position.y * float64(sprite.Bounds().Dy()))
		screen.DrawImage(sprite, op)
	}
}

// converts angle value from radians to degrees
func radToDeg(rad float64) (deg float64){
	deg = rad * 180.0 / math.Pi
	return deg
}

// converts angle value from degrees to radians
func degToRad(deg float64) (rad float64) {
	rad = deg / 180.0 * math.Pi
	return rad
}

// Rotates target point around pivot point by angle degrees
func rotatePointPoint(target Point, pivot Point, angle float64) (Point) {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	target.x -= pivot.x
	target.y -= pivot.y

	var xNew float64 = target.x * cos - target.y * sin
	var yNew float64 = target.x * sin + target.y * cos

	target.x = xNew + pivot.x
	target.y = yNew + pivot.y
	return target
}

/*
	Summary

because op.GeoM.Rotate rotates an image the wrong way this function
is needed to adjust rotation result.
It also moves image to cell's coordinates
*/
func MoveAfterRotation(c *Cell, op *ebiten.DrawImageOptions) {
	t := (*c).transform
	cellWidth := float64((*c).sprite.Bounds().Dx())
	cellHeight := float64((*c).sprite.Bounds().Dy())

	pivot := Point{t.position.x + t.width / 2, (*t).position.y + t.height / 2}
	newPos := rotatePointPoint(t.position, pivot, t.direction)
	op.GeoM.Translate(newPos.x * cellWidth, newPos.y * cellHeight)
}

// Called every frame to draw
func (g *Game) Draw(screen *ebiten.Image) {
	drawTiles(screen)
	drawCells(screen)
	drawButtons(screen)
}

// whatever
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return int(float64(screenWidth) / SCALE), int(float64(screenHeight) / SCALE)
}
