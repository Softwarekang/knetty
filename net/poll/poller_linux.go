//go:build linux && !race

/*
	Copyright 2022 Phoenix

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package poll

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/Softwarekang/knetty/pkg/log"
	msyscall "github.com/Softwarekang/knetty/pkg/syscall"
)

// Epoll poller for epoll.
type Epoll struct {
	fd int
}

// NewDefaultPoller return a  kqueue poller.
func NewDefaultPoller() Poll {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	return &Epoll{
		fd: fd,
	}
}

// Register implements Poll.
func (e *Epoll) Register(netFd *NetFileDesc, eventType EventType) error {
	var op int
	var events uint32
	switch eventType {
	case ReadToRW:
		op, events = syscall.EPOLL_CTL_MOD, syscall.EPOLLIN|syscall.EPOLLOUT
	case Read:
		op, events = syscall.EPOLL_CTL_ADD, syscall.EPOLLIN
	case RwToRead:
		op, events = syscall.EPOLL_CTL_MOD, syscall.EPOLLIN
	case DeleteRead:
		op, events = syscall.EPOLL_CTL_DEL, syscall.EPOLLIN
	case OnceWrite:
		// once write use et trigger
		op, events = syscall.EPOLL_CTL_ADD, uint32(msyscall.EpollET|syscall.EPOLLOUT)
	default:
		return fmt.Errorf("epoll not support the event type:%d", int(eventType))
	}

	return msyscall.EpollCtl(e.fd, op, netFd.FD, &msyscall.EpollEvent{
		Events: events | syscall.EPOLLHUP | syscall.EPOLLRDHUP | syscall.EPOLLERR,
		Udata:  *(*[8]byte)(unsafe.Pointer(&netFd)),
	})
}

// Wait implements Poll.
func (e *Epoll) Wait() error {
	events := make([]msyscall.EpollEvent, 1024)
	for {
		n, err := msyscall.EpollWait(e.fd, events, -1)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		for i := 0; i < n; i++ {
			event := events[i]
			netFD := *(**NetFileDesc)(unsafe.Pointer(&event.Udata))
			// check interrupt
			if event.Events&(syscall.EPOLLHUP|syscall.EPOLLRDHUP|syscall.EPOLLERR) != 0 {
				if netFD.OnInterrupt != nil {
					if err := netFD.OnInterrupt(); err != nil {
						log.Errorf("netFD onInterrupt err:%v", err)
					}
				}
				continue
			}

			// check read
			if event.Events&syscall.EPOLLIN != 0 {
				if netFD.OnRead != nil {
					if err := netFD.OnRead(); err != nil {
						log.Errorf("netFD OnRead err:%v", err)
					}
				}
				continue
			}

			// check write
			if event.Events&syscall.EPOLLOUT != 0 {
				if netFD.OnWrite != nil {
					if err := netFD.OnWrite(); err != nil {
						log.Errorf("netFD OnWrite err:%v", err)
					}
				}
				continue
			}
		}
	}
}

// Close implements Poll.
func (e *Epoll) Close() error {
	return syscall.Close(e.fd)
}
