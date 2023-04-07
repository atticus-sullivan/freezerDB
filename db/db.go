package db

import (
	"fmt"
	"os"
	"time"

	goMySql "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	migrateFile "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
)

type mySqlConf struct {
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Net            string `yaml:"net"`
	Addr           string `yaml:"addr"`
	DBName         string `yaml:"dbName"`
	MigrationsPath string `yaml:"migrationsPath"`
	Location       string `yaml:"location"`
	Loc            *time.Location
}

func newMySqlConf(fn string) (*mySqlConf, error) {
	r, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	d := yaml.NewDecoder(r)
	var c mySqlConf
	if err := d.Decode(&c); err != nil {
		return nil, err
	}

	c.Loc, err = time.LoadLocation(c.Location)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// extend by exec, query, prepare or embedd sql.DB so that this can be used directly
type DB struct {
	sqlx.DB
}

func NewDB(fn string) (*DB, error) {
	dbConf, err := newMySqlConf(fn)
	if err != nil {
		panic(err)
	}

	cfg := goMySql.NewConfig()
	cfg.User = dbConf.User
	cfg.Passwd = dbConf.Password
	cfg.Net = dbConf.Net
	cfg.Addr = dbConf.Addr
	cfg.DBName = dbConf.DBName
	cfg.Loc = dbConf.Loc
	cfg.MultiStatements = true
	cfg.Params = map[string]string{
		"charset":   "utf8mb4",
		"parseTime": "True",
	}

	db_, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db := &DB{DB: *db_}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping DB... %v", err)
	}

	// migration
	driver, err := mysql.WithInstance(db.DB.DB, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate instance: %v", err)
	}
	// open source
	src, err := (&migrateFile.File{}).Open("file://"+dbConf.MigrationsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open migration source: %w", err)
	}
    defer src.Close()
	m, err := migrate.NewWithInstance("file", src, "mysql", driver)
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
