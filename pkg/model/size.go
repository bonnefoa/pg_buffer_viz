package model

type Size struct {
	Width  int
	Height int
}

func (s *Size) Add(s2 Size) {
	s.Height += s2.Height
	s.Width += s2.Width
}

func (s *Size) AddWidthMaxHeight(s2 Size) {
	s.Width += s2.Width
	if s2.Height > s.Height {
		s.Height = s2.Height
	}
}

func (s *Size) AddHeightMaxWidth(s2 Size) {
	s.Height += s2.Height
	if s2.Width > s.Width {
		s.Width = s2.Width
	}
}

func (s *Size) GetMaxWidth(s2 Size) int {
	if s2.Width > s.Width {
		return s2.Width
	}
	return s.Width
}

func (s *Size) AddWidth(s2 Size) {
	s.Width += s2.Width
}

func (s *Size) AddHeight(s2 Size) {
	s.Width += s2.Width
}

func (s *Size) Mul(s2 Size) {
	s.Height *= s2.Height
	s.Width *= s2.Width
}
