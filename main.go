package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nuttech/bell"
)

type Game struct{}

const screenWidth int = 1024
const screenHeight int = 768
const SCALE float64 = 2.0
const updatesPerCall int = 30

// cells - those who iteract with each other.
// tiles - those who stay on a background.
// cell types
const moveStraightCell int = 0
const wallCell int = 1

// Tile types
const startTile int = 10
const exitTile int = 11

// map borders in tile witdths
var mapWidth int = 0
var mapHeight int = 0

var updatesElapsed int = 0
var isPaused bool = true

// This arrays contain all the cells from the game
var cells []*Cell
var backCells []*BackGroundCell

var pauseButton *Button

func pauseButtonPressed() {
	isPaused = !isPaused
}

// Called after RunGame is called
func init() {
	loadMapFromFile()
	playButtonImage := LoadImage("./textures/PlayButton.png", -1).image
	pauseButton = &Button{playButtonImage,
		Point{0.0, float64(mapHeight) + 1}, true, pauseButtonPressed}
	bell.Listen("LMB_pressed", pauseButton.PressDetect)
}

func handleCells() {
	if isPaused == true {
		return
	}
	for _, cell := range cells {
		switch (*cell).cellType {
		case moveStraightCell:
			(*cell).moveForwardOne(cells)
		}
	}
}

func handleInput() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	// fmt.Printf("lol")
	x, y := ebiten.CursorPosition()
	cursorPos := Point{float64(x), float64(y)}
	bell.Ring("LMB_pressed", cursorPos)
}

// Called every frame
func (g *Game) Update() error {

	handleInput()

	if updatesElapsed < updatesPerCall {
		updatesElapsed++
		return nil
	}
	updatesElapsed = 0

	handleCells()

	return nil
}

/*
	Summary

because op.GeoM.Rotate rotates an image the wrong way this function
is needed to adjust rotation result.
It also moves image to cell's coordinates
*/
func adjustAfterRotation(c *Cell, op *ebiten.DrawImageOptions) {
	cos := math.Cos((*c).direction)
	sin := math.Sin((*c).direction)

	rot := Point{cos - sin, sin + cos}
	if rot.x >= 0 {
		rot.x = 0.0
	}
	if rot.y >= 0 {
		rot.y = 0.0
	}

	op.GeoM.Translate(math.Round(((*c).position.x-rot.x)*(*c).sprite.size.x),
		math.Round(((*c).position.y-rot.y)*(*c).sprite.size.y))
}

// Called every frame to draw
func (g *Game) Draw(screen *ebiten.Image) {
	// background tiles
	for _, cell := range backCells {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate((*cell).position.x*(*cell).sprite.size.x,
			(*cell).position.y*(*cell).sprite.size.y)

		screen.DrawImage((*cell).sprite.image, op)
	}

	// cells
	for _, cell := range cells {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate((*cell).direction)
		adjustAfterRotation(cell, op)

		screen.DrawImage((*cell).sprite.image, op)
	}

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

func main() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GoGame")

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err.Error())
	}
}
