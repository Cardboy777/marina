package assets

import (
	_ "embed"
	"marina/ui/fonts"
)

//go:embed vendor/font-awesome/font-awesome-6-solid.otf
var glyphBytes []byte
var glyphFontName = "Font Awesome 6 Solid"

func LoadAssets() {
	fonts.AddFont(fonts.GlyphFontType, glyphFontName, glyphBytes)
}
