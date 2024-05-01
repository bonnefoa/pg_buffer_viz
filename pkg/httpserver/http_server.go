package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bonnefoa/pg_buffer_viz/pkg/bufferviz"
	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/bonnefoa/pg_buffer_viz/pkg/render"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
)

type PgBufferVizHttpServer struct {
	Port     int
	dbConfig db.DbConfigCli
}

func (f *PgBufferVizHttpServer) readinessRoute(c *gin.Context) {
	c.String(http.StatusOK, "Ok")
}

func (f *PgBufferVizHttpServer) statsRoute(c *gin.Context) {
	stats := make([]string, 0)
	logrus.Debugf("Sending stats: %v", stats)
	c.JSON(http.StatusOK, stats)
}

func (f *PgBufferVizHttpServer) bufferVizRoute(c *gin.Context) {
	canvas := render.NewCanvas(c.Writer)
	b := bufferviz.NewBufferViz(canvas.SVG, 30, 20)
	logrus.Info(c.Params)
	tableName := c.Params.ByName("table")

	ctx := c.Request.Context()
	d, err := db.Connect(ctx, f.dbConfig.ConnectUrl)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	table, err := d.FetchTable(ctx, tableName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	b.DrawTable(table)
	c.Header("Content-Type", "image/svg+xml")
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		err := c.Errors.Last()
		errJson := eris.ToJSON(err.Unwrap(), true)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errJson)
		logrus.Error(eris.ToString(err, true))
	}
}

func (f *PgBufferVizHttpServer) setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	skipLogs := []string{
		"/health",
	}
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, skipLogs...))
	router.Use(ErrorHandler)
	router.Use(gin.Recovery())
	router.GET("/readiness", f.readinessRoute)
	router.GET("/stats", f.statsRoute)

	router.GET("/buffer_viz/:table", f.bufferVizRoute)

	return router
}

func startHttpServer(ctx context.Context, listener net.Listener, srv *http.Server) {
	go func() {
		logrus.Infof("Starting http server on %s", srv.Addr)
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error listening: %s", eris.ToString(err, true))
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %s", eris.ToString(err, true))
	}
	logrus.Info("Exiting http server")
}

func StartHttpServer(ctx context.Context, h *HttpServerConfigCli) (*PgBufferVizHttpServer, error) {
	logrus.Infof("Starting http server on listen address '%s'", h.ListenAddress)
	if h.ListenAddress == "" {
		return nil, nil
	}
	if h.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	listener, err := net.Listen("tcp", h.ListenAddress)
	if err != nil {
		return nil, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	dbConfig := db.GetDbConfigCli()
	f := PgBufferVizHttpServer{
		Port:     port,
		dbConfig: dbConfig,
	}

	router := f.setupRouter()
	srv := &http.Server{
		Addr:    h.ListenAddress,
		Handler: router,
	}
	go startHttpServer(ctx, listener, srv)
	return &f, nil
}
