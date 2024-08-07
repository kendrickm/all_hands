module github.com/kendrickm/all_hands/ui

replace github.com/kendrickm/all_hands/game => ../game

go 1.19

require (
	github.com/kendrickm/all_hands/game v0.0.0-20240124234230-733e80e60cdc
	github.com/veandco/go-sdl2 v0.4.38
)