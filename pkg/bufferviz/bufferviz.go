package bufferviz

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/sirupsen/logrus"
)

type coordinate struct {
	x int
	y int
}

type BufferViz struct {
	canvas *svg.SVG

	BlockSize  Size
	MarginSize Size

	currentPos coordinate
}

func NewBufferViz(canvas *svg.SVG, blockSize Size, marginSize Size) BufferViz {
	b := BufferViz{
		canvas:     canvas,
		BlockSize:  blockSize,
		MarginSize: marginSize,
		currentPos: coordinate{x: 0, y: 0},
	}
	return b
}

func (b *BufferViz) getRelationSize(relation db.Relation) Size {
	numBuffers := float64(len(relation.Fsm))
	blocksPerLine := math.Ceil(math.Sqrt(float64(numBuffers)))
	lines := math.Ceil(numBuffers / blocksPerLine)
	res := Size{int(blocksPerLine), int(lines)}
	res.add(b.MarginSize)
	return res
}

func (b *BufferViz) getAncillarySizes(table db.Table) (sizes []Size) {
	for _, index := range table.Indexes {
		sizes = append(sizes, b.getRelationSize(index))
	}

	if table.Toast != nil {
		toast := table.Toast
		sizes = append(sizes, b.getRelationSize(toast.Relation))
		sizes = append(sizes, b.getRelationSize(toast.Index))
	}

	return sizes
}

func compareX(a Size, b Size) int {
	return cmp.Compare(a.Width, b.Height)
}

func compareY(a Size, b Size) int {
	return cmp.Compare(a.Width, b.Height)
}

func (b *BufferViz) getDrawSize(table db.Table) (res Size) {
	relationSize := b.getRelationSize(table.Relation)
	ancillarySizes := b.getAncillarySizes(table)

	// Names take 2 block height
	res.Height = relationSize.Height + 2*b.BlockSize.Height
	res.Height += slices.MaxFunc(ancillarySizes, compareY).Height + 2*b.BlockSize.Height

	bottomWidth := res.Width
	topWidth := len(ancillarySizes)
	for _, c := range ancillarySizes {
		topWidth += c.Width
	}

	res.Width = bottomWidth
	if topWidth > bottomWidth {
		res.Width = topWidth
	}

	return res
}

func (b *BufferViz) getFsmColor(fsmValue int16) string {
	percent := int((float64(fsmValue) / 8192) * 100)
	return fmt.Sprintf("fill: color-mix(in srgb, green %d%%, red)", percent)
}

func (b *BufferViz) drawName(blocksPerLine int, relation db.Relation) {
	relationWidth := blocksPerLine * b.BlockSize.Width
	xPos := b.currentPos.x + relationWidth/2
	yPos := b.currentPos.y + b.BlockSize.Height/2
	b.canvas.Text(xPos, yPos, relation.Name, "text-anchor:middle;font-size:20px")
}

func (b *BufferViz) drawRelation(relation db.Relation) (int, int) {
	numBuffers := len(relation.Fsm)
	blocksPerLine := int(math.Sqrt(float64(numBuffers))) + 1
	lines := numBuffers / blocksPerLine

	b.drawName(blocksPerLine, relation)

	pos := b.currentPos
	pos.y += b.BlockSize.Height

	b.canvas.Gstyle("stroke-width:2;stroke:black;fill:white")
	for line := range blocksPerLine {
		for column := range blocksPerLine {
			bufno := line*blocksPerLine + column
			if bufno >= numBuffers {
				break
			}
			x := pos.x + column*b.BlockSize.Width
			y := pos.y + line*b.BlockSize.Height
			style := b.getFsmColor(relation.Fsm[bufno])
			b.canvas.Rect(x, y, b.BlockSize.Width, b.BlockSize.Height, style)
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
				x := pos.x + column*b.BlockSize.Width + b.BlockSize.Width/2
				y := pos.y + line*b.BlockSize.Height + b.BlockSize.Height/2
				b.canvas.Text(x, y, fmt.Sprint(bufno))
			}
		}
	}
	b.canvas.Gend()
	return blocksPerLine * b.BlockSize.Width, lines * b.BlockSize.Height
}

func (b *BufferViz) DrawTable(table db.Table) {
	drawSize := b.getDrawSize(table)
	b.canvas.Start(drawSize.Width*b.BlockSize.Width, drawSize.Height*b.BlockSize.Height)

	// Track height to know the position for the relation
	maxHeight := 0

	for _, index := range table.Indexes {
		logrus.Infof("Drawing index %s", index.Name)
		width, height := b.drawRelation(index)
		b.currentPos.x += width + b.BlockSize.Width
		if height > maxHeight {
			maxHeight = height
		}
	}

	if table.Toast != nil {
		toast := table.Toast
		logrus.Infof("Drawing toast %s", toast.Name)
		width, height := b.drawRelation(toast.Relation)
		b.currentPos.x += width + b.MarginSize.Width
		if height > maxHeight {
			maxHeight = height
		}

		logrus.Infof("Drawing toast index %s", toast.Index.Name)
		width, height = b.drawRelation(toast.Index)
		b.currentPos.x += width + b.MarginSize.Width
		if height > maxHeight {
			maxHeight = height
		}
	}

	b.currentPos.x = 0
	b.currentPos.y = maxHeight + b.MarginSize.Height
	logrus.Infof("Drawing table %s", table.Name)
	b.drawRelation(table.Relation)

	b.canvas.End()
}
