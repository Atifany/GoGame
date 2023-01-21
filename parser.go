package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nuttech/bell"
)

// map tiles types
const wallSymbol byte = 'X'
const startSymbol byte = 'V'
const exitSymbol byte = 'E'

type MyImage struct {
	image    *ebiten.Image
	cellType int
}

func checkParseError(err error) {
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
}

func isTile(cellType int) bool {
	if cellType >= startTile {
		return true
	}
	return false
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
		img, err := getImageByType(cellImages, cellType)
		checkParseError(err)

		cell := &Cell{img.image, &Transform{
			Point{float64(col), float64(row)},
			Point{float64(col), float64(row)}, 1.0, 1.0},
			0.0, cellType, false, false}
		bell.Listen("LMB_pressed", cell.PressDetect)
		bell.Listen("LMB_released", cell.releaseDetect)
		cells = append(cells, cell)

		row++
		if row > mapHeight {
			row = 0
			col++
		}
		i++
	}
}

// readMap parses a map from a map file to backCells and walls in cells
func readMap(gameMap string, cellImages []*MyImage, tileTypeLookUp map[byte]int) {
	var row int = 0
	var col int = 0
	for i := 0; i < len(gameMap); i++ {
		if i != 0 && gameMap[i] == gameMap[i-1] && gameMap[i] == '\n' {break}
		if gameMap[i] == ' ' { col++; continue }
		if gameMap[i] == '\n' {
			row++
			if row > mapHeight { mapHeight = row }
			if col > mapWidth { mapWidth = col }
			col = 0
			continue
		}

		cellType := tileTypeLookUp[gameMap[i]]
		img, err := getImageByType(cellImages, tileTypeLookUp[gameMap[i]])
		checkParseError(err)

		if isTile(cellType) {
			tiles = append(tiles, &Tile{ img.image, &Transform{
				Point{float64(col), float64(row)},
				Point{float64(col), float64(row)}, 1.0, 1.0},
				cellType, false})
		} else {
			cell := &Cell{ img.image, &Transform{
				Point{float64(col), float64(row)},
				Point{float64(col), float64(row)}, 1.0, 1.0},
				0.0, cellType, false, false}
			cells = append(cells, cell)
			cellsReady++
		}
		col++;
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
	cellImages = append(cellImages,
		LoadImage("./textures/DublicationCell.png", dublicationCell))

	tileTypeLookUp := make(map[byte]int)
	tileTypeLookUp[wallSymbol] = wallCell
	tileTypeLookUp[startSymbol] = startTile
	tileTypeLookUp[exitSymbol] = exitTile

	data, err := os.ReadFile("./map.map")
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	gameMap := string(data)

	readMap(gameMap, cellImages, tileTypeLookUp)
	readAvailableTiles(gameMap, cellImages)
}

// Return image from cellImages wich corresponds to a cellType.
// In case getImageByType found nothing, it will return nil
func getImageByType(cellImages []*MyImage, cellType int) (*MyImage, error) {
	for _, img := range cellImages {
		if (*img).cellType == cellType {
			return img, nil
		}
	}
	return nil, errors.New("Parse error")
}

// Loads png file into ebiten.Image struct
func LoadImage(path string, cellType int) *MyImage {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	return &MyImage{img, cellType}
}
