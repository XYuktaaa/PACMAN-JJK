package main

import (
	"image/color"
	"math"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/basicfont"
	"embed"
	"log"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

// Embedded font files (if available)
//go:embed assets/fonts/*.ttf
var fontFiles embed.FS

type MenuState int

const (
	MenuStart MenuState = iota
	MenuPause
	MenuResume
	MenuQuit
)

type UIPage struct {
	// Core menu state
	selectedOption       int
	menuOptions          []string
	
	// Animation timers
	animationTime        float64
	particleTime         float64
	glowIntensity        float64
	titlePulse          float64
	menuPulse           float64
	energyWave          float64
	
	// Visual effects
	selectionTransition float64
	screenShake        float64
	
	// Assets
	logoImage          *ebiten.Image
	characterGif       *ebiten.Image  
	bgImage            *ebiten.Image
	gifFrames          []*ebiten.Image
	frameIndex         int
	frameTicker        int
	frameDelay         int
	
	// Font system - MUCH LARGER SIZES
	titleFont          font.Face
	menuFont           font.Face
	subtitleFont       font.Face
	uiFont            font.Face
	
	// Audio integration
	audioSystem       *AudioSystem
	
	// Pacman animation
	pacmanX             float64
	pacmanMouthAngle    float64
	
	// Simplified particle systems (removed nebula clouds)
	cursedEnergy        []CursedEnergyParticle
	backgroundParticles []BackgroundParticle
	hexagons           []HexagonElement
	floatingElements   []FloatingElement
	
	// JJK specific effects (simplified)
	domainExpansion    []DomainParticle
	malevolentAura     []AuraParticle
	
	// UI States
	showDomainExpansion bool
	domainTimer         int
}

// Simplified particle types
type CursedEnergyParticle struct {
	x, y         float64
	vx, vy       float64
	size         float64
	color        color.RGBA
	pulsePhase   float64
	energy       float64
	cursedType   string
}

type BackgroundParticle struct {
	x, y       float64
	vx, vy     float64
	size       float64
	rotation   float64
	rotSpeed   float64
	color      color.RGBA
	shape      int
	depth      float64
}

type HexagonElement struct {
	x, y        float64
	size        float64
	rotation    float64
	rotSpeed    float64
	alpha       float64
	pulsePhase  float64
	cursedLevel int
}

type FloatingElement struct {
	x, y          float64
	vx, vy        float64
	rotation      float64
	scale         float64
	alpha         float64
	pulsePhase    float64
	kanjiChar     string
}

type DomainParticle struct {
	x, y           float64
	radius         float64
	expansion      float64
	domainType     string
	alpha          float64
	rotationSpeed  float64
	currentAngle   float64
}

type AuraParticle struct {
	x, y       float64
	intensity  float64
	auraType   string
	pulsePhase float64
	size       float64
}

// Constructor
func NewUIPage() *UIPage {
	ui := &UIPage{
		selectedOption:      0,
		menuOptions:        []string{"呪術廻戦 START", "設定 SETTINGS", "ギャラリー GALLERY", "終了 EXIT"},
		pacmanX:           -150,
		frameDelay:        4,
		selectionTransition: 0,
		showDomainExpansion: false,
		domainTimer:         0,
	}
	
	ui.initializeParticles()
	ui.initializeJJKEffects()
	ui.initializeFonts()
	
	return ui
}

func (ui *UIPage) initializeParticles() {
	// Reduced cursed energy particles for cleaner look
	ui.cursedEnergy = make([]CursedEnergyParticle, 80) // Reduced from 180
	for i := range ui.cursedEnergy {
		cursedTypes := []string{"positive", "negative", "neutral"}
		cursedType := cursedTypes[i%len(cursedTypes)]
		
		var particleColor color.RGBA
		switch cursedType {
		case "positive":
			particleColor = color.RGBA{100, 200, 255, 150} // Reduced alpha for cleaner look
		case "negative":
			particleColor = color.RGBA{200, 50, 100, 150}
		case "neutral":
			particleColor = color.RGBA{150, 150, 255, 150}
		}
		
		ui.cursedEnergy[i] = CursedEnergyParticle{
			x:          math.Mod(float64(i*25), screenWidth),
			y:          math.Mod(float64(i*20), screenHeight),
			vx:         (math.Sin(float64(i)) * 1.5),
			vy:         (math.Cos(float64(i)) * 1.2),
			size:       2 + math.Sin(float64(i))*3,
			pulsePhase: float64(i) * 0.1,
			color:      particleColor,
			energy:     0.4 + math.Sin(float64(i))*0.6,
			cursedType: cursedType,
		}
	}
	
	// Minimal background particles
	ui.backgroundParticles = make([]BackgroundParticle, 40) // Reduced from 120
	for i := range ui.backgroundParticles {
		ui.backgroundParticles[i] = BackgroundParticle{
			x:        math.Mod(float64(i*40), screenWidth),
			y:        math.Mod(float64(i*30), screenHeight),
			vx:       (math.Sin(float64(i)*0.1) * 0.4),
			vy:       (math.Cos(float64(i)*0.1) * 0.3),
			size:     1 + math.Sin(float64(i))*3,
			rotation: float64(i) * 0.1,
			rotSpeed: 0.008 + math.Sin(float64(i))*0.005,
			shape:    i % 3, // Only 3 simple shapes
			depth:    0.3 + math.Sin(float64(i))*0.7,
			color:    color.RGBA{80, 120, 200, 80}, // More transparent
		}
	}
	
	// Fewer hexagons
	ui.hexagons = make([]HexagonElement, 12) // Reduced from 24
	for i := range ui.hexagons {
		angle := float64(i) * 2 * math.Pi / float64(len(ui.hexagons))
		radius := 300.0 + math.Sin(float64(i))*60
		ui.hexagons[i] = HexagonElement{
			x:           screenWidth/2 + math.Cos(angle)*radius,
			y:           screenHeight/2 + math.Sin(angle)*radius,
			size:        15 + math.Sin(float64(i))*10,
			rotation:    angle,
			rotSpeed:    0.003 + math.Sin(float64(i))*0.003,
			alpha:       0.3 + math.Sin(float64(i))*0.2, // More subtle
			pulsePhase:  float64(i) * 0.15,
			cursedLevel: (i % 3) + 1,
		}
	}
	
	// Minimal floating elements
	kanjiChars := []string{"呪", "術", "廻", "戦", "領", "域"}
	ui.floatingElements = make([]FloatingElement, 8) // Reduced from 25
	for i := range ui.floatingElements {
		ui.floatingElements[i] = FloatingElement{
			x:          math.Mod(float64(i*150), screenWidth),
			y:          math.Mod(float64(i*100), screenHeight),
			vx:         (math.Sin(float64(i)) * 0.5),
			vy:         (math.Cos(float64(i)) * 0.3),
			rotation:   float64(i) * 0.5,
			scale:      0.7 + math.Sin(float64(i))*0.3,
			alpha:      0.4 + math.Sin(float64(i))*0.3,
			pulsePhase: float64(i) * 0.25,
			kanjiChar:  kanjiChars[i%len(kanjiChars)],
		}
	}
}

func (ui *UIPage) initializeJJKEffects() {
	// Minimal JJK effects
	ui.domainExpansion = make([]DomainParticle, 3) // Reduced from 6
	for i := range ui.domainExpansion {
		ui.domainExpansion[i] = DomainParticle{
			x:             screenWidth/2 + math.Sin(float64(i)*1.0)*150,
			y:             screenHeight/2 + math.Cos(float64(i)*1.0)*150,
			radius:        60 + float64(i)*30,
			expansion:     0,
			domainType:    []string{"infinite_void", "malevolent_shrine", "coffin_of_iron"}[i%3],
			alpha:         0.6,
			rotationSpeed: 0.02 + float64(i)*0.005,
			currentAngle:  float64(i) * math.Pi / 2,
		}
	}
	
	// Minimal malevolent aura
	ui.malevolentAura = make([]AuraParticle, 20) // Reduced from 50
	for i := range ui.malevolentAura {
		ui.malevolentAura[i] = AuraParticle{
			x:          math.Mod(float64(i*60), screenWidth),
			y:          math.Mod(float64(i*40), screenHeight),
			intensity:  0.3 + math.Sin(float64(i))*0.4,
			auraType:   []string{"malevolent", "limitless"}[i%2],
			pulsePhase: float64(i) * 0.2,
			size:       4 + math.Sin(float64(i))*4,
		}
	}
}

func (ui *UIPage) initializeFonts() {
	// DRAMATICALLY LARGER FONT SIZES for maximum readability
	ui.titleFont = ui.loadFont("JujutsuTitle.ttf", 96)     // Increased from 72 to 96
	ui.menuFont = ui.loadFont("JujutsuMenu.ttf", 64)      // Increased from 48 to 64  
	ui.subtitleFont = ui.loadFont("JujutsuSubtitle.ttf", 48) // Increased from 36 to 48
	ui.uiFont = ui.loadFont("JujutsuUI.ttf", 36)          // Increased from 28 to 36
}

func (ui *UIPage) drawTitle(screen *ebiten.Image) {
	title := "呪術廻戦 × PAC-MAN"
	subtitle := "JUJUTSU KAISEN EDITION"
	
	titleY := 140
	titleX := screenWidth/2 - 480 // Adjusted for much larger font
	
	// Much larger title background
	panelWidth := float32(1000)  // Increased from 800
	panelHeight := float32(180)  // Increased from 140
	panelX := float32(titleX - 80)
	panelY := float32(titleY - 50)
	
	// Simple background panel
	panelColor := color.RGBA{15, 25, 60, 200}
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight, panelColor, false)
	
	// Clean borders
	borderColor := color.RGBA{150, 100, 255, uint8(200 * ui.glowIntensity)}
	accentColor := color.RGBA{255, 215, 0, uint8(180 * ui.glowIntensity)}
	
	vector.StrokeRect(screen, panelX-2, panelY-2, panelWidth+4, panelHeight+4, 3, borderColor, false)
	vector.StrokeRect(screen, panelX+3, panelY+3, panelWidth-6, panelHeight-6, 2, accentColor, false)
	
	// Much larger title text
	titleColor := color.RGBA{255, 215, 0, 255}
	ui.drawLargeText(screen, title, titleX, titleY, titleColor, 5.0, ui.titleFont)
	
	// Much larger subtitle
	subtitleY := titleY + 100  // More spacing for larger text
	subtitleX := screenWidth/2 - 320
	subtitleHue := 0.8 + 0.2*math.Sin(ui.titlePulse*2)
	subtitleColor := color.RGBA{
		uint8(180 * subtitleHue), 
		uint8(220 * subtitleHue), 
		255, 255,
	}
	ui.drawLargeText(screen, subtitle, subtitleX, subtitleY, subtitleColor, 4.0, ui.subtitleFont)
}

func (ui *UIPage) drawMenu(screen *ebiten.Image) {
	menuStartY := 420  // Adjusted for larger title
	menuSpacing := 140 // Much larger spacing for bigger text
	menuWidth := 700   // Much wider for bigger text
	menuX := screenWidth/2 - menuWidth/2
	
	panelHeight := float32(len(ui.menuOptions)*menuSpacing + 120)
	panelX := float32(menuX-100)
	panelY := float32(menuStartY-60)
	panelW := float32(menuWidth+200)
	
	// Draw GIF background (cleaner)
	if len(ui.gifFrames) > 0 && ui.gifFrames[ui.frameIndex] != nil {
		ui.drawGIFBackground(screen, panelX, panelY, panelW, panelHeight)
	}
	
	// Clean overlay for text readability
	overlayColor := color.RGBA{10, 15, 35, 120}
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelHeight, overlayColor, false)
	
	// Clean panel borders
	borderColor := color.RGBA{150, 100, 255, uint8(220 * ui.glowIntensity)}
	accentColor := color.RGBA{255, 215, 0, uint8(180 * ui.glowIntensity)}
	
	vector.StrokeRect(screen, panelX-3, panelY-3, panelW+6, panelHeight+6, 4, borderColor, false)
	vector.StrokeRect(screen, panelX+3, panelY+3, panelW-6, panelHeight-6, 2, accentColor, false)
	
	// Much larger menu options
	for i, option := range ui.menuOptions {
		y := menuStartY + i*menuSpacing
		x := menuX
		
		if i == ui.selectedOption {
			ui.drawSelectedOption(screen, option, x, y, menuWidth)
		} else {
			ui.drawUnselectedOption(screen, option, x, y, i)
		}
	}
}

func (ui *UIPage) drawSelectedOption(screen *ebiten.Image, option string, x, y, menuWidth int) {
	selectionWidth := float32(menuWidth + 100)
	selectionHeight := float32(120) // Much taller for bigger text
	selectionX := float32(x - 60)
	selectionY := float32(y - 40)
	
	pulseIntensity := 0.8 + 0.2*math.Sin(ui.menuPulse*3)
	
	// Clean selection background
	mainAlpha := uint8(160 * pulseIntensity)
	mainColor := color.RGBA{120, 0, 180, mainAlpha}
	vector.DrawFilledRect(screen, selectionX, selectionY, selectionWidth, selectionHeight, mainColor, false)
	
	// Selection borders
	borderGlow := color.RGBA{255, 255, 255, uint8(220 * pulseIntensity)}
	vector.StrokeRect(screen, selectionX, selectionY, selectionWidth, selectionHeight, 3, borderGlow, false)
	
	accentBorder := color.RGBA{255, 215, 0, uint8(180 * pulseIntensity)}
	vector.StrokeRect(screen, selectionX+2, selectionY+2, selectionWidth-4, selectionHeight-4, 2, accentBorder, false)
	
	// Much larger selected text
	ui.drawLargeText(screen, option, x, y, color.RGBA{255, 255, 255, 255}, 5.0, ui.menuFont)
	
	// Clean selection indicators
	ui.drawCleanSelectionIndicators(screen, selectionX, selectionY, selectionWidth, selectionHeight)
}

func (ui *UIPage) drawUnselectedOption(screen *ebiten.Image, option string, x, y, index int) {
	hoverIntensity := 0.9 + 0.1*math.Sin(ui.animationTime*1.5+float64(index)*0.5)
	textColor := color.RGBA{
		uint8(180 * hoverIntensity), 
		uint8(190 * hoverIntensity), 
		uint8(220 * hoverIntensity), 
		220,
	}
	// Much larger unselected text
	ui.drawLargeText(screen, option, x, y, textColor, 3.0, ui.menuFont)
}

func (ui *UIPage) drawCleanSelectionIndicators(screen *ebiten.Image, x, y, w, h float32) {
	// Larger, more visible indicators
	indicatorX := x - 40
	indicatorY := y + h/2
	
	pulsePhase := ui.menuPulse*4
	intensity := 0.7 + 0.3*math.Sin(pulsePhase)
	
	indicatorColor := color.RGBA{255, 100, 100, uint8(220 * intensity)}
	size := float32(12 + 6*intensity) // Much larger indicators
	
	vector.DrawFilledCircle(screen, indicatorX, indicatorY, size, indicatorColor, false)
	
	// Larger right indicator
	rightIndicatorX := x + w + 40
	rightIndicatorColor := color.RGBA{100, 150, 255, uint8(220 * intensity)}
	
	vector.DrawFilledCircle(screen, rightIndicatorX, indicatorY, size, rightIndicatorColor, false)
}

func (ui *UIPage) drawCharacterArea(screen *ebiten.Image) {
	charAreaX := float32(screenWidth - 350)
	charAreaY := float32(100)
	charAreaWidth := float32(300)
	charAreaHeight := float32(300)
	
	// Clean frame
	frameColor1 := color.RGBA{120, 80, 200, uint8(160 * ui.glowIntensity)}
	frameColor2 := color.RGBA{255, 215, 0, uint8(120 * ui.glowIntensity)}
	
	// Simple border
	vector.StrokeRect(screen, charAreaX-3, charAreaY-3, 
		charAreaWidth+6, charAreaHeight+6, 3, frameColor1, false)
	vector.StrokeRect(screen, charAreaX+2, charAreaY+2, 
		charAreaWidth-4, charAreaHeight-4, 2, frameColor2, false)
	
	// Much larger placeholder text
	centerX := int(charAreaX + charAreaWidth/2)
	centerY := int(charAreaY + charAreaHeight/2)
	
	placeholderText := "CHARACTER"
	ui.drawLargeText(screen, placeholderText, centerX-160, centerY-60, 
		color.RGBA{200, 220, 255, 200}, 4.0, ui.subtitleFont)
	
	subText := "DISPLAY AREA"
	ui.drawLargeText(screen, subText, centerX-140, centerY, 
		color.RGBA{160, 180, 210, 160}, 3.0, ui.uiFont)
}

func (ui *UIPage) drawFooter(screen *ebiten.Image) {
	instructions := []string{
		"↑↓ / W S  Navigate Menu",
		"ENTER / SPACE  Select Option", 
		"Experience Jujutsu Kaisen",
	}
	
	footerY := screenHeight - 100  // More space from bottom
	spacing := 380
	
	for i, instruction := range instructions {
		x := 100 + i*spacing
		y := footerY
		
		var instrColor color.RGBA
		switch i {
		case 0:
			instrColor = color.RGBA{150, 200, 255, 220}
		case 1:
			instrColor = color.RGBA{255, 215, 0, 220}
		case 2:
			pulse := 0.8 + 0.2*math.Sin(ui.animationTime*2)
			instrColor = color.RGBA{
				uint8(180 * pulse), 
				uint8(100 * pulse), 
				uint8(255 * pulse), 
				220,
			}
		}
		
		// Much larger footer text
		ui.drawLargeText(screen, instruction, x, y, instrColor, 3.0, ui.uiFont)
	}
	
	// Clean decorative line
	lineY := float32(footerY - 40)
	lineColor := color.RGBA{120, 80, 200, uint8(120 * ui.glowIntensity)}
	vector.DrawFilledRect(screen, 80, lineY, screenWidth-160, 2, lineColor, false)
}

func (ui *UIPage) drawEnhancedBasicText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	// For basic font, create MUCH larger text by drawing in a 4x4 grid
	gridSize := 4 // Increased from 2x2 to 4x4 for much larger appearance
	
	offsets := make([]struct{dx, dy int}, 0, gridSize*gridSize)
	for gx := 0; gx < gridSize; gx++ {
		for gy := 0; gy < gridSize; gy++ {
			offsets = append(offsets, struct{dx, dy int}{gx * 3, gy * 3}) // Larger spacing
		}
	}
	
	// Enhanced glow for basic text
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.3)}
		glowRadius := int(glowIntensity * 2)
		
		for dx := -glowRadius; dx <= glowRadius; dx++ {
			for dy := -glowRadius; dy <= glowRadius; dy++ {
				if dx != 0 || dy != 0 {
					distance := math.Sqrt(float64(dx*dx + dy*dy))
					if distance <= float64(glowRadius) {
						alpha := uint8(float64(glowColor.A) * (1.0 - distance/float64(glowRadius)))
						fadeColor := color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
						
						for _, offset := range offsets {
							text.Draw(screen, txt, basicfont.Face7x13, x+dx+offset.dx, y+dy+offset.dy+25, fadeColor)
						}
					}
				}
			}
		}
	}
	
	// Main text with much larger simulated size (4x4 grid)
	for _, offset := range offsets {
		text.Draw(screen, txt, basicfont.Face7x13, x+offset.dx, y+offset.dy+25, clr)
	}
	
	// Add extra thickness for better visibility
	for _, offset := range offsets {
		text.Draw(screen, txt, basicfont.Face7x13, x+offset.dx+1, y+offset.dy+26, clr)
		text.Draw(screen, txt, basicfont.Face7x13, x+offset.dx+2, y+offset.dy+27, clr)
	}
}

func (ui *UIPage) loadFont(filename string, size float64) font.Face {
	fontData, err := fontFiles.ReadFile("assets/fonts/" + filename)
	if err != nil {
		log.Printf("Font %s not found, using fallback", filename)
		return basicfont.Face7x13
	}
	
	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Printf("Failed to parse font %s: %v", filename, err)
		return basicfont.Face7x13
	}
	
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Printf("Failed to create font face for %s: %v", filename, err)
		return basicfont.Face7x13
	}
	
	return face
}

// Update function
func (ui *UIPage) Update() error {
	// Update animation timers
	ui.animationTime += 0.03
	ui.particleTime += 0.02
	ui.menuPulse += 0.08
	ui.titlePulse += 0.06
	ui.energyWave += 0.04
	
	ui.glowIntensity = 0.7 + 0.3*math.Sin(ui.animationTime*2.0)
	
	// Handle input
	ui.handleInput()
	ui.updateAnimations()
	ui.updateParticles()
	ui.updateJJKEffects()
	
	return nil
}

func (ui *UIPage) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		ui.selectedOption = (ui.selectedOption - 1 + len(ui.menuOptions)) % len(ui.menuOptions)
		ui.selectionTransition = 1.0
		ui.screenShake = 6.0
		if ui.audioSystem != nil {
			ui.audioSystem.PlaySFX("menu_move")
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ui.selectedOption = (ui.selectedOption + 1) % len(ui.menuOptions)
		ui.selectionTransition = 1.0
		ui.screenShake = 6.0
		if ui.audioSystem != nil {
			ui.audioSystem.PlaySFX("menu_move")
		}
	}
	
	// Domain expansion effect on START selection
	if ui.selectedOption == 0 && !ui.showDomainExpansion {
		ui.showDomainExpansion = true
		ui.domainTimer = 0
		if ui.audioSystem != nil {
			ui.audioSystem.PlaySFX("domain_expansion")
		}
	} else if ui.selectedOption != 0 && ui.showDomainExpansion {
		ui.showDomainExpansion = false
	}
}

func (ui *UIPage) updateAnimations() {
	// Decay transitions
	ui.selectionTransition *= 0.85
	ui.screenShake *= 0.9
	
	// Domain expansion timer
	if ui.showDomainExpansion {
		ui.domainTimer++
		if ui.domainTimer > 300 {
			ui.domainTimer = 0
		}
	}
	
	// Pacman animation
	ui.pacmanX += 2.0
	if ui.pacmanX > screenWidth+250 {
		ui.pacmanX = -250
	}
	ui.pacmanMouthAngle = math.Sin(ui.animationTime*8) * 1.0
	
	// GIF frame updating
	if len(ui.gifFrames) > 0 {
		ui.frameTicker++
		if ui.frameTicker >= ui.frameDelay {
			ui.frameIndex = (ui.frameIndex + 1) % len(ui.gifFrames)
			ui.frameTicker = 0
		}
	}
}

func (ui *UIPage) updateParticles() {
	// Update cursed energy particles (simplified)
	for i := range ui.cursedEnergy {
		p := &ui.cursedEnergy[i]
		
		// Gentle magnetic attraction to selected menu
		menuY := float64(400 + ui.selectedOption*120) // Updated for larger spacing
		menuX := float64(screenWidth/2)
		dx := menuX - p.x
		dy := menuY - p.y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		if distance > 0 && distance < 200 {
			force := 0.002 / (distance * 0.01)
			p.vx += (dx / distance) * force
			p.vy += (dy / distance) * force
		}
		
		p.vx *= 0.98
		p.vy *= 0.98
		
		p.x += p.vx
		p.y += p.vy
		
		// Wrapping
		if p.x < -50 { p.x = screenWidth + 50 }
		if p.x > screenWidth+50 { p.x = -50 }
		if p.y < -50 { p.y = screenHeight + 50 }
		if p.y > screenHeight+50 { p.y = -50 }
		
		p.energy = 0.4 + 0.4*math.Sin(ui.particleTime*2.0+p.pulsePhase)
	}
	
	// Update background particles (simplified)
	for i := range ui.backgroundParticles {
		p := &ui.backgroundParticles[i]
		p.x += p.vx * p.depth
		p.y += p.vy * p.depth
		p.rotation += p.rotSpeed
		
		if p.x < -30 { p.x = screenWidth + 30 }
		if p.x > screenWidth+30 { p.x = -30 }
		if p.y < -30 { p.y = screenHeight + 30 }
		if p.y > screenHeight+30 { p.y = -30 }
	}
	
	// Update hexagons (simplified)
	for i := range ui.hexagons {
		h := &ui.hexagons[i]
		h.rotation += h.rotSpeed
		h.alpha = 0.2 + 0.3*math.Sin(ui.animationTime*1.5+h.pulsePhase)
	}
	
	// Update floating elements (simplified)
	for i := range ui.floatingElements {
		fe := &ui.floatingElements[i]
		fe.x += fe.vx
		fe.y += fe.vy
		fe.scale = 0.8 + 0.2*math.Sin(ui.animationTime+fe.pulsePhase)
		fe.alpha = 0.4 + 0.3*math.Sin(ui.animationTime*1.3+fe.pulsePhase)
		
		if fe.x < -60 { fe.x = screenWidth + 60 }
		if fe.x > screenWidth+60 { fe.x = -60 }
		if fe.y < -60 { fe.y = screenHeight + 60 }
		if fe.y > screenHeight+60 { fe.y = -60 }
	}
}

func (ui *UIPage) updateJJKEffects() {
	// Update Domain Expansion (simplified)
	for i := range ui.domainExpansion {
		dp := &ui.domainExpansion[i]
		dp.currentAngle += dp.rotationSpeed
		
		if ui.showDomainExpansion {
			dp.expansion = math.Min(dp.expansion+0.025, 1.0)
		} else {
			dp.expansion = math.Max(dp.expansion-0.02, 0.0)
		}
	}
	
	// Update malevolent aura (simplified)
	for i := range ui.malevolentAura {
		ap := &ui.malevolentAura[i]
		ap.intensity = 0.3 + 0.4*math.Sin(ui.energyWave*2+ap.pulsePhase)
		ap.size = 3 + 3*math.Sin(ui.animationTime*2+ap.pulsePhase)
	}
}

// Drawing functions
func (ui *UIPage) Draw(screen *ebiten.Image) {
	// Minimal screen shake
	shakeX := math.Sin(ui.animationTime*15) * ui.screenShake
	shakeY := math.Cos(ui.animationTime*18) * ui.screenShake
	
	tempScreen := ebiten.NewImage(screenWidth, screenHeight)
	
	// Draw all layers
	ui.drawBackground(tempScreen)
	ui.drawParticleEffects(tempScreen)
	ui.drawJJKEffects(tempScreen)
	ui.drawUI(tempScreen)
	
	// Apply shake and draw to main screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(tempScreen, op)
}

func (ui *UIPage) drawBackground(screen *ebiten.Image) {
	// Clean background image
	if ui.bgImage != nil {
		imgBounds := ui.bgImage.Bounds()
		imgWidth := float64(imgBounds.Dx())
		imgHeight := float64(imgBounds.Dy())
		
		scaleX := float64(screenWidth) / imgWidth
		scaleY := float64(screenHeight) / imgHeight
		scale := math.Max(scaleX, scaleY)
		
		scaledWidth := imgWidth * scale
		scaledHeight := imgHeight * scale
		offsetX := (float64(screenWidth) - scaledWidth) / 2
		offsetY := (float64(screenHeight) - scaledHeight) / 2
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(offsetX, offsetY)
		op.ColorM.Scale(1, 1, 1, 0.7) // Slightly more visible than before
		
		screen.DrawImage(ui.bgImage, op)
	}
	
	// Clean gradient
	for y := 0; y < screenHeight; y++ {
		progress := float64(y) / float64(screenHeight)
		
		r1, g1, b1 := 10, 15, 40    // Cleaner dark blue
		r2, g2, b2 := 40, 25, 70    // Cleaner purple
		
		r := uint8(float64(r1) + (float64(r2-r1))*progress)
		g := uint8(float64(g1) + (float64(g2-g1))*progress)
		b := uint8(float64(b1) + (float64(b2-b1))*progress)
		
		// Subtle energy waves (much cleaner)
		waveR := math.Sin(float64(y)*0.008 + ui.energyWave*0.5) * 8
		waveG := math.Sin(float64(y)*0.012 + ui.energyWave*0.7) * 5
		waveB := math.Sin(float64(y)*0.006 + ui.energyWave*0.4) * 10
		
		r = uint8(math.Max(0, math.Min(255, float64(r)+waveR)))
		g = uint8(math.Max(0, math.Min(255, float64(g)+waveG)))
		b = uint8(math.Max(0, math.Min(255, float64(b)+waveB)))
		
		alpha := uint8(220)
		if ui.bgImage != nil {
			alpha = 160
		}
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{r, g, b, alpha}, false)
	}
}

func (ui *UIPage) drawParticleEffects(screen *ebiten.Image) {
	ui.drawBackgroundParticles(screen)
	ui.drawHexagonElements(screen)
	ui.drawCursedEnergy(screen)
	ui.drawFloatingElements(screen)
}

func (ui *UIPage) drawJJKEffects(screen *ebiten.Image) {
	// Draw Domain Expansion effects (simplified)
	if ui.showDomainExpansion {
		for _, dp := range ui.domainExpansion {
			if dp.expansion > 0 {
				ui.drawDomainExpansion(screen, &dp)
			}
		}
	}
}

func (ui *UIPage) drawUI(screen *ebiten.Image) {
	ui.drawTitle(screen)
	ui.drawMenu(screen)
	ui.drawPacman(screen)
	ui.drawCharacterArea(screen)
	ui.drawFooter(screen)
}

func (ui *UIPage) drawTitle(screen *ebiten.Image) {
	title := "呪術廻戦 × PAC-MAN"
	subtitle := "JUJUTSU KAISEN EDITION"
	
	titleY := 120
	titleX := screenWidth/2 - 360 // Adjusted for larger font
	
	// Clean title background
	panelWidth := float32(800)  // Larger for bigger text
	panelHeight := float32(140) // Taller for bigger text
	panelX := float32(titleX - 60)
	panelY := float32(titleY - 40)
	
	// Simple background panel
	panelColor := color.RGBA{15, 25, 60, 200}
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight, panelColor, false)
	
	// Clean borders
	borderColor := color.RGBA{150, 100, 255, uint8(200 * ui.glowIntensity)}
	accentColor := color.RGBA{255, 215, 0, uint8(180 * ui.glowIntensity)}
	
	vector.StrokeRect(screen, panelX-2, panelY-2, panelWidth+4, panelHeight+4, 3, borderColor, false)
	vector.StrokeRect(screen, panelX+3, panelY+3, panelWidth-6, panelHeight-6, 2, accentColor, false)
	
	// Large title text
	titleColor := color.RGBA{255, 215, 0, 255}
	ui.drawLargeText(screen, title, titleX, titleY, titleColor, 4.0, ui.titleFont)
	
	// Large subtitle
	subtitleY := titleY + 80  // More spacing for larger text
	subtitleX := screenWidth/2 - 250
	subtitleHue := 0.8 + 0.2*math.Sin(ui.titlePulse*2)
	subtitleColor := color.RGBA{
		uint8(180 * subtitleHue), 
		uint8(220 * subtitleHue), 
		255, 255,
	}
	ui.drawLargeText(screen, subtitle, subtitleX, subtitleY, subtitleColor, 3.0, ui.subtitleFont)
}

func (ui *UIPage) drawMenu(screen *ebiten.Image) {
	menuStartY := 380  // Adjusted for larger spacing
	menuSpacing := 120 // Much larger spacing for bigger text
	menuWidth := 600   // Wider for bigger text
	menuX := screenWidth/2 - menuWidth/2
	
	panelHeight := float32(len(ui.menuOptions)*menuSpacing + 100)
	panelX := float32(menuX-80)
	panelY := float32(menuStartY-50)
	panelW := float32(menuWidth+160)
	
	// Draw GIF background (cleaner)
	if len(ui.gifFrames) > 0 && ui.gifFrames[ui.frameIndex] != nil {
		ui.drawGIFBackground(screen, panelX, panelY, panelW, panelHeight)
	}
	
	// Clean overlay for text readability
	overlayColor := color.RGBA{10, 15, 35, 100}
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelHeight, overlayColor, false)
	
	// Clean panel borders
	borderColor := color.RGBA{150, 100, 255, uint8(220 * ui.glowIntensity)}
	accentColor := color.RGBA{255, 215, 0, uint8(180 * ui.glowIntensity)}
	
	vector.StrokeRect(screen, panelX-3, panelY-3, panelW+6, panelHeight+6, 4, borderColor, false)
	vector.StrokeRect(screen, panelX+3, panelY+3, panelW-6, panelHeight-6, 2, accentColor, false)
	
	// Large menu options
	for i, option := range ui.menuOptions {
		y := menuStartY + i*menuSpacing
		x := menuX
		
		if i == ui.selectedOption {
			ui.drawSelectedOption(screen, option, x, y, menuWidth)
		} else {
			ui.drawUnselectedOption(screen, option, x, y, i)
		}
	}
}

func (ui *UIPage) drawSelectedOption(screen *ebiten.Image, option string, x, y, menuWidth int) {
	selectionWidth := float32(menuWidth + 80)
	selectionHeight := float32(100) // Taller for bigger text
	selectionX := float32(x - 50)
	selectionY := float32(y - 30)
	
	pulseIntensity := 0.8 + 0.2*math.Sin(ui.menuPulse*3)
	
	// Clean selection background
	mainAlpha := uint8(160 * pulseIntensity)
	mainColor := color.RGBA{120, 0, 180, mainAlpha}
	vector.DrawFilledRect(screen, selectionX, selectionY, selectionWidth, selectionHeight, mainColor, false)
	
	// Selection borders
	borderGlow := color.RGBA{255, 255, 255, uint8(220 * pulseIntensity)}
	vector.StrokeRect(screen, selectionX, selectionY, selectionWidth, selectionHeight, 3, borderGlow, false)
	
	accentBorder := color.RGBA{255, 215, 0, uint8(180 * pulseIntensity)}
	vector.StrokeRect(screen, selectionX+2, selectionY+2, selectionWidth-4, selectionHeight-4, 2, accentBorder, false)
	
	// Large selected text
	ui.drawLargeText(screen, option, x, y, color.RGBA{255, 255, 255, 255}, 4.0, ui.menuFont)
	
	// Clean selection indicators
	ui.drawCleanSelectionIndicators(screen, selectionX, selectionY, selectionWidth, selectionHeight)
}

func (ui *UIPage) drawUnselectedOption(screen *ebiten.Image, option string, x, y, index int) {
	hoverIntensity := 0.9 + 0.1*math.Sin(ui.animationTime*1.5+float64(index)*0.5)
	textColor := color.RGBA{
		uint8(180 * hoverIntensity), 
		uint8(190 * hoverIntensity), 
		uint8(220 * hoverIntensity), 
		220,
	}
	ui.drawLargeText(screen, option, x, y, textColor, 2.0, ui.menuFont)
}

func (ui *UIPage) drawCleanSelectionIndicators(screen *ebiten.Image, x, y, w, h float32) {
	// Simple left indicator
	indicatorX := x - 30
	indicatorY := y + h/2
	
	pulsePhase := ui.menuPulse*4
	intensity := 0.7 + 0.3*math.Sin(pulsePhase)
	
	indicatorColor := color.RGBA{255, 100, 100, uint8(200 * intensity)}
	size := float32(8 + 4*intensity)
	
	vector.DrawFilledCircle(screen, indicatorX, indicatorY, size, indicatorColor, false)
	
	// Simple right indicator
	rightIndicatorX := x + w + 30
	rightIndicatorColor := color.RGBA{100, 150, 255, uint8(200 * intensity)}
	
	vector.DrawFilledCircle(screen, rightIndicatorX, indicatorY, size, rightIndicatorColor, false)
}

func (ui *UIPage) drawGIFBackground(screen *ebiten.Image, panelX, panelY, panelW, panelHeight float32) {
	currentFrame := ui.gifFrames[ui.frameIndex]
	
	frameBounds := currentFrame.Bounds()
	frameWidth := float64(frameBounds.Dx())
	frameHeight := float64(frameBounds.Dy())
	
	scaleX := float64(panelW) / frameWidth
	scaleY := float64(panelHeight) / frameHeight
	scale := math.Max(scaleX, scaleY)
	
	scaledWidth := frameWidth * scale
	scaledHeight := frameHeight * scale
	offsetX := float64(panelX) + (float64(panelW)-scaledWidth)/2
	offsetY := float64(panelY) + (float64(panelHeight)-scaledHeight)/2
	
	tempImg := ebiten.NewImage(int(panelW), int(panelHeight))
	tempOp := &ebiten.DrawImageOptions{}
	tempOp.GeoM.Scale(scale, scale)
	tempOp.GeoM.Translate(offsetX-float64(panelX), offsetY-float64(panelY))
	tempOp.ColorM.Scale(1, 1, 1, 0.8) // Slightly reduced opacity for cleaner look
	tempOp.CompositeMode = ebiten.CompositeModeSourceOver
	
	tempImg.DrawImage(currentFrame, tempOp)
	
	finalOp := &ebiten.DrawImageOptions{}
	finalOp.GeoM.Translate(float64(panelX), float64(panelY))
	screen.DrawImage(tempImg, finalOp)
}

func (ui *UIPage) drawPacman(screen *ebiten.Image) {
	pacmanY := float32(screenHeight/2 - 60)
	pacmanSize := float32(80) // Slightly smaller for cleaner look
	
	// Clean aura (fewer layers)
	auraLayers := 4 // Reduced from 7
	for layer := 0; layer < auraLayers; layer++ {
		layerSize := pacmanSize + float32(layer*20) + float32(25*ui.glowIntensity)
		layerAlpha := uint8(float64(80) / float64(layer + 1) * ui.glowIntensity)
		
		if layer < 2 {
			auraColor := color.RGBA{180, 100, 255, layerAlpha}
			vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, layerSize, auraColor, false)
		} else {
			auraColor := color.RGBA{120, 0, 180, layerAlpha/2}
			vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, layerSize, auraColor, false)
		}
	}
	
	// Main Pacman body
	pacmanColor := color.RGBA{255, 215, 0, 255}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
	// Clean highlights
	highlight1 := color.RGBA{255, 255, 220, 200}
	highlight2 := color.RGBA{255, 255, 255, 120}
	
	vector.DrawFilledCircle(screen, float32(ui.pacmanX-12), pacmanY-12, 
		pacmanSize*0.35, highlight1, false)
	vector.DrawFilledCircle(screen, float32(ui.pacmanX-18), pacmanY-18, 
		pacmanSize*0.2, highlight2, false)
	
	// Mouth animation
	mouthAngle := ui.pacmanMouthAngle
	if mouthAngle > 0 {
		bgColor := color.RGBA{10, 15, 40, 255}
		mouthWidth := pacmanSize * 1.0
		mouthHeight := float32(float64(pacmanSize) * math.Sin(mouthAngle))
		
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY-mouthHeight/2, 
			mouthWidth, mouthHeight/2, bgColor, false)
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY, 
			mouthWidth, mouthHeight/2, bgColor, false)
	}
	
	// Clean domain expansion preview
	if ui.selectedOption == 0 && ui.showDomainExpansion {
		ui.drawCleanDomainExpansion(screen, float32(ui.pacmanX), pacmanY)
	}
}

func (ui *UIPage) drawCharacterArea(screen *ebiten.Image) {
	charAreaX := float32(screenWidth - 320)
	charAreaY := float32(80)
	charAreaWidth := float32(280)
	charAreaHeight := float32(280)
	
	// Clean frame
	frameColor1 := color.RGBA{120, 80, 200, uint8(160 * ui.glowIntensity)}
	frameColor2 := color.RGBA{255, 215, 0, uint8(120 * ui.glowIntensity)}
	
	// Simple border
	vector.StrokeRect(screen, charAreaX-3, charAreaY-3, 
		charAreaWidth+6, charAreaHeight+6, 3, frameColor1, false)
	vector.StrokeRect(screen, charAreaX+2, charAreaY+2, 
		charAreaWidth-4, charAreaHeight-4, 2, frameColor2, false)
	
	// Large placeholder text
	centerX := int(charAreaX + charAreaWidth/2)
	centerY := int(charAreaY + charAreaHeight/2)
	
	placeholderText := "CHARACTER"
	ui.drawLargeText(screen, placeholderText, centerX-120, centerY-50, 
		color.RGBA{200, 220, 255, 200}, 3.0, ui.subtitleFont)
	
	subText := "DISPLAY AREA"
	ui.drawLargeText(screen, subText, centerX-100, centerY-10, 
		color.RGBA{160, 180, 210, 160}, 2.0, ui.uiFont)
}

func (ui *UIPage) drawFooter(screen *ebiten.Image) {
	instructions := []string{
		"↑↓ / W S  Navigate Menu",
		"ENTER / SPACE  Select Option", 
		"Experience Jujutsu Kaisen",
	}
	
	footerY := screenHeight - 80  // More space from bottom
	spacing := 380
	
	for i, instruction := range instructions {
		x := 90 + i*spacing
		y := footerY
		
		var instrColor color.RGBA
		switch i {
		case 0:
			instrColor = color.RGBA{150, 200, 255, 220}
		case 1:
			instrColor = color.RGBA{255, 215, 0, 220}
		case 2:
			pulse := 0.8 + 0.2*math.Sin(ui.animationTime*2)
			instrColor = color.RGBA{
				uint8(180 * pulse), 
				uint8(100 * pulse), 
				uint8(255 * pulse), 
				220,
			}
		}
		
		// Larger footer text
		ui.drawLargeText(screen, instruction, x, y, instrColor, 2.0, ui.uiFont)
	}
	
	// Clean decorative line
	lineY := float32(footerY - 30)
	lineColor := color.RGBA{120, 80, 200, uint8(120 * ui.glowIntensity)}
	vector.DrawFilledRect(screen, 80, lineY, screenWidth-160, 2, lineColor, false)
}

// Clean drawing helper functions
func (ui *UIPage) drawLargeText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64, fontFace font.Face) {
	if fontFace != nil && fontFace != basicfont.Face7x13 {
		// Use custom font with clean glow
		if glowIntensity > 1.0 {
			glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.3)}
			glowRadius := int(glowIntensity * 1.5) // Reduced glow radius for cleaner look
			
			// Cleaner glow layers (reduced)
			for layer := 0; layer < 2; layer++ {
				layerRadius := glowRadius - layer*int(glowRadius/2)
				if layerRadius <= 0 { continue }
				
				for dx := -layerRadius; dx <= layerRadius; dx += 3 {
					for dy := -layerRadius; dy <= layerRadius; dy += 3 {
						if dx != 0 || dy != 0 {
							distance := math.Sqrt(float64(dx*dx + dy*dy))
							if distance <= float64(layerRadius) {
								alpha := uint8(float64(glowColor.A) * 
									(1.0 - distance/float64(layerRadius)) / float64(layer+1))
								fadeColor := color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
								text.Draw(screen, txt, fontFace, x+dx, y+dy, fadeColor)
							}
						}
					}
				}
			}
		}
		
		// Main text
		text.Draw(screen, txt, fontFace, x, y, clr)
	} else {
		// Enhanced fallback with much larger scaling
		ui.drawEnhancedBasicText(screen, txt, x, y, clr, glowIntensity)
	}
}

func (ui *UIPage) drawEnhancedBasicText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	// For basic font, simulate larger size by drawing multiple times with offsets
	offsets := []struct{dx, dy int}{
		{0, 0}, {1, 0}, {0, 1}, {1, 1}, // 2x2 grid for thickness
		{2, 0}, {0, 2}, {2, 2}, {2, 1}, {1, 2}, // Additional points for larger appearance
	}
	
	// Clean glow for basic text
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.4)}
		glowRadius := int(glowIntensity)
		
		for dx := -glowRadius; dx <= glowRadius; dx++ {
			for dy := -glowRadius; dy <= glowRadius; dy++ {
				if dx != 0 || dy != 0 {
					distance := math.Sqrt(float64(dx*dx + dy*dy))
					if distance <= float64(glowRadius) {
						alpha := uint8(float64(glowColor.A) * (1.0 - distance/float64(glowRadius)))
						fadeColor := color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
						
						for _, offset := range offsets {
							text.Draw(screen, txt, basicfont.Face7x13, x+dx+offset.dx*2, y+dy+offset.dy*2+20, fadeColor)
						}
					}
				}
			}
		}
	}
	
	// Main text with simulated larger size
	for _, offset := range offsets {
		text.Draw(screen, txt, basicfont.Face7x13, x+offset.dx*2, y+offset.dy*2+20, clr)
	}
}

// Clean particle drawing functions
func (ui *UIPage) drawBackgroundParticles(screen *ebiten.Image) {
	for _, p := range ui.backgroundParticles {
		alpha := uint8(float64(p.color.A) * p.depth * 
			(0.6 + 0.3*math.Sin(ui.animationTime+p.rotation)))
		particleColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha}
		
		size := p.size * p.depth
		
		switch p.shape {
		case 0: // Circle
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(size), particleColor, false)
		case 1: // Diamond
			ui.drawSimpleDiamond(screen, p.x, p.y, size, particleColor)
		case 2: // Cross
			ui.drawSimpleCross(screen, p.x, p.y, size, particleColor)
		}
	}
}

func (ui *UIPage) drawHexagonElements(screen *ebiten.Image) {
	for _, h := range ui.hexagons {
		alpha := uint8(120 * h.alpha * ui.glowIntensity)
		hexColor := color.RGBA{150, 100, 200, alpha}
		
		ui.drawSimpleHexagon(screen, h.x, h.y, h.size, h.rotation, hexColor)
	}
}

func (ui *UIPage) drawCursedEnergy(screen *ebiten.Image) {
	for _, p := range ui.cursedEnergy {
		alpha := uint8(200 * p.energy)
		
		// Clean cursed energy without trails
		switch p.cursedType {
		case "positive":
			auraColor := color.RGBA{100, 180, 255, alpha/3}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size*2.5), auraColor, false)
		case "negative":
			auraColor := color.RGBA{255, 100, 120, alpha/3}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size*2.5), auraColor, false)
		case "neutral":
			auraColor := color.RGBA{180, 120, 255, alpha/3}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size*2.5), auraColor, false)
		}
		
		// Main particle
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(p.size), p.color, false)
	}
}

func (ui *UIPage) drawFloatingElements(screen *ebiten.Image) {
	for _, fe := range ui.floatingElements {
		alpha := uint8(200 * fe.alpha)
		elementColor := color.RGBA{160, 200, 255, alpha}
		
		// Large floating Kanji
		ui.drawLargeText(screen, fe.kanjiChar, int(fe.x-15), int(fe.y-15), 
			elementColor, 1.5, ui.subtitleFont)
	}
}

func (ui *UIPage) drawDomainExpansion(screen *ebiten.Image, dp *DomainParticle) {
	alpha := uint8(dp.alpha * dp.expansion * 180) // Reduced alpha for cleaner look
	
	switch dp.domainType {
	case "infinite_void":
		voidColor := color.RGBA{100, 150, 255, alpha}
		// Simple concentric circles
		for layer := 0; layer < 3; layer++ {
			layerRadius := float32(dp.radius * dp.expansion * float64(layer+1) * 0.4)
			layerAlpha := alpha / uint8(layer+1)
			layerColor := color.RGBA{voidColor.R, voidColor.G, voidColor.B, layerAlpha}
			
			vector.StrokeCircle(screen, float32(dp.x), float32(dp.y), 
				layerRadius, 2, layerColor, false)
		}
		
	case "malevolent_shrine":
		shrineColor := color.RGBA{255, 100, 100, alpha}
		numLines := int(dp.expansion * 8) // Fewer lines for cleaner look
		
		for i := 0; i < numLines; i++ {
			angle := dp.currentAngle + float64(i)*math.Pi*2/float64(numLines)
			length := dp.radius * dp.expansion * 0.8
			
			startX := dp.x + math.Cos(angle)*length*0.4
			startY := dp.y + math.Sin(angle)*length*0.4
			endX := dp.x + math.Cos(angle)*length
			endY := dp.y + math.Sin(angle)*length
			
			vector.StrokeLine(screen, float32(startX), float32(startY),
				float32(endX), float32(endY), 2, shrineColor, false)
		}
	}
}

func (ui *UIPage) drawCleanDomainExpansion(screen *ebiten.Image, centerX, centerY float32) {
	expansionProgress := float64(ui.domainTimer) / 300.0
	if expansionProgress > 1.0 { expansionProgress = 1.0 }
	
	// Clean domain barrier
	barrierRadius := float32(expansionProgress * 200)
	barrierAlpha := uint8(80 * expansionProgress)
	barrierColor := color.RGBA{100, 150, 255, barrierAlpha}
	
	vector.StrokeCircle(screen, centerX, centerY, barrierRadius, 3, barrierColor, false)
	
	// Simple inner effect
	innerRadius := barrierRadius * 0.6
	innerColor := color.RGBA{120, 180, 255, barrierAlpha/2}
	vector.StrokeCircle(screen, centerX, centerY, innerRadius, 2, innerColor, false)
}

// Simple shape drawing helpers
func (ui *UIPage) drawSimpleDiamond(screen *ebiten.Image, centerX, centerY, size float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cx := float32(centerX)
	cy := float32(centerY)
	
	// Simple diamond shape
	vector.StrokeLine(screen, cx, cy-halfSize, cx+halfSize, cy, 2, clr, false)
	vector.StrokeLine(screen, cx+halfSize, cy, cx, cy+halfSize, 2, clr, false)
	vector.StrokeLine(screen, cx, cy+halfSize, cx-halfSize, cy, 2, clr, false)
	vector.StrokeLine(screen, cx-halfSize, cy, cx, cy-halfSize, 2, clr, false)
}

func (ui *UIPage) drawSimpleCross(screen *ebiten.Image, centerX, centerY, size float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cx := float32(centerX)
	cy := float32(centerY)
	
	// Simple cross
	vector.StrokeLine(screen, cx-halfSize, cy, cx+halfSize, cy, 2, clr, false)
	vector.StrokeLine(screen, cx, cy-halfSize, cx, cy+halfSize, 2, clr, false)
}

func (ui *UIPage) drawSimpleHexagon(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	points := make([]float32, 12)
	for i := 0; i < 6; i++ {
		angle := rotation + float64(i)*math.Pi/3
		x := centerX + size*math.Cos(angle)
		y := centerY + size*math.Sin(angle)
		points[i*2] = float32(x)
		points[i*2+1] = float32(y)
	}
	
	// Simple hexagon outline
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen, points[i*2], points[i*2+1], 
			points[next*2], points[next*2+1], 1.5, clr, false)
	}
}

// Public interface functions
func (ui *UIPage) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (ui *UIPage) GetSelectedOption() int {
	return ui.selectedOption
}

func (ui *UIPage) IsEnterPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (ui *UIPage) SetImages(logo *ebiten.Image, gifFrames []*ebiten.Image, bg *ebiten.Image) {
	ui.logoImage = logo
	ui.gifFrames = gifFrames
	ui.bgImage = bg
	
	if len(gifFrames) > 0 {
		ui.frameIndex = 0
		ui.frameTicker = 0
	}
}

func (ui *UIPage) SetAudioSystem(audioSystem *AudioSystem) {
	ui.audioSystem = audioSystem
}

// Utility functions
func (ui *UIPage) GetMenuPulse() float64 {
	return ui.menuPulse
}

func (ui *UIPage) GetGlowIntensity() float64 {
	return ui.glowIntensity
}

func (ui *UIPage) SetScreenShake(intensity float64) {
	ui.screenShake = math.Max(ui.screenShake, intensity)
}

func (ui *UIPage) GetAnimationTime() float64 {
	return ui.animationTime
}

func (ui *UIPage) GetCurrentGIFFrame() int {
	return ui.frameIndex
}

func (ui *UIPage) SetGIFSpeed(delay int) {
	if delay > 0 {
		ui.frameDelay = delay
	}
}
