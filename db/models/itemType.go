package models

import (
	"io"
)

type ItemTypeList []ItemType
type ItemType struct {
	Name         string `json:"name" db:"name" arg:"-n,--name,required"`
	CategoryName string `json:"category_name" db:"category_name" arg:"-c,--cat,required"`
}

// writes a dot node table to the writer
func (itemTypes ItemTypeList) WriteDot(w io.Writer) {
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

// FillDefaults sets any zero values in the given ItemType to their
// default values based on the provided Defaults.
func (it *ItemType) FillDefaults(defaults *ItemType) {
	if it.Name == "" {
		it.Name = defaults.Name
	}
	if it.CategoryName == "" {
		it.CategoryName = defaults.CategoryName
	}
}
