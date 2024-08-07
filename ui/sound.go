package ui

import (
	"github.com/veandco/go-sdl2/mix"
	"fmt"
	"strconv"
	"math/rand"
)

type sounds struct {
	openingDoors []*mix.Chunk
	footsteps  []*mix.Chunk
}

func playRandomSound(chunks []*mix.Chunk, volume int){
	chunkIndex := rand.Intn(len(chunks))
	fmt.Println("Playing file " + strconv.Itoa(chunkIndex))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1,0)
}