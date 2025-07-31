package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "fmt"
    "image/color"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
    
)

var glowColor = color.RGBA{255, 50, 50, 255} // bright red center
var outerGlowColor = color.RGBA{255, 50, 50, 100} // transparent red for glow


type Game struct{
    Player *Player
    Ghosts []*Ghost
    Pellet []Pellet
    menuUI *UIPage
    showMenu bool
}

const TileSize =32

func NewGame() *Game {
    InitPellets(level, TileSize)
    return &Game{
        Player: NewPlayer(float64(1*TileSize), float64(1*TileSize),"assets/player.png"),
        Ghosts: []*Ghost{
		NewGhost(64, 64, "assets/sakuna.png", "Sukuna",50),
	    NewGhost(160, 96, "assets/jogo.png", "jogo",50),
		NewGhost(96, 160, "assets/kenjaku.png", "Kenjaku",50),
		NewGhost(160, 160, "assets/mahito.png", "Mahito",50),
			},

		menuUI: NewUIPage(),
        showMenu: true,
    }
}


func (g *Game) Update()error{
    if g.Player != nil{
    g.Player.Update(level, TileSize)}
    if g.showMenu {
        err := g.menuUI.Update()
        if err != nil {
            return err
        }
        if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
    g.showMenu = !g.showMenu // Toggle menu
}
        // Check if user selected an option
        if g.menuUI.IsEnterPressed() {
            switch g.menuUI.GetSelectedOption() {
            case 0: // START
                g.showMenu = false
                // Start your game logic
            case 1: // PAUSE
                // Pause game logic
            case 2: // RESUME
                g.showMenu = false
                // Resume game logic
            case 3: // QUIT
                return fmt.Errorf("quit game")
            }
        }
        return nil
    }
        // var jogo *Ghost
    // for _, ghost:= range g.Ghosts{
    //     if ghost.Name == "jogo"{
    //         jogo = ghost
    //         break
    //     }
    // }
    // for _, ghost := range g.Ghosts{
    //     ghost.Update(level, TileSize, g.Player.X, g.Player.Y)
    // }
    return nil
}


func DrawMaze(screen *ebiten.Image) {
    for y, row := range level {
        for x, tile := range row {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(x * TileSize), float64(y * TileSize))

            switch tile {
            case TileWall:
                screen.DrawImage(WallImage, op)
            case TilePellet:
                screen.DrawImage(FloorImage, op) // base layer
            default:
                screen.DrawImage(FloorImage, op)
            }
        }
    }
}


func (g *Game) Draw(screen *ebiten.Image) {
     if g.showMenu {
        g.menuUI.Draw(screen)
        return
    }
    DrawMaze(screen) // draw background or maze first
	// DrawPellets(screen)

    for y, row := range level {
        for x, tile := range row {
            op := &ebiten.DrawImageOptions{}
            op.GeoM.Translate(float64(x*TileSize), float64(y*TileSize))

            switch tile {
            case TileWall:
                screen.DrawImage(WallImage, op)
            case TilePellet:
                screen.DrawImage(FloorImage, op)
                DrawPellets(screen)
				CheckPelletCollision(g.Player.X, g.Player.Y,g.Player.Size)


                // screen.DrawImage(PelletImage, op)
                // cx := float64(x*TileSize + TileSize/2)
                // cy := float64(y*TileSize + TileSize/2)

                // ebitenutil.DrawRect(screen, cx-4, cy-4, 8, 8, outerGlowColor)
                // ebitenutil.DrawRect(screen, cx-2.5, cy-2.5, 5, 5, color.RGBA{255, 50, 50, 180})
                // ebitenutil.DrawRect(screen, cx-1.5, cy-1.5, 3, 3, glowColor)
            case TileEmpty:
                screen.DrawImage(FloorImage, op)
            }
        }
    }
    for _, ghost := range g.Ghosts {
 	   ghost.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.Player.Score))


    g.Player.Draw(screen) // draw player AFTER level so it appears on top
}


func (g *Game) Layout(outsideWidth, outsideHeight int)(int , int){
    width:= len(level[0])*TileSize
    height := len(level)*TileSize
    return width,height  
}

func isWallColliding(level [][]int, px, py float64, size, TileSize int) bool {
	margin := 6 // number of pixels to reduce from each side
	corners := [][2]int{
    	{int(px) + margin, int(py) + margin},
    	{int(px + float64(size)) - margin, int(py) + margin},
    	{int(px) + margin, int(py + float64(size)) - margin},
    	{int(px + float64(size)) - margin, int(py + float64(size)) - margin},
	}
	for _, corner := range corners {
		cx := corner[0] / TileSize
		cy := corner[1] / TileSize
		fmt.Printf("Checking tile (%d, %d) = %d\n", cx, cy, level[cy][cx])
		

		// Prevent out-of-bounds access
		if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
			fmt.Println("Out of bounds, blocking movement")
			return true // treat OOB as wall
		}

		if level[cy][cx] == TileWall {
    		fmt.Println("Hit wall at", cx, cy)
			return true
		}

		fmt.Printf("Ghost (%.2f, %.2f), checking tile (%d, %d) = %d\n", px, py, cx, cy, level[cy][cx])
	}
	

	return false
}

func isWallCollidingLenient(level [][]int, px, py float64, size, tileSize int) bool {
	offset := 22.0 // shrink hitbox by 5px from each side

	corners := [][2]int{
		{int(px + offset), int(py + offset)},
		{int(px + float64(size) - offset), int(py + offset)},
		{int(px + offset), int(py + float64(size) - offset)},
		{int(px + float64(size) - offset), int(py + float64(size) - offset)},
	}

	for _, corner := range corners {
		cx := corner[0] / tileSize
		cy := corner[1] / tileSize

		if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
			return true
		}

		if level[cy][cx] == TileWall {
			return true
		}
	}

	return false
}

func isWallCollidingStrict(level [][]int, px, py float64, size, tileSize int) bool {
	corners := [][2]int{
		{int(px), int(py)},
		{int(px + float64(size) - 1), int(py)},
		{int(px), int(py + float64(size) - 1)},
		{int(px + float64(size) - 1), int(py + float64(size) - 1)},
	}

	for _, corner := range corners {
		cx := corner[0] / tileSize
		cy := corner[1] / tileSize

		if cy < 0 || cy >= len(level) || cx < 0 || cx >= len(level[0]) {
			return true
		}

		if level[cy][cx] == TileWall {
			return true
		}
	}

	return false
}

