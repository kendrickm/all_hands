package game
//Base struct to use for other types of events
type Event struct {
	Triggered bool
	Name string
}

type GameStateEvent struct {
	Event

}