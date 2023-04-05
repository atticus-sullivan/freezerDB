package main

import (
	"github.com/atticus-sullivan/freezerDB/cli"
	"github.com/atticus-sullivan/freezerDB/db"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	goMySql "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

// get the directory for the configuration of this project
func getCfgDir() string {
	dir, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		usr, _ := user.Current()
		dir = filepath.Join(usr.HomeDir, ".config")
	}
	return filepath.Join(dir, "freezer")
}

type mySqlConf struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Net      string `yaml:"net"`
	Addr     string `yaml:"addr"`
	DBName   string `yaml:"dbName"`
	Location string `yaml:"location"`
	Loc      *time.Location
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

func main() {
	dbConf, err := newMySqlConf(filepath.Join(getCfgDir(), "freezer.yaml"))
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
	cfg.Params = map[string]string{
		"charset":   "utf8mb4",
		"parseTime": "True",
	}
	// Setup database connection
	db,err := db.NewDB(cfg.FormatDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	cli.Cli(os.Args[1:], db)
}
