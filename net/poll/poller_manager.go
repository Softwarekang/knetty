package poll

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

var PollerManager *pollerManager

func init() {
	var loops = runtime.GOMAXPROCS(0)/20 + 1
	PollerManager = &pollerManager{}
	_ = PollerManager.SetPollerNums(loops)
}

type pollerManager struct {
	NumLoops int
	pollers  []Poll // all the pollers
}

// SetPollerNums setup num for pollers
func (m *pollerManager) SetPollerNums(n int) error {
	if n < 1 {
		return fmt.Errorf("SetPollerNums(n int):@n < 0")
	}

	if n < m.NumLoops {
		var polls = make([]Poll, n)
		for idx := 0; idx < m.NumLoops; idx++ {
			if idx < n {
				polls[idx] = m.pollers[idx]
			} else {
				if err := m.pollers[idx].Close(); err != nil {
					log.Printf("close poller err: %v\n", err)
				}
			}
		}
		m.NumLoops = n
		m.pollers = polls
		return nil
	}

	m.NumLoops = n
	return m.Run()
}

// Close release all resources.
func (m *pollerManager) Close() error {
	for _, poller := range m.pollers {
		if err := poller.Close(); err != nil {
			log.Printf("close poller err:%v \n", err)
		}
	}
	m.NumLoops, m.pollers = 0, nil
	return nil
}

// Run all pollers
func (m *pollerManager) Run() error {
	for idx := len(m.pollers); idx < m.NumLoops; idx++ {
		var poller = NewDefaultPoller()
		m.pollers = append(m.pollers, poller)

		go func() {
			if err := poller.Wait(); err != nil {
				log.Println(err)
			}
		}()
	}

	return nil
}

// Pick rand get a poller
func (m *pollerManager) Pick() Poll {
	return m.pollers[rand.Intn(m.NumLoops)]
}
