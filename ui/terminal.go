package ui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/kendrickm/all_hands/game"
	"fmt"
)

func (ui *ui) DrawTerminal(terminal *game.Terminal) {
	switch terminal.Type {
	case game.SINGLE_BUTTON:
		ui.DrawSBTerminal(terminal)
	default:
		fmt.Println("Unknow Type")
		panic(terminal.Type)
	}
}

func (ui *ui) checkButton(buttonRect *sdl.Rect) bool {
	mousePos := ui.currentMouseState.pos
	return buttonRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y),int32(1),int32(1)})
}

func (ui *ui) DrawSBTerminal(terminal *game.Terminal) { // For drawing a simple, single button terminal
	// fmt.Println("Drawing terminal " + terminal.Name)
	ui.renderer.Clear()
	tState := terminal.GetCurrentState()
	terminalRect := ui.getTerminalRect()
	insetRect := getInsetRect(terminalRect)
	buttonRect := getButtonRect(insetRect)
	displayTextRect := getDisplayTextRect(insetRect)
	nameTexture := ui.stringToTexture(terminal.DisplayText, FontSmall,sdl.Color{255, 255, 255, 255})
	buttonPressed := tState.ButtonState
	if tState.ButtonActive {
		if ui.currentMouseState.leftButton && !ui.prevMouseState.leftButton{
			terminal.Buttons[0].PressButton()
		}
	}

	buttonTexture := ui.buttonTexture
	if buttonPressed {
		buttonTexture = ui.buttonTexturePressed
	}


	ui.renderer.Copy(ui.terminalBackground, nil, terminalRect)
	ui.renderer.Copy(ui.terminalForeground, nil, insetRect)
	ui.renderer.Copy(buttonTexture, nil, buttonRect)
	ui.renderer.Copy(ui.terminalTextboxTexture, nil, displayTextRect)
	ui.renderer.Copy(nameTexture, nil, displayTextRect)
}

func getInsetRect(outerTerminal *sdl.Rect) *sdl.Rect {
	terWidth  := int32(float32(outerTerminal.W) * 0.91)
	terHeight := int32(float32(outerTerminal.H) * 0.85)
	offsetX := outerTerminal.X+((outerTerminal.W - terWidth)/2)
	offsetY := outerTerminal.Y+((outerTerminal.H - terHeight)/2)
	return &sdl.Rect{X:offsetX, Y: offsetY, W: terWidth, H: terHeight} 
}


func (ui *ui) getTerminalRect() *sdl.Rect {
	terWidth  := int32(float32(ui.winWidth)*0.40)
	terHeight := int32(float32(ui.winHeight)*0.75)
	offsetX := (int32(ui.winWidth) - terWidth) / 2
	offsetY := (int32(ui.winHeight) - terHeight) / 2
	return &sdl.Rect{X:offsetX, Y: offsetY, W: terWidth, H: terHeight} 
}

func getButtonRect(insetTerminal *sdl.Rect) *sdl.Rect {
	terWidth  := int32(float32(insetTerminal.W) * 0.1)
	terHeight := int32(float32(insetTerminal.H) * 0.1)
	offsetX := insetTerminal.X+int32(float32(insetTerminal.W)*0.5) - terWidth/2
	offsetY := insetTerminal.Y+int32(float32(insetTerminal.H)*0.5) - terHeight/2
	return &sdl.Rect{X:offsetX, Y: offsetY, W: terWidth, H: terHeight} 
}

func getDisplayTextRect(insetTerminal *sdl.Rect) *sdl.Rect {
	terWidth  := int32(float32(insetTerminal.W) * 0.3)
	terHeight := int32(float32(insetTerminal.H) * 0.05)
	offsetX := insetTerminal.X+int32(float32(insetTerminal.W) * 0.02)
	offsetY := insetTerminal.Y+insetTerminal.H-terHeight-int32(float32(insetTerminal.H) * 0.02)
	return &sdl.Rect{X:offsetX, Y: offsetY, W: terWidth, H: terHeight} 
}