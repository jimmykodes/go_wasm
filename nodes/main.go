package main

import (
	"syscall/js"

	"github.com/jimmykodes/wasm/nodes/internal"
)

func main() {
	c := make(chan struct{}, 0)
	b := &internal.Board{}
	functions := map[string]js.Func{
		"init":         js.FuncOf(b.Init),
		"getPoints":    js.FuncOf(b.Points),
		"getLines":     js.FuncOf(b.Lines),
		"updatePoints": js.FuncOf(b.Update),
		"setLineCount": js.FuncOf(b.KLines),
		"setThreshold": js.FuncOf(b.Threshold),
	}
	for name, f := range functions {
		js.Global().Set(name, f)
	}
	<-c
}
