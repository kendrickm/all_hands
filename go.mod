module github.com/kendrickm/all_hands

replace github.com/kendrickm/all_hands/game => ./game

replace github.com/kendrickm/all_hands/ui => ./ui

go 1.19

require github.com/kendrickm/all_hands/game v0.0.0-20240124234230-733e80e60cdc

require (
	github.com/kendrickm/all_hands/ui v0.0.0-00010101000000-000000000000 // indirect
	github.com/veandco/go-sdl2 v0.4.40 // indirect
)
