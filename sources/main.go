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
var gameState int		= preparation
var isPaused bool		= true
var isTurnByTurn bool	= false

// map borders in tile witdths
var mapWidth int = 0
var mapHeight int = 0

var updatesElapsed int = 0
var cellsReady int = 0
var grabbedCell *Cell = nil

// This arrays contain all the cells from the game
var cells	[]*Cell
var tiles	[]*Tile

var pauseButton			*Button
var replayButton		*Button
var turnByTurnButton	*Button

func releaseAllCells() {
	for _, cell := range cells {
		cell.TryPlace()
	}
}

func releaseAllTiles() {
	for _, tile := range tiles {
		(*tile).isOccupied = false
	}
}

func initButtons() {
	var sprites []*Sprite

	playButtonImage := LoadImage("./textures/Buttons/PlayButton.png", -1).image
	pauseButtonImage := LoadImage("./textures/Buttons/PausePlayButton.png", -1).image
	inactiveButtonImage := LoadImage("./textures/Buttons/InactivePlayButton.png", -1).image
	sprites = append(sprites, &Sprite{ playButtonImage, "play" })
	sprites = append(sprites, &Sprite{ pauseButtonImage, "pause" })
	sprites = append(sprites, &Sprite{ inactiveButtonImage, "inactive" })
	pauseButton = &Button{ sprites,
		&Transform{ Point{0.0, float64(mapHeight) + 1},
		Point{0.0, float64(mapHeight) + 1}, 0.0, 1.0, 1.0},
		false, pauseButtonPressed, playButtonSpriteChooser}
	bell.Listen("LMB_pressed", pauseButton.PressDetect)
	sprites = nil

	replayButtonImage := LoadImage("./textures/Buttons/ReplayButton.png", -1).image
	inactiveReplayButtonImage := LoadImage("./textures/Buttons/InactiveReplayButton.png", -1).image
	sprites = append(sprites, &Sprite{ replayButtonImage, "replay" })
	sprites = append(sprites, &Sprite{ inactiveReplayButtonImage, "inactive" })
	replayButton = &Button{sprites, &Transform{
		Point{1.0, float64(mapHeight) + 1},
		Point{1.0, float64(mapHeight) + 1}, 0.0, 1.0, 1.0},
		true, restartLevel, replayButtonSpriteChooser}
	bell.Listen("LMB_pressed", replayButton.PressDetect)
	sprites = nil

	turnByTurnButtonImage := LoadImage("./textures/Buttons/TurnByTurnButton.png", -1).image
	inactiveTurnByTurnButtonImage := LoadImage("./textures/Buttons/InactiveTurnByTurnButton.png", -1).image
	sprites = append(sprites, &Sprite{ turnByTurnButtonImage, "turnByTurn" })
	sprites = append(sprites, &Sprite{ inactiveTurnByTurnButtonImage, "inactive"})
	turnByTurnButton = &Button{sprites, &Transform{
		Point{2.0, float64(mapHeight) + 1},
		Point{2.0, float64(mapHeight) + 1}, 0.0, 1.0, 1.0},
		false, switchTurnByTurn, turnByTurnButtonSpriteChooser}
	bell.Listen("LMB_pressed", turnByTurnButton.PressDetect)
	sprites = nil
}

// Called after RunGame is called
func init() {
	loadMapFromFile()
	initButtons()
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

func fixPositions() {
	for _, cell := range cells {
		m := (*cell).movement
		t := (*cell).transform
		(*t).position = m.target
		(*m).startPos = m.target
		SetRotation(&((*t).direction), m.endRot)
		SetRotation(&((*m).startRot), m.endRot)
		SetRotation(&((*m).endRot), m.endRot)
	}
}

func moveCells(k float64){
	for _, cell := range cells {
		m := (*cell).movement
		t := (*cell).transform
		(*t).LerpRotate(m.startRot, m.endRot, k)
		if cell.cellType == wallCell { continue }
		(*t).Lerp((*m).startPos, (*m).target, k)
	}
}

func handleCellsPlaying() {
	if updatesElapsed < updatesPerCall {
		updatesElapsed++
		moveCells(float64(updatesElapsed) / float64(updatesPerCall))
		return
	}
	fixPositions()
	if isPaused == true { return }
	updatesElapsed = 0

	for _, cell := range cells {
		if (*cell).cellType != dublicationCell { continue }
		(*cell).Dublicate()
	}
	for _, cell := range cells {
		if (*cell).cellType != rotationCell { continue }
		(*cell).Rotate()
	}
	for _, cell := range cells {
		if (*cell).cellType != moveStraightCell { continue }
		(*cell).moveOne(cell.transform.direction)
	}

	if isTurnByTurn == true {
		isPaused = true
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
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if grabbedCell != nil {
			t := (*grabbedCell).transform
			m := (*grabbedCell).movement
			SetRotation(&((*t).direction), t.direction + math.Pi / 2)
			SetRotation(&((*m).startRot), m.startRot + math.Pi / 2)
			SetRotation(&((*m).endRot), m.endRot + math.Pi / 2)
		}
	}
}

// Called every frame
func (g *Game) Update() error {
	handleInput()

	switch gameState{

	case preparation:
		handleCellsPreparation()
		turnByTurnButton.isActive = false
		if cellsReady == countNonWallCells() {
			pauseButton.isActive = true
		}

	case playing:
		turnByTurnButton.isActive = true
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
