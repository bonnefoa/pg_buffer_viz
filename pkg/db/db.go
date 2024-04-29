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

func (d *DbConnection) FetchFsm(ctx context.Context, relationName string) ([]int16, error) {
	rows, err := d.Query(ctx, "select avail from pg_freespace($1)", relationName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[int16])
}

func (d *DbConnection) FetchRelation(ctx context.Context, relationName string) (Relation, error) {
	avails, err := d.FetchFsm(ctx, relationName)
	r := Relation{
		Name: relationName,
		Fsm:  avails,
	}
	return r, err
}

func (d *DbConnection) FetchIndexes(ctx context.Context, relationName string) ([]Relation, error) {
	rows, err := d.Query(ctx, "select indexname from pg_indexes where tablename=$1", relationName)
	if err != nil {
		return nil, err
	}
	indexNames, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	indexes := make([]Relation, 0)
	for _, indexName := range indexNames {
		r, err := d.FetchRelation(ctx, indexName)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, r)
	}
	return indexes, nil
}

func (d *DbConnection) FetchTable(ctx context.Context, relationName string) (table Table, err error) {
	logrus.Infof("Fetch buffer information for table '%s'", relationName)
	r, err := d.FetchRelation(ctx, relationName)
	if err != nil {
		return
	}
	table.Relation = r
	indexes, err := d.FetchIndexes(ctx, relationName)
	if err != nil {
		return
	}
	table.Indexes = indexes
	return
}
