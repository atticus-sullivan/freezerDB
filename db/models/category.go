package models

import (
	"github.com/atticus-sullivan/freezerDB/db"
	"io"

	"github.com/jmoiron/sqlx"
)

type Categories []Category
type Category struct {
	Name string `json:"name" db:"name" arg:"-n,--name,required"`
}

// writes a dot node table to the writer
func (categories Categories) WriteDot(w io.Writer) {
	w.Write([]byte(`
digraph structs {
	node [shape=plaintext] struct [label=<
		<table cellspacing="2" border="0" rows="*" columns="*">
		<tr>`))
	// write header
	w.Write([]byte("<td><B><U>Name</U></B></td>\n"))
	w.Write([]byte("</tr>\n"))

	for _, c := range categories {
		w.Write([]byte("<tr>"))

		// TODO html escape!!!
		// TODO formatter like newlines (<BR/>) every 15 char
		w.Write([]byte("<td>  "))
		w.Write([]byte(c.Name))
		w.Write([]byte("  </td>"))

		w.Write([]byte("</tr>"))
	}
	w.Write([]byte("</table>>];\n"))
}

// Insert inserts a new category into the database.
func (category *Category) Insert(db *db.DB) error {
	_, err := db.NamedExec("INSERT INTO categories (name) VALUES (:name);", category)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCategory deletes a category from the database.
func DeleteCategory(db *db.DB, name string) error {
	_, err := db.Exec("DELETE FROM categories WHERE name = ?;", name)
	if err != nil {
		return err
	}
	return nil
}

// think about updating with whole category
// Update updates the given category in the database with non-zero values only
func (category *Category) Update(db *db.DB, idName string) error {
	// Prepare statement and execute with arguments
	query, args, err := sqlx.Named("UPDATE categories SET name = :name WHERE name = ?;", category)
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

func GetAllCategories(db *db.DB) (Categories, error) {
	var ret Categories
	err := db.Select(&ret, "SELECT * FROM categories;")
	if err != nil {
		return nil, err
	}
	return ret, nil
}
