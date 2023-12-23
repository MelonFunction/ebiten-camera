package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)
	ebiten.SetWindowSize(512, 512)
	// ebiten.SetTPS(4)

	game := &Game{
		p1:    NewEntity(64, 64, 4),
		enemy: NewEntity(32, 32, 4),
		cam:   camera.NewCamera(512, 512, 0, 0, 0, 1, &ebiten.DrawImageOptions{}),
	}

	game.cam.CamDrawOps.Filter = ebiten.FilterNearest
	// game.cam.Scale = 1.5
	game.spawnPositions = RandomPoints(-1000, 1000, 50)
	game.enemy.Img.Fill(color.RGBA{255, 0, 0, 255})
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

type Entity struct {
	X, Y, Width, Height float64
	Rotation            float64
	Img                 *ebiten.Image
	DrawOptions         *ebiten.DrawImageOptions
	Speed               float64
}

func NewEntity(w, h int, speed float64) *Entity {
	o := &Entity{
		Width:       float64(w),
		Height:      float64(h),
		Speed:       speed,
		Rotation:    0.0,
		Img:         ebiten.NewImage(w, h),
		DrawOptions: &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear},
	}
	o.Img.Fill(color.Gray{128})
	return o
}

func (e *Entity) GetCenter() (x, y float64) {
	return e.X + (e.Width * 0.5), e.Y + (e.Height * 0.5)
}

type Game struct {
	cam            *camera.Camera
	p1             *Entity
	enemy          *Entity
	spawnPositions []Point
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cam.SetZoom(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cam.SetZoom(2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.p1.X += g.p1.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.p1.X -= g.p1.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.p1.Y -= g.p1.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.p1.Y += g.p1.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.p1.Rotation += 0.1
	}

	g.cam.SetPosition(g.p1.GetCenter())
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.cam.Surface.Clear()
	// draw enemies
	for _, spawnPoint := range g.spawnPositions {
		g.enemy.DrawOptions.GeoM.Reset()
		g.enemy.DrawOptions = g.cam.GetTranslation(g.enemy.DrawOptions, spawnPoint.X, spawnPoint.Y)
		g.cam.Surface.DrawImage(g.enemy.Img, g.enemy.DrawOptions)
	}
	// reset
	g.p1.DrawOptions.GeoM.Reset()

	// locally rotate the player from center
	g.p1.DrawOptions.GeoM.Translate(-32, -32)
	g.p1.DrawOptions.GeoM.Rotate(g.p1.Rotation)
	g.p1.DrawOptions.GeoM.Translate(32, 32)

	// camera
	g.p1.DrawOptions = g.cam.GetTranslation(g.p1.DrawOptions, g.p1.X, g.p1.Y)
	g.cam.Surface.DrawImage(g.p1.Img, g.p1.DrawOptions)
	g.cam.Blit(screen)
}
func (g *Game) Layout(w, h int) (screenWidth, screenHeight int) {
	return 512, 512
}

func RandomPoints(min, max float64, n int) []Point {
	points := make([]Point, n)
	for i := range points {
		points[i] = Point{X: min + rand.Float64()*(max-min),
			Y: min + rand.Float64()*(max-min)}
	}
	return points
}

type Point struct {
	X, Y float64
}
