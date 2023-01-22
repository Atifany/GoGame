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
	hasMoved		bool
}

func (c *Cell) SetDirection(direction float64) {
	if direction >= 2*math.Pi {
		direction -= 2*math.Pi
	}
	(*c).direction = direction
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

func getOppositeDir(direction float64) float64 {
	if direction >= math.Pi {
		return direction - math.Pi
	} else {
		return direction + math.Pi
	}
}

// Moves its tile by one width in the direction pointed by direction
func (c *Cell) moveOne(direction float64) {
	if ((*c).hasMoved == true || (*c).cellType == wallCell) ||
		(c.cellType == moveStraightCell && 
		c.direction == getOppositeDir(direction)) {
	//
		return
	}

	t := (*c).transform
	target := Point{math.Round((*t).position.x + math.Cos(c.direction)),
		math.Round((*t).position.y + math.Sin(c.direction))}
	collision := checkCollisions(cells, target)
	if collision != nil {
		(*collision).moveOne(direction)
	}
	if checkCollisions(cells, target) == nil {
		(*t).position = target
	}
	(*c).hasMoved = true
}

func (c *Cell) Dublicate() {
	if (*c).cellType != dublicationCell { return }

	t := (*c).transform
	targetF := Point{math.Round((*t).position.x + math.Cos(c.direction)),
		math.Round((*t).position.y + math.Sin(c.direction))}
	collision := checkCollisions(cells, targetF)
	if collision != nil {
		(*collision).moveOne(c.direction)
	}
	if checkCollisions(cells, targetF) == nil {
		behindDir := getOppositeDir(c.direction)
		targetB := Point{math.Round((*t).position.x + math.Cos(behindDir)),
			math.Round((*t).position.y + math.Sin(behindDir))}
		collision = checkCollisions(cells, targetB)
		if collision == nil {
			return
		}
		newTransform := &Transform{targetF, targetF, 1.0, 1.0}
		cells = append(cells, &Cell{collision.sprite, newTransform,
			collision.direction, collision.cellType, true})
	}
}

func (c *Cell) Rotate() {
	if (*c).cellType != rotationCell { return }

	t := (*c).transform
	direction := 0.0
	for direction < 2 * math.Pi {
		target := Point{math.Round((*t).position.x + math.Cos(direction)),
			math.Round((*t).position.y + math.Sin(direction))}
		collision := checkCollisions(cells, target)
		direction += math.Pi / 2
		// also add here a check for duplication cell front a back surfaces
		if collision == nil || collision.cellType == wallCell { continue }
			
		(*collision).SetDirection(collision.direction + math.Pi / 2)
		//fmt.Println((*collision).cellType, " ", (*collision).direction / math.Pi * 2)
	}
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
		for _, tile := range tiles {
			if (*tile).transform.isPointInside(Point{pressedX, pressedY}) {
				(*tile).isOccupied= false
				break
			}
		}
		grabbedCell = c
	}
}

func (c *Cell) releaseDetect(message bell.Message){
	(*c).TryPlace()
}

func (c *Cell) TryPlace() {
	if grabbedCell != c { return }

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
				(*tile).isOccupied = true
				grabbedCell = nil
				cellsReady++
				return
			}
		}
	}
	(*cellT).position = (*cellT).defaultPosition
	grabbedCell = nil
}
