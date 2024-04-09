package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	gameMap             [][]GameCell
	tempMap             [][]int
	screenWidth         int
	screenHeight        int
	numberOfCellsWidth  int
	numberOfCellsHeight int
	cellSize            int
	testImage           *image.RGBA
}

type GameCell struct {
	value           int
	image           *ebiten.Image
	imageOptions    *ebiten.DrawImageOptions
	neighborsValues [8]*int
}

func (g *Game) Update() error {

	//fmt.Println(ebiten.ActualFPS())

	for i := 0; i < g.numberOfCellsWidth; i++ {
		for j := 0; j < g.numberOfCellsHeight; j++ {

			//neighbors := GetNeighbors(i, j, g.numberOfCellsWidth, g.numberOfCellsHeight)
			//aliveNeighbors := GetAliveNeighbors(g.gameMap, neighbors)
			aliveNeighbors := GetAliveNeighborsPointers(g.gameMap[i][j].neighborsValues)

			g.tempMap[i][j] = DecideCellFuture(g.gameMap[i][j].value, aliveNeighbors)
		}
	}

	for i := 0; i < g.numberOfCellsWidth; i++ {
		for j := 0; j < g.numberOfCellsHeight; j++ {

			g.gameMap[i][j].value = g.tempMap[i][j]

			/*
				if g.gameMap[i][j].value == 1 {
					g.gameMap[i][j].image.Fill(color.White)
				} else {
					g.gameMap[i][j].image.Fill(color.Black)
				}
			*/
		}
	}

	/*
		[0]
		0 1 2 3 4 5 6 7 8 9


	*/

	for i := 0; i < g.numberOfCellsWidth; i++ {
		for j := 0; j < g.numberOfCellsHeight; j++ {

			//isBlack := g.gameMap[i][j].value == 0

			for r := 0; r < 10; r++ {
				for t := 0; t < 10; t++ {

					/*
						g.testImage.Pix[4*(i*10+r+g.screenWidth*4*j+t)] = uint8(255)
						g.testImage.Pix[4*(i*10+1+r+g.screenWidth*4*j+t)] = uint8(0)
						g.testImage.Pix[4*(i*10+2+r+g.screenWidth*4*j+t)] = uint8(0)
						g.testImage.Pix[4*(i*10+3+r+g.screenWidth*4*j+t)] = 0xff
					*/

					g.testImage.Pix[4*(g.screenWidth*i*10+j*r)] = uint8(255)
					g.testImage.Pix[4*(g.screenWidth*i*10+j*r)+1] = uint8(0)
					g.testImage.Pix[4*(g.screenWidth*i*10+j*r)+2] = uint8(0)
					g.testImage.Pix[4*(g.screenWidth*i*10+j*r)+3] = 0xff
				}
			}

			/*
				g.testImage.Pix[4*(g.screenWidth*i*10+j*10)] = uint8(255)
				g.testImage.Pix[4*(g.screenWidth*i*10+j*10)+1] = uint8(0)
				g.testImage.Pix[4*(g.screenWidth*i*10+j*10)+2] = uint8(0)
				g.testImage.Pix[4*(g.screenWidth*i*10+j*10)+3] = 0xff
			*/

		}

	}

	/*
		var l = g.screenWidth * g.screenHeight
		fmt.Println(l)
		fmt.Println(len(g.testImage.Pix))

		for i := 0; i < l; i++ {
			g.testImage.Pix[4*i] = uint8(i / g.screenWidth)
			g.testImage.Pix[4*i+1] = uint8(i / g.screenHeight)
			g.testImage.Pix[4*i+2] = uint8(i % g.screenWidth)
			g.testImage.Pix[4*i+3] = 0xff
		}*/

	return nil
}

type CellPosition struct {
	i int
	j int
}

func GetNeighbors(i int, j int, maxWidth int, maxHeight int) [8]CellPosition {
	neighbors := [8]CellPosition{}

	/*
		(i-1, j-1) (i, j-1) (i+1, j-1)
		(i-1, j)   (i, j)   (i+1, j)
		(i-1, j+1) (i, j+1) (i+1, j+1)
	*/

	neighbors[0] = CellPosition{i: i - 1, j: j - 1}
	neighbors[1] = CellPosition{i: i, j: j - 1}
	neighbors[2] = CellPosition{i: i + 1, j: j - 1}
	neighbors[3] = CellPosition{i: i - 1, j: j}
	neighbors[4] = CellPosition{i: i + 1, j: j}
	neighbors[5] = CellPosition{i: i - 1, j: j + 1}
	neighbors[6] = CellPosition{i: i, j: j + 1}
	neighbors[7] = CellPosition{i: i + 1, j: j + 1}

	for i = 0; i < len(neighbors); i++ {
		if neighbors[i].i < 0 || neighbors[i].i >= maxWidth {
			neighbors[i].i = -1
			neighbors[i].j = -1
			continue
		}

		if neighbors[i].j < 0 || neighbors[i].j >= maxHeight {
			neighbors[i].i = -1
			neighbors[i].j = -1
			continue
		}
	}

	return neighbors
}

func GetNeighborsPointers(gameMap [][]GameCell, i int, j int, maxWidth int, maxHeight int) [8]*int {
	neighborsValuePointers := [8]*int{}

	var onLeftLimit bool = i-1 < 0
	var onRightLimit bool = i+1 >= maxWidth
	var onUpLimit bool = j-1 < 0
	var onDownLimit bool = j+1 >= maxHeight

	if !onUpLimit {
		if !onLeftLimit {
			neighborsValuePointers[0] = &gameMap[i-1][j-1].value
		}

		neighborsValuePointers[1] = &gameMap[i][j-1].value

		if !onRightLimit {
			neighborsValuePointers[2] = &gameMap[i+1][j-1].value
		}
	}

	if !onLeftLimit {
		neighborsValuePointers[3] = &gameMap[i-1][j].value
	}

	if !onRightLimit {
		neighborsValuePointers[4] = &gameMap[i+1][j].value
	}

	if !onDownLimit {
		if !onLeftLimit {
			neighborsValuePointers[5] = &gameMap[i-1][j+1].value
		}

		neighborsValuePointers[6] = &gameMap[i][j+1].value

		if !onRightLimit {
			neighborsValuePointers[7] = &gameMap[i+1][j+1].value
		}
	}

	return neighborsValuePointers
}

func GetAliveNeighborsPointers(neighborsPointers [8]*int) int {
	var alive int = 0

	for i := 0; i < len(neighborsPointers); i++ {
		if neighborsPointers[i] == nil {
			continue
		}

		var value int = *neighborsPointers[i]
		if value == 1 {
			alive++
		}
	}
	return alive
}

func GetAliveNeighbors(gameMap [][]GameCell, neighbors [8]CellPosition) int {
	var alive int = 0

	for i := 0; i < len(neighbors); i++ {
		if neighbors[i].i == -1 || neighbors[i].j == -1 {
			continue
		}

		var value int = gameMap[neighbors[i].i][neighbors[i].j].value

		if value == 1 {
			alive++
		}
	}

	return alive
}

func DecideCellFuture(cellValue int, aliveCells int) int {
	if cellValue == 1 {
		if aliveCells < 2 || aliveCells > 3 {
			return 0
		} else if aliveCells == 2 || aliveCells == 3 {
			return 1
		}
	}

	if aliveCells == 3 {
		return 1
	}

	return 0
}

func (g *Game) Draw(screen *ebiten.Image) {

	/*
		for i := 0; i < g.numberOfCellsWidth; i++ {
			for j := 0; j < g.numberOfCellsHeight; j++ {
				screen.DrawImage(g.gameMap[i][j].image, g.gameMap[i][j].imageOptions)

				screen.WritePixels()
			}
		}
	*/

	screen.WritePixels(g.testImage.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

func main() {
	fmt.Println("start")

	var cellSize int = 10
	var screenWidth int = 2000
	var screenHeight int = 1000
	var numberOfCellsWidth int = screenHeight / cellSize
	var numberOfCellsHeight int = screenWidth / cellSize

	game := &Game{
		gameMap:             make([][]GameCell, numberOfCellsWidth),
		tempMap:             make([][]int, numberOfCellsWidth),
		screenWidth:         screenWidth,
		screenHeight:        screenHeight,
		numberOfCellsWidth:  numberOfCellsWidth,
		numberOfCellsHeight: numberOfCellsHeight,
		cellSize:            cellSize,
		testImage:           image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}

	for i := 0; i < numberOfCellsWidth; i++ {
		for j := 0; j < numberOfCellsHeight; j++ {
			var cellValue int = rand.Intn(2)
			image := ebiten.NewImage(game.cellSize, game.cellSize)
			if cellValue == 1 {
				image.Fill(color.White)
			} else {
				image.Fill(color.Black)
			}

			position := ebiten.GeoM{}
			position.Translate(float64(j)*float64(cellSize), float64(i)*float64(cellSize))
			drawImageOptions := ebiten.DrawImageOptions{GeoM: position}

			gameCell := GameCell{
				value:        rand.Intn(2),
				image:        image,
				imageOptions: &drawImageOptions,
			}
			game.gameMap[i] = append(game.gameMap[i], gameCell)
			game.tempMap[i] = append(game.tempMap[i], 0)
		}
	}

	for i := 0; i < numberOfCellsWidth; i++ {
		for j := 0; j < numberOfCellsHeight; j++ {

			neighborsPointers := GetNeighborsPointers(game.gameMap, i, j, game.numberOfCellsWidth, game.numberOfCellsHeight)

			game.gameMap[i][j].neighborsValues = neighborsPointers

		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game Of Life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
