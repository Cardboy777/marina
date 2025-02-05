package marina

type Color struct {
	r uint32
	g uint32
	b uint32
	a uint32
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.r, c.g, c.b, c.a
}

func NewColor(r, g, b, a uint32) Color {
	return Color{r, g, b, a}
}
