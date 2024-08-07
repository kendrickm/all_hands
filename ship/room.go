package ship


type Room struct {
	Rooms []*Room
	Map      [][]Tile
	Player   *Player
	Debug     map[game.Pos]bool
	// Terminals map[game.Pos]*Terminal
	//Portals 
}