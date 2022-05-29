package main

import (
	"github.com/forrestjgq/klog"
	"os"
)

func main() {
	files := os.Args[1:]
	for _, f := range files {
		klog.Process(f, "")
	}
}
