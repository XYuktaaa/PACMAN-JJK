package main

import (
    "fmt"
    "image"
    "image/gif"
    "image/draw"
    "image/color"
    "os"
    "log"
    _ "image/png"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    _ "github.com/hajimehoshi/ebiten/v2/text"
    _ "embed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	//"io/ioutil"
)
//go:embed assets/PressStart2P-Regular.ttf
var fontBytes []byte

var (
	bigfont        font.Face
	PressStartFont *opentype.Font
)
const (
    TileEmpty       = 0
    TileWall        = 1
    TilePellet      = 2
    TilePlayer      = 3
    TilePowerPellet = 4
)
var level = [][]int{
    {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
    {1,2,2,2,2,2,2,2,2,2,2,2,1,2,2,2,2,2,2,2,2,2,2,2,2,2,1},
    {1,4,1,1,1,1,2,1,1,1,1,2,1,2,1,1,1,1,2,1,1,1,1,1,4,2,1},
    {1,2,1,1,1,1,2,1,1,1,1,2,1,2,1,1,1,1,2,1,1,1,1,1,2,2,1},
    {1,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,1},
    {1,2,1,1,1,1,2,1,1,2,1,1,1,1,1,2,1,1,2,1,1,1,1,1,2,2,1},
    {1,2,1,1,1,1,2,1,1,2,1,1,1,1,1,2,1,1,2,1,1,1,1,1,2,2,1},
    {1,2,2,2,2,2,2,1,1,2,2,2,1,2,2,2,1,1,2,2,2,2,2,2,2,2,1},
    {1,1,1,1,1,1,2,1,1,1,1,0,1,0,1,1,1,1,2,1,1,1,1,1,1,1,1},
    {0,0,0,0,0,1,2,1,1,1,1,0,1,0,1,1,1,1,2,1,0,0,0,0,0,0,0},
    {0,0,0,0,0,1,2,1,1,0,0,0,0,0,0,0,1,1,2,1,0,0,0,0,0,0,0},
    {0,0,0,0,0,1,2,1,1,0,1,1,0,1,1,0,1,1,2,1,0,0,0,0,0,0,0},
    {1,1,1,1,1,1,2,1,1,0,1,0,0,0,1,0,1,1,2,1,1,1,1,1,1,1,1},
    {0,0,0,0,0,0,2,0,0,0,1,0,0,0,1,0,0,0,2,0,0,0,0,0,0,0,0},
    {1,1,1,1,1,1,2,1,1,0,1,0,0,0,1,0,1,1,2,1,1,1,1,1,1,1,1},
    {0,0,0,0,0,1,2,1,1,0,1,1,1,1,1,0,1,1,2,1,0,0,0,0,0,0,0},
    {0,0,0,0,0,1,2,1,1,0,0,0,0,0,0,0,1,1,2,1,0,0,0,0,0,0,0},
    {0,0,0,0,0,1,2,1,1,1,1,0,1,0,1,1,1,1,2,1,0,0,0,0,0,0,0},
    {1,1,1,1,1,1,2,1,1,1,1,0,1,0,1,1,1,1,2,1,1,1,1,1,1,1,1},
    {1,2,2,2,2,2,2,2,2,2,2,2,1,2,2,2,2,2,2,2,2,2,2,2,2,2,1},
    {1,2,1,1,1,1,2,1,1,1,1,2,1,2,1,1,1,1,2,1,1,1,1,1,2,2,1},
    {1,2,1,1,1,1,2,1,1,1,1,2,1,2,1,1,1,1,2,1,1,1,1,1,2,2,1},
    {1,4,2,2,1,1,2,2,2,2,2,2,2,2,2,2,2,2,2,1,1,2,2,2,4,2,1},
    {1,1,1,2,1,1,2,1,1,2,1,1,1,1,1,2,1,1,2,1,1,2,1,1,1,1,1},
    {1,1,1,2,1,1,2,1,1,2,1,1,1,1,1,2,1,1,2,1,1,2,1,1,1,1,1},
    {1,2,2,2,2,2,2,1,1,2,2,2,1,2,2,2,1,1,2,2,2,2,2,2,2,2,1},
    {1,2,1,1,1,1,1,1,1,1,1,2,1,2,1,1,1,1,1,1,1,1,1,1,2,2,1},
    {1,2,1,1,1,1,1,1,1,1,1,2,1,2,1,1,1,1,1,1,1,1,1,1,2,2,1},
    {1,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,1},
    {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
}
var (
    WallImage    *ebiten.Image
    // PelletImage  *ebiten.Image
    FloorImage   *ebiten.Image
    PlayerImage  *ebiten.Image
)

func LoadAssets() {
    var err error
    fmt.Println("Loading wall image from assets/wall.png")

    WallImage, _, err = ebitenutil.NewImageFromFile("assets/wall.png")
    if err != nil {
        log.Fatal("failed to load wall.png",err)
    }
    FloorImage, _, err = ebitenutil.NewImageFromFile("assets/floor.png")
    if err != nil {
        log.Fatal("failed to load floor.png",err)
    }
    // PelletImage, _, err = ebitenutil.NewImageFromFile("assets/pellet.png")
    // if err != nil {
    //     log.Fatal("failed to load pellet",err)
    // }
    PlayerImage, _, err = ebitenutil.NewImageFromFile("assets/player.png")
    if err != nil {
        log.Fatal("failed to load player",err)
    }
    ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	bigfont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}

    // PelletImage = ebiten.NewImage(8,8)
    // PelletImage.Fill(color.White)
}

func loadImage(path string) *ebiten.Image {
    fmt.Printf("Trying to load image: %s\n", path)

    // Check if file exists first
    if _, err := os.Stat(path); os.IsNotExist(err) {
        log.Printf("Warning: Image file does not exist: %s", path)
        // Return a placeholder image
        placeholder := ebiten.NewImage(100, 100)
        placeholder.Fill(color.RGBA{128, 128, 128, 255}) // Gray placeholder
        return placeholder
    }

    f, err := os.Open(path)
    if err != nil {
        log.Printf("Warning: Failed to open image file %s: %v", path, err)
        // Return a placeholder image
        placeholder := ebiten.NewImage(100, 100)
        placeholder.Fill(color.RGBA{128, 128, 128, 255}) // Gray placeholder
        return placeholder
    }
    defer f.Close()

    img, format, err := image.Decode(f)
    if err != nil {
        log.Printf("Warning: Failed to decode image %s: %v", path, err)
        // Return a placeholder image
        placeholder := ebiten.NewImage(100, 100)
        placeholder.Fill(color.RGBA{128, 128, 128, 255}) // Gray placeholder
        return placeholder
    }

    fmt.Printf("Successfully loaded image: %s (format: %s)\n", path, format)
    return ebiten.NewImageFromImage(img)
}

func LoadFont() {
    var err error
    PressStartFont, err = opentype.Parse(fontBytes)
    if err != nil {
        log.Printf("Warning: Failed to parse embedded font: %v", err)
        return
    }

    bigfont, err = opentype.NewFace(PressStartFont, &opentype.FaceOptions{
        Size:    24,
        DPI:     72,
        Hinting: font.HintingFull,
    })
    if err != nil {
        log.Printf("Warning: Failed to create font face: %v", err)
    }
}

func LoadGIF(path string) ([]*ebiten.Image, error) {
    fmt.Printf("Trying to load GIF: %s\n", path)

    // Check if file exists first
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("GIF file does not exist: %s", path)
    }

    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open GIF file %s: %v", path, err)
    }
    defer f.Close()

    g, err := gif.DecodeAll(f)
    if err != nil {
        return nil, fmt.Errorf("failed to decode GIF %s: %v", path, err)
    }

    if len(g.Image) == 0 {
        return nil, fmt.Errorf("GIF has no frames: %s", path)
    }

    var frames []*ebiten.Image
    for i, src := range g.Image {
        rgba := image.NewRGBA(src.Bounds())
        draw.Draw(rgba, rgba.Bounds(), src, image.Point{}, draw.Over)
        ebitenImg := ebiten.NewImageFromImage(rgba)
        frames = append(frames, ebitenImg)
        fmt.Printf("Loaded GIF frame %d/%d\n", i+1, len(g.Image))
    }

    fmt.Printf("Successfully loaded GIF: %s (%d frames)\n", path, len(frames))
    return frames, nil
}
