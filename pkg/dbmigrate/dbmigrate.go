package dbmigrate

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	// file is necessary to read migrations from filesystem.
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DbMigrate struct {
	db           *sql.DB
	migrationPah string
	driver       database.Driver
}

func New(db *sql.DB, migrationPah string) (*DbMigrate, error) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, err
	}
	return &DbMigrate{
		db:           db,
		migrationPah: migrationPah,
		driver:       driver,
	}, nil
}

func (m *DbMigrate) Up() error {
	srcURL := fmt.Sprintf("file://%s", m.migrationPah)

	mx, err := migrate.NewWithDatabaseInstance(srcURL, "sqlite", m.driver)
	if err != nil {
		return err
	}

	if err := mx.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
