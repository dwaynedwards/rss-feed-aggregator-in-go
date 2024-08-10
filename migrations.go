package rf

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/pressly/goose/v3"
)

type Migration struct {
	db *sql.DB
}

func NewMigration(db *sql.DB, embedMigrations embed.FS, printLog bool) (*Migration, error) {
	if db == nil {
		return &Migration{}, errors.New("db is nil")
	}

	goose.SetBaseFS(embedMigrations)
	if !printLog {
		goose.SetLogger(goose.NopLogger())
	}

	return &Migration{db: db}, nil
}

func (m *Migration) Up() error {
	if err := goose.Up(m.db, "migrations"); err != nil {
		return err
	}
	return nil
}

func (m *Migration) Down() error {
	if err := goose.Down(m.db, "migrations"); err != nil {
		return err
	}
	return nil
}

func (m *Migration) Reset() error {
	if err := goose.Reset(m.db, "migrations"); err != nil {
		return err
	}
	return nil
}

func (m *Migration) Close() error {
	if err := m.db.Close(); err != nil {
		return err
	}
	return nil
}
