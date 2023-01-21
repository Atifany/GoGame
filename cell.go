package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nuttech/bell"
)

type Point struct {
	x float64
	y float64
}

type Tile struct {
	sprite		*ebiten.Image
	transform	*Transform
	cellType	int
	isOccupied	bool
}

// Cells could not be grabbed when canPlay() triggers

type Cell struct {
	sprite			*ebiten.Image
	transform		*Transform
	direction		float64
	cellType		int
	isGrabbed		bool
}

type Transform struct {
	position		Point
	defaultPosition Point
	width			float64
	height			float64
}

// Returns a pointer to a tile which occupies targetPos. Returns nil of none were found
func checkCollisions(cells []*Cell, targetPos Point) *Cell {
	for _, cell := range cells {
		t := (*cell).transform
		if (*t).position.x == targetPos.x && (*t).position.y == targetPos.y {
			return cell
		}
	}
	return nil
}

// Moves its tile by one width in the direction pointed by Cell.direction
func (c *Cell) moveForwardOne(cells []*Cell) {
	t := (*c).transform
	target := Point{math.Round((*t).position.x + math.Cos(c.direction)),
		math.Round((*t).position.y + math.Sin(c.direction))}
	if checkCollisions(cells, target) != nil {
		return
	}

	(*t).position = target
}

func (t *Transform) isPointInside(point Point) bool {
	if point.x < (*t).position.x + (*t).width &&
		point.x > (*t).position.x &&
		point.y < (*t).position.y + (*t).height &&
		point.y > (*t).position.y {
	//
		return true
	}
	return false
}

// PressDetect is called on LMB clicked and detects whether
// a click landed on a parent button.
func (c *Cell) PressDetect(message bell.Message) {
	if gameState != preparation { return }
	pressedX := message.Value.(Point).x / float64((*c).sprite.Bounds().Dx())
	pressedY := message.Value.(Point).y / float64((*c).sprite.Bounds().Dy())

	if (*c).transform.isPointInside(Point{pressedX, pressedY}) {
		(*c).isGrabbed = true
	}
}

func (c *Cell) releaseDetect(message bell.Message){
	(*c).TryPlace()
}

func (c *Cell) TryPlace() {
	if (*c).isGrabbed == false { return }

	cellT := (*c).transform
	releasedX := (*cellT).position.x + 0.5
	releasedY := (*cellT).position.y + 0.5
	
	for _, tile := range tiles {
		tileT := (*tile).transform
		if tile.cellType == startTile &&
			(*tile).transform.isPointInside(Point{releasedX, releasedY}) {
		// 
			if (*tile).isOccupied == true {
				break
			} else {
				(*cellT).position = (*tileT).position
				(*cellT).defaultPosition = (*tileT).position
				(*c).isGrabbed = false
				cellsReady++
				return
			}
		}
	}
	(*cellT).position = (*cellT).defaultPosition
	(*c).isGrabbed = false
}
