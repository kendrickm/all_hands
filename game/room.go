package game
type Room struct {
	Map    [][]Tile
	Player *Player
	Debug  map[Pos]bool
	Ship   *Ship
	Terminals map[Pos]*Terminal
	Stations map[Pos]*Station
	//Portals
}

// TODO Load full images instead of building tile by tile
type Tile struct {
	Rune        rune
	OverlayRune rune
}

type StationType int
const (
	Reactor StationType = iota
	AuxPower
)

type Station struct {
	Type StationType
	Active bool
	Name string
}

const (
	Bulkhead   rune = '#'
	ShipFloor  rune = '.'
	ClosedDoor       rune = '|'
	OpenDoor         rune = '/'
	TerminalAccess   rune = 'T'
	UnpoweredReactor rune = 'r'
	PoweredReactor   rune = 'R'
	Blank      rune = 0
	Pending    rune = -1
)

func createReactorStation() *Station {
	return &Station{Type:Reactor, Active:false, Name:"Main Reactor"}
}

func (room *Room) Update() {
	for _, ter := range room.Terminals{
		if ter.Buttons[0].state {
			ter.LinkedStation.Active = true
		}
	}

	for pos, station := range room.Stations {
		if station.Active{
			if room.Map[pos.Y][pos.X].Rune == UnpoweredReactor {
				room.Map[pos.Y][pos.X].Rune = PoweredReactor
			}
		}
	}
}
