package bufferviz

import (
	"github.com/bonnefoa/pg_buffer_viz/pkg/model"
	"github.com/sirupsen/logrus"
)

func (b *BufferViz) getRelationSize(relation model.Relation) (res model.Size) {
	res = relation.GetRelationSize()
	res.Add(b.MarginSize)
	if res.Width <= 5 {
		res.Width = 10
	}
	logrus.Infof("Size of relation %s: %v", relation.Name, res)
	return res
}

func (b *BufferViz) getAncillarySize(table model.Table) (res model.Size) {
	for _, index := range table.Indexes {
		res.AddWidthMaxHeight(b.getRelationSize(index))
	}

	if table.Toast != nil {
		toast := table.Toast
		res.AddWidthMaxHeight(b.getRelationSize(toast.Relation))
		res.AddWidthMaxHeight(b.getRelationSize(toast.Index))
	}

	return res
}

func (b *BufferViz) getDrawSize(table model.Table) (res model.Size) {
	res = b.getRelationSize(table.Relation)
	ancillarySize := b.getAncillarySize(table)
	res.AddHeightMaxWidth(ancillarySize)
	return res
}
