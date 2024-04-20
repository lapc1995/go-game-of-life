package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
)

type Game struct {
	ui                  *ebitenui.UI
	gameMap             [][]GameCell
	tempMap             [][]int
	screenWidth         int
	screenHeight        int
	numberOfCellsWidth  int
	numberOfCellsHeight int
	cellSize            int
	image               *image.RGBA
	isPaused            bool
	startingCellSize    int

	zoomLevel   int
	zoomCenterX int
	zoomCenterY int
	zoomedImage ebiten.Image
	zoomX       int
	zoomY       int
	minX        int
	maxX        int
	minY        int
	maxY        int
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

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.StartGame()
	}

	_, dy := ebiten.Wheel()

	if !g.isPaused {
		for i := 0; i < g.numberOfCellsWidth; i++ {
			for j := 0; j < g.numberOfCellsHeight; j++ {
				aliveNeighbors := GetAliveNeighborsPointers(g.gameMap[i][j].neighborsValues)
				g.tempMap[i][j] = DecideCellFuture(g.gameMap[i][j].value, aliveNeighbors)
			}
		}
	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			for i := 0; i < g.numberOfCellsWidth; i++ {
				for j := 0; j < g.numberOfCellsHeight; j++ {
					aliveNeighbors := GetAliveNeighborsPointers(g.gameMap[i][j].neighborsValues)
					g.tempMap[i][j] = DecideCellFuture(g.gameMap[i][j].value, aliveNeighbors)
				}
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

	hasZoomed := false

	// zoom in
	if dy == 1 {
		g.zoomLevel += 1
		hasZoomed = true
	} else if dy == -1 {
		hasZoomed = true
		if g.zoomLevel-0 < 0 {
			g.zoomLevel = 0
		} else {
			g.zoomLevel -= 1
		}
	}

	if g.zoomLevel > 0 {

		if hasZoomed {
			g.zoomX = 1000 - (10 * g.zoomLevel)
			g.zoomY = 1000 - (10 * g.zoomLevel)

			var zoomXHalf float32 = float32(g.zoomX) / 2.0
			var zoomYHalf float32 = float32(g.zoomY) / 2.0

			mouseX, mouseY := ebiten.CursorPosition()

			g.minX = mouseX - int(zoomXHalf)
			g.maxX = mouseX + int(zoomXHalf)

			g.minY = mouseY - int(zoomYHalf)
			g.maxY = mouseY + int(zoomYHalf)

		}

		g.zoomedImage = *ebiten.NewImageFromImage(g.image.SubImage(image.Rect(g.minX, g.minY, g.maxX, g.maxY)))
	}

	g.ui.Update()

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

func (g *Game) StartGame() {

	g.gameMap = make([][]GameCell, g.numberOfCellsWidth)
	g.tempMap = make([][]int, g.numberOfCellsWidth)
	g.image = image.NewRGBA(image.Rect(0, 0, g.screenWidth, g.screenHeight))

	for i := 0; i < g.numberOfCellsWidth; i++ {
		for j := 0; j < g.numberOfCellsHeight; j++ {
			gameCell := GameCell{
				value: rand.Intn(2),
			}
			g.gameMap[i] = append(g.gameMap[i], gameCell)
			g.tempMap[i] = append(g.tempMap[i], gameCell.value)
		}
	}

	for i := 0; i < g.numberOfCellsWidth; i++ {
		for j := 0; j < g.numberOfCellsHeight; j++ {
			neighborsPointers := GetNeighborsPointers(g.gameMap, i, j, g.numberOfCellsWidth, g.numberOfCellsHeight)
			g.gameMap[i][j].neighborsValues = neighborsPointers

		}
	}

}

func (g *Game) Draw(screen *ebiten.Image) {

	drawImageOptions := ebiten.DrawImageOptions{}
	drawImageOptions.GeoM.Translate(2000/2-800/2, 1000/2-800/2)
	if g.zoomLevel > 0 {

		new_scale := 1000 / float64(g.zoomX)
		fmt.Println("scale", new_scale)
		drawImageOptions.GeoM.Scale(new_scale, new_scale)
		screen.DrawImage(&g.zoomedImage, &drawImageOptions)
	} else {
		screen.DrawImage(ebiten.NewImageFromImage(g.image), &drawImageOptions)
	}

	// Draw the message.
	tutorial := "Space: Move forward\nLeft/Right: Rotate"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), tutorial)
	ebitenutil.DebugPrint(screen, msg)

	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	var cellSize int = 10
	var screenWidth int = 800
	var screenHeight int = 800
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
		startingCellSize:    cellSize,

		zoomLevel:   0,
		zoomCenterX: 0,
		zoomCenterY: 0,
	}

	game.StartGame()

	/*
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
		}*/

	ebiten.SetWindowSize(2000, 1000)
	ebiten.SetWindowTitle("Game Of Life")

	ebiten.SetFullscreen(false)

	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer()

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Hello World!", fontFace, color.White),
	)

	rootContainer.AddChild(helloWorldLabel)

	game.ui = eui

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
