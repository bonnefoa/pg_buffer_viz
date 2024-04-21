package db

import (
	"context"

	"github.com/bonnefoa/pg_buffer_viz/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type DbConnection struct {
	*pgx.Conn
}

type FreeSpace struct {
	Blkno int
	Avail int
}

func SetDbConfigFlags(fs *pflag.FlagSet) {
	fs.String("connect-url", "", "Connection url to PostgreSQL db")
	fs.String("relation", "", "Target relation")
}

func Connect(ctx context.Context, connectUrl string) *DbConnection {
	logrus.Debugf("Connecting to PostgreSQL using connecturl '%v'", connectUrl)
	conn, err := pgx.Connect(ctx, connectUrl)
	util.FatalIf(err)

	return &DbConnection{conn}
}

func (d *DbConnection) FetchFSM(ctx context.Context, relation string) []FreeSpace {
	rows, err := d.Query(ctx, "select blkno, avail from pg_freespace($1)", relation)
	util.FatalIf(err)
	fs, err := pgx.CollectRows(rows, pgx.RowToStructByName[FreeSpace])
	util.FatalIf(err)

	return fs
}
