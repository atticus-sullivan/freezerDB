package models

import (
	"io"
	"strconv"
	"time"
)

type FreezerItemList []FreezerItem
type FreezerItem struct {
	ID         uint      `json:"id" db:"id"`
	Date       time.Time `json:"date" db:"date" arg:"-d,--date"` // default to today
	Identifier string    `json:"identifier,omitempty" db:"identifier" arg:"--identifier"`
	Amount     string    `json:"amount" db:"amount" arg:"-a,--amount,required"`
	Misc       string    `json:"misc,omitempty" db:"misc" arg:"-m,--misc"`
	ItemName   string    `json:"item_name" db:"item_name" arg:"-n,--name,required"`
}

// writes a dot node table to the writer
func (freezerItems FreezerItemList) WriteDot(w io.Writer) {
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

// FillDefaults sets any zero values in the given FreezerItem to their
// default values based on the provided Defaults.
func (fi *FreezerItem) FillDefaults(defaults *FreezerItem) {
    if fi.Date.IsZero() {
        fi.Date = defaults.Date
    }
    if fi.Identifier == "" {
        fi.Identifier = defaults.Identifier
    }
    if fi.Amount == "" {
        fi.Amount = defaults.Amount
    }
    if fi.Misc == "" {
        fi.Misc = defaults.Misc
    }
    if fi.ItemName == "" {
        fi.ItemName = defaults.ItemName
    }
}
