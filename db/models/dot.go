package models

import (
	"bufio"
	"fmt"
)

// writes the header row for a dot node table
func writeDotHdr(w *bufio.Writer, headers []string) error {
	_, err := w.Write([]byte("<tr>"))
	if err != nil {
		return err
	}
	for _, header := range headers {
		_, err = fmt.Fprintf(w, "<td><b>%s</b></td>", header)
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte("</tr>"))
	if err != nil {
		return err
	}
	return nil
}
