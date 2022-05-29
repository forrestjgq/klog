package klog

import (
	"bufio"
	"fmt"
	"github.com/forrestjgq/glog"
	"os"
	"strings"
)

func Process(in, out string) {
	f, err := os.Open(in)
	if err != nil {
		glog.Fatalf("Open file %s fail, err: %v", in, err)
	}
	defer func() {
		_ = f.Close()
	}()
	cb := func(s string) {
		fmt.Print(s)
	}

	m := newManager(cb, 8)
	scan := bufio.NewScanner(f)

	for {
		if scan.Scan() {
			t := scan.Text()
			t = strings.TrimSpace(t)
			if len(t) > 0 {
				m.input(t)
			}
		} else {
			break
		}
	}
	m.stop()
}
