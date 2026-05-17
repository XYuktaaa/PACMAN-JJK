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
