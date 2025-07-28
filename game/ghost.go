package main

import(
 "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)
type Ghost struct{
    X,Y       float64
    Speed     float64
    Image     *ebiten.Image
    Width     int
    Height    int
    Direction string
    Name      string
    visible   bool
    Size      int
}

func NewGhost(x,y float64, spritePath string,name string,Size int) *Ghost{
    img, _, err := ebitenutil.NewImageFromFile(spritePath)
    if err != nil{
        panic(err)
    }

    w,h := img.Size()
    return &Ghost{
        X:         x,
        Y:         y,
        Speed:     1.5,
        Name:      name,
        visible:   true,
        Image:     img,
        Width:     w, 
        Height:    h,
        Direction: "left",
        Size:      55,
    }
}

func(g *Ghost) Update (level[][]int, TileSize int){
    nextX, nextY := g.X, g.Y

    switch g.Direction{
        case "left":
        nextX -= g.Speed
    case "right":
        nextX += g.Speed
    case "up":
        nextY -= g.Speed
    case "down":
        nextY += g.Speed
    }
   if !isWallColliding(level, nextX, nextY, g.Size, TileSize) {
        g.X = nextX
        g.Y = nextY
    } else {
        // reverse direction on wall hit
        switch g.Direction {
        case "left":
            g.Direction = "right"
        case "right":
            g.Direction = "left"
        case "up":
            g.Direction = "down"
        case "down":
            g.Direction = "up"
        }
    }
}
func (g *Ghost) Draw(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    w, h := g.Image.Size()

    scaleX := float64(g.Size)/float64(w)
	scaleY := float64(g.Size)/float64(h)

    op.GeoM.Scale(scaleX,scaleY)

    op.GeoM.Translate(g.X, g.Y)
    screen.DrawImage(g.Image, op)
}

