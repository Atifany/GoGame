package main

import (
	"math"

	"github.com/nuttech/bell"
)

type Point struct {
	x float64
	y float64
}

type BackGroundCell struct {
	sprite *MyImage
	position Point
	cellType int
}

// Cells could not be grabbed when canPlay() triggers

type Cell struct {
	sprite			*MyImage
	position		Point
	defaultPosition Point
	direction		float64
	cellType		int
	isGrabbed		bool
	isReady			bool
}

// Returns a pointer to a tile which occupies targetPos. Returns nil of none were found
func checkCollisions(cells []*Cell, targetPos Point) *Cell {
	for _, cell := range cells {
		if (*cell).position.x == targetPos.x && (*cell).position.y == targetPos.y {
			return cell
		}
	}
	return nil
}

// Moves its tile by one width in the direction pointed by Cell.direction
func (c *Cell) moveForwardOne(cells []*Cell) {
	target := Point{math.Round(c.position.x + 1*math.Cos(c.direction)),
		math.Round(c.position.y + 1*math.Sin(c.direction))}
	if checkCollisions(cells, target) != nil {
		return
	}

	c.position = target
}

// Clear all consts properly
// PressDetect is called on LMB clicked and detects whether
// a click landed on a parent button.
func (c *Cell) PressDetect(message bell.Message) {
	pressedX := message.Value.(Point).x / 16
	pressedY := message.Value.(Point).y / 16
	cWidth := 1.0
	cHegith := 1.0
	cPosX := (*c).position.x + cWidth / 2
	cPosY := (*c).position.y + cHegith / 2

	if pressedX < cPosX + cWidth / 2 &&
		pressedX > cPosX - cWidth / 2 &&
		pressedY < cPosY + cHegith / 2 &&
		pressedY >  cPosY - cHegith / 2 {
		(*c).isGrabbed = true
	}
}

func (c *Cell) releaseDetect(message bell.Message){
	(*c).TryPlace()
}

func (c *Cell) TryPlace() {
	if (*c).isGrabbed == false { return }

	releasedX := (*c).position.x + 0.5
	releasedY := (*c).position.y + 0.5
	cWidth := 1.0
	cHegith := 1.0
	
	for _, tile := range backCells {
		cPosX := (*tile).position.x + cWidth / 2
		cPosY := (*tile).position.y + cHegith / 2
		if tile.cellType == startTile &&
			releasedX <= cPosX + cWidth / 2 &&
			releasedX >= cPosX - cWidth / 2 &&
			releasedY <= cPosY + cHegith / 2 &&
			releasedY >= cPosY - cHegith / 2 {
			tPosX := (*tile).position.x
			tPosY := (*tile).position.y
			for _, cell := range cells {
				if	math.Round(tPosX) == math.Round((*cell).position.x) &&
					math.Round(tPosY) == math.Round((*cell).position.y) &&
					(*cell) != (*c){
					(*c).position = (*c).defaultPosition
					(*c).isGrabbed = false
					return
				}
			}
			(*c).position = tile.position
			(*c).defaultPosition = tile.position
			(*c).isReady = true
			(*c).isGrabbed = false
			return
		}
	}

	(*c).position = (*c).defaultPosition
	(*c).isGrabbed = false
}
