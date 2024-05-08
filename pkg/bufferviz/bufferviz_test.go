package bufferviz

import (
	"testing"

	"github.com/bonnefoa/pg_buffer_viz/pkg/model"
	"github.com/stretchr/testify/require"
)

func getTestRelation(sizeRelation int) model.Relation {
	fsm := make([]int16, sizeRelation)
	for i := 0; i < sizeRelation; i++ {
		fsm[i] = int16(i)
	}
	return model.Relation{
		Name: "TestRelation", Fsm: fsm,
	}
}

func getTestTable(sizeRelation int, sizeIndexes []int, toastSize int, toastIndexSize int) model.Table {
	var table model.Table
	table.Indexes = make([]model.Relation, 0)
	for _, sizeIndex := range sizeIndexes {
		table.Indexes = append(table.Indexes, getTestRelation(sizeIndex))
	}
	table.Relation = getTestRelation(sizeRelation)
	toast := model.Toast{
		Relation: getTestRelation(toastSize),
		Index:    getTestRelation(toastIndexSize),
	}
	table.Toast = &toast
	return table
}

func TestRelationSize(t *testing.T) {
	testCases := []struct {
		desc         string
		relation     model.Relation
		expectedSize model.Size
		margin       model.Size
	}{
		{"Test 3 elements no margin", getTestRelation(3), model.Size{Width: 2, Height: 2}, model.Size{Width: 0, Height: 0}},
		{"Test 4 elements no margin", getTestRelation(4), model.Size{Width: 2, Height: 2}, model.Size{Width: 0, Height: 0}},
		{"Test 5 elements no margin", getTestRelation(5), model.Size{Width: 3, Height: 2}, model.Size{Width: 0, Height: 0}},
		{"Test 9 elements no margin", getTestRelation(9), model.Size{Width: 3, Height: 3}, model.Size{Width: 0, Height: 0}},
		{"Test 256 elements no margin", getTestRelation(256), model.Size{Width: 16, Height: 16}, model.Size{Width: 0, Height: 0}},
		{"Test 257 elements no margin", getTestRelation(257), model.Size{Width: 17, Height: 16}, model.Size{Width: 0, Height: 0}},

		{"Test 3 elements 1_2 margin", getTestRelation(3), model.Size{Width: 3, Height: 4}, model.Size{Width: 1, Height: 2}},
		{"Test 4 elements 1_2 margin", getTestRelation(4), model.Size{Width: 3, Height: 4}, model.Size{Width: 1, Height: 2}},
		{"Test 5 elements 1_2 margin", getTestRelation(5), model.Size{Width: 4, Height: 4}, model.Size{Width: 1, Height: 2}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bv := NewBufferViz(nil, model.Size{Width: 1, Height: 1}, tC.margin)
			relSize := bv.getRelationSize(tC.relation)
			require.Equal(t, tC.expectedSize, relSize)
		})
	}
}

func TestAncillarySize(t *testing.T) {
	testCases := []struct {
		desc         string
		table        model.Table
		expectedSize model.Size
		margin       model.Size
	}{
		{"Test only index", getTestTable(257, []int{9}, 0, 0), model.Size{Width: 3, Height: 3}, model.Size{Width: 0, Height: 0}},
		{"Test index plus toast no margin", getTestTable(257, []int{9}, 3, 3), model.Size{Width: 3 + 2 + 2, Height: 3}, model.Size{Width: 0, Height: 0}},
		{"Test index plus toast 1_1 margin", getTestTable(257, []int{9}, 3, 3), model.Size{Width: 3 + 2 + 2 + 3, Height: 3 + 1}, model.Size{Width: 1, Height: 1}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bv := NewBufferViz(nil, model.Size{Width: 1, Height: 1}, tC.margin)
			tableSize := bv.getAncillarySize(tC.table)
			require.Equal(t, tC.expectedSize, tableSize)
		})
	}
}
