

package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "log"
)

func main() {
    // First load basic assets for the game
    LoadAssets()
    
    // Create the game instance
    game := NewGame()
    
    // Load images/gifs for the menu with error handling
    var logo *ebiten.Image
    var characterFrames []*ebiten.Image
    var bg *ebiten.Image
    
    // Load logo with error handling
    logo = loadImage("assets/jogo.png")
    if logo == nil {
        log.Printf("Warning: Failed to load logo image from assets/jogo.png")
    }
    
    // Load character GIF with error handling
    var err error
    characterFrames, err = LoadGIF("assets/gojo.gif")
    if err != nil {
        log.Printf("Warning: Failed to load character GIF: %v", err)
        // Try loading as a static image instead
        staticChar := loadImage("assets/gojo.png")
        if staticChar != nil {
            characterFrames = []*ebiten.Image{staticChar}
            log.Printf("Loaded static character image as fallback")
        }
    } else if len(characterFrames) == 0 {
        log.Printf("Warning: GIF loaded but contains no frames")
    } else {
        log.Printf("Successfully loaded %d GIF frames", len(characterFrames))
    }
    
    // Load background with error handling
    bg = loadImage("assets/cursed_bg.png")
    if bg == nil {
        log.Printf("Warning: Failed to load background image from assets/cursed_bg.png")
        // Try alternative background names
        bg = loadImage("assets/136290597-satoru-gojou-purple.png")
        if bg == nil {
            log.Printf("Warning: Failed to load alternative background image")
        }
    }
    
    // Set the images in the UI (even if some are nil)
    game.menuUI.SetImages(logo, characterFrames, bg)
    
    // Log what was successfully loaded
    assetsLoaded := 0
    if logo != nil { assetsLoaded++ }
    if len(characterFrames) > 0 { assetsLoaded++ }
    if bg != nil { assetsLoaded++ }
    
    log.Printf("Successfully loaded %d out of 3 menu assets", assetsLoaded)
    
    // Set window properties
    ebiten.SetWindowSize(1200, 800)
    ebiten.SetWindowTitle("Jujutsu Kaisen Pac-Man")
    
    // Run the game
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}

// Helper function to safely load an image
func safeLoadImage(path string) *ebiten.Image {
    img := loadImage(path)
    if img != nil {
        log.Printf("Successfully loaded image: %s", path)
    } else {
        log.Printf("Failed to load image: %s", path)
    }
    return img
}

// Helper function to safely load a GIF
func safeLoadGIF(path string) ([]*ebiten.Image, error) {
    frames, err := LoadGIF(path)
    if err != nil {
        log.Printf("Failed to load GIF: %s - Error: %v", path, err)
        return nil, err
    }
    if len(frames) == 0 {
        log.Printf("GIF loaded but contains no frames: %s", path)
        return nil, nil
    }
    log.Printf("Successfully loaded GIF with %d frames: %s", len(frames), path)
    return frames, nil
}
