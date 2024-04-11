package main

import (
	"fmt"
	"image"
	"log"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	gameMap             [][]GameCell
	tempMap             [][]int
	screenWidth         int
	screenHeight        int
	numberOfCellsWidth  int
	numberOfCellsHeight int
	cellSize            int
	image               *image.RGBA
	isPaused            bool
}

type GameCell struct {
	value           int
	neighborsValues [8]*int
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.isPaused = !g.isPaused
	}

	if !g.isPaused {
		for i := 0; i < g.numberOfCellsWidth; i++ {
			for j := 0; j < g.numberOfCellsHeight; j++ {
				aliveNeighbors := GetAliveNeighborsPointers(g.gameMap[i][j].neighborsValues)
				g.tempMap[i][j] = DecideCellFuture(g.gameMap[i][j].value, aliveNeighbors)
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		var i int = y / g.cellSize
		var j int = x / g.cellSize

		if g.tempMap[i][j] == 0 {
			g.tempMap[i][j] = 1
		} else {
			g.tempMap[i][j] = 0
		}
	}

	for i := 0; i < g.numberOfCellsWidth; i++ {
		for j := 0; j < g.numberOfCellsHeight; j++ {

			g.gameMap[i][j].value = g.tempMap[i][j]

			var colour uint8 = 0

			if g.gameMap[i][j].value == 1 {
				colour = 255
			}

			for r := 0; r < g.cellSize; r++ {
				for t := 0; t < g.cellSize; t++ {

					g.image.Pix[4*((g.screenWidth*i*g.cellSize+(t*g.screenWidth))+(j*g.cellSize+r))] = uint8(colour)
					g.image.Pix[4*((g.screenWidth*i*g.cellSize+(t*g.screenWidth))+(j*g.cellSize+r))+1] = uint8(colour)
					g.image.Pix[4*((g.screenWidth*i*g.cellSize+(t*g.screenWidth))+(j*g.cellSize+r))+2] = uint8(colour)
					g.image.Pix[4*((g.screenWidth*i*g.cellSize+(t*g.screenWidth))+(j*g.cellSize+r))+3] = 0xff

				}
			}

			g.image.Pix[4*(g.screenWidth*i*g.cellSize+j*g.cellSize)] = uint8(0)
			g.image.Pix[4*(g.screenWidth*i*g.cellSize+j*g.cellSize)+1] = uint8(0)
			g.image.Pix[4*(g.screenWidth*i*g.cellSize+j*g.cellSize)+2] = uint8(255)
			g.image.Pix[4*(g.screenWidth*i*g.cellSize+j*g.cellSize)+3] = 0xff

		}

	}

	return nil
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
	screen.WritePixels(g.image.Pix)

	// Draw the message.
	tutorial := "Space: Move forward\nLeft/Right: Rotate"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), tutorial)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

func main() {
	var cellSize int = 10
	var screenWidth int = 1000
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
		image:               image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		isPaused:            false,
	}

	for i := 0; i < numberOfCellsWidth; i++ {
		for j := 0; j < numberOfCellsHeight; j++ {
			gameCell := GameCell{
				value: rand.Intn(2),
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

	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
