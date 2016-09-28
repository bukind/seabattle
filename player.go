package seabattle

import (
	"bytes"
	"fmt"
)

type Player struct {
	self *Board
	peer *Board
}

func NewPlayer(size int) *Player {
	p := new(Player)
	p.self = NewBoard(size)
	p.peer = NewBoard(size)
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
		p.peer.HtmlShow(true))
	return out.String()
}

func (p *Player) Hit(x, y int) Result {
	return p.self.Hit(x, y)
}

func (p *Player) ApplyResult(x, y int, res Result) {
	p.peer.ApplyResult(x, y, res)
}
