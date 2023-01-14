package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

const screenWidth int = 1024
const screenHeight int = 768
const SCALE int = 1

var sampleImage *ebiten.Image

var cell Cell

var entities []interface{}

// Called after RunGame is called
func init() {
	sampleImage = LoadImage("./Sample-PNG-Image.png")

	cell = Cell{sampleImage, 0, 0, 0.15, 0.15}
	entities = append(entities, cell)
}

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return img
}

// Called every frame
func (g *Game) Update() error {
	cell.x += 0.5
	cell.y += 0.4
	return nil
}

// Called every frame to draw
func (g *Game) Draw(screen *ebiten.Image) {
	for _, entity := range entities {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(entity.(Cell).x), float64(entity.(Cell).y))
		op.GeoM.Scale(entity.(Cell).scaleX, entity.(Cell).scaleY)
		screen.DrawImage(sampleImage, op)
	}
}

// whatever
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth / SCALE, screenHeight / SCALE
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
