package model

import "math"

type Relation struct {
	Name string
	Fsm  []int16
}

type Table struct {
	Relation
	Indexes []Relation
	Toast   *Toast
}

type Toast struct {
	Relation
	Index Relation
}

func (r *Relation) GetRelationSize() Size {
	numBuffers := len(r.Fsm)
	width := int(math.Sqrt(float64(numBuffers))) + 1
	height := numBuffers / width

	return Size{width, height}
}

func (r *Relation) GetNumbBuffers() int {
	return len(r.Fsm)
}
