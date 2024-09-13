//TODO 
//Event management
//Ship status
//** Reactor status
//** Aux Power Banks

package game

import (
	"fmt"
)

type Game struct {
	//Is there a better way to handle different game states communicating between game and ui threads?
	GameStateChan chan *StateChange //True when a terminal is active
	InputChan   chan *Input
	Ships       []*Ship
	CurrentShip *Ship
	CurrentRoom *Room
	Chapter     int32 //Index of ship
	ActiveTerminal *Terminal
}

func NewGame() *Game {
	//TODO Create channels for event management

	inputChan := make(chan *Input)
	stateChan := make(chan *StateChange)
	ships := loadChapter0(newPlayer())
	var ct *Terminal
	game := &Game{stateChan, inputChan, ships, ships[0], ships[0].Rooms[0], 0, ct}
	if ct == nil {
		fmt.Println("Starting with nil")
	}
	return game

}

type StateChange struct {
	TerminalActive bool
	Terminal *Terminal
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	TerminalInteract
	QuitGame
	CloseWindow
)

type Input struct {
	Typ InputType
}

type Pos struct {
	X, Y int
}

func checkDoor(room *Room, pos Pos) {
	t := room.Map[pos.Y][pos.X]
	if t.OverlayRune == ClosedDoor {
		room.Map[pos.Y][pos.X].OverlayRune = OpenDoor
	}
}

func checkTerminal(room *Room, playerPos Pos) *Terminal { //Since we aren't changing facing direction, check each direction currently
	nPos := Pos{playerPos.X,playerPos.Y-1}
	sPos := Pos{playerPos.X,playerPos.Y+1}
	ePos := Pos{playerPos.X+1,playerPos.Y}
	wPos := Pos{playerPos.X-1,playerPos.Y}
	var terminal *Terminal
	fmt.Println(room.Terminals)

	if room.Map[nPos.Y][nPos.X].Rune == TerminalAccess {
		terminal = room.Terminals[nPos]
		fmt.Println("Found Terminal north")
		fmt.Println(nPos)
	} else if room.Map[sPos.Y][sPos.X].Rune == TerminalAccess {
		terminal = room.Terminals[sPos]
		fmt.Println("Found Terminal south")
	} else if room.Map[ePos.Y][ePos.X].Rune == TerminalAccess {
		terminal = room.Terminals[ePos]
		fmt.Println("Found Terminal east")
	} else if room.Map[wPos.Y][wPos.X].Rune == TerminalAccess {
		terminal = room.Terminals[wPos]
		fmt.Println("Found Terminal west")
	}

	return terminal
}

func (game *Game) Move(to Pos) {
	game.CurrentRoom.Player.Pos = to
	fmt.Println("Moving")
}

func (game *Game) resolveMovement(pos Pos) {
	room := game.CurrentRoom
	if canWalk(room, pos) {
		game.Move(pos)
	} else {
		checkDoor(room, pos)
	}


}

func canWalk(room *Room, pos Pos) bool {
	if inRange(room, pos) {
		t := room.Map[pos.Y][pos.X]
		switch t.Rune {
		case TerminalAccess,PoweredReactor,UnpoweredReactor,Bulkhead, Blank:
			return false
		}
		return true
	}
	return false
}

func inRange(room *Room, pos Pos) bool {
	return pos.X < len(room.Map[0]) && pos.Y < len(room.Map) && pos.X >= 0 && pos.Y >= 0
}

func (game *Game) handleInput(input *Input) {
	room := game.CurrentRoom
	p := room.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.resolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.resolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.resolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.resolveMovement(newPos)
	case CloseWindow:
		//Handle closing terminals here
	case TerminalInteract:
		if game.ActiveTerminal == nil {
			t := checkTerminal(room, p.Pos)
			if t != nil{
				game.ActiveTerminal = t
				fmt.Println("Setting active terminal")
			}
		} else {
			game.ActiveTerminal = nil
			fmt.Println("Unsetting terminal")
		}
	}
}

func getNeighbors(room *Room, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if canWalk(room, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(room, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(room, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(room, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func (game *Game) Run() {
	// for _, r := range game.RoomChans {
	// 	r <- game.CurrentRoom
	// }

	for input := range game.InputChan {
		if input.Typ == QuitGame {
			fmt.Println("Quitting")
			return
		}
		game.handleInput(input)
		state := &StateChange{}
		if game.ActiveTerminal != nil {
			state.TerminalActive = true
			state.Terminal = game.ActiveTerminal
		}
		game.GameStateChan <- state
		game.CurrentRoom.Update()
	}

	

}
