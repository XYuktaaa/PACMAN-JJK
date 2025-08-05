package main

import (
	"image/color"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

type MenuState int

const (
	MenuStart MenuState = iota
	MenuPause
	MenuResume
	MenuQuit
)

type UIPage struct {
	selectedOption       int
	animationTime        float64
	particleTime         float64
	glowIntensity        float64
	menuOptions          []string
	pacmanX             float64
	pacmanMouthAngle    float64
	cursedEnergy        []CursedEnergyParticle
	backgroundParticles []BackgroundParticle
	hexagons           []HexagonElement
	transitionOffset    float64
	logoImage          *ebiten.Image  // Space for logo image
	characterGif       *ebiten.Image  // Space for character GIF
	backgroundTexture  *ebiten.Image  // Space for background texture
}

type CursedEnergyParticle struct {
	x, y         float64
	vx, vy       float64
	life         float64
	maxLife      float64
	size         float64
	color        color.RGBA
	pulsePhase   float64
	trail        []TrailPoint
}

type TrailPoint struct {
	x, y  float64
	alpha float64
}

type BackgroundParticle struct {
	x, y       float64
	vx, vy     float64
	size       float64
	rotation   float64
	rotSpeed   float64
	color      color.RGBA
	shape      int // 0=circle, 1=diamond, 2=cross
}

type HexagonElement struct {
	x, y        float64
	size        float64
	rotation    float64
	rotSpeed    float64
	alpha       float64
	pulsePhase  float64
}

func NewUIPage() *UIPage {
	ui := &UIPage{
		selectedOption:      0,
		menuOptions:        []string{"START GAME", "SETTINGS", "GALLERY", "EXIT"},
		pacmanX:           -150,
		cursedEnergy:      make([]CursedEnergyParticle, 80),
		backgroundParticles: make([]BackgroundParticle, 60),
		hexagons:          make([]HexagonElement, 12),
	}
	
	// Initialize enhanced cursed energy particles with trails
	for i := range ui.cursedEnergy {
		ui.cursedEnergy[i] = CursedEnergyParticle{
			x:          math.Mod(float64(i*20), screenWidth),
			y:          math.Mod(float64(i*15), screenHeight),
			vx:         (math.Sin(float64(i)) * 1.2),
			vy:         (math.Cos(float64(i)) * 0.8),
			life:       1.0,
			maxLife:    1.0,
			size:       3 + math.Sin(float64(i))*2,
			pulsePhase: float64(i) * 0.1,
			color:      color.RGBA{148, 0, 211, 180}, // Deep purple
			trail:      make([]TrailPoint, 8),
		}
	}
	
	// Initialize modern background particles
	for i := range ui.backgroundParticles {
		ui.backgroundParticles[i] = BackgroundParticle{
			x:        math.Mod(float64(i*25), screenWidth),
			y:        math.Mod(float64(i*18), screenHeight),
			vx:       (math.Sin(float64(i)*0.1) * 0.3),
			vy:       (math.Cos(float64(i)*0.1) * 0.2),
			size:     2 + math.Sin(float64(i))*3,
			rotation: float64(i) * 0.1,
			rotSpeed: 0.01 + math.Sin(float64(i))*0.005,
			shape:    i % 3,
			color:    color.RGBA{65, 105, 225, 60}, // Steel blue with transparency
		}
	}
	
	// Initialize hexagonal UI elements
	for i := range ui.hexagons {
		angle := float64(i) * 2 * math.Pi / float64(len(ui.hexagons))
		radius := 200.0 + math.Sin(float64(i))*50
		ui.hexagons[i] = HexagonElement{
			x:          screenWidth/2 + math.Cos(angle)*radius,
			y:          screenHeight/2 + math.Sin(angle)*radius,
			size:       20 + math.Sin(float64(i))*10,
			rotation:   angle,
			rotSpeed:   0.005 + math.Sin(float64(i))*0.003,
			alpha:      0.3 + math.Sin(float64(i))*0.2,
			pulsePhase: float64(i) * 0.2,
		}
	}
	
	return ui
}

func (ui *UIPage) Update() error {
	ui.animationTime += 0.04
	ui.particleTime += 0.03
	ui.glowIntensity = 0.6 + 0.4*math.Sin(ui.animationTime*1.5)
	ui.transitionOffset = math.Sin(ui.animationTime*0.5) * 20
	
	// Update enhanced Pacman animation
	ui.pacmanX += 1.5
	if ui.pacmanX > screenWidth+200 {
		ui.pacmanX = -200
	}
	ui.pacmanMouthAngle = math.Sin(ui.animationTime*6) * 0.8
	
	// Handle input with smooth transitions
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		ui.selectedOption = (ui.selectedOption - 1 + len(ui.menuOptions)) % len(ui.menuOptions)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ui.selectedOption = (ui.selectedOption + 1) % len(ui.menuOptions)
	}
	
	// Update cursed energy particles with trails
	for i := range ui.cursedEnergy {
		p := &ui.cursedEnergy[i]
		
		// Update trail
		for j := len(p.trail) - 1; j > 0; j-- {
			p.trail[j] = p.trail[j-1]
			p.trail[j].alpha *= 0.9
		}
		p.trail[0] = TrailPoint{x: p.x, y: p.y, alpha: 1.0}
		
		p.x += p.vx * (1 + 0.3*math.Sin(ui.animationTime+p.pulsePhase))
		p.y += p.vy * (1 + 0.2*math.Cos(ui.animationTime+p.pulsePhase))
		
		// Wrap around with smooth transitions
		if p.x < -50 { p.x = screenWidth + 50 }
		if p.x > screenWidth+50 { p.x = -50 }
		if p.y < -50 { p.y = screenHeight + 50 }
		if p.y > screenHeight+50 { p.y = -50 }
		
		// Enhanced pulsing
		p.life = 0.4 + 0.6*math.Sin(ui.particleTime*2+p.pulsePhase)
	}
	
	// Update background particles
	for i := range ui.backgroundParticles {
		p := &ui.backgroundParticles[i]
		p.x += p.vx
		p.y += p.vy
		p.rotation += p.rotSpeed
		
		if p.x < -20 { p.x = screenWidth + 20 }
		if p.x > screenWidth+20 { p.x = -20 }
		if p.y < -20 { p.y = screenHeight + 20 }
		if p.y > screenHeight+20 { p.y = -20 }
	}
	
	// Update hexagons
	for i := range ui.hexagons {
		h := &ui.hexagons[i]
		h.rotation += h.rotSpeed
		h.alpha = 0.2 + 0.3*math.Sin(ui.animationTime*1.2+h.pulsePhase)
	}
	
	return nil
}

func (ui *UIPage) Draw(screen *ebiten.Image) {
	// Draw modern gradient background
	ui.drawModernBackground(screen)
	
	// Draw background particles
	ui.drawBackgroundParticles(screen)
	
	// Draw hexagonal UI elements
	ui.drawHexagonElements(screen)
	
	// Draw cursed energy with trails
	ui.drawEnhancedCursedEnergy(screen)
	
	// Draw character area (reserved for GIFs/images)
	ui.drawCharacterArea(screen)
	
	// Draw modern title with effects
	ui.drawModernTitle(screen)
	
	// Draw sleek menu
	ui.drawSleekMenu(screen)
	
	// Draw animated Pacman with modern effects
	ui.drawModernPacman(screen)
	
	// Draw modern UI elements
	ui.drawModernUIElements(screen)
	
	// Draw footer with instructions
	ui.drawModernFooter(screen)
}

func (ui *UIPage) drawModernBackground(screen *ebiten.Image) {
	// Create a sophisticated gradient background
	for y := 0; y < screenHeight; y++ {
		progress := float64(y) / float64(screenHeight)
		
		// Multi-layered gradient
		r1, g1, b1 := 8, 12, 25    // Dark navy at top
		r2, g2, b2 := 25, 15, 45   // Deep purple at bottom
		
		r := uint8(float64(r1) + (float64(r2-r1))*progress)
		g := uint8(float64(g1) + (float64(g2-g1))*progress)
		b := uint8(float64(b1) + (float64(b2-b1))*progress)
		
		// Add subtle animation
		waveOffset := math.Sin(float64(y)*0.01 + ui.animationTime*0.5) * 5
		r = uint8(math.Max(0, math.Min(255, float64(r)+waveOffset)))
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{r, g, b, 255}, false)
	}
	
	// Add subtle scan lines for modern effect
	for y := 0; y < screenHeight; y += 4 {
		alpha := uint8(20 + 10*math.Sin(ui.animationTime*2+float64(y)*0.1))
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{100, 150, 255, alpha}, false)
	}
}

func (ui *UIPage) drawBackgroundParticles(screen *ebiten.Image) {
	for _, p := range ui.backgroundParticles {
		alpha := uint8(p.color.A * (uint8)(0.5 + 0.5*math.Sin(ui.animationTime+p.rotation)))
		particleColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha}
		
		switch p.shape {
		case 0: // Circle
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size), particleColor, false)
		case 1: // Diamond
			ui.drawDiamond(screen, p.x, p.y, p.size, p.rotation, particleColor)
		case 2: // Cross
			ui.drawCross(screen, p.x, p.y, p.size, p.rotation, particleColor)
		}
	}
}

func (ui *UIPage) drawHexagonElements(screen *ebiten.Image) {
	for _, h := range ui.hexagons {
		alpha := uint8(255 * h.alpha * ui.glowIntensity)
		hexColor := color.RGBA{148, 0, 211, alpha}
		
		ui.drawHexagon(screen, h.x, h.y, h.size, h.rotation, hexColor)
	}
}

func (ui *UIPage) drawEnhancedCursedEnergy(screen *ebiten.Image) {
	for _, p := range ui.cursedEnergy {
		// Draw trail
		for i, trail := range p.trail {
			if i > 0 && trail.alpha > 0.1 {
				size := p.size * trail.alpha * 0.8
				alpha := uint8(150 * trail.alpha * p.life)
				trailColor := color.RGBA{148, 0, 211, alpha}
				
				vector.DrawFilledCircle(screen, float32(trail.x), float32(trail.y), 
					float32(size), trailColor, false)
			}
		}
		
		// Draw main particle with enhanced effects
		alpha := uint8(200 * p.life)
		mainColor := color.RGBA{148, 0, 211, alpha}
		
		// Outer glow
		glowSize := p.size * 2.5
		glowAlpha := uint8(alpha / 3)
		glowColor := color.RGBA{200, 100, 255, glowAlpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(glowSize), glowColor, false)
		
		// Main particle
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(p.size), mainColor, false)
		
		// Inner core
		if p.life > 0.8 {
			coreColor := color.RGBA{255, 255, 255, uint8(alpha/2)}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size*0.3), coreColor, false)
		}
	}
}

func (ui *UIPage) drawCharacterArea(screen *ebiten.Image) {
	// Reserved area for character GIF/image (top-right)
	charAreaX := float32(screenWidth - 300)
	charAreaY := float32(50)
	charAreaWidth := float32(250)
	charAreaHeight := float32(250)
	
	// Draw placeholder frame with modern styling
	frameColor := color.RGBA{100, 150, 255, 80}
	vector.StrokeRect(screen, charAreaX, charAreaY, charAreaWidth, charAreaHeight, 3, frameColor, false)
	
	// Add corner accents
	cornerSize := float32(20)
	accentColor := color.RGBA{255, 215, 0, 200}
	
	// Top-left corner
	vector.DrawFilledRect(screen, charAreaX-2, charAreaY-2, cornerSize, 3, accentColor, false)
	vector.DrawFilledRect(screen, charAreaX-2, charAreaY-2, 3, cornerSize, accentColor, false)
	
	// Top-right corner
	vector.DrawFilledRect(screen, charAreaX+charAreaWidth-cornerSize+2, charAreaY-2, cornerSize, 3, accentColor, false)
	vector.DrawFilledRect(screen, charAreaX+charAreaWidth-1, charAreaY-2, 3, cornerSize, accentColor, false)
	
	// Placeholder text
	if ui.characterGif == nil {
		placeholderText := "CHARACTER"
		textX := int(charAreaX + charAreaWidth/2 - float32(len(placeholderText)*4))
		textY := int(charAreaY + charAreaHeight/2)
		ui.drawGlowText(screen, placeholderText, textX, textY, 
			color.RGBA{150, 150, 200, 150}, 1.0)
		
		subText := "ANIMATION AREA"
		subTextX := int(charAreaX + charAreaWidth/2 - float32(len(subText)*3))
		subTextY := textY + 20
		ui.drawGlowText(screen, subText, subTextX, subTextY, 
			color.RGBA{100, 100, 150, 100}, 0.8)
	}
}

func (ui *UIPage) drawModernTitle(screen *ebiten.Image) {
	// Main title with enhanced effects
	title := "呪術廻戦 × PAC-MAN"
	subtitle := "JUJUTSU KAISEN EDITION"
	
	titleY := 100
	titleX := screenWidth/2 - len(title)*16
	
	// Title background panel
	panelWidth := float32(len(title) * 32 + 60)
	panelHeight := float32(80)
	panelX := float32(titleX - 30)
	panelY := float32(titleY - 20)
	
	// Draw panel with glassmorphism effect
	panelColor := color.RGBA{20, 30, 60, 120}
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight, panelColor, false)
	
	// Panel border with glow
	borderColor := color.RGBA{100, 150, 255, uint8(150 * ui.glowIntensity)}
	vector.StrokeRect(screen, panelX, panelY, panelWidth, panelHeight, 2, borderColor, false)
	
	// Draw title with enhanced glow
	ui.drawGlowText(screen, title, titleX, titleY, 
		color.RGBA{255, 215, 0, 255}, 4.0)
	
	// Subtitle with modern styling
	subtitleY := titleY + 45
	subtitleX := screenWidth/2 - len(subtitle)*6
	ui.drawGlowText(screen, subtitle, subtitleX, subtitleY, 
		color.RGBA{150, 200, 255, 255}, 2.0)
	
	// Logo area (top-left)
	logoAreaX := float32(50)
	logoAreaY := float32(50)
	logoSize := float32(80)
	
	if ui.logoImage == nil {
		// Placeholder logo frame
		logoColor := color.RGBA{255, 215, 0, 150}
		vector.StrokeRect(screen, logoAreaX, logoAreaY, logoSize, logoSize, 2, logoColor, false)
		
		// Placeholder text
		ui.drawGlowText(screen, "LOGO", int(logoAreaX+20), int(logoAreaY+40), 
			color.RGBA{255, 215, 0, 200}, 1.5)
	}
}

func (ui *UIPage) drawSleekMenu(screen *ebiten.Image) {
	menuStartY := 350
	menuSpacing := 80
	menuWidth := 400
	menuX := screenWidth/2 - menuWidth/2
	
	// Draw menu background panel
	panelHeight := float32(len(ui.menuOptions)*menuSpacing + 40)
	panelColor := color.RGBA{15, 25, 45, 180}
	vector.DrawFilledRect(screen, float32(menuX-40), float32(menuStartY-20), 
		float32(menuWidth+80), panelHeight, panelColor, false)
	
	// Panel border
	borderColor := color.RGBA{100, 150, 255, 100}
	vector.StrokeRect(screen, float32(menuX-40), float32(menuStartY-20), 
		float32(menuWidth+80), panelHeight, 1, borderColor, false)
	
	for i, option := range ui.menuOptions {
		y := menuStartY + i*menuSpacing
		x := menuX
		
		if i == ui.selectedOption {
			// Selected option with modern styling
			selectionWidth := float32(menuWidth)
			selectionHeight := float32(60)
			selectionX := float32(x - 20)
			selectionY := float32(y - 15)
			
			// Animated selection background
			glowIntensity := ui.glowIntensity
			selectionColor := color.RGBA{148, 0, 211, uint8(150 * glowIntensity)}
			vector.DrawFilledRect(screen, selectionX, selectionY, 
				selectionWidth, selectionHeight, selectionColor, false)
			
			// Selection border with glow
			borderGlow := color.RGBA{255, 215, 0, uint8(200 * glowIntensity)}
			vector.StrokeRect(screen, selectionX, selectionY, 
				selectionWidth, selectionHeight, 2, borderGlow, false)
			
			// Option text with enhanced glow
			ui.drawGlowText(screen, option, x, y, 
				color.RGBA{255, 255, 255, 255}, 3.0)
			
			// Modern selection indicator
			indicatorX := x - 40
			indicatorY := float32(y + 5)
			ui.drawModernArrow(screen, float32(indicatorX), indicatorY, 
				color.RGBA{255, 215, 0, 255})
			
			// Add side accent lines
			accentColor := color.RGBA{255, 215, 0, uint8(200 * glowIntensity)}
			vector.DrawFilledRect(screen, selectionX-5, selectionY, 3, selectionHeight, accentColor, false)
			vector.DrawFilledRect(screen, selectionX+selectionWidth+2, selectionY, 3, selectionHeight, accentColor, false)
		} else {
			// Unselected options with subtle styling
			ui.drawGlowText(screen, option, x, y, 
				color.RGBA{180, 180, 220, 200}, 1.2)
		}
	}
}

func (ui *UIPage) drawModernPacman(screen *ebiten.Image) {
	pacmanY := float32(screenHeight/2 - 50)
	pacmanSize := float32(60)
	
	// Enhanced aura effect
	auraLayers := 3
	for layer := 0; layer < auraLayers; layer++ {
		layerSize := pacmanSize + float32(layer*15) + float32(20*ui.glowIntensity)
		layerAlpha := uint8((60) /float64 (layer + 1) * ui.glowIntensity)
		auraColor := color.RGBA{148, 0, 211, layerAlpha}
		
		vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, 
			layerSize, auraColor, false)
	}
	
	// Main Pacman body with gradient effect
	pacmanColor := color.RGBA{255, 215, 0, 255}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
	// Inner highlight
	highlightColor := color.RGBA{255, 255, 200, 180}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX-10), pacmanY-10, 
		pacmanSize*0.3, highlightColor, false)
	
	// Enhanced mouth animation
	mouthAngle := ui.pacmanMouthAngle
	if mouthAngle > 0 {
		bgColor := color.RGBA{15, 25, 45, 255}
		
		// Create more realistic mouth shape
		mouthWidth := pacmanSize * 0.8
		mouthHeight := float32(float64(pacmanSize) * math.Sin(mouthAngle))
		
		// Upper mouth
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY-mouthHeight/2, 
			mouthWidth, mouthHeight/2, bgColor, false)
		// Lower mouth  
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY, 
			mouthWidth, mouthHeight/2, bgColor, false)
	}
	
	// Add cursed energy particles around Pacman
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4 + ui.animationTime*2
		distance := 80 + 20*math.Sin(ui.animationTime*3+float64(i))
		px := ui.pacmanX + math.Cos(angle)*distance
		py := float64(pacmanY) + math.Sin(angle)*distance
		
		particleSize := float32(5 + 3*math.Sin(ui.animationTime*4+float64(i)))
		particleAlpha := uint8(150 + 100*math.Sin(ui.animationTime*2+float64(i)))
		particleColor := color.RGBA{148, 0, 211, particleAlpha}
		
		vector.DrawFilledCircle(screen, float32(px), float32(py), 
			particleSize, particleColor, false)
	}
}

func (ui *UIPage) drawModernUIElements(screen *ebiten.Image) {
	// Draw modern side panels
	ui.drawSidePanels(screen)
	
	// Draw tech-style decorative elements
	ui.drawTechDecorations(screen)
	
	// Draw status indicators
	ui.drawStatusIndicators(screen)
}

func (ui *UIPage) drawSidePanels(screen *ebiten.Image) {
	panelWidth := float32(120)
	panelHeight := float32(screenHeight - 100)
	
	// Left panel
	leftPanelColor := color.RGBA{10, 20, 40, 100}
	vector.DrawFilledRect(screen, 25, 50, panelWidth, panelHeight, leftPanelColor, false)
	
	leftBorderColor := color.RGBA{100, 150, 255, 80}
	vector.StrokeRect(screen, 25, 50, panelWidth, panelHeight, 1, leftBorderColor, false)
	
	// Right panel
	rightX := float32(screenWidth - 145)
	rightPanelColor := color.RGBA{10, 20, 40, 100}
	vector.DrawFilledRect(screen, rightX, 50, panelWidth, panelHeight, rightPanelColor, false)
	
	rightBorderColor := color.RGBA{100, 150, 255, 80}
	vector.StrokeRect(screen, rightX, 50, panelWidth, panelHeight, 1, rightBorderColor, false)
}

func (ui *UIPage) drawTechDecorations(screen *ebiten.Image) {
	// Corner tech elements
	corners := []struct{ x, y float32 }{
		{60, 80}, {screenWidth - 180, 80},
		{60, screenHeight - 120}, {screenWidth - 180, screenHeight - 120},
	}
	
	for _, corner := range corners {
		techColor := color.RGBA{100, 150, 255, uint8(120 * ui.glowIntensity)}
		
		// Tech corner brackets
		bracketSize := float32(25)
		vector.DrawFilledRect(screen, corner.x, corner.y, bracketSize, 2, techColor, false)
		vector.DrawFilledRect(screen, corner.x, corner.y, 2, bracketSize, techColor, false)
		
		vector.DrawFilledRect(screen, corner.x, corner.y+50, bracketSize, 2, techColor, false)
		vector.DrawFilledRect(screen, corner.x, corner.y+30, 2, bracketSize, techColor, false)
	}
}

func (ui *UIPage) drawStatusIndicators(screen *ebiten.Image) {
	// Bottom status bar
	statusY := float32(screenHeight - 60)
	statusColor := color.RGBA{20, 30, 50, 150}
	vector.DrawFilledRect(screen, 50, statusY, screenWidth-100, 40, statusColor, false)
	
	// Status border
	borderColor := color.RGBA{100, 150, 255, 100}
	vector.StrokeRect(screen, 50, statusY, screenWidth-100, 40, 1, borderColor, false)
	
	// Animated status dots
	for i := 0; i < 5; i++ {
		dotX := float32(70 + i*30)
		dotY := statusY + 20
		dotSize := float32(4 + 2*math.Sin(ui.animationTime*3+float64(i)*0.5))
		
		dotAlpha := uint8(150 + 100*math.Sin(ui.animationTime*2+float64(i)*0.3))
		dotColor := color.RGBA{148, 0, 211, dotAlpha}
		
		vector.DrawFilledCircle(screen, dotX, dotY, dotSize, dotColor, false)
	}
}

func (ui *UIPage) drawModernFooter(screen *ebiten.Image) {
	instructions := []string{
		"↑↓ / W S  Navigate",
		"ENTER  Select",
		"Experience the Cursed Energy",
	}
	
	startY := screenHeight - 40
	for i, instruction := range instructions {
		x := 80 + i*300
		y := startY
		
		instrColor := color.RGBA{150, 180, 220, 180}
		if i == 2 {
			instrColor = color.RGBA{148, 0, 211, 200}
		}
		
		ui.drawGlowText(screen, instruction, x, y, instrColor, 1.0)
	}
}

// Helper drawing functions
func (ui *UIPage) drawGlowText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.3)}
		glowRadius := int(glowIntensity)
		for dx := -glowRadius; dx <= glowRadius; dx++ {
			for dy := -glowRadius; dy <= glowRadius; dy++ {
				if dx != 0 || dy != 0 {
					distance := math.Sqrt(float64(dx*dx + dy*dy))
					if distance <= float64(glowRadius) {
						alpha := uint8(float64(glowColor.A) * (1.0 - distance/float64(glowRadius)))
						fadeColor := color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
						text.Draw(screen, txt, basicfont.Face7x13, x+dx, y+dy+13, fadeColor)
					}
				}
			}
		}
	}
	
	text.Draw(screen, txt, basicfont.Face7x13, x, y+13, clr)
}

func (ui *UIPage) drawHexagon(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	points := make([]float32, 12)
	for i := 0; i < 6; i++ {
		angle := rotation + float64(i)*math.Pi/3
		x := centerX + size*math.Cos(angle)
		y := centerY + size*math.Sin(angle)
		points[i*2] = float32(x)
		points[i*2+1] = float32(y)
	}
	
	// Draw hexagon outline
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen, points[i*2], points[i*2+1], 
			points[next*2], points[next*2+1], 2, clr, false)
	}
}

func (ui *UIPage) drawDiamond(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cosR := float32(math.Cos(rotation))
	sinR := float32(math.Sin(rotation))
	
	// Diamond vertices
	vertices := []struct{ x, y float32 }{
		{0, -halfSize},  // Top
		{halfSize, 0},   // Right
		{0, halfSize},   // Bottom
		{-halfSize, 0},  // Left
	}
	
	// Rotate and draw diamond
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		
		// Rotate current vertex
		x1 := vertices[i].x*cosR - vertices[i].y*sinR + float32(centerX)
		y1 := vertices[i].x*sinR + vertices[i].y*cosR + float32(centerY)
		
		// Rotate next vertex
		x2 := vertices[next].x*cosR - vertices[next].y*sinR + float32(centerX)
		y2 := vertices[next].x*sinR + vertices[next].y*cosR + float32(centerY)
		
		vector.StrokeLine(screen, x1, y1, x2, y2, 1.5, clr, false)
	}
}

func (ui *UIPage) drawCross(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cosR := float32(math.Cos(rotation))
	sinR := float32(math.Sin(rotation))
	
	cx := float32(centerX)
	cy := float32(centerY)
	
	// Horizontal line
	x1 := -halfSize*cosR + cx
	y1 := -halfSize*sinR + cy
	x2 := halfSize*cosR + cx
	y2 := halfSize*sinR + cy
	vector.StrokeLine(screen, x1, y1, x2, y2, 2, clr, false)
	
	// Vertical line
	x3 := halfSize*sinR + cx
	y3 := -halfSize*cosR + cy
	x4 := -halfSize*sinR + cx
	y4 := halfSize*cosR + cy
	vector.StrokeLine(screen, x3, y3, x4, y4, 2, clr, false)
}

func (ui *UIPage) drawModernArrow(screen *ebiten.Image, x, y float32, clr color.RGBA) {
	// Modern triangular arrow indicator
	arrowSize := float32(12)
	
	// Triangle points
	x1 := x
	y1 := y - arrowSize/2
	x2 := x
	y2 := y + arrowSize/2
	x3 := x + arrowSize
	y3 := y
	
	// Draw filled triangle
	// Since we don't have a direct filled triangle function, we'll use lines
	vector.StrokeLine(screen, x1, y1, x2, y2, 2, clr, false)
	vector.StrokeLine(screen, x1, y1, x3, y3, 2, clr, false)
	vector.StrokeLine(screen, x2, y2, x3, y3, 2, clr, false)
	
	// Add inner fill lines
	for i := 0; i < int(arrowSize); i++ {
		progress := float32(i) / arrowSize
		startY := y1 + (y2-y1)*progress
		endX := x + arrowSize*progress
		vector.StrokeLine(screen, x, startY, endX, y, 1, clr, false)
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
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// SetImages allows setting the images/gifs for the UI
func (ui *UIPage) SetImages(logo, character, background *ebiten.Image) {
	ui.logoImage = logo
	ui.characterGif = character
	ui.backgroundTexture = background
}

/*
Enhanced Integration Example:

type Game struct {
    menuUI *UIPage
    showMenu bool
    logoImg *ebiten.Image
    characterGif *ebiten.Image
    bgTexture *ebiten.Image
}

func main() {
    // Load your images/gifs
    logo, _ := LoadImage("assets/jjk_logo.png")
    character, _ := LoadGIF("assets/gojo_animation.gif") 
    bg, _ := LoadImage("assets/cursed_energy_bg.png")
    
    game := &Game{
        menuUI: NewUIPage(),
        showMenu: true,
        logoImg: logo,
        characterGif: character,
        bgTexture: bg,
    }
    
    // Set the images in the UI
    game.menuUI.SetImages(logo, character, bg)
    
    ebiten.SetWindowSize(1200, 800)
    ebiten.SetWindowTitle("Jujutsu Kaisen Pac-Man")
    ebiten.RunGame(game)
}

Key Features Added:
- Modern glassmorphism UI panels
- Enhanced particle systems with trails
- Hexagonal UI elements for tech aesthetic
- Reserved spaces for logo, character GIFs, and background
- Improved gradients and lighting effects
- Better input handling (WASD + Arrow keys)
- Status indicators and decorative elements
- Responsive design elements
- Enhanced cursed energy effects
- Modern selection indicators
*/
