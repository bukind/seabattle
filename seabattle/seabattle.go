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

	return showState(ctx, "Started.")
}

func hit(ctx *web.Context) string {
	if self == nil {
		return start(ctx)
	}
	var msg string
	for firstpass := true; firstpass; firstpass = false {
		var err error
		x, err := getInt(ctx, "x")
		if err != nil {
			msg = err.Error()
			break
		}
		y, err := getInt(ctx, "y")
		if err != nil {
			msg = err.Error()
			break
		}

		result := peer.Hit(x, y)
		self.ApplyResult(x, y, result)

		switch result {
		case seabattle.ResultOut:
			msg = fmt.Sprintf("Invalid position (%d,%d), please strike again", x, y)
		case seabattle.ResultHitAgain:
			msg = fmt.Sprintf("You've strike this cell already, please strike again")
		case seabattle.ResultMiss:
			msg = fmt.Sprintf("You've missed :(")
		case seabattle.ResultHit:
			msg = fmt.Sprintf("You've just hit a target! Go on.")
		case seabattle.ResultKill:
			msg = fmt.Sprintf("The target sunk, find another...")
		case seabattle.ResultGameOver:
			msg = fmt.Sprintf("YOU WON THE BATTLE.")
		default:
			panic("Unknown result of strike")
		}
		break
	}
	return showState(ctx, msg)
}

func showState(ctx *web.Context, msg string) string {
	out := &bytes.Buffer{}
	fmt.Fprintf(out, "<!DOCTYPE html>\n<html>\n<head>\n")
	fmt.Fprintf(out, "<meta charset=\"UTF-8\"/>\n")
	fmt.Fprintf(out, "<title>Some title</title>\n</head>\n")
	fmt.Fprintf(out, "<body>\n")
	fmt.Fprintf(out, "<a href=\"/\">[ restart ]</a>\n")
	fmt.Fprintf(out, "<h2>%s</h2>\n", msg)
	out.WriteString(self.HtmlShow())
	fmt.Fprintf(out, "</body>\n")
	fmt.Fprintf(out, "</html>\n")
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
