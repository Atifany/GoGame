package main

import (
	"math"
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

type Cell struct {
	sprite    *MyImage
	position  Point
	direction float64
	cellType  int
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
