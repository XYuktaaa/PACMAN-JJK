package main

import (
   // "fmt"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "math"
)

type Player struct {
    X, Y float64 // coordinates of player
    Speed float64
    Direction string
    NextDirection string // for buffered input
    Image *ebiten.Image
    Width int
    Height int
    Score  int
    Size   int
    
    // Animation properties
    AnimationTimer int
    MouthOpenAngle float64  // For mouth animation
    IsMoving bool
    
    // Movement properties
    ContinuousMovement bool
    MovementBuffer string // Store next intended direction
    CollectionRadius float64
}

func NewPlayer(x, y float64, spritePath string) *Player {
    img, _, err := ebitenutil.NewImageFromFile(spritePath)
    if err != nil {
        panic(err)
    }

    // Make player size 45x45 for better movement
    targetSize := 45
    
    return &Player{
        X:         x,
        Y:         y,
        Speed:     2,
        Image:     img,
        Direction: "right",
        NextDirection: "right",
        Width:     targetSize,
        Height:    targetSize,
        Size:      targetSize,
        AnimationTimer: 0,
        MouthOpenAngle: 0,
        IsMoving: false,
        ContinuousMovement: false,
        MovementBuffer: "",
        CollectionRadius: float64(targetSize)*0.8,
    }
}

func (p *Player) CheckPelletCollection(level [][]int, pellets [][]bool, TileSize int) (int, bool) {
    points := 0
    powerPelletEaten := false
    
    // Check multiple points around player for more lenient collection
    checkPoints := []struct{x, y float64}{
        {p.X + float64(p.Size)/2, p.Y + float64(p.Size)/2}, // Center
        {p.X + float64(p.Size)*0.3, p.Y + float64(p.Size)*0.3}, // Top-left area
        {p.X + float64(p.Size)*0.7, p.Y + float64(p.Size)*0.3}, // Top-right area
        {p.X + float64(p.Size)*0.3, p.Y + float64(p.Size)*0.7}, // Bottom-left area
        {p.X + float64(p.Size)*0.7, p.Y + float64(p.Size)*0.7}, // Bottom-right area
    }
    
    for _, point := range checkPoints {
        tileX := int(point.x / float64(TileSize))
        tileY := int(point.y / float64(TileSize))
        
        // Bounds check
        if tileY >= 0 && tileY < len(level) && tileX >= 0 && tileX < len(level[0]) {
            if pellets[tileY][tileX] {
                pellets[tileY][tileX] = false
                
                if level[tileY][tileX] == TilePowerPellet {
                    points += 50
                    powerPelletEaten = true
                } else if level[tileY][tileX] == TilePellet {
                    points += 10
                }
                
                // Only collect one pellet per frame
                break
            }
        }
    }
    
    return points, powerPelletEaten
}


func (p *Player) Draw(screen *ebiten.Image) {
    if p.Image == nil {
        return
    }

    op := &ebiten.DrawImageOptions{}
    
    // Scale the image to match ghost size (55x55)
    originalW, originalH := p.Image.Size()
    scaleX := float64(p.Width) / float64(originalW)
    scaleY := float64(p.Height) / float64(originalH)
    
    // Apply scaling
    op.GeoM.Scale(scaleX, scaleY)
    
    // Handle left/right flipping
    if p.Direction == "left" {
        // Flip horizontally by scaling by -1 and adjusting position
        op.GeoM.Scale(-1, 1)
        op.GeoM.Translate(float64(p.Width), 0)
    }
    
    // Mouth animation effect (optional visual enhancement)
    if p.IsMoving {
        // Create a subtle pulsing effect for mouth animation
        pulse := math.Sin(float64(p.AnimationTimer) * 0.3)
        mouthScale := 0.95 + pulse*0.05
        
        // Apply mouth animation scaling from center
        op.GeoM.Translate(-float64(p.Width)/2, -float64(p.Height)/2)
        op.GeoM.Scale(mouthScale, mouthScale)
        op.GeoM.Translate(float64(p.Width)/2, float64(p.Height)/2)
    }
    
    // Final position translation
    op.GeoM.Translate(p.X, p.Y)
    
    screen.DrawImage(p.Image, op)
}

func (p *Player) Update(level [][]int, TileSize int) {
    p.AnimationTimer++
    
    // Simple movement - just like your original but with relaxed collision
    nextX, nextY := p.X, p.Y
    moved := false

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

    // Check horizontal movement first
    if nextX != p.X && !isWallCollidingRelaxed(level, nextX, p.Y, p.Size, TileSize) {
        p.X = nextX
        moved = true
    }
    // Check vertical movement
    if nextY != p.Y && !isWallCollidingRelaxed(level, p.X, nextY, p.Size, TileSize) {
        p.Y = nextY
        moved = true
    }
    
    p.IsMoving = moved
    
    p.handleScreenWrapping(level, TileSize)
}

// Slightly relaxed collision detection - not too lenient
func isWallCollidingRelaxed(level [][]int, px, py float64, size, tileSize int) bool {
    margin := 15// Small margin for 45px player - just enough to move smoothly
    
    // Check the four corners with small margin
    corners := [][2]int{
        {int(px) + margin, int(py) + margin},                    // top-left
        {int(px + float64(size)) - margin, int(py) + margin},    // top-right
        {int(px) + margin, int(py + float64(size)) - margin},    // bottom-left
        {int(px + float64(size)) - margin, int(py + float64(size)) - margin}, // bottom-right
    }
    
    for _, corner := range corners {
        cx := corner[0] / tileSize
        cy := corner[1] / tileSize
        
        // Prevent out-of-bounds access
        if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
            return true
        }
        
        if level[cy][cx] == TileWall {
            return true
        }
    }
    
    return false
}

// Check if player can change direction at current position
func (p *Player) canChangeDirection(level [][]int, TileSize int, newDirection string) bool {
    testX, testY := p.X, p.Y
    
    // Test a small movement in the new direction
    switch newDirection {
    case "right":
        testX += p.Speed
    case "left":
        testX -= p.Speed
    case "up":
        testY -= p.Speed
    case "down":
        testY += p.Speed
    }
    
    return !isWallCollidingRelaxed(level, testX, testY, p.Size, TileSize)
}


func (p *Player) handleScreenWrapping(level [][]int, TileSize int) {
    screenWidth := float64(len(level[0]) * TileSize)
    
    // Left tunnel
    if p.X < -float64(p.Size) {
        p.X = screenWidth
    }
    // Right tunnel
    if p.X > screenWidth {
        p.X = -float64(p.Size)
    }
}

// Reset player position (called when player dies)
func (p *Player) Reset(startX, startY float64) {
    p.X = startX
    p.Y = startY
    p.Direction = "right"
    p.NextDirection = "right"
    p.ContinuousMovement = false
    p.MovementBuffer = ""
    p.IsMoving = false
    p.AnimationTimer = 0
}

// Stop player movement (useful for game states like pause)
func (p *Player) StopMovement() {
    p.ContinuousMovement = false
    p.IsMoving = false
    p.MovementBuffer = ""
}


