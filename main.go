package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/atticus-sullivan/freezerDB/api"
	"github.com/atticus-sullivan/freezerDB/cli"
	dbMod "github.com/atticus-sullivan/freezerDB/db"

	"gopkg.in/yaml.v3"
)

type ServerConf struct {
	Key  string `yaml:"key"`
	Addr string `yaml:"addr"`
}

type Args struct {
	Cli     *cli.CmdArgs `arg:"subcommand:cli"`
	Rest    *struct{}    `arg:"subcommand:rest"`
	RestDoc *struct{}    `arg:"subcommand:restDoc"`
}

// get the directory for the configuration of this project
func getCfgDir() string {
	dir, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		usr, _ := user.Current()
		dir = filepath.Join(usr.HomeDir, ".config")
	}
	return filepath.Join(dir, "freezer")
}

func main() {
	var args Args
	parser, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		panic(err)
	}

	if err := parser.Parse(os.Args[1:]); err != nil {
		switch err {
		case arg.ErrVersion:
			println("Version is not implemented")
			return
		case arg.ErrHelp:
			parser.WriteHelp(os.Stdout)
			return
		default:
			panic(err)
		}
	}

	// initialize database connection if needed
	var db *dbMod.DB
	if args.RestDoc == nil {
		// Setup database connection
		db, err = dbMod.NewDB(filepath.Join(getCfgDir(), "freezer.yaml"))
		if err != nil {
			panic(err)
		}
		defer db.Close()
	} else {
		db = &dbMod.DB{}
	}

	switch {
	case args.Cli != nil:
		cli.Cli(args.Cli, db)
		return

	case args.RestDoc != nil:
		s := api.CreateNewServer(&db.DB)
		s.MountHandlers("")
		s.Doc()

	case args.Rest != nil:
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		var serverConf ServerConf
		f, err := os.Open(filepath.Join(getCfgDir(), "server.yaml"))
		if err != nil {
			panic(err)
		}
		if err := yaml.NewDecoder(f).Decode(&serverConf); err != nil {
			panic(err)
		}

		s := api.CreateNewServer(&db.DB)
		s.MountHandlers(serverConf.Key)
		srv := &http.Server{
			Addr:    serverConf.Addr,
			Handler: s.Router,
		}
		go func(srv *http.Server) {
			fmt.Println(srv.ListenAndServe())
		}(srv)
		fmt.Printf("Server Started on %s\n", serverConf.Addr)

		// Block until a signal is received.
		<-c

		// shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel() // cancel context if returning
		}()
		if err := srv.Shutdown(ctx); err != nil {
			panic(err)
		}
		fmt.Println("Server Stopped")
	}
}
