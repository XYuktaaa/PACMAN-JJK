package main

import(
 "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "math"
    "fmt"
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
    TargetX   int
    TargetY   int

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

func (g *Ghost) Update(level [][]int, TileSize int, playerX, playerY float64, blinkyX, blinkyY float64) {
    g.updateTarget(playerX, playerY, blinkyX, blinkyY, TileSize)

    // Only recalculate direction at tile centers
    if int(g.X)%TileSize == 0 && int(g.Y)%TileSize == 0 {
        g.moveTowardTarget(level, TileSize)
    }

    // Try moving in current direction
    nextX, nextY := g.X, g.Y
    switch g.Direction {
    case "left":
        nextX -= g.Speed
    case "right":
        nextX += g.Speed
    case "up":
        nextY -= g.Speed
    case "down":
        nextY += g.Speed
    }

    // Only move if no wall
    if !isWallCollidingLenient(level, nextX, nextY, g.Size, TileSize) {
        g.X = nextX
        g.Y = nextY
    }
    mapWidth := len(level[0]) * TileSize
	mapHeight := len(level) * TileSize

	if g.X < 0 { g.X = 0 }
	if g.Y < 0 { g.Y = 0 }
	if g.X > float64(mapWidth - g.Size) { g.X = float64(mapWidth - g.Size) }
	if g.Y > float64(mapHeight - g.Size) { g.Y = float64(mapHeight - g.Size) }


    fmt.Printf("Ghost pos: %.2f, %.2f | direction: %s\n", g.X, g.Y, g.Direction)
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

func (g *Ghost) moveTowardTarget(level [][]int, TileSize int){
    type Direction struct{
        Name    string
        offsetX float64
        offsetY float64
    }
    
    directions := []Direction{
        {"left",-g.Speed,0},
        {"right",g.Speed,0},
        {"up",0,-g.Speed},
        {"down",0,g.Speed},
    }
    //prevent reversing unless stuck
    opposite := map[string]string{
        "left": "right",
        "right":"left",
        "up":   "down",
        "down": "up",

    }
   
    bestDir := g.Direction
    shortestDist := math.MaxFloat64
    moved := false

    for _, dir := range directions{
        // Skip the reverse direction unless stuck
		if dir.Name == opposite[g.Direction] {
			continue
		}
        nx := g.X + dir.offsetX
        ny := g.Y + dir.offsetY
        fmt.Printf("%s trying %s → (%.2f, %.2f)\n", g.Name, dir.Name, nx, ny)


        if isWallCollidingLenient(level, nx, ny, g.Size, TileSize) {
			continue // Skip this direction if it hits a wall
		

        	dist := distance(nx, ny, float64(g.TargetX), float64(g.TargetY))
			if dist < shortestDist {
				shortestDist = dist
				bestDir = dir.Name
			}
			moved =true 
			fmt.Printf("%s chose direction: %s\n", g.Name, g.Direction)
			fmt.Printf("Trying %s → (%.2f, %.2f), dist: %.2f\n", dir.Name, nx, ny, dist)
		}
    }
    if !moved {
		// If stuck, allow reversing
		for _, dir := range directions {
			nx := g.X + dir.offsetX
			ny := g.Y + dir.offsetY
			if !isWallColliding(level, nx, ny, g.Size, TileSize) {
				bestDir = dir.Name
				break
			}
		}
	}
	g.Direction = bestDir
	// Move in chosen direction
	switch g.Direction {
	case "left":
		g.X -= g.Speed
	case "right":
		g.X += g.Speed
	case "up":
		g.Y -= g.Speed
	case "down":
		g.Y += g.Speed
	}
}

func distance(x1, y1, x2, y2 float64)float64{
    return math.Hypot(x2-x1, y2-y1)
}
