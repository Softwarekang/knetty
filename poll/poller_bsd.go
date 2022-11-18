package poll

import (
	"fmt"
	"syscall"
	"unsafe"
)

// KqueuePoller poller for
type KqueuePoller struct {
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

	return &KqueuePoller{fd: fd}
}

// Register .
func (k KqueuePoller) Register(netFd *NetFileDesc, eventType PollEventType) error {
	var filter int16
	var flags uint16
	switch eventType {
	case Write:
		filter, flags = syscall.EVFILT_WRITE, syscall.EV_ADD|syscall.EV_ENABLE
	case Read:
		filter, flags = syscall.EVFILT_READ, syscall.EV_ADD|syscall.EV_ENABLE
	case DeleteWrite:
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
func (k KqueuePoller) Wait() error {
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
						fmt.Printf("netFD onInterrupt err:%v", err)
					}
				}
				continue
			}

			// check read
			if event.Filter == syscall.EVFILT_READ && event.Flags&syscall.EV_ENABLE != 0 {
				if netFD.OnRead != nil {
					if err := netFD.OnRead(); err != nil {
						fmt.Printf("netFD OnRead err:%v", err)
					}
				}
				continue
			}

			// check write
			if event.Filter == syscall.EVFILT_WRITE && event.Flags&syscall.EV_ENABLE != 0 {
				if netFD.OnWrite != nil {
					if err := netFD.OnWrite(); err != nil {
						fmt.Printf("netFD OnWrite err:%v", err)
					}
				}
				continue
			}
		}
	}
}

// Close .
func (k KqueuePoller) Close() error {
	return syscall.Close(k.fd)
}
