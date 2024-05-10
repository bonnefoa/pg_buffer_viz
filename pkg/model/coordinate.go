package model

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) AddWidth(s Size) {
	c.X += s.Width
}

func (c *Coordinate) AddHeight(s Size) {
	c.Y += s.Height
}

func (c *Coordinate) AddSize(s Size) {
	c.X += s.Width
	c.Y += s.Height
}
