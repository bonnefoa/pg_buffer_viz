package bufferviz

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/model"
	"github.com/sirupsen/logrus"
)

type BufferViz struct {
	canvas *svg.SVG

	BlockSize  model.Size
	MarginSize model.Size

	currentPos model.Coordinate
}

func NewBufferViz(canvas *svg.SVG, blockSize model.Size, marginSize model.Size) BufferViz {
	b := BufferViz{
		canvas:     canvas,
		BlockSize:  blockSize,
		MarginSize: marginSize,
		currentPos: model.Coordinate{X: 0, Y: 0},
	}
	return b
}

func (b *BufferViz) getFsmColor(fsmValue int16) string {
	percent := int((float64(fsmValue) / 8192) * 100)
	return fmt.Sprintf("fill: color-mix(in srgb, green %d%%, red)", percent)
}

func (b *BufferViz) drawName(relation model.Relation) {
	xPos := b.currentPos.X * b.BlockSize.Width
	yPos := b.currentPos.Y * b.BlockSize.Height
	b.canvas.Text(xPos, yPos, relation.Name,
		// "text-align:left;dominant-baseline=hanging;font-size:20px")
		"text-align:left;font-size:20px")
}

func (b *BufferViz) drawRelation(relation model.Relation) model.Size {
	relationSize := relation.GetRelationSize()
	numBuffers := relation.GetNumbBuffers()

	b.drawName(relation)

	pos := b.currentPos
	pos.Y += 1

	// Draw one rect per block
	b.canvas.Gstyle("stroke-width:2;stroke:black;fill:white")
	for line := range relationSize.Width {
		for column := range relationSize.Width {
			bufno := line*relationSize.Width + column
			if bufno >= numBuffers {
				break
			}
			x := (pos.X + column) * b.BlockSize.Width
			y := (pos.Y + line) * b.BlockSize.Height
			style := b.getFsmColor(relation.Fsm[bufno])
			b.canvas.Rect(x, y, b.BlockSize.Width, b.BlockSize.Height, style)
		}
	}
	b.canvas.Gend()
	relationSize.Add(b.MarginSize)

	//	b.canvas.Gstyle("text-anchor:middle;font-size:20px;fill:black;dominant-baseline=middle")
	//	for line := range relationSize.Width {
	//		for column := range relationSize.Width {
	//			bufno := line*relationSize.Width + column
	//			if bufno > numBuffers {
	//				break
	//			}
	//			if bufno%50 == 0 {
	//				x := pos.x + column*b.BlockSize.Width + b.BlockSize.Width/2
	//				y := pos.y + line*b.BlockSize.Height + b.BlockSize.Height/2
	//				b.canvas.Text(x, y, fmt.Sprint(bufno))
	//			}
	//		}
	//	}
	//	b.canvas.Gend()
	return relationSize
}

func (b *BufferViz) DrawTable(table model.Table) {
	drawSize := b.getDrawSize(table)
	b.canvas.Start(
		drawSize.Width*b.BlockSize.Width,
		drawSize.Height*b.BlockSize.Height,
	)

	// Track height to know the position for the relation
	totalSize := model.Size{Width: 0, Height: 0}
	initialPos := b.currentPos

	for _, index := range table.Indexes {
		logrus.Infof("Drawing index %s", index.Name)
		relationSize := b.drawRelation(index)

		b.currentPos.X += relationSize.Width
		totalSize.AddWidthMaxHeight(relationSize)
	}

	if table.Toast != nil {
		toast := table.Toast
		logrus.Infof("Drawing toast %s", toast.Name)
		toastSize := b.drawRelation(toast.Relation)
		b.currentPos.X += toastSize.Width
		totalSize.AddWidthMaxHeight(toastSize)

		logrus.Infof("Drawing toast index %s at pos %v", toast.Index.Name, b.currentPos)
		toastIndexSize := b.drawRelation(toast.Index)
		b.currentPos.X += toastIndexSize.Width
		totalSize.AddWidthMaxHeight(toastIndexSize)
	}

	b.currentPos = initialPos
	b.currentPos.AddHeight(totalSize)

	logrus.Infof("Drawing table %s at pos %v", table.Name, b.currentPos)
	b.drawRelation(table.Relation)
}
