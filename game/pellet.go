package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    //"github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "math"

    "image/color"
    
)

type Pellet struct{
    X,Y float64
    Eaten bool
    IsPower bool
    Image *ebiten.Image
}

var(
    Pellets          []*Pellet
    PowerPelletImage *ebiten.Image
    PelletImage      *ebiten.Image
)

func InitPellets(level [][]int, TileSize int) {
	Pellets = []*Pellet{}

	PelletImage = createRegularPelletImage()
	PowerPelletImage = createPowerPelletImage()

	for y := 0; y < len(level); y++ {
		for x := 0; x < len(level[y]); x++ {
			tile := level[y][x]
			px := float64(x*TileSize + TileSize/2 - PelletImage.Bounds().Dx()/2)
			py := float64(y*TileSize + TileSize/2 - PelletImage.Bounds().Dy()/2)

			switch tile {
			case 2: // Normal Pellet
				Pellets = append(Pellets, &Pellet{
					X:     px,
					Y:     py,
					Image: PelletImage,
					IsPower: false,
					Eaten: false,
				})
			case 4: // Power Pellet
				Pellets = append(Pellets, &Pellet{
					X:     px,
					Y:     py,
					Image: PowerPelletImage,
					IsPower: true,
					Eaten: false,
				})
			}
		}
	}
}


func createRegularPelletImage() *ebiten.Image {
    img := ebiten.NewImage(8, 8)
    for y := 0; y < 8; y++ {
        for x := 0; x < 8; x++ {
            dx := float64(x - 4)
            dy := float64(y - 4)
            dist := math.Hypot(dx, dy)
            if dist < 2 {
                img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Core red
            } else if dist < 3.5 {
                img.Set(x, y, color.RGBA{255, 0, 0, 128}) // Outer glow
            }
        }
    }
    return img
}
func createPowerPelletImage() *ebiten.Image {
    img := ebiten.NewImage(10, 10)
    for y := 0; y < 10; y++ {
        for x := 0; x < 10; x++ {
            dx := float64(x - 5)
            dy := float64(y - 5)
            dist := math.Hypot(dx, dy)
            if dist > 3.2 && dist < 4.8 {
                img.Set(x, y, color.RGBA{128, 0, 255, 220}) // Hollow purple ring
            }
        }
    }
    return img
}


func UpdatePellets(playerX, playerY float64, playerWidth, playerHeight int) {
	for _, pellet := range Pellets {
		if pellet.Eaten {
			continue
		}

		pelletCenterX := pellet.X + float64(PelletImage.Bounds().Dx())/2
		pelletCenterY := pellet.Y + float64(PelletImage.Bounds().Dy())/2

		if pelletCenterX >= playerX &&
			pelletCenterX <= playerX+float64(playerWidth) &&
			pelletCenterY >= playerY &&
			pelletCenterY <= playerY+float64(playerHeight) {

			pellet.Eaten = true
			// Add scoring logic here if needed
		}
	}
}

func DrawPellets(screen *ebiten.Image) {
    	for _, pellet := range Pellets {
		if pellet.Eaten {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pellet.X, pellet.Y)
		if pellet.IsPower {
			screen.DrawImage(PowerPelletImage, op)
		} else {
			screen.DrawImage(PelletImage, op)
		}
	}
}

func CheckPelletCollision(playerX, playerY float64, size int) {
	for _, pellet := range Pellets {
		if pellet.Eaten {
			continue
		}
		// Simple bounding box overlap
if playerX < float64(pellet.X + float64(size)) &&
   playerX + float64(size) > float64(pellet.X) &&
   playerY < float64(pellet.Y + float64(size)) &&
   playerY + float64(size) > float64(pellet.Y){
			pellet.Eaten = true

			if pellet.IsPower {
				// Trigger power pellet effect
			}
		}
	}
}

