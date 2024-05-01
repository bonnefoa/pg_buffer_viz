package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rotisserie/eris"
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
	Toast   *Toast
}

type Toast struct {
	Relation
	Index Relation
}

type RelationFreeSpace struct {
	Name string
	Fsm  []int
}

func Connect(ctx context.Context, connectUrl string) (*DbConnection, error) {
	logrus.Debugf("Connecting to PostgreSQL using connecturl '%v'", connectUrl)
	conn, err := pgx.Connect(ctx, connectUrl)
	if err != nil {
		return nil, eris.Wrap(err, "Connection error")
	}
	return &DbConnection{conn}, nil
}

func (d *DbConnection) FetchFsm(ctx context.Context, relationName string) ([]int16, error) {
	logrus.Debugf("Fetch FSM for relation '%s'", relationName)
	rows, err := d.Query(ctx, "select avail from pg_freespace($1)", relationName)
	if err != nil {
		return nil, eris.Wrap(err, "Fetch FSM failed")
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
	logrus.Debugf("Fetch indexes for relation '%s'", relationName)
	rows, err := d.Query(ctx, "select indexname from pg_indexes where tablename=$1", relationName)
	if err != nil {
		return nil, eris.Wrap(err, "Fetch index name failed")
	}
	indexNames, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, eris.Wrap(err, "Reading index failed")
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

func (d *DbConnection) FetchToast(ctx context.Context, relationName string) (*Toast, error) {
	logrus.Debugf("Fetch toast for relation '%s'", relationName)
	rows, err := d.Query(ctx, "SELECT relname FROM pg_class WHERE oid=(select reltoastrelid from pg_class where relname=$1);", relationName)
	if err != nil {
		return nil, eris.Wrap(err, "Toast query failed")
	}
	if !rows.Next() {
		rows.Close()
		return nil, nil
	}
	var toastRelationName string
	err = rows.Scan(&toastRelationName)
	if err != nil {
		rows.Close()
		return nil, eris.Wrap(err, "Scanning toast name failed")
	}
	rows.Close()

	toastIndexName := fmt.Sprintf("%s_index", toastRelationName)
	relation, err := d.FetchRelation(ctx, toastRelationName)
	if err != nil {
		return nil, eris.Wrap(err, "Fetch toast failed")
	}

	index, err := d.FetchRelation(ctx, toastIndexName)
	if err != nil {
		return nil, err
	}
	return &Toast{relation, index}, nil
}

func (d *DbConnection) FetchTable(ctx context.Context, relationName string) (table Table, err error) {
	logrus.Infof("Fetch buffer information for table '%s'", relationName)
	table.Relation, err = d.FetchRelation(ctx, relationName)
	if err != nil {
		return
	}
	table.Indexes, err = d.FetchIndexes(ctx, relationName)
	if err != nil {
		return
	}
	table.Toast, err = d.FetchToast(ctx, relationName)
	return
}
