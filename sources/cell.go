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
	sprite		*ebiten.Image
	transform	*Transform
	cellType	int
	isSummoned	bool
	movement	*MovementPlaceHolder
}

func SetRotation(dst *float64, newValue float64) {
	if newValue >= 2 * math.Pi {
		newValue -= 2 * math.Pi
	}
	*dst = newValue
}

type Transform struct {
	position		Point
	defaultPosition Point
	direction		float64
	width			float64
	height			float64
}

type MovementPlaceHolder struct {
	startPos	Point
	target		Point
	startRot	float64
	endRot		float64
}

func (t *Transform) Lerp(startPos Point, endPos Point, k float64) {
	shiftX := endPos.x - startPos.x
	shiftY := endPos.y - startPos.y

	(*t).position.x = startPos.x + k * shiftX
	(*t).position.y = startPos.y + k * shiftY
}

func (t *Transform) LerpRotate(startRot float64, endRot float64, k float64) {
	shift := endRot - startRot
	
	(*t).direction = startRot + k * shift
	//SetRotation(&((*t).direction), startRot + k * shift)
}

// Returns a pointer to a tile which occupies targetPos. Returns nil of none were found
func checkCollisions(cells []*Cell, targetPos Point) *Cell {
	for _, cell := range cells {
		m := (*cell).movement
		if (*m).target == targetPos {
		//
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
	if (*c).cellType == wallCell { return }
	
	m := (*c).movement
	t := (*c).transform
	
	target := Point{math.Round((*m).target.x + math.Cos(direction)),
		math.Round((*m).target.y + math.Sin(direction))}
	collision := checkCollisions(cells, target)
	if collision != nil {
		if collision.cellType == moveStraightCell &&
			(*t).direction == getOppositeDir((*t).direction) { return }
		(*collision).moveOne(direction)
	}

	collision = checkCollisions(cells, target)
	if collision == nil {
		(*m).target = target
	}
}

func (c *Cell) Dublicate() {
	if (*c).cellType != dublicationCell { return }

	m := (*c).movement
	t := (*c).transform
	
	behindDir := getOppositeDir(t.direction)
	targetB := Point{math.Round((*m).target.x + math.Cos(behindDir)),
		math.Round((*m).target.y + math.Sin(behindDir))}
	collisionB := checkCollisions(cells, targetB)
	if collisionB == nil { return }

	targetF := Point{math.Round((*m).target.x + math.Cos(t.direction)),
		math.Round((*m).target.y + math.Sin(t.direction))}
	collisionF := checkCollisions(cells, targetF)
	if collisionF != nil {
		(*collisionF).moveOne(t.direction)
	}
	if checkCollisions(cells, targetF) != nil { return }
	newTransform := &Transform{targetF, targetF, collisionB.transform.direction, 1.0, 1.0}
	cells = append(cells, &Cell{collisionB.sprite, newTransform,
		collisionB.cellType, true,
		&MovementPlaceHolder{newTransform.position, newTransform.position,
		collisionB.transform.direction, collisionB.transform.direction}})
}

func (c *Cell) Rotate() {
	if (*c).cellType != rotationCell { return }

	// t := (*c).transform
	m := (*c).movement
	direction := 0.0
	for direction < 2 * math.Pi {
		target := Point{math.Round((*m).target.x + math.Cos(direction)),
			math.Round((*m).target.y + math.Sin(direction))}
		collision := checkCollisions(cells, target)
		direction += math.Pi / 2
		if collision == nil { continue }
		collisionM := (*collision).movement
		(*collisionM).endRot = collisionM.endRot - math.Pi / 2
	}
	(*m).endRot = m.endRot - math.Pi / 2
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
				(*tile).isOccupied = false
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
				m := (*c).movement
				(*m).startPos = (*tileT).position
				(*m).target = (*tileT).position
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

func countNonWallCells() (int) {
	res := 0
	for _, cell := range cells {
		if cell.cellType == wallCell { continue }
		res++
	}
	return res
}
