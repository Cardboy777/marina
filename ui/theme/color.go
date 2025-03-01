package theme

import (
	"fmt"
	"image/color"
)

func NewColor(s string) (c color.RGBA) {
	_, err := fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	if err != nil {
		panic(fmt.Errorf("Unable to create color: %w", err))
	}
	return
}
