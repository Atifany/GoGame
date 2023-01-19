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
var tiles []*Tile

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
		for _, tile := range tiles {
			cellT := (*cell).transform
			tileT := (*tile).transform
			if (*tile).cellType != exitTile {continue}
			if	math.Round((*cellT).position.x) == math.Round((*tileT).position.x) &&
				math.Round((*cellT).position.y) == math.Round((*tileT).position.y) {
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
				cursorX := float64(x) / float64((*cell).sprite.Bounds().Dx())
				cursorY := float64(y) / float64((*cell).sprite.Bounds().Dy())
				cellT := (*cell).transform

				(*cellT).position.x = cursorX - 0.5
				(*cellT).position.y = cursorY - 0.5
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

func main() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GoGame")

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err.Error())
	}
}
