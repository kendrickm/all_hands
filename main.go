package main

import (
	"github.com/kendrickm/all_hands/game"
	"github.com/kendrickm/all_hands/ui"
)

func main() {
	// TODO When we need multiple UI Support refactor event polling to it's own component
	// and run only on main thread
	game := game.NewGame()
	go func() {
		game.Run()

	}()
	ui := ui.NewUI(game.InputChan, game.RoomChans[0])
	ui.Run()

}
