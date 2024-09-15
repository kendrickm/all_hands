package game

type ButtonAction int

const (
	UnusedButton ButtonAction = iota
	ReactorOn
	ReactorOff
)

type ButtonType int 
const (
	ToggleSwitch ButtonType = iota
	TriggerButton
	ToggleButton
)

type Button struct {
	Type ButtonType
	state bool
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
}

func (button *Button) PressButton() {
	button.state = true
	// fmt.Println("Pressed button")
}


func singleButtonTerminalFactory() *Terminal {
	t := &Terminal{}
	t.Powered = true
	t.Buttons = make([]*Button, 1)
	b := &Button{TriggerButton, false}
	t.Buttons[0] = b
	t.Name = ""

	return t
}