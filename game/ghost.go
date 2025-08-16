package main

import (
	"math"
	"math/rand"
	"time"
	"github.com/hajimehoshi/ebiten/v2"

"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"fmt"
	"image/color"
)


// Constants for game balance
const (
	FRIGHT_DURATION = 600 // 10 seconds at 60fps
	CHASE_DURATION = 1200 // 20 seconds
	SCATTER_DURATION = 420 // 7 seconds
	
	// Ghost house positions
	GHOST_HOUSE_X = 13
	GHOST_HOUSE_Y = 13
	GHOST_HOUSE_EXIT_Y = 11
)

// GhostMode represents the current state of a ghost
type GhostMode int

const (
	ChaseMode GhostMode = iota
	ScatterMode
	FrightenedMode
	DeadMode
	InHouseMode
)

// Ghost represents a ghost entity with advanced AI
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
	ScatterTarget [2]int
	LastTileX     int
	LastTileY     int
	// MovementSmooth bool
	
	// Advanced AI properties
	GhostType        string        // "blinky", "pinky", "inky", "clyde"
	BaseSpeed        float64       // Original speed for mode calculations
	FrightTimer      int           // Frames remaining in frightened mode
	ScatterTimer     int           // Frames remaining in scatter mode
	ChaseTimer       int           // Frames remaining in chase mode
	ReleaseTimer     int           // Frames until release from house
	PersonalityMode  int           // Counter for personality behaviors
	CruiseElroyMode  int           // Blinky's speed boost level (0, 1, 2)
	PreviousDirection string       // For avoiding reverse unless forced
	StuckCounter     int           // Detect when ghost is stuck
	LastUpdate       time.Time     // For timing calculations
	
	// Additional production features
	LastPosition     [2]int        // For stuck detection
	// drawDebugInfo    bool          // Debug visualization toggle
	// soundEnabled     bool          // Audio trigger toggle
	//networkSync      bool          // Network synchronization flag
	//difficultyLevel  int           // Current difficulty level
}


// Ghost manager for coordinated AI behavior
type GhostManager struct {
	ghosts []*Ghost
	gameState *GameStateStruct
	globalModeTimer int
	waveNumber int // For scatter/chase wave patterns
}

//NewGhostManager creates a coordinated ghost management system
func NewGhostManager(gameState *GameStateStruct) *GhostManager {
	return &GhostManager{
		ghosts: make([]*Ghost, 0),
		gameState: gameState,
		globalModeTimer: 0,
		waveNumber: 0,
	}
}

// AddGhost adds a ghost to the manager
func (gm *GhostManager) AddGhost(ghost *Ghost) {
	gm.ghosts = append(gm.ghosts, ghost)
	// gm.gameState.Ghosts = gm.ghosts
}

// UpdateAll updates all ghosts with coordinated behavior
func (gm *GhostManager) UpdateAll() {
	gm.globalModeTimer++
	gm.handleGlobalModeChanges()
	
	for _, ghost := range gm.ghosts {
		ghost.Update(gm.gameState)
	}
}

// Handle global mode changes (scatter/chase waves)
func (gm *GhostManager) handleGlobalModeChanges() {
	// Classic Pacman wave pattern
	wavePattern := []struct{
		duration int
		mode GhostMode
	}{
		{420, ScatterMode},   // 7 seconds scatter
		{1200, ChaseMode},    // 20 seconds chase
		{420, ScatterMode},   // 7 seconds scatter
		{1200, ChaseMode},    // 20 seconds chase
		{300, ScatterMode},   // 5 seconds scatter
		{1200, ChaseMode},    // 20 seconds chase
		{300, ScatterMode},   // 5 seconds scatter
		{-1, ChaseMode},      // Indefinite chase
	}
	
	if gm.waveNumber < len(wavePattern) {
		wave := wavePattern[gm.waveNumber]
		if wave.duration > 0 && gm.globalModeTimer >= wave.duration {
			gm.advanceWave()
		}
	}
}

func (gm *GhostManager) advanceWave() {
	gm.waveNumber++
	gm.globalModeTimer = 0
	
	// Update all ghosts to new mode
for _, ghost := range gm.ghosts {
		if ghost.Mode != FrightenedMode && ghost.Mode != DeadMode && ghost.Mode != InHouseMode {
			if gm.waveNumber < 8 {
				if gm.waveNumber%2 == 0 {
					ghost.Mode = ScatterMode
					ghost.ScatterTimer = 420
				} else {
					ghost.Mode = ChaseMode
					ghost.ChaseTimer = 1200
				}
			} else {
				ghost.Mode = ChaseMode
				ghost.ChaseTimer = -1 // Indefinite
			}
			
			// Force direction reversal on mode change (except first wave)
			if gm.waveNumber > 0 {
				ghost.reverseDirection()
			}
		}
	}
}
	
// TriggerFrightMode activates fright mode for all ghosts
func (gm *GhostManager) TriggerFrightMode() {
	for _, ghost := range gm.ghosts {
		ghost.SetFrightened(FRIGHT_DURATION)
	}
	gm.gameState.FrightModeActive = true
}

// CheckCollisions handles all ghost-player collisions
func (gm *GhostManager) CheckCollisions(playerX, playerY float64) string {
	for _, ghost := range gm.ghosts {
		result := ghost.CollideWithPlayer(playerX, playerY)
		if result != "no_collision" {
			return result
		}
	}
	return "no_collision"
}

// NewGhost creates a new ghost with advanced AI capabilities
func NewGhost(x, y float64, imagePath, ghostType string, size int) *Ghost {
	// Load image (you'll need to implement this based on your image loading system)
	var image *ebiten.Image
    image = loadImage(imagePath) // Implement this function
	
	ghost := &Ghost{
		X:               x,
		Y:               y,
		Speed:           0.8,
		BaseSpeed:       0.8,
		Image:           image,
		Direction:       "up",
		PreviousDirection: "up",
		Name:            ghostType,
		GhostType:       ghostType,
		Visible:         true,
		Size:            size,
		Mode:            InHouseMode,
		PathIndex:       0,
		LastUpdate:      time.Now(),
	}
	// Set ghost-specific properties
	switch ghostType {
	case "jogo"://Blinky
		ghost.ScatterTarget = [2]int{25, 0}    // Top-right
		ghost.ReleaseTimer = 0
		ghost.Mode = ChaseMode
		ghost.ChaseTimer=1200
	case "sakuna"://Pinky
		ghost.ScatterTarget = [2]int{2, 0}     // Top-left
		//ghost.Mode=ChaseMode
		ghost.ReleaseTimer = 300
	case "kenjaku"://Inky
		ghost.ScatterTarget = [2]int{25, 30}// Bottom-right
		//ghost.Mode=ChaseMode
		ghost.ReleaseTimer = 600
	case "mahito"://clyde
		ghost.ScatterTarget = [2]int{2, 30}    // Bottom-left
		ghost.ReleaseTimer = 900
		//ghost.Mode=ChaseMode
	}
	fmt.Printf("Created ghost %s at (%.1f, %.1f) with speed %.2f\n", 
		ghostType, x, y, ghost.Speed)
		return ghost
}

// Quick test function to verify tile constants
func TestTileConstants() {
	fmt.Println("Checking tile constants:")
	fmt.Printf("TileEmpty = %d\n", TileEmpty)
	fmt.Printf("TileWall = %d\n", TileWall)
	fmt.Printf("TilePellet = %d\n", TilePellet)
	fmt.Printf("TilePowerPellet = %d\n", TilePowerPellet)
}

// Update handles the main ghost AI logic
func (g *Ghost) Update(gameState *GameStateStruct) {
	now := time.Now()
	// deltaTime := now.Sub(g.LastUpdate).Seconds()
	g.LastUpdate = now


// Debug output (remove this in production)
	if g.PersonalityMode%60 == 0 { // Print every second
		pacmanTileX := int(gameState.PacmanX / TileSize)
		pacmanTileY := int(gameState.PacmanY / TileSize)
		ghostTileX := int(g.X / TileSize)
		ghostTileY := int(g.Y / TileSize)
		
		fmt.Printf("Ghost %s: Mode=%d, Pos=(%d,%d), Target=(%d,%d), Pacman=(%d,%d)\n",
			g.GhostType, g.Mode, ghostTileX, ghostTileY, g.TargetX, g.TargetY, pacmanTileX, pacmanTileY)
	}

	// Performance optimization - only update AI logic periodically
	// if !g.shouldUpdateAI() {
	// 	return
	// }


	if g.PersonalityMode%60 == 0 {
		fmt.Printf("=== GHOST DEBUG %s ===\n", g.GhostType)
		fmt.Printf("Position: (%.1f, %.1f) -> Tile: (%d, %d)\n", 
			g.X, g.Y, int(g.X/TileSize), int(g.Y/TileSize))
		fmt.Printf("Speed: %.2f, Direction: %s\n", g.Speed, g.Direction)
		fmt.Printf("Mode: %d, Target: (%d, %d)\n", g.Mode, g.TargetX, g.TargetY)
		fmt.Printf("Pacman: (%.1f, %.1f) -> Tile: (%d, %d)\n", 
			gameState.PacmanX, gameState.PacmanY, 
			int(gameState.PacmanX/TileSize), int(gameState.PacmanY/TileSize))
		
		// Check what's around the ghost
		currentTileX := int(g.X / TileSize)
		currentTileY := int(g.Y / TileSize)
		fmt.Printf("Surrounding tiles:\n")
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				tx := currentTileX + dx
				ty := currentTileY + dy
				if ty >= 0 && ty < len(gameState.Level) && tx >= 0 && tx < len(gameState.Level[0]) {
					fmt.Printf("%d ", gameState.Level[ty][tx])
				} else {
					fmt.Printf("X ")
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("==================\n")
	}
	// Update timers
	g.updateTimers()
	
	// Handle stuck state detection and recovery
	g.handleStuckState(gameState)
	
	// Update difficulty scaling
	g.updateDifficultyScaling(gameState)
	
	// Update mode based on timers and game state
	g.updateMode(gameState)
	
	// Update target based on current mode and ghost type
	g.UpdateTarget(gameState)

		
	// Test movement in all directions to see what's valid
	if g.PersonalityMode%120 == 0 { // Every 2 seconds
		fmt.Printf("Movement test for %s:\n", g.GhostType)
		directions := []string{"up", "down", "left", "right"}
		for _, dir := range directions {
			testX, testY := g.X, g.Y
			switch dir {
			case "up":
				testY -= g.Speed
			case "down":
				testY += g.Speed
			case "left":
				testX -= g.Speed
			case "right":
				testX += g.Speed
			}
			
			valid := g.isValidPosition(gameState, testX, testY)
			fmt.Printf("  %s: (%.1f,%.1f) -> %v\n", dir, testX, testY, valid)
		}
	}

	// Handle screen wrapping (tunnels)
	g.handleTunnels(gameState)
	
	// Check for collisions before moving
		g.moveToTarget(gameState)
	
	
	// Update personality counter
	g.PersonalityMode++
	
	
	// // Trigger sound effects for state changes
	// if g.soundEnabled {
	// 	g.checkSoundTriggers()
	// }
}

// Simplified movement function for debugging
func (g *Ghost) moveToTargetSimple(gameState *GameStateStruct) {
	// Very simple movement - just try to move towards target
	currentTileX := int((g.X + float64(g.Size)/2) / TileSize)
	currentTileY := int((g.Y + float64(g.Size)/2) / TileSize)
	
	fmt.Printf("Ghost %s at tile (%d,%d) targeting (%d,%d)\n", 
		g.GhostType, currentTileX, currentTileY, g.TargetX, g.TargetY)
	
	// Choose direction
	dx := g.TargetX - currentTileX
	dy := g.TargetY - currentTileY
	
	var moveX, moveY float64 = 0, 0
	
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		if dx > 0 {
			moveX = g.Speed
			g.Direction = "right"
		} else {
			moveX = -g.Speed
			g.Direction = "left"
		}
	} else {
		if dy > 0 {
			moveY = g.Speed
			g.Direction = "down"
		} else {
			moveY = -g.Speed
			g.Direction = "up"
		}
	}
	
	// Test the move
	newX := g.X + moveX
	newY := g.Y + moveY
	
	// Simple collision check - just check center point
	centerX := newX + float64(g.Size)/2
	centerY := newY + float64(g.Size)/2
	
	tileX := int(centerX / TileSize)
	tileY := int(centerY / TileSize)
	
	// Bounds check
	if tileY < 0 || tileY >= len(gameState.Level) || 
		tileX < 0 || tileX >= len(gameState.Level[0]) {
		fmt.Printf("  Move blocked: out of bounds (%d,%d)\n", tileX, tileY)
		return
	}
	
	// Wall check
	if gameState.Level[tileY][tileX] == TileWall {
		fmt.Printf("  Move blocked: wall at (%d,%d) = %d\n", 
			tileX, tileY, gameState.Level[tileY][tileX])
		return
	}
	
	// Move is valid
	g.X = newX
	g.Y = newY
	fmt.Printf("  Moved to (%.1f,%.1f)\n", g.X, g.Y)
}
// Check what your tile constants are
func (g *Ghost) debugTileConstants(gameState *GameStateStruct) {
	// Print some sample tiles to see what values you're using
	fmt.Println("Sample tile values:")
	for y := 0; y < min(5, len(gameState.Level)); y++ {
		for x := 0; x < min(10, len(gameState.Level[0])); x++ {
			fmt.Printf("%d ", gameState.Level[y][x])
		}
		fmt.Println()
	}
}
// Helper function for Go versions that don't have min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Ghost) moveToTarget(gameState *GameStateStruct) {
    // Debug info
    if g.PersonalityMode%120 == 0 {
        fmt.Printf("Moving %s: pos(%.1f,%.1f) target(%d,%d) mode=%d\n",
                   g.GhostType, g.X, g.Y, g.TargetX, g.TargetY, g.Mode)
    }
    
    targetPixelX := float64(g.TargetX * TileSize)
    targetPixelY := float64(g.TargetY * TileSize)
    
    // Calculate direction to target
    dx := targetPixelX - g.X
    dy := targetPixelY - g.Y
    
    // If close enough to target, we're done
    if math.Abs(dx) < 4 && math.Abs(dy) < 4 {
        return
    }
    
    // Determine primary and secondary movement directions
    var primaryMove, secondaryMove struct {
        name string
        dx, dy float64
    }
    
    if math.Abs(dx) > math.Abs(dy) {
        // Horizontal movement is primary
        if dx > 0 {
            primaryMove = struct{name string; dx, dy float64}{"right", g.Speed, 0}
        } else {
            primaryMove = struct{name string; dx, dy float64}{"left", -g.Speed, 0}
        }
        
        if dy > 0 {
            secondaryMove = struct{name string; dx, dy float64}{"down", 0, g.Speed}
        } else {
            secondaryMove = struct{name string; dx, dy float64}{"up", 0, -g.Speed}
        }
    } else {
        // Vertical movement is primary
        if dy > 0 {
            primaryMove = struct{name string; dx, dy float64}{"down", 0, g.Speed}
        } else {
            primaryMove = struct{name string; dx, dy float64}{"up", 0, -g.Speed}
        }
        
        if dx > 0 {
            secondaryMove = struct{name string; dx, dy float64}{"right", g.Speed, 0}
        } else {
            secondaryMove = struct{name string; dx, dy float64}{"left", -g.Speed, 0}
        }
    }
    
    // Try primary movement first
    newX := g.X + primaryMove.dx
    newY := g.Y + primaryMove.dy
    
    if g.isValidPosition(gameState, newX, newY) {
        g.X = newX
        g.Y = newY
        g.Direction = primaryMove.name
        return
    }
    
    // Try secondary movement
    newX = g.X + secondaryMove.dx
    newY = g.Y + secondaryMove.dy
    
    if g.isValidPosition(gameState, newX, newY) {
        g.X = newX
        g.Y = newY
        g.Direction = secondaryMove.name
        return
    }
    
    // If both failed, try all 4 directions as backup
    allDirections := []struct{name string; dx, dy float64}{
        {"up", 0, -g.Speed},
        {"down", 0, g.Speed},
        {"left", -g.Speed, 0},
        {"right", g.Speed, 0},
    }
    
    for _, dir := range allDirections {
        newX := g.X + dir.dx
        newY := g.Y + dir.dy
        
        if g.isValidPosition(gameState, newX, newY) {
            g.X = newX
            g.Y = newY
            g.Direction = dir.name
            return
        }
    }
    
    // Last resort: try smaller movements to get unstuck
    for _, dir := range allDirections {
        newX := g.X + dir.dx * 0.1 // Very small movement
        newY := g.Y + dir.dy * 0.1
        
        if g.isValidPosition(gameState, newX, newY) {
            g.X = newX
            g.Y = newY
            g.Direction = dir.name
            fmt.Printf("Ghost %s unstuck with micro-movement %s\n", g.GhostType, dir.name)
            return
        }
    }
    
    fmt.Printf("Ghost %s completely stuck at (%.1f,%.1f)!\n", g.GhostType, g.X, g.Y)
}

func (g *Ghost) UpdateTarget(gameState *GameStateStruct) {
    oldTargetX, oldTargetY := g.TargetX, g.TargetY
    
    switch g.Mode {
    case ChaseMode:
        g.updateChaseTarget(gameState)
    case ScatterMode:
        g.TargetX = g.ScatterTarget[0]
        g.TargetY = g.ScatterTarget[1]
    case FrightenedMode:
        g.updateFrightenedTarget(gameState)
    case DeadMode:
        g.TargetX = GHOST_HOUSE_X
        g.TargetY = GHOST_HOUSE_Y
    case InHouseMode:
        if g.ReleaseTimer <= 0 {
            // Move toward exit
            g.TargetX = GHOST_HOUSE_X
            g.TargetY = GHOST_HOUSE_EXIT_Y
        } else {
            // Stay in house area - set a position inside the house
            g.TargetX = GHOST_HOUSE_X
            g.TargetY = GHOST_HOUSE_Y
        }
    }
    
    // Clamp target to valid bounds
    g.clampTarget(gameState)
    
    // Debug: Log target changes
    if (g.TargetX != oldTargetX || g.TargetY != oldTargetY) && g.PersonalityMode%60 == 0 {
        fmt.Printf("Ghost %s target changed: (%d,%d) -> (%d,%d) [mode=%d]\n",
                   g.GhostType, oldTargetX, oldTargetY, g.TargetX, g.TargetY, g.Mode)
    }
}

// Advanced chase AI for each ghost type
// func (g *Ghost) updateChaseTarget(gameState *GameStateStruct) {
// 	// Use player center for more accurate targeting
// 	pacmanTileX := int((gameState.PacmanX + 16) / TileSize) // Assuming player size ~32
// 	pacmanTileY := int((gameState.PacmanY + 16) / TileSize)

// 	switch g.GhostType {
// 	case "jogo": // Direct chaser
// 		g.TargetX = pacmanTileX
// 		g.TargetY = pacmanTileY
		
// 		// Cruise Elroy mode
// 		if gameState.DotsRemaining < 20 && g.CruiseElroyMode == 0 {
// 			g.CruiseElroyMode = 1
// 			g.Speed = g.BaseSpeed * 1.05
// 		} else if gameState.DotsRemaining < 10 && g.CruiseElroyMode == 1 {
// 			g.CruiseElroyMode = 2
// 			g.Speed = g.BaseSpeed * 1.1
// 		}
		
// 	case "sukuna": // Ambush 4 tiles ahead
// 		targetX, targetY := pacmanTileX, pacmanTileY
		
// 		switch gameState.PacmanDirection {
// 		case "up":
// 			targetY -= 4
// 			targetX -= 4 // Original bug behavior
// 		case "down":
// 			targetY += 4
// 		case "left":
// 			targetX -= 4
// 		case "right":
// 			targetX += 4
// 		}
		
	// 	g.TargetX = targetX
	// 	g.TargetY = targetY
		
	// case "kenjaku": // Complex flanking
	// 	targetX, targetY := pacmanTileX, pacmanTileY
	// 	switch gameState.PacmanDirection {
	// 	case "up":
	// 		targetY -= 2
	// 	case "down":
	// 		targetY += 2
	// 	case "left":
	// 		targetX -= 2
	// 	case "right":
	// 		targetX += 2
	// 	}
		
	// 	// Find jogo for vector calculation
	// 	var blinky *Ghost
	// 	for _, ghost := range gameState.Ghosts {
	// 		if ghost.GhostType == "jogo" {
	// 			blinky = ghost
	// 			break
	// 		}
	// 	}
		
	// 	if blinky != nil {
	// 		blinkyTileX := int((blinky.X + float64(blinky.Size)/2) / TileSize)
	// 		blinkyTileY := int((blinky.Y + float64(blinky.Size)/2) / TileSize)
			
	// 		vectorX := targetX - blinkyTileX
	// 		vectorY := targetY - blinkyTileY
			
	// 		g.TargetX = targetX + vectorX
	// 		g.TargetY = targetY + vectorY
	// 	} else {
	// 		g.TargetX = targetX
	// 		g.TargetY = targetY
	// 	}
		
	// case "mahito": // Shy behavior
	// 	distance := g.getDistance(float64(pacmanTileX), float64(pacmanTileY))
		
	// 	if distance > 8 {
	// 		g.TargetX = pacmanTileX
	// 		g.TargetY = pacmanTileY
	// 	} else {
	// 		g.TargetX = g.ScatterTarget[0]
// 			g.TargetY = g.ScatterTarget[1]
// 		}
// 	}
// }


func (g *Ghost) updateChaseTarget(gameState *GameStateStruct) {
	// Use player center for more accurate targeting
	pacmanTileX := int((gameState.PacmanX + 16) / TileSize)
	pacmanTileY := int((gameState.PacmanY + 16) / TileSize)

	switch g.GhostType {
	case "jogo": 
		// Direct chase
		g.TargetX = pacmanTileX
		g.TargetY = pacmanTileY
		fmt.Printf("Ghost jogo chase mode: target=(%d,%d) pacman=(%d,%d)\n", 
			g.TargetX, g.TargetY, pacmanTileX, pacmanTileY)
		
	case "sukuna": 
		// 2 tiles ahead (simplified)
		targetX, targetY := pacmanTileX, pacmanTileY
		switch gameState.PacmanDirection {
		case "up":
			targetY -= 2
		case "down":
			targetY += 2
		case "left":
			targetX -= 2
		case "right":
			targetX += 2
		}
		g.TargetX = targetX
		g.TargetY = targetY
		
	case "kenjaku": 
		// Just target pacman for now (simplified)
		g.TargetX = pacmanTileX
		g.TargetY = pacmanTileY
		
	case "mahito": 
		// Just target pacman for now (simplified)
		g.TargetX = pacmanTileX
		g.TargetY = pacmanTileY
	}
}

// Frightened mode target - random movement with bias away from Pacman
func (g *Ghost) updateFrightenedTarget(gameState *GameStateStruct) {
	currentTileX := int(g.X / TileSize)
	currentTileY := int(g.Y / TileSize)
	
	// Get valid directions
	validTargets := make([]Node, 0)
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	
	for _, dir := range directions {
		newX := currentTileX + dir[0]
		newY := currentTileY + dir[1]
		
		// Use the tile-based validation function
		if g.isValidTilePosition(gameState, newX, newY) {
			// Bias away from Pacman
			pacmanDistance := math.Sqrt(math.Pow(float64(newX)-gameState.PacmanX/TileSize, 2) + 
				math.Pow(float64(newY)-gameState.PacmanY/TileSize, 2))
			
			// Add multiple entries for positions further from Pacman
			weight := int(pacmanDistance) + 1
			for i := 0; i < weight; i++ {
				validTargets = append(validTargets, Node{X: newX, Y: newY})
			}
		}
	}
	
	if len(validTargets) > 0 {
		target := validTargets[rand.Intn(len(validTargets))]
		g.TargetX = target.X
		g.TargetY = target.Y
	}
}


// // Move towards target using pathfinding from ghostAI.go
// func (g *Ghost) moveToTarget(gameState *GameStateStruct) {
// 	// Try to move in small steps toward target
// 	targetPixelX := float64(g.TargetX * TileSize)
// 	targetPixelY := float64(g.TargetY * TileSize)
	
// 	// Calculate direction to target
// 	dx := targetPixelX - g.X
// 	dy := targetPixelY - g.Y
	
// 	// If very close to target, we're done
// 	if math.Abs(dx) < 2 && math.Abs(dy) < 2 {
// 		return
// 	}
	
// 	// Normalize direction
// 	distance := math.Sqrt(dx*dx + dy*dy)
// 	if distance == 0 {
// 		return
// 	}
	
// 	moveX := (dx / distance) * g.Speed
// 	moveY := (dy / distance) * g.Speed
	
// 	// Try direct movement first
// 	newX := g.X + moveX
// 	newY := g.Y + moveY
	
// 	if g.isValidPosition(gameState, newX, newY) {
// 		g.X = newX
// 		g.Y = newY
		
// 		// Update direction for animation
// 		if math.Abs(moveX) > math.Abs(moveY) {
// 			if moveX > 0 {
// 				g.Direction = "right"
// 			} else {
// 				g.Direction = "left"
// 			}
// 		} else {
// 			if moveY > 0 {
// 				g.Direction = "down"
// 			} else {
// 				g.Direction = "up"
// 			}
// 		}
// 		return
// 	}
	
// 	// If direct movement fails, try axis-aligned movement
// 	// Try horizontal first
// 	if math.Abs(dx) > math.Abs(dy) {
// 		if dx > 0 && g.isValidPosition(gameState, g.X + g.Speed, g.Y) {
// 			g.X += g.Speed
// 			g.Direction = "right"
// 			return
// 		} else if dx < 0 && g.isValidPosition(gameState, g.X - g.Speed, g.Y) {
// 			g.X -= g.Speed
// 			g.Direction = "left"
// 			return
// 		}
// 	}
	
// 	// Try vertical
// 	if dy > 0 && g.isValidPosition(gameState, g.X, g.Y + g.Speed) {
// 		g.Y += g.Speed
// 		g.Direction = "down"
// 		return
// 	} else if dy < 0 && g.isValidPosition(gameState, g.X, g.Y - g.Speed) {
// 		g.Y -= g.Speed
// 		g.Direction = "up"
// 		return
// 	}
	
// 	// If horizontal failed, try the other horizontal direction
// 	if math.Abs(dx) <= math.Abs(dy) {
// 		if dx > 0 && g.isValidPosition(gameState, g.X + g.Speed, g.Y) {
// 			g.X += g.Speed
// 			g.Direction = "right"
// 			return
// 		} else if dx < 0 && g.isValidPosition(gameState, g.X - g.Speed, g.Y) {
// 			g.X -= g.Speed
// 			g.Direction = "left"
// 			return
// 		}
// 	}
	
// 	// Last resort: try any valid direction
// 	directions := []struct{
// 		name string
// 		dx, dy float64
// 	}{
// 		{"up", 0, -g.Speed},
// 		{"down", 0, g.Speed},
// 		{"left", -g.Speed, 0},
// 		{"right", g.Speed, 0},
// 	}
	
// 	for _, dir := range directions {
// 		if g.isValidPosition(gameState, g.X + dir.dx, g.Y + dir.dy) {
// 			g.X += dir.dx
// 			g.Y += dir.dy
// 			g.Direction = dir.name
// 			return
// 		}
// 	}
// }

// // Choose random valid direction when pathfinding fails
// func (g *Ghost) chooseRandomDirection(gameState *GameStateStruct) {
// 	currentTileX := int(g.X / TILE_SIZE)
// 	currentTileY := int(g.Y / TILE_SIZE)
	
// 	validDirections := []string{}
// 	directions := []struct{
// 		name string
// 		dx, dy int
// 	}{
// 		{"up", 0, -1},
// 		{"down", 0, 1},
// 		{"left", -1, 0},
// 		{"right", 1, 0},
// 	}
	
// 	for _, dir := range directions {
// 		newX := currentTileX + dir.dx
// 		newY := currentTileY + dir.dy
		
// 		if g.isValidPosition(gameState, newX, newY) {
// 			// Don't reverse unless forced
// 			if g.getReverseDirection(dir.name) != g.Direction || len(validDirections) == 0 {
// 				validDirections = append(validDirections, dir.name)
// 			}
// 		}
// 	}
	
// 	if len(validDirections) > 0 {
// 		g.Direction = validDirections[rand.Intn(len(validDirections))]
		
// 		// Create a simple path in the chosen direction
// 		targetX := currentTileX
// 		targetY := currentTileY
// 		switch g.Direction {
// 		case "up":
// 			targetY--
// 		case "down":
// 			targetY++
// 		case "left":
// 			targetX--
// 		case "right":
// 			targetX++
// 		}
		
// 		g.Path = []Node{
// 			{X: currentTileX, Y: currentTileY},
// 			{X: targetX, Y: targetY},
// 		}
// 		g.PathIndex = 0
// 	}
// }

// Get reverse direction
func (g *Ghost) getReverseDirection(direction string) string {
	switch direction {
	case "up":
		return "down"
	case "down":
		return "up"
	case "left":
		return "right"
	case "right":
		return "left"
	default:
		return ""
	}
}

// Follow the calculated path (uses Node from ghostAI.go)
func (g *Ghost) followPath(gameState *GameStateStruct) {
	if len(g.Path) == 0 || g.PathIndex >= len(g.Path) {
		return
	}
	
	// Get next waypoint
	if g.PathIndex < len(g.Path)-1 {
		g.PathIndex++
	}
	
	nextNode := g.Path[g.PathIndex]
	targetPixelX := float64(nextNode.X * TileSize)
	targetPixelY := float64(nextNode.Y * TileSize)
	
	// Move towards next waypoint
	dx := targetPixelX - g.X
	dy := targetPixelY - g.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	if distance > 1 {
		// Normalize and apply speed
		moveX := (dx / distance) * g.Speed
		moveY := (dy / distance) * g.Speed
		
		g.X += moveX
		g.Y += moveY
		
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
	}
}

// SetFrightened activates frightened mode
func (g *Ghost) SetFrightened(duration int) {
	if g.Mode != DeadMode && g.Mode != InHouseMode {
		g.Mode = FrightenedMode
		g.FrightTimer = duration
		g.Speed = g.BaseSpeed * 0.5
		
		// Reverse direction immediately
		g.reverseDirection()
	}
}


// CanBeEaten checks if ghost can be eaten in current state
func (g *Ghost) CanBeEaten() bool {
	return g.Mode == FrightenedMode && g.FrightTimer > 60 // Give some warning time
}

// SetVisible controls ghost visibility
func (g *Ghost) SetVisible(visible bool) {
	g.Visible = visible
}

// ResetMode resets ghost to default mode
func (g *Ghost) ResetMode() {
	switch g.GhostType {
	case "jogo":
		g.Mode = ChaseMode
		g.ChaseTimer = 1200
	default:
		g.Mode = ScatterMode
		g.ScatterTimer = 420 // 7 seconds
	}
	
	g.Speed = g.BaseSpeed
	g.FrightTimer = 0
	g.CruiseElroyMode = 0
}
// CollideWithPlayer handles collision with player
func (g *Ghost) CollideWithPlayer(playerX, playerY float64) string {
	// Calculate distance between centers
	ghostCenterX := g.X + float64(g.Size)/2
	ghostCenterY := g.Y + float64(g.Size)/2
	
	// Assume player is also centered (adjust if needed)
	playerCenterX := playerX + 16 // Assuming player size ~32, so center at +16
	playerCenterY := playerY + 16
	
	distance := math.Sqrt(math.Pow(ghostCenterX-playerCenterX, 2) + 
						 math.Pow(ghostCenterY-playerCenterY, 2))
	
	// Collision threshold based on combined sizes
	threshold := float64(g.Size)/2 + 16 // Half ghost size + half player size
	
	if distance < threshold {
    	fmt.Printf("Ghost %s collision! Mode: %d, Distance: %.2f, Threshold: %.2f\n", 
                   g.GhostType, g.Mode, distance, threshold)
		switch g.Mode {
		case FrightenedMode:
			if g.CanBeEaten() {
				g.Mode = DeadMode
				g.Speed = g.BaseSpeed * 2
				g.Visible = false // Hide dead ghost temporarily
                fmt.Printf("Ghost %s eaten!\n", g.GhostType)
				return "ghost_eaten"
			}
		case ChaseMode, ScatterMode:
    		fmt.Printf("Player caught by ghost %s!\n", g.GhostType)
			return "player_caught"
		}
	}
	
	return "no_collision"
}

// Helper function to get smoother movement
func (g *Ghost) smoothToGrid() {
	// Gradually align ghost to grid for smoother movement
	tileSize := float64(TileSize)
	
	// Calculate ideal grid-aligned position
	idealX := math.Round(g.X/tileSize) * tileSize
	idealY := math.Round(g.Y/tileSize) * tileSize
	
	// Smoothly move towards grid alignment (optional)
	alignmentSpeed := 0.1
	if math.Abs(g.X - idealX) < 2 {
		g.X += (idealX - g.X) * alignmentSpeed
	}
	if math.Abs(g.Y - idealY) < 2 {
		g.Y += (idealY - g.Y) * alignmentSpeed
	}
}
// Compatibility method for old interface
func (g *Ghost) CollidesWith(playerX, playerY float64, playerSize int) bool {
	result := g.CollideWithPlayer(playerX, playerY)
	return result != "no_collision"
}

func (g *Ghost) Draw(screen *ebiten.Image) {
    if !g.Visible || g.Image == nil {
        return
    }
    
    op := &ebiten.DrawImageOptions{}
    
    // Frightened mode visual effects
    if g.Mode == FrightenedMode {
        // Blue color for frightened mode
        if g.FrightTimer > 120 {
            // Solid blue
            op.ColorM.Scale(0.2, 0.2, 1.0, 1.0) // Blue tint
        } else {
            // Flashing white/blue when fright mode is ending
            if (g.FrightTimer/10)%2 == 0 {
                op.ColorM.Scale(0.8, 0.8, 1.0, 1.0) // Light blue/white
            } else {
                op.ColorM.Scale(0.2, 0.2, 1.0, 1.0) // Blue
            }
        }
    }
    
    // Death animation - darker appearance with eyes only
    if g.Mode == DeadMode {
        op.ColorM.Scale(0.3, 0.3, 0.3, 0.8) // Dark and semi-transparent
    }
    
    // Scale factor to make image exactly g.Size × g.Size
    scaleX := float64(g.Size) / float64(g.Image.Bounds().Dx())
    scaleY := float64(g.Size) / float64(g.Image.Bounds().Dy())
    op.GeoM.Scale(scaleX, scaleY)

    // Position after scaling
    op.GeoM.Translate(g.X, g.Y)
    screen.DrawImage(g.Image, op)
    
    // Debug: Draw hitbox in frightened mode
    if g.Mode == FrightenedMode {
        // Draw a blue rectangle outline to show hitbox
        ebitenutil.DrawRect(screen, g.X, g.Y, float64(g.Size), 2, color.RGBA{0, 0, 255, 128}) // Top
        ebitenutil.DrawRect(screen, g.X, g.Y, 2, float64(g.Size), color.RGBA{0, 0, 255, 128}) // Left  
        ebitenutil.DrawRect(screen, g.X+float64(g.Size)-2, g.Y, 2, float64(g.Size), color.RGBA{0, 0, 255, 128}) // Right
        ebitenutil.DrawRect(screen, g.X, g.Y+float64(g.Size)-2, float64(g.Size), 2, color.RGBA{0, 0, 255, 128}) // Bottom
    }
}

func (g *Ghost) debugSurroundingTiles(gameState *GameStateStruct) {
    currentTileX := int(g.X / TileSize)
    currentTileY := int(g.Y / TileSize)
    
    fmt.Printf("=== TILE DEBUG for %s at (%.1f, %.1f) = tile (%d, %d) ===\n", 
               g.GhostType, g.X, g.Y, currentTileX, currentTileY)
    
    // Check 3x3 area around ghost
    for dy := -1; dy <= 1; dy++ {
        for dx := -1; dx <= 1; dx++ {
            tx := currentTileX + dx
            ty := currentTileY + dy
            
            if ty >= 0 && ty < len(gameState.Level) && tx >= 0 && tx < len(gameState.Level[0]) {
                tileType := gameState.Level[ty][tx]
                symbol := "?"
                switch tileType {
                case TileEmpty:
                    symbol = "."
                case TileWall:
                    symbol = "#"
                case TilePellet:
                    symbol = "o"
                case TilePowerPellet:
                    symbol = "O"
                }
                fmt.Printf("%s ", symbol)
            } else {
                fmt.Printf("X ")
            }
        }
        fmt.Printf("\n")
    }
    fmt.Printf("Current tile type: %d\n", gameState.Level[currentTileY][currentTileX])
    fmt.Printf("===================\n")
}

// Tunnel handling for screen wrapping
func (g *Ghost) handleTunnels(gameState *GameStateStruct) {
	screenWidth := float64(len(gameState.Level[0]) * TileSize)
	ghostSize := float64(g.Size)
	
	// Left tunnel - ghost completely off screen
	if g.X + ghostSize < 0 {
		g.X = screenWidth
	} 
	// Right tunnel - ghost completely off screen
	if g.X > screenWidth {
		g.X = -ghostSize
	}
}
// Collision detection with walls and boundaries
func (g *Ghost) checkCollisions(gameState *GameStateStruct) bool {
	nextX := g.X + g.getDirectionVector().X
	nextY := g.Y + g.getDirectionVector().Y
	
	tileX := int(nextX / 16)
	tileY := int(nextY / 16)
	
	if tileY < 0 || tileY >= len(gameState.Level) || 
		tileX < 0 || tileX >= len(gameState.Level[0]) {
		return false
	}
	
	return gameState.Level[tileY][tileX] != 1
}

// Get movement vector for current direction
func (g *Ghost) getDirectionVector() struct{ X, Y float64 } {
	switch g.Direction {
	case "up":
		return struct{ X, Y float64 }{0, -g.Speed}
	case "down":
		return struct{ X, Y float64 }{0, g.Speed}
	case "left":
		return struct{ X, Y float64 }{-g.Speed, 0}
	case "right":
		return struct{ X, Y float64 }{g.Speed, 0}
	default:
		return struct{ X, Y float64 }{0, 0}
	}
}

// Advanced stuck detection and recovery
func (g *Ghost) handleStuckState(gameState *GameStateStruct) {
	currentTile := [2]int{int(g.X / TileSize), int(g.Y / TileSize)}
	
	
	if g.LastPosition == currentTile {
		g.StuckCounter++
	} else {
		g.StuckCounter = 0
		g.LastPosition = currentTile
	}
	
	// If stuck for too long, force a new path or random direction
	if g.StuckCounter > 30 {
		g.Path = nil // Force new pathfinding
		g.PathIndex = 0
		g.StuckCounter = 0
		
		// Emergency random direction
		directions := []string{"up", "down", "left", "right"}
		g.Direction = directions[rand.Intn(len(directions))]
	}
}


// Performance optimization - only update AI every few frames
func (g *Ghost) shouldUpdateAI() bool {
	return true// Update every 4 frames
}

// Advanced difficulty scaling
func (g *Ghost) updateDifficultyScaling(gameState *GameStateStruct) {
	baseSpeed := 0.8
	
	// Increase speed and reduce scatter time as level increases
	speedMultiplier := 1.0 + (float64(gameState.CurrentLevel-1) * 0.1)
	g.BaseSpeed = baseSpeed * speedMultiplier
	
	if g.Mode != FrightenedMode {
		g.Speed = g.BaseSpeed
	}
	
	// Reduce fright mode duration on higher levels
	maxFrightTime := 600 - (gameState.CurrentLevel-1)*30
	if maxFrightTime < 120 {
		maxFrightTime = 120
	}
	
	if g.Mode == FrightenedMode && g.FrightTimer > maxFrightTime {
		g.FrightTimer = maxFrightTime
	}
}

// Debug visualization for AI development
func (g *Ghost) drawAIDebug(screen *ebiten.Image) {
	// Draw target position
	// Draw current path
	// Draw AI state information
	// This would use your text/debug rendering system
}

// Sound effect triggers
func (g *Ghost) triggerSoundEffect(effect string) {
	// Integration point for your audio system
	switch effect {
	case "ghost_eaten":
		// Play ghost eaten sound
	case "mode_change":
		// Play mode change sound
	case "fright_warning":
		// Play warning sound when fright mode ending
	}
}

// Save/Load ghost state for game persistence
func (g *Ghost) SaveState() map[string]interface{} {
	return map[string]interface{}{
		"x":               g.X,
		"y":               g.Y,
		"direction":       g.Direction,
		"mode":            g.Mode,
		"frightTimer":     g.FrightTimer,
		"cruiseElroyMode": g.CruiseElroyMode,
		"pathIndex":       g.PathIndex,
	}
}

func (g *Ghost) LoadState(state map[string]interface{}) {
	if x, ok := state["x"].(float64); ok {
		g.X = x
	}
	if y, ok := state["y"].(float64); ok {
		g.Y = y
	}
	// ... load other state variables
}

// Network synchronization for multiplayer
func (g *Ghost) GetNetworkState() []byte {
	// Serialize minimal state for network transmission
	// Return compressed ghost state
	return nil
}

func (g *Ghost) SetNetworkState(data []byte) {
	// Deserialize and apply network state
}

// Add these fields to Ghost struct
var lastPosition [2]int
var drawDebugInfo bool

// Helper functions
func (g *Ghost) updateTimers() {
	if g.FrightTimer > 0 {
		g.FrightTimer--
	}
	if g.ScatterTimer > 0 {
		g.ScatterTimer--
	}
	if g.ChaseTimer > 0 {
		g.ChaseTimer--
	}
	if g.ReleaseTimer > 0 {
		g.ReleaseTimer--
	}
}

func (g *Ghost) updateMode(gameState *GameStateStruct) {
    oldMode := g.Mode
    
    // Global fright mode overrides everything except dead mode
    if gameState.FrightModeActive && g.Mode != DeadMode {
        if g.Mode != FrightenedMode {
            g.SetFrightened(600)
        }
        return
    }
    
    // Handle release from ghost house
    if g.Mode == InHouseMode {
        if g.ReleaseTimer <= 0 {
            // Released! Start in scatter mode
            g.Mode = ScatterMode
            g.ScatterTimer = 420 // 7 seconds
            if oldMode != g.Mode {
                fmt.Printf("Ghost %s released from house, entering scatter mode\n", g.GhostType)
            }
        } else {
            // Still waiting to be released
            g.ReleaseTimer--
            return
        }
    }
    
    // Normal mode transitions
    if g.FrightTimer <= 0 && g.Mode == FrightenedMode {
        g.Mode = ScatterMode // Start with scatter after fright
        g.ScatterTimer = 420
        g.Speed = g.BaseSpeed
        if oldMode != g.Mode {
            fmt.Printf("Ghost %s exiting fright mode, entering scatter\n", g.GhostType)
        }
    }
     if g.ChaseTimer <= 0 && g.Mode == ChaseMode {
        g.Mode = ScatterMode
        g.ScatterTimer = 420
        if oldMode != g.Mode {
            fmt.Printf("Ghost %s switching to scatter mode\n", g.GhostType)
        }
    }
    
    if g.ScatterTimer <= 0 && g.Mode == ScatterMode {
        g.Mode = ChaseMode
        g.ChaseTimer = 1200
        if oldMode != g.Mode {
            fmt.Printf("Ghost %s switching to chase mode\n", g.GhostType)
        }
    }
}

func (g *Ghost) isValidPosition(gameState *GameStateStruct, pixelX, pixelY float64) bool {
    // Check all four corners of the ghost, not just center
    margin := 1.0 // Small margin to prevent wall clipping
    size:=float64(g.Size)
    checkPoints := [][2]float64{
        {pixelX + size/2, pixelY + size/2},           // Center (most important)
        {pixelX + margin, pixelY + margin},           // Top-left
        {pixelX + size - margin, pixelY + margin},    // Top-right
        {pixelX + margin, pixelY + size - margin},    // Bottom-left
        {pixelX + size - margin, pixelY + size - margin}, // Bottom-right
    }
    
    for i, point := range checkPoints {
        tileX := int(point[0] / TileSize)
        tileY := int(point[1] / TileSize)
        
        // Bounds check
        if tileY < 0 || tileY >= len(gameState.Level) || 
           tileX < 0 || tileX >= len(gameState.Level[0]) {
            return false
        }
        
        // Wall check
        if gameState.Level[tileY][tileX] == TileWall {
            // For center point, this is definitely invalid
            if i == 0 {
                return false
            }
            // For corner points, allow some wall overlap if center is clear
            // This makes movement more forgiving
            continue
        }
    }
    
    return true
}


    // Helper function for tile-based position checking
func (g *Ghost) isValidTilePosition(gameState *GameStateStruct, tileX, tileY int) bool {
	// Convert tile coordinates to pixel coordinates
	pixelX := float64(tileX * TileSize)
	pixelY := float64(tileY * TileSize)
	return g.isValidPosition(gameState, pixelX, pixelY)
}

func (g *Ghost) isValidPositionTile(gameState *GameStateStruct, x, y int) bool {
	return g.isValidTilePosition(gameState, x, y)
}
func (g *Ghost) clampTarget(gameState *GameStateStruct) {
	if g.TargetX < 0 {
		g.TargetX = 0
	} else if g.TargetX >= len(gameState.Level[0]) {
		g.TargetX = len(gameState.Level[0]) - 1
	}
	
	if g.TargetY < 0 {
		g.TargetY = 0
	} else if g.TargetY >= len(gameState.Level) {
		g.TargetY = len(gameState.Level) - 1
	}
}

func (g *Ghost) predictPacmanPosition(gameState *GameStateStruct, steps int) (int, int) {
	x := gameState.PacmanX / TileSize
	y := gameState.PacmanY / TileSize
	
	for i := 0; i < steps; i++ {
		switch gameState.PacmanDirection {
		case "up":
			y--
		case "down":
			y++
		case "left":
			x--
		case "right":
			x++
		}
		
		// Clamp to bounds and check for walls
		tileX := int(math.Max(0, math.Min(float64(len(gameState.Level[0])-1), x)))
		tileY := int(math.Max(0, math.Min(float64(len(gameState.Level)-1), y)))
		
		if gameState.Level[tileY][tileX] == TileWall {
			break
		}
	}
	
	return int(x), int(y)
}

func (g *Ghost) isPacmanCornered(gameState *GameStateStruct) bool {
	pacmanX := int(gameState.PacmanX / TileSize)
	pacmanY := int(gameState.PacmanY / TileSize)
	
	validMoves := 0
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	
	for _, dir := range directions {
		newX := pacmanX + dir[0]
		newY := pacmanY + dir[1]
		if g.isValidTilePosition(gameState, newX, newY) {
			validMoves++
		}
	}
	
	return validMoves <= 1
}

func (g *Ghost) getEscapeRoutes(gameState *GameStateStruct, x, y int) []Node {
	routes := make([]Node, 0)
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	
	for _, dir := range directions {
		newX := x + dir[0]
		newY := y + dir[1]
		if g.isValidTilePosition(gameState, newX, newY) {
			routes = append(routes, Node{X: newX, Y: newY})
		}
	}
	
	return routes
}
func (g *Ghost) hasDirectPath(gameState *GameStateStruct, targetX, targetY int) bool {
	currentX := int(g.X / TileSize)
	currentY := int(g.Y / TileSize)
	
	path := findPath(gameState.Level, currentX, currentY, targetX, targetY)
	return len(path) > 0
}

func (g *Ghost) getFlankingPositions(gameState *GameStateStruct, targetX, targetY int) []Node {
	positions := make([]Node, 0)
	offsets := [][2]int{{-2, -2}, {2, -2}, {-2, 2}, {2, 2}, {-3, 0}, {3, 0}, {0, -3}, {0, 3}}
	
	for _, offset := range offsets {
		newX := targetX + offset[0]
		newY := targetY + offset[1]
		if g.isValidTilePosition(gameState, newX, newY) {
			positions = append(positions, Node{X: newX, Y: newY})
		}
	}
	
	return positions
}
func (g *Ghost) reverseDirection() {
	switch g.Direction {
	case "up":
		g.Direction = "down"
	case "down":
		g.Direction = "up"
	case "left":
		g.Direction = "right"
	case "right":
		g.Direction = "left"
	}
	g.PreviousDirection = g.Direction
}


// Helper function to get distance using ghost center
func (g *Ghost) getDistance(targetX, targetY float64) float64 {
	ghostCenterX := (g.X + float64(g.Size)/2) / TileSize
	ghostCenterY := (g.Y + float64(g.Size)/2) / TileSize
	dx := targetX - ghostCenterX
	dy := targetY - ghostCenterY
	return math.Sqrt(dx*dx + dy*dy)
}

func (g *Ghost) printLocalArea(gameState *GameStateStruct) {
    currentTileX := int(g.X / TileSize)
    currentTileY := int(g.Y / TileSize)
    
    fmt.Printf("=== AREA AROUND %s at (%d,%d) ===\n", g.GhostType, currentTileX, currentTileY)
    
    for dy := -3; dy <= 3; dy++ {
        for dx := -3; dx <= 3; dx++ {
            tx := currentTileX + dx
            ty := currentTileY + dy
            
            if ty >= 0 && ty < len(gameState.Level) && tx >= 0 && tx < len(gameState.Level[0]) {
                tile := gameState.Level[ty][tx]
                switch {
                case dx == 0 && dy == 0:
                    fmt.Printf("[%d]", tile) // Ghost position
                case tile == TileWall:
                    fmt.Printf(" # ")
                case tile == TileEmpty:
                    fmt.Printf(" . ")
                case tile == TilePellet:
                    fmt.Printf(" o ")
                case tile == TilePowerPellet:
                    fmt.Printf(" O ")
                default:
                    fmt.Printf(" ? ")
                }
            } else {
                fmt.Printf(" X ") // Out of bounds
            }
        }
        fmt.Printf("\n")
    }
    fmt.Printf("===============================\n")
}
