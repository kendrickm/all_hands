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
	Level int16
	MinLevel int16
	MaxLevel int16
	changeRate float32 //When active, the rate per frame(1/60 of second) to grow Level
	tmp float32 //Holds growth between numbers
}

const (
	Bulkhead   rune = '#'
	ShipFloor  rune = '.'
	ClosedDoor       rune = '|'
	OpenDoor         rune = '/'
	TerminalAccess   rune = 'T'
	UnpoweredReactor rune = 'r'
	PoweredReactor   rune = 'R'
	GraphDisplay rune = 'G'
	Blank      rune = 0
	Pending    rune = -1
)

func createReactorStation() *Station {
	changeRate := float32(0.0008)//5/s grow rate
	return &Station{Type:Reactor, Active:false, Name:"Main Reactor", Level: 0, MinLevel:0, MaxLevel:100, changeRate:changeRate, tmp: 0.0}
}

func (room *Room) Update() {
	for _, ter := range room.Terminals{
		if ter.Type == GRAPH_DISPLAY{
			// skip for now
		} else {
			if ter.Buttons[0].state {
			ter.LinkedStation.Active = true
		}
		}
	}
	for pos, station := range room.Stations {
		// fmt.Println(station.Name + ": " + strconv.Itoa(int(station.Level)))
		if station.Active{
			if station.Level < station.MaxLevel {
				// fmt.Println("Updating levels")
				// fmt.Println(strconv.FormatFloat(float64(station.tmp) ,'f', -1, 64))
				station.tmp += station.changeRate
				if station.tmp >= 1.0 {
					station.Level += int16(station.tmp)
					if station.Level > station.MaxLevel{
						station.Level = station.MaxLevel
					}
					station.tmp = 0.0
				}
			}
			if station.Level == station.MaxLevel {
				if room.Map[pos.Y][pos.X].Rune == UnpoweredReactor {
				room.Map[pos.Y][pos.X].Rune = PoweredReactor
				}
			}
			
		}
	}
}
