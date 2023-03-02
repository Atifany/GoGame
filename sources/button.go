package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nuttech/bell"
)

type Button struct {
	sprites			[]*Sprite
	transform		*Transform
	isActive		bool
	pressed			func()
	spriteChooser	func(Button) *ebiten.Image
}

type Sprite struct {
	sprite	*ebiten.Image
	name	string
}

func (b Button) findSpriteByName(name string) (*ebiten.Image) {
	for _, sprite := range b.sprites {
		if sprite.name == name {
			return sprite.sprite
		}
	}
	return nil
}

func playButtonSpriteChooser(b Button) (*ebiten.Image) {
	if b.isActive == false {
		return b.findSpriteByName("inactive")
	}
	if isPaused == true {
		return b.findSpriteByName("pause")
	}
	if isPaused == false {
		return b.findSpriteByName("play")
	}
	return b.findSpriteByName("play")
}

func replayButtonSpriteChooser(b Button) (*ebiten.Image) {
	if b.isActive == false {
		return b.findSpriteByName("inactive")
	}
	if b.isActive == true {
		return b.findSpriteByName("replay")
	}
	return b.findSpriteByName("replay")
}

func turnByTurnButtonSpriteChooser(b Button) (*ebiten.Image) {
	if b.isActive == false {
		return b.findSpriteByName("inactive")
	}
	if b.isActive == true {
		return b.findSpriteByName("turnByTurn")
	}
	return b.findSpriteByName("turnByTurn")
}

func pauseButtonPressed() {
	isPaused = !isPaused
	if gameState == preparation {
		gameState = playing
		releaseAllCells()
		releaseAllTiles()
	}
}

func removeFromCells(index int){
	cells[index] = cells[len(cells) - 1]
}

func removeSummonedCells() {
	var newCells []*Cell
	for _, cell := range cells {
		if cell.isSummoned == true { continue }
		t := (*cell).transform
		m := (*cell).movement
		(*m).startRot = 0.0
		(*m).endRot = 0.0
		(*t).direction = 0.0
		newCells = append(newCells, cell)
	}
	cells = newCells
}

func restartLevel() {
	gameState = preparation
	pauseButton.isActive = false
	isPaused = true
	isTurnByTurn = false
	cellsReady = 0
	removeSummonedCells()
	releaseAllCells()
	releaseAllTiles()
	placeCellsAtStart()
}

func switchTurnByTurn() {
	isTurnByTurn = !isTurnByTurn
	if gameState == playing {
		isPaused = !isPaused
	}
}

// Clear all consts properly
// PressDetect is called on LMB clicked and detects whether
// a click landed on a parent button.
func (b *Button) PressDetect(message bell.Message) {
	if (*b).isActive == false { return }

	pressedX := message.Value.(Point).x / 16
	pressedY := message.Value.(Point).y / 16
	if (*b).transform.isPointInside(Point{pressedX, pressedY}) == true {
		(*b).pressed()
	}
}
