package fonts

import (
	_ "embed"

	g "github.com/AllenDang/giu"
)

const (
	IconSize      = 22
	HeaderSize    = 19
	SubHeaderSize = 14
	CaptionSize   = 11
)

type FontType int

const (
	GlyphFontType FontType = iota
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
		case GlyphFontType:
			GlyphFont = f.fontInfo
		}
	}
}
