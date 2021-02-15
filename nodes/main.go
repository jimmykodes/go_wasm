package main

import (
	"syscall/js"

	"github.com/jimmykodes/wasm/nodes/internal"
)

func main() {
	c := make(chan struct{}, 0)
	functions := map[string]js.Func{
		"newBoard": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return internal.NewBoard(args[0]).Serializer()
		}),
	}
	for name, f := range functions {
		js.Global().Set(name, f)
	}
	<-c
}
