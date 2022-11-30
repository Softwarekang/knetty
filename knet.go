package knet

import "github.com/Softwarekang/knet/net/poll"

// SetPollerNums set reactor goroutine nums
func SetPollerNums(n int) error {
	return poll.PollerManager.SetPollerNums(n)
}
