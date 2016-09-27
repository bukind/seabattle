package main

import (
	"bytes"
	"fmt"
	"github.com/bukind/seabattle"
	"github.com/hoisie/web"
	// "log"
	// "os"
)

type Game struct {
	self *seabattle.Board
	peer *seabattle.Board
}

var game Game

func start(ctx *web.Context, page string) string {
	out := &bytes.Buffer{}
	fmt.Fprintf(out, "<!DOCTYPE html>\n<html>\n<head>\n")
	fmt.Fprintf(out, "<meta charset=\"UTF-8\"/>\n")
	fmt.Fprintf(out, "<title>Some title</title>\n</head>\n")
	fmt.Fprintf(out, "<body>\n")
	fmt.Fprintf(out, "<table id=\"selfboard\">\n")
	fmt.Fprintf(out, "%s</table>\n", game.self.HtmlShow())
	fmt.Fprintf(out, "<table id=\"peerboard\">\n")
	fmt.Fprintf(out, "%s</table>\n", game.peer.HtmlShow())
	fmt.Fprintf(out, "</body>\n")
	fmt.Fprintf(out, "</html>\n")
	return out.String()
}

func main() {
	/*
	  f, err := os.Create("seabattle.log")
		if err != nil {
		  fmt.Println("%v", err)
			return
		}
		logger := log.New(f, "", log.Ldate | log.Ltime)
		web.SetLogger(logger)
	*/
	game.self = seabattle.NewBoard(10)
	game.peer = seabattle.NewBoard(10)

	web.Get("/(.*)", start)
	web.Run("0.0.0.0:9999")
}
