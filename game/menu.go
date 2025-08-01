// package main

// import (
// 	//"fmt"
// 	"image/color"
// 	"math"
// 	//"strings"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/hajimehoshi/ebiten/v2/inpututil"
// 	"github.com/hajimehoshi/ebiten/v2/text"
// 	"github.com/hajimehoshi/ebiten/v2/vector"
// 	//"golang.org/x/image/font"
// 	"golang.org/x/image/font/basicfont"
// )

// const (
// 	screenWidth  = 800
// 	screenHeight = 600
// )

// type MenuState int

// const (
// 	MenuStart MenuState = iota
// 	MenuPause
// 	MenuResume
// 	MenuQuit
// )

// type UIPage struct {
// 	selectedOption   int
// 	animationTime    float64
// 	particleTime     float64
// 	glowIntensity    float64
// 	menuOptions      []string
// 	pacmanX         float64
// 	pacmanMouthAngle float64
// 	cursedEnergy     []CursedEnergyParticle
// }

// type CursedEnergyParticle struct {
// 	x, y     float64
// 	vx, vy   float64
// 	life     float64
// 	maxLife  float64
// 	size     float64
// 	color    color.RGBA
// }

// func NewUIPage() *UIPage {
// 	ui := &UIPage{
// 		selectedOption: 0,
// 		menuOptions:    []string{"START", "PAUSE", "RESUME", "QUIT"},
// 		pacmanX:        -100,
// 		cursedEnergy:   make([]CursedEnergyParticle, 50),
// 	}
	
// 	// Initialize cursed energy particles
// 	for i := range ui.cursedEnergy {
// 		ui.cursedEnergy[i] = CursedEnergyParticle{
// 			x:       math.Mod(float64(i*16), screenWidth),
// 			y:       math.Mod(float64(i*12), screenHeight),
// 			vx:      (math.Sin(float64(i)) * 0.5),
// 			vy:      (math.Cos(float64(i)) * 0.3),
// 			life:    1.0,
// 			maxLife: 1.0,
// 			size:    2 + math.Sin(float64(i))*1,
// 			color:   color.RGBA{138, 43, 226, 100}, // Purple cursed energy
// 		}
// 	}
	
// 	return ui
// }

// func (ui *UIPage) Update() error {
// 	ui.animationTime += 0.05
// 	ui.particleTime += 0.02
// 	ui.glowIntensity = 0.5 + 0.5*math.Sin(ui.animationTime*2)
	
// 	// Update Pacman position
// 	ui.pacmanX += 2
// 	if ui.pacmanX > screenWidth+50 {
// 		ui.pacmanX = -100
// 	}
	
// 	// Update Pacman mouth animation
// 	ui.pacmanMouthAngle = math.Sin(ui.animationTime*8) * 0.7
	
// 	// Handle input
// 	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
// 		ui.selectedOption = (ui.selectedOption - 1 + len(ui.menuOptions)) % len(ui.menuOptions)
// 	}
// 	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
// 		ui.selectedOption = (ui.selectedOption + 1) % len(ui.menuOptions)
// 	}
	
// 	// Update cursed energy particles
// 	for i := range ui.cursedEnergy {
// 		p := &ui.cursedEnergy[i]
// 		p.x += p.vx
// 		p.y += p.vy
		
// 		// Wrap around screen
// 		if p.x < 0 {
// 			p.x = screenWidth
// 		}
// 		if p.x > screenWidth {
// 			p.x = 0
// 		}
// 		if p.y < 0 {
// 			p.y = screenHeight
// 		}
// 		if p.y > screenHeight {
// 			p.y = 0
// 		}
		
// 		// Pulsing effect
// 		p.life = 0.5 + 0.5*math.Sin(ui.particleTime*3+float64(i)*0.1)
// 	}
	
// 	return nil
// }

// func (ui *UIPage) Draw(screen *ebiten.Image) {
// 	// Draw dark background with gradient effect
// 	ui.drawBackground(screen)
	
// 	// Draw cursed energy particles
// 	ui.drawCursedEnergyParticles(screen)
	
// 	// Draw animated Pacman
// 	ui.drawAnimatedPacman(screen)
	
// 	// Draw title with glow effect
// 	ui.drawTitle(screen)
	
// 	// Draw menu options
// 	ui.drawMenu(screen)
	
// 	// Draw decorative elements
// 	ui.drawDecorations(screen)
	
// 	// Draw instructions
// 	ui.drawInstructions(screen)
// }

// func (ui *UIPage) drawBackground(screen *ebiten.Image) {
// 	// Create a dark gradient background
// 	for y := 0; y < screenHeight; y++ {
// 		intensity := float64(y) / float64(screenHeight)
// 		r := uint8(10 + intensity*20)
// 		g := uint8(5 + intensity*15)
// 		b := uint8(25 + intensity*40)
		
// 		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
// 			color.RGBA{r, g, b, 255}, false)
// 	}
// }

// func (ui *UIPage) drawCursedEnergyParticles(screen *ebiten.Image) {
// 	for _, p := range ui.cursedEnergy {
// 		alpha := uint8(p.life * 150)
// 		particleColor := color.RGBA{138, 43, 226, alpha}
		
// 		// Draw particle with glow effect
// 		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
// 			float32(p.size), particleColor, false)
		
// 		// Add inner glow
// 		if p.life > 0.7 {
// 			innerColor := color.RGBA{200, 100, 255, uint8(alpha/2)}
// 			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
// 				float32(p.size*0.5), innerColor, false)
// 		}
// 	}
// }

// func (ui *UIPage) drawAnimatedPacman(screen *ebiten.Image) {
// 	// Draw Pacman with Jujutsu Kaisen colors (golden/yellow)
// 	pacmanY := float32(screenHeight/2 - 100)
// 	pacmanSize := float32(40)
	
// 	// Main body
// 	pacmanColor := color.RGBA{255, 215, 0, 255} // Gold color
// 	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
// 	// Mouth (create pac-man effect)
// 	mouthAngle := ui.pacmanMouthAngle
// 	bgColor := color.RGBA{15, 10, 35, 255}
	
// 	// Draw mouth triangle
// 	if mouthAngle > 0 {
// 		x1 := float32(ui.pacmanX)
// 		//y1 := pacmanY
// 		x2 := float32(ui.pacmanX + float64(pacmanSize)*math.Cos(mouthAngle))
// 		y2 := float32(float64(pacmanY) - float64(pacmanSize)*math.Sin(mouthAngle))
// 		//x3 := float32(ui.pacmanX + float64(pacmanSize)*math.Cos(-mouthAngle))
// 		y3 := float32(float64(pacmanY) - float64(pacmanSize)*math.Sin(-mouthAngle))
		
// 		// Fill the mouth area with background color
// 		vector.DrawFilledRect(screen, x1, y2, x2-x1, y3-y2, bgColor, false)
// 	}
	
// 	// Add cursed energy aura around Pacman
// 	auraIntensity := ui.glowIntensity
// 	auraColor := color.RGBA{138, 43, 226, uint8(50 * auraIntensity)}
// 	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, 
// 		pacmanSize+float32(10*auraIntensity), auraColor, false)
// }

// func (ui *UIPage) drawTitle(screen *ebiten.Image) {
// 	title := "呪術廻戦 PAC-MAN"
// 	subtitle := "JUJUTSU KAISEN EDITION"
	
// 	// Main title with glow effect
// 	titleY := 80
// 	ui.drawGlowText(screen, title, screenWidth/2-len(title)*12, titleY, 
// 		color.RGBA{255, 215, 0, 255}, 3.0)
	
// 	// Subtitle
// 	subtitleY := titleY + 40
// 	ui.drawGlowText(screen, subtitle, screenWidth/2-len(subtitle)*6, subtitleY, 
// 		color.RGBA{200, 200, 255, 255}, 1.5)
	
// 	// Draw cursed technique symbols
// 	ui.drawCursedSymbols(screen)
// }

// func (ui *UIPage) drawMenu(screen *ebiten.Image) {
// 	menuStartY := 300
// 	menuSpacing := 60
	
// 	for i, option := range ui.menuOptions {
// 		y := menuStartY + i*menuSpacing
// 		x := screenWidth/2 - len(option)*8
		
// 		// Highlight selected option
// 		if i == ui.selectedOption {
// 			// Draw selection background with cursed energy effect
// 			glowColor := color.RGBA{138, 43, 226, uint8(100 * ui.glowIntensity)}
// 			vector.DrawFilledRect(screen, float32(x-20), float32(y-10), 
// 				float32(len(option)*16+40), 30, glowColor, false)
			
// 			// Draw option text with glow
// 			ui.drawGlowText(screen, option, x, y, 
// 				color.RGBA{255, 255, 255, 255}, 2.0)
			
// 			// Draw selection indicator
// 			indicator := "►"
// 			ui.drawGlowText(screen, indicator, x-30, y, 
// 				color.RGBA{255, 215, 0, 255}, 2.0)
// 		} else {
// 			// Draw normal option
// 			ui.drawGlowText(screen, option, x, y, 
// 				color.RGBA{150, 150, 200, 255}, 1.0)
// 		}
// 	}
// }

// func (ui *UIPage) drawGlowText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
// 	// Draw glow effect
// 	if glowIntensity > 1.0 {
// 		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.3)}
// 		for dx := -2; dx <= 2; dx++ {
// 			for dy := -2; dy <= 2; dy++ {
// 				if dx != 0 || dy != 0 {
// 					text.Draw(screen, txt, basicfont.Face7x13, x+dx, y+dy+13, glowColor)
// 				}
// 			}
// 		}
// 	}
	
// 	// Draw main text
// 	text.Draw(screen, txt, basicfont.Face7x13, x, y+13, clr)
// }

// func (ui *UIPage) drawCursedSymbols(screen *ebiten.Image) {
// 	// Draw some decorative cursed technique symbols
// 	symbols := []string{"呪", "術", "式"}
	
// 	for i, symbol := range symbols {
// 		x := 100 + i*600/len(symbols)
// 		y := 150
		
// 		// Add floating animation
// 		floatY := y + int(10*math.Sin(ui.animationTime+float64(i)))
		
// 		symbolColor := color.RGBA{138, 43, 226, uint8(150 * ui.glowIntensity)}
// 		ui.drawGlowText(screen, symbol, x, floatY, symbolColor, 2.0)
// 	}
// }

// func (ui *UIPage) drawDecorations(screen *ebiten.Image) {
// 	// Draw decorative border lines with cursed energy effect
// 	borderColor := color.RGBA{138, 43, 226, uint8(200 * ui.glowIntensity)}
	
// 	// Top border
// 	vector.DrawFilledRect(screen, 50, 50, screenWidth-100, 3, borderColor, false)
// 	// Bottom border  
// 	vector.DrawFilledRect(screen, 50, screenHeight-53, screenWidth-100, 3, borderColor, false)
	
// 	// Side decorative elements
// 	for i := 0; i < 5; i++ {
// 		y := 100 + i*100
// 		intensity := 0.5 + 0.5*math.Sin(ui.animationTime*2+float64(i)*0.5)
// 		decorColor := color.RGBA{255, 215, 0, uint8(100 * intensity)}
		
// 		// Left side
// 		vector.DrawFilledCircle(screen, 30, float32(y), 5, decorColor, false)
// 		// Right side
// 		vector.DrawFilledCircle(screen, screenWidth-30, float32(y), 5, decorColor, false)
// 	}
// }

// func (ui *UIPage) drawInstructions(screen *ebiten.Image) {
// 	instructions := []string{
// 		"↑↓ Navigate Menu",
// 		"ENTER Select Option",
// 		"Embrace the Cursed Energy!",
// 	}
	
// 	startY := screenHeight - 100
// 	for i, instruction := range instructions {
// 		x := screenWidth/2 - len(instruction)*4
// 		y := startY + i*20
		
// 		instrColor := color.RGBA{150, 150, 200, 200}
// 		if i == 2 { // Special color for the last line
// 			instrColor = color.RGBA{138, 43, 226, 200}
// 		}
		
// 		text.Draw(screen, instruction, basicfont.Face7x13, x, y, instrColor)
// 	}
// }

// func (ui *UIPage) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

// // GetSelectedOption returns the currently selected menu option
// func (ui *UIPage) GetSelectedOption() int {
// 	return ui.selectedOption
// }

// // IsEnterPressed checks if enter key was just pressed
// func (ui *UIPage) IsEnterPressed() bool {
// 	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
// }

// // Integration example for your game.go:
// /*
// // Add this to your main game struct:
// type Game struct {
//     menuUI *UIPage
//     showMenu bool
//     // ... your existing fields
// }

// // In your game's Update() method:
// func (g *Game) Update() error {
//     if g.showMenu {
//         err := g.menuUI.Update()
//         if err != nil {
//             return err
//         }
        
//         // Check if user selected an option
//         if g.menuUI.IsEnterPressed() {
//             switch g.menuUI.GetSelectedOption() {
//             case 0: // START
//                 g.showMenu = false
//                 // Start your game logic
//             case 1: // PAUSE
//                 // Pause game logic
//             case 2: // RESUME
//                 g.showMenu = false
//                 // Resume game logic
//             case 3: // QUIT
//                 return fmt.Errorf("quit game")
//             }
//         }
//         return nil
//     }
    
//     // Your existing game update logic here
//     return nil
// }

// // In your game's Draw() method:
// func (g *Game) Draw(screen *ebiten.Image) {
//     if g.showMenu {
//         g.menuUI.Draw(screen)
//         return
//     }
    
//     // Your existing game draw logic here
// }

// // Initialize in your main function:
// func main() {
//     game := &Game{
//         menuUI: NewUIPage(),
//         showMenu: true,
//     }
//     // ... rest of your ebiten.RunGame setup
// }
// */
package main

import (
	//"fmt"
	"image/color"
	"math"
	"math/rand"
	//"strings"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	//"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "io/ioutil"
)

var bigFont font.Face 

const (
	screenWidth  = 1200  // Increased screen size
	screenHeight = 800   // Increased screen size
	
	// Background image/gif placement area
	// BACKGROUND_IMAGE_X = 0
	// BACKGROUND_IMAGE_Y = 0  
	// BACKGROUND_IMAGE_WIDTH = 1200
	// BACKGROUND_IMAGE_HEIGHT = 800
	// You can place a full-screen background image/gif here
	
	// Character image/gif placement areas
	// CHARACTER_LEFT_X = 50     // Left side character area
	// CHARACTER_LEFT_Y = 200
	// CHARACTER_LEFT_WIDTH = 300
	// CHARACTER_LEFT_HEIGHT = 400
	
	// CHARACTER_RIGHT_X = 850   // Right side character area  
	// CHARACTER_RIGHT_Y = 200
	// CHARACTER_RIGHT_WIDTH = 300
	// CHARACTER_RIGHT_HEIGHT = 400
)

type MenuState int

const (
	MenuStart MenuState = iota
	MenuPause
	MenuResume
	MenuQuit
)

type UIPage struct {
	selectedOption     int
	animationTime      float64
	particleTime       float64
	glowIntensity      float64
	menuOptions        []string
	pacmanX           float64
	pacmanMouthAngle  float64
	cursedEnergy      []CursedEnergyParticle
	domainParticles   []DomainParticle
	selectedGlow      float64
	pulseIntensity    float64
}

type CursedEnergyParticle struct {
	x, y      float64
	vx, vy    float64
	life      float64
	size      float64
	color     color.RGBA
	intensity float64
}

type DomainParticle struct {
	x, y        float64
	angle       float64
	radius      float64
	speed       float64
	life        float64
	glowRadius  float64
}

func init() {
    fontBytes, _ := ioutil.ReadFile("assets/PressStart2P-Regular.ttf")
    ttf, _ := opentype.Parse(fontBytes)
    bigFont, _ = opentype.NewFace(ttf, &opentype.FaceOptions{
        Size: 32,
        DPI:  72,
    })
}

func NewUIPage() *UIPage {
	ui := &UIPage{
		selectedOption: 0,
		menuOptions:    []string{"START GAME", "PAUSE", "RESUME", "QUIT"},
		pacmanX:        -150,
		cursedEnergy:   make([]CursedEnergyParticle, 40), // Reduced for cleaner look
		domainParticles: make([]DomainParticle, 20),      // Reduced particles
	}
	
	// Initialize cursed energy particles (fewer for cleaner look)
	for i := range ui.cursedEnergy {
		ui.cursedEnergy[i] = CursedEnergyParticle{
			x:         rand.Float64() * screenWidth,
			y:         rand.Float64() * screenHeight,
			vx:        (rand.Float64() - 0.5) * 1.5,
			vy:        (rand.Float64() - 0.5) * 1.0,
			life:      rand.Float64(),
			size:      2 + rand.Float64()*4,
			color:     color.RGBA{138, 43, 226, uint8(30 + rand.Intn(70))},
			intensity: rand.Float64(),
		}
	}
	
	// Initialize domain expansion particles
	for i := range ui.domainParticles {
		ui.domainParticles[i] = DomainParticle{
			x:          screenWidth / 2,
			y:          screenHeight / 2,
			angle:      rand.Float64() * 2 * math.Pi,
			radius:     100 + rand.Float64()*200,
			speed:      0.3 + rand.Float64()*0.8,
			life:       rand.Float64(),
			glowRadius: 3 + rand.Float64()*6,
		}
	}
	
	return ui
}

func (ui *UIPage) Update() error {
	ui.animationTime += 0.03
	ui.particleTime += 0.02
	ui.glowIntensity = 0.6 + 0.4*math.Sin(ui.animationTime*2.5)
	ui.pulseIntensity = 0.4 + 0.6*math.Sin(ui.animationTime*3)
	ui.selectedGlow = 0.6 + 0.4*math.Sin(ui.animationTime*4)
	
	// Update Pacman movement
	ui.pacmanX += 3
	if ui.pacmanX > screenWidth+100 {
		ui.pacmanX = -150
	}
	
	// Pacman mouth animation
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
		if p.x < 0 { p.x = screenWidth }
		if p.x > screenWidth { p.x = 0 }
		if p.y < 0 { p.y = screenHeight }
		if p.y > screenHeight { p.y = 0 }
		
		// Pulsing effect
		p.intensity = 0.4 + 0.6*math.Sin(ui.particleTime*3+float64(i)*0.2)
	}
	
	// Update domain particles
	for i := range ui.domainParticles {
		p := &ui.domainParticles[i]
		p.angle += p.speed * 0.015
		p.x = screenWidth/2 + math.Cos(p.angle)*p.radius
		p.y = screenHeight/2 + math.Sin(p.angle)*p.radius
		p.life = 0.3 + 0.7*math.Sin(ui.animationTime*1.5+float64(i)*0.4)
	}
	
	return nil
}

func (ui *UIPage) Draw(screen *ebiten.Image) {
	// STEP 1: Draw background image/gif here if you have one
	// Example: screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})
	
	// Draw gradient background
	ui.drawCleanBackground(screen)
	
	// Draw subtle cursed energy particles
	ui.drawCursedEnergyParticles(screen)
	
	// Draw domain expansion effect (subtle)
	ui.drawDomainParticles(screen)
	
	// STEP 2: Draw character images here
	// Left side character area (50, 200, 300x400)
	// Right side character area (850, 200, 300x400)
	
	// Draw enhanced Pacman
	ui.drawEnhancedPacman(screen)
	
	// Draw clean title
	ui.drawCleanTitle(screen)
	
	// Draw menu with better spacing
	ui.drawCleanMenu(screen)
	
	// Draw instructions
	ui.drawInstructions(screen)
}

func (ui *UIPage) drawCleanBackground(screen *ebiten.Image) {
	// Clean gradient background
	for y := 0; y < screenHeight; y++ {
		progress := float64(y) / float64(screenHeight)
		
		// Darker, cleaner gradient
		r := uint8(8 + progress*20)
		g := uint8(4 + progress*15)
		b := uint8(25 + progress*45)
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{r, g, b, 255}, false)
	}
	
	// Subtle vignette effect
	centerX, centerY := screenWidth/2, screenHeight/2
	maxDist := math.Sqrt(float64(centerX*centerX + centerY*centerY))
	
	for y := 0; y < screenHeight; y += 8 {
		for x := 0; x < screenWidth; x += 8 {
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist > maxDist*0.7 {
				alpha := uint8((dist - maxDist*0.7) / (maxDist*0.3) * 60)
				vector.DrawFilledRect(screen, float32(x), float32(y), 8, 8,
					color.RGBA{0, 0, 0, alpha}, false)
			}
		}
	}
}

func (ui *UIPage) drawCursedEnergyParticles(screen *ebiten.Image) {
	for _, p := range ui.cursedEnergy {
		alpha := uint8(p.intensity * float64(p.color.A))
		
		// Main particle
		mainColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y),
			float32(p.size), mainColor, false)
		
		// Subtle glow
		if p.intensity > 0.7 {
			glowColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha/3}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y),
				float32(p.size*2), glowColor, false)
		}
	}
}

func (ui *UIPage) drawDomainParticles(screen *ebiten.Image) {
	for _, p := range ui.domainParticles {
		if p.life < 0.3 { continue } // Only draw visible particles
		
		alpha := uint8(p.life * 100)
		
		// Main particle
		particleColor := color.RGBA{200, 100, 255, alpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y),
			float32(p.glowRadius), particleColor, false)
		
		// Connection line to center (very subtle)
		if p.life > 0.8 {
			lineColor := color.RGBA{138, 43, 226, uint8(alpha/6)}
			ui.drawLine(screen, float32(p.x), float32(p.y),
				screenWidth/2, screenHeight/2, 1, lineColor)
		}
	}
}

func (ui *UIPage) drawLine(screen *ebiten.Image, x1, y1, x2, y2, thickness float32, clr color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	
	if length > 0 {
		steps := int(length / 4) // Fewer steps for performance
		for i := 0; i <= steps; i++ {
			t := float32(i) / float32(steps)
			x := x1 + dx*t
			y := y1 + dy*t
			vector.DrawFilledCircle(screen, x, y, thickness/2, clr, false)
		}
	}
}

func (ui *UIPage) drawEnhancedPacman(screen *ebiten.Image) {
	pacmanY := float32(150) // Moved higher up
	pacmanSize := float32(60) // Bigger size
	
	// Cursed energy aura (subtle)
	auraIntensity := ui.pulseIntensity * 0.6
	auraColor := color.RGBA{138, 43, 226, uint8(40 * auraIntensity)}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY,
		pacmanSize+float32(20*auraIntensity), auraColor, false)
	
	// Main Pacman body
	pacmanColor := color.RGBA{255, 215, 0, 255}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
	// Mouth animation
	mouthAngle := ui.pacmanMouthAngle
	if mouthAngle > 0 {
		bgColor := color.RGBA{15, 10, 35, 255}
		centerX := float32(ui.pacmanX)
		centerY := pacmanY
		
		mouthSize := pacmanSize * 0.8
		x2 := centerX + float32(math.Cos(mouthAngle)*float64(mouthSize))
		y2 := centerY - float32(math.Sin(mouthAngle)*float64(mouthSize))
		x3 := centerX + float32(math.Cos(-mouthAngle)*float64(mouthSize))
		y3 := centerY - float32(math.Sin(-mouthAngle)*float64(mouthSize))
		
		ui.drawTriangle(screen, centerX, centerY, x2, y2, x3, y3, bgColor)
	}
}

func (ui *UIPage) drawTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, clr color.RGBA) {
	// Simple triangle fill
	minY := int(math.Min(float64(y1), math.Min(float64(y2), float64(y3))))
	maxY := int(math.Max(float64(y1), math.Max(float64(y2), float64(y3))))
	
	for y := minY; y <= maxY; y++ {
		vector.DrawFilledRect(screen, x1-15, float32(y), 30, 1, clr, false)
	}
}

func (ui *UIPage) drawCleanTitle(screen *ebiten.Image) {
	// Main title - larger and cleaner
	title := "JUJUTSU KAISEN"
	subtitle := "PAC-MAN EDITION"
	
	titleY := 100
	
	// Subtle domain expansion circle behind title
	domainRadius := 200 + math.Sin(ui.animationTime)*30
	domainColor := color.RGBA{138, 43, 226, uint8(25*ui.glowIntensity)}
	vector.DrawFilledCircle(screen, screenWidth/2, float32(titleY+30),
		float32(domainRadius), domainColor, false)
	
	// Title with enhanced glow
	titleX := screenWidth/2 - len(title)*20 // Bigger font spacing
	ui.drawLargeGlowText(screen, title, titleX, titleY,
		color.RGBA{255, 215, 0, 255}, 3.0)
	
	// Subtitle
	subtitleY := titleY + 80 // More spacing
	subtitleX := screenWidth/2 - len(subtitle)*16
	ui.drawLargeGlowText(screen, subtitle, subtitleX, subtitleY,
		color.RGBA{200, 200, 255, 255}, 2.0)
}

func (ui *UIPage) drawCleanMenu(screen *ebiten.Image) {
	menuStartY := 400 // Moved down for better spacing
	menuSpacing := 100 // Increased spacing
	
	for i, option := range ui.menuOptions {
		y := menuStartY + i*menuSpacing
		x := screenWidth/2 - len(option)*20 // Bigger font spacing
		
		if i == ui.selectedOption {
			// Clean selection background
			selectionWidth := float32(len(option)*40 + 100)
			selectionHeight := float32(70)
			
			// Glow effect
			glowColor := color.RGBA{138, 43, 226, uint8(120 * ui.selectedGlow)}
			vector.DrawFilledRect(screen,
				float32(x-50), float32(y-25),
				selectionWidth, selectionHeight, glowColor, false)
			
			// Enhanced selected text
			ui.drawLargeGlowText(screen, option, x, y,
				color.RGBA{255, 255, 255, 255}, 3.0)
			
			// Selection indicator
			indicator := ">"
			ui.drawLargeGlowText(screen, indicator, x-70, y,
				color.RGBA{255, 215, 0, 255}, 2.5)
		} else {
			// Normal menu option
			ui.drawLargeGlowText(screen, option, x, y,
				color.RGBA{150, 150, 200, 255}, 1.0)
		}
	}
}

func (ui *UIPage) drawLargeGlowText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	// Draw multiple text instances for larger appearance and glow
	offsets := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0},           {1, 0},
		{-1, 1},  {0, 1},  {1, 1},
	}
	
	// Glow effect
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.4)}
		for _, offset := range offsets {
			for scale := 1; scale <= int(glowIntensity); scale++ {
				dx := offset.dx * scale * 2
				dy := offset.dy * scale * 2
				// Draw multiple times for "larger" font effect
				text.Draw(screen, txt, basicfont.Face7x13, x+dx, y+dy+15, glowColor)
				text.Draw(screen, txt, basicfont.Face7x13, x+dx+1, y+dy+15, glowColor)
				text.Draw(screen, txt, basicfont.Face7x13, x+dx, y+dy+16, glowColor)
				text.Draw(screen, txt, basicfont.Face7x13, x+dx+1, y+dy+16, glowColor)
			}
		}
	}
	
	// Main text (draw multiple times for thickness/size)
	text.Draw(screen, txt, basicfont.Face7x13, x, y+15, clr)
	text.Draw(screen, txt, basicfont.Face7x13, x+1, y+15, clr)
	text.Draw(screen, txt, basicfont.Face7x13, x, y+16, clr)
	text.Draw(screen, txt, basicfont.Face7x13, x+1, y+16, clr)
}

func (ui *UIPage) drawInstructions(screen *ebiten.Image) {
	instructions := []string{
		"Arrow Keys: Navigate Menu",
		"Enter: Select Option",
		"Experience the Cursed Energy!",
	}
	
	startY := screenHeight - 120
	for i, instruction := range instructions {
		x := screenWidth/2 - len(instruction)*8
		y := startY + i*30
		
		instrColor := color.RGBA{120, 120, 160, 180}
		if i == 2 { // Special color for the last line
			instrColor = color.RGBA{138, 43, 226, 180}
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

/*
IMAGE/GIF PLACEMENT GUIDE:
========================

1. FULL BACKGROUND IMAGE/GIF:
   - Position: (0, 0)
   - Size: 1200x800 (full screen)
   - Place in drawCleanBackground() function
   - Example: screen.DrawImage(backgroundImg, &ebiten.DrawImageOptions{})

2. LEFT CHARACTER IMAGE/GIF:
   - Position: (50, 200)  
   - Size: 300x400
   - Good for: Gojo, Yuji, or other main characters
   - Place after drawCleanBackground() call

3. RIGHT CHARACTER IMAGE/GIF:
   - Position: (850, 200)
   - Size: 300x400  
   - Good for: Sukuna, Megumi, or villain characters
   - Place after drawCleanBackground() call

4. SMALLER DECORATIVE IMAGES:
   - Various positions around the UI
   - Keep them small to avoid clutter
   - Recommended size: 64x64 to 128x128

The menu is positioned in the center (400-750 area) so it won't overlap
with side character images. The title area (top 250px) is clear for 
background images to show through.
*/

// Integration example for your game.go:
/*
type Game struct {
    menuUI *UIPage
    showMenu bool
    backgroundImg *ebiten.Image  // Add your background image/gif
    leftCharImg *ebiten.Image    // Add your left character
    rightCharImg *ebiten.Image   // Add your right character
    // ... your existing fields
}

func (g *Game) Update() error {
    if g.showMenu {
        err := g.menuUI.Update()
        if err != nil {
            return err
        }
        
        if g.menuUI.IsEnterPressed() {
            switch g.menuUI.GetSelectedOption() {
            case 0: // START GAME
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

func (g *Game) Draw(screen *ebiten.Image) {
    if g.showMenu {
        g.menuUI.Draw(screen)
        
        // Draw your images after the UI base
        if g.backgroundImg != nil {
            op := &ebiten.DrawImageOptions{}
            screen.DrawImage(g.backgroundImg, op)
        }
        
        if g.leftCharImg != nil {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(50, 200)  // Left character position
            screen.DrawImage(g.leftCharImg, op)
        }
        
        if g.rightCharImg != nil {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(850, 200)  // Right character position
            screen.DrawImage(g.rightCharImg, op)
        }
        
        return
    }
    
    // Your existing game draw logic here
}
*/
