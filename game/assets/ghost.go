
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
	"fmt"
	
)

type GhostMode int

const (
	Chase GhostMode = iota
	Scatter
	Frightened
)

type Ghost struct {
	X, Y          float64
	Speed         float64
	Image         *ebiten.Image
	Direction     string
	Name          string
	Visible       bool
	Size          int
	TargetX       int
	TargetY       int
	Path          []Node
	PathIndex     int
	Mode          GhostMode
	ModeTimer     int
	ScatterTarget [2]int // Fixed scatter corner for this ghost
	LastTileX     int
	LastTileY     int
	MovementSmooth bool
	TargetUpdateTimer int
}

func NewGhost(x, y float64, spritePath string, name string, size int) *Ghost {
	img, _, err := ebitenutil.NewImageFromFile(spritePath)
	if err != nil {
		panic(err)
	}

	// Set scatter corners for each ghost
	var scatterTarget [2]int
	switch name {
	case "jogo":
		scatterTarget = [2]int{25, 1} // top-right
	case "sukuna":
		scatterTarget = [2]int{2, 1}  // top-left
	case "kenjaku":
		scatterTarget = [2]int{25, 29} // bottom-right
	case "mahito":
		scatterTarget = [2]int{2, 29}  // bottom-left
	}

	return &Ghost{
		X:             x,
		Y:             y,
		Speed:         1.2,
		Name:          name,
		Visible:       true,
		Image:         img,
		Direction:     "left",
		Size:          size,
		Mode:          Scatter,
		ModeTimer:     0,
		ScatterTarget: scatterTarget,
		LastTileX:     -1,
		LastTileY:     -1,
		MovementSmooth: true,
	}
}

// Update ghost behavior and movement
func (g *Ghost) Update(level [][]int, TileSize int, playerX, playerY float64, playerDirection string, blinkyX, blinkyY float64) {
	// Update mode timing (simplified - you might want more complex mode switching)
 	fmt.Printf("Ghost %s: Before update at (%.1f, %.1f)\n", g.Name, g.X, g.Y)
    
	g.ModeTimer++
	g.TargetUpdateTimer++
	// Only update target every 30 frames (0.5 seconds) to prevent erratic behavior
    if g.TargetUpdateTimer >= 30 {
        g.updateTarget(level, playerX, playerY, playerDirection, blinkyX, blinkyY, TileSize)
        g.TargetUpdateTimer = 0
    }
    // Always try to move towards current target
    g.moveToTarget(level, TileSize)
	if g.ModeTimer > 1200 { // Switch modes every 20 seconds at 60 FPS
		if g.Mode == Chase {
			g.Mode = Scatter
		} else {
			g.Mode = Chase
		}
		g.ModeTimer = 0
		g.Path = nil // Clear path when switching modes
	}

	// Update target based on current mode and ghost type
	g.updateTarget(level, playerX, playerY, playerDirection, blinkyX, blinkyY, TileSize)
	fmt.Printf("Ghost %s: Target set to (%d, %d)\n", g.Name, g.TargetX, g.TargetY)
    
	// Move towards target
	g.moveToTarget(level, TileSize)
	fmt.Printf("Ghost %s: After update at (%.1f, %.1f)\n", g.Name, g.X, g.Y)

}

// Set target tile based on ghost personality and current mode
func (g *Ghost) updateTarget(level [][]int, playerX, playerY float64, playerDirection string, blinkyX, blinkyY float64, TileSize int) {
    playerTileX := int(playerX) / TileSize
    playerTileY := int(playerY) / TileSize
    
    // Ensure player tile is within bounds
    if playerTileX < 0 { playerTileX = 0 }
    if playerTileX >= len(level[0]) { playerTileX = len(level[0]) - 1 }
    if playerTileY < 0 { playerTileY = 0 }
    if playerTileY >= len(level) { playerTileY = len(level) - 1 }

    if g.Mode == Scatter {
        // Safe scatter corners
        switch g.Name {
        case "jogo":
            g.TargetX = len(level[0]) - 2  // Top right corner
            g.TargetY = 1
        case "sukuna":
            g.TargetX = 1   // Top left corner  
            g.TargetY = 1
        case "kenjaku":
            g.TargetX = len(level[0]) - 2  // Bottom right corner
            g.TargetY = len(level) - 2
        case "mahito":
            g.TargetX = 1   // Bottom left corner
            g.TargetY = len(level) - 2
        }
        return
    }

    if g.Mode == Frightened {
        // Simple random movement - just pick a nearby empty tile
        currentTileX := int(g.X) / TileSize
        currentTileY := int(g.Y) / TileSize
        
        // Try to move in a random direction
        directions := [][2]int{{0,1}, {1,0}, {0,-1}, {-1,0}}
        for _, dir := range directions {
            newX := currentTileX + dir[0]*3
            newY := currentTileY + dir[1]*3
            
            // Bounds check
            if newX >= 0 && newX < len(level[0]) && newY >= 0 && newY < len(level) {
                if level[newY][newX] != TileWall {
                    g.TargetX = newX
                    g.TargetY = newY
                    return
                }
            }
        }
        // Fallback - just target current position
        g.TargetX = currentTileX
        g.TargetY = currentTileY
        return
    }

    // Chase mode - simplified and safe
    switch g.Name {
    case "jogo": // Direct chase
        g.TargetX = playerTileX
        g.TargetY = playerTileY

    case "sukuna": // Target 2 tiles ahead (reduced from 4 to prevent out of bounds)
        g.TargetX = playerTileX
        g.TargetY = playerTileY
        
        switch playerDirection {
        case "left":
            g.TargetX = playerTileX - 2
        case "right":
            g.TargetX = playerTileX + 2
        case "up":
            g.TargetY = playerTileY - 2
        case "down":
            g.TargetY = playerTileY + 2
        }

    case "kenjaku": // Simple chase with slight offset
        g.TargetX = playerTileX + 1
        g.TargetY = playerTileY + 1

    case "mahito": // Chase with distance check
        currentTileX := int(g.X) / TileSize
        currentTileY := int(g.Y) / TileSize
        distance := abs(playerTileX-currentTileX) + abs(playerTileY-currentTileY)
        
        if distance > 8 {
            g.TargetX = playerTileX
            g.TargetY = playerTileY
        } else {
            g.TargetX = 1
            g.TargetY = len(level) - 2
        }
    }

    // CRITICAL: Ensure target is always within bounds
    if g.TargetX < 0 { g.TargetX = 0 }
    if g.TargetX >= len(level[0]) { g.TargetX = len(level[0]) - 1 }
    if g.TargetY < 0 { g.TargetY = 0 }
    if g.TargetY >= len(level) { g.TargetY = len(level) - 1 }
    
    // If target is a wall, default to player position
    if level[g.TargetY][g.TargetX] == TileWall {
        g.TargetX = playerTileX
        g.TargetY = playerTileY
    }
}
// Add this helper function to your ghost.go:
func (g *Ghost) findNearestEmptySpace(level [][]int) {
    for radius := 1; radius < 10; radius++ {
        for dx := -radius; dx <= radius; dx++ {
            for dy := -radius; dy <= radius; dy++ {
                newX := g.TargetX + dx
                newY := g.TargetY + dy
                
                if newY >= 0 && newY < len(level) && newX >= 0 && newX < len(level[0]) {
                    if level[newY][newX] == TileEmpty || level[newY][newX] == TilePellet {
                        g.TargetX = newX
                        g.TargetY = newY
                        return
                    }
                }
            }
        }
    }
}

// Add this helper function too:
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
// Improved movement with smoother pathfinding
func (g *Ghost) moveToTarget(level [][]int, TileSize int) {
	currentTileX := int(g.X) / TileSize
	currentTileY := int(g.Y) / TileSize
	fmt.Printf("Ghost %s: moveToTarget called, current tile (%d,%d), target (%d,%d)\n", 
               g.Name, currentTileX, currentTileY, g.TargetX, g.TargetY)
	// Only recalculate path if we've moved to a new tile or don't have a path
	// if g.Path == nil || g.PathIndex >= len(g.Path) || 
	//    currentTileX != g.LastTileX || currentTileY != g.LastTileY {
if g.Path == nil || g.PathIndex >= len(g.Path) {
		
		path := findPath(level, currentTileX, currentTileY, g.TargetX, g.TargetY)
		fmt.Printf("Ghost %s: findPath returned %d nodes\n", g.Name, len(path))
		if len(path) > 1 {
			g.Path = path
			g.PathIndex = 1 // Skip current position
			g.LastTileX = currentTileX
			g.LastTileY = currentTileY
			fmt.Printf("Ghost %s: Path set, next node (%d,%d)\n", g.Name, path[1].X, path[1].Y)
		}else {
            fmt.Printf("Ghost %s: No valid path found!\n", g.Name)
        }
	}

	// Follow the path if we have one
	if g.Path != nil && g.PathIndex < len(g.Path) {
		g.followPath(TileSize)
	}else {
        fmt.Printf("Ghost %s: No path to follow\n", g.Name)
    }
}

// Smooth movement along the calculated path
func (g *Ghost) followPath(TileSize int) {
	if g.PathIndex >= len(g.Path) {
    	g.Path=nil
		return
	}

	next := g.Path[g.PathIndex]
	
	// Calculate target pixel position (center of tile)
	targetX := float64(next.X*TileSize) + float64(TileSize)/2 - float64(g.Size)/2
	targetY := float64(next.Y*TileSize) + float64(TileSize)/2 - float64(g.Size)/2

	// Calculate movement vector
	dx := targetX - g.X
	dy := targetY - g.Y
	distance := math.Hypot(dx, dy)

	// Update direction for animation
	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			g.Direction = "right"
		} else {
			g.Direction = "left"
		}
	} else {
		if dy > 0 {
			g.Direction = "down"
		} else {
			g.Direction = "up"
		}
	}

	// Move towards target
	if distance <= g.Speed {
		// Snap to target and move to next path node
		g.X = targetX
		g.Y = targetY
		g.PathIndex++
	} else {
		// Move smoothly towards target
		g.X += (dx / distance) * g.Speed
		g.Y += (dy / distance) * g.Speed
	}
}

// Draw ghost with proper scaling and optional direction-based sprite changes
func (g *Ghost) Draw(screen *ebiten.Image) {
	if !g.Visible {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Scale image to fit ghost size
	iw, ih := g.Image.Bounds().Dx(), g.Image.Bounds().Dy()
	scaleX := float64(g.Size) / float64(iw)
	scaleY := float64(g.Size) / float64(ih)
	op.GeoM.Scale(scaleX, scaleY)

	// Apply color tint based on mode
	switch g.Mode {
	case Frightened:
		op.ColorM.Scale(0.5, 0.5, 1.0, 1.0) // Blue tint for frightened
	case Scatter:
		op.ColorM.Scale(1.0, 1.0, 1.0, 0.8) // Slightly transparent in scatter mode
	}

	op.GeoM.Translate(g.X, g.Y)
	screen.DrawImage(g.Image, op)
}

// Set ghost to frightened mode (when player eats power pellet)
func (g *Ghost) SetFrightened(duration int) {
	g.Mode = Frightened
	g.ModeTimer = -duration // Negative timer to count up to 0
	g.Path = nil // Clear current path
	g.Speed = 0.6 // Slower when frightened
}

// Check if ghost can be eaten (is frightened)
func (g *Ghost) CanBeEaten() bool {
	return g.Mode == Frightened
}

func (g *Ghost) SetVisible(visible bool) {
	g.Visible = visible
}


// Reset ghost to normal chase/scatter behavior
func (g *Ghost) ResetMode() {
	g.Mode = Chase
	g.ModeTimer = 0
	g.Speed = 1.2
	g.Path = nil
}

// Get current tile position
func (g *Ghost) GetTilePosition(TileSize int) (int, int) {
	return int(g.X) / TileSize, int(g.Y) / TileSize
}

// Check collision with player (circular collision detection)
func (g *Ghost) CollidesWith(playerX, playerY float64, playerSize int) bool {
	centerX := g.X + float64(g.Size)/2
	centerY := g.Y + float64(g.Size)/2
	playerCenterX := playerX + float64(playerSize)/2
	playerCenterY := playerY + float64(playerSize)/2
	
	distance := math.Hypot(centerX-playerCenterX, centerY-playerCenterY)
	return distance < float64(g.Size+playerSize)/3 // Adjust collision sensitivity
}

