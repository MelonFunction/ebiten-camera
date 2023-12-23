// Package camera provides a simple camera system for use with ebiten
package camera

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera can look at positions, zoom and rotate.
type Camera struct {
	X, Y, Rot, Scale float64
	Width, Height    int
	Surface          *ebiten.Image
	CamDrawOps       *ebiten.DrawImageOptions
}

// NewCamera returns a new Camera
func NewCamera(width, height int, x, y, rotation, zoom float64, op *ebiten.DrawImageOptions) *Camera {
	return &Camera{
		X:          x,
		Y:          y,
		Width:      width,
		Height:     height,
		Rot:        rotation,
		Scale:      zoom,
		Surface:    ebiten.NewImage(width, height),
		CamDrawOps: op,
	}
}

// SetPosition looks at a position
func (c *Camera) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
}

// Translate translates the Camera by x and y.
// Use SetPosition if you want to set the position
func (c *Camera) Translate(x, y float64) {
	c.X += x
	c.Y += y
}

// Rotate rotates by phi
func (c *Camera) Rotate(phi float64) {
	c.Rot += phi
}

// SetRotation sets the rotation to rot
func (c *Camera) SetRotation(rot float64) {
	c.Rot = rot
}

// Zoom *= the current zoom
func (c *Camera) Zoom(mul float64) *Camera {
	c.Scale *= mul
	if c.Scale <= 0.01 {
		c.Scale = 0.01
	}
	c.Resize(c.Width, c.Height)
	return c
}

// SetZoom sets the zoom
func (c *Camera) SetZoom(zoom float64) {
	c.Translate(float64(c.Width)/2, float64(c.Width)/2)
	c.Scale = zoom
	if c.Scale <= 0.01 {
		c.Scale = 0.01
	}
	c.Resize(c.Width, c.Height)
}

// Resize resizes the camera Surface
func (c *Camera) Resize(w, h int) *Camera {
	c.Width = w
	c.Height = h
	newW := int(float64(w) * 1.0 / c.Scale)
	newH := int(float64(h) * 1.0 / c.Scale)
	if newW <= 16384 && newH <= 16384 {
		c.Surface.Dispose()
		c.Surface = ebiten.NewImage(newW, newH)
	}
	return c
}

// GetTranslation alters the provided *ebiten.DrawImageOptions' translation based on the given x,y offset and the
// camera's position
func (c *Camera) GetTranslation(ops *ebiten.DrawImageOptions, x, y float64) *ebiten.DrawImageOptions {
	ops.GeoM.Translate(float64(c.Width)/2.0, float64(c.Height)/2.0)
	ops.GeoM.Translate(-c.X+x, -c.Y+y)
	return ops
}

// GetRotation alters the provided *ebiten.DrawImageOptions' rotation using the provided originX and originY args
func (c *Camera) GetRotation(ops *ebiten.DrawImageOptions, rot, originX, originY float64) *ebiten.DrawImageOptions {
	ops.GeoM.Translate(originX, originY)
	ops.GeoM.Rotate(rot)
	ops.GeoM.Translate(-originX, -originY)
	return ops
}

// GetScale alters the provided *ebiten.DrawImageOptions' scale
func (c *Camera) GetScale(ops *ebiten.DrawImageOptions, scaleX, scaleY float64) *ebiten.DrawImageOptions {

	ops.GeoM.Scale(scaleX, scaleY)
	return ops
}

// GetSkew alters the provided *ebiten.DrawImageOptions' skew
func (c *Camera) GetSkew(ops *ebiten.DrawImageOptions, skewX, skewY float64) *ebiten.DrawImageOptions {
	ops.GeoM.Skew(skewX, skewY)
	return ops
}

// Blit draws the camera's surface to the screen and applies zoom
func (c *Camera) Blit(screen *ebiten.Image) {
	centerX := float64(c.Width) / 2.0
	centerY := float64(c.Height) / 2.0
	c.CamDrawOps.GeoM.Reset()
	c.CamDrawOps.GeoM.Translate(-centerX, -centerY)
	c.CamDrawOps.GeoM.Scale(c.Scale, c.Scale)
	c.CamDrawOps.GeoM.Rotate(c.Rot)
	c.CamDrawOps.GeoM.Translate(centerX*c.Scale, centerY*c.Scale)
	screen.DrawImage(c.Surface, c.CamDrawOps)
}

// WorldToScreenCoords converts world coords into screen coords
func (c *Camera) WorldToScreenCoords(worldX, worldY float64) (float64, float64) {
	co, si := math.Cos(-c.Rot), math.Sin(-c.Rot)
	worldX, worldY = worldX-c.X, worldY-c.Y
	worldX, worldY = co*worldX-si*worldY, si*worldX+co*worldY
	return worldX*c.Scale + float64(c.Width)/2.0, worldY*c.Scale + float64(c.Height)/2.0
}

// ScreenToWorldCoords converts screen coords into world coords
func (c *Camera) ScreenToWorldCoords(screenX, screenY float64) (float64, float64) {
	co, si := math.Cos(-c.Rot), math.Sin(-c.Rot)
	screenX, screenY = (screenX-float64(c.Width)/2.0)/c.Scale, (screenY-float64(c.Height)/2.0)/c.Scale
	screenX, screenY = co*screenX-si*screenY, si*screenX+co*screenY
	return screenX + c.X, screenY + c.Y
}
