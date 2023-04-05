package models

import (
	"github.com/atticus-sullivan/freezerDB/db"
	"io"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type FreezerItems []FreezerItem
type FreezerItem struct {
	ID         uint      `json:"id" db:"id"`
	Date       time.Time `json:"date" db:"date" arg:"-d,--date"` // default to today
	Identifier string    `json:"identifier" db:"identifier" arg:"--identifier"`
	Amount     string    `json:"amount" db:"amount" arg:"-a,--amount"`
	Misc       string    `json:"misc" db:"misc" arg:"-m,--misc"`
	ItemName   string    `json:"item_name" db:"item_name" arg:"-n,--name,required"`
}

// writes a dot node table to the writer
func (freezerItems FreezerItems) WriteDot(w io.Writer) {
	w.Write([]byte(`
digraph structs {
	node [shape=plaintext] struct [label=<
		<table cellspacing="2" border="0" rows="*" columns="*">
		<tr>`))
	// write header TODO localize?
	w.Write([]byte("<td><B>ID</B></td>\n"))
	w.Write([]byte("<td><B>Name</B></td>\n"))
	w.Write([]byte("<td><B>Identifier</B></td>\n"))
	w.Write([]byte("<td><B>Amount</B></td>\n"))
	w.Write([]byte("<td><B>Date</B></td>\n"))
	w.Write([]byte("<td><B>Misc</B></td>\n"))
	w.Write([]byte("</tr>\n"))

	for _, f := range freezerItems {
		w.Write([]byte("<tr>"))

		// TODO html escape!!!
		// TODO formatter like newlines (<BR/>) every 15 char
		w.Write([]byte("<td>  "))
		w.Write([]byte(strconv.Itoa(int(f.ID))))
		w.Write([]byte("  </td>"))

		w.Write([]byte("<td>  "))
		w.Write([]byte(f.ItemName))
		w.Write([]byte("  </td>"))

		w.Write([]byte("<td>  "))
		w.Write([]byte(f.Identifier))
		w.Write([]byte("  </td>"))

		w.Write([]byte("<td>  "))
		w.Write([]byte(f.Amount))
		w.Write([]byte("  </td>"))

		w.Write([]byte("<td>  "))
		w.Write([]byte(f.Date.Format("2006-01-02")))
		w.Write([]byte("  </td>"))

		w.Write([]byte("<td>  "))
		w.Write([]byte(f.Misc))
		w.Write([]byte("  </td>"))

		w.Write([]byte("</tr>"))
	}
	w.Write([]byte("</table>>];\n"))
}

// InsertFreezerItem inserts a new freezer item into the database
func (item *FreezerItem) Insert(db *db.DB) error {
	if item.Date.IsZero() {
		item.Date = time.Now()
	}
	_, err := db.NamedExec("INSERT INTO freezer_items (date, identifier, amount, misc, item_name) VALUES (:date, :identifier, :amount, :misc, :item_name);", item)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFreezerItem deletes a freezer item from the database
func DeleteFreezerItem(db *db.DB, id uint) error {
	_, err := db.Exec("DELETE FROM freezer_items WHERE id = ?;", id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateFreezerItem updates the given freezer item in the database with non-zero values only
func (item *FreezerItem) Update(db *db.DB, id uint) error {
	if item.Date.IsZero() {
		item.Date = time.Now()
	}
	// Prepare statement and execute with arguments
	query, args, err := sqlx.Named("UPDATE freezer_items SET id = :id, date = :date, identifier = :identifier, amount = :amount, misc = :misc, item_name = :item_name WHERE id = ?;", item)
	if err != nil {
		return err
	}
	query = db.Rebind(query)
	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(append(args, id)...)
	return err
}

func GetAllFreezerItems(db *db.DB) (FreezerItems, error) {
	var ret FreezerItems
	err := db.Select(&ret, "SELECT * FROM freezer_items;")
	if err != nil {
		return nil, err
	}
	return ret, nil
}
