package bufferviz

import (
	"fmt"
	"math"

	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/bonnefoa/pg_buffer_viz/pkg/render"
)

func GetSize(co render.CanvasOptions, fsr db.RelationFreeSpace) (int, int) {
	numBuffers := len(fsr.FreeSpace)
	blocksPerLine := int(math.Sqrt(float64(numBuffers))) + 1

	relationWidth := blocksPerLine * co.BlockWidth
	relationHeight := blocksPerLine * (co.BlockHeight + 1)
	return relationWidth + 2, relationHeight + 2
}

func DrawRelation(c *render.Canvas, fsr db.RelationFreeSpace) {
	o := c.Options

	numBuffers := len(fsr.FreeSpace)
	blocksPerLine := int(math.Sqrt(float64(numBuffers))) + 1

	relationWidth := blocksPerLine * o.BlockWidth

	c.Text(relationWidth/2, o.BlockHeight/2, fsr.Name, "text-anchor:middle;font-size:20px;fill:black")

	xOffset := 0
	yOffset := o.BlockHeight

	c.Gstyle("stroke-width:2;stroke:black;fill:white")
	for line := range blocksPerLine {
		for column := range blocksPerLine {
			bufno := line*blocksPerLine + column
			if bufno > numBuffers {
				break
			}
			x := xOffset + column*o.BlockWidth
			y := yOffset + line*o.BlockHeight
			c.Rect(x, y, o.BlockWidth, o.BlockHeight)
		}
	}
	c.Gend()

	c.Gstyle("text-anchor:middle;font-size:20px;fill:black;dominant-baseline=middle")
	for line := range blocksPerLine {
		for column := range blocksPerLine {
			bufno := line*blocksPerLine + column
			if bufno > numBuffers {
				break
			}
			if bufno%50 == 0 {
				x := xOffset + column*o.BlockWidth + o.BlockWidth/2
				y := yOffset + line*o.BlockHeight + o.BlockHeight/2
				c.Text(x, y, fmt.Sprint(bufno))
			}

		}
	}
	c.Gend()
}
