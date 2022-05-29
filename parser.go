package klog

import "strings"
import "github.com/tidwall/gjson"

type request struct {
	line uint64
	str  string
}
type ResultCallback func(line uint64, str string)

const kChanSize = 10000

type LogParser struct {
	rx ResultCallback
	c  chan *request
}

func (p *LogParser) Parse(line uint64, str string) {
	p.c <- &request{
		line: line,
		str:  str,
	}
}

func (p *LogParser) extract(s string) string {
	idx := strings.Index(s, "{")
	if idx < 0 {
		return ""
	}
	s = s[idx:]
	return gjson.Get(s, "log").String()
}
func (p *LogParser) run() {
	for r := range p.c {
		if r == nil {
			return
		}

		s := p.extract(r.str)
		p.rx(r.line, s)
	}
}

func NewParser(rx ResultCallback) *LogParser {
	p := &LogParser{
		rx: rx,
		c:  make(chan *request, kChanSize),
	}

	go p.run()
	return p
}
