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
  CellEmpty: "__",
	CellMiss: "..",
	CellShip: "[]",
	CellHit: "++",
	CellDebris: "XX",
	CellShadow: "~~",
}

func HtmlShowCell(c Cell) string {
  return htmlCellRep[c]
}

func (b *Board) HtmlShow() string {
  out := &bytes.Buffer{}
  for i, row := range b.Cells {
	  out.WriteString("<tr>\n")
	  for j, cell := range row {
		  fmt.Fprintf(out, "<td id=\"%c%c\">%s</td>", 'A'+j, '1'+i, HtmlShowCell(cell))
		}
		out.WriteString("\n</tr>\n")
	}
	return out.String()
}
