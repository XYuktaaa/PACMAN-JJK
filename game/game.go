package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    // "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "fmt"
)


type Game struct{
    Player *Player
}

func NewGame() *Game {
    return &Game{
        Player: NewPlayer(float64(1*TileSize), float64(1*TileSize),"assets/player.png"),
    }
}


func (g *Game) Update()error{
    g.Player.Update(level, TileSize)
    return nil
}

const TileSize =32

func DrawMaze(screen *ebiten.Image) {
    for y, row := range level {
        for x, tile := range row {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(x * TileSize), float64(y * TileSize))

            switch tile {
            case TileWall:
                screen.DrawImage(WallImage, op)
            case TilePellet:
                screen.DrawImage(FloorImage, op) // base layer
                // draw the dot centered
                pelletOp := &ebiten.DrawImageOptions{}
                dx := float64(x*TileSize) + float64(TileSize/2-4)
                dy := float64(y*TileSize) + float64(TileSize/2-4)
                pelletOp.GeoM.Translate(dx, dy)
                screen.DrawImage(PelletImage, pelletOp)
            default:
                screen.DrawImage(FloorImage, op)
            }
        }
    }
}


func (g *Game) Draw(screen *ebiten.Image) {
    DrawMaze(screen) // draw background or maze first

    for y, row := range level {
        for x, tile := range row {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(x*TileSize), float64(y*TileSize))

            switch tile {
            case TileWall:
                screen.DrawImage(WallImage, op)
            case TilePellet:
                screen.DrawImage(FloorImage, op)
                screen.DrawImage(PelletImage, op)
            case TileEmpty:
                screen.DrawImage(FloorImage, op)
            }
        }
    }

    g.Player.Draw(screen) // draw player AFTER level so it appears on top
}


func (g *Game) Layout(outsideWidth, outsideHeight int)(int , int){
    width:= len(level[0])*TileSize
    height := len(level)*TileSize
    return width,height  
}
func isWallColliding(level [][]int, px, py float64, size, tileSize int) bool {
	corners := [][2]int{
		{int(px), int(py)},
		{int(px + float64(size) - 1), int(py)},
		{int(px), int(py + float64(size) - 1)},
		{int(px + float64(size) - 1), int(py + float64(size) - 1)},
	}

	for _, corner := range corners {
		cx := corner[0] / tileSize
		cy := corner[1] / tileSize
		fmt.Printf("Checking tile (%d, %d) = %d\n", cx, cy, level[cy][cx])
		

		// Prevent out-of-bounds access
		if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
			fmt.Println("Out of bounds, blocking movement")
			return true // treat OOB as wall
		}

		if level[cy][cx] == TileWall {
    		fmt.Println("Hit wall at", cx, cy)
			return true
		}
	}

	return false
}

