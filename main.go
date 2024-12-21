package main

import (
	"github.com/kendrickm/all_hands/game"
	"github.com/kendrickm/all_hands/ui"
)

func main() {
	// TODO When we need multiple UI Support refactor event polling to it's own component
	// and run only on main thread
	g := game.NewGame()
	sm := game.NewStateMachine(g)
	go func() {
		g.Run(sm)

	}()
	ui := ui.NewUI(g.InputChan, g.CurrentRoom, g.GameStateChan)
	ui.Run()

}
