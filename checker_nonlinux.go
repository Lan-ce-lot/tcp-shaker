//go:build !linux
// +build !linux

package tcp

import (
	"context"
	"net"
	"time"
)

// Checker is a fake implementation.
type Checker struct {
	zeroLinger bool
	isReady    chan struct{}
}

// NewChecker creates a Checker with linger set to zero.
func NewChecker() *Checker {
	return NewCheckerZeroLinger(true)
}

// NewCheckerZeroLinger creates a Checker with zeroLinger set to given value.
func NewCheckerZeroLinger(zeroLinger bool) *Checker {
	isReady := make(chan struct{})
	close(isReady)
	return &Checker{zeroLinger: zeroLinger, isReady: isReady}
}

// CheckingLoop is unnecessary on this platform.
func (c *Checker) CheckingLoop(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// CheckAddr performs a TCP check with given TCP address and timeout.
// NOTE: zeroLinger is ignored on non-POSIX operating systems because
// net.TCPConn.SetLinger is only implemented in src/net/sockopt_posix.go.
func (c *Checker) CheckAddr(addr string, timeout time.Duration) error {
	return c.CheckAddrZeroLinger(addr, timeout, c.zeroLinger)
}

// CheckAddrZeroLinger is CheckerAddr with a zeroLinger parameter.
func (c *Checker) CheckAddrZeroLinger(addr string, timeout time.Duration, zeroLinger bool) error {
	return c.CheckAddrWithOptions(addr, Options{
		Timeout:    timeout,
		Network:    "tcp",
		ZeroLinger: zeroLinger,
		Mark:       0,
	})
}

// CheckAddrWithOptions performs a TCP check with given options
// Supported network types: "tcp" (try IPv4 before IPv6), "tcp4" (IPv4-only),
// "tcp6" (IPv6-only), "tcp6-then-tcp4" (try IPv6 before IPv4)
// NOTE: Mark option is ignored on non-Linux platforms
func (c *Checker) CheckAddrWithOptions(addr string, opts Options) error {
	// Set default network if not specified
	if opts.Network == "" {
		opts.Network = "tcp"
	}

	// Handle advanced network types
	switch opts.Network {
	case "tcp6-then-tcp4":
		// Try IPv6 first, then IPv4 on failure
		if err := c.tryCheckAddr(addr, "tcp6", opts); err == nil {
			return nil
		}
		// If IPv6 fails, try IPv4
		return c.tryCheckAddr(addr, "tcp4", opts)
	default:
		// Standard network types: tcp, tcp4, tcp6
		return c.tryCheckAddr(addr, opts.Network, opts)
	}
}

// tryCheckAddr performs the actual TCP check with a specific network type
func (c *Checker) tryCheckAddr(addr, network string, opts Options) error {
	conn, err := net.DialTimeout(network, addr, opts.Timeout)
	if conn != nil {
		if opts.ZeroLinger {
			// Simply ignore the error since this is a fake implementation.
			_ = conn.(*net.TCPConn).SetLinger(0)
		}
		_ = conn.Close()
	}
	if opErr, ok := err.(*net.OpError); ok {
		if opErr.Timeout() {
			return ErrTimeout
		}
	}
	return err
}

// IsReady is always true on this platform.
func (c *Checker) IsReady() bool { return true }

// WaitReady returns a closed chan on this platform.
func (c *Checker) WaitReady() <-chan struct{} {
	return c.isReady
}

// Close is unnecessary on this platform.
func (c *Checker) Close() error { return nil }
