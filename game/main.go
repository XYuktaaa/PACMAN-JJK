

// package main

// import (
//     "github.com/hajimehoshi/ebiten/v2"
//     "log"
// )

// func main() {
//     // First load basic assets for the game
//     LoadAssets()
    
//     // Create the game instance
//     game := NewGame()
    
//     // Load images/gifs for the menu with error handling
//     var logo *ebiten.Image
//     var characterFrames []*ebiten.Image
//     var bg *ebiten.Image
    
//     // Load logo with error handling
//     logo = loadImage("assets/jogo.png")
//     if logo == nil {
//         log.Printf("Warning: Failed to load logo image from assets/jogo.png")
//     }
    
//     // Load character GIF with error handling
//     var err error
//     characterFrames, err = LoadGIF("assets/gojo.gif")
//     if err != nil {
//         log.Printf("Warning: Failed to load character GIF: %v", err)
//         // Try loading as a static image instead
//         staticChar := loadImage("assets/gojo.png")
//         if staticChar != nil {
//             characterFrames = []*ebiten.Image{staticChar}
//             log.Printf("Loaded static character image as fallback")
//         }
//     } else if len(characterFrames) == 0 {
//         log.Printf("Warning: GIF loaded but contains no frames")
//     } else {
//         log.Printf("Successfully loaded %d GIF frames", len(characterFrames))
//     }
    
//     // Load background with error handling
//     bg = loadImage("assets/cursed_bg.png")
//     if bg == nil {
//         log.Printf("Warning: Failed to load background image from assets/cursed_bg.png")
//         // Try alternative background names
//         bg = loadImage("assets/136290597-satoru-gojou-purple.png")
//         if bg == nil {
//             log.Printf("Warning: Failed to load alternative background image")
//         }
//     }
    
//     // Set the images in the UI (even if some are nil)
//     game.menuUI.SetImages(logo, characterFrames, bg)
    
//     // Log what was successfully loaded
//     assetsLoaded := 0
//     if logo != nil { assetsLoaded++ }
//     if len(characterFrames) > 0 { assetsLoaded++ }
//     if bg != nil { assetsLoaded++ }
    
//     log.Printf("Successfully loaded %d out of 3 menu assets", assetsLoaded)
    
//     // Set window properties
//     ebiten.SetWindowSize(1200, 800)
//     ebiten.SetWindowTitle("Jujutsu Kaisen Pac-Man")
    
//     // Run the game
//     if err := ebiten.RunGame(game); err != nil {
//         log.Fatal(err)
//     }
// }

// // Helper function to safely load an image
// func safeLoadImage(path string) *ebiten.Image {
//     img := loadImage(path)
//     if img != nil {
//         log.Printf("Successfully loaded image: %s", path)
//     } else {
//         log.Printf("Failed to load image: %s", path)
//     }
//     return img
// }

// // Helper function to safely load a GIF
// func safeLoadGIF(path string) ([]*ebiten.Image, error) {
//     frames, err := LoadGIF(path)
//     if err != nil {
//         log.Printf("Failed to load GIF: %s - Error: %v", path, err)
//         return nil, err
//     }
//     if len(frames) == 0 {
//         log.Printf("GIF loaded but contains no frames: %s", path)
//         return nil, nil
//     }
//     log.Printf("Successfully loaded GIF with %d frames: %s", len(frames), path)
//     return frames, nil
// }


package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "log"
)

func main() {
    // First load basic game assets
    LoadAssets()
    
    // Create the game instance
    game := NewGame()
    
    // Load UI-specific images/gifs for the menu
    logo := loadImage("assets/jogo.png")
    
    // Load character GIF frames for menu background
    characterFrames, err := LoadGIF("assets/gojo.gif")
    if err != nil {
        log.Printf("Failed to load character GIF: %v", err)
        // Continue without the GIF - will show placeholder
        characterFrames = []*ebiten.Image{}
    }
    
    // Load background image (this will be the main background)
    bg := loadImage("assets/cursed_bg.png")
    if bg == nil {
        log.Printf("Warning: Background image not found, using gradient background")
    }
    
    // Set the images in the UI
    // The GIF will now appear as background in the menu area
    game.menuUI.SetImages(logo, characterFrames, bg)
    
    // Set window properties
    ebiten.SetWindowSize(1200, 800)
    ebiten.SetWindowTitle("Jujutsu Kaisen Pac-Man")
    
    // Run the game
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}

/*
Asset Loading Tips:

1. Background Image (cursed_bg.png):
   - Should be large enough to cover 1200x800 screen
   - Will be scaled to fit screen while maintaining aspect ratio
   - Recommended size: 1920x1080 or larger
   - Used as main background with semi-transparent gradient overlay

2. Character GIF (gojo.gif):
   - Now used as background for the menu panel
   - Will be scaled and cropped to fit menu dimensions (480x360 approx)
   - Rendered at 40% opacity so menu text remains readable
   - Recommended size: 400x400 to 800x800 for best quality
   - Animation speed controlled by frameDelay (currently 8 ticks = ~133ms per frame)

3. Logo Image (jogo.png):
   - Used in top-left corner
   - Scaled to fit 80x80 pixel area while maintaining aspect ratio
   - Recommended size: 128x128 or 256x256

File Structure:
assets/
├── cursed_bg.png     (Main background - large image)
├── gojo.gif          (Menu background animation - medium size)
├── jogo.png          (Logo - small square image)
└── ... (other game assets)

If images are too large and causing performance issues:
- Resize cursed_bg.png to 1920x1080 maximum
- Resize gojo.gif frames to 512x512 maximum
- Resize jogo.png to 256x256 maximum

The code now handles:
- Automatic scaling to fit screen/areas
- Proper aspect ratio maintenance
- Fallback to placeholders if images fail to load
- GIF animation in menu background with proper opacity
- Smooth background blending
*/
