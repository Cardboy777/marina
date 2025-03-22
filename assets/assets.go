package assets

import (
	_ "embed"
	"marina/ui/fonts"
)

//go:embed vendor/fonts/ProFontIIxNerdFont-Regular.ttf
var fontBytes []byte
var fontName = "ProggyClean Nerd Font"

func LoadAssets() {
	fonts.AddFont(fonts.DefaultFontType, fontName, fontBytes)
}
