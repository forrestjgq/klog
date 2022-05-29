package klog

import (
	"github.com/forrestjgq/glog"
	"sync"
	"sync/atomic"
)

type output func(s string)
type response struct {
	line uint64
	str  string
}
type manager struct {
	end       int32
	out       output
	instances []*LogParser
	nextInst  int    // next instance idx
	seq       uint64 // next sequence
	c         chan *response
	cache     map[uint64]string
	wg        *sync.WaitGroup
}

func (m *manager) run() {
	next := uint64(1) // start from 1
	for r := range m.c {
		if r.line == next {
			m.out(r.str)
			next++
			for {
				if s, exist := m.cache[next]; exist {
					m.out(s)
					next++
				} else {
					break
				}
			}
		} else {
			m.cache[r.line] = r.str
		}
		if m.seq == next && atomic.LoadInt32(&m.end) == 1 {
			m.wg.Done()
			return
		}
	}
}
func (m *manager) nextInstance() *LogParser {
	inst := m.instances[m.nextInst]
	m.nextInst = (m.nextInst + 1) % len(m.instances)
	return inst
}
func (m *manager) input(str string) {
	inst := m.nextInstance()
	inst.Parse(m.seq, str)
	m.seq++
}
func (m *manager) stop() {
	atomic.StoreInt32(&m.end, 1)
	m.wg.Wait()
}

func newManager(out output, nrInstance int) *manager {
	if nrInstance <= 0 {
		glog.Fatalf("instance nr should be > 0")
	}
	m := &manager{
		end:       0,
		out:       out,
		instances: []*LogParser{},
		c:         make(chan *response, 100000),
		seq:       0,
		cache:     make(map[uint64]string),
		wg:        &sync.WaitGroup{},
	}
	m.wg.Add(1)
	for i := 0; i < nrInstance; i++ {
		inst := NewParser(func(line uint64, str string) {
			m.c <- &response{
				line: line,
				str:  str,
			}
		})
		m.instances = append(m.instances, inst)
	}
	go m.run()
	return m
}
