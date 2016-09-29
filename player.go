package seabattle

import (
	"bytes"
	"fmt"
	"math/rand"
)

type Player struct {
	self   *Board
	peer   *Board
	ingame bool
}

func NewPlayer(size int) *Player {
	p := new(Player)
	p.self = NewBoard(size)
	p.peer = NewBoard(size)
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
	if res == ResultGameOver {
		p.ingame = false
	}
}

func (p *Player) FindHit() (int, int) {
	const attempts = 50
	var x, y int
	for i := 0; i < attempts; i++ {
		x = rand.Intn(len(p.peer.Cells[0]))
		y = rand.Intn(len(p.peer.Cells))
		if p.peer.Cells[y][x] == CellEmpty {
			break
		}
	}
	return x, y
}
