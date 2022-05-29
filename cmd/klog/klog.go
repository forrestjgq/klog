package main

import (
	"fmt"
	"github.com/forrestjgq/klog"
	"os"
)

func main() {
	files := os.Args[1:]
	var in, out string
	if len(files) == 0 || len(files) > 2 {
		fmt.Print("klog <input> [<output>]")
		os.Exit(1)
	}
	in = files[0]
	if len(files) > 1 {
		out = files[1]
	}
	klog.Process(in, out)
}
