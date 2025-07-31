package main

import (
	//"fmt"
	"image/color"
	"math"
	//"strings"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	//"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type MenuState int

const (
	MenuStart MenuState = iota
	MenuPause
	MenuResume
	MenuQuit
)

type UIPage struct {
	selectedOption   int
	animationTime    float64
	particleTime     float64
	glowIntensity    float64
	menuOptions      []string
	pacmanX         float64
	pacmanMouthAngle float64
	cursedEnergy     []CursedEnergyParticle
}

type CursedEnergyParticle struct {
	x, y     float64
	vx, vy   float64
	life     float64
	maxLife  float64
	size     float64
	color    color.RGBA
}

func NewUIPage() *UIPage {
	ui := &UIPage{
		selectedOption: 0,
		menuOptions:    []string{"START", "PAUSE", "RESUME", "QUIT"},
		pacmanX:        -100,
		cursedEnergy:   make([]CursedEnergyParticle, 50),
	}
	
	// Initialize cursed energy particles
	for i := range ui.cursedEnergy {
		ui.cursedEnergy[i] = CursedEnergyParticle{
			x:       math.Mod(float64(i*16), screenWidth),
			y:       math.Mod(float64(i*12), screenHeight),
			vx:      (math.Sin(float64(i)) * 0.5),
			vy:      (math.Cos(float64(i)) * 0.3),
			life:    1.0,
			maxLife: 1.0,
			size:    2 + math.Sin(float64(i))*1,
			color:   color.RGBA{138, 43, 226, 100}, // Purple cursed energy
		}
	}
	
	return ui
}

func (ui *UIPage) Update() error {
	ui.animationTime += 0.05
	ui.particleTime += 0.02
	ui.glowIntensity = 0.5 + 0.5*math.Sin(ui.animationTime*2)
	
	// Update Pacman position
	ui.pacmanX += 2
	if ui.pacmanX > screenWidth+50 {
		ui.pacmanX = -100
	}
	
	// Update Pacman mouth animation
	ui.pacmanMouthAngle = math.Sin(ui.animationTime*8) * 0.7
	
	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		ui.selectedOption = (ui.selectedOption - 1 + len(ui.menuOptions)) % len(ui.menuOptions)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		ui.selectedOption = (ui.selectedOption + 1) % len(ui.menuOptions)
	}
	
	// Update cursed energy particles
	for i := range ui.cursedEnergy {
		p := &ui.cursedEnergy[i]
		p.x += p.vx
		p.y += p.vy
		
		// Wrap around screen
		if p.x < 0 {
			p.x = screenWidth
		}
		if p.x > screenWidth {
			p.x = 0
		}
		if p.y < 0 {
			p.y = screenHeight
		}
		if p.y > screenHeight {
			p.y = 0
		}
		
		// Pulsing effect
		p.life = 0.5 + 0.5*math.Sin(ui.particleTime*3+float64(i)*0.1)
	}
	
	return nil
}

func (ui *UIPage) Draw(screen *ebiten.Image) {
	// Draw dark background with gradient effect
	ui.drawBackground(screen)
	
	// Draw cursed energy particles
	ui.drawCursedEnergyParticles(screen)
	
	// Draw animated Pacman
	ui.drawAnimatedPacman(screen)
	
	// Draw title with glow effect
	ui.drawTitle(screen)
	
	// Draw menu options
	ui.drawMenu(screen)
	
	// Draw decorative elements
	ui.drawDecorations(screen)
	
	// Draw instructions
	ui.drawInstructions(screen)
}

func (ui *UIPage) drawBackground(screen *ebiten.Image) {
	// Create a dark gradient background
	for y := 0; y < screenHeight; y++ {
		intensity := float64(y) / float64(screenHeight)
		r := uint8(10 + intensity*20)
		g := uint8(5 + intensity*15)
		b := uint8(25 + intensity*40)
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{r, g, b, 255}, false)
	}
}

func (ui *UIPage) drawCursedEnergyParticles(screen *ebiten.Image) {
	for _, p := range ui.cursedEnergy {
		alpha := uint8(p.life * 150)
		particleColor := color.RGBA{138, 43, 226, alpha}
		
		// Draw particle with glow effect
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(p.size), particleColor, false)
		
		// Add inner glow
		if p.life > 0.7 {
			innerColor := color.RGBA{200, 100, 255, uint8(alpha/2)}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size*0.5), innerColor, false)
		}
	}
}

func (ui *UIPage) drawAnimatedPacman(screen *ebiten.Image) {
	// Draw Pacman with Jujutsu Kaisen colors (golden/yellow)
	pacmanY := float32(screenHeight/2 - 100)
	pacmanSize := float32(40)
	
	// Main body
	pacmanColor := color.RGBA{255, 215, 0, 255} // Gold color
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
	// Mouth (create pac-man effect)
	mouthAngle := ui.pacmanMouthAngle
	bgColor := color.RGBA{15, 10, 35, 255}
	
	// Draw mouth triangle
	if mouthAngle > 0 {
		x1 := float32(ui.pacmanX)
		//y1 := pacmanY
		x2 := float32(ui.pacmanX + float64(pacmanSize)*math.Cos(mouthAngle))
		y2 := float32(float64(pacmanY) - float64(pacmanSize)*math.Sin(mouthAngle))
		//x3 := float32(ui.pacmanX + float64(pacmanSize)*math.Cos(-mouthAngle))
		y3 := float32(float64(pacmanY) - float64(pacmanSize)*math.Sin(-mouthAngle))
		
		// Fill the mouth area with background color
		vector.DrawFilledRect(screen, x1, y2, x2-x1, y3-y2, bgColor, false)
	}
	
	// Add cursed energy aura around Pacman
	auraIntensity := ui.glowIntensity
	auraColor := color.RGBA{138, 43, 226, uint8(50 * auraIntensity)}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, 
		pacmanSize+float32(10*auraIntensity), auraColor, false)
}

func (ui *UIPage) drawTitle(screen *ebiten.Image) {
	title := "呪術廻戦 PAC-MAN"
	subtitle := "JUJUTSU KAISEN EDITION"
	
	// Main title with glow effect
	titleY := 80
	ui.drawGlowText(screen, title, screenWidth/2-len(title)*12, titleY, 
		color.RGBA{255, 215, 0, 255}, 3.0)
	
	// Subtitle
	subtitleY := titleY + 40
	ui.drawGlowText(screen, subtitle, screenWidth/2-len(subtitle)*6, subtitleY, 
		color.RGBA{200, 200, 255, 255}, 1.5)
	
	// Draw cursed technique symbols
	ui.drawCursedSymbols(screen)
}

func (ui *UIPage) drawMenu(screen *ebiten.Image) {
	menuStartY := 300
	menuSpacing := 60
	
	for i, option := range ui.menuOptions {
		y := menuStartY + i*menuSpacing
		x := screenWidth/2 - len(option)*8
		
		// Highlight selected option
		if i == ui.selectedOption {
			// Draw selection background with cursed energy effect
			glowColor := color.RGBA{138, 43, 226, uint8(100 * ui.glowIntensity)}
			vector.DrawFilledRect(screen, float32(x-20), float32(y-10), 
				float32(len(option)*16+40), 30, glowColor, false)
			
			// Draw option text with glow
			ui.drawGlowText(screen, option, x, y, 
				color.RGBA{255, 255, 255, 255}, 2.0)
			
			// Draw selection indicator
			indicator := "►"
			ui.drawGlowText(screen, indicator, x-30, y, 
				color.RGBA{255, 215, 0, 255}, 2.0)
		} else {
			// Draw normal option
			ui.drawGlowText(screen, option, x, y, 
				color.RGBA{150, 150, 200, 255}, 1.0)
		}
	}
}

func (ui *UIPage) drawGlowText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	// Draw glow effect
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.3)}
		for dx := -2; dx <= 2; dx++ {
			for dy := -2; dy <= 2; dy++ {
				if dx != 0 || dy != 0 {
					text.Draw(screen, txt, basicfont.Face7x13, x+dx, y+dy+13, glowColor)
				}
			}
		}
	}
	
	// Draw main text
	text.Draw(screen, txt, basicfont.Face7x13, x, y+13, clr)
}

func (ui *UIPage) drawCursedSymbols(screen *ebiten.Image) {
	// Draw some decorative cursed technique symbols
	symbols := []string{"呪", "術", "式"}
	
	for i, symbol := range symbols {
		x := 100 + i*600/len(symbols)
		y := 150
		
		// Add floating animation
		floatY := y + int(10*math.Sin(ui.animationTime+float64(i)))
		
		symbolColor := color.RGBA{138, 43, 226, uint8(150 * ui.glowIntensity)}
		ui.drawGlowText(screen, symbol, x, floatY, symbolColor, 2.0)
	}
}

func (ui *UIPage) drawDecorations(screen *ebiten.Image) {
	// Draw decorative border lines with cursed energy effect
	borderColor := color.RGBA{138, 43, 226, uint8(200 * ui.glowIntensity)}
	
	// Top border
	vector.DrawFilledRect(screen, 50, 50, screenWidth-100, 3, borderColor, false)
	// Bottom border  
	vector.DrawFilledRect(screen, 50, screenHeight-53, screenWidth-100, 3, borderColor, false)
	
	// Side decorative elements
	for i := 0; i < 5; i++ {
		y := 100 + i*100
		intensity := 0.5 + 0.5*math.Sin(ui.animationTime*2+float64(i)*0.5)
		decorColor := color.RGBA{255, 215, 0, uint8(100 * intensity)}
		
		// Left side
		vector.DrawFilledCircle(screen, 30, float32(y), 5, decorColor, false)
		// Right side
		vector.DrawFilledCircle(screen, screenWidth-30, float32(y), 5, decorColor, false)
	}
}

func (ui *UIPage) drawInstructions(screen *ebiten.Image) {
	instructions := []string{
		"↑↓ Navigate Menu",
		"ENTER Select Option",
		"Embrace the Cursed Energy!",
	}
	
	startY := screenHeight - 100
	for i, instruction := range instructions {
		x := screenWidth/2 - len(instruction)*4
		y := startY + i*20
		
		instrColor := color.RGBA{150, 150, 200, 200}
		if i == 2 { // Special color for the last line
			instrColor = color.RGBA{138, 43, 226, 200}
		}
		
		text.Draw(screen, instruction, basicfont.Face7x13, x, y, instrColor)
	}
}

func (ui *UIPage) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// GetSelectedOption returns the currently selected menu option
func (ui *UIPage) GetSelectedOption() int {
	return ui.selectedOption
}

// IsEnterPressed checks if enter key was just pressed
func (ui *UIPage) IsEnterPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Integration example for your game.go:
/*
// Add this to your main game struct:
type Game struct {
    menuUI *UIPage
    showMenu bool
    // ... your existing fields
}

// In your game's Update() method:
func (g *Game) Update() error {
    if g.showMenu {
        err := g.menuUI.Update()
        if err != nil {
            return err
        }
        
        // Check if user selected an option
        if g.menuUI.IsEnterPressed() {
            switch g.menuUI.GetSelectedOption() {
            case 0: // START
                g.showMenu = false
                // Start your game logic
            case 1: // PAUSE
                // Pause game logic
            case 2: // RESUME
                g.showMenu = false
                // Resume game logic
            case 3: // QUIT
                return fmt.Errorf("quit game")
            }
        }
        return nil
    }
    
    // Your existing game update logic here
    return nil
}

// In your game's Draw() method:
func (g *Game) Draw(screen *ebiten.Image) {
    if g.showMenu {
        g.menuUI.Draw(screen)
        return
    }
    
    // Your existing game draw logic here
}

// Initialize in your main function:
func main() {
    game := &Game{
        menuUI: NewUIPage(),
        showMenu: true,
    }
    // ... rest of your ebiten.RunGame setup
}
*/
