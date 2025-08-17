package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "fmt"
    "image/color"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"math"
)

var glowColor = color.RGBA{255, 50, 50, 255} // bright red center
var outerGlowColor = color.RGBA{255, 50, 50, 100} // transparent red for glow

type GameState int

const(
    StateMenu GameState = iota
    StatePlaying
    StateGameOver
    StatePaused
    StateIntro
    RoundReady
    StateRoundReady
)

type GameStateStruct struct{
    PacmanX, PacmanY   float64
    PacmanDirection    string
    Level              [][]int
    DotsRemaining      int
    PowerPelletActive  bool
    FrightModeActive   bool
    GlobalTimer        int
    Ghosts             []*Ghost
    CurrentLevel       int
    GhostManager *GhostManager
    NewGhostManager *GhostManager
}

type Game struct{
    Player *Player
    Ghosts []*Ghost
    Pellet []Pellet
    menuUI *UIPage
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

    ghostManager *GhostManager
    gameState *GameStateStruct
    globalTimer int
    IntroSystem  *IntroSystem
    SoundManager *SoundManager
    RoundReadyTimer  int
    RoundNumber  int
    ShowRoundReady bool
    AudioSystem *AudioSystem

}

const TileSize = 32
const tileSize=32

func NewGame() *Game {
    playerStartX := float64(1 * TileSize)
    playerStartY := float64(1 * TileSize)

    gameState := &GameStateStruct{
        Level: level,
        CurrentLevel: 1,
    }
    AudioSystem:= NewAudioSystem()
    
    if AudioSystem != nil {
        fmt.Println("🎵 Initializing audio system...")
        AudioSystem.LoadAllAudio()
    } else {
        fmt.Println("❌ Failed to initialize audio system")
    }

    g := &Game{
        Player: NewPlayer(playerStartX, playerStartY, "assets/player.png"),
		menuUI: NewUIPage(),
        State: StateMenu,
        lives: 3,
        playerStartX: playerStartX,
        playerStartY: playerStartY,
        gameState: gameState,
        globalTimer: 0,
        IntroSystem: NewIntroSystem(),		
        SoundManager: NewSoundManager(),		
        RoundReadyTimer:0,
        RoundNumber: 1,
        ShowRoundReady: false,
        AudioSystem: AudioSystem,
    }

    g.ghostManager = NewGhostManager(gameState) 
    
	ghosts := []*Ghost{
		NewGhost(13*TileSize, 13*TileSize, "assets/jogo.png", "jogo", 55),
    	NewGhost(12*TileSize, 13*TileSize, "assets/sakuna.png", "sukuna", 55),    
    	NewGhost(14*TileSize, 13*TileSize, "assets/kenjaku.png", "kenjaku", 55),  
    	NewGhost(13*TileSize, 14*TileSize, "assets/mahito.png", "mahito", 55),
	}

	//g.AudioSystem.LoadAllAudio()

    // Add ghosts to manager
    for _, ghost := range ghosts {
        g.ghostManager.AddGhost(ghost)
    }
    
    g.Ghosts = ghosts
	g.resetGhosts()

    InitPellets(level, TileSize)
	g.countPellets()
	fmt.Println("🎵 Starting intro music...")
	if g.SoundManager!=nil{
		g.SoundManager.PlayBGM("Intro_theme")
	}else{
    	fmt.Println("⚠️  SoundManager is nil")
	}
	if g.AudioSystem != nil {
        fmt.Println("🎵 AudioSystem is available, playing intro music")
	g.AudioSystem.PlayIntroMusic()
	} else {
        fmt.Println("⚠️  AudioSystem is nil")
    }
    return g
}

func (g *Game) Update() error {
    g.globalTimer++
    g.updateGameState()
    if g.AudioSystem != nil {
        g.AudioSystem.Update()
    }
     g.handleSoundControls()

    switch g.State {
    case StateMenu:
        return g.updateMenu()
    case StatePlaying:
        return g.updateGame()
    case StateGameOver:
        return g.updateGameOver()
    case StatePaused:
        return g.updatePaused()
    case StateIntro:
        return g.updateIntro()
    case StateRoundReady:
        return g.updateRoundReady()
    }
    return nil
}
func (g *Game) updateIntro() error {

  if g.IntroSystem != nil {
        err := g.IntroSystem.Update()
        if err != nil {
            return err
        }
        
      // Check if intro is complete

       if g.IntroSystem.IsComplete() {
            fmt.Println("Intro complete, transitioning to menu")
            g.State = StateMenu
            g.SoundManager.PlayBGM("menu_theme")
        }
    }
    return nil
}

func (g *Game) updateRoundReady() error {
    g.RoundReadyTimer++
    
  // Auto-advance after 3 seconds or on input

   if g.RoundReadyTimer > 180 || 
       inpututil.IsKeyJustPressed(ebiten.KeySpace) || 
       inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
        
      g.State = StatePlaying

       g.ShowRoundReady = false
        g.RoundReadyTimer = 0
        
      // Start game music

       g.SoundManager.PlayBGM("game_theme")
        g.SoundManager.PlaySFX("round_start")
        
      fmt.Printf("Starting round %d\n", g.RoundNumber)
    }
    
    
  return nil
}


func (g *Game) updateGameState() {
    if g.Player != nil && g.gameState != nil {
        g.gameState.PacmanX = g.Player.X
        g.gameState.PacmanY = g.Player.Y
        g.gameState.PacmanDirection = g.Player.Direction
        g.gameState.DotsRemaining = g.pelletCount
        g.gameState.PowerPelletActive = g.powerPelletActive
        g.gameState.FrightModeActive = g.powerPelletActive
        g.gameState.GlobalTimer = g.globalTimer
        g.gameState.Ghosts = g.Ghosts
    }
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
                g.State = StateRoundReady
                g.ShowRoundReady=true
                g.RoundReadyTimer=0
                g.RoundNumber=1 
                g.resetGame()
                g.SoundManager.PlaySFX("menu_selected")
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
        // Debug: Print game state occasionally
    if g.globalTimer%120 == 0 { // Every 2 seconds
        fmt.Printf("=== GAME DEBUG ===\n")
        fmt.Printf("Global Timer: %d\n", g.globalTimer)
        fmt.Printf("Player Position: (%.1f, %.1f)\n", g.Player.X, g.Player.Y)
        fmt.Printf("Number of ghosts: %d\n", len(g.Ghosts))
        fmt.Printf("Ghost Manager exists: %v\n", g.ghostManager != nil)
        
        for i, ghost := range g.Ghosts {
            if ghost != nil {
                fmt.Printf("Ghost %d (%s): (%.1f, %.1f) Speed: %.2f Visible: %v\n", 
                    i, ghost.GhostType, ghost.X, ghost.Y, ghost.Speed, ghost.Visible)
            } else {
                fmt.Printf("Ghost %d: nil\n", i)
            }
        }
        fmt.Printf("=================\n")
    }



    
    // Update player
    if g.Player != nil {
        g.Player.Update(level, TileSize)
    }
    
    // Update power pellet timer
    if g.powerPelletActive {
        g.powerPelletTimer--
        if g.powerPelletTimer <= 0 {
            g.powerPelletActive = false
            g.gameState.FrightModeActive = false
             fmt.Println("Power pellet mode ended") // Debug
            // Reset all ghosts to normal mode
            for _, ghost := range g.Ghosts {
                if ghost.Mode == FrightenedMode {
                ghost.ResetMode()}
            }
        }
    }
    // Update ghosts - try both methods to see which works
    fmt.Printf("Updating ghosts... Manager: %v, Direct: %d ghosts\n", 
        g.ghostManager != nil, len(g.Ghosts))
    // Update ghosts using ghost manager
    if g.ghostManager != nil {
        g.ghostManager.UpdateAll()
    }

    // Method 2: Update ghosts directly (for debugging)
    for i, ghost := range g.Ghosts {
        if ghost != nil {
            fmt.Printf("Updating ghost %d (%s)\n", i, ghost.GhostType)
            ghost.Update(g.gameState)
        }
    }
    
    // Check for collisions using ghost manager
    if g.ghostManager != nil {
        result := g.ghostManager.CheckCollisions(g.Player.X, g.Player.Y)
        switch result {
        case "ghost_eaten":
            g.Player.Score += 200
            g.SoundManager.PlaySFX("ghost_eaten")
            // Ghost manager already handles the ghost state change
        case "player_caught":
            g.lives--
            g.SoundManager.PlaySFX("player_death")
            g.resetPlayerPosition()
            
            
            if g.lives <= 0 {
                g.State = StateGameOver
                g.SoundManager.PlaySFX("game_over")
                g.SoundManager.StopBGM()
                return nil
            }
            
            g.resetGhosts()
        }
    }
    
    // Check pellet collection
    g.checkPelletCollection()
    
    // Check win condition
    if g.pelletCount <= 0 {
        g.RoundNumber++
        g.State=StateRoundReady
        g.ShowRoundReady=true
        g.RoundReadyTimer=0
        g.resetGame()
        g.SoundManager.PlaySFX("round_complete")

        fmt.Printf("Round %d completed! Advancing to round %d \n",g.RoundNumber-1,g.RoundNumber)
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
        
    case TilePowerPellet:
        level[playerTileY][playerTileX] = TileEmpty
        g.Player.Score += 50
        g.pelletCount--
        
        // Activate power pellet mode
        g.powerPelletActive = true
        g.powerPelletTimer = 600 // 10 seconds at 60 FPS
        g.gameState.FrightModeActive=true

        g.SoundManager.PlaySFX("power_pellet")
        // Set all visible ghosts to frightened mode
        for _, ghost := range g.Ghosts {
            if ghost.Visible {
                ghost.SetFrightened(600)
            }
        }
        // Also trigger through ghost manager
        if g.ghostManager != nil {
            g.ghostManager.TriggerFrightMode()
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



func (g *Game) resetGame() {
    if g.Player != nil {
        g.Player.Score = 0
        g.resetPlayerPosition()
    } else {
        log.Println("⚠️ Warning: g.Player is nil during resetGame")
    }
    g.lives = 3
    g.resetGhosts()
    
    // Reset level pellets
    InitPellets(level, TileSize)
    g.countPellets()
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

    case StateIntro:
        if g.IntroSystem != nil {
            g.IntroSystem.Draw(screen)
        }
        return

     case StateRoundReady:      // Add round ready drawing
        g.drawRoundReady(screen)
        return
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
// func isValidPosition(level [][]int, x, y, size, tileSize int) bool {
//     return !isWallColliding(level, float64(x), float64(y), size, tileSize)
// }

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
func (g *Game) resetGhosts() {
    // Use VERIFIED empty tile positions (check your level array first!)
    // These positions should be in the CENTER of empty tiles
    ghostStartPositions := [][2]float64{
        // Place ghosts in known empty areas - ADJUST THESE TO YOUR ACTUAL LEVEL
        {13*TileSize, 13*TileSize}, // jogo - center of ghost house
        {12*TileSize, 13*TileSize}, // sukuna - left of center
        {14*TileSize, 13*TileSize}, // kenjaku - right of center
        {13*TileSize, 14*TileSize}, // mahito - below center
    }
    
    fmt.Println("=== RESETTING GHOSTS ===")
    
    // Verify all positions are actually empty before placing ghosts
    for i, pos := range ghostStartPositions {
        if i >= len(g.Ghosts) {
            break
        }
        
        // Test if this tile is actually empty
        tileX := int(pos[0] / TileSize)
        tileY := int(pos[1] / TileSize)
        
        fmt.Printf("Checking position %d: pixel(%.1f,%.1f) = tile(%d,%d)\n", 
                   i, pos[0], pos[1], tileX, tileY)
        
        if tileY >= 0 && tileY < len(level) && tileX >= 0 && tileX < len(level[0]) {
            tileValue := level[tileY][tileX]
            fmt.Printf("  Tile value: %d\n", tileValue)
            
            if tileValue == TileWall {
                fmt.Printf("  ERROR: Position %d is a wall! Finding alternative...\n", i)
                // Find the nearest empty tile
                found := false
                for radius := 1; radius <= 5 && !found; radius++ {
                    for dy := -radius; dy <= radius && !found; dy++ {
                        for dx := -radius; dx <= radius && !found; dx++ {
                            checkX := tileX + dx
                            checkY := tileY + dy
                            
                            if checkY >= 0 && checkY < len(level) && 
                               checkX >= 0 && checkX < len(level[0]) && 
                               level[checkY][checkX] != TileWall {
                                ghostStartPositions[i][0] = float64(checkX * TileSize)
                                ghostStartPositions[i][1] = float64(checkY * TileSize)
                                fmt.Printf("  Found alternative: tile(%d,%d) = pixel(%.1f,%.1f)\n",
                                          checkX, checkY, ghostStartPositions[i][0], ghostStartPositions[i][1])
                                found = true
                            }
                        }
                    }
                }
                if !found {
                    fmt.Printf("  WARNING: Could not find safe position for ghost %d!\n", i)
                }
            }
        }
        
        // Set ghost position
        ghost := g.Ghosts[i]
        ghost.X = ghostStartPositions[i][0]
        ghost.Y = ghostStartPositions[i][1]
        
        // Reset ghost state properly
        ghost.Speed = ghost.BaseSpeed
        ghost.FrightTimer = 0
        ghost.SetVisible(true)
        
        // Set proper initial modes based on ghost type
        switch ghost.GhostType {
        case "jogo":
            ghost.Mode = ChaseMode
            ghost.ChaseTimer = 1200
            ghost.ReleaseTimer = 0
            ghost.Direction = "left" // Start moving left from house
        default:
            ghost.Mode = InHouseMode
            ghost.ReleaseTimer = 300 * (i + 1) // Stagger release times
            ghost.Direction = "up" // Default direction
        }
        
        fmt.Printf("Reset ghost %s: pos(%.1f,%.1f), mode=%d, release_timer=%d\n",
                   ghost.GhostType, ghost.X, ghost.Y, ghost.Mode, ghost.ReleaseTimer)
    }
    
    g.powerPelletActive = false
    g.powerPelletTimer = 0
    g.gameState.FrightModeActive = false
    fmt.Println("=== RESET COMPLETE ===")
}

func (g *Game) drawRoundReady(screen *ebiten.Image) {
    // Draw the game field first (dimmed)
    g.drawGame(screen)
    
    // Dark overlay
    ebitenutil.DrawRect(screen, 0, 0, float64(len(level[0])*TileSize), float64(len(level)*TileSize), 
                       color.RGBA{0, 0, 0, 180})
    
    // Pulsing effect
    // pulse := 0.7 + 0.3*math.Sin(float64(g.RoundReadyTimer)*0.15)
    
    // Round number
    width := len(level[0]) * TileSize
    height := len(level) * TileSize
    
    roundText := fmt.Sprintf("ROUND %d", g.RoundNumber)
    if bigfont != nil {
        // Use the existing font system
        ebitenutil.DebugPrintAt(screen, roundText, width/2-60, height/2-40)
    } else {
        ebitenutil.DebugPrintAt(screen, roundText, width/2-60, height/2-40)
    }
    
    // Ready text
    // readyAlpha := uint8(pulse * 255)
    readyText := "READY?"
    ebitenutil.DebugPrintAt(screen, readyText, width/2-30, height/2)
    
    // Instructions
    if g.RoundReadyTimer > 60 {
        instrText := "PRESS SPACE TO BEGIN"
        ebitenutil.DebugPrintAt(screen, instrText, width/2-80, height/2+40)
    }
    
    // Cursed energy effects around the text
    g.drawCursedEnergyEffects(screen, width/2, height/2)
}

// Add cursed energy visual effects
func (g *Game) drawCursedEnergyEffects(screen *ebiten.Image, centerX, centerY int) {
    time := float64(g.RoundReadyTimer)
    
    // Draw swirling energy particles
    for i := 0; i < 12; i++ {
        angle := float64(i)*math.Pi/6 + time*0.05
        radius := 80 + 20*math.Sin(time*0.1+float64(i))
        
        x := float64(centerX) + radius*math.Cos(angle)
        y := float64(centerY) + radius*math.Sin(angle)
        
        // Energy particle color (purple/red cursed energy)
        intensity := (math.Sin(time*0.1+float64(i)) + 1) / 2
        red := uint8(100 + intensity*155)
        blue := uint8(50 + intensity*100)
        alpha := uint8(100 + intensity*100)
        
        ebitenutil.DrawRect(screen, x-2, y-2, 4, 4, color.RGBA{red, 50, blue, alpha})
    }
}

func (g *Game) handleSoundControls() {
    if inpututil.IsKeyJustPressed(ebiten.KeyM) {
        g.SoundManager.ToggleBGM()
    }
    
    if inpututil.IsKeyJustPressed(ebiten.KeyN) {
        g.SoundManager.ToggleSFX()
    }
    
    if inpututil.IsKeyJustPressed(ebiten.KeyEqual) { // Plus key
        currentVol := g.SoundManager.Volume
        g.SoundManager.SetVolume(currentVol + 0.1)
    }
    
    if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
        currentVol := g.SoundManager.Volume
        g.SoundManager.SetVolume(currentVol - 0.1)
    }
}
func (g *Game) updateGameEnhanced() error {
    // Handle pause
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
        g.State = StatePaused
        g.SoundManager.PlaySFX("pause")
        return nil
    }

    // Update player
    if g.Player != nil {
        oldScore := g.Player.Score
        g.Player.Update(level, TileSize)
        
        // Check if score changed to play sound
        if g.Player.Score > oldScore {
            scoreDiff := g.Player.Score - oldScore
            if scoreDiff == 10 {
                g.SoundManager.PlaySFX("pellet_eat")
            } else if scoreDiff == 50 {
                g.SoundManager.PlaySFX("power_pellet")
            } else if scoreDiff >= 200 {
                g.SoundManager.PlaySFX("ghost_eaten")
            }
        }
    }
    
    // Update power pellet timer with sound
    if g.powerPelletActive {
        g.powerPelletTimer--
        
        // Warning sound when power pellet is about to end
        if g.powerPelletTimer == 120 { // 2 seconds left
            g.SoundManager.PlaySFX("power_pellet_warning")
        }
        
        if g.powerPelletTimer <= 0 {
            g.powerPelletActive = false
            g.gameState.FrightModeActive = false
            g.SoundManager.PlaySFX("power_pellet_end")
            
            for _, ghost := range g.Ghosts {
                if ghost.Mode == FrightenedMode {
                    ghost.ResetMode()
                }
            }
        }
    }
    
    // Update ghosts
    if g.ghostManager != nil {
        g.ghostManager.UpdateAll()
    }
    
    for i, ghost := range g.Ghosts {
        _=i
        if ghost != nil {
            ghost.Update(g.gameState)
        }
    }
    
    // Check for collisions
    if g.ghostManager != nil {
        result := g.ghostManager.CheckCollisions(g.Player.X, g.Player.Y)
        switch result {
        case "ghost_eaten":
            g.Player.Score += 200
            g.SoundManager.PlaySFX("ghost_eaten")
        case "player_caught":
            g.lives--
            g.SoundManager.PlaySFX("player_death")
            g.resetPlayerPosition()
            
            if g.lives <= 0 {
                g.State = StateGameOver
                g.SoundManager.PlaySFX("game_over")
                g.SoundManager.StopBGM()
                return nil
            }
            
            g.resetGhosts()
        }
    }
    
    // Check pellet collection
    g.checkPelletCollection()
    
    // Check win condition
    if g.pelletCount <= 0 {
        g.RoundNumber++
        g.State = StateRoundReady
        g.ShowRoundReady = true
        g.RoundReadyTimer = 0
        g.resetGame()
        g.SoundManager.PlaySFX("round_complete")
        
        // Increase difficulty slightly each round
        if g.RoundNumber > 1 {
            for _, ghost := range g.Ghosts {
                ghost.BaseSpeed += 0.1 // Slightly faster each round
            }
        }
    }
    
    return nil
}

// Enhanced UI drawing with sound visualization
func (g *Game) drawUIEnhanced(screen *ebiten.Image) {
    // Score with glow effect for high scores
    scoreColor := color.RGBA{255, 255, 255, 255}
    if g.Player.Score > 1000 {
        pulse := math.Sin(float64(g.globalTimer)*0.2) * 0.3 + 0.7
        scoreColor = color.RGBA{uint8(255 * pulse), uint8(255 * pulse), 100, 255}
    }
    
    scoreText := fmt.Sprintf("Score: %d", g.Player.Score)
    ebitenutil.DebugPrintAt(screen, scoreText, 10, 10)
    
    // Lives with heart symbols (or curse symbols for JJK theme)
    livesText := fmt.Sprintf("Lives: %d", g.lives)
    for i := 0; i < g.lives; i++ {
        // Draw curse symbol or use text
        ebitenutil.DrawRect(screen, float64(10+i*20), 30, 15, 15, color.RGBA{200, 50, 50, 255})
    }
    ebitenutil.DebugPrintAt(screen, livesText, 10, 50)
    
    // Power pellet timer with dramatic countdown
    if g.powerPelletActive {
        timeLeft := g.powerPelletTimer / 60
        timerColor := color.RGBA{255, 255, 100, 255}
_ = scoreColor
_ = timerColor
        
        // Change color as time runs out
        if timeLeft <= 2 {
            pulse := math.Sin(float64(g.globalTimer)*0.5) * 0.5 + 0.5
            timerColor = color.RGBA{uint8(255 * pulse), 50, 50, 255}
        }
        
        powerText := fmt.Sprintf("CURSED POWER: %ds", timeLeft)
        ebitenutil.DebugPrintAt(screen, powerText, 10, 70)
    }
    
    // Round number
    roundText := fmt.Sprintf("Round: %d", g.RoundNumber)
    ebitenutil.DebugPrintAt(screen, roundText, 10, 90)
    
    // Pellets remaining
    pelletsText := fmt.Sprintf("Pellets: %d", g.pelletCount)
    ebitenutil.DebugPrintAt(screen, pelletsText, 10, 110)
    
    // Sound indicator
    if g.SoundManager.BGMEnabled {
        ebitenutil.DebugPrintAt(screen, "♪", len(level[0])*TileSize-30, 10)
    }
}

