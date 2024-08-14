package postgres

import (
	"context"
	"time"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db     *pgxpool.Pool
	ctx    context.Context
	cancel func()

	DBURL string
	Now   func() time.Time
}

func NewDB(dbURL string) *DB {
	db := &DB{
		DBURL: dbURL,
		Now:   time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
	if db.DBURL == "" {
		return rf.NewAppError(rf.ECIntenal, "db url required")
	}

	if db.db, err = pgxpool.New(db.ctx, db.DBURL); err != nil {
		return err
	}

	if err := db.db.Ping(db.ctx); err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	db.cancel()
	db.db.Close()
	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts pgx.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

type Tx struct {
	pgx.Tx
	db  *DB
	now time.Time
}
