package ship

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Ship struct {
	Rooms []*Room
}


//TODO Load full images instead of building tile by tile
type Tile struct {
	Rune        rune
	OverlayRune rune
}

const (
	Bulkhead  rune = '#'
	ShipFloor  rune = '.'
	Blank      rune = 0
	Pending    rune = -1
)


func loadChapter0(player *game.Player) {
	// TODO refactor the depreacted function
	// TODO abstract loading that will span chapters
	// TODO load series of rooms from a single ship file
	// TODO load a series of ships from chapter logic
	filenames, err := filepath.Glob("ship/rooms/*.room")
	if err != nil {
		panic(err)
	}

	rooms := []*Room

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

		room := &Level{}
		level.Debug = make(map[game.Pos]bool)
		level.Events = make([]string, 10)
		level.Player = player

		level.Map = make([][]Tile, len(levelLines))
		level.Monsters = make(map[Pos]*Monster)
		level.Items = make(map[Pos][]*Item)
		level.Portal = make(map[Pos]*LevelPos)

		for i := range level.Map {
			level.Map[i] = make([]Tile, longestRow)
		}

		for y := 0; y < len(level.Map); y++ {
			line := levelLines[y]

			for x, c := range line {
				var t Tile
				p := Pos{x,y}
				t.OverlayRune = Blank
				switch c {
				case ' ', '\n', '\t', '\r':
					t.Rune = Blank
				case '#':
					t.Rune = StoneWall
				case '|':
					t.OverlayRune = ClosedDoor
					t.Rune = Pending
				case '/':
					t.Rune = Pending
					t.OverlayRune = OpenDoor
				case 'u':
					t.Rune = Pending
					t.OverlayRune = UpStair
				case 'd':
					t.Rune = Pending
					t.OverlayRune = DownStair
				case 's':
					t.Rune = Pending
					level.Items[p] = append(level.Items[p], NewSword(p))
					level.Items[p] = append(level.Items[p], NewHelmet(p))
				case 'h':
					t.Rune = Pending
					level.Items[p] = append(level.Items[p], NewHelmet(p))
				case '.':
					t.Rune = DirtFloor
				case '@':
					level.Player.X = x
					level.Player.Y = y
					t.Rune = Pending //The tile under this is filled in later
				case 'R':
					level.Monsters[p] = NewRat(p)
					t.Rune = Pending
				case 'S':
					level.Monsters[p] = NewSpider(p)
					t.Rune = Pending
				default:
					panic("Invalid character in map")
				}
				level.Map[y][x] = t

			}
		}

		for y, row := range level.Map {
			for x, tile := range row {
				if tile.Rune == Pending {
					level.Map[y][x].Rune = level.bfsFloor(Pos{x, y})
				}
			}
		}
		levels[levelName] = level
	}
	return levels
}
