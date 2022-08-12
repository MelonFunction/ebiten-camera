package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	camera "github.com/melonfunction/ebiten-camera"
)

//go:embed sprites.png
var embedded embed.FS

// vars
var (
	cam              *camera.Camera
	sprites          *ebiten.Image
	spriteSize       = 16
	spriteScale      = 8
	LastWindowWidth  int
	LastWindowHeight int
	rotation         float64

	ErrNormalExit = errors.New("Normal exit")
)

var (
	rectRed    = image.Rect(0, 0, spriteSize, spriteSize)
	rectGreen  = image.Rect(spriteSize, 0, spriteSize*2, spriteSize)
	rectBlue   = image.Rect(spriteSize*2, 0, spriteSize*3, spriteSize)
	rectYellow = image.Rect(spriteSize*3, 0, spriteSize*4, spriteSize)
)

// Point represents a point in space
type Point struct {
	X, Y float64
}

// NewPoint returns a new *Point
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// Rotate rotates a point about 0,0
func (p *Point) Rotate(phi float64) *Point {
	c, s := math.Cos(phi), math.Sin(phi)
	return &Point{
		X: c*p.X - s*p.Y,
		Y: s*p.X + c*p.Y,
	}
}

// AngleTo returns the angle between p to other
func (p *Point) AngleTo(other *Point) float64 {
	return math.Atan2(p.Y-other.Y, p.X-other.X)
}

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrNormalExit
	}
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		rotation += math.Pi / 100
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		rotation -= math.Pi / 100
	}

	rotation = math.Atan2(math.Sin(rotation), math.Cos(rotation))

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	cam.Surface.Clear()

	// face logic
	x, y := -float64(spriteSize*spriteScale)/2, -float64(spriteSize*spriteScale)/2
	sideLength := float64(spriteSize * spriteScale)
	tl := NewPoint(x, y).Rotate(rotation)
	tr := NewPoint(x+float64(spriteScale*spriteScale)*2, y).Rotate(rotation)
	bl := NewPoint(x, y+float64(spriteScale*spriteScale)*2).Rotate(rotation)
	br := NewPoint(x+float64(spriteScale*spriteScale)*2, y+float64(spriteScale*spriteScale)*2).Rotate(rotation)

	op := &ebiten.DrawImageOptions{}

	// bottom face of the cube (for debug, but this can also be used to draw floor tiles)
	// op.ColorM.Scale(1, 1, 1, 0.2)
	// op = camera.GetRotation(op, rotation, -float64(spriteSize)/2, -float64(spriteSize)/2)
	// op = camera.GetScale(op, 8, 8)
	// op = camera.GetTranslation(op, x, y)
	// camera.Surface.DrawImage(
	// 	sprites.SubImage(
	// 		image.Rect(0, 0, spriteSize, spriteSize)).(*ebiten.Image),
	// 	op)

	// draw faces clockwise to prevent image flipping
	drawFace := func(p1, p2 *Point, rect image.Rectangle) {
		op = &ebiten.DrawImageOptions{}
		op.ColorM.Scale(1, 1, 1, 0.5)
		op = cam.GetScale(op, 8*(p2.X-p1.X)/sideLength, 8)
		op = cam.GetSkew(op, 0, p1.AngleTo(p2))
		op = cam.GetTranslation(op, p1.X, p1.Y-float64(spriteScale*spriteSize))
		cam.Surface.DrawImage(
			sprites.SubImage(rect).(*ebiten.Image),
			op)
	}

	if math.Abs(rotation) <= math.Pi/2 {
		drawFace(bl, br, rectBlue)
	}
	if math.Abs(rotation) >= math.Pi/2 {
		drawFace(tr, tl, rectBlue)
	}
	if rotation >= -math.Pi && rotation < 0 {
		drawFace(tl, bl, rectGreen)
	}
	if rotation > 0 && rotation <= math.Pi {
		drawFace(br, tr, rectYellow)
	}

	// top face of the cube
	op = &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	op = cam.GetRotation(op, rotation, -float64(spriteSize)/2, -float64(spriteSize)/2)
	op = cam.GetScale(op, 8, 8)
	op = cam.GetTranslation(op, x, y-float64(spriteScale*spriteSize))
	cam.Surface.DrawImage(
		sprites.SubImage(
			image.Rect(0, 0, spriteSize, spriteSize)).(*ebiten.Image),
		op)

	cam.Blit(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Rotation: %0.1f", rotation))
}

// Layout sets window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if LastWindowWidth != outsideWidth || LastWindowHeight != outsideHeight {
		cam.Resize(outsideWidth, outsideHeight)
		log.Println("resize", outsideWidth, outsideHeight)
		LastWindowWidth = outsideWidth
		LastWindowHeight = outsideHeight
	}
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}
	w, h := 640*2, 480*2
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Spinning cube example")
	ebiten.SetWindowResizable(true)
	cam = camera.NewCamera(w, h, 0, 0, 0, 1)

	if b, err := embedded.ReadFile("sprites.png"); err == nil {
		if s, err := png.Decode(bytes.NewReader(b)); err == nil {
			sprites = ebiten.NewImageFromImage(s)
		}
	} else {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
