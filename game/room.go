package game

type Room struct {
	Map    [][]Tile
	Player *Player
	Debug  map[Pos]bool
	Ship   *Ship
	Terminals map[Pos]*Terminal
	//Portals
}

// TODO Load full images instead of building tile by tile
type Tile struct {
	Rune        rune
	OverlayRune rune
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

func (room *Room) Update() {

}
