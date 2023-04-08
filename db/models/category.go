package models

import (
	"bufio"
	"fmt"
	"html"
	"io"
)

type CategoryList []Category
type Category struct {
	Name string `json:"name" db:"name" arg:"-n,--name,required"`
}

var categoryHdr [1]string = [1]string{"Name"}

// writes a dot node table to the writer
func (categories CategoryList) WriteDot(wParam io.Writer) error {
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
	if err := writeDotHdr(w, categoryHdr[:]); err != nil {
		return err
	}

	// write rows
	for _, category := range categories {
		_, err = fmt.Fprintf(w, "<tr><td>%s</td></tr>",
			html.EscapeString(category.Name),
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

// FillDefaults sets any zero values in the given Category to their
// default values based on the provided Defaults.
func (c *Category) FillDefaults(defaults *Category) {
	if defaults != nil {
		if c.Name == "" {
			c.Name = defaults.Name
		}
	}
}
