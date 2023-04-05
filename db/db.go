package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// TODO config
const migrationsPath = "file://home/pi/freezerDB/migrations"

// extend by exec, query, prepare or embedd sql.DB so that this can be used directly
type DB struct {
	sqlx.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	// source.Register("file2", &file.File{})
	db_, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db := &DB{DB: *db_}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping DB... %v", err)
	}

	// run migration
	driver, err := mysql.WithInstance(db.DB.DB, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file:///home/pi/freezerDB/migrations", "mysql", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to apply migrations: %v", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
