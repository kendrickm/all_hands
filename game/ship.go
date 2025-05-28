package game

import (
		"path/filepath"
		"strings"
		"fmt"
		"os"
		"bufio"
		"strconv"
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
			text := scanner.Text()
			if text == "-------" {
				break
			} else {
				roomLines = append(roomLines, text)
				if len(roomLines[index]) > longestRow {
					longestRow = len(roomLines[index])
				}
				index++
			}
		}


		room := &Room{}
		room.Debug = make(map[Pos]bool)
		room.Player = player

		room.Map = make([][]Tile, len(roomLines))
		room.Terminals = make(map[Pos]*Terminal, 1)
		room.Stations = make(map[Pos]*Station, 1)

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
				case 'L':
					t.Rune = ULCornerBulkhead
				case 'l':
					t.Rune = LLCornerBulkhead
				case 'R':
					t.Rune = URCornerBulkhead
				case 'r':
					t.Rune = LRCornerBulkhead
				case 'g':
					t.Rune = LeftBulkhead
				case 'h':
					t.Rune = RightBulkhead
				case 'C':
					t.Rune = CenterTBulkhead
				case 'c':
					t.Rune = CenterBBulkhead
				case '#':
					t.Rune = FillerBulkhead
				case '.':
					t.Rune = ShipFloor
				case 'T': 
					 t.Rune = TerminalAccess
				case 'G':
					  t.Rune = GraphDisplay
			    case 'p':
			    	 t.Rune = UnpoweredReactor
			    	 // room.Stations[Pos{x,y}] = createReactorStation()
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
		// for _, x := range room.Terminals{
		// 	for _, y := range room.Stations {
		// 		x.LinkedStation = y
		// 	}
		// }
		scanner.Scan()
		testStr := scanner.Text()
		if testStr != "STATIONS" {
			panic("Missing data in room file" + testStr)
		}

		for scanner.Scan() { //Start STATION block
			text := scanner.Text()
			if text == "-------" {
				break
			} else {
				//Start with POS
				// xy := text[1:]
				splitXYCount := strings.Split(text, ",")
				x, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[0]), 10, 64)
				if err != nil {
					panic(err)
				}
				y, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[1]), 10, 64)
				scanner.Scan()
				name := scanner.Text() //Get name
				scanner.Scan()
				text = scanner.Text() //Get type
				var typ StationType 
				switch text {
				case "REACTOR":
					fmt.Println("Creating reactor station")
					typ = Reactor
				case "AUX_POWER":
					fmt.Println("Creating Aux power station")
					typ = AuxPower
				default:
					panic("Invalid type " + text)
				}
				scanner.Scan()
				text = scanner.Text() //Get powered status
				var active bool
				if text == "0"{
					active = false
				} else {
					active = true
				}

				room.Stations[Pos{int(x),int(y)}] = &Station{Type: typ, Active: active, Name: name, Level: 0, MinLevel:0, MaxLevel:100, changeRate:0.08, tmp: 0.0}

			}
		}

		scanner.Scan()
		testStr = scanner.Text()
		if testStr != "TERMINALS" {
			panic("Missing data in room file" + testStr)
		}

		for scanner.Scan() { //Start TERMINAL block
			text := scanner.Text()
			if text == "-------" {
				break
			} else {
				var ter *Terminal
				name := scanner.Text() //Get name
				fmt.Println(name)
				scanner.Scan()
				text = scanner.Text() //Get pos
				splitXYCount := strings.Split(text, ",")
				x, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[0]), 10, 64)
				if err != nil {
					panic(err)
				}
				y, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[1]), 10, 64)				
				scanner.Scan()
				text = scanner.Text() //Get type
				switch text {
				case "SINGLE_BUTTON":
					ter = singleButtonTerminalFactory()
				case "GRAPH_DISPLAY":
					ter = graphDisplayTerminalFactor()
				default:
					panic("Invalid type " + text)
				}
				ter.Name = name
				scanner.Scan()
				//TODO: Handle different number of possible buttons
				text = scanner.Text() //Get button type 
				fmt.Println(text)
				switch text {
				case "TRIGGER":
					ter.Buttons[0]= &Button{TriggerButton, false,""}
				case "TOGGLE":
					ter.Buttons[0]= &Button{ToggleSwitch, false,""}
				case "0":
				default:
					panic("Invalid type " + text)
				}
				for x := range ter.Buttons{
					scanner.Scan()
					text = scanner.Text() //Get button display text
					ter.Buttons[x].DisplayText = text
				}

				scanner.Scan()
				text = scanner.Text() //Get status
				if text == "0"{
					ter.Powered = false
				} else {
					ter.Powered = true
				}
				scanner.Scan()
				text = scanner.Text() //Get linked station name
				for _, station := range room.Stations {
					if station.Name == text {
						ter.LinkedStation = station
						break
					}
				}
				scanner.Scan()
				text = scanner.Text() //Get flavor text
				ter.DisplayText = text
				pos := Pos{int(x),int(y)}
				fmt.Println("Creating " + ter.Name + " at ")
				fmt.Println(pos)
				room.Terminals[Pos{int(x),int(y)}] = ter

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
