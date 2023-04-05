package models

import (
	"github.com/atticus-sullivan/freezerDB/db"
	"io"

	"github.com/jmoiron/sqlx"
)

type ItemTypes []ItemType
type ItemType struct {
	Name         string `json:"name" db:"name" arg:"-n,--name,required"`
	CategoryName string `json:"category_name" db:"category_name" arg:"-c,--cat,required"`
}

// writes a dot node table to the writer
func (itemTypes ItemTypes) WriteDot(w io.Writer) {
	w.Write([]byte(`
digraph structs {
	node [shape=plaintext] struct [label=<
		<table cellspacing="2" border="0" rows="*" columns="*">
		<tr>`))
	// write header
	w.Write([]byte("<td><B><U>Name</U></B></td>\n"))
	w.Write([]byte("<td><B>Category</B></td>\n"))
	w.Write([]byte("</tr>\n"))

	for _, i := range itemTypes {
		w.Write([]byte("<tr>"))

		// TODO html escape!!!
		// TODO formatter like newlines (<BR/>) every 15 char
		w.Write([]byte("<td>  "))
		w.Write([]byte(i.Name))
		w.Write([]byte("  </td>"))

		w.Write([]byte("<td>  "))
		w.Write([]byte(i.CategoryName))
		w.Write([]byte("  </td>"))

		w.Write([]byte("</tr>"))
	}
	w.Write([]byte("</table>>];\n"))
}

// InsertItemType inserts a new item type into the database
func (itemType *ItemType) Insert(db *db.DB) error {
	_, err := db.NamedExec("INSERT INTO item_types (name, category_name) VALUES (:name, :category_name);", itemType)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItemType deletes an item type from the database
func DeleteItemType(db *db.DB, name string) error {
	_, err := db.Exec("DELETE FROM item_types WHERE name=?;", name)
	if err != nil {
		return err
	}
	return nil
}

// UpdateItemType updates the given item type in the database with non-zero values only
func (itemType *ItemType) UpdateItemType(db *db.DB, idName string) error {
	// Prepare statement and execute with arguments
	query, args, err := sqlx.Named("UPDATE item_types SET category_name = :category_name, name = :name WHERE name = ?;", itemType)
	if err != nil {
		return err
	}

	query = db.Rebind(query)
	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(append(args, idName)...)
	return err
}

func GetAllItemTypes(db *db.DB) (ItemTypes, error) {
	var ret ItemTypes
	err := db.Select(&ret, "SELECT * FROM item_types;")
	if err != nil {
		return nil, err
	}
	return ret, nil
}
