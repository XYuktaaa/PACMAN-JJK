// audio.go - Fixed audio system implementation
package main

import (
    "math"
    "bytes"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
    
    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/vorbis"
    "github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
    SampleRate   = 48000  // PipeWire default sample rate
    MaxVolume    = 1.0
    DefaultVolume = 0.7
)

type AudioSystem struct {
    AudioContext *audio.Context
    
    // Background Music
    BGMPlayer    *audio.Player
    CurrentBGM   string
    BGMVolume    float64
    BGMEnabled   bool
    BGMFadeOut   bool
    BGMFadeIn    bool
    BGMFadeTimer int
    
    // Sound Effects
    SFXEnabled bool
    SFXVolume  float64
    SFXPlayers map[string]*audio.Player
    
    // Audio data cache
    BGMData map[string][]byte
    SFXData map[string][]byte
    
    // Loading status
    AssetsLoaded bool
    LoadErrors   []string
}

var globalAudioCtx *audio.Context
var globalAudioSystem *AudioSystem

func NewAudioSystem() *AudioSystem {
    // Initialize audio context with error handling for PipeWire systems
    if globalAudioCtx == nil {
        fmt.Println("🎵 Creating new audio context for PipeWire...")
        
        // Try PipeWire-friendly sample rates
        sampleRates := []int{48000, 44100, 22050}
        var audioCtx *audio.Context
        
        for _, rate := range sampleRates {
            fmt.Printf("🎵 Trying sample rate: %d Hz\n", rate)
            audioCtx = audio.NewContext(rate)
            if audioCtx != nil {
                fmt.Printf("✅ Successfully created audio context at %d Hz\n", rate)
                globalAudioCtx = audioCtx
                break
            }
        }
        
        if globalAudioCtx == nil {
            fmt.Println("❌ Failed to create audio context with any sample rate")
            fmt.Println("💡 Try running with: SDL_AUDIODRIVER=pulse ./your-game")
            // Return a dummy audio system that won't crash
            return &AudioSystem{
                AudioContext: nil,
                BGMVolume:    DefaultVolume,
                SFXVolume:    DefaultVolume,
                BGMEnabled:   false, // Disable audio if context creation failed
                SFXEnabled:   false,
                SFXPlayers:   make(map[string]*audio.Player),
                BGMData:      make(map[string][]byte),
                SFXData:      make(map[string][]byte),
                AssetsLoaded: true, // Set to true to avoid loading attempts
                LoadErrors:   []string{"Audio context creation failed - try SDL_AUDIODRIVER=pulse"},
            }
        }
    }
    
    // Create audio system
    audioSystem := &AudioSystem{
        AudioContext: globalAudioCtx,
        BGMVolume:    DefaultVolume,
        SFXVolume:    DefaultVolume,
        BGMEnabled:   true,
        SFXEnabled:   true,
        SFXPlayers:   make(map[string]*audio.Player),
        BGMData:      make(map[string][]byte),
        SFXData:      make(map[string][]byte),
        AssetsLoaded: false,
        LoadErrors:   make([]string, 0),
    }
    
    fmt.Printf("🎵 AudioSystem created successfully (Sample Rate: %d Hz)\n", globalAudioCtx.SampleRate())
    return audioSystem
}

func (a *AudioSystem) LoadAllAudio() {
    // Safety check
    if a == nil {
        fmt.Println("❌ Error: AudioSystem is nil")
        return
    }
    
    // Check if audio context is available
    if a.AudioContext == nil {
        fmt.Println("⚠️  Audio context is nil, skipping audio loading")
        a.AssetsLoaded = true
        a.BGMEnabled = false
        a.SFXEnabled = false
        return
    }
    
    fmt.Println("🎵 Loading audio assets...")
     
    
    // Define audio files to load
    bgmFiles := map[string]string{
        "intro_theme":     "assets/audio/intro_theme.ogg",
        "menu_theme":      "assets/audio/menu_theme.ogg", 
        "game_theme":      "assets/audio/game_theme.ogg",
        "power_mode":      "assets/audio/power_mode.ogg",
        "boss_theme":      "assets/audio/boss_theme.ogg",
    }
    
    sfxFiles := map[string]string{
        "menu_select":         "assets/audio/menu_select.wav",
        "transition":          "assets/audio/transition.wav",
        "character_reveal":    "assets/audio/character_reveal.wav",
        "game_start":          "assets/audio/game_start.wav",
        "round_start":         "assets/audio/round_start.wav",
        "round_complete":      "assets/audio/round_complete.wav",
        "pellet_eat":          "assets/audio/pellet_eat.wav",
        "power_pellet":        "assets/audio/power_pellet.wav",
        "power_pellet_warning": "assets/audio/power_warning.wav",
        "power_pellet_end":    "assets/audio/power_end.wav",
        "ghost_eaten":         "assets/audio/ghost_eaten.wav",
        "player_death":        "assets/audio/player_death.wav",
        "game_over":           "assets/audio/game_over.wav",
        "pause":               "assets/audio/pause.wav",
        "unpause":             "assets/audio/unpause.wav",
    }
    
    // Load BGM files
    for name, path := range bgmFiles {
        data, err := a.loadAudioFile(path)
        if err != nil {
        a.LoadErrors = append(a.LoadErrors, fmt.Sprintf("BGM %s: %v", name, err))
            // Create silence as fallback using actual context sample rate
            contextSampleRate := 48000
            if a.AudioContext != nil {
                contextSampleRate = a.AudioContext.SampleRate()
            }
            data = a.createSilence(5 * contextSampleRate) // 5 seconds of silence
        }
        a.BGMData[name] = data
    }
    
    // Load SFX files
    for name, path := range sfxFiles {
        data, err := a.loadAudioFile(path)
        if err != nil {
            a.LoadErrors = append(a.LoadErrors, fmt.Sprintf("SFX %s: %v", name, err))
            // Create short beep as fallback
            data = a.createBeep(0.2, 440) // 0.2 second beep at 440Hz
        }
        a.SFXData[name] = data
    }
    
    a.AssetsLoaded = true
    
    if len(a.LoadErrors) > 0 {
        fmt.Printf("⚠️  Audio loading completed with %d errors:\n", len(a.LoadErrors))
        for _, err := range a.LoadErrors {
            fmt.Printf("   - %s\n", err)
        }
    } else {
        fmt.Printf("✅ All audio assets loaded successfully!\n")
    }
}

func (a *AudioSystem) loadAudioFile(path string) ([]byte, error) {
    // Check if file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("file not found: %s", path)
    }
    
    // Open file
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    // Determine file format by extension
    ext := filepath.Ext(path)
    var stream io.Reader
    
    // Get the actual sample rate being used by the audio context
    contextSampleRate := 22050
    if a.AudioContext != nil {
        contextSampleRate = a.AudioContext.SampleRate()
    }
    
    switch ext {
    case ".ogg":
        stream, err = vorbis.DecodeWithSampleRate(contextSampleRate, file)
        if err != nil {
            return nil, fmt.Errorf("failed to decode OGG: %v", err)
        }
    case ".wav":
        stream, err = wav.DecodeWithSampleRate(contextSampleRate, file)
        if err != nil {
            return nil, fmt.Errorf("failed to decode WAV: %v", err)
        }
    default:
        return nil, fmt.Errorf("unsupported audio format: %s", ext)
    }
    
    // Read all data
    data, err := io.ReadAll(stream)
    if err != nil {
        return nil, fmt.Errorf("failed to read audio data: %v", err)
    }
    
    fmt.Printf("   Loaded: %s (%d bytes)\n", filepath.Base(path), len(data))
    return data, nil
}

func (a *AudioSystem) createSilence(samples int) []byte {
    // Create silence buffer (16-bit stereo)
    data := make([]byte, samples*4) // 2 channels * 2 bytes per sample
    return data
}

func (a *AudioSystem) createBeep(duration float64, frequency float64) []byte {
    samples := int(duration * SampleRate)
    data := make([]byte, samples*4) // 16-bit stereo
    
    for i := 0; i < samples; i++ {
        // Generate sine wave
        t := float64(i) / SampleRate
        sample := int16(16000 * math.Sin(2*math.Pi*frequency*t))
        
        // Write stereo samples (little endian)
        pos := i * 4
        data[pos] = byte(sample)
        data[pos+1] = byte(sample >> 8)
        data[pos+2] = byte(sample)
        data[pos+3] = byte(sample >> 8)
    }
    
    return data
}

func (a *AudioSystem) PlayBGM(name string) {
    if a == nil || !a.BGMEnabled || !a.AssetsLoaded {
        return
    }
    
    // Stop current BGM if playing
    if a.BGMPlayer != nil && a.BGMPlayer.IsPlaying() {
        a.BGMPlayer.Pause()
    }
    
    // Get audio data
    data, exists := a.BGMData[name]
    if !exists {
        log.Printf("BGM not found: %s", name)
        return
    }
    
    // Create new player
    reader := bytes.NewReader(data)
    player, err := a.AudioContext.NewPlayer(reader)
    if err != nil {
        log.Printf("Failed to create BGM player for %s: %v", name, err)
        return
    }
    
    // Set volume and loop
    player.SetVolume(a.BGMVolume)
    
    // Start playing
    a.BGMPlayer = player
    a.CurrentBGM = name
    a.BGMPlayer.Play()
    
    // For looping BGM, we'll need to handle this in Update()
    fmt.Printf("🎵 Playing BGM: %s\n", name)
}

func (a *AudioSystem) PlaySFX(name string) {
    if a == nil || !a.SFXEnabled || !a.AssetsLoaded {
        return
    }
    
    // Get audio data
    data, exists := a.SFXData[name]
    if !exists {
        log.Printf("SFX not found: %s", name)
        return
    }
    
    // Create new player for this sound effect
    reader := bytes.NewReader(data)
    player, err := a.AudioContext.NewPlayer(reader)
    if err != nil {
        log.Printf("Failed to create SFX player for %s: %v", name, err)
        return
    }
    
    // Set volume
    player.SetVolume(a.SFXVolume)
    
    // Play the sound
    player.Play()
    
    // Store reference (will be cleaned up in Update())
    a.SFXPlayers[name+fmt.Sprint(time.Now().UnixNano())] = player
    
    fmt.Printf("🔊 Playing SFX: %s\n", name)
}

func (a *AudioSystem) Update() {
    if a == nil {
        return
    }
    
    // Handle BGM looping
    if a.BGMPlayer != nil && !a.BGMPlayer.IsPlaying() && a.BGMEnabled && a.CurrentBGM != "" {
        // Restart the BGM for looping
        if data, exists := a.BGMData[a.CurrentBGM]; exists {
            reader := bytes.NewReader(data)
            if newPlayer, err := a.AudioContext.NewPlayer(reader); err == nil {
                newPlayer.SetVolume(a.BGMVolume)
                newPlayer.Play()
                a.BGMPlayer = newPlayer
            }
        }
    }
    
    // Clean up finished SFX players
    for key, player := range a.SFXPlayers {
        if !player.IsPlaying() {
            player.Close()
            delete(a.SFXPlayers, key)
        }
    }
    
    // Handle fade effects
    if a.BGMFadeOut {
        a.BGMFadeTimer++
        if a.BGMFadeTimer >= 60 { // 1 second fade
            a.StopBGM()
            a.BGMFadeOut = false
            a.BGMFadeTimer = 0
        } else {
            fadeVolume := a.BGMVolume * (1.0 - float64(a.BGMFadeTimer)/60.0)
            if a.BGMPlayer != nil {
                a.BGMPlayer.SetVolume(fadeVolume)
            }
        }
    }
    
    if a.BGMFadeIn {
        a.BGMFadeTimer++
        if a.BGMFadeTimer >= 60 {
            a.BGMFadeIn = false
            a.BGMFadeTimer = 0
        } else {
            fadeVolume := a.BGMVolume * (float64(a.BGMFadeTimer) / 60.0)
            if a.BGMPlayer != nil {
                a.BGMPlayer.SetVolume(fadeVolume)
            }
        }
    }
}

func (a *AudioSystem) StopBGM() {
    if a == nil {
        return
    }
    
    if a.BGMPlayer != nil {
        a.BGMPlayer.Pause()
        a.BGMPlayer = nil
    }
    a.CurrentBGM = ""
    fmt.Println("🎵 BGM stopped")
}

func (a *AudioSystem) FadeBGM() {
    if a == nil {
        return
    }
    
    if a.BGMPlayer != nil && a.BGMPlayer.IsPlaying() {
        a.BGMFadeOut = true
        a.BGMFadeTimer = 0
    }
}

func (a *AudioSystem) SetBGMVolume(volume float64) {
    if a == nil {
        return
    }
    
    a.BGMVolume = math.Max(0.0, math.Min(MaxVolume, volume))
    if a.BGMPlayer != nil {
        a.BGMPlayer.SetVolume(a.BGMVolume)
    }
    fmt.Printf("🎵 BGM Volume: %.1f\n", a.BGMVolume)
}

func (a *AudioSystem) SetSFXVolume(volume float64) {
    if a == nil {
        return
    }
    
    a.SFXVolume = math.Max(0.0, math.Min(MaxVolume, volume))
    fmt.Printf("🔊 SFX Volume: %.1f\n", a.SFXVolume)
}

func (a *AudioSystem) ToggleBGM() {
    if a == nil {
        return
    }
    
    a.BGMEnabled = !a.BGMEnabled
    if !a.BGMEnabled {
        a.StopBGM()
    }
    fmt.Printf("🎵 BGM %s\n", map[bool]string{true: "enabled", false: "disabled"}[a.BGMEnabled])
}

func (a *AudioSystem) ToggleSFX() {
    if a == nil {
        return
    }
    
    a.SFXEnabled = !a.SFXEnabled
    
    if !a.SFXEnabled {
        // Stop all current SFX
        for key, player := range a.SFXPlayers {
            player.Pause()
            player.Close()
            delete(a.SFXPlayers, key)
        }
    }
    
    fmt.Printf("🔊 SFX %s\n", map[bool]string{true: "enabled", false: "disabled"}[a.SFXEnabled])
}

func (a *AudioSystem) GetStatus() string {
    if a == nil {
        return "Audio system not initialized"
    }
    
    status := fmt.Sprintf("BGM: %s (%.1f) | SFX: %s (%.1f)",
        map[bool]string{true: "ON", false: "OFF"}[a.BGMEnabled],
        a.BGMVolume,
        map[bool]string{true: "ON", false: "OFF"}[a.SFXEnabled],
        a.SFXVolume)
    
    if a.CurrentBGM != "" {
        status += fmt.Sprintf(" | Playing: %s", a.CurrentBGM)
    }
    
    return status
}

// Convenience methods for specific game events
func (a *AudioSystem) PlayMenuMusic() {
    a.PlayBGM("menu_theme")
}

func (a *AudioSystem) PlayGameMusic() {
    a.PlayBGM("game_theme")
}

func (a *AudioSystem) PlayIntroMusic() {
    a.PlayBGM("intro_theme")
}

func (a *AudioSystem) PlayPowerMode() {
    a.PlayBGM("power_mode")
}

func (a *AudioSystem) EndPowerMode() {
    a.PlayBGM("game_theme")
}

// Special effect: Play power pellet music temporarily
func (a *AudioSystem) TriggerPowerPelletMode(duration int) {
    if a == nil || !a.SFXEnabled {
        return
    }
    
    // Play power mode music
    a.PlayBGM("power_mode")
    
    // TODO: Set timer to restore previous music after duration
    // This would require integration with the game's update loop
}

func (a *AudioSystem) Cleanup() {
    if a == nil {
        return
    }
    
    // Stop all audio
    a.StopBGM()
    
    // Close all SFX players
    for key, player := range a.SFXPlayers {
        player.Close()
        delete(a.SFXPlayers, key)
    }
    
    fmt.Println("🎵 Audio system cleaned up")
}

// Integration helper: Replace the simple SoundManager in your existing code
type SoundManager struct {
    *AudioSystem
    BGMEnabled bool
    SFXEnabled bool
    Volume     float64
}

func NewSoundManager() *SoundManager {
    fmt.Println("🔊 Creating new SoundManager...")
    audioSys := NewAudioSystem()
    
    if audioSys == nil {
        fmt.Println("❌ Failed to create AudioSystem for SoundManager")
        return &SoundManager{
            AudioSystem: nil,
            BGMEnabled:  true,
            SFXEnabled:  true,
            Volume:      DefaultVolume,
        }
    }
    
    fmt.Println("✅ SoundManager created successfully")
    return &SoundManager{
        AudioSystem: audioSys,
        BGMEnabled:  true,
        SFXEnabled:  true,
        Volume:      DefaultVolume,
    }
}

// Maintain compatibility with existing interface
func (s *SoundManager) PlayBGM(track string) {
    if s.AudioSystem != nil {
        s.AudioSystem.PlayBGM(track)
    }
}

func (s *SoundManager) PlaySFX(effect string) {
    if s.AudioSystem != nil {
        s.AudioSystem.PlaySFX(effect)
    }
}

func (s *SoundManager) StopBGM() {
    if s.AudioSystem != nil {
        s.AudioSystem.StopBGM()
    }
}

func (s *SoundManager) SetVolume(volume float64) {
    s.Volume = volume
    if s.AudioSystem != nil {
        s.AudioSystem.SetBGMVolume(volume)
        s.AudioSystem.SetSFXVolume(volume)
    }
}

func (s *SoundManager) ToggleBGM() {
    s.BGMEnabled = !s.BGMEnabled
    if s.AudioSystem != nil {
        s.AudioSystem.ToggleBGM()
    }
}

func (s *SoundManager) ToggleSFX() {
    s.SFXEnabled = !s.SFXEnabled
    if s.AudioSystem != nil {
        s.AudioSystem.ToggleSFX()
    }
}
