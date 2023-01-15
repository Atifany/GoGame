package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MyImage struct {
	image *ebiten.Image
	size  Point
}

type Game struct{}

const screenWidth int = 1024
const screenHeight int = 768
const SCALE float64 = 2.0
const updatesPerCall int = 30

const moveStraightCell int = 0
const wallCell int = 1

var updatesElapsed int = 0

var cells []*Cell

// Called after RunGame is called
func init() {
	redTileImage := LoadImage("./textures/RedTile.png")
	blackTileImage := LoadImage("./textures/BlackTile.png")

	cells = append(cells, &Cell{redTileImage,
		Point{1, 1}, 0.0, moveStraightCell})

	cells = append(cells, &Cell{redTileImage,
		Point{2, 2}, math.Pi / 2, moveStraightCell})

	cells = append(cells, &Cell{blackTileImage,
		Point{2, 3}, 0.0, wallCell})

	cells = append(cells, &Cell{blackTileImage,
		Point{4, 1}, 0.0, wallCell})

	cells = append(cells, &Cell{redTileImage,
		Point{6, 6}, 0.0, moveStraightCell})

	cells = append(cells, &Cell{redTileImage,
		Point{8, 6}, math.Pi, moveStraightCell})

	cells = append(cells, &Cell{redTileImage,
		Point{6, 7}, 0.0, moveStraightCell})

	cells = append(cells, &Cell{redTileImage,
		Point{9, 7}, math.Pi, moveStraightCell})

	cells = append(cells, &Cell{redTileImage,
		Point{6, 8}, 0.0, moveStraightCell})

	cells = append(cells, &Cell{redTileImage,
		Point{7, 8}, 0.0, moveStraightCell})
}

// Loads png file into ebiten.Image struct
func LoadImage(path string) *MyImage {
	img, i, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return &MyImage{img, Point{float64(i.Bounds().Max.X), float64(i.Bounds().Max.Y)}}
}

func ProcedeCells() {
	for _, cell := range cells {
		switch (*cell).cellType {
		case moveStraightCell:
			(*cell).moveForwardOne(cells)
		}
	}
}

// Called every frame
func (g *Game) Update() error {
	if updatesElapsed < updatesPerCall {
		updatesElapsed++
		return nil
	}
	updatesElapsed = 0

	ProcedeCells()

	return nil
}

/*
	Summary

because op.GeoM.Rotate rotates an image the wrong way this function
is needed to adjust rotation result.
It also moves image to cell's coordinates
*/
func adjustAfterRotation(c *Cell, op *ebiten.DrawImageOptions) {
	cos := math.Cos((*c).direction)
	sin := math.Sin((*c).direction)

	rot := Point{cos - sin, sin + cos}
	if rot.x >= 0 {
		rot.x = 0.0
	}
	if rot.y >= 0 {
		rot.y = 0.0
	}

	op.GeoM.Translate(math.Round(((*c).position.x-rot.x)*(*c).sprite.size.x),
		math.Round(((*c).position.y-rot.y)*(*c).sprite.size.y))
}

// Called every frame to draw
func (g *Game) Draw(screen *ebiten.Image) {
	for _, cell := range cells {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate((*cell).direction)
		adjustAfterRotation(cell, op)

		screen.DrawImage((*cell).sprite.image, op)
	}
}

// whatever
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return int(float64(screenWidth) / SCALE), int(float64(screenHeight) / SCALE)
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
