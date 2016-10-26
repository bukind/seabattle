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
	CellMistery // the one is not hit yet
)

type Cell int

type HtmlCellRep struct {
	class string
	text  string
}

var htmlCellRep = map[Cell]HtmlCellRep{
	CellEmpty:   HtmlCellRep{"empty", "  "},
	CellMiss:    HtmlCellRep{"miss", ".."},
	CellShip:    HtmlCellRep{"ship", "\\/"},
	CellHit:     HtmlCellRep{"hit", "++"},
	CellDebris:  HtmlCellRep{"debris", "xx"},
	CellShadow:  HtmlCellRep{"shadow", "~~"},
	CellMistery: HtmlCellRep{"mist", "__"},
}

func PosToStr(x, y int) string {
	return fmt.Sprintf("%c%d", 'A'+x, y+1)
}

func (c Cell) String() string {
	return htmlCellRep[c].class
}

func (c Cell) htmlClass() string {
	return htmlCellRep[c].class
}

func (c Cell) htmlShow(x, y int, active bool) string {
	if active && c == CellMistery {
		return fmt.Sprintf("<a href=\"/hit?x=%d&y=%d\">%s</a>",
			x, y, htmlCellRep[c].text)
	}
	return htmlCellRep[c].text
}
