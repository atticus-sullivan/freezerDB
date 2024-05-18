package db

// Copyright (c) 2023, Lukas Heindl
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
