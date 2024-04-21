package db

import (
	"context"

	"github.com/bonnefoa/pg_buffer_viz/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func SetDbConfigFlags(fs *pflag.FlagSet) {
	fs.String("connect-url", "", "Connection url to PostgreSQL db")
}

func Connect(ctx context.Context, connectUrl string) *pgx.Conn {
	logrus.Debugf("Connecting to PostgreSQL using connecturl '%v'", connectUrl)
	conn, err := pgx.Connect(ctx, connectUrl)
	util.FatalIf(err)

	return conn
}
