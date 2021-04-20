package main

import (
	"encoding/binary"
	"encoding/json"
	"image/png"
	"log"
	"os"

	"github.com/dantidwell/game/assets"
)

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
}

func main() {
	var info []assets.AssetInfo
	var data []byte
	for _, t := range tiles {
		f, err := os.Open("data/" + t.sheet)
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
