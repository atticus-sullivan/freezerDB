package cli

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
	"github.com/atticus-sullivan/freezerDB/db"
	"github.com/atticus-sullivan/freezerDB/db/models"
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

// allowed to panic
func Cli(args *CmdArgs, db *db.DB) {
	switch {
	case args.Cat != nil:
		handleCat(args.Cat, db)
	case args.Item != nil:
		handleItem(args.Item, db)
	case args.Type != nil:
		handleType(args.Type, db)
	}
}

// allowed to panic
func handleType(args *typeArgs, db *db.DB) {
	switch {
	case args.AddArgs != nil:
		a := args.AddArgs
		_, err := db.NamedExec("INSERT INTO item_types (name, category_name) VALUES (:name, :category_name)", a)
		if err != nil {
			panic(err)
		}
	case args.LsArgs != nil:
		a := args.LsArgs
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
				fmt.Printf("%+v\n", &cat)
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
				fmt.Printf("%+v\n", &cat)
			}
		}
	case args.RmArgs != nil:
		a := args.RmArgs
		_, err := db.Exec("DELETE FROM item_types WHERE name = ?", a.Name)
		if err != nil {
			panic(err)
		}
	}
}

// allowed to panic
func handleItem(args *itemArgs, db *db.DB) {
	switch {
	case args.AddArgs != nil:
		a := args.AddArgs
		_, err := db.NamedExec("INSERT INTO freezer_items (date, identifier, amount, misc, item_name) VALUES (:date, :identifier, :amount, :misc, :item_name)", a)
		if err != nil {
			panic(err)
		}
	case args.LsArgs != nil:
		a := args.LsArgs
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
				fmt.Printf("%+v\n", &item)
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
				fmt.Printf("%+v\n", &item)
			}
		}
	case args.RmArgs != nil:
		a := args.RmArgs
		_, err := db.Exec("DELETE FROM freezer_items WHERE id = ?", a.ID)
		if err != nil {
			panic(err)
		}
	}
}

// allowed to panic
func handleCat(args *catArgs, db *db.DB) {
	switch {
	case args.AddArgs != nil:
		a := args.AddArgs
		_, err := db.NamedExec("INSERT INTO categories (name) VALUES (:name)", a)
		if err != nil {
			panic(err)
		}
	case args.LsArgs != nil:
		a := args.LsArgs
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
				fmt.Printf("%+v\n", &cat)
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
				fmt.Printf("%+v\n", &cat)
			}
		}
	case args.RmArgs != nil:
		a := args.RmArgs
		_, err := db.Exec("DELETE FROM categories WHERE name=?", a.Name)
		if err != nil {
			panic(err)
		}
	}
}
