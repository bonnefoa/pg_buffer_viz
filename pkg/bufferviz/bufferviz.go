package bufferviz

import (
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/sirupsen/logrus"
)

type BufferViz struct {
	canvas *svg.SVG

	BlockHeight int
	BlockWidth  int

	x int
	y int
}

func NewBufferViz(canvas *svg.SVG, blockHeight int, blockWidth int) BufferViz {
	b := BufferViz{BlockHeight: blockHeight, BlockWidth: blockWidth}
	b.canvas = canvas
	b.x = 0
	b.y = 0
	return b
}

func (b *BufferViz) GetSize(table db.Table) (int, int) {
	numBuffers := len(table.Fsm)
	blocksPerLine := int(math.Sqrt(float64(numBuffers))) + 1

	relationWidth := blocksPerLine * b.BlockWidth
	relationHeight := blocksPerLine * (b.BlockHeight + 1)
	return relationWidth + 2, relationHeight + 2
}

func (b *BufferViz) GetFsmColor(fsmValue int16) string {
	percent := int((float64(fsmValue) / 8192) * 100)
	return fmt.Sprintf("fill: color-mix(in srgb, green %d%%, red)", percent)
}

func (b *BufferViz) DrawName(blocksPerLine int, relation db.Relation) {
	relationWidth := blocksPerLine * b.BlockWidth
	xPos := b.x + relationWidth/2
	yPos := b.y + b.BlockHeight/2
	b.canvas.Text(xPos, yPos, relation.Name, "text-anchor:middle;font-size:20px")
}

func (b *BufferViz) DrawRelation(relation db.Relation) (int, int) {
	numBuffers := len(relation.Fsm)
	blocksPerLine := int(math.Sqrt(float64(numBuffers))) + 1
	lines := numBuffers / blocksPerLine

	b.DrawName(blocksPerLine, relation)

	xOffset := b.x
	yOffset := b.y + b.BlockHeight

	b.canvas.Gstyle("stroke-width:2;stroke:black;fill:white")
	for line := range blocksPerLine {
		for column := range blocksPerLine {
			bufno := line*blocksPerLine + column
			if bufno >= numBuffers {
				break
			}
			x := xOffset + column*b.BlockWidth
			y := yOffset + line*b.BlockHeight
			style := b.GetFsmColor(relation.Fsm[bufno])
			b.canvas.Rect(x, y, b.BlockWidth, b.BlockHeight, style)
		}
	}
	b.canvas.Gend()

	b.canvas.Gstyle("text-anchor:middle;font-size:20px;fill:black;dominant-baseline=middle")
	for line := range blocksPerLine {
		for column := range blocksPerLine {
			bufno := line*blocksPerLine + column
			if bufno > numBuffers {
				break
			}
			if bufno%50 == 0 {
				x := xOffset + column*b.BlockWidth + b.BlockWidth/2
				y := yOffset + line*b.BlockHeight + b.BlockHeight/2
				b.canvas.Text(x, y, fmt.Sprint(bufno))
			}
		}
	}
	b.canvas.Gend()
	return blocksPerLine * b.BlockWidth, lines * b.BlockHeight
}

func (b *BufferViz) DrawTable(table db.Table) {
	w, h := b.GetSize(table)
	b.canvas.Start(w, h)

	// Track height to know the position for the relation
	maxHeight := 0

	for _, index := range table.Indexes {
		logrus.Infof("Drawing index %s", index.Name)
		width, height := b.DrawRelation(index)
		b.x += width + b.BlockWidth
		if height > maxHeight {
			maxHeight = height
		}
	}

	//logrus.Infof("Drawing toast %s", index.Name)
	//width, height := b.DrawRelation(index)
	//b.x += width + b.BlockWidth
	//if height > maxHeight {
	//	maxHeight = height
	//}

	b.x = 0
	b.y = maxHeight + b.BlockHeight*2
	logrus.Infof("Drawing table %s", table.Name)
	b.DrawRelation(table.Relation)

	b.canvas.End()
}
