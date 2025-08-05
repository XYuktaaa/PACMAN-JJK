// package main

// import (
//     "github.com/hajimehoshi/ebiten/v2"
//     "log"
// )

// func main (){
//     LoadAssets()
//     game:= NewGame()
//     ebiten.SetWindowSize(640, 480)
//     ebiten.SetWindowTitle("JJK PACMAN: GOJO vs CURSES")
//     if err := ebiten.RunGame(game); err != nil {
//         log.Fatal(err)
//     }




// // Load your images/gifs
//     logo := loadImage("assets/jogo.png")
//     characterFrame,_ := LoadGIF("assets/gojo.gif")
//     character := characterFrame[0]//this is the fix i added but not sure
//     bg := loadImage("assets/136290597-satoru-gojou-purple.png")
    
//     game = &Game{
//         menuUI: NewUIPage(),
//         //showMenu: true,
//         logoImg: logo,
//         characterGif: character,
//         bgTexture: bg,
//     }
    
//     // Set the images in the UI
//     game.menuUI.SetImages(logo, character, bg)
// //    game.menuUI.SetImages(logo, character[0], bg)

//     ebiten.SetWindowSize(1200, 800)
//     ebiten.SetWindowTitle("Jujutsu Kaisen Pac-Man")
//     ebiten.RunGame(game)
// }

package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "log"
)

func main() {
    // First load basic assets
    LoadAssets()
    
    // Load your images/gifs for the menu
    logo := loadImage("assets/jogo.png")
    characterFrames, err := LoadGIF("assets/gojo.gif")
    if err != nil {
        log.Printf("Failed to load character GIF: %v", err)
        // You can continue without the GIF or use a fallback image
        characterFrames = []*ebiten.Image{nil}
    }
    
    var character *ebiten.Image
    if len(characterFrames) > 0 && characterFrames[0] != nil {
        character = characterFrames[0]
    }
    
    bg := loadImage("assets/cursed_bg.png")
    
    // Create the game with menu
//     game := &Game{
//         menuUI:       NewUIPage(),
// //        showMenu:     true,  // Make sure to show menu initially
//         logoImg:      logo,
//         characterGif: character,
//         bgTexture:    bg,
//     }

game := NewGame()

    // Set the images in the UI
    game.menuUI.SetImages(logo, character, bg)
    
    // Set window properties
    ebiten.SetWindowSize(1200, 800)
    ebiten.SetWindowTitle("Jujutsu Kaisen Pac-Man")
    
    // Run the game
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
