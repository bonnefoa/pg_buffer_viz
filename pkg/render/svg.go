package render

import (
	"bufio"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
)

type Canvas struct {
	*svg.SVG
	file *os.File
	bw   *bufio.Writer
}

func NewCanvas(filename string) *Canvas {
	var c Canvas
	var err error

	c.file, err = os.Create(filename)
	util.FatalIf(err)
	c.bw = bufio.NewWriter(c.file)

	c.SVG = svg.New(c.bw)
	return &c
}

func (c *Canvas) End() {
	c.SVG.End()
	err := c.bw.Flush()
	util.FatalIf(err)
	err = c.file.Close()
	util.FatalIf(err)
}
