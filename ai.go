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
