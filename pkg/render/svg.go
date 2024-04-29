package render

import (
	"bufio"
	"io"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
)

type FileCanvas struct {
	*svg.SVG
	file *os.File
	bw   *bufio.Writer
}

type Canvas struct {
	*svg.SVG
}

func NewFileCanvas(filename string) *FileCanvas {
	var c FileCanvas
	var err error

	c.file, err = os.Create(filename)
	util.FatalIf(err)
	c.bw = bufio.NewWriter(c.file)

	c.SVG = svg.New(c.bw)
	return &c
}

func NewCanvas(w io.Writer) *Canvas {
	var c Canvas
	c.SVG = svg.New(w)
	return &c
}

func (c *Canvas) End() {
	c.SVG.End()
}

func (c *FileCanvas) End() {
	c.SVG.End()
	err := c.bw.Flush()
	util.FatalIf(err)
	err = c.file.Close()
	util.FatalIf(err)
}
