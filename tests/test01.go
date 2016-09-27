package main

import (
  "fmt"
	"github.com/hoisie/web"
)

func hello(ctx *web.Context, val string) string {
	for k, v := range ctx.Params {
	   fmt.Println(k, v)
	}
	return fmt.Sprintf("%s: %v\n", val, ctx.Params)
}

func main() {
	web.Get("/(.*)", hello)
	web.Run("0.0.0.0:9999")
}
