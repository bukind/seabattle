package seabattle

import (
	"bytes"
	"fmt"
	"math/rand"
)

type Row []Cell

type Board struct {
	Cells []Row
}

func NewBoard(size int, isPeer bool) *Board {
	b := new(Board)
	b.init(size, isPeer)
	return b
}

func (b *Board) init(size int, isPeer bool) {
	b.Cells = make([]Row, size)
	cellType := CellEmpty
	if isPeer {
	  cellType = CellMistery
	}
	cellProto := Cell(cellType)
	for i := 0; i < size; i++ {
		row := make(Row, size)
		b.Cells[i] = row
		for j := 0; j < size; j++ {
			row[j] = cellProto
		}
	}
}

func (b *Board) AddRandomShips() bool {
	const maxshipsize = 5
	const maxattempt = 50
	num := 1
	for s := maxshipsize; s > 0; s-- {
		for n := 0; n < num; n++ {
			for attempt := 0; attempt < maxattempt; attempt++ {
				if b.placeShip(s) {
					break
				}
				if attempt+1 >= maxattempt {
					return false
				}
			}
		}
		num++
	}
	return true
}

func (b *Board) placeShip(s int) bool {
	dx := 1
	dy := s
	if r := rand.Intn(2); r != 0 {
		dx = s
		dy = 1
	}
	x0 := rand.Intn(len(b.Cells[0]) - dx + 1)
	y0 := rand.Intn(len(b.Cells) - dy + 1)
	if dx > 1 {
		if !(b.isCellsEmptyY(y0-1, x0, x0+dx) &&
			b.isCellsEmptyY(y0+1, x0, x0+dx) &&
			b.isCellsEmptyY(y0, x0-1, x0+dx+1)) {
			return false
		}
		for i := x0; i < x0+dx; i++ {
			b.Cells[y0][i] = CellShip
		}
	} else {
		if !(b.isCellsEmptyX(x0-1, y0, y0+dy) &&
			b.isCellsEmptyX(x0+1, y0, y0+dy) &&
			b.isCellsEmptyX(x0, y0-1, y0+dy+1)) {
			return false
		}
		for i := y0; i < y0+dy; i++ {
			b.Cells[i][x0] = CellShip
		}
	}
	return true
}

func (b *Board) isCellsEmptyY(y, x0, x1 int) bool {
	if y < 0 || y >= len(b.Cells) {
		return true
	}
	if x0 < 0 {
		x0 = 0
	}
	if x1 > len(b.Cells[0]) {
		x1 = len(b.Cells[0])
	}
	for i := x0; i < x1; i++ {
		c := b.Cells[y][i]
		if c == CellShip || c == CellHit || c == CellDebris {
			return false
		}
	}
	return true
}

func (b *Board) isCellsEmptyX(x, y0, y1 int) bool {
	if x < 0 || x >= len(b.Cells[0]) {
		return true
	}
	if y0 < 0 {
		y0 = 0
	}
	if y1 > len(b.Cells) {
		y1 = len(b.Cells)
	}
	for i := y0; i < y1; i++ {
		c := b.Cells[i][x]
		if c == CellShip || c == CellHit || c == CellDebris {
			return false
		}
	}
	return true
}

func (b *Board) HtmlShow(active bool) string {
	out := &bytes.Buffer{}
	size := len(b.Cells)
	cid := "s"
	if active {
		cid = "p"
	}
	for y := size; y >= -1; y-- {
		out.WriteString("<tr>\n")
		if y < 0 || y >= size {
			// create a row of letters
			out.WriteString("<th></th>")
			for x := 0; x < size; x++ {
				fmt.Fprintf(out, "<th>%c</th>", 'A'+x)
			}
			out.WriteString("<th></th>")
		} else {
			fmt.Fprintf(out, "<th>%d</th>", y+1)
			for x := 0; x < size; x++ {
				c := b.Cells[y][x]
				fmt.Fprintf(out, "<td id=\"%s%s\" class=\"%s\">%s</td>",
					cid, PosToStr(x,y), c.htmlClass(),
					c.htmlShow(x, y, active))
			}
			fmt.Fprintf(out, "<th>%d</th>", y+1)
		}
		out.WriteString("\n</tr>\n")
	}
	return out.String()
}

// A peer hits the board.
func (b *Board) Hit(x, y int) Result {
	if x < 0 || x >= len(b.Cells[0]) ||
		y < 0 || y >= len(b.Cells) {
		return ResultOut
	}
	c := b.Cells[y][x]
	switch c {
	case CellEmpty:
		b.Cells[y][x] = CellMiss
		return ResultMiss
	case CellMiss, CellHit, CellDebris, CellShadow:
		return ResultHitAgain
	case CellShip:
		break
	default:
		panic("Something is wrong - bad cell")
	}
	b.Cells[y][x] = CellHit
	if sunk := b.isShipSunk(x, y); sunk == nil {
		return ResultHit
	} else {
		b.markShipSunk(x, y, sunk, false)
	}
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell == CellShip {
				return ResultKill
			}
		}
	}
	return ResultGameOver
}

func (b *Board) isShipSunk(x, y int) *Rect {
	x0 := b.isCellShipX(x, y, -1)
	x1 := b.isCellShipX(x, y, 1)
	y0 := b.isCellShipY(x, y, -1)
	y1 := b.isCellShipY(x, y, 1)
	if x0 < 0 || x1 < 0 || y0 < 0 || y1 < 0 {
		return nil
	}
	return &Rect{x0, y0, x1, y1}
}

func (b *Board) markShipSunk(x, y int, r *Rect, makeShadow bool) {
	for i := r.x0; i <= r.x1; i++ {
		b.Cells[y][i] = CellDebris
	}
	for i := r.y0; i <= r.y1; i++ {
		b.Cells[i][x] = CellDebris
	}
	if makeShadow {
		ym := r.y0 - 1
		yp := r.y1 + 1
		if ym < 0 {
			ym = yp
		}
		if yp >= len(b.Cells) {
			yp = ym
		}
		for i := r.x0; i <= r.x1; i++ {
			if b.Cells[ym][i] == CellEmpty || b.Cells[ym][i] == CellMistery {
				b.Cells[ym][i] = CellShadow
			}
			if b.Cells[yp][i] == CellEmpty || b.Cells[yp][i] == CellMistery {
				b.Cells[yp][i] = CellShadow
			}
		}
		xm := r.x0 - 1
		xp := r.x1 + 1
		if xm < 0 {
			xm = xp
		}
		if xp >= len(b.Cells[0]) {
			xp = xm
		}
		for i := r.y0; i <= r.y1; i++ {
			if b.Cells[i][xm] == CellEmpty || b.Cells[i][xm] == CellMistery {
				b.Cells[i][xm] = CellShadow
			}
			if b.Cells[i][xp] == CellEmpty || b.Cells[i][xp] == CellMistery {
				b.Cells[i][xp] = CellShadow
			}
		}
	}
}

func (b *Board) isCellShipX(x, y, inc int) int {
	for {
		i := x + inc
		if i < 0 || i >= len(b.Cells[0]) {
			return x
		}
		c := b.Cells[y][i]
		if c == CellShip {
			return -1
		}
		if c != CellHit {
			return x
		}
		x = i
	}
}

func (b *Board) isCellShipY(x, y, inc int) int {
	for {
		i := y + inc
		if i < 0 || i >= len(b.Cells) {
			return y
		}
		c := b.Cells[i][x]
		if c == CellShip {
			return -1
		}
		if c != CellHit {
			return y
		}
		y = i
	}
}

func (b *Board) ApplyResult(x, y int, res Result) {
	switch res {
	case ResultOut, ResultHitAgain:
		// do nothing
	case ResultMiss:
		b.Cells[y][x] = CellMiss
	case ResultHit:
		b.Cells[y][x] = CellHit
	case ResultKill, ResultGameOver:
		b.Cells[y][x] = CellHit
		sunk := b.isShipSunk(x, y)
		if sunk == nil {
			panic("ship is not sunk")
		}
		b.markShipSunk(x, y, sunk, true)
	default:
		panic("unknown result - cannot apply")
	}
}
