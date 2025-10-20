package main

import (
	"image/color"
	"math"
	// "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/basicfont"
	"embed"
	"log"
	// "math/rand"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

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
	selectedOption       int
	menuOptions          []string
	
	animationTime        float64
	particleTime         float64
	glowIntensity        float64
	titlePulse          float64
	menuPulse           float64
	energyWave          float64
	
	selectionTransition float64
	screenShake        float64
	
	logoImage          *ebiten.Image
	characterGif       *ebiten.Image  
	bgImage            *ebiten.Image
	gifFrames          []*ebiten.Image
	frameIndex         int
	frameTicker        int
	frameDelay         int
	
	titleFont          font.Face
	menuFont           font.Face
	subtitleFont       font.Face
	uiFont            font.Face
	
	audioSystem       *AudioSystem
	
	pacmanX             float64
	pacmanMouthAngle    float64
	
	// JJK themed particles
	cursedEnergy        []CursedEnergyParticle
	domainParticles     []DomainParticle
	ryoikiTenkai        []RyoikiEffect
	hexagons           []HexagonElement
	kanjiSymbols       []KanjiSymbol
	
	showDomainExpansion bool
	domainTimer         int
}

type CursedEnergyParticle struct {
	x, y         float64
	vx, vy       float64
	size         float64
	color        color.RGBA
	pulsePhase   float64
	energy       float64
	cursedType   string
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

type RyoikiEffect struct {
	x, y        float64
	radius      float64
	intensity   float64
	rotation    float64
	active      bool
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

type KanjiSymbol struct {
	x, y       float64
	char       string
	alpha      float64
	rotation   float64
	scale      float64
	pulsePhase float64
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

type FloatingElement struct {
	x, y          float64
	vx, vy        float64
	rotation      float64
	scale         float64
	alpha         float64
	pulsePhase    float64
	kanjiChar     string
}

type AuraParticle struct {
	x, y       float64
	intensity  float64
	auraType   string
	pulsePhase float64
	size       float64
}

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
	 img, _, err := ebitenutil.NewImageFromFile("assets/character.png")
    if err != nil {
        log.Printf("Failed to load character image: %v", err)
        return ui
    }
    ui.characterGif = img
	
	ui.initializeParticles()
	ui.initializeJJKEffects()
	ui.initializeFonts()
	
	return ui
}

// LoadCharacter loads the character image once and stores it in ui.characterGif
func (ui *UIPage) LoadCharacter() {
    img, _, err := ebitenutil.NewImageFromFile("assets/character.png")
    if err != nil {
        log.Printf("Failed to load character image: %v", err)
        ui.characterGif = nil
        return
    }
    ui.characterGif = img
}


func (ui *UIPage) initializeFonts() {
	// Use actual fonts from assets - BIGGER SIZES
	ui.titleFont = ui.loadFont("StalinistOne-Regular.ttf", 44)      // Pixel font for title
	ui.menuFont = ui.loadFont("menu-font.ttf", 36)              // Bold for menu
	ui.subtitleFont = ui.loadFont("NotoSansJP-VariableFont_wght.ttf", 28)       // Regular for subtitle
	ui.uiFont = ui.loadFont("Roboto-Regular.ttf", 22)             // UI text
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

func (ui *UIPage) initializeParticles() {
	// Cursed Energy particles - Purple/Blue theme
	ui.cursedEnergy = make([]CursedEnergyParticle, 100)
	for i := range ui.cursedEnergy {
		cursedTypes := []string{"positive", "negative", "neutral"}
		cursedType := cursedTypes[i%len(cursedTypes)]
		
		var particleColor color.RGBA
		switch cursedType {
		case "positive":
			particleColor = color.RGBA{100, 150, 255, 180} // Blue cursed energy
		case "negative":
			particleColor = color.RGBA{180, 50, 255, 180}  // Purple cursed energy
		case "neutral":
			particleColor = color.RGBA{255, 100, 200, 180} // Pink cursed energy
		}
		
		ui.cursedEnergy[i] = CursedEnergyParticle{
			x:          math.Mod(float64(i*25), screenWidth),
			y:          math.Mod(float64(i*20), screenHeight),
			vx:         (math.Sin(float64(i)) * 1.2),
			vy:         (math.Cos(float64(i)) * 0.8),
			size:       2 + math.Sin(float64(i))*2,
			pulsePhase: float64(i) * 0.1,
			color:      particleColor,
			energy:     0.5 + math.Sin(float64(i))*0.5,
			cursedType: cursedType,
		}
	}
	
	// Hexagon barriers (JJK trademark)
	ui.hexagons = make([]HexagonElement, 18)
	for i := range ui.hexagons {
		angle := float64(i) * 2 * math.Pi / float64(len(ui.hexagons))
		radius := 350.0 + math.Sin(float64(i))*80
		ui.hexagons[i] = HexagonElement{
			x:           screenWidth/2 + math.Cos(angle)*radius,
			y:           screenHeight/2 + math.Sin(angle)*radius,
			size:        20 + math.Sin(float64(i))*12,
			rotation:    angle,
			rotSpeed:    0.004 + math.Sin(float64(i))*0.003,
			alpha:       0.4 + math.Sin(float64(i))*0.3,
			pulsePhase:  float64(i) * 0.15,
			cursedLevel: (i % 3) + 1,
		}
	}
	
	// Floating Kanji symbols
	kanjiChars := []string{"呪", "術", "廻", "戦", "領", "域", "展", "開"}
	ui.kanjiSymbols = make([]KanjiSymbol, len(kanjiChars))
	for i := range ui.kanjiSymbols {
		ui.kanjiSymbols[i] = KanjiSymbol{
			x:          math.Mod(float64(i*150), screenWidth),
			y:          math.Mod(float64(i*100), screenHeight),
			char:       kanjiChars[i],
			alpha:      0.3,
			rotation:   float64(i) * 0.5,
			scale:      0.8 + math.Sin(float64(i))*0.2,
			pulsePhase: float64(i) * 0.3,
		}
	}
}

func (ui *UIPage) initializeJJKEffects() {
	// Domain Expansion effects
	ui.domainParticles = make([]DomainParticle, 3)
	for i := range ui.domainParticles {
		ui.domainParticles[i] = DomainParticle{
			x:             screenWidth/2,
			y:             screenHeight/2,
			radius:        80 + float64(i)*40,
			expansion:     0,
			domainType:    []string{"infinite_void", "malevolent_shrine", "chimera_shadow"}[i%3],
			alpha:         0.7,
			rotationSpeed: 0.02 + float64(i)*0.008,
			currentAngle:  float64(i) * math.Pi / 2,
		}
	}
	
	// Ryoiki Tenkai (Domain Expansion) effects
	ui.ryoikiTenkai = make([]RyoikiEffect, 4)
	for i := range ui.ryoikiTenkai {
		angle := float64(i) * math.Pi / 2
		ui.ryoikiTenkai[i] = RyoikiEffect{
			x:         screenWidth/2 + math.Cos(angle)*200,
			y:         screenHeight/2 + math.Sin(angle)*200,
			radius:    50,
			intensity: 0.6,
			rotation:  angle,
			active:    false,
		}
	}
}

func (ui *UIPage) Update() error {
	ui.animationTime += 0.03
	ui.particleTime += 0.02
	ui.menuPulse += 0.08
	ui.titlePulse += 0.06
	ui.energyWave += 0.04
	
	ui.glowIntensity = 0.7 + 0.3*math.Sin(ui.animationTime*2.0)
	
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
		ui.screenShake = 4.0
		if ui.audioSystem != nil {
			ui.audioSystem.PlaySFX("menu_move")
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ui.selectedOption = (ui.selectedOption + 1) % len(ui.menuOptions)
		ui.selectionTransition = 1.0
		ui.screenShake = 4.0
		if ui.audioSystem != nil {
			ui.audioSystem.PlaySFX("menu_move")
		}
	}
	
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
	ui.selectionTransition *= 0.85
	ui.screenShake *= 0.9
	
	if ui.showDomainExpansion {
		ui.domainTimer++
		if ui.domainTimer > 300 {
			ui.domainTimer = 0
		}
	}
	
	ui.pacmanX += 2.5
	if ui.pacmanX > screenWidth+250 {
		ui.pacmanX = -250
	}
	ui.pacmanMouthAngle = math.Sin(ui.animationTime*8) * 1.0
	
	if len(ui.gifFrames) > 0 {
		ui.frameTicker++
		if ui.frameTicker >= ui.frameDelay {
			ui.frameIndex = (ui.frameIndex + 1) % len(ui.gifFrames)
			ui.frameTicker = 0
		}
	}
}

func (ui *UIPage) updateParticles() {
	// Update cursed energy
	for i := range ui.cursedEnergy {
		p := &ui.cursedEnergy[i]
		
		// Attraction to selected menu
		menuY := float64(340 + ui.selectedOption*110)
		menuX := float64(screenWidth / 2)
		dx := menuX - p.x
		dy := menuY - p.y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		if distance > 0 && distance < 250 {
			force := 0.003 / (distance * 0.01)
			p.vx += (dx / distance) * force
			p.vy += (dy / distance) * force
		}
		
		p.vx *= 0.98
		p.vy *= 0.98
		p.x += p.vx
		p.y += p.vy
		
		if p.x < -50 { p.x = screenWidth + 50 }
		if p.x > screenWidth+50 { p.x = -50 }
		if p.y < -50 { p.y = screenHeight + 50 }
		if p.y > screenHeight+50 { p.y = -50 }
		
		p.energy = 0.5 + 0.5*math.Sin(ui.particleTime*2.0+p.pulsePhase)
	}
	
	// Update hexagons
	for i := range ui.hexagons {
		h := &ui.hexagons[i]
		h.rotation += h.rotSpeed
		h.alpha = 0.3 + 0.4*math.Sin(ui.animationTime*1.5+h.pulsePhase)
	}
	
	// Update Kanji symbols
	for i := range ui.kanjiSymbols {
		k := &ui.kanjiSymbols[i]
		k.rotation += 0.002
		k.alpha = 0.2 + 0.3*math.Sin(ui.animationTime+k.pulsePhase)
		k.scale = 0.8 + 0.2*math.Sin(ui.animationTime*1.2+k.pulsePhase)
		
		// Gentle float
		k.y += math.Sin(ui.animationTime*0.5+k.pulsePhase) * 0.3
	}
}

func (ui *UIPage) updateJJKEffects() {
	// Update Domain Expansion
	for i := range ui.domainParticles {
		dp := &ui.domainParticles[i]
		dp.currentAngle += dp.rotationSpeed
		
		if ui.showDomainExpansion {
			dp.expansion = math.Min(dp.expansion+0.03, 1.0)
		} else {
			dp.expansion = math.Max(dp.expansion-0.02, 0.0)
		}
	}
	
	// Update Ryoiki Tenkai
	for i := range ui.ryoikiTenkai {
		rt := &ui.ryoikiTenkai[i]
		rt.rotation += 0.01
		rt.active = ui.showDomainExpansion
		if rt.active {
			rt.intensity = 0.6 + 0.4*math.Sin(ui.animationTime*3+float64(i))
		}
	}
}

func (ui *UIPage) Draw(screen *ebiten.Image) {
	shakeX := math.Sin(ui.animationTime*15) * ui.screenShake
	shakeY := math.Cos(ui.animationTime*18) * ui.screenShake
	
	tempScreen := ebiten.NewImage(screenWidth, screenHeight)
	
	ui.drawBackground(tempScreen)
	ui.drawParticleEffects(tempScreen)
	ui.drawJJKEffects(tempScreen)
	ui.drawUI(tempScreen)
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(tempScreen, op)
}

func (ui *UIPage) drawBackground(screen *ebiten.Image) {
	// JJK themed gradient - darker with purple tones
	for y := 0; y < screenHeight; y++ {
		progress := float64(y) / float64(screenHeight)
		
		// Deep purple to dark blue gradient
		r := uint8(15 + float64(30)*progress)
		g := uint8(10 + float64(20)*progress)
		b := uint8(40 + float64(60)*progress)
		
		// Cursed energy waves
		wave := math.Sin(float64(y)*0.01 + ui.energyWave) * 10
		r = uint8(math.Max(0, math.Min(255, float64(r)+wave)))
		b = uint8(math.Max(0, math.Min(255, float64(b)+wave*1.2)))
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{r, g, b, 255}, false)
	}
	
	// Background image if available
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
		op.ColorM.Scale(1, 1, 1, 0.2)
		
		screen.DrawImage(ui.bgImage, op)
	}
// 	if int(ui.animationTime*2)%400 == 0 {
//     // Flash cursed energy streak
//     randY := rand.Float32() * float32(screenHeight)
//     vector.DrawFilledRect(screen, 0, randY, screenWidth, 2, color.RGBA{255, 100, 200, 120}, false)
// }
}

func (ui *UIPage) drawParticleEffects(screen *ebiten.Image) {
	// Draw cursed energy particles
	for _, p := range ui.cursedEnergy {
		alpha := uint8(float64(p.color.A) * p.energy)
		
		// Cursed energy aura
		auraSize := p.size * 3
		auraColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha / 4}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), float32(auraSize), auraColor, false)
		
		// Core particle
		coreColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), float32(p.size), coreColor, false)
	}
	
	// Draw hexagon barriers
	for _, h := range ui.hexagons {
		alpha := uint8(180 * h.alpha)
		hexColor := color.RGBA{150, 100, 255, alpha}
		
		ui.drawSimpleHexagon(screen, h.x, h.y, h.size, h.rotation, hexColor)
		
		// Inner glow
		innerColor := color.RGBA{200, 150, 255, alpha / 2}
		ui.drawSimpleHexagon(screen, h.x, h.y, h.size*0.7, h.rotation, innerColor)
	}
	
	// Draw floating Kanji
	for _, k := range ui.kanjiSymbols {
		alpha := uint8(220 * k.alpha)
		kanjiColor := color.RGBA{180, 150, 255, alpha}
		
		// Scale and rotate (simplified - just draw at position)
		ui.drawLargeText(screen, k.char, int(k.x), int(k.y), kanjiColor, 2.0, ui.subtitleFont)
	}
}

func (ui *UIPage) drawJJKEffects(screen *ebiten.Image) {
	// Draw Domain Expansion
	if ui.showDomainExpansion {
		for _, dp := range ui.domainParticles {
			if dp.expansion > 0 {
				ui.drawDomainExpansion(screen, &dp)
			}
		}
		
		// Draw Ryoiki Tenkai effects
		for _, rt := range ui.ryoikiTenkai {
			if rt.active {
				alpha := uint8(150 * rt.intensity)
				rtColor := color.RGBA{150, 100, 255, alpha}
				
				vector.StrokeCircle(screen, float32(rt.x), float32(rt.y), 
					float32(rt.radius), 3, rtColor, false)
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
	title := "JUJUTSU KAISEN"
	subtitle := "呪術廻戦 × PAC-MAN"
	
	titleY := 90
	titleX := screenWidth/2 - 280
	
	// Title panel with cursed energy border
	panelWidth := float32(820)
	panelHeight := float32(160)
	panelX := float32(screenWidth/2) - panelWidth/2
	panelY := float32(titleY - 50)
	// Dark panel background
	panelColor := color.RGBA{20, 15, 40, 220}
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight, panelColor, false)
	
	// // Cursed energy borders (multiple layers)
	// borderPulse := 0.8 + 0.2*math.Sin(ui.titlePulse*3)
	
	// // Outer purple glow
	// outerColor := color.RGBA{150, 80, 255, uint8(180 * borderPulse)}
	// vector.StrokeRect(screen, panelX-2, panelY-2, panelWidth+4, panelHeight+4, 3, outerColor, false)
	
	// // Inner blue glow
	// innerColor := color.RGBA{100, 150, 255, uint8(200 * borderPulse)}
	// vector.StrokeRect(screen, panelX+2, panelY+2, panelWidth-4, panelHeight-4, 2, innerColor, false)

// Enhanced cursed energy frame with layered glow
borderPulse := 0.8 + 0.2*math.Sin(ui.titlePulse*3)

outerGlow := color.RGBA{180, 120, 255, uint8(180 * borderPulse)} // bright purple
innerGlow := color.RGBA{120, 180, 255, uint8(220 * borderPulse)} // blue accent
accentColor := color.RGBA{255, 100, 200, uint8(240 * borderPulse)} // pink edge

// Outer neon glow
for i := 0; i < 4; i++ {
    alpha := uint8(80 + i*40)
    colorLayer := color.RGBA{outerGlow.R, outerGlow.G, outerGlow.B, alpha}
    vector.StrokeRect(screen, panelX-float32(i*2), panelY-float32(i*2), panelWidth+float32(i*4), panelHeight+float32(i*4), 2, colorLayer, false)
}

// Inner energy border
vector.StrokeRect(screen, panelX+3, panelY+3, panelWidth-6, panelHeight-6, 2, innerGlow, false)

// Accent corners — stylized like cursed seals
cornerSize := float32(25)
for i := 0; i < 4; i++ {
    pulse := 1.0 + 0.1*math.Sin(ui.titlePulse*4+float64(i))
    clr := color.RGBA{accentColor.R, accentColor.G, accentColor.B, uint8(200 * pulse)}

    switch i {
    case 0: // top-left
        vector.StrokeLine(screen, panelX, panelY, panelX+cornerSize, panelY, 3, clr, false)
        vector.StrokeLine(screen, panelX, panelY, panelX, panelY+cornerSize, 3, clr, false)
    case 1: // top-right
        vector.StrokeLine(screen, panelX+panelWidth, panelY, panelX+panelWidth-cornerSize, panelY, 3, clr, false)
        vector.StrokeLine(screen, panelX+panelWidth, panelY, panelX+panelWidth, panelY+cornerSize, 3, clr, false)
    case 2: // bottom-left
        vector.StrokeLine(screen, panelX, panelY+panelHeight, panelX+cornerSize, panelY+panelHeight, 3, clr, false)
        vector.StrokeLine(screen, panelX, panelY+panelHeight, panelX, panelY+panelHeight-cornerSize, 3, clr, false)
    case 3: // bottom-right
        vector.StrokeLine(screen, panelX+panelWidth, panelY+panelHeight, panelX+panelWidth-cornerSize, panelY+panelHeight, 3, clr, false)
        vector.StrokeLine(screen, panelX+panelWidth, panelY+panelHeight, panelX+panelWidth, panelY+panelHeight-cornerSize, 3, clr, false)
    }
}
glowColor := color.RGBA{100, 50, 180, 60}
vector.DrawFilledRect(screen, panelX+6, panelY+6, panelWidth-12, panelHeight-12, glowColor, false)

	
	// Accent corners (JJK style)
	cornerSize = float32(15)
	accentColor = color.RGBA{255, 100, 200, uint8(220 * borderPulse)}
	
	// Top-left corner
	vector.StrokeLine(screen, panelX, panelY, panelX+cornerSize, panelY, 3, accentColor, false)
	vector.StrokeLine(screen, panelX, panelY, panelX, panelY+cornerSize, 3, accentColor, false)
	
	// Top-right corner
	vector.StrokeLine(screen, panelX+panelWidth, panelY, panelX+panelWidth-cornerSize, panelY, 3, accentColor, false)
	vector.StrokeLine(screen, panelX+panelWidth, panelY, panelX+panelWidth, panelY+cornerSize, 3, accentColor, false)
	
	// Bottom corners
	vector.StrokeLine(screen, panelX, panelY+panelHeight, panelX+cornerSize, panelY+panelHeight, 3, accentColor, false)
	vector.StrokeLine(screen, panelX, panelY+panelHeight, panelX, panelY+panelHeight-cornerSize, 3, accentColor, false)
	
	vector.StrokeLine(screen, panelX+panelWidth, panelY+panelHeight, panelX+panelWidth-cornerSize, panelY+panelHeight, 3, accentColor, false)
	vector.StrokeLine(screen, panelX+panelWidth, panelY+panelHeight, panelX+panelWidth, panelY+panelHeight-cornerSize, 3, accentColor, false)
	
	// Title text - BIGGER
	titleColor := color.RGBA{255, 255, 255, 255}
	ui.drawLargeText(screen, title, titleX, titleY+30, titleColor, 4.0, ui.titleFont)
	
	// Subtitle - BIGGER
	subtitleY := titleY + 60
	subtitleX := screenWidth/2 - 270
	subtitleColor := color.RGBA{150, 200, 255, 255}
	ui.drawLargeText(screen, subtitle, subtitleX, subtitleY, subtitleColor, 2.5, ui.subtitleFont)

	// for i, r := range []rune(title) {
 //    alpha := uint8(255 * math.Min(1, ui.animationTime - float64(i)*0.2))
 //    text.Draw(screen, string(r), ui.titleFont, titleX+i*28, titleY, color.RGBA{255,255,255,alpha})
		// }


}

func (ui *UIPage) drawMenu(screen *ebiten.Image) {
	menuStartY := 280
	menuSpacing := 100
	menuWidth := 580
	menuX := screenWidth/2 - menuWidth/2
	
	panelHeight := float32(len(ui.menuOptions)*menuSpacing + 80)
	panelX := float32(menuX - 70)
	panelY := float32(menuStartY - 40)
	panelW := float32(menuWidth + 140)
	
	// GIF background - OPAQUE
	if len(ui.gifFrames) > 0 && ui.gifFrames[ui.frameIndex] != nil {
		ui.drawGIFBackground(screen, panelX, panelY, panelW, panelHeight)
	}
	
	// Dark overlay (minimal)
	overlayColor := color.RGBA{15, 10, 30, 180}
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelHeight, overlayColor, false)
	
	// Cursed energy borders
	borderPulse := 0.8 + 0.2*math.Sin(ui.menuPulse*2)
	borderColor := color.RGBA{150, 100, 255, uint8(200 * borderPulse)}
	vector.StrokeRect(screen, panelX, panelY, panelW, panelHeight, 3, borderColor, false)
	
	accentColor := color.RGBA{255, 100, 200, uint8(180 * borderPulse)}
	vector.StrokeRect(screen, panelX+3, panelY+3, panelW-6, panelHeight-6, 2, accentColor, false)
	
	// Menu options - BIGGER TEXT
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
	selectionHeight := float32(80)
	selectionX := float32(x - 45)
	selectionY := float32(y - 28)
	
	pulseIntensity := 0.85 + 0.15*math.Sin(ui.menuPulse*4)
	
	// Cursed energy selection background
	mainAlpha := uint8(200 * pulseIntensity)
	
	// Gradient with cursed energy colors
	for i := 0; i < int(selectionHeight); i++ {
		progress := float64(i) / float64(selectionHeight)
		r := uint8(100 + 50*progress)
		g := uint8(30 + 70*progress)
		b := uint8(180 + 40*progress)
		lineColor := color.RGBA{r, g, b, mainAlpha}
		vector.DrawFilledRect(screen, selectionX, selectionY+float32(i), selectionWidth, 1, lineColor, false)
	}
	
	// Multiple border layers (JJK style)
	neonColor := color.RGBA{150, 100, 255, uint8(255 * pulseIntensity)}
	vector.StrokeRect(screen, selectionX, selectionY, selectionWidth, selectionHeight, 3, neonColor, false)
	
	accentBorder := color.RGBA{255, 100, 200, uint8(220 * pulseIntensity)}
	vector.StrokeRect(screen, selectionX+2, selectionY+2, selectionWidth-4, selectionHeight-4, 2, accentBorder, false)
	
	// Selected text - BIGGER
	textColor := color.RGBA{255, 255, 255, 255}
	// ui.drawLargeText(screen, option, x, y, textColor, 4.0, ui.menuFont)
	
	// Cursed energy indicators
	ui.drawCleanSelectionIndicators(screen, selectionX, selectionY, selectionWidth, selectionHeight)
	 ui.drawCenteredText(screen, option, screenWidth/2, y+25, textColor, ui.menuFont)
}

func (ui *UIPage) drawUnselectedOption(screen *ebiten.Image, option string, x, y, index int) {
	hoverIntensity := 0.75 + 0.15*math.Sin(ui.animationTime*1.5+float64(index)*0.5)
	textColor := color.RGBA{
		uint8(170 * hoverIntensity),
		uint8(180 * hoverIntensity),
		uint8(230 * hoverIntensity),
		210,
	}
	// BIGGER unselected text
	// ui.drawLargeText(screen, option, x, y, textColor, 2.5, ui.menuFont)
	 ui.drawCenteredText(screen, option, screenWidth/2, y+25, textColor, ui.menuFont)
}

func (ui *UIPage) drawCleanSelectionIndicators(screen *ebiten.Image, x, y, w, h float32) {
	indicatorY := y + h/2
	
	pulsePhase := ui.menuPulse * 5
	intensity := 0.8 + 0.2*math.Sin(pulsePhase)
	
	// Left cursed energy indicator
	leftX := x - 25
	
	// Outer glow
	glowColor := color.RGBA{150, 100, 255, uint8(120 * intensity)}
	vector.DrawFilledCircle(screen, leftX, indicatorY, 12, glowColor, false)
	
	// Core
	coreColor := color.RGBA{255, 100, 200, uint8(255 * intensity)}
	vector.DrawFilledCircle(screen, leftX, indicatorY, 7, coreColor, false)
	
	// Center highlight
	highlightColor := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, leftX, indicatorY, 3, highlightColor, false)
	
	// Right indicator
	rightX := x + w + 25
	vector.DrawFilledCircle(screen, rightX, indicatorY, 12, glowColor, false)
	vector.DrawFilledCircle(screen, rightX, indicatorY, 7, coreColor, false)
	vector.DrawFilledCircle(screen, rightX, indicatorY, 3, highlightColor, false)
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
	tempOp.ColorM.Scale(1, 1, 1, 1.0) // FULLY OPAQUE
	
	tempImg.DrawImage(currentFrame, tempOp)
	
	finalOp := &ebiten.DrawImageOptions{}
	finalOp.GeoM.Translate(float64(panelX), float64(panelY))
	screen.DrawImage(tempImg, finalOp)
}

func (ui *UIPage) drawPacman(screen *ebiten.Image) {
	pacmanY := float32(screenHeight/2 + 80)
	pacmanSize := float32(70)
	
	// Cursed energy trail
	trailLength := 8
	for i := 0; i < trailLength; i++ {
		trailX := float32(ui.pacmanX) - float32(i*20)
		trailAlpha := uint8(140 - i*15)
		trailSize := pacmanSize * (1.0 - float32(i)*0.1)
		
		// Purple cursed energy trail
		trailColor := color.RGBA{150, 100, 255, trailAlpha}
		vector.DrawFilledCircle(screen, trailX, pacmanY, trailSize*1.2, trailColor, false)
	}
	
	// Main cursed energy aura
	auraLayers := 4
	for layer := 0; layer < auraLayers; layer++ {
		layerSize := pacmanSize + float32(layer*15) + float32(20*ui.glowIntensity)
		layerAlpha := uint8(float64(100) / float64(layer + 1) * ui.glowIntensity)
		
		auraColor := color.RGBA{180, 120, 255, layerAlpha}
		vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, layerSize, auraColor, false)
	}
	
	// Pac-Man body with JJK colors
	pacmanColor := color.RGBA{255, 220, 80, 255}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
	// Highlight
	highlightColor := color.RGBA{255, 255, 255, 200}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX-12), pacmanY-12, pacmanSize*0.3, highlightColor, false)
	
	// Mouth animation
	mouthAngle := ui.pacmanMouthAngle
	if mouthAngle > 0 {
		bgColor := color.RGBA{15, 10, 30, 255}
		mouthWidth := pacmanSize * 1.3
		mouthHeight := float32(float64(pacmanSize) * math.Sin(mouthAngle))
		
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY-mouthHeight/2, 
			mouthWidth, mouthHeight/2, bgColor, false)
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY, 
			mouthWidth, mouthHeight/2, bgColor, false)
	}
	
	// Domain expansion around Pac-Man
	if ui.selectedOption == 0 && ui.showDomainExpansion {
		ui.drawCleanDomainExpansion(screen, float32(ui.pacmanX), pacmanY)
	}
}

func (ui *UIPage) drawCleanDomainExpansion(screen *ebiten.Image, centerX, centerY float32) {
	expansionProgress := float64(ui.domainTimer) / 300.0
	if expansionProgress > 1.0 {
		expansionProgress = 1.0
	}
	
	// Multiple expanding rings (Domain Expansion effect)
	for ring := 0; ring < 4; ring++ {
		ringProgress := math.Max(0, expansionProgress-float64(ring)*0.15)
		radius := float32(ringProgress * 220)
		alpha := uint8((1.0 - ringProgress) * 180)
		
		// Alternate colors for visual depth
		var ringColor color.RGBA
		if ring%2 == 0 {
			ringColor = color.RGBA{150, 100, 255, alpha}
		} else {
			ringColor = color.RGBA{255, 100, 200, alpha}
		}
		
		vector.StrokeCircle(screen, centerX, centerY, radius, 3, ringColor, false)
	}
	
	// Hexagon pattern for domain
	if expansionProgress > 0.3 {
		hexSize := float32(expansionProgress * 100)
		hexColor := color.RGBA{200, 150, 255, uint8(expansionProgress * 100)}
		ui.drawSimpleHexagon(screen, float64(centerX), float64(centerY), float64(hexSize), ui.animationTime, hexColor)
	}
}

func (ui *UIPage) drawCharacterArea(screen *ebiten.Image) {
	charAreaX := float32(screenWidth - 320)
	charAreaY := float32(75)
	charAreaWidth := float32(270)
	charAreaHeight := float32(270)
	
	// Dark background
	// bgColor := color.RGBA{20, 15, 40, 220}
	// vector.DrawFilledRect(screen, charAreaX, charAreaY, charAreaWidth, charAreaHeight, bgColor, false)
	
	// Cursed energy frame
	borderPulse := 0.8 + 0.2*math.Sin(ui.animationTime*2)
	// frameColor1 := color.RGBA{150, 100, 255, uint8(200 * borderPulse)}
	// frameColor2 := color.RGBA{255, 100, 200, uint8(180 * borderPulse)}
	
	// vector.StrokeRect(screen, charAreaX, charAreaY, charAreaWidth, charAreaHeight, 3, frameColor1, false)
	// vector.StrokeRect(screen, charAreaX+3, charAreaY+3, charAreaWidth-6, charAreaHeight-6, 2, frameColor2, false)
	
	// Corner accents (JJK style)
	cornerSize := float32(20)
	accentColor := color.RGBA{255, 255, 255, uint8(200 * borderPulse)}
	
	// Top-left
	vector.StrokeLine(screen, charAreaX, charAreaY, charAreaX+cornerSize, charAreaY, 3, accentColor, false)
	vector.StrokeLine(screen, charAreaX, charAreaY, charAreaX, charAreaY+cornerSize, 3, accentColor, false)
	
	// Top-right
	vector.StrokeLine(screen, charAreaX+charAreaWidth, charAreaY, charAreaX+charAreaWidth-cornerSize, charAreaY, 3, accentColor, false)
	vector.StrokeLine(screen, charAreaX+charAreaWidth, charAreaY, charAreaX+charAreaWidth, charAreaY+cornerSize, 3, accentColor, false)
	
	// Bottom-left
	vector.StrokeLine(screen, charAreaX, charAreaY+charAreaHeight, charAreaX+cornerSize, charAreaY+charAreaHeight, 3, accentColor, false)
	vector.StrokeLine(screen, charAreaX, charAreaY+charAreaHeight, charAreaX, charAreaY+charAreaHeight-cornerSize, 3, accentColor, false)
	
	// Bottom-right
	vector.StrokeLine(screen, charAreaX+charAreaWidth, charAreaY+charAreaHeight, charAreaX+charAreaWidth-cornerSize, charAreaY+charAreaHeight, 3, accentColor, false)
	vector.StrokeLine(screen, charAreaX+charAreaWidth, charAreaY+charAreaHeight, charAreaX+charAreaWidth, charAreaY+charAreaHeight-cornerSize, 3, accentColor, false)
	
	// Placeholder text - BIGGER
	centerX := int(charAreaX + charAreaWidth/2)
	centerY := int(charAreaY + charAreaHeight/2)

	// characterImg, _,_ := ebitenutil.NewImageFromFile("assets/character.png")
	// ui.characterGif = characterImg

	
	if ui.characterGif != nil {
    	op := &ebiten.DrawImageOptions{}
    	imgBounds := ui.characterGif.Bounds()
    	imgWidth := float64(imgBounds.Dx())
    	imgHeight := float64(imgBounds.Dy())

    	// Scale and center
    	scale := 330.0 / math.Max(imgWidth, imgHeight)
  	    op.GeoM.Scale(scale, scale)
    	op.GeoM.Translate(float64(centerX)-imgWidth*scale/2, float64(centerY)-imgHeight*scale/2)
    	screen.DrawImage(ui.characterGif, op)
	} else {
    	// Fallback placeholder text
    	// ui.drawLargeText(screen, "CHARACTER", centerX-100, centerY-35,
     //    color.RGBA{200, 220, 255, 200}, 2.5, ui.subtitleFont)
	}
		
	// subText := "DISPLAY"
	// ui.drawLargeText(screen, subText, centerX-55, centerY+5,
	// 	color.RGBA{160, 180, 210, 160}, 2.0, ui.uiFont)
}

func (ui *UIPage) drawFooter(screen *ebiten.Image) {
	instructions := []string{
		"↑↓ / W S  NAVIGATE",
		"ENTER / SPACE  SELECT",
		"JUJUTSU KAISEN EXPERIENCE",
	}
	
	footerY := screenHeight - 70
	spacing := 380
	
	for i, instruction := range instructions {
		x := 90 + i*spacing
		y := footerY
		
		var instrColor color.RGBA
		switch i {
		case 0:
			instrColor = color.RGBA{150, 200, 255, 230}
		case 1:
			instrColor = color.RGBA{255, 200, 100, 230}
		case 2:
			pulse := 0.85 + 0.15*math.Sin(ui.animationTime*3)
			instrColor = color.RGBA{
				uint8(170 * pulse),
				uint8(120 * pulse),
				uint8(255 * pulse),
				230,
			}
		}
		
		// BIGGER footer text
		ui.drawLargeText(screen, instruction, x, y, instrColor, 2.0, ui.uiFont)
	}
	
	// Cursed energy separator line
	lineY := float32(footerY - 30)
	segments := 120
	
	for i := 0; i < segments; i++ {
		progress := float64(i) / float64(segments)
		x := float32(80 + i*(screenWidth-160)/segments)
		
		// Animated wave pattern
		wave := math.Sin(progress*math.Pi*4 + ui.energyWave)
		alpha := uint8(120 + 80*math.Abs(wave))
		
		// Alternate cursed energy colors
		var lineColor color.RGBA
		if int(progress*10)%2 == 0 {
			lineColor = color.RGBA{150, 100, 255, alpha}
		} else {
			lineColor = color.RGBA{255, 100, 200, alpha}
		}
		
		vector.DrawFilledRect(screen, x, lineY, float32((screenWidth-160)/segments)+1, 2, lineColor, false)
	}
}

func (ui *UIPage) drawLargeText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64, fontFace font.Face) {
	if fontFace != nil && fontFace != basicfont.Face7x13 {
		// Cursed energy glow effect
		if glowIntensity > 2.5 {
			// Strong glow for selected items
			glowColor := color.RGBA{150, 100, 255, 100}
			
			for dx := -3; dx <= 3; dx++ {
				for dy := -3; dy <= 3; dy++ {
					if dx != 0 || dy != 0 {
						distance := math.Sqrt(float64(dx*dx + dy*dy))
						if distance <= 3 {
							alpha := uint8(100 * (1.0 - distance/3.0))
							fadeColor := color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
							text.Draw(screen, txt, fontFace, x+dx, y+dy, fadeColor)
						}
					}
				}
			}
		} else if glowIntensity > 1.5 {
			// Subtle glow
			glowColor := color.RGBA{clr.R, clr.G, clr.B, 80}
			text.Draw(screen, txt, fontFace, x+1, y+1, glowColor)
			text.Draw(screen, txt, fontFace, x-1, y-1, glowColor)
		}
		
		// Main text
		text.Draw(screen, txt, fontFace, x, y, clr)
	} else {
		ui.drawEnhancedBasicText(screen, txt, x, y, clr, glowIntensity)
	}
}

func (ui *UIPage) drawEnhancedBasicText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	offsets := []struct{dx, dy int}{
		{0, 0}, {1, 0}, {0, 1}, {1, 1},
		{2, 0}, {0, 2}, {2, 2}, {2, 1}, {1, 2},
	}
	
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.5)}
		glowRadius := int(glowIntensity * 1.5)
		
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
	
	for _, offset := range offsets {
		text.Draw(screen, txt, basicfont.Face7x13, x+offset.dx*2, y+offset.dy*2+20, clr)
	}
}

func (ui *UIPage) drawDomainExpansion(screen *ebiten.Image, dp *DomainParticle) {
	alpha := uint8(dp.alpha * dp.expansion * 200)
	
	switch dp.domainType {
	case "infinite_void":
		// Gojo's Infinite Void - Blue/white
		voidColor := color.RGBA{100, 180, 255, alpha}
		for layer := 0; layer < 4; layer++ {
			layerRadius := float32(dp.radius * dp.expansion * float64(layer+1) * 0.35)
			layerAlpha := alpha / uint8(layer+1)
			layerColor := color.RGBA{voidColor.R, voidColor.G, voidColor.B, layerAlpha}
			
			vector.StrokeCircle(screen, float32(dp.x), float32(dp.y), layerRadius, 2, layerColor, false)
		}
		
	case "malevolent_shrine":
		// Sukuna's Malevolent Shrine - Red/pink
		shrineColor := color.RGBA{255, 80, 120, alpha}
		numLines := int(dp.expansion * 16)
		
		for i := 0; i < numLines; i++ {
			angle := dp.currentAngle + float64(i)*math.Pi*2/float64(numLines)
			length := dp.radius * dp.expansion * 0.8
			
			startX := dp.x + math.Cos(angle)*length*0.25
			startY := dp.y + math.Sin(angle)*length*0.25
			endX := dp.x + math.Cos(angle)*length
			endY := dp.y + math.Sin(angle)*length
			
			vector.StrokeLine(screen, float32(startX), float32(startY),
				float32(endX), float32(endY), 2, shrineColor, false)
		}
		
	case "chimera_shadow":
		// Megumi's Shadow - Dark purple
		shadowColor := color.RGBA{120, 80, 180, alpha}
		size := float32(dp.radius * dp.expansion)
		
		for ring := 0; ring < 4; ring++ {
			ringSize := size * float32(ring+1) * 0.25
			ui.drawSimpleHexagon(screen, dp.x, dp.y, float64(ringSize), dp.currentAngle+float64(ring)*0.3, shadowColor)
		}
	}
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
	
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen, points[i*2], points[i*2+1], 
			points[next*2], points[next*2+1], 2, clr, false)
	}
}

func (ui *UIPage) drawSimpleDiamond(screen *ebiten.Image, centerX, centerY, size float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cx := float32(centerX)
	cy := float32(centerY)
	
	vector.StrokeLine(screen, cx, cy-halfSize, cx+halfSize, cy, 2, clr, false)
	vector.StrokeLine(screen, cx+halfSize, cy, cx, cy+halfSize, 2, clr, false)
	vector.StrokeLine(screen, cx, cy+halfSize, cx-halfSize, cy, 2, clr, false)
	vector.StrokeLine(screen, cx-halfSize, cy, cx, cy-halfSize, 2, clr, false)
}

func (ui *UIPage) drawSimpleCross(screen *ebiten.Image, centerX, centerY, size float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cx := float32(centerX)
	cy := float32(centerY)
	
	vector.StrokeLine(screen, cx-halfSize, cy, cx+halfSize, cy, 2, clr, false)
	vector.StrokeLine(screen, cx, cy-halfSize, cx, cy+halfSize, 2, clr, false)
}

// Legacy compatibility functions
// func (ui *UIPage) drawTitle(screen *ebiten.Image) {
// 	ui.drawTitle(screen)
// }

// func (ui *UIPage) drawMenu(screen *ebiten.Image) {
// 	ui.drawMenu(screen)
// }

// func (ui *UIPage) drawSelectedOption(screen *ebiten.Image, option string, x, y, menuWidth int) {
// 	ui.drawSelectedOption(screen, option, x, y, menuWidth)
// }

// func (ui *UIPage) drawUnselectedOption(screen *ebiten.Image, option string, x, y, index int) {
// 	ui.drawUnselectedOption(screen, option, x, y, index)
// }

// func (ui *UIPage) drawCharacterArea(screen *ebiten.Image) {
// 	ui.drawCharacterArea(screen)
// }

// func (ui *UIPage) drawFooter(screen *ebiten.Image) {
// 	ui.drawFooter(screen)
// }

// func (ui *UIPage) drawPacman(screen *ebiten.Image) {
// 	ui.drawPacman(screen)
// }

func (ui *UIPage) drawCenteredText(screen *ebiten.Image, txt string, centerX, centerY int, clr color.RGBA, fontFace font.Face) {
    b, _ := font.BoundString(fontFace, txt)
    width := (b.Max.X - b.Min.X).Ceil()
    height := (b.Max.Y - b.Min.Y).Ceil()
    text.Draw(screen, txt, fontFace, centerX - width/2, centerY + height/2, clr)
}

func (ui *UIPage) drawBackgroundParticles(screen *ebiten.Image) {}
func (ui *UIPage) drawHexagonElements(screen *ebiten.Image) {}
func (ui *UIPage) drawCursedEnergy(screen *ebiten.Image) {}
func (ui *UIPage) drawFloatingElements(screen *ebiten.Image) {}

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

func (ui *UIPage) SetImages(logo *ebiten.Image, gifFrames []*ebiten.Image, bg *ebiten.Image,screen *ebiten.Image) {
	ui.logoImage = logo
	ui.gifFrames = gifFrames
	ui.bgImage = bg

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.6 + 0.2*math.Sin(ui.animationTime))
	// screen.DrawImage(ui.characterGif, op)

	if len(gifFrames) > 0 {
		ui.frameIndex = 0
		ui.frameTicker = 0
	}
}

func (ui *UIPage) SetAudioSystem(audioSystem *AudioSystem) {
	ui.audioSystem = audioSystem
}

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
