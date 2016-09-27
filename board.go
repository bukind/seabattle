package seabattle

import (
	"bytes"
	"fmt"
)

const (
	CellEmpty = iota
	CellMiss
	CellShip
	CellHit
	CellDebris
	CellShadow
)

type Cell int

type Row []Cell

type Board struct {
	Cells []Row
}

func NewBoard(size int) *Board {
	b := new(Board)
	b.Cells = make([]Row, size)
	for i := 0; i < size; i++ {
		row := make(Row, size)
		b.Cells[i] = row
		for j := 0; j < size; j++ {
			row[j] = CellEmpty
		}
	}
	return b
}

var htmlCellRep = map[Cell]string{
	CellEmpty:  "__",
	CellMiss:   "..",
	CellShip:   "[]",
	CellHit:    "++",
	CellDebris: "XX",
	CellShadow: "~~",
}

func HtmlShowCell(c Cell) string {
	return htmlCellRep[c]
}

func (b *Board) HtmlShow() string {
	out := &bytes.Buffer{}
	size := len(b.Cells)
	for i := 0; i < size+2; i++ {
		out.WriteString("<tr>\n")
		if i < 1 || i > size {
			// create a row of letters
			out.WriteString("<th></th>")
			for j := 1; j < size+1; j++ {
				fmt.Fprintf(out, "<th>%c</th>", 'A'+j-1)
			}
			out.WriteString("<th></th>")
		} else {
			idx := size - i
			row := b.Cells[idx]
			fmt.Fprintf(out, "<th>%d</th>", idx+1)
			for j := 1; j < size+1; j++ {
				cell := row[j-1]
				fmt.Fprintf(out, "<td id=\"%c%d\">%s</td>",
					'A'+j-1, idx+1, HtmlShowCell(cell))
			}
			fmt.Fprintf(out, "<th>%d</th>", idx+1)
		}
		out.WriteString("\n</tr>\n")
	}
	return out.String()
}
