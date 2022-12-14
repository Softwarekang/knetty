// Package knetty .
package knetty

import "github.com/Softwarekang/knetty/net/poll"

// SetPollerNums set reactor goroutine nums
func SetPollerNums(n int) error {
	return poll.PollerManager.SetPollerNums(n)
}
