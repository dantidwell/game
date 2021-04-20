package assets

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten"
)

type ImageInfo struct {
	Width       int           `json:"width"`
	Height      int           `json:"height"`
	EngineImage *ebiten.Image `json:"-"`
}

type AssetType int

type AssetInfo struct {
	Name   string    `json:"name"`
	Length int       `json:"length"`
	Offset int       `json:"offset"`
	Type   AssetType `json:"type"`

	Image *ImageInfo `json:"image"`
}

const (
	TypeImage = iota
)

type Pack struct {
	Info []AssetInfo
	Data []byte
}

func Load(r io.ReadSeeker) *Pack {
	var p Pack

	_, _ = r.Seek(0, io.SeekStart)
	l, _ := r.Seek(0, io.SeekEnd)
	_, _ = r.Seek(0, io.SeekStart)

	var n, pos int
	var dataLenBuf [8]byte

	n, _ = r.Read(dataLenBuf[:])
	pos += n

	p.Data = make([]byte, binary.LittleEndian.Uint64(dataLenBuf[:]))
	n, _ = r.Read(p.Data)
	pos += n

	infoBuf := make([]byte, int(l)-pos)
	n, _ = r.Read(infoBuf)
	pos += n

	if err := json.Unmarshal(infoBuf, &p.Info); err != nil {
		panic(err)
	}
	return &p
}

func (p *Pack) GetImage(name string) *ebiten.Image {
	var info *AssetInfo
	for i := range p.Info {
		if p.Info[i].Name == name {
			info = &p.Info[i]
			break
		}
	}
	if info == nil {
		panic(fmt.Errorf("asset not found: %s", name))
	}
	if info.Image.EngineImage == nil {
		info.Image.EngineImage, _ = ebiten.NewImageFromImage(&image.RGBA{
			Pix:    p.Data[info.Offset : info.Offset+info.Length],
			Stride: 4 * info.Image.Width,
			Rect:   image.Rect(0, 0, info.Image.Width, info.Image.Height),
		}, ebiten.FilterDefault)
	}
	return info.Image.EngineImage
}
