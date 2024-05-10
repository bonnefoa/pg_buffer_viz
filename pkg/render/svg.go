package render

import (
	"bufio"
	"io"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
	"github.com/sirupsen/logrus"
)

type CanvasFile struct {
	*svg.SVG
	file *os.File
	bw   *bufio.Writer
}

type CanvasIo struct {
	*svg.SVG
}

type Canvas interface {
	End()
}

func NewCanvasFile(filename string) *CanvasFile {
	var c CanvasFile
	var err error

	c.file, err = os.Create(filename)
	util.FatalIf(err)
	c.bw = bufio.NewWriter(c.file)

	c.SVG = svg.New(c.bw)
	return &c
}

func NewCanvasIo(w io.Writer) *CanvasIo {
	c := CanvasIo{svg.New(w)}
	return &c
}

func (c *CanvasIo) End() {
	c.SVG.End()
}

func (c *CanvasFile) End() {
	logrus.Infof("Flushing file")
	c.SVG.End()
	err := c.bw.Flush()
	util.FatalIf(err)
	err = c.file.Close()
	util.FatalIf(err)
}
