package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "fmt"
    "image/color"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
    
)

var glowColor = color.RGBA{255, 50, 50, 255} // bright red center
var outerGlowColor = color.RGBA{255, 50, 50, 100} // transparent red for glow

type GameState int

const(
    StateMenu GameState = iota
    StatePlaying
    StateGameOver
    StatePaused
)

type Game struct{
    Player *Player
    Ghosts []*Ghost
    Pellet []Pellet
    menuUI *UIPage
    //showMenu bool
    lives    int
    playerStartX float64
    playerStartY float64
    State  GameState //main state variable
    pelletCount  int
    powerPelletActive bool
	powerPelletTimer int
    logoImg *ebiten.Image
    characterGif *ebiten.Image
    bgTexture *ebiten.Image
}

const TileSize =32

func NewGame() *Game {
    playerStartX := float64(1 * TileSize)
    playerStartY := float64(1 * TileSize)
        g:= &Game{
        Player: NewPlayer(playerStartX,playerStartY,"assets/player.png"),
		menuUI: NewUIPage(),
        State: StateMenu,
        lives:  3,
        playerStartX: playerStartX,
        playerStartY: playerStartY,
        Ghosts: []*Ghost{
		NewGhost(13*TileSize, 13*TileSize, "assets/jogo.png", "jogo", 55),
    	NewGhost(12*TileSize, 13*TileSize, "assets/sakuna.png", "sukuna", 55),    
    	NewGhost(14*TileSize, 13*TileSize, "assets/kenjaku.png", "kenjaku", 55),  
    	NewGhost(13*TileSize, 14*TileSize, "assets/mahito.png", "mahito", 55),
        },


    }
	g.resetGhosts()

    InitPellets(level, TileSize)
	g.countPellets()
return g
}


func (g *Game) Update() error {
    switch g.State {
    case StateMenu:
        return g.updateMenu()
    case StatePlaying:
        return g.updateGame()
    case StateGameOver:
        return g.updateGameOver()
    case StatePaused:
        return g.updatePaused()
    }
    return nil
}
func (g *Game) updateMenu() error {
    if g.menuUI != nil {
        err := g.menuUI.Update()
        if err != nil {
            return err
        }
        
        // Handle menu selection
        if g.menuUI.IsEnterPressed() {
            switch g.menuUI.GetSelectedOption() {
            case 0: // START GAME
                fmt.Println("Starting game...")
                g.State = StatePlaying
                g.resetGame()
            case 1: // SETTINGS (you can implement later)
                fmt.Println("Settings selected")
                // For now, do nothing or show a message
            case 2: // GALLERY (you can implement later)  
                fmt.Println("Gallery selected")
                // For now, do nothing or show a message
            case 3: // EXIT
                return fmt.Errorf("quit game")
            }
        }
    }
    return nil
}


func (g *Game) updateGame() error {
    // Handle pause
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
        g.State = StatePaused
        return nil
    }
    
    // Update player
    if g.Player != nil {
        g.Player.Update(level, TileSize)
    }
    
    // Get Blinky's position for Inky's AI
    var blinkyX, blinkyY float64
    if len(g.Ghosts) > 0 {
        blinkyX = g.Ghosts[0].X
        blinkyY = g.Ghosts[0].Y
    }
    
    // Update power pellet timer
    if g.powerPelletActive {
        g.powerPelletTimer--
        if g.powerPelletTimer <= 0 {
            g.powerPelletActive = false
            // Reset all ghosts to normal mode
            for _, ghost := range g.Ghosts {
                ghost.ResetMode()
            }
        }
    }
    
    // Update each ghost
    // for _, ghost := range g.Ghosts {
    //     if ghost.Visible {
    //         ghost.Update(level, TileSize, g.Player.X, g.Player.Y, g.Player.Direction, blinkyX, blinkyY)
    //     }
    // }
    // Update each ghost
for i, ghost := range g.Ghosts {
    if ghost.Visible {
        fmt.Printf("Updating ghost %d (%s) at (%.1f, %.1f), player at (%.1f, %.1f)\n", 
                   i, ghost.Name, ghost.X, ghost.Y, g.Player.X, g.Player.Y)
        ghost.Update(level, TileSize, g.Player.X, g.Player.Y, g.Player.Direction, blinkyX, blinkyY)
    }
}
    
    // Check for collisions between player and ghosts
    for _, ghost := range g.Ghosts {
        if ghost.Visible && ghost.CollidesWith(g.Player.X, g.Player.Y, g.Player.Size) {
            if ghost.CanBeEaten() {
                // Ghost is frightened, player eats ghost
                g.Player.Score += 200
                ghost.SetVisible(false)
                
                // Respawn ghost after delay (you might want to add a respawn timer)
                go func(gh *Ghost) {
                    // In a real game, you'd use a proper timer system
                    // For now, just immediately respawn at start position
                    gh.X = g.playerStartX + float64(TileSize*12) // Ghost starting area
                    gh.Y = g.playerStartY + float64(TileSize*12)
                    gh.ResetMode()
                    gh.SetVisible(true)
                }(ghost)
                
            } else {
                // Player dies
                g.lives--
                g.resetPlayerPosition()
                
                if g.lives <= 0 {
                    g.State = StateGameOver
                    return nil
                }
                
                // Reset all ghosts to normal mode and positions
                g.resetGhosts()
            }
        }
    }
    
    // Check pellet collection
    g.checkPelletCollection()
    
    // Check win condition
    if g.pelletCount <= 0 {
        // Level complete - you could add level progression here
        g.resetGame()
    }
    
    return nil
}

func (g *Game) updatePaused() error {
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
        g.State = StatePlaying
    }
    return nil
}

func (g *Game) updateGameOver() error {
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
        g.resetGame()
        g.State = StateMenu
    }
    return nil
}

func (g *Game) checkPelletCollection() {
    playerTileX := int(g.Player.X) / TileSize
    playerTileY := int(g.Player.Y) / TileSize
    
    // Bounds checking
    if playerTileY < 0 || playerTileY >= len(level) || 
       playerTileX < 0 || playerTileX >= len(level[0]) {
        return
    }
    
    currentTile := level[playerTileY][playerTileX]
    
    switch currentTile {
    case TilePellet:
        level[playerTileY][playerTileX] = TileEmpty
        g.Player.Score += 10
        g.pelletCount--
        
        // Update player score if it has one
        if g.Player != nil {
            g.Player.Score = g.Player.Score
        }
        
    case TilePowerPellet:
        level[playerTileY][playerTileX] = TileEmpty
        g.Player.Score += 50
        g.pelletCount--
        
        // Activate power pellet mode
        g.powerPelletActive = true
        g.powerPelletTimer = 600 // 10 seconds at 60 FPS
        
        // Set all visible ghosts to frightened mode
        for _, ghost := range g.Ghosts {
            if ghost.Visible {
                ghost.SetFrightened(600)
            }
        }
        
        // Update player score
        if g.Player != nil {
            g.Player.Score = g.Player.Score
        }
    }
}

func (g *Game) resetPlayerPosition() {
    if g.Player != nil {
        g.Player.X = g.playerStartX
        g.Player.Y = g.playerStartY
        g.Player.Direction = "right" // or whatever default direction
    }
}

func (g *Game) resetGhosts() {
    // Force ghosts to proper starting positions in empty spaces
    ghostStartPositions := [][2]float64{
        {13 * TileSize, 13 * TileSize}, // jogo - center of ghost house
        {12 * TileSize, 13 * TileSize}, // sukuna - left of center
        {14 * TileSize, 13 * TileSize}, // kenjaku - right of center  
        {13 * TileSize, 11 * TileSize}, // mahito - above center
    }
    
    for i, ghost := range g.Ghosts {
        if i < len(ghostStartPositions) {
            ghost.X = ghostStartPositions[i][0]
            ghost.Y = ghostStartPositions[i][1]
            fmt.Printf("Reset ghost %s to position (%.1f, %.1f) = tile (%d, %d)\n", 
                      ghost.Name, ghost.X, ghost.Y, 
                      int(ghost.X)/TileSize, int(ghost.Y)/TileSize)
        }
        ghost.ResetMode()
        ghost.SetVisible(true)
    }
    
    g.powerPelletActive = false
    g.powerPelletTimer = 0
}
func (g *Game) resetGame() {
    //g.Player.Score = 0
    if g.Player != nil {
        g.Player.Score = 0
        g.resetPlayerPosition()
    } else {
        log.Println("⚠️ Warning: g.Player is nil during resetGame")
    }
    g.lives = 3
    g.resetPlayerPosition()
    g.resetGhosts()
    
    // Reset level pellets
    InitPellets(level, TileSize)
    g.countPellets()
    
    if g.Player != nil {
        g.Player.Score = 0
    }
}

func (g *Game) countPellets() {
    count := 0
    for _, row := range level {
        for _, tile := range row {
            if tile == TilePellet || tile == TilePowerPellet {
                count++
            }
        }
    }
    g.pelletCount = count
}

func (g *Game) Draw(screen *ebiten.Image) {
    switch g.State {
    case StateMenu:
        if g.menuUI != nil {
            g.menuUI.Draw(screen)
        }
        return
        
    case StatePlaying, StatePaused:
        g.drawGame(screen)
        
        if g.State == StatePaused {
            g.drawPauseOverlay(screen)
        }
        
    case StateGameOver:
        g.drawGame(screen)
        g.drawGameOverOverlay(screen)
    }
}

func (g *Game) drawGame(screen *ebiten.Image) {
    // Draw maze background
    DrawMaze(screen)
    
    // Draw level tiles with pellets
    for y, row := range level {
        for x, tile := range row {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(x*TileSize), float64(y*TileSize))

            switch tile {
            case TileWall:
                screen.DrawImage(WallImage, op)
            case TilePellet:
                screen.DrawImage(FloorImage, op)
                g.drawPellet(screen, x, y, false)
            case TilePowerPellet:
                screen.DrawImage(FloorImage, op)
                g.drawPellet(screen, x, y, true)
            case TileEmpty:
                screen.DrawImage(FloorImage, op)
            }
        }
    }
    
    // Draw ghosts
    for _, ghost := range g.Ghosts {
        if ghost.Visible {
            ghost.Draw(screen)
        }
    }
    
    // Draw player on top
    if g.Player != nil {
        g.Player.Draw(screen)
    }
    
    // Draw UI
    g.drawUI(screen)
}

func (g *Game) drawPellet(screen *ebiten.Image, x, y int, isPowerPellet bool) {
    cx := float64(x*TileSize + TileSize/2)
    cy := float64(y*TileSize + TileSize/2)
    
    if isPowerPellet {
        // Draw power pellet with glow effect
        size := 8.0
        if g.powerPelletActive {
            // Pulsing effect when active
            size += 2.0 * float64(g.powerPelletTimer%30) / 30.0
        }
        
        ebitenutil.DrawRect(screen, cx-size, cy-size, size*2, size*2, outerGlowColor)
        ebitenutil.DrawRect(screen, cx-size/2, cy-size/2, size, size, glowColor)
    } else {
        // Draw regular pellet
        pelletSize := 3.0
        ebitenutil.DrawRect(screen, cx-pelletSize/2, cy-pelletSize/2, pelletSize, pelletSize, color.RGBA{255, 255, 0, 255})
    }
}

func (g *Game) drawUI(screen *ebiten.Image) {
    // Score
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.Player.Score), 10, 10)
    
    // Lives
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lives: %d", g.lives), 10, 30)
    
    // Power pellet timer
    if g.powerPelletActive {
        timeLeft := g.powerPelletTimer / 60 // Convert to seconds
        ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Power: %ds", timeLeft), 10, 50)
    }
    
    // Pellets remaining
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pellets: %d", g.pelletCount), 10, 70)
}

func (g *Game) drawPauseOverlay(screen *ebiten.Image) {
    // Semi-transparent overlay
    ebitenutil.DrawRect(screen, 0, 0, float64(len(level[0])*TileSize), float64(len(level)*TileSize), 
                       color.RGBA{0, 0, 0, 128})
    
    // Pause text
    width := len(level[0]) * TileSize
    height := len(level) * TileSize
    ebitenutil.DebugPrintAt(screen, "PAUSED", width/2-30, height/2)
    ebitenutil.DebugPrintAt(screen, "Press ESC to resume", width/2-70, height/2+20)
}

func (g *Game) drawGameOverOverlay(screen *ebiten.Image) {
    // Semi-transparent overlay
    ebitenutil.DrawRect(screen, 0, 0, float64(len(level[0])*TileSize), float64(len(level)*TileSize), 
                       color.RGBA{0, 0, 0, 128})
    
    // Game over text
    width := len(level[0]) * TileSize
    height := len(level) * TileSize
    ebitenutil.DebugPrintAt(screen, "GAME OVER", width/2-40, height/2-10)
    ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Final Score: %d", g.Player.Score), width/2-60, height/2+10)
    ebitenutil.DebugPrintAt(screen, "Press SPACE to return to menu", width/2-100, height/2+30)
}

func DrawMaze(screen *ebiten.Image) {
    for y, row := range level {
        for x, tile := range row {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(x*TileSize), float64(y*TileSize))

            switch tile {
            case TileWall:
                screen.DrawImage(WallImage, op)
            default:
                screen.DrawImage(FloorImage, op)
            }
        }
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    // Use the modern menu size when in menu state
    if g.State == StateMenu {
        return 1200, 800  // Match the modern menu size
    }
    width := len(level[0]) * TileSize
    height := len(level) * TileSize
    return width, height
}

// Improved collision detection functions
func isWallColliding(level [][]int, px, py float64, size, tileSize int) bool {
    margin := 4 // Reduced margin for better movement feel
    
    corners := [][2]int{
        {int(px) + margin, int(py) + margin},
        {int(px + float64(size)) - margin, int(py) + margin},
        {int(px) + margin, int(py + float64(size)) - margin},
        {int(px + float64(size)) - margin, int(py + float64(size)) - margin},
    }
    
    for _, corner := range corners {
        cx := corner[0] / tileSize
        cy := corner[1] / tileSize
        
        // Prevent out-of-bounds access
        if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
            return true // treat out-of-bounds as wall
        }
        
        if level[cy][cx] == TileWall {
            return true
        }
    }
    
    return false
}

// Center-point collision for more forgiving gameplay
func isWallCollidingCenter(level [][]int, px, py float64, size, tileSize int) bool {
    centerX := int(px + float64(size)/2) / tileSize
    centerY := int(py + float64(size)/2) / tileSize
    
    if centerY < 0 || centerY >= len(level) || centerX < 0 || centerX >= len(level[0]) {
        return true
    }
    
    return level[centerY][centerX] == TileWall
}

// Tunnel handling for screen wrapping (if your maze has tunnels)
func handleTunnels(px, py *float64, screenWidth, screenHeight float64) {
    // Left tunnel
    if *px < -float64(TileSize) {
        *px = screenWidth
    }
    // Right tunnel  
    if *px > screenWidth {
        *px = -float64(TileSize)
    }
    // You can add vertical tunnels if needed
}

// Helper function to check if a position is valid for movement
func isValidPosition(level [][]int, x, y, size, tileSize int) bool {
    return !isWallColliding(level, float64(x), float64(y), size, tileSize)
}

// Get tile type at pixel coordinates
func getTileAt(level [][]int, px, py float64, tileSize int) int {
    tileX := int(px) / tileSize
    tileY := int(py) / tileSize
    
    if tileY < 0 || tileY >= len(level) || tileX < 0 || tileX >= len(level[0]) {
        return TileWall // Out of bounds treated as wall
    }
    
    return level[tileY][tileX]
}

func isWallCollidingStrict(level [][]int, px, py float64, size, TileSize int) bool {
    corners := [][2]int{
        {int(px), int(py)},
        {int(px + float64(size) - 1), int(py)},
        {int(px), int(py + float64(size) - 1)},
        {int(px + float64(size) - 1), int(py + float64(size) - 1)},
    }

    for _, corner := range corners {
        cx := corner[0] / TileSize
        cy := corner[1] / TileSize

        if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
            return true
        }

        if level[cy][cx] == TileWall {
            return true
        }
    }

    return false
}
