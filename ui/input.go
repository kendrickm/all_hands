package ui

import (
	"github.com/kendrickm/all_hands/game"
	"github.com/veandco/go-sdl2/sdl"
)

type mouseState struct {
	leftButton  bool
	rightButton bool
	pos        game.Pos
}

func getMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()

	var result mouseState
	result.pos = game.Pos{int(mouseX), int(mouseY)}
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return &result
}


func (ui *ui) keyDownOnce(key uint8) bool {
	return ui.keyboardState[key] == 1 && ui.preKeyboardState[key] == 0
}

// Check for key pressed and then released
func (ui *ui) keyPressed(key uint8) bool {
	return ui.keyboardState[key] == 0 && ui.preKeyboardState[key] == 1
}