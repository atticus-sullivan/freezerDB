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
)

type ItemTypeList []ItemType
type ItemType struct {
	Name         string `json:"name" db:"name" arg:"-n,--name,required"`
	CategoryName string `json:"category_name" db:"category_name" arg:"-c,--cat,required"`
}

var itemTypeHdr [2]string = [2]string{"Name", "Kategorie"}

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
