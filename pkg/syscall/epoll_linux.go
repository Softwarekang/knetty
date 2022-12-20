//go:build !arm64

package syscall

import (
	"syscall"
	"unsafe"
)

var (
	// EpollET et for epoll
	EpollET = -syscall.EPOLLET
)

// EpollEvent epoll event
type EpollEvent struct {
	Events uint32
	Udata  [8]byte
}

// EpollCtl ctl for epoll
func EpollCtl(epfd int, op int, fd int, event *EpollEvent) (err error) {
	_, _, err = syscall.RawSyscall6(syscall.SYS_EPOLL_CTL, uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event)), 0, 0)
	if err == syscall.Errno(0) {
		err = nil
	}
	return err
}

// EpollWait wait for epoll
func EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error) {
	var r0 uintptr
	var _p0 = unsafe.Pointer(&events[0])
	if msec == 0 {
		r0, _, err = syscall.RawSyscall6(syscall.SYS_EPOLL_PWAIT, uintptr(epfd), uintptr(_p0), uintptr(len(events)), 0, 0, 0)
	} else {
		r0, _, err = syscall.Syscall6(syscall.SYS_EPOLL_PWAIT, uintptr(epfd), uintptr(_p0), uintptr(len(events)), uintptr(msec), 0, 0)
	}
	if err == syscall.Errno(0) {
		err = nil
	}
	return int(r0), err
}
