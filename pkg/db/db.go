package db

import (
	"context"

	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type DbConnection struct {
	*pgx.Conn
}

type Relation struct {
	Name string
	Fsm  []int16
}

type Table struct {
	Relation
	Indexes []Relation
	Toast   Toast
}

type Toast struct {
	Relation
	Index Relation
}

type RelationFreeSpace struct {
	Name string
	Fsm  []int
}

func Connect(ctx context.Context, connectUrl string) *DbConnection {
	logrus.Debugf("Connecting to PostgreSQL using connecturl '%v'", connectUrl)
	conn, err := pgx.Connect(ctx, connectUrl)
	util.FatalIf(err)

	return &DbConnection{conn}
}

func (d *DbConnection) FetchFsm(ctx context.Context, relationName string) []int16 {
	rows, err := d.Query(ctx, "select avail from pg_freespace($1)", relationName)
	util.FatalIf(err)
	avails, err := pgx.CollectRows(rows, pgx.RowTo[int16])
	util.FatalIf(err)
	return avails
}

func (d *DbConnection) FetchRelation(ctx context.Context, relationName string) Relation {
	avails := d.FetchFsm(ctx, relationName)
	r := Relation{
		Name: relationName,
		Fsm:  avails,
	}
	return r
}

func (d *DbConnection) FetchIndexes(ctx context.Context, relationName string) []Relation {
	rows, err := d.Query(ctx, "select indexname from pg_indexes where tablename=$1", relationName)
	util.FatalIf(err)
	indexNames, err := pgx.CollectRows(rows, pgx.RowTo[string])
	util.FatalIf(err)
	indexes := make([]Relation, 0)
	for _, indexName := range indexNames {
		r := d.FetchRelation(ctx, indexName)
		indexes = append(indexes, r)
	}
	return indexes
}

func (d *DbConnection) FetchTable(ctx context.Context, relationName string) Table {
	r := d.FetchRelation(ctx, relationName)
	indexes := d.FetchIndexes(ctx, relationName)
	return Table{Relation: r, Indexes: indexes}
}

// func (d *DbConnection) GetIndexes(ctx context.Context, relation string) RelationFreeSpace {
// 	rows, err := d.Query(ctx, "select blkno, avail from pg_freespace($1)", relation)
// 	util.FatalIf(err)
// 	fs, err := pgx.CollectRows(rows, pgx.RowToStructByName[FreeSpace])
// 	util.FatalIf(err)
//
// 	return RelationFreeSpace{
// 		Name: relation,
// 		Fsm:  fs,
// 	}
// }
