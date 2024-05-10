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
	width := math.Ceil(math.Sqrt(float64(numBuffers)))
	height := math.Ceil(float64(numBuffers) / width)

	return Size{int(width), int(height)}
}

func (r *Relation) GetNumbBuffers() int {
	return len(r.Fsm)
}
