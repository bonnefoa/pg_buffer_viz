package httpserver

import (
	"net/http"

	"github.com/bonnefoa/pg_buffer_viz/pkg/bufferviz"
	"github.com/bonnefoa/pg_buffer_viz/pkg/render"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
)

func (s *HttpServer) readinessRoute(c *gin.Context) {
	c.String(http.StatusOK, "Ok")
}

func (s *HttpServer) statsRoute(c *gin.Context) {
	stats := make([]string, 0)
	logrus.Debugf("Sending stats: %v", stats)
	c.JSON(http.StatusOK, stats)
}

func (s *HttpServer) bufferVizRoute(c *gin.Context) {
	canvas := render.NewCanvas(c.Writer)
	b := bufferviz.NewBufferViz(canvas.SVG, 30, 20)
	logrus.Info(c.Params)
	tableName := c.Params.ByName("table")

	ctx := c.Request.Context()
	table, err := s.db.FetchTable(ctx, tableName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	b.DrawTable(table)
	c.Header("Content-Type", "image/svg+xml")
}

func (s *HttpServer) listRelations(c *gin.Context) {
	relations, err := s.db.ListRelationNames(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"relations": relations,
	})
}

func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		err := c.Errors.Last()
		errJson := eris.ToJSON(err.Unwrap(), true)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errJson)
		logrus.Error(eris.ToString(err, true))
	}
}

func (s *HttpServer) setupRouter() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.Use(gin.Recovery())
	skipLogs := []string{
		"/health",
	}
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, skipLogs...))
	router.Use(errorHandler)
	router.Use(gin.Recovery())
	router.GET("/readiness", s.readinessRoute)
	router.GET("/stats", s.statsRoute)

	router.GET("/", s.listRelations)
	router.GET("/buffer_viz/:table", s.bufferVizRoute)

	return router
}
