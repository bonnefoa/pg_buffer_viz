package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PgBufferVizHttpServer struct {
	Port int
}

func (f *PgBufferVizHttpServer) readinessRoute(c *gin.Context) {
	c.String(http.StatusOK, "Ok")
}

func (f *PgBufferVizHttpServer) statsRoute(c *gin.Context) {
	stats := make([]string, 0)
	logrus.Debugf("Sending stats: %v", stats)
	c.JSON(http.StatusOK, stats)
}

func (f *PgBufferVizHttpServer) setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	skipLogs := []string{
		"/health",
	}
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, skipLogs...))
	router.Use(gin.Recovery())
	router.GET("/readiness", f.readinessRoute)
	router.GET("/stats", f.statsRoute)

	router.GET("/buffer_viz", f.statsRoute)

	return router
}

func startHttpServer(ctx context.Context, listener net.Listener, srv *http.Server) {
	go func() {
		logrus.Infof("Starting http server on %s", srv.Addr)
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error listening: %s", err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %s", err)
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
	f := PgBufferVizHttpServer{
		Port: port,
	}
	router := f.setupRouter()
	srv := &http.Server{
		Addr:    h.ListenAddress,
		Handler: router,
	}
	go startHttpServer(ctx, listener, srv)
	return &f, nil
}
