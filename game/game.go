package game

import (
	"fmt"
)

type Game struct {
	RoomChans   []chan *Room
	InputChan   chan *Input
	Ships       []*Ship
	CurrentShip *Ship
	CurrentRoom *Room
	Chapter     int32 //Index of ship
}

func NewGame() *Game {
	//TODO Create channels for event management

	roomChans := make([]chan *Room, 1)
	for i := range roomChans {
		roomChans[i] = make(chan *Room)
	}

	inputChan := make(chan *Input)
	ships := loadChapter0(newPlayer())

	game := &Game{roomChans, inputChan, ships, ships[0], ships[0].Rooms[0], 0}
	return game

}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
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

func (game *Game) Move(to Pos) {
	game.CurrentRoom.Player.Pos = to
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
		case Bulkhead, Blank:
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
	for _, r := range game.RoomChans {
		r <- game.CurrentRoom
	}

	for input := range game.InputChan {
		if input.Typ == QuitGame {
			fmt.Println("Quitting")
			return
		}

		game.handleInput(input)

	}

}
