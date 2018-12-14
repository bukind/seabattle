package main

import (
	"bytes"
	"fmt"
	"github.com/bukind/seabattle"
	"github.com/hoisie/web"
	"io"
	"text/template"
	// "log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var self *seabattle.Player
var peer *seabattle.Player
var ai seabattle.AI

func start(ctx *web.Context) string {

	self = seabattle.NewPlayer("self", 10)
	peer = seabattle.NewPlayer("peer", 10)
	ai = seabattle.TrackingAI(peer)

	if !self.AddRandomShips() || !peer.AddRandomShips() {
		return "Cannot place ships, <a href=\"/\">[try again]</a>"
	}

	return showState(ctx, []string{"Started <b>OR STARTED</b>."})
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
			x, y = ai.FindHit()
			result = self.Hit(x, y)
			peer.ApplyResult(x, y, result)
			msg = ""
			switch result {
			case seabattle.ResultOut, seabattle.ResultHitAgain:
				// hit again
				panic("cannot proceed!")
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
					fmt.Sprintf("Peer strikes %s. ", seabattle.PosToStr(x, y))+msg)
			}
		}

	}
	return showState(ctx, msgs)
}

func showState(ctx *web.Context, msgs []string) string {
	var err error
	f, err := os.Open("main.tmpl")
	if err != nil {
		panic("Failed to open template file: " + err.Error())
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(f)
	if err != nil && err != io.EOF {
		panic("Failed to read template file: " + err.Error())
	}

	t := template.New("main")
	if t, err = t.Parse(buf.String()); err != nil {
		panic("Failed to parse template:" + err.Error())
	}
	type mainPage struct {
		Boards string
		Msgs   []string
	}
	m := mainPage{Boards: self.HtmlShow(), Msgs: msgs}
	out := &bytes.Buffer{}
	if err = t.Execute(out, m); err != nil {
		panic("Failed to execute template:" + err.Error())
	}
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
	retcode := 0
	defer func() {
		os.Exit(retcode)
	}()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot find current directory", err)
		retcode = 1
		return
	}
	defer os.Chdir(cwd)

	wantFile := "github.com/bukind/seabattle/seabattle/main.tmpl"
	for fp := wantFile; ; {
		f, err := os.Open(fp)
		defer f.Close()
		if err == nil {
			wd := path.Dir(fp)
			if err := os.Chdir(wd); err != nil {
				retcode = 1
				fmt.Fprintln(os.Stderr, "Cannot chdir to", wd)
				return
			}
			fmt.Println("Changed dir to", wd)
			break
		}
		// continue to search
		words := strings.SplitN(fp, "/", 2)
		if len(words) != 2 {
			// the last attempt was made and we haven't found the file yet.
			retcode = 1
			fmt.Fprintln(os.Stderr, "Cannot find", wantFile)
			return
		}
		fp = words[1]
	}

	rand.Seed(time.Now().Unix())

	web.Get("/", start)
	web.Get("/hit", hit)
	web.SetLogger(seabattle.Logger())

	web.Run("0.0.0.0:9999")
}
