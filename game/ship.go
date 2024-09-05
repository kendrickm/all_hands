package game

import (
		"path/filepath"
		"strings"
		"fmt"
		"os"
		"bufio"
)

type Ship struct {
	Rooms []*Room
}

func (ship *Ship) Update() {
	for _, room := range ship.Rooms {
		room.Update()
	}
}

func loadChapter0(player *Player) []*Ship {
	// TODO refactor the depreacted function
	// TODO abstract loading that will span chapters
	// TODO load series of rooms from a single ship file
	// TODO load a series of ships from chapter logic
	filenames, err := filepath.Glob("game/rooms/*.room")
	if err != nil {
		panic(err)
	}

	var rooms []*Room
	var ships []*Ship

	for _, filename := range filenames {

		roomName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		fmt.Println("loading:", roomName)
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		roomLines := make([]string, 0)
		longestRow := 0
		index := 0
		for scanner.Scan() {
			roomLines = append(roomLines, scanner.Text())
			if len(roomLines[index]) > longestRow {
				longestRow = len(roomLines[index])
			}
			index++
		}

		room := &Room{}
		room.Debug = make(map[Pos]bool)
		room.Player = player

		room.Map = make([][]Tile, len(roomLines))
		room.Terminals = make(map[Pos]*Terminal, 1)

		for i := range room.Map {
			room.Map[i] = make([]Tile, longestRow)
		}

		for y := 0; y < len(room.Map); y++ {
			line := roomLines[y]

			for x, c := range line {
				var t Tile
				t.OverlayRune = Blank
				switch c {
				case ' ', '\n', '\t', '\r':
					t.Rune = Blank
				case '#':
					t.Rune = Bulkhead
				case '.':
					t.Rune = ShipFloor
				case 'T': 
					 t.Rune = TerminalAccess
					 room.Terminals[Pos{x,y}] = createReactorTerminal()
			    case 'r':
			    	 t.Rune = UnpoweredReactor
				case '@':
					room.Player.X = x
					room.Player.Y = y
					t.Rune = Pending //The tile under this is filled in later
				default:
					panic("Invalid character in map")
				}
				room.Map[y][x] = t

			}
		}

		for y, row := range room.Map {
			for x, tile := range row {
				if tile.Rune == Pending {
					room.Map[y][x].Rune = room.bfsFloor(Pos{x, y})
				}
			}
		}
		rooms = append(rooms, room)
	}
	ship := &Ship{rooms}
	ships = append(ships, ship)


	
	
	return ships
}

// type GameEvent int
// const(
// 	Move GameEvent = iota
// 	DoorOpen
// 	Portal
// 	AccessTerminal
// )

// func (ship *Ship) AddEvent(event string) {
// 	level.Events[level.EventPos] = event
// 	level.EventPos++
// 	if level.EventPos == len(level.Events) {
// 		level.EventPos = 0
// 	}
// }
