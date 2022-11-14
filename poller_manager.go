package knet

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

func setNumLoops(numLoops int) error {
	return PollerManager.SetNumLoops(numLoops)
}

var PollerManager *pollerManager

func init() {
	var loops = runtime.GOMAXPROCS(0)/20 + 1
	PollerManager = &pollerManager{}
	PollerManager.SetNumLoops(loops)
}

type pollerManager struct {
	NumLoops int
	polls    []Poll // all the polls
}

// SetNumLoops setup num for pollers
func (m *pollerManager) SetNumLoops(numLoops int) error {
	if numLoops < 1 {
		return fmt.Errorf("set invalid numLoops[%d]", numLoops)
	}

	if numLoops < m.NumLoops {
		var polls = make([]Poll, numLoops)
		for idx := 0; idx < m.NumLoops; idx++ {
			if idx < numLoops {
				polls[idx] = m.polls[idx]
			} else {
				if err := m.polls[idx].Close(); err != nil {
					log.Printf("poller Close failed: %v\n", err)
				}
			}
		}
		m.NumLoops = numLoops
		m.polls = polls
		return nil
	}

	m.NumLoops = numLoops
	return m.Run()
}

// Close release all resources.
func (m *pollerManager) Close() error {
	for _, poll := range m.polls {
		poll.Close()
	}
	m.NumLoops = 0
	m.polls = nil
	return nil
}

// Run all pollers.
func (m *pollerManager) Run() error {
	for idx := len(m.polls); idx < m.NumLoops; idx++ {
		var poll = NewDefaultPoller()
		m.polls = append(m.polls, poll)
		go poll.Wait()
	}

	return nil
}

// Pick rand get a poller
func (m *pollerManager) Pick() Poll {
	return m.polls[rand.Intn(m.NumLoops)]
}
