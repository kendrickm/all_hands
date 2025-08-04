package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) DrawMenu() {
	// ui.renderer.Clear()
	menuRect := ui.GetMenuRect()
	ui.renderer.Copy(ui.menuBackground, nil, menuRect)

	//Render buttons

}

func (ui *ui) GetMenuRect() *sdl.Rect {
	terWidth  := int32(float32(ui.winWidth)*0.40)
	terHeight := int32(float32(ui.winHeight)*0.75)
	offsetX := (int32(ui.winWidth) - terWidth) / 2
	offsetY := (int32(ui.winHeight) - terHeight) / 2
	return &sdl.Rect{X:offsetX, Y: offsetY, W: terWidth, H: terHeight} 
}