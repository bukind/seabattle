package seabattle

import (
	"bytes"
	"fmt"
)

type Pos struct {
	x int
	y int
}

type Player struct {
	name    string
	self    *Board
	peer    *Board
	ingame  bool
	lasthit *Pos
}

func NewPlayer(name string, size int) *Player {
	p := new(Player)
	p.name = name
	p.self = NewBoard(size, false)
	p.peer = NewBoard(size, true)
	p.ingame = true
	return p
}

func (p *Player) AddRandomShips() bool {
	out.Printf("%s ships placed\n", p.name)
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
	out.Printf("%s.Hit(%d,%d/%s) => %v", p.name, x, y, PosToStr(x,y), res)
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
