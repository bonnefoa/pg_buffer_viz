package model

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

type RelationFreeSpace struct {
	Name string
	Fsm  []int
}

//func (r *Relation) GetRelationSize() int {
//	return math.Sqrt(len(r.Fsm))
//}
