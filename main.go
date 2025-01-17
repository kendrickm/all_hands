package main

import (
	"github.com/kendrickm/all_hands/game"
	"github.com/kendrickm/all_hands/ui"
	"time"
)

func main() {
	// TODO When we need multiple UI Support refactor event polling to it's own component
	// and run only on main thread
	g := game.NewGame()
	sm := game.NewStateMachine(g)
	ticker := time.NewTicker(16 * time.Millisecond)

	go func() {
		// for t := range ticker.C {
		// 	fmt.Println("Tick at", t)
        //     g.Run(sm,ticker)
        // }
        g.Run(sm,ticker)
	}()
	ui := ui.NewUI(g.InputChan, g.CurrentRoom, g.GameStateChan)
	ui.Run()

}
