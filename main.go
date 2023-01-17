package main

import (
	"fmt"
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
const (
	// cell types
	moveStraightCell int	= 0
	wallCell int			= 1
	// Tile types
	startTile int			= 10
	exitTile int			= 11
)

// game states
const (
	preparation int	= 0
	playing int		= 1
)
var gameState int = preparation

// map borders in tile witdths
var mapWidth int = 0
var mapHeight int = 0

var updatesElapsed int = 0
var isPaused bool = true

// This arrays contain all the cells from the game
var cells []*Cell
var backCells []*BackGroundCell

var pauseButton *Button

func canPlay() bool {
	for _, cell := range cells {
		if cell.isReady == false {return false}
	}
	return true
}

func releaseAllCells() {
	for _, cell := range cells {
		cell.TryPlace()
	}
}

func pauseButtonPressed() {
	isPaused = !isPaused
}

// Called after RunGame is called
func init() {
	loadMapFromFile()

	playButtonImage := LoadImage("./textures/PlayButton.png", -1).image
	pauseButton = &Button{playButtonImage,
		Point{0.0, float64(mapHeight) + 1}, false, pauseButtonPressed}
	bell.Listen("LMB_pressed", pauseButton.PressDetect)
}

func checkWinCondition() {
	for _, cell := range cells {
		for _, tile := range backCells {
			if (*tile).cellType != exitTile {continue}
			if	math.Round((*cell).position.x) == math.Round((*tile).position.x) &&
				math.Round((*cell).position.y) == math.Round((*tile).position.y) {
				fmt.Println("Win condition triggered")
				//os.Exit(0)
			}
		}
	}
}

func handleCells() {

	if gameState == playing {
		if updatesElapsed < updatesPerCall {
			updatesElapsed++
			return
		}
		updatesElapsed = 0
	}

	for _, cell := range cells {
		switch gameState {
		case preparation:
			if (*cell).isGrabbed {
				x, y := ebiten.CursorPosition()
				cursorX := float64(x) / 16
				cursorY := float64(y) / 16

				(*cell).position.x = cursorX - 0.5
				(*cell).position.y = cursorY - 0.5
			}
		case playing:
			if isPaused == true { return }

			switch (*cell).cellType {
			case moveStraightCell:
				(*cell).moveForwardOne(cells)
			}
		}
	}
}

func handleInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPos := Point{float64(x), float64(y)}
		bell.Ring("LMB_pressed", cursorPos)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPos := Point{float64(x), float64(y)}
		bell.Ring("LMB_released", cursorPos)
	}
}

// Called every frame
func (g *Game) Update() error {

	checkWinCondition()
	handleInput()
	handleCells()

	if gameState != playing && canPlay() {
		releaseAllCells()
		gameState = playing
		pauseButton.isActive = true
	}

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
