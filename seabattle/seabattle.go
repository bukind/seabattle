package main

import (
	"bytes"
	"fmt"
	"github.com/bukind/seabattle"
	"github.com/hoisie/web"
	// "log"
	// "os"
	"math/rand"
	"strconv"
	"time"
)

var self *seabattle.Player
var peer *seabattle.Player

func start(ctx *web.Context) string {

	self = seabattle.NewPlayer(10)
	peer = seabattle.NewPlayer(10)

	if !self.AddRandomShips() || !peer.AddRandomShips() {
		return "Cannot place ships"
	}

	return showState(ctx, []string{"Started."})
}

func hit(ctx *web.Context) string {
	if self == nil {
		return start(ctx)
	}
	var msgs []string
	for firstpass := true; firstpass; firstpass = false {
		var err error
		x, err := getInt(ctx, "x")
		if err != nil {
			msgs = append(msgs, err.Error())
			break
		}
		y, err := getInt(ctx, "y")
		if err != nil {
			msgs = append(msgs, err.Error())
			break
		}

		result := peer.Hit(x, y)
		self.ApplyResult(x, y, result)

		peerTurn := false
		var msg string
		switch result {
		case seabattle.ResultOut:
			msg = fmt.Sprintf("Invalid position (%d,%d), please strike again", x, y)
		case seabattle.ResultHitAgain:
			msg = fmt.Sprintf("You've strike this cell already, please strike again")
		case seabattle.ResultMiss:
			msg = fmt.Sprintf("You've missed :(")
			peerTurn = true
		case seabattle.ResultHit:
			msg = fmt.Sprintf("You've just hit a target! Go on.")
		case seabattle.ResultKill:
			msg = fmt.Sprintf("The target sunk, find another...")
		case seabattle.ResultGameOver:
			msg = fmt.Sprintf("YOU WON THE BATTLE.")
		default:
			panic("Unknown result of strike")
		}

		msgs = append(msgs, msg)

		for peerTurn {
			x, y = peer.FindHit()
			result = self.Hit(x, y)
			peer.ApplyResult(x, y, result)
			msg = ""
			switch result {
			case seabattle.ResultOut, seabattle.ResultHitAgain:
				// hit again
			case seabattle.ResultMiss:
				msg = "It missed."
				peerTurn = false
			case seabattle.ResultHit:
				msg = "Your ship is damaged!"
			case seabattle.ResultKill:
				msg = "Your ship is destroyed!!"
			case seabattle.ResultGameOver:
				msg = "YOU'VE LOST THE BATTLE."
				peerTurn = false
			}
			if len(msg) > 0 {
				msgs = append(msgs,
					fmt.Sprintf("Peer strikes %s. ", seabattle.CellToStr(x, y))+msg)
			}
		}

	}
	return showState(ctx, msgs)
}

func showState(ctx *web.Context, msgs []string) string {
	out := &bytes.Buffer{}
	out.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	out.WriteString("<meta charset=\"UTF-8\"/>\n")
	out.WriteString("<link rel=\"stylesheet\" href=\"style.css\"/>\n")
	out.WriteString("<title>Some title</title>\n</head>\n")
	out.WriteString("<body>\n")
	out.WriteString("<a href=\"/\">[ restart ]</a>\n")
	out.WriteString(self.HtmlShow())
	out.WriteString("<ul id=\"messages\">\n")
	for _, msg := range msgs {
		fmt.Fprintf(out, "<li>%s</li>\n", msg)
	}
	out.WriteString("</ul>\n")
	out.WriteString("</body>\n")
	out.WriteString("</html>\n")
	return out.String()
}

func getInt(ctx *web.Context, name string) (int, error) {
	if _, ok := ctx.Params[name]; !ok {
		return 0, fmt.Errorf("Param '%s' must be supplied", name)
	}
	i, err := strconv.ParseInt(ctx.Params[name], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Cannot convert param '%s' to integer: %s", name, err.Error())
	}
	return int(i), nil
}

func main() {
	rand.Seed(time.Now().Unix())

	web.Get("/", start)
	web.Get("/hit", hit)
	web.Run("0.0.0.0:9999")
}
