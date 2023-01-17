package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// map tiles types
const wallSymbol byte = 'X'
const startSymbol byte = 'V'
const exitSymbol byte = 'E'
const forwardMovingName string = "ForwardMoving"

type MyImage struct {
	image    *ebiten.Image
	size     Point
	cellType int
}

// readAvailableTiles parses given available cells to cells
func readAvailableTiles(gameMap string, cellImages []*MyImage) {
	var col int = mapWidth + 1
	var row int = 0
	mapLines := strings.Split(gameMap, "\n")

	var i int = 0
	for i < len(mapLines) && mapLines[i] != "Available tiles:" {
		i++
	}
	i++ // skip "Available tiles:" line

	for i < len(mapLines) {
		if mapLines[i] == "" {
			i++
			continue
		}
		cellNames := strings.Split(mapLines[i], ":")
		cellType, _ := strconv.Atoi(cellNames[1])
		img := getImageByType(cellImages, cellType)
		if img == nil {
			log.Fatal("Parse error")
			panic("Parse error")
		}

		switch cellNames[0] {
		case forwardMovingName:
			cells = append(cells, &Cell{img,
				Point{float64(col), float64(row)}, 0.0, cellType})
		}
		row++
		if row > mapHeight {
			row = 0
			col++
		}
		i++
	}
}

// readMap parses a map from a map file to backCells and walls in cells
func readMap(gameMap string, cellImages []*MyImage) {
	var endFlag bool = false
	var row int = 0
	var col int = 0
	for i := 0; i < len(gameMap) && !endFlag; i++ {
		switch gameMap[i] {
		case '\n':
			if gameMap[i] == gameMap[i-1] {
				endFlag = true
			}
			row++
			col = 0
			continue
		case wallSymbol:
			cells = append(cells, &Cell{
				getImageByType(cellImages, wallCell),
				Point{float64(col), float64(row)},
				0.0, wallCell})
		case startSymbol:
			backCells = append(backCells, &BackGroundCell{
				getImageByType(cellImages, startTile),
				Point{float64(col), float64(row)}, startTile})
		case exitSymbol:
			backCells = append(backCells, &BackGroundCell{
				getImageByType(cellImages, exitTile),
				Point{float64(col), float64(row)}, exitTile})
		}
		col++
		if row > mapHeight {
			mapHeight = row
		}
		if col > mapWidth {
			mapWidth = col
		}
	}
}

// loadMapFromFile reads and parses map file
func loadMapFromFile() {
	var cellImages []*MyImage
	cellImages = append(cellImages,
		LoadImage("./textures/BlackTile.png", wallCell))
	cellImages = append(cellImages,
		LoadImage("./textures/RedTile.png", moveStraightCell))
	cellImages = append(cellImages,
		LoadImage("./textures/StartTile.png", startTile))
	cellImages = append(cellImages,
		LoadImage("./textures/ExitTile.png", exitTile))

	data, err := os.ReadFile("./map.map")
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	gameMap := string(data)

	readMap(gameMap, cellImages)
	readAvailableTiles(gameMap, cellImages)
}

// Return image from cellImages wich corresponds to a cellType.
// In case getImageByType found nothing, it will return nil
func getImageByType(cellImages []*MyImage, cellType int) *MyImage {
	for _, img := range cellImages {
		if (*img).cellType == cellType {
			return img
		}
	}
	return nil
}

// Loads png file into ebiten.Image struct
func LoadImage(path string, cellType int) *MyImage {
	img, i, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	return &MyImage{img, Point{float64(i.Bounds().Max.X),
		float64(i.Bounds().Max.Y)}, cellType}
}
