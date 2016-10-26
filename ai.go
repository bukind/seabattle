package seabattle

import (
	"math/rand"
)

type AI interface {
	FindHit() (int, int)
}

type simpleAI struct {
	peer *Player
}

type trackingAI struct {
  peer *Player
}

// ---

func SimpleAI(peer *Player) AI {
	return &simpleAI{peer}
}

func (a *simpleAI) FindHit() (int, int) {
	const attempts = 50
	p := a.peer
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

func TrackingAI(peer *Player) AI {
	return &trackingAI{peer}
}

func (a *trackingAI) FindHit() (int, int) {
	const attempts = 50
	p := a.peer
	var x, y int
	out.Printf("FindHit(lasthit=%v)\n", p.lasthit)
	if p.lasthit != nil {
	  x = p.lasthit.x
		y = p.lasthit.y
		for d := -1; d < 2; d += 2 {
			x0 := p.peer.GetCellMisteryX(x, y, d)
			out.Printf("GetCellMisteryX(%s,%d) = %d\n", PosToStr(x, y), d, x0)
			if x0 != -1 {
				return x0, y
			}
		}
		for d := -1; d < 2; d += 2 {
			y0 := p.peer.GetCellMisteryY(x, y, d)
			out.Printf("GetCellMisteryY(%s,%d) = %d\n", PosToStr(x, y), d, y0)
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
