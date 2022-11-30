package syscall

import "syscall"

// SetConnectionNoBlock set conn read/set no block
func SetConnectionNoBlock(fd int) error {
	return syscall.SetNonblock(fd, true)
}
