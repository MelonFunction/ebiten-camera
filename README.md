# ebiten-camera

A simple camera implementation based on [vrld's hump for LÖVE](https://github.com/vrld/hump)

Look at [cmd](https://github.com/melonfunction/ebiten-camera/tree/master/cmd) to see a very basic implementation.
Understand that the code is terrible because I wanted to keep the logic as simple as possible 😅

## Usage 

[📖 Docs](https://pkg.go.dev/github.com/melonfunction/ebiten-camera)

Here is a stripped-down version of the example code, highlighting the most important functions. 

It won't run, so please don't bother trying it 😄

```go

var (
    cam *ebitenCamera.Camera
    // other vars excluded, such as tiles, PlayerX etc
)


func main() {
    w,h := 640, 480
    // excluding normal ebiten setup
    cam = camera.NewCamera(w, h, 0, 0, 0, 1)
}

func (g *Game) Update() error {
    // Follows the player
    cam.SetPosition(PlayerX+float64(PlayerSize)/2, PlayerY+float64(PlayerSize)/2)
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Clear camera surface
    cam.Surface.Clear()
    cam.Surface.Fill(color.RGBA{255, 128, 128, 255})
    // Draw tiles
    tileOps := &ebiten.DrawImageOptions{}
    cam.Surface.DrawImage(tiles, cam.GetTranslation(tileOps, 0, 0))
    // Draw the player
    playerOps := &ebiten.DrawImageOptions{}
    playerOps = cam.GetRotation(playerOps, PlayerRot, -float64(PlayerSize)/2, -float64(PlayerSize)/2)
    playerOps = cam.GetScale(playerOps, 0.5, 0.5)
    playerOps = cam.GetTranslation(playerOps, PlayerX, PlayerY)
    cam.Surface.DrawImage(player, playerOps)
    
    // Draw to screen and zoom
    cam.Blit(screen)
}
```
## Considerations 

1) When setting the `*ebiten.DrawImageOptions` in the `ebitenCamera.Surface.DrawImage` function, the order of **operation is 
important**!  
**Rotate**, **Scale** and then **Translate**!

2) My example doesn't include a range for the camera's zoom to highlight what happens if you zoom too far in either direction. This is because I use a render texture to draw everything to, and I simply resize this when zooming in or out. This texture has a max size, causing the positioning logic to stop centering it. I'll probably change this at some point (unless you beat me to it with a PR, of course 😆)
