package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bonnefoa/pg_buffer_viz/pkg/bufferviz"
	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	db        *db.DbPool
	bufferViz *bufferviz.BufferViz
}

func newHttpServer(ctx context.Context) (*HttpServer, error) {
	dbConfig := db.GetDbConfigCli()
	dbConnection, err := db.NewDbPool(ctx, dbConfig.ConnectUrl)
	if err != nil {
		return nil, err
	}
	b := bufferviz.NewBufferViz(
		nil,
		util.GetBlockSize(),
		util.GetMarginSize(),
	)

	server := &HttpServer{bufferViz: &b, db: dbConnection}
	return server, nil
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

func StartHttpServer(ctx context.Context, h *HttpServerConfigCli) (*HttpServer, error) {
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
		return nil, eris.Wrapf(err, "Failed to listen on %s", h.ListenAddress)
	}
	server, err := newHttpServer(ctx)
	if err != nil {
		return nil, err
	}

	router := server.setupRouter()
	srv := &http.Server{
		Addr:    h.ListenAddress,
		Handler: router,
	}
	go startHttpServer(ctx, listener, srv)
	return server, nil
}
