package ui

import (
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"os"
	"log"
	"unsafe"
	"strconv"
	"strings"
	"bufio"
)

type FontSize int

const (
	FontSmall FontSize = iota
	FontMedium
	FontLarge
)


func (ui *ui) GetSinglePixelTex(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}
	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A

	tex.Update(nil, unsafe.Pointer(&pixels[0]), 4)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func (ui *ui) stringToTexture(s string, size FontSize, color sdl.Color) *sdl.Texture {

	var font *ttf.Font
	switch size {
	case FontSmall:
		font = ui.fontSmall
		tex, exists := ui.string2TexSmall[s]
		if exists {
			return tex
		}
	case FontMedium:
		font = ui.fontMedium
		tex, exists := ui.string2TexMed[s]
		if exists {
			return tex
		}
	case FontLarge:
		font = ui.fontLarge
		tex, exists := ui.string2TexLarge[s]
		if exists {
			return tex
		}
	}
	fontSurface, err := font.RenderUTF8Blended(s, sdl.Color{255, 0, 0, 0})
	if err != nil {
		panic(err)
	}
	fontTexture, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic(err)
	}

	switch size {
	case FontSmall:
		ui.string2TexSmall[s] = fontTexture
	case FontMedium:
		ui.string2TexMed[s] = fontTexture
	case FontLarge:
		ui.string2TexLarge[s] = fontTexture
	}

	return fontTexture
}

func (ui *ui) loadTextureIndex() {
	ui.textureIndex = make(map[rune][]sdl.Rect)
	file, err := os.Open("ui/assets/asset-index.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := rune(line[0])
		xy := line[1:]
		splitXYCount := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYCount[2]), 10, 64) //Supports randomly picking from a batch of tiles
		if err != nil {
			panic(err)
		}
		var rects []sdl.Rect
		for i := int64(0); i < variationCount; i++ {
			rects = append(rects, sdl.Rect{int32(x * 32), int32(y * 32), 32, 32})
			x++
			if x > 62 { //handles wrap arounds
				x = 0
				y++
			}

		}

		ui.textureIndex[tileRune] = rects
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	image, err := img.Load(filename)
	if err != nil {
		panic(err)
	}

	tex, err := ui.renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}