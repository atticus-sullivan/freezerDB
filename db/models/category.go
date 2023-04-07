package models

import (
	"io"
)

type CategoryList []Category
type Category struct {
	Name string `json:"name" db:"name" arg:"-n,--name,required"`
}

// writes a dot node table to the writer
func (categories CategoryList) WriteDot(w io.Writer) {
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

// FillDefaults sets any zero values in the given Category to their
// default values based on the provided Defaults.
func (c *Category) FillDefaults(defaults *Category) {
	if c.Name == "" {
		c.Name = defaults.Name
	}
}
