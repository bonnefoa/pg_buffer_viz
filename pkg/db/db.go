package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
)

type DbPool struct {
	*pgxpool.Pool
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

func NewDbPool(ctx context.Context, connectUrl string) (*DbPool, error) {
	config, err := pgxpool.ParseConfig(connectUrl)
	if err != nil {
		return nil, eris.Wrap(err, "Error parsing db configuration")
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, eris.Wrap(err, "Error creating pgxpool")
	}
	return &DbPool{pool}, nil
}

func (d *DbPool) FetchFsmFromOid(ctx context.Context, oid uint32) ([]int16, error) {
	logrus.Debugf("Fetch FSM for oid '%d'", oid)
	rows, err := d.Query(ctx, "select avail from pg_freespace($1)", oid)
	if err != nil {
		return nil, eris.Wrap(err, "Fetch FSM failed")
	}
	return pgx.CollectRows(rows, pgx.RowTo[int16])
}

func (d *DbPool) FetchFsm(ctx context.Context, relationName string) ([]int16, error) {
	logrus.Debugf("Fetch FSM for relation '%s'", relationName)
	rows, err := d.Query(ctx, "select avail from pg_freespace($1)", relationName)
	if err != nil {
		return nil, eris.Wrap(err, "Fetch FSM failed")
	}
	return pgx.CollectRows(rows, pgx.RowTo[int16])
}

func (d *DbPool) FetchRelationFromOid(ctx context.Context, relationName string, oid uint32) (Relation, error) {
	avails, err := d.FetchFsmFromOid(ctx, oid)
	r := Relation{
		Name: relationName,
		Fsm:  avails,
	}
	return r, err
}

func (d *DbPool) FetchRelation(ctx context.Context, relationName string) (Relation, error) {
	avails, err := d.FetchFsm(ctx, relationName)
	r := Relation{
		Name: relationName,
		Fsm:  avails,
	}
	return r, err
}

func (d *DbPool) ListRelationNames(ctx context.Context) ([]string, error) {
	rows, err := d.Query(ctx, "select relname from pg_class where relkind='r' order by oid desc")
	if err != nil {
		return nil, eris.Wrap(err, "Error fetching the list of relation names")
	}
	relationNames, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, eris.Wrap(err, "Error collecting rows for relation names")
	}
	return relationNames, err
}

func (d *DbPool) FetchIndexes(ctx context.Context, relationName string) ([]Relation, error) {
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

type ToastResponse struct {
	ToastOid     uint32
	RelationName string
	IndexOid     uint32
	IndexName    string
}

func (d *DbPool) FetchToast(ctx context.Context, relationName string) (*Toast, error) {
	logrus.Debugf("Fetch toast for relation '%s'", relationName)
	rows, err := d.Query(ctx, `WITH toast_ids AS (
    SELECT c.reltoastrelid as oid, i.indexrelid as idx_oid
    FROM pg_class c, pg_index i
    WHERE relname=$1
    AND i.indrelid = c.reltoastrelid
) SELECT t_oids.oid, t.relname, t_oids.idx_oid, ti.relname
FROM pg_class t, pg_class ti, toast_ids t_oids
WHERE t.oid = t_oids.oid AND ti.oid = t_oids.idx_oid`, relationName)
	if err != nil {
		return nil, eris.Wrap(err, "Toast query failed")
	}

	toastResponse, err := pgx.CollectOneRow[ToastResponse](rows, pgx.RowTo[ToastResponse])
	if err != nil {
		// No toast found
		return nil, nil
	}

	relation, err := d.FetchRelationFromOid(ctx, toastResponse.RelationName, toastResponse.ToastOid)
	if err != nil {
		return nil, err
	}
	index, err := d.FetchRelationFromOid(ctx, toastResponse.IndexName, toastResponse.IndexOid)
	if err != nil {
		return nil, err
	}
	return &Toast{relation, index}, nil
}

func (d *DbPool) FetchTable(ctx context.Context, relationName string) (table Table, err error) {
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
