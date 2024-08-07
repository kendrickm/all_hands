package game

type Entity struct {
	Pos
	Name string
	Rune rune
}

type Character struct {
	Entity
}

type Player struct {
	Character
}

func newPlayer() *Player {
	player := &Player{}
	player.Name = "SpaceMan"
	player.Rune = '@'

	return player
}
