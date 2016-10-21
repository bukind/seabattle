package seabattle

import (
	"bytes"
	"fmt"
	"math/rand"
)

type Pos struct {
	x int
	y int
}

type Player struct {
	self    *Board
	peer    *Board
	ingame  bool
	lasthit *Pos
}

func NewPlayer(size int) *Player {
	p := new(Player)
	p.self = NewBoard(size, false)
	p.peer = NewBoard(size, true)
	p.ingame = true
	return p
}

func (p *Player) AddRandomShips() bool {
	return p.self.AddRandomShips()
}

func (p *Player) HtmlShow() string {
	out := &bytes.Buffer{}
	fmt.Fprintf(out, "<table id=\"selfboard\">\n%s</table>\n",
		p.self.HtmlShow(false))
	fmt.Fprintf(out, "<table id=\"peerboard\">\n%s</table>\n",
		p.peer.HtmlShow(p.ingame))
	return out.String()
}

func (p *Player) Hit(x, y int) Result {
	res := p.self.Hit(x, y)
	if res == ResultGameOver {
		p.ingame = false
	}
	return res
}

func (p *Player) ApplyResult(x, y int, res Result) {
	p.peer.ApplyResult(x, y, res)
	switch res {
	case ResultHit:
		p.lasthit = &Pos{x, y}
	case ResultKill:
		p.lasthit = nil
	case ResultGameOver:
		p.lasthit = nil
		p.ingame = false
	}
}

// Find a place where to hit.
func (p *Player) FindHit() (int, int) {
	const attempts = 50
	var x, y int
	if p.lasthit != nil {
	  x = p.lasthit.x
		y = p.lasthit.y
		for d := -1; d < 2; d += 2 {
			x0 := p.peer.GetCellMisteryX(x, y, d)
			if x0 != -1 {
				return x0, y
			}
		}
		for d := -1; d < 2; d += 2 {
			y0 := p.peer.GetCellMisteryY(x, y, d)
			if y0 != -1 {
				return x, y0
			}
		}
	}
	for i := 0; i < attempts; i++ {
		x = rand.Intn(len(p.peer.Cells[0]))
		y = rand.Intn(len(p.peer.Cells))
		if p.peer.Cells[y][x] == CellMistery {
			break
		}
	}
	return x, y
}
