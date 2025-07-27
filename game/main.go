package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "log"
)

func main (){
    LoadAssets()
    game:= NewGame()
    ebiten.SetWindowSize(640, 480)
    ebiten.SetWindowTitle("JJK PACMAN: GOJO vs CURSES")
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
