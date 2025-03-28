package game

import (
	"fmt"
)

type State interface {
	Enter()
	Exit()
	Update(sm *StateMachine, input *Input)
}

type GameStart struct{
	game *Game
}
func (g *GameStart) Enter(){
	fmt.Println("Welcome to the game!")
}
func (g *GameStart) Exit(){}
func (g *GameStart) Update(sm *StateMachine, input *Input){
	sm.setState(&MainGame{g.game})
}

type MainGame struct{
	game *Game
	// events []*Events
}
func (g *MainGame) Enter(){
	fmt.Println("Here we go!")
	g.game.stateChange(g)
}
func (g *MainGame) Exit(){
}
func (g *MainGame) Update(sm *StateMachine, input *Input){
	if input == nil { // No input
		return
	}
	if input.Typ == QuitGame {
		sm.setState(&GameOver{})
		return
	}
	g.game.handleInput(input)
	g.game.CurrentRoom.Update()
	if g.game.ActiveTerminal != nil {
		sm.setState(&TerminalState{g.game})
	}
}

type TerminalState struct{
	game *Game
}
func (g *TerminalState) Enter(){
	g.game.stateChange(g)
}
func (g *TerminalState) Exit(){}
func (g *TerminalState) Update(sm *StateMachine, input *Input){
	if g.game.ActiveTerminal == nil{
		sm.setState(&MainGame{g.game})
	}
	if input.Typ == TerminalInteract {
		g.game.ActiveTerminal = nil
	}
	g.game.CurrentRoom.Update()
}


type GameOver struct{
	game *Game
}
func (g *GameOver) Enter(){
	fmt.Println("Quitting")
}
func (g *GameOver) Exit(){}
func (g *GameOver) Update(sm *StateMachine, input *Input){

}


type MsgRecieved struct{
	game *Game
}
func (g *MsgRecieved) Enter(){}
func (g *MsgRecieved) Exit(){}
func (g *MsgRecieved) Update(sm *StateMachine, input *Input){

}

type StateMachine struct {
    currentState State
    states       map[string]State
}

func (sm *StateMachine) setState(s State) {
	if sm.currentState != nil {
		sm.currentState.Exit()
	}
	sm.currentState = s
  	sm.currentState.Enter()
}

func (sm *StateMachine) Update(input *Input) {
	sm.currentState.Update(sm, input)

	if input != nil && input.handled{
		input = nil
	}
}

func NewStateMachine(g *Game) *StateMachine {
	
	sm := &StateMachine{
		currentState: nil,
		states:       make(map[string]State),
	}
	sm.setState(&GameStart{g})
	return sm
}
