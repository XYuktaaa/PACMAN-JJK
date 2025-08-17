// intro.go
package main

import (
    "fmt"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "image/color"
    "math"
    //"time"
)

type IntroState int

const (
    IntroTitle IntroState = iota
    IntroCharacterReveal
    IntroReadyScreen
    IntroComplete
)

type IntroSystem struct {
    State       IntroState
    Timer       int
    Alpha       float64
    TextAlpha   float64
    
    // Title screen
    TitleImage      *ebiten.Image
    LogoImage       *ebiten.Image
    BackgroundImage *ebiten.Image
    
    // Character reveal
    CharacterImages []*ebiten.Image
    CharacterNames  []string
    CurrentChar     int
    CharRevealTimer int
    
    // Ready screen
    ReadyTimer      int
    ReadyPulse      float64
    
    // Fonts
    TitleFont       font.Face
    SubtitleFont    font.Face
    CharNameFont    font.Face
    
    // Effects
    Particles       []Particle
    FlashEffect     float64
    
    // Sound
    SoundManager    *SoundManager
}

type Particle struct {
    X, Y        float64
    VX, VY      float64
    Life        int
    MaxLife     int
    Size        float64
    Color       color.RGBA
    ParticleType string
}

func NewIntroSystem() *IntroSystem {
    intro := &IntroSystem{
        State:           IntroTitle,
        Timer:          0,
        Alpha:          0.0,
        TextAlpha:      0.0,
        CharacterNames: []string{"Gojo Satoru", "Sukuna", "Kenjaku", "Mahito"},
        CurrentChar:    0,
        CharRevealTimer: 0,
        ReadyTimer:     0,
        ReadyPulse:     0.0,
        SoundManager:   NewSoundManager(),
    }
    
    intro.loadAssets()
    intro.initParticles()
    
    return intro
}

func (i *IntroSystem) loadAssets() {
    var err error
    
    // Load images (create placeholder if files don't exist)
    i.TitleImage = loadImageWithFallback("assets/jjk_title.png", 400, 200, color.RGBA{255, 100, 100, 255})
    i.LogoImage = loadImageWithFallback("assets/jjk_logo.png", 300, 150, color.RGBA{200, 50, 50, 255})
    i.BackgroundImage = loadImageWithFallback("assets/intro_bg.png", 1200, 800, color.RGBA{20, 20, 40, 255})
    
    // Load character images
    charPaths := []string{
        "assets/gojo_intro.png",
        "assets/sukuna_intro.png", 
        "assets/kenjaku_intro.png",
        "assets/mahito_intro.png",
    }
    
    for _, path := range charPaths {
        img := loadImageWithFallback(path, 300, 400, color.RGBA{100, 100, 200, 255})
        i.CharacterImages = append(i.CharacterImages, img)
    }
    
    // Create fonts
    if PressStartFont != nil {
        i.TitleFont, err = opentype.NewFace(PressStartFont, &opentype.FaceOptions{
            Size:    36,
            DPI:     72,
            Hinting: font.HintingFull,
        })
        if err == nil {
            i.SubtitleFont, err = opentype.NewFace(PressStartFont, &opentype.FaceOptions{
                Size:    18,
                DPI:     72,
                Hinting: font.HintingFull,
            })
        }
        if err == nil {
            i.CharNameFont, err = opentype.NewFace(PressStartFont, &opentype.FaceOptions{
                Size:    24,
                DPI:     72,
                Hinting: font.HintingFull,
            })
        }
    }
}

func loadImageWithFallback(path string, width, height int, fallbackColor color.RGBA) *ebiten.Image {
    img, _, err := ebitenutil.NewImageFromFile(path)
    if err != nil {
        // Create fallback image
        fallback := ebiten.NewImage(width, height)
        fallback.Fill(fallbackColor)
        return fallback
    }
    return img
}

func (i *IntroSystem) initParticles() {
    // Create cursed energy particles
    for j := 0; j < 50; j++ {
        particle := Particle{
            X:           float64(j%10) * 120,
            Y:           float64(j/10) * 80,
            VX:          (float64(j%7) - 3) * 0.5,
            VY:          (float64(j%5) - 2) * 0.3,
            Life:        60 + j%120,
            MaxLife:     60 + j%120,
            Size:        2 + float64(j%3),
            Color:       color.RGBA{100 + uint8(j%155), 50, 150 + uint8(j%105), 200},
            ParticleType: "cursed_energy",
        }
        i.Particles = append(i.Particles, particle)
    }
}

func (i *IntroSystem) Update() error {
    i.Timer++
    
    // Update particles
    for idx := range i.Particles {
        p := &i.Particles[idx]
        p.X += p.VX
        p.Y += p.VY
        p.Life--
        
        // Wrap around screen
        if p.X < 0 { p.X = 1200 }
        if p.X > 1200 { p.X = 0 }
        if p.Y < 0 { p.Y = 800 }
        if p.Y > 800 { p.Y = 0 }
        
        // Reset particle
        if p.Life <= 0 {
            p.Life = p.MaxLife
            p.X = float64(idx%10) * 120
            p.Y = float64(idx/10) * 80
        }
    }
    
    switch i.State {
    case IntroTitle:
        return i.updateTitleScreen()
    case IntroCharacterReveal:
        return i.updateCharacterReveal()
    case IntroReadyScreen:
        return i.updateReadyScreen()
    }
    
    return nil
}

func (i *IntroSystem) updateTitleScreen() error {
    // Fade in effect
    if i.Timer < 120 {
        i.Alpha = float64(i.Timer) / 120.0
    } else {
        i.Alpha = 1.0
    }
    
    // Text fade in after logo
    if i.Timer > 60 {
        i.TextAlpha = math.Min(1.0, float64(i.Timer-60)/60.0)
    }
    
    // Flash effect
    i.FlashEffect = math.Max(0, i.FlashEffect-0.05)
    
    // Skip to character reveal on input or after delay
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) || 
       inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
       i.Timer > 300 {
        
        i.State = IntroCharacterReveal
        i.Timer = 0
        i.CurrentChar = 0
        i.FlashEffect = 1.0
        
        // Play sound effect
        i.SoundManager.PlaySFX("transition")
        
        fmt.Println("Transitioning to character reveal")
    }
    
    return nil
}

func (i *IntroSystem) updateCharacterReveal() error {
    i.CharRevealTimer++
    
    // Character reveal timing
    charDisplayTime := 180 // 3 seconds per character at 60 FPS
    
    if i.CharRevealTimer >= charDisplayTime {
        i.CurrentChar++
        i.CharRevealTimer = 0
        i.FlashEffect = 0.8
        
        // Play character sound
        i.SoundManager.PlaySFX("character_reveal")
        
        if i.CurrentChar >= len(i.CharacterImages) {
            i.State = IntroReadyScreen
            i.Timer = 0
            i.ReadyTimer = 0
            fmt.Println("All characters revealed, moving to ready screen")
        }
    }
    
    // Skip on input
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) || 
       inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
        i.State = IntroReadyScreen
        i.Timer = 0
        i.ReadyTimer = 0
    }
    
    return nil
}

func (i *IntroSystem) updateReadyScreen() error {
    i.ReadyTimer++
    
    // Pulsing effect for "READY" text
    i.ReadyPulse = 0.8 + 0.2*math.Sin(float64(i.ReadyTimer)*0.2)
    
    // Auto-advance after delay or on input
    if i.ReadyTimer > 180 || 
       inpututil.IsKeyJustPressed(ebiten.KeySpace) || 
       inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
        
        i.State = IntroComplete
        fmt.Println("Intro complete, starting game")
        
        // Play game start sound
        i.SoundManager.PlaySFX("game_start")
    }
    
    return nil
}

func (i *IntroSystem) Draw(screen *ebiten.Image) {
    switch i.State {
    case IntroTitle:
        i.drawTitleScreen(screen)
    case IntroCharacterReveal:
        i.drawCharacterReveal(screen)
    case IntroReadyScreen:
        i.drawReadyScreen(screen)
    }
    
    // Draw particles on top
    i.drawParticles(screen)
    
    // Flash effect
    if i.FlashEffect > 0 {
        flashAlpha := uint8(i.FlashEffect * 100)
        ebitenutil.DrawRect(screen, 0, 0, 1200, 800, color.RGBA{255, 255, 255, flashAlpha})
    }
}

func (i *IntroSystem) drawTitleScreen(screen *ebiten.Image) {
    // Background
    op := &ebiten.DrawImageOptions{}
    op.ColorM.Scale(1, 1, 1, i.Alpha)
    screen.DrawImage(i.BackgroundImage, op)
    
    // Main logo
    op = &ebiten.DrawImageOptions{}
    op.GeoM.Translate(450, 200)
    op.ColorM.Scale(1, 1, 1, i.Alpha)
    screen.DrawImage(i.LogoImage, op)
    
    // Title text
    if i.TitleFont != nil && i.TextAlpha > 0 {
        titleColor := color.RGBA{255, 255, 255, uint8(i.TextAlpha * 255)}
        drawTextCentered(screen, "JUJUTSU KAISEN", i.TitleFont, 600, 450, titleColor)
        
        subtitleColor := color.RGBA{200, 200, 200, uint8(i.TextAlpha * 180)}
        drawTextCentered(screen, "PACMAN", i.SubtitleFont, 600, 490, subtitleColor)
        
        // Instructions
        if i.Timer > 180 {
            instrAlpha := math.Sin(float64(i.Timer-180) * 0.1) * 0.5 + 0.5
            instrColor := color.RGBA{255, 255, 100, uint8(instrAlpha * 200)}
            drawTextCentered(screen, "PRESS SPACE TO CONTINUE", i.SubtitleFont, 600, 650, instrColor)
        }
    }
}

func (i *IntroSystem) drawCharacterReveal(screen *ebiten.Image) {
    // Dark background
    ebitenutil.DrawRect(screen, 0, 0, 1200, 800, color.RGBA{10, 10, 20, 255})
    
    if i.CurrentChar < len(i.CharacterImages) {
        // Character image
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(450, 150)
        
        // Scaling effect
        scale := 0.8 + 0.2*math.Sin(float64(i.CharRevealTimer)*0.1)
        op.GeoM.Scale(scale, scale)
        
        screen.DrawImage(i.CharacterImages[i.CurrentChar], op)
        
        // Character name
        if i.CharNameFont != nil && i.CurrentChar < len(i.CharacterNames) {
            nameColor := color.RGBA{255, 200, 100, 255}
            drawTextCentered(screen, i.CharacterNames[i.CurrentChar], i.CharNameFont, 600, 600, nameColor)
        }
        
        // Progress indicator
        progressText := fmt.Sprintf("%d / %d", i.CurrentChar+1, len(i.CharacterImages))
        if i.SubtitleFont != nil {
            drawTextCentered(screen, progressText, i.SubtitleFont, 600, 700, color.RGBA{150, 150, 150, 255})
        }
    }
}

func (i *IntroSystem) drawReadyScreen(screen *ebiten.Image) {
    // Gradient background
    for y := 0; y < 800; y++ {
        alpha := uint8(float64(y) / 800.0 * 100)
        ebitenutil.DrawRect(screen, 0, float64(y), 1200, 1, color.RGBA{20, 0, 40, alpha})
    }
    
    // "READY?" text with pulsing effect
    if i.TitleFont != nil {
        readyAlpha := uint8(i.ReadyPulse * 255)
        readyColor := color.RGBA{255, 100, 100, readyAlpha}
        drawTextCentered(screen, "READY?", i.TitleFont, 600, 350, readyColor)
        
        // Subtitle
        subColor := color.RGBA{255, 255, 255, 200}
        drawTextCentered(screen, "Prepare to face the cursed spirits!", i.SubtitleFont, 600, 420, subColor)
    }
    
    // Countdown or instruction
    if i.ReadyTimer > 60 {
        instrText := "PRESS SPACE TO BEGIN"
        if i.ReadyTimer > 120 {
            countdown := 4 - ((i.ReadyTimer - 120) / 60)
            if countdown > 0 {
                instrText = fmt.Sprintf("STARTING IN %d...", countdown)
            }
        }
        
        if i.SubtitleFont != nil {
            instrAlpha := math.Sin(float64(i.ReadyTimer) * 0.15) * 0.3 + 0.7
            instrColor := color.RGBA{255, 255, 100, uint8(instrAlpha * 255)}
            drawTextCentered(screen, instrText, i.SubtitleFont, 600, 500, instrColor)
        }
    }
}

func (i *IntroSystem) drawParticles(screen *ebiten.Image) {
    for _, p := range i.Particles {
        alpha := float64(p.Life) / float64(p.MaxLife)
        particleColor := p.Color
        particleColor.A = uint8(alpha * float64(particleColor.A))
        
        ebitenutil.DrawRect(screen, p.X, p.Y, p.Size, p.Size, particleColor)
    }
}

// Helper function to draw centered text
func drawTextCentered(screen *ebiten.Image, text string, font font.Face, x, y float64, clr color.Color) {
    if font == nil {
        // Fallback to debug print if font is not available
        ebitenutil.DebugPrintAt(screen, text, int(x-float64(len(text)*6)), int(y))
        return
    }
    
    // Calculate text width for centering (approximate)
    textWidth := float64(len(text)) * 12 // Rough estimation
    drawX := x - textWidth/2
    
    // For now, use debug print with offset - in a real implementation you'd use text.Draw
    ebitenutil.DebugPrintAt(screen, text, int(drawX), int(y))
}

func (i *IntroSystem) IsComplete() bool {
    return i.State == IntroComplete
}

func (i *IntroSystem) Reset() {
    i.State = IntroTitle
    i.Timer = 0
    i.Alpha = 0.0
    i.TextAlpha = 0.0
    i.CurrentChar = 0
    i.CharRevealTimer = 0
    i.ReadyTimer = 0
    i.ReadyPulse = 0.0
    i.FlashEffect = 0.0
}

// Sound Manager for handling audio
// type SoundManager struct {
//     BGMEnabled bool
//     SFXEnabled bool
//     Volume     float64
    
//     // Sound channels/sources would go here
//     // For now, we'll just simulate with debug prints
// }

// func NewSoundManager() *SoundManager {
//     return &SoundManager{
//         BGMEnabled: true,
//         SFXEnabled: true,
//         Volume:     0.7,
//     }
// }

// func (s *SoundManager) PlayBGM(track string) {
//     if !s.BGMEnabled {
//         return
//     }
//     fmt.Printf("🎵 Playing BGM: %s (Volume: %.1f)\n", track, s.Volume)
    
//     // In a real implementation, you'd use a library like:
//     // - github.com/hajimehoshi/oto for low-level audio
//     // - github.com/faiface/beep for higher-level audio
//     // - Or integrate with Ebiten's audio capabilities
// }

// func (s *SoundManager) PlaySFX(effect string) {
//     if !s.SFXEnabled {
//         return
//     }
//     fmt.Printf("🔊 Playing SFX: %s\n", effect)
    
    // Sound effects you might want:
    // - "transition" - for screen transitions
    // - "character_reveal" - when showing each character
    // - "game_start" - when the game begins
    // - "pellet_eat" - when eating pellets
    // - "power_pellet" - when eating power pellets
    // - "ghost_eaten" - when eating a ghost
    // - "death" - when player dies
    // - "game_over" - game over sound
// }

// func (s *SoundManager) StopBGM() {
//     fmt.Println("🎵 Stopping BGM")
// }

// func (s *SoundManager) SetVolume(volume float64) {
//     s.Volume = math.Max(0.0, math.Min(1.0, volume))
//     fmt.Printf("🔊 Volume set to: %.1f\n", s.Volume)
// }

// func (s *SoundManager) ToggleBGM() {
//     s.BGMEnabled = !s.BGMEnabled
//     fmt.Printf("🎵 BGM %s\n", map[bool]string{true: "enabled", false: "disabled"}[s.BGMEnabled])
// }

// func (s *SoundManager) ToggleSFX() {
//     s.SFXEnabled = !s.SFXEnabled  
//     fmt.Printf("🔊 SFX %s\n", map[bool]string{true: "enabled", false: "disabled"}[s.SFXEnabled])
// }
