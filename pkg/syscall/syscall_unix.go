//go:build linux && !race

package syscall

import (
	"golang.org/x/sys/unix"
)

func Writev(fd int, bs [][]byte) (int, error) {
	if len(bs) == 0 {
		return 0, nil
	}

	return unix.Writev(fd, bs)
}

func Readv(fd int, bs [][]byte) (int, error) {
	if len(bs) == 0 {
		return 0, nil
	}

	return unix.Readv(fd, bs)
}
