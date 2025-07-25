package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	tileSize     = 20
	mapWidth     = screenWidth / tileSize
	mapHeight    = screenHeight / tileSize
)

// Game represents the main game state
type Game struct {
	player Player
	ghosts []Ghost
	maze   [][]int
	dots   [][]bool
	score  int
}

// Player represents Pacman
type Player struct {
	x, y    float64
	dir     Direction
	nextDir Direction
	speed   float64
}

// Ghost represents a ghost enemy
type Ghost struct {
	x, y  float64
	dir   Direction
	color color.RGBA
	speed float64
}

// Direction represents movement direction
type Direction int

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

// NewGame initializes a new game
func NewGame() *Game {
	g := &Game{
		player: Player{
			x:     10 * tileSize,
			y:     10 * tileSize,
			dir:   None,
			speed: 2.0,
		},
		score: 0,
	}

	g.initMaze()
	g.initDots()
	g.initGhosts()

	return g
}

// initMaze creates a simple maze layout
func (g *Game) initMaze() {
	g.maze = make([][]int, mapHeight)
	for i := range g.maze {
		g.maze[i] = make([]int, mapWidth)
	}

	// Create walls around the border
	for x := 0; x < mapWidth; x++ {
		g.maze[0][x] = 1
		g.maze[mapHeight-1][x] = 1
	}
	for y := 0; y < mapHeight; y++ {
		g.maze[y][0] = 1
		g.maze[y][mapWidth-1] = 1
	}

	// Add some internal walls for maze structure
	for y := 2; y < mapHeight-2; y += 4 {
		for x := 2; x < mapWidth-2; x += 4 {
			g.maze[y][x] = 1
			g.maze[y][x+1] = 1
			g.maze[y+1][x] = 1
		}
	}
}

// initDots places dots throughout the maze
func (g *Game) initDots() {
	g.dots = make([][]bool, mapHeight)
	for i := range g.dots {
		g.dots[i] = make([]bool, mapWidth)
	}

	// Place dots in empty spaces
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if g.maze[y][x] == 0 {
				g.dots[y][x] = true
			}
		}
	}
}

// initGhosts creates the ghost entities
func (g *Game) initGhosts() {
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Red
		{255, 192, 203, 255}, // Pink
		{0, 255, 255, 255},   // Cyan
		{255, 165, 0, 255},   // Orange
	}

	g.ghosts = make([]Ghost, 4)
	for i := 0; i < 4; i++ {
		g.ghosts[i] = Ghost{
			x:     float64(5+i*2) * tileSize,
			y:     5 * tileSize,
			dir:   Right,
			color: colors[i],
			speed: 1.5,
		}
	}
}

// Update handles game logic each frame
func (g *Game) Update() error {
	g.handleInput()
	g.updatePlayer()
	g.updateGhosts()
	g.checkCollisions()
	return nil
}

// handleInput processes player input
func (g *Game) handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.nextDir = Up
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.nextDir = Down
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.nextDir = Left
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.nextDir = Right
	}
}

// updatePlayer handles player movement and direction changes
func (g *Game) updatePlayer() {
	// Try to change direction if possible
	if g.canMove(g.player.x, g.player.y, g.player.nextDir) {
		g.player.dir = g.player.nextDir
	}

	// Move player
	if g.canMove(g.player.x, g.player.y, g.player.dir) {
		switch g.player.dir {
		case Up:
			g.player.y -= g.player.speed
		case Down:
			g.player.y += g.player.speed
		case Left:
			g.player.x -= g.player.speed
		case Right:
			g.player.x += g.player.speed
		}
	}

	// Screen wrapping
	if g.player.x < 0 {
		g.player.x = screenWidth - tileSize
	} else if g.player.x >= screenWidth {
		g.player.x = 0
	}
}

// updateGhosts handles ghost movement with simple AI
func (g *Game) updateGhosts() {
	for i := range g.ghosts {
		ghost := &g.ghosts[i]

		// Simple AI: change direction at intersections
		if g.isAtIntersection(ghost.x, ghost.y) {
			directions := []Direction{Up, Down, Left, Right}
			for _, dir := range directions {
				if g.canMove(ghost.x, ghost.y, dir) && dir != g.oppositeDir(ghost.dir) {
					ghost.dir = dir
					break
				}
			}
		}

		// Move ghost
		if g.canMove(ghost.x, ghost.y, ghost.dir) {
			switch ghost.dir {
			case Up:
				ghost.y -= ghost.speed
			case Down:
				ghost.y += ghost.speed
			case Left:
				ghost.x -= ghost.speed
			case Right:
				ghost.x += ghost.speed
			}
		} else {
			// Try a different direction if blocked
			directions := []Direction{Up, Down, Left, Right}
			for _, dir := range directions {
				if g.canMove(ghost.x, ghost.y, dir) {
					ghost.dir = dir
					break
				}
			}
		}

		// Screen wrapping for ghosts
		if ghost.x < 0 {
			ghost.x = screenWidth - tileSize
		} else if ghost.x >= screenWidth {
			ghost.x = 0
		}
	}
}

// canMove checks if movement is possible in given direction
func (g *Game) canMove(x, y float64, dir Direction) bool {
	newX, newY := x, y
	margin := float64(tileSize) * 0.4

	switch dir {
	case Up:
		newY -= margin
	case Down:
		newY += margin
	case Left:
		newX -= margin
	case Right:
		newX += margin
	default:
		return false
	}

	// Allow wrapping at screen edges
	if newX < 0 || newY < 0 || newX >= screenWidth || newY >= screenHeight {
		return true
	}

	// Check for walls
	tileX := int(newX / tileSize)
	tileY := int(newY / tileSize)

	if tileX >= 0 && tileX < mapWidth && tileY >= 0 && tileY < mapHeight {
		return g.maze[tileY][tileX] == 0
	}

	return true
}

// isAtIntersection checks if position is at a maze intersection
func (g *Game) isAtIntersection(x, y float64) bool {
	tileX := int(x / tileSize)
	tileY := int(y / tileSize)

	if tileX <= 0 || tileX >= mapWidth-1 || tileY <= 0 || tileY >= mapHeight-1 {
		return false
	}

	paths := 0
	directions := []Direction{Up, Down, Left, Right}
	for _, dir := range directions {
		if g.canMove(x, y, dir) {
			paths++
		}
	}

	return paths > 2
}

// oppositeDir returns the opposite direction
func (g *Game) oppositeDir(dir Direction) Direction {
	switch dir {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return None
}

// checkCollisions handles dot collection and ghost collisions
func (g *Game) checkCollisions() {
	playerTileX := int(g.player.x / tileSize)
	playerTileY := int(g.player.y / tileSize)

	// Check dot collection
	if playerTileX >= 0 && playerTileX < mapWidth && playerTileY >= 0 && playerTileY < mapHeight {
		if g.dots[playerTileY][playerTileX] {
			g.dots[playerTileY][playerTileX] = false
			g.score += 10
		}
	}

	// Check ghost collisions
	for _, ghost := range g.ghosts {
		distance := math.Sqrt(math.Pow(g.player.x-ghost.x, 2) + math.Pow(g.player.y-ghost.y, 2))
		if distance < tileSize*0.8 {
			// Reset player position on collision
			g.player.x = 10 * tileSize
			g.player.y = 10 * tileSize
			g.player.dir = None
		}
	}
}

// Draw renders the game to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with black background
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw all game elements
	g.drawMaze(screen)
	g.drawDots(screen)
	g.drawPlayer(screen)
	g.drawGhosts(screen)
	g.drawUI(screen)
}

// drawMaze renders the maze walls
func (g *Game) drawMaze(screen *ebiten.Image) {
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if g.maze[y][x] == 1 {
				vector.DrawFilledRect(screen,
					float32(x*tileSize), float32(y*tileSize),
					float32(tileSize), float32(tileSize),
					color.RGBA{0, 0, 255, 255}, false)
			}
		}
	}
}

// drawDots renders the collectible dots
func (g *Game) drawDots(screen *ebiten.Image) {
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if g.dots[y][x] {
				centerX := float32(x*tileSize + tileSize/2)
				centerY := float32(y*tileSize + tileSize/2)
				vector.DrawFilledCircle(screen, centerX, centerY, 3,
					color.RGBA{255, 255, 0, 255}, false)
			}
		}
	}
}

// drawPlayer renders Pacman
func (g *Game) drawPlayer(screen *ebiten.Image) {
	centerX := float32(g.player.x + tileSize/2)
	centerY := float32(g.player.y + tileSize/2)
	radius := float32(tileSize/2 - 2)

	// Draw Pacman's body
	vector.DrawFilledCircle(screen, centerX, centerY, radius,
		color.RGBA{255, 255, 0, 255}, false)

	// Draw mouth based on direction
	if g.player.dir != None {
		mouthSize := float32(4)
		switch g.player.dir {
		case Right:
			vector.DrawFilledRect(screen, centerX, centerY-mouthSize/2,
				radius, mouthSize, color.RGBA{0, 0, 0, 255}, false)
		case Left:
			vector.DrawFilledRect(screen, centerX-radius, centerY-mouthSize/2,
				radius, mouthSize, color.RGBA{0, 0, 0, 255}, false)
		case Up:
			vector.DrawFilledRect(screen, centerX-mouthSize/2, centerY-radius,
				mouthSize, radius, color.RGBA{0, 0, 0, 255}, false)
		case Down:
			vector.DrawFilledRect(screen, centerX-mouthSize/2, centerY,
				mouthSize, radius, color.RGBA{0, 0, 0, 255}, false)
		}
	}
}

// drawGhosts renders the ghost enemies
func (g *Game) drawGhosts(screen *ebiten.Image) {
	for _, ghost := range g.ghosts {
		centerX := float32(ghost.x + tileSize/2)
		centerY := float32(ghost.y + tileSize/2)
		radius := float32(tileSize/2 - 2)

		// Draw ghost body
		vector.DrawFilledCircle(screen, centerX, centerY, radius,
			ghost.color, false)

		// Draw eyes
		eyeSize := float32(3)
		vector.DrawFilledCircle(screen, centerX-radius/3, centerY-radius/3,
			eyeSize, color.RGBA{255, 255, 255, 255}, false)
		vector.DrawFilledCircle(screen, centerX+radius/3, centerY-radius/3,
			eyeSize, color.RGBA{255, 255, 255, 255}, false)
	}
}

// drawUI renders the user interface elements
func (g *Game) drawUI(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrint(screen, scoreText)
}

// Layout defines the screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pacman Game")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
