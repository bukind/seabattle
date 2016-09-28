package seabattle

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
