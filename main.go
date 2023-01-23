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
	dublicationCell int		= 2
	rotationCell int		= 3
	// Tile types
	startTile int			= 10
	exitTile int			= 11
)

// game states
const (
	preparation int	= 0
	playing int		= 1
	win int			= 2
	lose int		= 3
)
var gameState int = preparation

// map borders in tile witdths
var mapWidth int = 0
var mapHeight int = 0

var updatesElapsed int = 0
var isPaused bool = true
var cellsReady int = 0
var grabbedCell *Cell = nil

// This arrays contain all the cells from the game
var cells []*Cell
var tiles []*Tile

var pauseButton *Button

func releaseAllCells() {
	for _, cell := range cells {
		cell.TryPlace()
	}
}

// Called after RunGame is called
func init() {
	loadMapFromFile()

	playButtonImage := LoadImage("./textures/PlayButton.png", -1).image
	pauseButton = &Button{playButtonImage, &Transform{
		Point{0.0, float64(mapHeight) + 1},
		Point{0.0, float64(mapHeight) + 1}, 1.0, 1.0},
		false, pauseButtonPressed}
	bell.Listen("LMB_pressed", pauseButton.PressDetect)
}

func checkWinCondition() {
	if gameState != playing { return }
	for _, cell := range cells {
		for _, tile := range tiles {
			cellT := (*cell).transform
			tileT := (*tile).transform
			if (*tile).cellType != exitTile {continue}
			if	math.Round((*cellT).position.x) == math.Round((*tileT).position.x) &&
				math.Round((*cellT).position.y) == math.Round((*tileT).position.y) {
				fmt.Println("Win condition triggered")
				gameState = win
			}
		}
	}
}

// drags cell with the cursor if cell is grabbed
func dragCell() {
	if grabbedCell == nil { return }
	x, y := ebiten.CursorPosition()
	cursorX := float64(x) / float64((*grabbedCell).sprite.Bounds().Dx())
	cursorY := float64(y) / float64((*grabbedCell).sprite.Bounds().Dy())
	cellT := (*grabbedCell).transform

	(*cellT).position.x = cursorX - 0.5
	(*cellT).position.y = cursorY - 0.5
}

func handleCellsPreparation() {
	dragCell()
}

func releaseIsMovingFlag() {
	for _, cell := range cells {
		m := (*cell).movement
		(*m).isMoving = false
	}
}

func moveCells(k float64){
	for _, cell := range cells {
		m := (*cell).movement
		if cell.cellType == wallCell || m.isMoving == false { continue }
		t := (*cell).transform
		(*t).Lerp((*m).startPos, (*m).target, k)
	}
}

func handleCellsPlaying() {
	if updatesElapsed < updatesPerCall {
		updatesElapsed++
		moveCells(float64(updatesElapsed) / float64(updatesPerCall))
		return
	}
	releaseIsMovingFlag()
	if isPaused == true { return }
	updatesElapsed = 0

	for _, cell := range cells {
		switch (*cell).cellType {
		case moveStraightCell:
			(*cell).moveOne((*cell).direction)
		case dublicationCell:
			(*cell).Dublicate()
		case rotationCell:
			(*cell).Rotate()
		}
	}
	isPaused = true
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
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if grabbedCell != nil {
			(*grabbedCell).SetDirection((*grabbedCell).direction + math.Pi / 2)
		}
	}
}

// Called every frame
func (g *Game) Update() error {
	handleInput()

	switch gameState{

	case preparation:
		handleCellsPreparation()
		if cellsReady == len(cells) {
			pauseButton.isActive = true
		}

	case playing:
		handleCellsPlaying()

	case win:

	case lose:
		
	}
	checkWinCondition()

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
