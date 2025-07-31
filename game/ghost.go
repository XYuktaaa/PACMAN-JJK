package main

import(
 "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "math"
    //"fmt"

)
type Ghost struct{
    X,Y       float64
    Speed     float64
    Image     *ebiten.Image
    //Width     int
    //Height    int
    Direction string
    Name      string
    visible   bool
    Size      int
    TargetX   int
    TargetY   int
 	Path      []Node
 	PathIndex int
}


func NewGhost(x,y float64, spritePath string,name string,Size int, ) *Ghost{
    img, _, err := ebitenutil.NewImageFromFile(spritePath)
    if err != nil{
        panic(err)
    }

    //w,h := img.Size()
    return &Ghost{
        X:         x,
        Y:         y,
        Speed:     1.5,
        Name:      name,
        visible:   true,
        Image:     img,
//        Width:     w, 
  //      Height:    h,
        Direction: "left",
        Size:      55,
        
    }
}

func (g *Ghost) updateTarget(playerX, playerY, blinkyX, blinkyY float64, TileSize int) {
    switch g.Name {
    case "jogo": // 🔴 Blinky
        g.TargetX = int(playerX)
        g.TargetY = int(playerY)

    case "kenjaku": // 🟣 Pinky
        // Predictive movement: 4 tiles ahead of player
        offset := 4 * TileSize
        switch g.Direction {
        case "up":
            g.TargetX = int(playerX)
            g.TargetY = int(playerY - float64(offset))
        case "down":
            g.TargetX = int(playerX)
            g.TargetY = int(playerY + float64(offset))
        case "left":
            g.TargetX = int(playerX - float64(offset))
            g.TargetY = int(playerY)
        case "right":
            g.TargetX = int(playerX + float64(offset))
            g.TargetY = int(playerY)
        }

    case "mahito": // 🟠 Clyde
        // If far, chase; else retreat
        dist := distance(g.X, g.Y, playerX, playerY)
        if dist > 8*float64(TileSize) {
            g.TargetX = int(playerX)
            g.TargetY = int(playerY)
        } else {
            g.TargetX = 0
            g.TargetY = len(level)*TileSize - 1 // bottom-left corner
        }

    case "sakuna": // 🔵 Inky
        // Vector from Blinky through Player (P + (P - B))
        dx := playerX - blinkyX
        dy := playerY - blinkyY
        g.TargetX = int(playerX + dx)
        g.TargetY = int(playerY + dy)

    default:
        g.TargetX = int(playerX)
        g.TargetY = int(playerY)
    }
}


func (g *Ghost) Update(level [][]int, TileSize int, playerX, playerY float64) {
	// Set target based on ghost type
	switch g.Name {
	case "jogo": // blinky
		g.TargetX = int(playerX) / TileSize
		g.TargetY = int(playerY) / TileSize
	case "kenjaku": // inky - random wander (for now)
		if g.Path == nil || g.PathIndex >= len(g.Path) {
			g.TargetX = (g.TargetX + 5) % len(level[0])
			g.TargetY = (g.TargetY + 3) % len(level)
		}
	case "mahito": // clyde - patrol behavior
		g.TargetX = 1
		g.TargetY = 1
	case "sakuna": // pinky - ambush
		dx := int(playerX+24) / TileSize
		dy := int(playerY+24) / TileSize
		g.TargetX = dx
		g.TargetY = dy
	}


	g.moveToTarget(level, TileSize)
}

func (g *Ghost) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.X, g.Y)
	screen.DrawImage(g.Image, op)
}



func (g *Ghost) moveToTarget(level [][]int, TileSize int) {
	gx := int(g.X) / TileSize
	gy := int(g.Y) / TileSize

	if g.Path == nil || g.PathIndex >= len(g.Path) {
			path := findPath(level, gx, gy, g.TargetX, g.TargetY)
			// if len(path) > 0 {
   			//  g.Path = path
   			//  path := findPath(level, gx, gy, g.TargetX, g.TargetY)
   		if len(path) > 1 {
    		next := path[1] // path[0] is current tile
    		g.X = float64(next.X * TileSize)
    		g.Y = float64(next.Y * TileSize)
    		g.Path = path          // ← Save the path for reuse
            g.PathIndex = 1
            return
		}


	}

	if g.PathIndex < len(g.Path) {
		next := g.Path[g.PathIndex]
		// dx := float64(next.x*TileSize) + float64(TileSize/2) - g.X - float64(g.Size)/2
		// dy := float64(next.y*TileSize) + float64(TileSize/2) - g.Y - float64(g.Size)/2
		dx := float64(next.X*TileSize) + float64(TileSize/2) - g.X - float64(g.Size)/2
		dy := float64(next.Y*TileSize) + float64(TileSize/2) - g.Y - float64(g.Size)/2

		dist := math.Hypot(dx, dy)
		if dist < g.Speed {
			// g.X = float64(next.x*TileSize) + float64(TileSize-g.Size)/2
			// g.Y = float64(next.y*TileSize) + float64(TileSize-g.Size)/2
			g.X = float64(next.X*TileSize) + float64(TileSize-g.Size)/2
			g.Y = float64(next.Y*TileSize) + float64(TileSize-g.Size)/2

			g.PathIndex++
		} else {
			dx /= dist
			dy /= dist
			g.X += dx * g.Speed
			g.Y += dy * g.Speed
		}
	}
}


func distance(x1, y1, x2, y2 float64)float64{
    return math.Hypot(x2-x1, y2-y1)
}
