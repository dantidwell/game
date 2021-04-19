package main

import (
	"image/png"
	"log"
	"os"
)

var sprites = []struct {
	x, y, w, h int
}{
	{1, 26, 16, 32},  // super_mario_stand
	{22, 26, 16, 32}, // super_mario_duck
	{43, 26, 16, 32}, // super_mario_walk_f0
	{60, 26, 16, 32}, // super_mario_walk_f1
	{77, 26, 16, 32}, // super_mario_walk_f2
}

// standing still: (1,26,16,32)
// ducking: (22,26,16,32)
// walk_one: (43, 26, 16, 32)
// walk_two: (60,26,16,32)
// walk_three: (77,26,16,32)

func main() {
	img, err := png.Decode(os.Stdin)
	if err != nil {
		log.Fatal("assets: failed to decode source image")
	}

	for _, s := range sprites {
		for y := s.y; y < s.y+s.h; y++ {
			for x := s.x; x < s.x+s.w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				os.Stdout.Write([]byte{
					uint8(r),
					uint8(g),
					uint8(b),
					uint8(a),
				})
			}
		}
	}
}
