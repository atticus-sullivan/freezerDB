package models

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
	"bufio"
	"fmt"
	"html"
	"io"
	"strconv"
	"time"
)

type FreezerItemList []FreezerItem
type FreezerItem struct {
	ID         uint      `json:"id" db:"id"`
	Date       time.Time `json:"date" db:"date" arg:"-d,--date"` // default to today
	Identifier string    `json:"identifier" db:"identifier" arg:"--identifier"`
	Amount     string    `json:"amount" db:"amount" arg:"-a,--amount,required"`
	Misc       string    `json:"misc" db:"misc" arg:"-m,--misc"`
	ItemName   string    `json:"item_name" db:"item_name" arg:"-n,--name,required"`
}

const (
	dateFormat = "2006-01-02"
)

var frItemHdr [6]string = [6]string{"ID", "Name", "Datum", "Merkmal", "Anzahl", "Bemerkung"}

// writes a dot node table to the writer
func (items FreezerItemList) WriteDot(wParam io.Writer) error {
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
	if err := writeDotHdr(w, frItemHdr[:]); err != nil {
		return err
	}

	// write rows
	for _, item := range items {
		_, err = fmt.Fprintf(w, "<tr><td>  %s  </td><td>  %s  </td><td>  %s  </td><td>  %s  </td><td>  %s  </td><td>  %s  </td></tr>",
			html.EscapeString(strconv.FormatUint(uint64(item.ID), 10)),
			html.EscapeString(item.ItemName),
			html.EscapeString(item.Date.Format(dateFormat)),
			html.EscapeString(item.Identifier),
			html.EscapeString(item.Amount),
			html.EscapeString(item.Misc),
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

// FillDefaults sets any zero values in the given FreezerItem to their
// default values based on the provided Defaults.
func (fi *FreezerItem) FillDefaults(defaults *FreezerItem) {
	if defaults != nil {
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
}
