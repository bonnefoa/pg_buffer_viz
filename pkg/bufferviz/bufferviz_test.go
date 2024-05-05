package bufferviz

import (
	"testing"

	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/stretchr/testify/assert"
)

func getTestRelation(sizeRelation int) db.Relation {
	fsm := make([]int16, sizeRelation)
	for i := 0; i < sizeRelation; i++ {
		fsm = append(fsm, int16(i))
	}
	return db.Relation{
		Name: "TestRelation", Fsm: fsm,
	}
}

func getTestTable(sizeRelation int, sizeIndexes []int, toastSize int, toastIndexSize int) db.Table {
	var table db.Table
	table.Indexes = make([]db.Relation, 0)
	for _, sizeIndex := range sizeIndexes {
		table.Indexes = append(table.Indexes, getTestRelation(sizeIndex))
	}
	table.Relation = getTestRelation(sizeRelation)
	toast := db.Toast{
		Relation: getTestRelation(toastSize),
		Index:    getTestRelation(toastIndexSize),
	}
	table.Toast = &toast
	return table
}

func TestRelationSize(t *testing.T) {
	testCases := []struct {
		desc         string
		relation     db.Relation
		expectedSize Size
	}{
		{"Test 3 elements", getTestRelation(3), Size{3, 3}},
		{"Test 4 elements", getTestRelation(4), Size{3, 3}},
		{"Test 5 elements", getTestRelation(5), Size{4, 3}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bv := NewBufferViz(nil, Size{1, 1}, Size{1, 1})
			relSize := bv.getRelationSize(tC.relation)
			assert.Equal(t, tC.expectedSize, relSize)
		})
	}
}

func TestAncillarySize(t *testing.T) {
	testCases := []struct {
		desc         string
		table        db.Table
		expectedSize Size
	}{
		{"Test 3 elements", getTestTable(1, []int{1, 2, 3}, 4, 5), Size{2, 2}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bv := NewBufferViz(nil, Size{1, 1}, Size{1, 1})
			tableSize := bv.getDrawSize(tC.table)
			assert.Equal(t, tC.expectedSize, tableSize)
		})
	}
}
