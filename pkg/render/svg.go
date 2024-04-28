package render

import (
	"bufio"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
)

type CanvasOptions struct {
	FileName string

	BlockHeight int
	BlockWidth  int
}

type Canvas struct {
	*svg.SVG
	Options CanvasOptions
	file    *os.File
	bw      *bufio.Writer
}

func Start(options CanvasOptions, width int, height int) *Canvas {
	var c Canvas
	var err error

	c.file, err = os.Create(options.FileName)
	util.FatalIf(err)
	c.bw = bufio.NewWriter(c.file)

	c.SVG = svg.New(c.bw)
	c.Options = options

	c.Start(width, height)
	return &c
}

func (c *Canvas) End() {
	c.SVG.End()
	err := c.bw.Flush()
	util.FatalIf(err)
	err = c.file.Close()
	util.FatalIf(err)
}
