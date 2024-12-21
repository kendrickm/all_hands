package game

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
	OUTPUT_GRAPH
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
}

//TODO Will need to refactor this once multi button terminals are introduced
func (t *Terminal) GetCurrentState() *TerminalStatus {
	active := true
	//Trigger buttons can only be pressed once so if its pressed then lock it
	if t.Buttons[0].Type == TriggerButton && t.Buttons[0].state == true { 
		active = false
	}
	return &TerminalStatus{t.Powered, t.Buttons[0].state, active}
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