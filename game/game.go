package game

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{
    player Player
    ghosts []Ghost
    maze [][]int
    int score
    
}

func NewGame() *Game {
    return &Game{}
}

func (g *Game) Update() error {
    g.player.Update()
    for i := range g.ghosts{
        g.ghosts[i].Update()
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    ebitenutil.DebugPrint(screen, "PACMAN - GOJO VS CURSES")
    g.drawMaze(screen)
    g.player.Draw(screen)
    for _,ghost := range g.ghosts {
        ghost.Draw(screen)
    }

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 640, 640
}

