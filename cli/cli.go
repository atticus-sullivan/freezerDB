package cli

import (
	"fmt"
	"github.com/atticus-sullivan/freezerDB/db"
	"github.com/atticus-sullivan/freezerDB/db/models"
	"os"

	"github.com/alexflint/go-arg"
)

type CmdArgs struct {
	Cat  *catArgs  `arg:"subcommand:cat"`
	Type *typeArgs `arg:"subcommand:type"`
	Item *itemArgs `arg:"subcommand:item"`
}

type catArgs struct {
	AddArgs *struct {
		models.Category
	} `arg:"subcommand:add"`
	LsArgs *struct {
		Name string `arg:"-n,--name"`
	} `arg:"subcommand:ls"`
	RmArgs *struct {
		Name string `arg:"-n,--name,required"`
	} `arg:"subcommand:rm"`
}

type typeArgs struct {
	AddArgs *struct {
		models.ItemType
	} `arg:"subcommand:add"`
	LsArgs *struct {
		Name string `arg:"-n,--name"`
	} `arg:"subcommand:ls"`
	RmArgs *struct {
		Name string `arg:"-n,--name,required"`
	} `arg:"subcommand:rm"`
}
type itemArgs struct {
	AddArgs *struct {
		models.FreezerItem
	} `arg:"subcommand:add"`
	LsArgs *struct {
		ID uint `arg:"-i,--id" default:"0"`
	} `arg:"subcommand:ls"`
	RmArgs *struct {
		ID uint `arg:"-i,--id,required"`
	} `arg:"subcommand:rm"`
}

// might panic
func Cli(vargs []string, db *db.DB) {
	args := CmdArgs{}
	parser, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		panic(err)
	}
	err = parser.Parse(vargs)
	switch err {
	case arg.ErrHelp:
		err := parser.WriteHelpForSubcommand(os.Stdout, parser.SubcommandNames()...)
		if err != nil {
			panic(err)
		}
		return
	case arg.ErrVersion:
		fmt.Println("version not implemented")
		return
	default:
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO are there arguments which need to be checked if they are valid?

	switch {
	case args.Cat != nil:
		switch {
		case args.Cat.AddArgs != nil:
			a := args.Cat.AddArgs
			err := a.Category.Insert(db)
			if err != nil {
				panic(err)
			}
		case args.Cat.LsArgs != nil:
			a := args.Cat.LsArgs
			if a.Name == "" {
				rows, err := db.Queryx("select * from categories;")
				if err != nil {
					panic(err)
				}
				defer rows.Close()
				var cat models.Category
				for rows.Next() {
					if err := rows.StructScan(&cat); err != nil {
						panic(err)
					}
					fmt.Println(&cat)
				}
			} else {
				rows, err := db.Queryx("select * from categories where name = ?;", a.Name)
				if err != nil {
					panic(err)
				}
				defer rows.Close()
				var cat models.Category
				for rows.Next() {
					if err := rows.StructScan(&cat); err != nil {
						panic(err)
					}
					fmt.Println(&cat)
				}
			}
		case args.Cat.RmArgs != nil:
			a := args.Cat.RmArgs
			err := models.DeleteCategory(db, a.Name)
			if err != nil {
				panic(err)
			}
		}

	case args.Item != nil:
		switch {
		case args.Item.AddArgs != nil:
			a := args.Item.AddArgs
			err := a.FreezerItem.Insert(db)
			if err != nil {
				panic(err)
			}
		case args.Item.LsArgs != nil:
			a := args.Item.LsArgs
			if a.ID == 0 {
				rows, err := db.Queryx("select * from freezer_items;")
				if err != nil {
					panic(err)
				}
				defer rows.Close()
				var item models.FreezerItem
				for rows.Next() {
					if err := rows.StructScan(&item); err != nil {
						panic(err)
					}
					fmt.Println(&item)
				}
			} else {
				rows, err := db.Queryx("select * from freezer_items where id = ?;", a.ID)
				if err != nil {
					panic(err)
				}
				defer rows.Close()
				var item models.FreezerItem
				for rows.Next() {
					if err := rows.StructScan(&item); err != nil {
						panic(err)
					}
					fmt.Println(&item)
				}
			}
		case args.Item.RmArgs != nil:
			a := args.Item.RmArgs
			err := models.DeleteFreezerItem(db, a.ID)
			if err != nil {
				panic(err)
			}
		}
	case args.Type != nil:
		switch {
		case args.Type.AddArgs != nil:
			a := args.Type.AddArgs
			err := a.ItemType.Insert(db)
			if err != nil {
				panic(err)
			}
		case args.Type.LsArgs != nil:
			a := args.Type.LsArgs
			if a.Name == "" {
				rows, err := db.Queryx("select * from item_types;")
				if err != nil {
					panic(err)
				}
				defer rows.Close()
				var cat models.ItemType
				for rows.Next() {
					if err := rows.StructScan(&cat); err != nil {
						panic(err)
					}
					fmt.Println(&cat)
				}
			} else {
				rows, err := db.Queryx("select * from item_types where name = ?;", a.Name)
				if err != nil {
					panic(err)
				}
				defer rows.Close()
				var cat models.Category
				for rows.Next() {
					if err := rows.StructScan(&cat); err != nil {
						panic(err)
					}
					fmt.Println(&cat)
				}
			}
		case args.Type.RmArgs != nil:
			a := args.Type.RmArgs
			err := models.DeleteItemType(db, a.Name)
			if err != nil {
				panic(err)
			}
		}
	}
}
