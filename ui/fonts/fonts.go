package fonts

import (
	_ "embed"
	"image/color"

	g "github.com/AllenDang/giu"
)

var (
	ColorDefault     = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	ColorDestructive = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	ColorCaption     = color.RGBA{0xD9, 0xD9, 0xD9, 0xFF}
)

const (
	HeaderSize    = 18
	SubHeaderSize = 14
	CaptionSize   = 11
)

type FontType int

const (
	DefaultFontType FontType = iota
)

var fonts = make([]font, 1)

type font struct {
	fontType FontType
	name     string
	fontInfo *g.FontInfo
	bytes    []byte
}

var GlyphFont *g.FontInfo

func AddFont(fontType FontType, fontName string, fontBytes []byte) {
	fonts[fontType].fontType = fontType
	fonts[fontType].name = fontName
	fonts[fontType].bytes = fontBytes
}

func InitializeFonts() {
	for _, f := range fonts {
		f.fontInfo = g.Context.FontAtlas.AddFontFromBytes(f.name, f.bytes, 11)
		switch f.fontType {
		case DefaultFontType:
			GlyphFont = f.fontInfo
		}
	}
}
