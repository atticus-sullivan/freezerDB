package models

import (
	"bufio"
	"fmt"
	"html"
	"io"
)

type ItemTypeList []ItemType
type ItemType struct {
	Name         string `json:"name" db:"name" arg:"-n,--name,required"`
	CategoryName string `json:"category_name" db:"category_name" arg:"-c,--cat,required"`
}

var itemTypeHdr [2]string = [2]string{"Name", "Category"}

// writes a dot node table to the writer
func (items ItemTypeList) WriteDot(wParam io.Writer) error {
	w := bufio.NewWriter(wParam)
	defer w.Flush()

	_, err := w.Write([]byte(`
digraph structs {
	node [shape=plaintext] struct [label=<
		<table cellspacing="2" border="0" rows="*" columns="*">
`))
	if err != nil {
		return err
	}

	// write header
	if err := writeDotHdr(w, itemTypeHdr[:]); err != nil {
		return err
	}

	// write rows
	for _, item := range items {
		_, err = fmt.Fprintf(w, "<tr><td>  %s  </td><td>  %s  </td></tr>",
			html.EscapeString(item.Name),
			html.EscapeString(item.CategoryName),
		)
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("</table>>];}\n"))
	if err != nil {
		return err
	}
	return nil
}

// FillDefaults sets any zero values in the given ItemType to their
// default values based on the provided Defaults.
func (it *ItemType) FillDefaults(defaults *ItemType) {
	if defaults != nil {
		if it.Name == "" {
			it.Name = defaults.Name
		}
		if it.CategoryName == "" {
			it.CategoryName = defaults.CategoryName
		}
	}
}
