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
	var fd *os.File
	if len(out) > 0 {
		fd, err = os.Create(out)
		if err != nil {
			glog.Fatalf("create output file %s fail, err: %v", out, err)
		}
		cb = func(s string) {
			_, _ = fd.WriteString(s)
		}
		defer func() {
			_ = fd.Close()
		}()
	}

	m := newManager(cb, 20)
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
