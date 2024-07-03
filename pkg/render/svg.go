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

func StartSVG(s *svg.SVG, width int, height int) {
	s.Start(width, height, "onload=\"init(evt)\"")
	svg_css, err := os.ReadFile("resources/svg_css.css")
	util.FatalIf(err)
	svg_js, err := os.ReadFile("resources/svg_functions.js")
	util.FatalIf(err)
	s.Style("text/css", string(svg_css))
	s.Script("text/ecmascript", string(svg_js))

	offColors := make([]svg.Offcolor, 0)
	offColors = append(offColors, svg.Offcolor{Offset: 5, Color: "#eeeeee", Opacity: 1})
	offColors = append(offColors, svg.Offcolor{Offset: 95, Color: "#eeeeb0", Opacity: 1})
	s.Def()
	s.LinearGradient("background", 0, 0, 0, 100, offColors)
	s.DefEnd()
	s.Rect(0, 0, width, height, "fill=\"url(#background)\"")
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
