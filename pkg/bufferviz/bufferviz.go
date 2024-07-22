package bufferviz

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/model"
	"github.com/bonnefoa/pg_buffer_viz/pkg/render"
	"github.com/sirupsen/logrus"
)

type BufferViz struct {
	canvas *svg.SVG

	BlockSize  model.Size
	MarginSize model.Size

	currentCoordinate model.Coordinate
}

func NewBufferViz(canvas *svg.SVG, blockSize model.Size, marginSize model.Size) BufferViz {
	b := BufferViz{
		canvas:            canvas,
		BlockSize:         blockSize,
		MarginSize:        marginSize,
		currentCoordinate: model.Coordinate{X: 1, Y: 1},
	}
	return b
}

func (b *BufferViz) SetCanvas(canvas *svg.SVG) {
	b.canvas = canvas
	b.currentCoordinate = model.Coordinate{X: 1, Y: 1}
}

func (b *BufferViz) getFsmColor(fsmValue int16) string {
	percent := int((float64(fsmValue) / 8192) * 100)
	return fmt.Sprintf("fill: color-mix(in srgb, green %d%%, red)", percent)
}

func (b *BufferViz) coordinateToPosition(c model.Coordinate) (x, y int) {
	return c.X * b.BlockSize.Width, c.Y * b.BlockSize.Height
}

func (b *BufferViz) drawName(relation model.Relation) {
	xPos := b.currentCoordinate.X * b.BlockSize.Width
	yPos := (float64(b.currentCoordinate.Y) + float64(0.5)) * float64(b.BlockSize.Height)
	b.canvas.Text(xPos, int(yPos), relation.Name, "text-align:left;font-size:10px")
}

func (b *BufferViz) drawRelation(relation model.Relation) model.Size {
	relationSize := relation.GetRelationSize()
	numBuffers := relation.GetNumbBuffers()

	b.drawName(relation)

	coordinate := b.currentCoordinate
	coordinate.Y += 1

	// Draw one rect per block
	for line := range relationSize.Width {
		for column := range relationSize.Width {
			bufno := line*relationSize.Width + column
			if bufno >= numBuffers {
				break
			}
			x := (coordinate.X + column) * b.BlockSize.Width
			y := (coordinate.Y + line) * b.BlockSize.Height
			fsmBucket := relation.Fsm[bufno] / 32
			blockId := fmt.Sprintf("id=\"%s_%d\"", relation.Name, bufno)

			b.canvas.Rect(x+2, y+2, b.BlockSize.Width-1, b.BlockSize.Height-1, blockId,
				fmt.Sprintf("class=\"block fsm%d\"", fsmBucket))
		}
	}
	relationSize.Add(b.MarginSize)

	return relationSize
}

func (b *BufferViz) DrawTable(table model.Table) {
	drawSize := b.getDrawSize(table)
	width := drawSize.Width * b.BlockSize.Width
	height := drawSize.Height * b.BlockSize.Height

	render.StartSVG(b.canvas, width, height)

	// Track height to know the position for the relation
	totalSize := model.Size{Width: 0, Height: 0}
	initialPos := b.currentCoordinate

	for _, index := range table.Indexes {
		logrus.Infof("Drawing index %s", index.Name)
		relationSize := b.drawRelation(index)

		b.currentCoordinate.X += relationSize.Width
		totalSize.AddWidthMaxHeight(relationSize)
	}

	if table.Toast != nil {
		toast := table.Toast
		logrus.Infof("Drawing toast %s", toast.Name)
		toastSize := b.drawRelation(toast.Relation)
		b.currentCoordinate.X += toastSize.Width
		totalSize.AddWidthMaxHeight(toastSize)

		logrus.Infof("Drawing toast index %s at coord %v", toast.Index.Name, b.currentCoordinate)
		toastIndexSize := b.drawRelation(toast.Index)
		b.currentCoordinate.X += toastIndexSize.Width
		totalSize.AddWidthMaxHeight(toastIndexSize)
	}

	b.currentCoordinate = initialPos
	b.currentCoordinate.AddHeight(totalSize)

	logrus.Infof("Drawing table %s at coord %v", table.Name, b.currentCoordinate)
	relationSize := b.drawRelation(table.Relation)
	b.currentCoordinate.AddHeight(relationSize)
}

func (b *BufferViz) AddFooter() {
	x, y := b.coordinateToPosition(b.currentCoordinate)
	b.canvas.Text(x, y, "Details: ", "id=\"details\"", "text-align:left;font-size:10px")
}
