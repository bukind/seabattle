package seabattle

const (
	ResultOut      = iota // out of bounds
	ResultHitAgain        // hit the cell that was hit already
	ResultMiss            // miss the ship
	ResultHit             // hit the ship
	ResultKill            // kill the ship
	ResultGameOver        // end of game - no more ships
)

type Result int
