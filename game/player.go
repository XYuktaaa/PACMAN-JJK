package main

import (
    // "image"
    // "log"
    "fmt"
    "github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
	 
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	
)


type Player struct {
    X, Y float64 //cordinates of player
    Speed float64
    Direction string
    Image *ebiten.Image
    Width int
    Height int
    Score  int
    Size   int
}

func NewPlayer(x, y float64, spritePath string) *Player {
    img, _, err := ebitenutil.NewImageFromFile(spritePath)
    if err != nil {
        panic(err)
    }

    w, h := img.Size()

    return &Player{
        X:         x,
        Y:         y,
        Speed:     2,
        Image:     img,
        Direction: "right",
        Width:     w,
        Height:    h,
        Size:      w,
    }
    
}

func (p *Player) Draw(screen *ebiten.Image){
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(p.X, p.Y)
    screen.DrawImage(p.Image, op)
        


}

func (p *Player) Update(level [][]int, TileSize int) {
    nextX, nextY := p.X, p.Y

    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        nextX += p.Speed
        p.Direction = "right"
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        nextX -= p.Speed
        p.Direction = "left"
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        nextY -= p.Speed
        p.Direction = "up"
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        nextY += p.Speed
        p.Direction = "down"
    }

    //  collision check
    if !isWallCollidingStrict(level, nextX, p.Y, p.Width, TileSize) {
        p.X = nextX
    }
    if !isWallCollidingStrict(level, p.X, nextY, p.Height, TileSize) {
        p.Y = nextY
    }

    gridX := int(p.X + float64(TileSize)/2)/TileSize
    gridY := int(p.Y + float64(TileSize)/2)/TileSize

    if level[gridY][gridX]== TilePellet{
        level[gridY][gridX] = TileEmpty
        p.Score++
        fmt.Println("Pellet eaten! Score:", p.Score)
    }

    fmt.Println("Player pos:", p.X, p.Y)
}

