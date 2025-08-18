package main

import (
	"image/color"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	logoImage          *ebiten.Image
	characterGif       *ebiten.Image  
	bgImage            *ebiten.Image
	
	// Enhanced GIF handling
	gifFrames          []*ebiten.Image
	frameIndex         int
	frameTicker        int
	frameDelay         int
	
	// New UI enhancements
	menuPulse          float64
	selectionTransition float64
	energyField        []EnergyFieldParticle
	lightRays          []LightRay
	floatingElements   []FloatingElement
	screenShake        float64
	nebulaClouds       []NebulaCloud
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
	energy       float64
	magneticForce float64
}

type TrailPoint struct {
	x, y  float64
	alpha float64
	size  float64
}

type BackgroundParticle struct {
	x, y       float64
	vx, vy     float64
	size       float64
	rotation   float64
	rotSpeed   float64
	color      color.RGBA
	shape      int
	depth      float64 // For parallax effect
}

type HexagonElement struct {
	x, y        float64
	size        float64
	rotation    float64
	rotSpeed    float64
	alpha       float64
	pulsePhase  float64
	glowRadius  float64
}

type EnergyFieldParticle struct {
	x, y           float64
	targetX, targetY float64
	speed          float64
	phase          float64
	intensity      float64
	connectionLines []ConnectionLine
}

type ConnectionLine struct {
	targetIndex int
	strength    float64
	pulsePhase  float64
}

type LightRay struct {
	startX, startY float64
	endX, endY     float64
	width          float64
	alpha          float64
	color          color.RGBA
	animPhase      float64
}

type FloatingElement struct {
	x, y          float64
	vx, vy        float64
	rotation      float64
	rotSpeed      float64
	scale         float64
	scaleSpeed    float64
	elementType   int // 0=kanji, 1=symbol, 2=geometric
	alpha         float64
	pulsePhase    float64
}

type NebulaCloud struct {
	x, y       float64
	size       float64
	alpha      float64
	color      color.RGBA
	driftSpeed float64
	pulsePhase float64
}

func NewUIPage() *UIPage {
	ui := &UIPage{
		selectedOption:      0,
		menuOptions:        []string{"START GAME", "SETTINGS", "GALLERY", "EXIT"},
		pacmanX:           -150,
		cursedEnergy:      make([]CursedEnergyParticle, 120),
		backgroundParticles: make([]BackgroundParticle, 80),
		hexagons:          make([]HexagonElement, 18),
		energyField:       make([]EnergyFieldParticle, 25),
		lightRays:         make([]LightRay, 8),
		floatingElements:  make([]FloatingElement, 15),
		nebulaClouds:      make([]NebulaCloud, 6),
		frameDelay:        6, // Faster GIF animation
		selectionTransition: 0,
	}
	
	// Initialize enhanced cursed energy particles
	for i := range ui.cursedEnergy {
		ui.cursedEnergy[i] = CursedEnergyParticle{
			x:          math.Mod(float64(i*20), screenWidth),
			y:          math.Mod(float64(i*15), screenHeight),
			vx:         (math.Sin(float64(i)) * 1.5),
			vy:         (math.Cos(float64(i)) * 1.2),
			life:       1.0,
			maxLife:    1.0,
			size:       2 + math.Sin(float64(i))*3,
			pulsePhase: float64(i) * 0.1,
			color:      color.RGBA{148, 0, 211, 200},
			trail:      make([]TrailPoint, 12),
			energy:     0.5 + math.Sin(float64(i))*0.5,
			magneticForce: 0.02 + math.Sin(float64(i)*0.3)*0.01,
		}
	}
	
	// Initialize background particles with depth
	for i := range ui.backgroundParticles {
		ui.backgroundParticles[i] = BackgroundParticle{
			x:        math.Mod(float64(i*25), screenWidth),
			y:        math.Mod(float64(i*18), screenHeight),
			vx:       (math.Sin(float64(i)*0.1) * 0.4),
			vy:       (math.Cos(float64(i)*0.1) * 0.3),
			size:     1 + math.Sin(float64(i))*4,
			rotation: float64(i) * 0.1,
			rotSpeed: 0.008 + math.Sin(float64(i))*0.006,
			shape:    i % 4, // Added more shapes
			depth:    0.3 + math.Sin(float64(i))*0.7,
			color:    color.RGBA{65, 105, 225, 80},
		}
	}
	
	// Initialize enhanced hexagons
	for i := range ui.hexagons {
		angle := float64(i) * 2 * math.Pi / float64(len(ui.hexagons))
		radius := 180.0 + math.Sin(float64(i))*60
		ui.hexagons[i] = HexagonElement{
			x:          screenWidth/2 + math.Cos(angle)*radius,
			y:          screenHeight/2 + math.Sin(angle)*radius,
			size:       15 + math.Sin(float64(i))*12,
			rotation:   angle,
			rotSpeed:   0.003 + math.Sin(float64(i))*0.004,
			alpha:      0.4 + math.Sin(float64(i))*0.3,
			pulsePhase: float64(i) * 0.15,
			glowRadius: 25 + math.Sin(float64(i))*15,
		}
	}
	
	// Initialize energy field
	for i := range ui.energyField {
		ui.energyField[i] = EnergyFieldParticle{
			x:         math.Mod(float64(i*60), screenWidth),
			y:         math.Mod(float64(i*45), screenHeight),
			targetX:   math.Mod(float64(i*60), screenWidth),
			targetY:   math.Mod(float64(i*45), screenHeight),
			speed:     0.02 + math.Sin(float64(i))*0.015,
			phase:     float64(i) * 0.2,
			intensity: 0.6 + math.Sin(float64(i))*0.4,
		}
	}
	
	// Initialize light rays
	for i := range ui.lightRays {
		angle := float64(i) * 2 * math.Pi / float64(len(ui.lightRays))
		ui.lightRays[i] = LightRay{
			startX:    screenWidth / 2,
			startY:    screenHeight / 2,
			endX:      screenWidth/2 + math.Cos(angle)*600,
			endY:      screenHeight/2 + math.Sin(angle)*600,
			width:     2 + math.Sin(float64(i))*3,
			alpha:     0.3 + math.Sin(float64(i))*0.2,
			color:     color.RGBA{255, 215, 0, 100},
			animPhase: float64(i) * 0.3,
		}
	}
	
	// Initialize floating elements
	for i := range ui.floatingElements {
		ui.floatingElements[i] = FloatingElement{
			x:           math.Mod(float64(i*80), screenWidth),
			y:           math.Mod(float64(i*60), screenHeight),
			vx:          (math.Sin(float64(i)) * 0.5),
			vy:          (math.Cos(float64(i)) * 0.3),
			rotation:    float64(i) * 0.5,
			rotSpeed:    0.01 + math.Sin(float64(i))*0.008,
			scale:       0.5 + math.Sin(float64(i))*0.3,
			scaleSpeed:  0.005 + math.Sin(float64(i))*0.003,
			elementType: i % 3,
			alpha:       0.4 + math.Sin(float64(i))*0.3,
			pulsePhase:  float64(i) * 0.25,
		}
	}
	
	// Initialize nebula clouds
	for i := range ui.nebulaClouds {
		ui.nebulaClouds[i] = NebulaCloud{
			x:          math.Mod(float64(i*200), screenWidth),
			y:          math.Mod(float64(i*150), screenHeight),
			size:       80 + math.Sin(float64(i))*40,
			alpha:      0.15 + math.Sin(float64(i))*0.1,
			color:      color.RGBA{75, 0, 130, 30}, // Indigo
			driftSpeed: 0.2 + math.Sin(float64(i))*0.15,
			pulsePhase: float64(i) * 0.4,
		}
	}
	
	return ui
}

func (ui *UIPage) Update() error {
	ui.animationTime += 0.035
	ui.particleTime += 0.025
	ui.menuPulse += 0.08
	ui.glowIntensity = 0.7 + 0.3*math.Sin(ui.animationTime*1.8)
	ui.transitionOffset = math.Sin(ui.animationTime*0.6) * 15
	
	// Smooth selection transition
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		ui.selectedOption = (ui.selectedOption - 1 + len(ui.menuOptions)) % len(ui.menuOptions)
		ui.selectionTransition = 1.0
		ui.screenShake = 5.0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ui.selectedOption = (ui.selectedOption + 1) % len(ui.menuOptions)
		ui.selectionTransition = 1.0
		ui.screenShake = 5.0
	}
	
	// Decay transitions
	ui.selectionTransition *= 0.85
	ui.screenShake *= 0.9
	
	// Enhanced Pacman animation
	ui.pacmanX += 2.0
	if ui.pacmanX > screenWidth+250 {
		ui.pacmanX = -250
	}
	ui.pacmanMouthAngle = math.Sin(ui.animationTime*8) * 0.9
	
	// Update enhanced cursed energy with magnetic attraction
	for i := range ui.cursedEnergy {
		p := &ui.cursedEnergy[i]
		
		// Magnetic attraction to selected menu item
		menuY := float64(350 + ui.selectedOption*80)
		menuX := float64(screenWidth/2)
		dx := menuX - p.x
		dy := menuY - p.y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		if distance > 0 && distance < 200 {
			force := p.magneticForce / (distance * 0.01)
			p.vx += (dx / distance) * force
			p.vy += (dy / distance) * force
		}
		
		// Apply velocity damping
		p.vx *= 0.98
		p.vy *= 0.98
		
		// Update trail with size variation
		for j := len(p.trail) - 1; j > 0; j-- {
			p.trail[j] = p.trail[j-1]
			p.trail[j].alpha *= 0.88
			p.trail[j].size *= 0.95
		}
		p.trail[0] = TrailPoint{
			x: p.x, 
			y: p.y, 
			alpha: 1.0,
			size: p.size,
		}
		
		p.x += p.vx * (1 + 0.4*math.Sin(ui.animationTime+p.pulsePhase))
		p.y += p.vy * (1 + 0.3*math.Cos(ui.animationTime+p.pulsePhase))
		
		// Enhanced wrapping
		if p.x < -80 { p.x = screenWidth + 80 }
		if p.x > screenWidth+80 { p.x = -80 }
		if p.y < -80 { p.y = screenHeight + 80 }
		if p.y > screenHeight+80 { p.y = -80 }
		
		// Dynamic energy levels
		p.energy = 0.3 + 0.7*math.Sin(ui.particleTime*1.5+p.pulsePhase)
		p.life = 0.5 + 0.5*math.Sin(ui.particleTime*2.2+p.pulsePhase)
	}
	
	// Update background particles with parallax
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
	
	// Update enhanced hexagons
	for i := range ui.hexagons {
		h := &ui.hexagons[i]
		h.rotation += h.rotSpeed
		h.alpha = 0.3 + 0.4*math.Sin(ui.animationTime*1.4+h.pulsePhase)
		h.glowRadius = 20 + 15*math.Sin(ui.animationTime*0.8+h.pulsePhase)
	}
	
	// Update energy field
	for i := range ui.energyField {
		ef := &ui.energyField[i]
		ef.phase += 0.02
		
		// Orbital motion around center
		centerX := screenWidth / 2
		centerY := screenHeight / 2
		radius := 150 + math.Sin(ef.phase)*50
		angle := ef.phase + float64(i)*0.4
		
		ef.targetX = float64(centerX) + math.Cos(angle)*radius
		ef.targetY = float64(centerY) + math.Sin(angle)*radius
		
		// Smooth movement toward target
		ef.x += (ef.targetX - ef.x) * ef.speed
		ef.y += (ef.targetY - ef.y) * ef.speed
		
		ef.intensity = 0.4 + 0.6*math.Sin(ui.animationTime*1.6+float64(i)*0.3)
	}
	
	// Update light rays
	for i := range ui.lightRays {
		lr := &ui.lightRays[i]
		lr.animPhase += 0.04
		lr.alpha = 0.2 + 0.3*math.Sin(ui.animationTime*2+lr.animPhase)
		
		// Rotate rays slowly
		angle := lr.animPhase + float64(i)*math.Pi/4
		lr.endX = screenWidth/2 + math.Cos(angle)*700
		lr.endY = screenHeight/2 + math.Sin(angle)*700
	}
	
	// Update floating elements
	for i := range ui.floatingElements {
		fe := &ui.floatingElements[i]
		fe.x += fe.vx
		fe.y += fe.vy
		fe.rotation += fe.rotSpeed
		fe.scale = 0.6 + 0.4*math.Sin(ui.animationTime*fe.scaleSpeed+fe.pulsePhase)
		fe.alpha = 0.3 + 0.4*math.Sin(ui.animationTime*1.3+fe.pulsePhase)
		
		// Wrap around
		if fe.x < -50 { fe.x = screenWidth + 50 }
		if fe.x > screenWidth+50 { fe.x = -50 }
		if fe.y < -50 { fe.y = screenHeight + 50 }
		if fe.y > screenHeight+50 { fe.y = -50 }
	}
	
	// Update nebula clouds
	for i := range ui.nebulaClouds {
		nc := &ui.nebulaClouds[i]
		nc.x += nc.driftSpeed
		nc.alpha = 0.1 + 0.15*math.Sin(ui.animationTime*0.7+nc.pulsePhase)
		
		if nc.x > screenWidth+nc.size {
			nc.x = -nc.size
		}
	}

	// Enhanced GIF frame updating
	if len(ui.gifFrames) > 0 {
		ui.frameTicker++
		if ui.frameTicker >= ui.frameDelay {
			ui.frameIndex = (ui.frameIndex + 1) % len(ui.gifFrames)
			ui.frameTicker = 0
		}
	}
	
	return nil
}

func (ui *UIPage) Draw(screen *ebiten.Image) {
	// Add screen shake effect
	shakeX := (math.Sin(ui.animationTime*15) * ui.screenShake)
	shakeY := (math.Cos(ui.animationTime*18) * ui.screenShake)
	
	// Create temporary image for shake effect
	tempScreen := ebiten.NewImage(screenWidth, screenHeight)
	
	// Draw all elements to temp screen
	ui.drawEnhancedBackground(tempScreen)
	ui.drawNebulaClouds(tempScreen)
	ui.drawLightRays(tempScreen)
	ui.drawBackgroundParticles(tempScreen)
	ui.drawEnergyField(tempScreen)
	ui.drawHexagonElements(tempScreen)
	ui.drawEnhancedCursedEnergy(tempScreen)
	ui.drawFloatingElements(tempScreen)
	ui.drawCharacterArea(tempScreen)
	ui.drawEnhancedTitle(tempScreen)
	ui.drawEnhancedMenu(tempScreen)
	ui.drawEnhancedPacman(tempScreen)
	ui.drawEnhancedUIElements(tempScreen)
	ui.drawEnhancedFooter(tempScreen)
	
	// Apply shake and draw to main screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shakeX, shakeY)
	screen.DrawImage(tempScreen, op)
}

func (ui *UIPage) drawEnhancedBackground(screen *ebiten.Image) {
	// Background image with enhanced blending
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
		op.ColorM.Scale(1, 1, 1, 0.5) // Reduced opacity for better text readability
		
		screen.DrawImage(ui.bgImage, op)
	}
	
	// Enhanced gradient with multiple layers
	for y := 0; y < screenHeight; y++ {
		progress := float64(y) / float64(screenHeight)
		
		// Primary gradient
		r1, g1, b1 := 5, 8, 20    // Darker navy
		r2, g2, b2 := 30, 10, 50  // Deeper purple
		
		r := uint8(float64(r1) + (float64(r2-r1))*progress)
		g := uint8(float64(g1) + (float64(g2-g1))*progress)
		b := uint8(float64(b1) + (float64(b2-b1))*progress)
		
		// Add animated color shifts
		waveR := math.Sin(float64(y)*0.008 + ui.animationTime*0.4) * 8
		waveG := math.Sin(float64(y)*0.012 + ui.animationTime*0.6) * 5
		waveB := math.Sin(float64(y)*0.006 + ui.animationTime*0.3) * 12
		
		r = uint8(math.Max(0, math.Min(255, float64(r)+waveR)))
		g = uint8(math.Max(0, math.Min(255, float64(g)+waveG)))
		b = uint8(math.Max(0, math.Min(255, float64(b)+waveB)))
		
		alpha := uint8(180)
		if ui.bgImage != nil {
			alpha = 120
		}
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, 
			color.RGBA{r, g, b, alpha}, false)
	}
	
	// Enhanced scan lines with variation
	for y := 0; y < screenHeight; y += 3 {
		intensity := math.Sin(ui.animationTime*3+float64(y)*0.02)
		alpha := uint8(15 + 20*intensity)
		if ui.bgImage != nil {
			alpha /= 2
		}
		
		scanColor := color.RGBA{120, 180, 255, alpha}
		if y%6 == 0 {
			scanColor = color.RGBA{255, 215, 0, alpha/2}
		}
		
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, scanColor, false)
	}
}

func (ui *UIPage) drawNebulaClouds(screen *ebiten.Image) {
	for _, nc := range ui.nebulaClouds {
		// Create soft cloud effect with multiple circles
		for layer := 0; layer < 5; layer++ {
			layerOffset := float64(layer) * 15
			layerAlpha := nc.alpha / float64(layer+1)
			layerSize := nc.size + layerOffset
			
			cloudColor := color.RGBA{
				nc.color.R, 
				nc.color.G, 
				nc.color.B, 
				uint8(255 * layerAlpha),
			}
			
			vector.DrawFilledCircle(screen, 
				float32(nc.x + math.Sin(ui.animationTime*0.3+float64(layer))*10), 
				float32(nc.y + math.Cos(ui.animationTime*0.4+float64(layer))*8), 
				float32(layerSize), cloudColor, false)
		}
	}
}

func (ui *UIPage) drawLightRays(screen *ebiten.Image) {
	for _, lr := range ui.lightRays {
		alpha := uint8(255 * lr.alpha * ui.glowIntensity)
		rayColor := color.RGBA{lr.color.R, lr.color.G, lr.color.B, alpha}
		
		// Draw ray with gradient effect
		segments := 20
		for i := 0; i < segments; i++ {
			t := float64(i) / float64(segments)
			x := float32(lr.startX + (lr.endX-lr.startX)*t)
			y := float32(lr.startY + (lr.endY-lr.startY)*t)
			
			segmentAlpha := uint8(float64(alpha) * (1.0 - t*0.7))
			segmentColor := color.RGBA{rayColor.R, rayColor.G, rayColor.B, segmentAlpha}
			segmentWidth := float32(lr.width * (1.0 - t*0.5))
			
			vector.DrawFilledCircle(screen, x, y, segmentWidth, segmentColor, false)
		}
	}
}

func (ui *UIPage) drawBackgroundParticles(screen *ebiten.Image) {
	for _, p := range ui.backgroundParticles {
		// Parallax alpha based on depth
		alpha := uint8(float64(p.color.A) * p.depth * 
			(0.6 + 0.4*math.Sin(ui.animationTime+p.rotation)))
		particleColor := color.RGBA{p.color.R, p.color.G, p.color.B, alpha}
		
		// Scale based on depth
		size := p.size * p.depth
		
		switch p.shape {
		case 0: // Circle
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(size), particleColor, false)
		case 1: // Diamond
			ui.drawDiamond(screen, p.x, p.y, size, p.rotation, particleColor)
		case 2: // Cross
			ui.drawCross(screen, p.x, p.y, size, p.rotation, particleColor)
		case 3: // Star
			ui.drawStar(screen, p.x, p.y, size, p.rotation, particleColor)
		}
	}
}

func (ui *UIPage) drawEnergyField(screen *ebiten.Image) {
	// Draw connections between energy field particles
	for i := range ui.energyField {
		ef1 := &ui.energyField[i]
		
		for j := i + 1; j < len(ui.energyField); j++ {
			ef2 := &ui.energyField[j]
			distance := math.Sqrt(math.Pow(ef1.x-ef2.x, 2) + math.Pow(ef1.y-ef2.y, 2))
			
			if distance < 150 {
				alpha := uint8(100 * (1.0 - distance/150) * ef1.intensity * ef2.intensity)
				lineColor := color.RGBA{148, 0, 211, alpha}
				
				vector.StrokeLine(screen, float32(ef1.x), float32(ef1.y),
					float32(ef2.x), float32(ef2.y), 1, lineColor, false)
			}
		}
		
		// Draw energy field particle
		particleSize := float32(6 + 4*math.Sin(ui.animationTime*2+ef1.phase))
		particleAlpha := uint8(200 * ef1.intensity)
		particleColor := color.RGBA{255, 100, 255, particleAlpha}
		
		vector.DrawFilledCircle(screen, float32(ef1.x), float32(ef1.y), 
			particleSize, particleColor, false)
	}
}

func (ui *UIPage) drawHexagonElements(screen *ebiten.Image) {
	for _, h := range ui.hexagons {
		// Enhanced glow effect
		glowAlpha := uint8(100 * h.alpha * ui.glowIntensity)
	    // glowColor := color.RGBA{148, 0, 211, glowAlpha}
		
		// Draw glow layers
		for layer := 0; layer < 3; layer++ {
			layerAlpha := glowAlpha / uint8(layer+1)
			layerSize := h.size + h.glowRadius*float64(layer)*0.3
			layerColor := color.RGBA{200, 100, 255, layerAlpha}
			
			ui.drawHexagon(screen, h.x, h.y, layerSize, h.rotation, layerColor)
		}
		
		// Main hexagon
		mainAlpha := uint8(255 * h.alpha)
		hexColor := color.RGBA{148, 0, 211, mainAlpha}
		ui.drawHexagon(screen, h.x, h.y, h.size, h.rotation, hexColor)
	}
}

func (ui *UIPage) drawEnhancedCursedEnergy(screen *ebiten.Image) {
	for _, p := range ui.cursedEnergy {
		// Enhanced trail with varying sizes
		for i, trail := range p.trail {
			if i > 0 && trail.alpha > 0.05 {
				alpha := uint8(120 * trail.alpha * p.life * p.energy)
				trailColor := color.RGBA{148, 0, 211, alpha}
				
				// Gradient trail effect
				if i < len(p.trail)/2 {
					trailColor = color.RGBA{255, 100, 255, alpha}
				}
				
				vector.DrawFilledCircle(screen, float32(trail.x), float32(trail.y), 
					float32(trail.size), trailColor, false)
			}
		}
		
		// Main particle with multi-layer glow
		alpha := uint8(240 * p.life * p.energy)
		
		// Outer aura
		auraSize := p.size * 3.5
		auraAlpha := uint8(alpha / 4)
		auraColor := color.RGBA{200, 100, 255, auraAlpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(auraSize), auraColor, false)
		
		// Middle glow
		glowSize := p.size * 2.2
		glowAlpha := uint8(alpha / 2)
		glowColor := color.RGBA{220, 120, 255, glowAlpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(glowSize), glowColor, false)
		
		// Main particle
		mainColor := color.RGBA{148, 0, 211, alpha}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
			float32(p.size), mainColor, false)
		
		// Bright core
		if p.energy > 0.7 {
			coreAlpha := uint8(float64(alpha) * 0.8*255)
			coreColor := color.RGBA{255, 255, 255, coreAlpha}
			vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), 
				float32(p.size*0.4), coreColor, false)
		}
	}
}

func (ui *UIPage) drawFloatingElements(screen *ebiten.Image) {
	for _, fe := range ui.floatingElements {
		alpha := uint8(255 * fe.alpha)
		elementColor := color.RGBA{150, 200, 255, alpha}
		
		size := 15.0 * fe.scale
		
		switch fe.elementType {
		case 0: // Kanji-style symbols
			ui.drawKanjiSymbol(screen, fe.x, fe.y, size, fe.rotation, elementColor)
		case 1: // Mystical symbols
			ui.drawMysticalSymbol(screen, fe.x, fe.y, size, fe.rotation, elementColor)
		case 2: // Geometric patterns
			ui.drawGeometricPattern(screen, fe.x, fe.y, size, fe.rotation, elementColor)
		}
	}
}

func (ui *UIPage) drawCharacterArea(screen *ebiten.Image) {
	charAreaX := float32(screenWidth - 320)
	charAreaY := float32(40)
	charAreaWidth := float32(270)
	charAreaHeight := float32(270)
	
	// Enhanced frame with multiple layers
	frameColor1 := color.RGBA{100, 150, 255, 120}
	frameColor2 := color.RGBA{255, 215, 0, 80}
	
	// Outer frame
	vector.StrokeRect(screen, charAreaX-5, charAreaY-5, 
		charAreaWidth+10, charAreaHeight+10, 3, frameColor1, false)
	
	// Inner frame
	vector.StrokeRect(screen, charAreaX, charAreaY, 
		charAreaWidth, charAreaHeight, 2, frameColor2, false)
	
	// Corner tech elements
	cornerSize := float32(30)
	accentColor := color.RGBA{255, 215, 0, uint8(255 * ui.glowIntensity)}
	
	corners := []struct{ x, y float32 }{
		{charAreaX-8, charAreaY-8},
		{charAreaX+charAreaWidth-cornerSize+8, charAreaY-8},
		{charAreaX-8, charAreaY+charAreaHeight-cornerSize+8},
		{charAreaX+charAreaWidth-cornerSize+8, charAreaY+charAreaHeight-cornerSize+8},
	}
	
	for _, corner := range corners {
		// L-shaped corner brackets
		vector.DrawFilledRect(screen, corner.x, corner.y, cornerSize, 4, accentColor, false)
		vector.DrawFilledRect(screen, corner.x, corner.y, 4, cornerSize, accentColor, false)
	}
	
	// Enhanced placeholder with better styling
	placeholderText := "CHARACTER"
	textX := int(charAreaX + charAreaWidth/2 - float32(len(placeholderText)*5))
	textY := int(charAreaY + charAreaHeight/2 - 10)
	ui.drawGlowText(screen, placeholderText, textX, textY, 
		color.RGBA{200, 220, 255, 200}, 2.0)
	
	subText := "DISPLAY AREA"
	subTextX := int(charAreaX + charAreaWidth/2 - float32(len(subText)*4))
	subTextY := textY + 25
	ui.drawGlowText(screen, subText, subTextX, subTextY, 
		color.RGBA{150, 170, 200, 150}, 1.5)
}

func (ui *UIPage) drawEnhancedTitle(screen *ebiten.Image) {
	title := "呪術廻戦 × PAC-MAN"
	subtitle := "JUJUTSU KAISEN EDITION"
	
	titleY := 90
	titleX := screenWidth/2 - len(title)*18
	
	// Enhanced title background with multiple layers
	panelWidth := float32(len(title)*36 + 80)
	panelHeight := float32(100)
	panelX := float32(titleX - 40)
	panelY := float32(titleY - 25)
	
	// Background blur effect
	blurColor := color.RGBA{5, 10, 25, 150}
	vector.DrawFilledRect(screen, panelX-10, panelY-10, 
		panelWidth+20, panelHeight+20, blurColor, false)
	
	// Main panel with glassmorphism
	panelColor := color.RGBA{15, 25, 50, 140}
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight, panelColor, false)
	
	// Animated border with pulse effect
	borderIntensity := ui.glowIntensity
	borderColor := color.RGBA{100, 150, 255, uint8(180 * borderIntensity)}
	vector.StrokeRect(screen, panelX, panelY, panelWidth, panelHeight, 3, borderColor, false)
	
	// Secondary border
	innerBorderColor := color.RGBA{255, 215, 0, uint8(120 * borderIntensity)}
	vector.StrokeRect(screen, panelX+3, panelY+3, panelWidth-6, panelHeight-6, 1, innerBorderColor, false)
	
	// Title with enhanced effects
	ui.drawGlowText(screen, title, titleX, titleY, 
		color.RGBA{255, 215, 0, 255}, 5.0)
	
	// Subtitle with color animation
	subtitleY := titleY + 50
	subtitleX := screenWidth/2 - len(subtitle)*7
	subtitleHue := math.Sin(ui.animationTime*2) * 0.3 + 0.7
	subtitleColor := color.RGBA{
		uint8(150 * subtitleHue), 
		uint8(200 * subtitleHue), 
		255, 
		255,
	}
	ui.drawGlowText(screen, subtitle, subtitleX, subtitleY, subtitleColor, 2.5)
	
	// Logo area with enhanced effects
	ui.drawEnhancedLogoArea(screen)
}

func (ui *UIPage) drawEnhancedLogoArea(screen *ebiten.Image) {
	logoAreaX := float32(40)
	logoAreaY := float32(40)
	logoSize := float32(100)
	
	// Enhanced logo frame
	frameColor := color.RGBA{255, 215, 0, uint8(200 * ui.glowIntensity)}
	vector.StrokeRect(screen, logoAreaX-5, logoAreaY-5, 
		logoSize+10, logoSize+10, 3, frameColor, false)
	
	if ui.logoImage != nil {
		logoBounds := ui.logoImage.Bounds()
		logoWidth := float64(logoBounds.Dx())
		logoHeight := float64(logoBounds.Dy())
		
		scaleX := float64(logoSize) / logoWidth
		scaleY := float64(logoSize) / logoHeight
		scale := math.Min(scaleX, scaleY)
		
		scaledWidth := logoWidth * scale
		scaledHeight := logoHeight * scale
		offsetX := float64(logoAreaX) + (float64(logoSize)-scaledWidth)/2
		offsetY := float64(logoAreaY) + (float64(logoSize)-scaledHeight)/2
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(offsetX, offsetY)
		
		// Add glow effect to logo
		op.ColorM.Scale(1, 1, 1, float64(ui.glowIntensity))
		
		screen.DrawImage(ui.logoImage, op)
	} else {
		// Enhanced placeholder
		ui.drawGlowText(screen, "LOGO", int(logoAreaX+25), int(logoAreaY+50), 
			color.RGBA{255, 215, 0, 255}, 2.0)
	}
}

func (ui *UIPage) drawEnhancedMenu(screen *ebiten.Image) {
	menuStartY := 320
	menuSpacing := 90
	menuWidth := 450
	menuX := screenWidth/2 - menuWidth/2
	
	// Calculate enhanced menu panel
	panelHeight := float32(len(ui.menuOptions)*menuSpacing + 60)
	panelX := float32(menuX-50)
	panelY := float32(menuStartY-30)
	panelW := float32(menuWidth+100)
	
	// Draw GIF background with FULL OPACITY
	if len(ui.gifFrames) > 0 && ui.gifFrames[ui.frameIndex] != nil {
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
		
		// Create clipping mask for menu area
		tempImg := ebiten.NewImage(int(panelW), int(panelHeight))
		tempOp := &ebiten.DrawImageOptions{}
		tempOp.GeoM.Scale(scale, scale)
		tempOp.GeoM.Translate(offsetX-float64(panelX), offsetY-float64(panelY))
		
		// FULL OPACITY - No transparency applied to GIF
		tempOp.ColorM.Scale(1, 1, 1, 1.0) // Changed from 0.4 to 1.0 for full opacity
		tempOp.CompositeMode = ebiten.CompositeModeSourceOver // Better blending
		
		tempImg.DrawImage(currentFrame, tempOp)
		
		// Draw the GIF to screen
		finalOp := &ebiten.DrawImageOptions{}
		finalOp.GeoM.Translate(float64(panelX), float64(panelY))
		screen.DrawImage(tempImg, finalOp)
	}
	
	// Semi-transparent overlay for text readability (reduced opacity)
	overlayColor := color.RGBA{10, 15, 30, 60} // Reduced from 120 to 60
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelHeight, overlayColor, false)
	
	// Enhanced panel borders
	borderColor1 := color.RGBA{100, 150, 255, uint8(200 * ui.glowIntensity)}
	borderColor2 := color.RGBA{255, 215, 0, uint8(150 * ui.glowIntensity)}
	
	vector.StrokeRect(screen, panelX-2, panelY-2, panelW+4, panelHeight+4, 3, borderColor1, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelHeight, 2, borderColor2, false)
	
	// Menu options with enhanced effects
	for i, option := range ui.menuOptions {
		y := menuStartY + i*menuSpacing
		x := menuX
		
		if i == ui.selectedOption {
			// Enhanced selection with animation
			selectionWidth := float32(menuWidth + 40)
			selectionHeight := float32(70)
			selectionX := float32(x - 30)
			selectionY := float32(y - 20)
			
			// Pulsing selection background
			pulseIntensity := 0.7 + 0.3*math.Sin(ui.menuPulse*3)
			selectionAlpha := uint8(120 * pulseIntensity * ui.glowIntensity)
			selectionColor := color.RGBA{148, 0, 211, selectionAlpha}
			
			// Multi-layer selection effect
			vector.DrawFilledRect(screen, selectionX-5, selectionY-5, 
				selectionWidth+10, selectionHeight+10, 
				color.RGBA{255, 215, 0, selectionAlpha/3}, false)
			
			vector.DrawFilledRect(screen, selectionX, selectionY, 
				selectionWidth, selectionHeight, selectionColor, false)
			
			// Animated selection borders
			borderGlow := color.RGBA{255, 255, 255, uint8(255 * pulseIntensity)}
			vector.StrokeRect(screen, selectionX, selectionY, 
				selectionWidth, selectionHeight, 3, borderGlow, false)
			
			// Option text with maximum glow
			ui.drawGlowText(screen, option, x, y, 
				color.RGBA{255, 255, 255, 255}, 4.0)
			
			// Enhanced selection indicators
			indicatorX := x - 50
			indicatorY := float32(y + 8)
			ui.drawEnhancedArrow(screen, float32(indicatorX), indicatorY, 
				color.RGBA{255, 215, 0, 255})
			
			// Side energy effects
			ui.drawSelectionEnergyEffects(screen, selectionX, selectionY, 
				selectionWidth, selectionHeight)
			
			// Transition effect
			if ui.selectionTransition > 0 {
				ui.drawSelectionTransition(screen, selectionX, selectionY, 
					selectionWidth, selectionHeight)
			}
		} else {
			// Unselected options with hover effect
			hoverIntensity := 0.8 + 0.2*math.Sin(ui.animationTime*1.5+float64(i)*0.5)
			textColor := color.RGBA{
				uint8(180 * hoverIntensity), 
				uint8(190 * hoverIntensity), 
				uint8(230 * hoverIntensity), 
				220,
			}
			ui.drawGlowText(screen, option, x, y, textColor, 1.5)
		}
	}
}

func (ui *UIPage) drawSelectionEnergyEffects(screen *ebiten.Image, x, y, w, h float32) {
	// Side energy streams
	streamCount := 8
	for i := 0; i < streamCount; i++ {
		streamY := y + h*float32(i)/float32(streamCount)
		
		// Left side
		leftStartX := x - 10
		leftEndX := x - 30 - float32(15*math.Sin(ui.animationTime*4+float64(i)*0.5))
		
		streamAlpha := uint8(150 * ui.glowIntensity * 
			(0.6 + 0.4*math.Sin(ui.animationTime*3+float64(i)*0.3)))
		streamColor := color.RGBA{255, 215, 0, streamAlpha}
		
		vector.StrokeLine(screen, leftStartX, streamY, leftEndX, streamY, 2, streamColor, false)
		
		// Right side
		rightStartX := x + w + 10
		rightEndX := x + w + 30 + float32(15*math.Sin(ui.animationTime*4+float64(i)*0.5))
		
		vector.StrokeLine(screen, rightStartX, streamY, rightEndX, streamY, 2, streamColor, false)
	}
}

func (ui *UIPage) drawSelectionTransition(screen *ebiten.Image, x, y, w, h float32) {
	// Explosion effect on selection change
	transitionAlpha := uint8(255 * ui.selectionTransition)
	
	// Expanding rings
	for ring := 0; ring < 3; ring++ {
		ringRadius := float32(ui.selectionTransition * 100 * float64(ring+1))
		ringAlpha := transitionAlpha / uint8(ring+1)
		ringColor := color.RGBA{255, 255, 255, ringAlpha}
		
		centerX := x + w/2
		centerY := y + h/2
		
		vector.StrokeCircle(screen, centerX, centerY, ringRadius, 3, ringColor, false)
	}
}

func (ui *UIPage) drawEnhancedPacman(screen *ebiten.Image) {
	pacmanY := float32(screenHeight/2 - 60)
	pacmanSize := float32(70)
	
	// Enhanced multi-layer aura
	auraLayers := 5
	for layer := 0; layer < auraLayers; layer++ {
		layerSize := pacmanSize + float32(layer*18) + float32(25*ui.glowIntensity)
		layerAlpha := uint8(float64(80) / float64(layer + 1) * ui.glowIntensity)
		
		// Gradient aura colors
		if layer < 2 {
			auraColor := color.RGBA{200, 100, 255, layerAlpha}
			vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, 
				layerSize, auraColor, false)
		} else {
			auraColor := color.RGBA{148, 0, 211, layerAlpha/2}
			vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, 
				layerSize, auraColor, false)
		}
	}
	
	// Main Pacman with enhanced shading
	pacmanColor := color.RGBA{255, 215, 0, 255}
	vector.DrawFilledCircle(screen, float32(ui.pacmanX), pacmanY, pacmanSize, pacmanColor, false)
	
	// Multiple highlight layers
	highlightColor1 := color.RGBA{255, 255, 220, 200}
	highlightColor2 := color.RGBA{255, 255, 255, 120}
	
	vector.DrawFilledCircle(screen, float32(ui.pacmanX-12), pacmanY-12, 
		pacmanSize*0.35, highlightColor1, false)
	vector.DrawFilledCircle(screen, float32(ui.pacmanX-18), pacmanY-18, 
		pacmanSize*0.2, highlightColor2, false)
	
	// Enhanced mouth animation
	mouthAngle := ui.pacmanMouthAngle
	if mouthAngle > 0 {
		bgColor := color.RGBA{5, 8, 20, 255}
		
		mouthWidth := pacmanSize * 0.9
		mouthHeight := float32(float64(pacmanSize) * math.Sin(mouthAngle))
		
		// Create mouth shadow
		shadowColor := color.RGBA{0, 0, 0, 100}
		vector.DrawFilledRect(screen, float32(ui.pacmanX+2), pacmanY-mouthHeight/2+2, 
			mouthWidth, mouthHeight, shadowColor, false)
		
		// Main mouth
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY-mouthHeight/2, 
			mouthWidth, mouthHeight/2, bgColor, false)
		vector.DrawFilledRect(screen, float32(ui.pacmanX), pacmanY, 
			mouthWidth, mouthHeight/2, bgColor, false)
	}
	
	// Enhanced surrounding energy with orbital patterns
	orbitalLayers := 3
	for layer := 0; layer < orbitalLayers; layer++ {
		particlesInLayer := 6 + layer*2
		layerRadius := 90 + float64(layer)*30
		
		for i := 0; i < particlesInLayer; i++ {
			angle := float64(i)*2*math.Pi/float64(particlesInLayer) + 
				ui.animationTime*float64(2-layer) + float64(layer)*0.5
			distance := layerRadius + 25*math.Sin(ui.animationTime*3+float64(i)+float64(layer))
			
			px := ui.pacmanX + math.Cos(angle)*distance
			py := float64(pacmanY) + math.Sin(angle)*distance
			
			particleSize := float32(4 + 2*math.Sin(ui.animationTime*5+float64(i)+float64(layer)))
			particleAlpha := uint8(120 + 80*math.Sin(ui.animationTime*2.5+float64(i)))
			
			// Layer-specific colors
			var particleColor color.RGBA
			switch layer {
			case 0:
				particleColor = color.RGBA{255, 215, 0, particleAlpha}
			case 1:
				particleColor = color.RGBA{255, 100, 255, particleAlpha}
			case 2:
				particleColor = color.RGBA{148, 0, 211, particleAlpha}
			}
			
			vector.DrawFilledCircle(screen, float32(px), float32(py), 
				particleSize, particleColor, false)
		}
	}
}

func (ui *UIPage) drawEnhancedUIElements(screen *ebiten.Image) {
	ui.drawEnhancedSidePanels(screen)
	ui.drawAdvancedTechDecorations(screen)
	ui.drawEnhancedStatusIndicators(screen)
	ui.drawEnergyReadouts(screen)
}

func (ui *UIPage) drawEnhancedSidePanels(screen *ebiten.Image) {
	panelWidth := float32(140)
	panelHeight := float32(screenHeight - 80)
	
	// Left panel with gradient
	for i := 0; i < int(panelHeight); i++ {
		progress := float64(i) / float64(panelHeight)
		alpha := uint8(60 + 40*progress)
		rowColor := color.RGBA{8, 15, 35, alpha}
		vector.DrawFilledRect(screen, 20, float32(40+i), panelWidth, 1, rowColor, false)
	}
	
	leftBorderColor := color.RGBA{100, 150, 255, uint8(120 * ui.glowIntensity)}
	vector.StrokeRect(screen, 20, 40, panelWidth, panelHeight, 2, leftBorderColor, false)
	
	// Right panel
	rightX := float32(screenWidth - 160)
	for i := 0; i < int(panelHeight); i++ {
		progress := float64(i) / float64(panelHeight)
		alpha := uint8(60 + 40*progress)
		rowColor := color.RGBA{8, 15, 35, alpha}
		vector.DrawFilledRect(screen, rightX, float32(40+i), panelWidth, 1, rowColor, false)
	}
	
	rightBorderColor := color.RGBA{100, 150, 255, uint8(120 * ui.glowIntensity)}
	vector.StrokeRect(screen, rightX, 40, panelWidth, panelHeight, 2, rightBorderColor, false)
}

func (ui *UIPage) drawAdvancedTechDecorations(screen *ebiten.Image) {
	// Enhanced corner elements
	corners := []struct{ x, y float32 }{
		{50, 70}, {screenWidth - 190, 70},
		{50, screenHeight - 140}, {screenWidth - 190, screenHeight - 140},
	}
	
	for _, corner := range corners {
		intensity := ui.glowIntensity
		techColor := color.RGBA{100, 150, 255, uint8(150 * intensity)}
		accentColor := color.RGBA{255, 215, 0, uint8(180 * intensity)}
		
		// Multi-layer corner brackets
		bracketSize := float32(35)
		
		// Outer bracket
		vector.DrawFilledRect(screen, corner.x-2, corner.y-2, bracketSize+4, 3, techColor, false)
		vector.DrawFilledRect(screen, corner.x-2, corner.y-2, 3, bracketSize+4, techColor, false)
		
		// Inner bracket
		vector.DrawFilledRect(screen, corner.x, corner.y, bracketSize, 2, accentColor, false)
		vector.DrawFilledRect(screen, corner.x, corner.y, 2, bracketSize, accentColor, false)
		
		// Corner dots
		vector.DrawFilledCircle(screen, corner.x+bracketSize+8, corner.y+8, 3, accentColor, false)
		vector.DrawFilledCircle(screen, corner.x+8, corner.y+bracketSize+8, 3, techColor, false)
	}
	
	// Center HUD elements
	ui.drawCenterHUD(screen)
}

func (ui *UIPage) drawCenterHUD(screen *ebiten.Image) {
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)
	
	// Rotating HUD ring
	ringRadius := float32(250)
	ringSegments := 24
	
	for i := 0; i < ringSegments; i++ {
		angle := float64(i)*2*math.Pi/float64(ringSegments) + ui.animationTime*0.5
		
		x1 := centerX + float32(math.Cos(angle)*(float64(ringRadius)-10))
		y1 := centerY + float32(math.Sin(angle)*(float64(ringRadius)-10))
		x2 := centerX + float32(math.Cos(angle)*(float64(ringRadius)+10))
		y2 := centerY + float32(math.Sin(angle)*(float64(ringRadius)+10))
		
		segmentAlpha := uint8(80 + 40*math.Sin(ui.animationTime*2+float64(i)*0.3))
		segmentColor := color.RGBA{100, 150, 255, segmentAlpha}
		
		if i%4 == 0 {
			segmentColor = color.RGBA{255, 215, 0, segmentAlpha}
		}
		
		vector.StrokeLine(screen, x1, y1, x2, y2, 2, segmentColor, false)
	}
}

func (ui *UIPage) drawEnhancedStatusIndicators(screen *ebiten.Image) {
	// Enhanced status bar
	statusY := float32(screenHeight - 80)
	statusHeight := float32(50)
	
	// Gradient status background
	for i := 0; i < int(statusHeight); i++ {
		progress := float64(i) / float64(statusHeight)
		alpha := uint8(100 + 50*progress)
		rowColor := color.RGBA{15, 25, 45, alpha}
		vector.DrawFilledRect(screen, 40, statusY+float32(i), screenWidth-80, 1, rowColor, false)
	}
	
	// Status border with glow
	borderColor := color.RGBA{100, 150, 255, uint8(150 * ui.glowIntensity)}
	vector.StrokeRect(screen, 40, statusY, screenWidth-80, statusHeight, 2, borderColor, false)
	
	// Enhanced status indicators
	for i := 0; i < 8; i++ {
		dotX := float32(70 + i*40)
		dotY := statusY + statusHeight/2
		
		baseSize := float32(5)
		pulseSize := baseSize + float32(4*math.Sin(ui.animationTime*4+float64(i)*0.4))
		
		dotAlpha := uint8(120 + 135*math.Sin(ui.animationTime*2.5+float64(i)*0.4))
		
		// Alternate colors
		var dotColor color.RGBA
		if i%2 == 0 {
			dotColor = color.RGBA{148, 0, 211, dotAlpha}
		} else {
			dotColor = color.RGBA{255, 215, 0, dotAlpha}
		}
		
		// Outer glow
		vector.DrawFilledCircle(screen, dotX, dotY, pulseSize*1.8, 
			color.RGBA{dotColor.R, dotColor.G, dotColor.B, dotAlpha/3}, false)
		
		// Main dot
		vector.DrawFilledCircle(screen, dotX, dotY, pulseSize, dotColor, false)
	}
}

func (ui *UIPage) drawEnergyReadouts(screen *ebiten.Image) {
	// Left side energy readout
	readoutX := float32(60)
	readoutY := float32(200)
	
	ui.drawGlowText(screen, "CURSED ENERGY", int(readoutX), int(readoutY), 
		color.RGBA{148, 0, 211, 200}, 1.5)
	
	// Energy bar
	barWidth := float32(80)
	barHeight := float32(8)
	energyLevel := 0.7 + 0.3*math.Sin(ui.animationTime*1.5)
	
	// Background bar
	vector.DrawFilledRect(screen, readoutX, readoutY+20, barWidth, barHeight, 
		color.RGBA{50, 50, 50, 150}, false)
	
	// Energy fill
	fillWidth := barWidth * float32(energyLevel)
	energyColor := color.RGBA{148, 0, 211, 200}
	vector.DrawFilledRect(screen, readoutX, readoutY+20, fillWidth, barHeight, energyColor, false)
	
	// Energy bar glow
	glowColor := color.RGBA{200, 100, 255, uint8(100 * ui.glowIntensity)}
	vector.DrawFilledRect(screen, readoutX, readoutY+18, fillWidth, barHeight+4, glowColor, false)
	
	// Right side technique readout
	rightReadoutX := float32(screenWidth - 180)
	ui.drawGlowText(screen, "DOMAIN", int(rightReadoutX), int(readoutY), 
		color.RGBA{255, 215, 0, 200}, 1.5)
	ui.drawGlowText(screen, "EXPANSION", int(rightReadoutX), int(readoutY+20), 
		color.RGBA{255, 215, 0, 180}, 1.2)
}

func (ui *UIPage) drawEnhancedFooter(screen *ebiten.Image) {
	instructions := []string{
		"↑↓ / W S  Navigate Menu",
		"ENTER / SPACE  Select Option",
		"Experience Infinite Cursed Energy",
	}
	
	footerY := screenHeight - 50
	spacing := 280
	
	for i, instruction := range instructions {
		x := 60 + i*spacing
		y := footerY
		
		var instrColor color.RGBA
		switch i {
		case 0:
			instrColor = color.RGBA{150, 200, 255, 200}
		case 1:
			instrColor = color.RGBA{255, 215, 0, 200}
		case 2:
			pulse := 0.7 + 0.3*math.Sin(ui.animationTime*2)
			instrColor = color.RGBA{
				uint8(148 * pulse), 
				0, 
				uint8(211 * pulse), 
				200,
			}
		}
		
		ui.drawGlowText(screen, instruction, x, y, instrColor, 1.8)
	}
	
	// Footer decorative line
	lineY := float32(footerY - 15)
	lineColor := color.RGBA{100, 150, 255, uint8(100 * ui.glowIntensity)}
	vector.DrawFilledRect(screen, 50, lineY, screenWidth-100, 2, lineColor, false)
}

// Enhanced helper drawing functions
func (ui *UIPage) drawGlowText(screen *ebiten.Image, txt string, x, y int, clr color.RGBA, glowIntensity float64) {
	if glowIntensity > 1.0 {
		glowColor := color.RGBA{clr.R, clr.G, clr.B, uint8(float64(clr.A) * 0.4)}
		glowRadius := int(glowIntensity * 1.2)
		
		// Multiple glow layers for smoother effect
		for layer := 0; layer < 3; layer++ {
			layerRadius := glowRadius - layer*int(glowRadius/3)
			if layerRadius <= 0 { continue }
			
			for dx := -layerRadius; dx <= layerRadius; dx++ {
				for dy := -layerRadius; dy <= layerRadius; dy++ {
					if dx != 0 || dy != 0 {
						distance := math.Sqrt(float64(dx*dx + dy*dy))
						if distance <= float64(layerRadius) {
							alpha := uint8(float64(glowColor.A) * 
								(1.0 - distance/float64(layerRadius)) / float64(layer+1))
							fadeColor := color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}
							text.Draw(screen, txt, basicfont.Face7x13, x+dx, y+dy+13, fadeColor)
						}
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
	
	// Enhanced hexagon with fill and stroke
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		vector.StrokeLine(screen, points[i*2], points[i*2+1], 
			points[next*2], points[next*2+1], 2.5, clr, false)
	}
	
	// Add center dot
	centerColor := color.RGBA{clr.R, clr.G, clr.B, clr.A/2}
	vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), 
		float32(size*0.2), centerColor, false)
}

func (ui *UIPage) drawDiamond(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cosR := float32(math.Cos(rotation))
	sinR := float32(math.Sin(rotation))
	
	vertices := []struct{ x, y float32 }{
		{0, -halfSize}, {halfSize, 0}, {0, halfSize}, {-halfSize, 0},
	}
	
	// Draw diamond with enhanced thickness
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		
		x1 := vertices[i].x*cosR - vertices[i].y*sinR + float32(centerX)
		y1 := vertices[i].x*sinR + vertices[i].y*cosR + float32(centerY)
		x2 := vertices[next].x*cosR - vertices[next].y*sinR + float32(centerX)
		y2 := vertices[next].x*sinR + vertices[next].y*cosR + float32(centerY)
		
		vector.StrokeLine(screen, x1, y1, x2, y2, 2, clr, false)
	}
}

func (ui *UIPage) drawCross(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	halfSize := float32(size / 2)
	cosR := float32(math.Cos(rotation))
	sinR := float32(math.Sin(rotation))
	
	cx := float32(centerX)
	cy := float32(centerY)
	
	// Enhanced cross with thickness variation
	thickness := float32(2 + math.Sin(ui.animationTime*3)*0.5)
	
	// Horizontal line
	x1 := -halfSize*cosR + cx
	y1 := -halfSize*sinR + cy
	x2 := halfSize*cosR + cx
	y2 := halfSize*sinR + cy
	vector.StrokeLine(screen, x1, y1, x2, y2, thickness, clr, false)
	
	// Vertical line
	x3 := halfSize*sinR + cx
	y3 := -halfSize*cosR + cy
	x4 := -halfSize*sinR + cx
	y4 := halfSize*cosR + cy
	vector.StrokeLine(screen, x3, y3, x4, y4, thickness, clr, false)
}

func (ui *UIPage) drawStar(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	points := 5
	outerRadius := float32(size)
	innerRadius := outerRadius * 0.4
	
	for i := 0; i < points*2; i++ {
		angle := rotation + float64(i)*math.Pi/float64(points)
		var radius float32
		if i%2 == 0 {
			radius = outerRadius
		} else {
			radius = innerRadius
		}
		
		x := float32(centerX) + radius*float32(math.Cos(angle))
		y := float32(centerY) + radius*float32(math.Sin(angle))
		
		if i > 0 {
			prevAngle := rotation + float64(i-1)*math.Pi/float64(points)
			var prevRadius float32
			if (i-1)%2 == 0 {
				prevRadius = outerRadius
			} else {
				prevRadius = innerRadius
			}
			
			prevX := float32(centerX) + prevRadius*float32(math.Cos(prevAngle))
			prevY := float32(centerY) + prevRadius*float32(math.Sin(prevAngle))
			
			vector.StrokeLine(screen, prevX, prevY, x, y, 1.5, clr, false)
		}
	}
}

func (ui *UIPage) drawKanjiSymbol(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	// Simplified kanji-like symbol
	halfSize := float32(size / 2)
	cx := float32(centerX)
	cy := float32(centerY)
	
	// Horizontal strokes
	vector.StrokeLine(screen, cx-halfSize, cy-halfSize/2, cx+halfSize, cy-halfSize/2, 2, clr, false)
	vector.StrokeLine(screen, cx-halfSize*0.7, cy, cx+halfSize*0.7, cy, 2, clr, false)
	vector.StrokeLine(screen, cx-halfSize*0.5, cy+halfSize/2, cx+halfSize*0.5, cy+halfSize/2, 2, clr, false)
	
	// Vertical stroke
	vector.StrokeLine(screen, cx, cy-halfSize, cx, cy+halfSize, 2, clr, false)
}

func (ui *UIPage) drawMysticalSymbol(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	// Pentagram-like mystical symbol
	points := 5
	radius := float32(size * 0.8)
	cx := float32(centerX)
	cy := float32(centerY)
	
	// Draw connecting lines between every second point
	for i := 0; i < points; i++ {
		angle1 := rotation + float64(i)*2*math.Pi/float64(points)
		angle2 := rotation + float64((i+2)%points)*2*math.Pi/float64(points)
		
		x1 := cx + radius*float32(math.Cos(angle1))
		y1 := cy + radius*float32(math.Sin(angle1))
		x2 := cx + radius*float32(math.Cos(angle2))
		y2 := cy + radius*float32(math.Sin(angle2))
		
		vector.StrokeLine(screen, x1, y1, x2, y2, 1.5, clr, false)
	}
}

func (ui *UIPage) drawGeometricPattern(screen *ebiten.Image, centerX, centerY, size, rotation float64, clr color.RGBA) {
	// Geometric mandala pattern
	cx := float32(centerX)
	cy := float32(centerY)
	radius := float32(size)
	
	// Inner circle
	vector.StrokeCircle(screen, cx, cy, radius*0.3, 1.5, clr, false)
	
	// Radiating lines
	for i := 0; i < 8; i++ {
		angle := rotation + float64(i)*math.Pi/4
		x1 := cx + radius*0.4*float32(math.Cos(angle))
		y1 := cy + radius*0.4*float32(math.Sin(angle))
		x2 := cx + radius*float32(math.Cos(angle))
		y2 := cy + radius*float32(math.Sin(angle))
		
		vector.StrokeLine(screen, x1, y1, x2, y2, 1, clr, false)
	}
}

func (ui *UIPage) drawEnhancedArrow(screen *ebiten.Image, x, y float32, clr color.RGBA) {
	// Modern arrow with glow effect
	arrowSize := float32(16)
	
	// Glow effect
	glowAlpha := uint8(100 * ui.glowIntensity)
	//glowColor := color.RGBA{clr.R, clr.G, clr.B, glowAlpha}
	
	for offset := 3; offset >= 0; offset-- {
		alpha := glowAlpha / uint8(offset+1)
		layerColor := color.RGBA{clr.R, clr.G, clr.B, alpha}
		layerSize := arrowSize + float32(offset)*2
		
		// Triangle points
		x1 := x - float32(offset)
		y1 := y - layerSize/2
		x2 := x - float32(offset)
		y2 := y + layerSize/2
		x3 := x + layerSize - float32(offset)
		y3 := y
		
		// Draw triangle outline
		vector.StrokeLine(screen, x1, y1, x2, y2, 2, layerColor, false)
		vector.StrokeLine(screen, x1, y1, x3, y3, 2, layerColor, false)
		vector.StrokeLine(screen, x2, y2, x3, y3, 2, layerColor, false)
		
		// Fill triangle
		for i := 0; i < int(layerSize); i++ {
			progress := float32(i) / layerSize
			startY := y1 + (y2-y1)*progress
			endX := x + layerSize*progress - float32(offset)
			vector.StrokeLine(screen, x-float32(offset), startY, endX, y, 1, layerColor, false)
		}
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

// SetImages allows setting the images/gifs for the UI with enhanced handling
func (ui *UIPage) SetImages(logo *ebiten.Image, gifFrames []*ebiten.Image, bg *ebiten.Image) {
	ui.logoImage = logo
	ui.gifFrames = gifFrames
	ui.bgImage = bg
	
	// Reset GIF animation when new frames are set
	if len(gifFrames) > 0 {
		ui.frameIndex = 0
		ui.frameTicker = 0
	}
}

// Additional utility functions for enhanced UI

// GetMenuPulse returns the current menu pulse value for external use
func (ui *UIPage) GetMenuPulse() float64 {
	return ui.menuPulse
}

// GetGlowIntensity returns current glow intensity
func (ui *UIPage) GetGlowIntensity() float64 {
	return ui.glowIntensity
}

// SetScreenShake allows external triggers for screen shake
func (ui *UIPage) SetScreenShake(intensity float64) {
	ui.screenShake = math.Max(ui.screenShake, intensity)
}

// GetAnimationTime returns current animation time for synchronization
func (ui *UIPage) GetAnimationTime() float64 {
	return ui.animationTime
}

// ResetSelectionTransition resets the selection transition effect
func (ui *UIPage) ResetSelectionTransition() {
	ui.selectionTransition = 0
}

// GetCurrentGIFFrame returns the current GIF frame index
func (ui *UIPage) GetCurrentGIFFrame() int {
	return ui.frameIndex
}

// SetGIFSpeed allows adjusting GIF animation speed
func (ui *UIPage) SetGIFSpeed(delay int) {
	if delay > 0 {
		ui.frameDelay = delay
	}
}
