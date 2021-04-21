package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/dantidwell/game/assets"
)

var glyphs = []byte{
	0:  ' ',
	1:  '!',
	2:  '"',
	3:  '#',
	4:  '$',
	5:  '%',
	6:  '&',
	7:  '\'',
	8:  '(',
	9:  ')',
	10: '*',
	11: '+',
	12: ',',
	13: '-',
	14: '.',
	15: '/',
	16: '0',
	17: '1',
	18: '2',
	19: '3',
	20: '4',
	21: '5',
	22: '6',
	23: '7',
	24: '8',
	25: '9',
	26: ':',
	27: ';',
	28: '<',
	29: '=',
	30: '>',
	31: '?',
	32: '@',
	33: 'A',
	34: 'B',
	35: 'C',
	36: 'D',
	37: 'E',
	38: 'F',
	39: 'G',
	40: 'H',
	41: 'I',
	42: 'J',
	43: 'K',
	44: 'L',
	45: 'M',
	46: 'N',
	47: 'O',
	48: 'P',
	49: 'Q',
	50: 'R',
	51: 'S',
	52: 'T',
	53: 'U',
	54: 'V',
	55: 'W',
	56: 'X',
	57: 'Y',
	58: 'Z',
	59: '[',
	60: '\\',
	61: ']',
	62: '^',
	63: '_',
	64: '`',
	65: 'a',
	66: 'b',
	67: 'c',
	68: 'd',
	69: 'e',
	70: 'f',
	71: 'g',
	72: 'h',
	73: 'i',
	74: 'j',
	75: 'k',
	76: 'l',
	77: 'm',
	78: 'n',
	79: 'o',
	80: 'p',
	81: 'q',
	82: 'r',
	83: 's',
	84: 't',
	85: 'u',
	86: 'v',
	87: 'w',
	88: 'x',
	89: 'y',
	90: 'z',
	91: '{',
	92: '|',
	93: '}',
	94: '~',
}

var songs = []struct {
	name string
}{}

var tiles = []struct {
	sheet, name string
	x, y, w, h  int
}{
	{"dan.png", "dan_right", 0, 10, 16, 22},
	{"dan.png", "dan_up", 16, 10, 16, 22},
	{"dan.png", "dan_left", 32, 10, 16, 22},
	{"dan.png", "dan_down", 48, 10, 16, 22},
	{"room.png", "floor0", 544, 128, 16, 16},
	{"room.png", "border_ne", 80, 32, 16, 16},
	{"room.png", "border_nw", 96, 32, 16, 16},
	{"room.png", "border_se", 80, 6, 16, 16},
	{"room.png", "border_sw", 96, 6, 16, 16},
	{"room.png", "border_top", 135, 16, 16, 16},
	{"room.png", "border_bottom", 135, 6, 16, 16},
	{"room.png", "border_right", 112, 0, 16, 16},
	{"room.png", "border_left", 128, 0, 16, 16},
	{"room.png", "wall", 0, 96, 16, 32},
	{"interiors.png", "shelf_full_vertical", 96, 6808, 32, 48},
}

func main() {
	var info []assets.AssetInfo
	var data []byte

	for _, t := range tiles {
		f, err := os.Open("data/image/" + t.sheet)
		if err != nil {
			log.Fatal("assets: failed to load sprite sheet", t.sheet)
		}
		img, err := png.Decode(f)
		if err != nil {
			log.Fatal("assets: failed to decode source image")
		}

		var buf []byte
		for y := t.y; y < t.y+t.h; y++ {
			for x := t.x; x < t.x+t.w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				buf = append(buf, []byte{
					uint8(r),
					uint8(g),
					uint8(b),
					uint8(a),
				}...)
			}
		}
		info = append(info, assets.AssetInfo{
			Name:   t.name,
			Length: len(buf),
			Offset: len(data),
			Type:   assets.TypeImage,

			Image: &assets.ImageInfo{
				Width:  t.w,
				Height: t.h,
			},
		})
		data = append(data, buf...)

		f.Close()
	}

	f, err := os.Open("data/font/font_8x8.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal("assets: failed to decode source image")
	}
	for i, g := range glyphs {
		w, h := 8, 8
		x0, y0 := ((i * w) % 128), (h * (i / 16))
		os.Stderr.WriteString(fmt.Sprintf("x: %d, y: %d, glyph: %v, idx: %v\n", x0, y0, string(g), i))
		if g == 0 {
			continue
		}

		var buf []byte
		for y := y0; y < y0+h; y++ {
			for x := x0; x < x0+w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				buf = append(buf, []byte{
					uint8(r),
					uint8(g),
					uint8(b),
					uint8(a),
				}...)
			}
		}
		info = append(info, assets.AssetInfo{
			Name:   "font_" + string(g),
			Length: len(buf),
			Offset: len(data),
			Type:   assets.TypeImage,
			Image: &assets.ImageInfo{
				Width:  w,
				Height: h,
			},
		})
		data = append(data, buf...)
	}

	_ = f.Close()

	infoBuf, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}

	var lData [8]byte
	binary.LittleEndian.PutUint64(lData[:], uint64(len(data)))

	os.Stdout.Write(lData[:])
	os.Stdout.Write(data)
	os.Stdout.Write(infoBuf)
}
