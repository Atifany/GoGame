package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nuttech/bell"
)

type Button struct {
	sprite		*ebiten.Image
	transform	*Transform
	isActive	bool
	pressed		func()
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
		(*cell).direction = 0.0
		newCells = append(newCells, cell)
	}
	cells = newCells
}

func restartLevel() {
	gameState = preparation
	removeSummonedCells()
	placeCellsAtStart()
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
