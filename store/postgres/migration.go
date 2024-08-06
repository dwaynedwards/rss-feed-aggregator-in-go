package postgres

import (
	"embed"

	rf "github.com/dwaynedwards/rss-feed-aggregator-in-go"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewMigration(db *DB, printNoop bool) (*rf.Migration, error) {
	return rf.NewMigration("postgres", stdlib.OpenDBFromPool(db.db), embedMigrations, printNoop)
}
