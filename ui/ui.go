package ui

import (
	"fmt"
	"unsafe"
	"math/rand"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/mix"
	"github.com/kendrickm/all_hands/game"
)

func f(p unsafe.Pointer) {}

func init() {
	fmt.Println("Init innit")
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	err = mix.Init(mix.INIT_OGG)
		if err != nil {
		panic(err)
	}

}


//Migrate to a state machine
type uiState int
const (
	UIMain uiState = iota
	UITerminal
)


type ui struct {
	state uiState
	sounds sounds
	winWidth  int
	winHeight int

	renderer         *sdl.Renderer
	window           *sdl.Window
	textureAtlas     *sdl.Texture
	textureIndex     map[rune][]sdl.Rect
	preKeyboardState []uint8
	keyboardState    []uint8
	r                *rand.Rand
	centerX          int
	centerY          int

	roomChan chan *game.Room
	inputChan chan *game.Input

	fontMedium *ttf.Font
	fontLarge  *ttf.Font
	fontSmall  *ttf.Font

	string2TexSmall map[string]*sdl.Texture
	string2TexMed   map[string]*sdl.Texture
	string2TexLarge map[string]*sdl.Texture

	currentMouseState *mouseState
	prevMouseState *mouseState
}

func NewUI(inputChan chan *game.Input, roomChan chan *game.Room) *ui {
	ui := &ui{}
	ui.state = UIMain
	ui.inputChan = inputChan
	ui.roomChan = roomChan
	ui.string2TexSmall = make(map[string]*sdl.Texture)
	ui.string2TexMed = make(map[string]*sdl.Texture)
	ui.string2TexLarge = make(map[string]*sdl.Texture)
	ui.winHeight = 720
	ui.winWidth = 1080
	ui.r = rand.New(rand.NewSource(1))
	window, err := sdl.CreateWindow("All Hands!", 200, 200,
		int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(ui.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	//sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.textureAtlas = ui.imgFileToTexture("ui/assets/tiles.png")
	ui.loadTextureIndex()

	ui.keyboardState = sdl.GetKeyboardState()
	ui.preKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.preKeyboardState[i] = v
	}

	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("ui/assets/gothic.ttf", int(float64(ui.winWidth)*.015))
	if err != nil {
		panic(err)
	}

	ui.fontMedium, err = ttf.OpenFont("ui/assets/gothic.ttf", 32)
	if err != nil {
		panic(err)
	}

	ui.fontLarge, err = ttf.OpenFont("ui/assets/gothic.ttf", 64)
	if err != nil {
		panic(err)
	}

	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT,2,4096)
	if err != nil {
		panic(err)
	}

	return ui
}

func (ui *ui) Run() {

	ui.prevMouseState = getMouseState()
	var newRoom  *game.Room

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Typ: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow}
					return
				}

			}
		}

		ui.currentMouseState = getMouseState()
		var ok bool
		select {
		case newRoom, ok = <-ui.roomChan:
			if ok {
				ui.Draw(newRoom)
			}
		default:
		}

		ui.Draw(newRoom)
		var input game.Input

		ui.renderer.Present()

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			
			if ui.keyDownOnce(sdl.SCANCODE_UP) {
				input.Typ = game.Up
			} else if ui.keyDownOnce(sdl.SCANCODE_DOWN) {
				input.Typ = game.Down
			} else if ui.keyDownOnce(sdl.SCANCODE_RIGHT) {
				input.Typ = game.Right
			} else if ui.keyDownOnce(sdl.SCANCODE_LEFT) {
				input.Typ = game.Left
			} 
		}

			for i, v := range ui.keyboardState {
				ui.preKeyboardState[i] = v
			}

			if input.Typ != game.None {
				ui.inputChan <- &input
			}
		}
		ui.prevMouseState = ui.currentMouseState
		sdl.Delay(10)
}

func (ui *ui) Draw(room *game.Room) {

	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = room.Player.X
		ui.centerY = room.Player.Y
	}

	limit := 5
	if room.Player.X > ui.centerX+limit {
		diff := room.Player.X - (ui.centerX+limit)
		ui.centerX += diff
	} else if room.Player.X < ui.centerX-limit {
		diff := (ui.centerX-limit) - room.Player.X
		ui.centerX -= diff
	} else if room.Player.Y > ui.centerY+limit {
		diff :=  room.Player.Y - (ui.centerY+limit)
		ui.centerY += diff
	} else if room.Player.Y < ui.centerY-limit {
		diff := (ui.centerY-limit) - room.Player.Y
		ui.centerY -= diff
	}

	offsetX := int32((ui.winWidth / 2) - ui.centerX*32)
	offsetY := int32((ui.winHeight / 2) - ui.centerY*32)
	ui.renderer.Clear()
	ui.r.Seed(2)
	for y, row := range room.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRects := ui.textureIndex[tile.Rune]
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
					dstRect := sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}
					pos := game.Pos{x, y}
					if room.Debug[pos] {
						ui.textureAtlas.SetColorMod(128, 0, 0)
					} else {
						ui.textureAtlas.SetColorMod(255, 255, 255)
					}
					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)

					if tile.OverlayRune != game.Blank {
						// Todo what if there are multiple varients for overlay images?
						srcRect := ui.textureIndex[tile.OverlayRune][0]
						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					}
			}
		}
	}
	ui.textureAtlas.SetColorMod(255, 255, 255)

	//Draw Player
	playerSrcRect := ui.textureIndex[room.Player.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{X: int32(room.Player.X)*32 + offsetX, Y: int32(room.Player.Y)*32 + offsetY, W: 32, H: 32})

}