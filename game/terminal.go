package game
import (
	"strconv"
)

type ButtonAction int
const (
	UnusedButton ButtonAction = iota
	ReactorOn
	ReactorOff
)

type ButtonType int 
const (
	ToggleSwitch ButtonType = iota //Switch can do different functions in up or down
	TriggerButton //Button that can be pressed only once
	ToggleButton //Button that can be pressed multiple times
)

type Button struct {
	Type ButtonType
	state bool
	DisplayText string
}

func (button *Button) PressButton() {
	button.state = true
	// fmt.Println("Pressed button")
}

type TerminalType int
const(
	SINGLE_BUTTON TerminalType = iota
	GRAPH_DISPLAY
)

type Terminal struct {
	Powered bool
	Buttons []*Button
	Name string
	LinkedStation *Station
	DisplayText string
	Type TerminalType
}

//TODO Will need to refactor this once multi button terminals are introduced
type TerminalStatus struct {
	Powered bool
	ButtonState bool
	ButtonActive bool
	Info string
}

//TODO Will need to refactor this once multi button terminals are introduced
func (t *Terminal) GetCurrentState() *TerminalStatus {
	active := true
	if t.Type == SINGLE_BUTTON {
		//Trigger buttons can only be pressed once so if its pressed then lock it
		if t.Buttons[0].Type == TriggerButton && t.Buttons[0].state == true { 
			active = false
		}
		return &TerminalStatus{t.Powered, t.Buttons[0].state, active, ""}
	} else if t.Type == GRAPH_DISPLAY {
		info := t.LinkedStation.Level
		return &TerminalStatus{t.Powered, false, false, strconv.Itoa(int(info))}
	}
	panic("Terminal type not implemented yet")
	return nil

}

func singleButtonTerminalFactory() *Terminal {
	t := &Terminal{}
	t.Powered = true
	t.Buttons = make([]*Button, 1)
	t.Name = ""
	t.DisplayText = ""
	t.Type = SINGLE_BUTTON

	return t
}

func graphDisplayTerminalFactor() *Terminal {
	t := &Terminal{}
	t.Powered = true
	t.Buttons = nil
	t.Name = ""
	t.DisplayText = ""
	t.Type = GRAPH_DISPLAY

	return t
}