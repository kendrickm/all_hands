package game

type Room struct {
	Rooms  []*Room
	Map    [][]Tile
	Player *Player
	Debug  map[Pos]bool
	Ship   *Ship
	// Terminals map[game.Pos]*Terminal
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
	ClosedDoor rune = '|'
	OpenDoor   rune = '/'
	Blank      rune = 0
	Pending    rune = -1
)
