package postgresstore

import (
	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go/internal"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresMigration(db *DB, migrationPath string) (*rf.Migration, error) {
	return rf.NewMigration(stdlib.OpenDBFromPool(db.db), migrationPath)
}
