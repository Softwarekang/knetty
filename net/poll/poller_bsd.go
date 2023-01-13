//go:build (darwin || netbsd || freebsd || openbsd || dragonfly) && !race

/*
	Copyright 2022 ankangan

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
)

// Kqueue poller for
type Kqueue struct {
	fd int
}

// NewDefaultPoller .
func NewDefaultPoller() Poll {
	fd, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}

	if _, err := syscall.Kevent(fd, []syscall.Kevent_t{{
		Ident:  0,
		Filter: syscall.EVFILT_USER,
		Flags:  syscall.EV_ADD | syscall.EV_CLEAR,
	}}, nil, nil); err != nil {
		panic(err)
	}

	return &Kqueue{fd: fd}
}

// Register .
func (k Kqueue) Register(netFd *NetFileDesc, eventType EventType) error {
	var filter int16
	var flags uint16
	switch eventType {
	case ReadToRW:
		filter, flags = syscall.EVFILT_WRITE, syscall.EV_ADD|syscall.EV_ENABLE
	case Read:
		filter, flags = syscall.EVFILT_READ, syscall.EV_ADD|syscall.EV_ENABLE
	case RwToRead:
		filter, flags = syscall.EVFILT_WRITE, syscall.EV_DELETE|syscall.EV_ONESHOT
	case DeleteRead:
		filter, flags = syscall.EVFILT_READ, syscall.EV_DELETE|syscall.EV_ONESHOT
	case OnceWrite:
		filter, flags = syscall.EVFILT_WRITE, syscall.EV_ADD|syscall.EV_ENABLE|syscall.EV_ONESHOT
	default:
		return fmt.Errorf("kqueue not support the event type:%d", int(eventType))
	}
	if _, err := syscall.Kevent(k.fd, []syscall.Kevent_t{{
		Ident:  uint64(netFd.FD),
		Filter: filter,
		Flags:  flags,
		Udata:  *(**byte)(unsafe.Pointer(&netFd)),
	}}, nil, nil); err != nil {
		return err
	}

	return nil
}

// Wait .
func (k Kqueue) Wait() error {
	events := make([]syscall.Kevent_t, 1024)
	for {
		n, err := syscall.Kevent(k.fd, nil, events, nil)
		if err != nil {
			// kqueue fd is illegal
			if err == syscall.EBADF {
				return nil
			}
			continue
		}

		for i := 0; i < n; i++ {
			event := events[i]
			netFD := *(**NetFileDesc)(unsafe.Pointer(&event.Udata))
			// check interrupt
			if event.Flags&syscall.EV_EOF != 0 {
				if netFD.OnInterrupt != nil {
					if err := netFD.OnInterrupt(); err != nil {
						log.Errorf("netFD onInterrupt err:%v", err)
					}
				}
				continue
			}

			// check read
			if event.Filter == syscall.EVFILT_READ && event.Flags&syscall.EV_ENABLE != 0 {
				if netFD.OnRead != nil {
					if err := netFD.OnRead(); err != nil {
						log.Errorf("netFD OnRead err:%v", err)
					}
				}
				continue
			}

			// check write
			if event.Filter == syscall.EVFILT_WRITE && event.Flags&syscall.EV_ENABLE != 0 {
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

// Close .
func (k Kqueue) Close() error {
	return syscall.Close(k.fd)
}
