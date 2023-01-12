// Package knetty .
package knetty

import (
	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/log"
)

// SetPollerNums set reactor goroutine nums
func SetPollerNums(n int) error {
	return poll.PollerManager.SetPollerNums(n)
}

// SetLogger set custom log
func SetLogger(logger log.Logger) {
	log.DefaultLogger = logger
}
