package seabattle

import (
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

var htmlCellRep = map[Cell]string{
	CellEmpty:  "__",
	CellMiss:   "..",
	CellShip:   "\\/",
	CellHit:    "++",
	CellDebris: "xx",
	CellShadow: "~~",
}

func CellToStr(x, y int) string {
	return fmt.Sprintf("%c%d", 'A'+x, y+1)
}
