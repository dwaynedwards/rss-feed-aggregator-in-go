package rf

import (
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
)

type Migration struct {
	db            *sql.DB
	migrationPath string
}

func NewMigration(db *sql.DB, migrationPath string) (*Migration, error) {
	if db == nil {
		return &Migration{}, errors.New("db is nil")
	}

	// goose.SetBaseFS(embedMigrations)

	return &Migration{
		db:            db,
		migrationPath: migrationPath,
	}, nil
}

func (m *Migration) Up() error {
	if err := goose.Up(m.db, m.migrationPath); err != nil {
		return err
	}
	return nil
}

func (m *Migration) Down() error {
	if err := goose.Down(m.db, m.migrationPath); err != nil {
		return err
	}
	return nil
}

func (m *Migration) Reset() error {
	if err := goose.Reset(m.db, m.migrationPath); err != nil {
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
