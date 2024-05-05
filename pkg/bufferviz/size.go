package bufferviz

type Size struct {
	Width  int
	Height int
}

func (s *Size) add(s2 Size) {
	s.Height += s2.Height
	s.Width += s2.Width
}

func (s *Size) mul(s2 Size) {
	s.Height *= s2.Height
	s.Width *= s2.Width
}
