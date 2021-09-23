package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	camera "github.com/scarycoffee/ebiten-camera"
)

// Vars
var (
	cam    *camera.Camera
	tiles  *ebiten.Image
	player *ebiten.Image

	// When keyF is pressed, change follow mode (implemented by using
	// cam.SetPosition() or cam.MovePosition())
	CamFollowPlayer = true

	LastWindowWidth  int
	LastWindowHeight int

	LastMouseX      = 0
	LastMouseY      = 0
	MouseWasDown    bool
	MouseDownStartX int
	MouseDownStartY int
	MouseDownAt     time.Time
	MousePanAfter   = time.Millisecond * 100

	ErrNormalExit = errors.New("Normal exit")

	PlayerSize  int
	TileSize    int
	LevelWidth  = 30
	LevelHeight = 5
	Level       = []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}

	PlayerX float64 = 400.0
	PlayerY float64
	VelX    float64
	VelY    float64
	Gravity = 5.0
	JumpVel = -40.0
	Jumping = false

	mx, my float64
	px, py int
)

// Game required by ebiten
type Game struct{}

// Update updates the Game
func (g *Game) Update() error {

	VelX = 0
	if !Jumping {
		VelY = 0
	}

	// Keyboard controls
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		CamFollowPlayer = !CamFollowPlayer
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyH) {
		VelX = -5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyN) {
		VelX = 5
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !Jumping {
			VelY = JumpVel
			Jumping = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrNormalExit
	}

	// Physics
	if TileSize != 0 {
		VelY += Gravity

		// Cursor tile position
		mx, my = cam.GetMouseCoords()
		my = float64((int(my)) / int(TileSize))
		mx = float64((int(mx)) / int(TileSize))

		// Player tile position
		px = (int(PlayerX) + PlayerSize/2) / int(TileSize)
		py = (int(PlayerY) + PlayerSize) / int(TileSize)

		// Absolutely terrible collision detection and physics 🤫
		index := py*LevelWidth + px
		if index < LevelWidth*LevelHeight && index >= 0 {
			if Level[index] == 1 {
				// Touching tile
				Jumping = false
				VelY -= Gravity
				PlayerY = float64(py*TileSize) - float64(PlayerSize)
			}
		}

		PlayerX += VelX
		PlayerY += VelY

		if CamFollowPlayer {
			cam.SetPosition(PlayerX+float64(PlayerSize)/2, PlayerY+float64(PlayerSize)/2)
		}
	}

	// Panning, setting up for click events
	cx, cy := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !MouseWasDown {
			// First frame mouse is down
			MouseWasDown = true
			MouseDownStartX = cx
			MouseDownStartY = cy
			MouseDownAt = time.Now()
		} else {
			// Pan when pressed for long enough
			if time.Now().Sub(MouseDownAt) > MousePanAfter && !CamFollowPlayer {
				cam.MovePosition(
					float64(LastMouseX)-float64(cx),
					float64(LastMouseY)-float64(cy))
			}
		}

	} else if MouseWasDown {
		MouseWasDown = false
		// Only call mouse up event if the cursor didn't move more than a certain amount
		triggerMoveAmount := 1.0
		if math.Abs(float64(MouseDownStartX-cx)) < triggerMoveAmount && math.Abs(float64(MouseDownStartY-cy)) < triggerMoveAmount {

		}
	}

	LastMouseX = cx
	LastMouseY = cy

	// Zoom
	_, scrollAmount := ebiten.Wheel()
	if scrollAmount > 0 {
		cam.Zoom(1.1)
	} else if scrollAmount < 0 {
		cam.Zoom(0.9)
	}

	return nil
}

// Draw renders everything to screen
func (g *Game) Draw(screen *ebiten.Image) {

	w, _ := ebiten.WindowSize()
	if tiles == nil {
		// Draw tiles image
		TileSize = w / 10
		tiles = ebiten.NewImage(TileSize*LevelWidth, TileSize*LevelHeight)
		PlayerSize = int(float64(TileSize) * 0.75)

		for y := 0; y < LevelHeight; y++ {
			for x := 0; x < LevelWidth; x++ {
				switch Level[y*LevelWidth+x] {
				case 0:
				case 1:
					ebitenutil.DrawRect(
						tiles,
						float64(x*TileSize),
						float64(y*TileSize),
						float64(TileSize),
						float64(TileSize),
						color.RGBA{0, 255, 0, 255})
				}
			}
		}
	}
	if player == nil {
		// Draw player image
		player = ebiten.NewImage(PlayerSize, PlayerSize)
		player.Fill(color.RGBA{128, 0, 128, 255})
	}

	// Clear camera surface
	cam.Surface.Clear()
	// Draw tiles
	cam.Surface.DrawImage(tiles, cam.GetTranslation(0, 0))
	// Draw the player
	cam.Surface.DrawImage(player, cam.GetTranslation(PlayerX, PlayerY))

	// Draw to screen and zoom
	cam.Blit(screen)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"Camera:\n  X: %3.3f\n  Y: %3.3f\n  Rot: %3.3f\n  Zoom: %3.3f\n"+
				"Tiles:\n  PlayerX: %d\n  PlayerY: %d\n  MouseX: %1.0f\n  MouseY: %1.0f",
			cam.X, cam.Y, cam.Rot, cam.Scale,
			px, py, mx, my,
		))
}

// Layout sets window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if LastWindowWidth != outsideWidth || LastWindowHeight != outsideHeight {
		cam.Resize(outsideWidth, outsideHeight)
		LastWindowWidth = outsideWidth
		LastWindowHeight = outsideHeight
	}
	return outsideWidth, outsideHeight
}

func main() {
	log.SetFlags(log.Lshortfile)

	ebiten.SetWindowSize(640*2, 480*2)
	ebiten.SetWindowTitle("Camera Test")
	ebiten.SetWindowResizable(true)

	cam = camera.NewCamera(0, 0, 0, 1)

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}