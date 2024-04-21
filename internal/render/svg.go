package render

import (
	"bufio"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/internal/util"
)

type Canvas struct {
	*svg.SVG
	file *os.File
	bw   *bufio.Writer
}

type CanvasOptions struct {
	Width    int
	Height   int
	FileName string
}

func Start(options CanvasOptions) *Canvas {
	var c Canvas
	var err error

	c.file, err = os.Create(options.FileName)
	util.FatalIf(err)
	c.bw = bufio.NewWriter(c.file)

	c.SVG = svg.New(c.bw)

	c.Start(options.Width, options.Width)
	width := options.Width
	height := options.Height
	c.Circle(width/2, height/2, 100)
	c.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")

	return &c
}

func (c *Canvas) End() {
	c.SVG.End()
	err := c.bw.Flush()
	util.FatalIf(err)
	err = c.file.Close()
	util.FatalIf(err)
}
