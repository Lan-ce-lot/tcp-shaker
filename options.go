package tcp

import "time"

// Options contains parameters for CheckAddrWithOptions
type Options struct {
	// Timeout for the connection attempt
	Timeout time.Duration
	// Network type: "tcp" (default, try IPv4 before IPv6), "tcp4" (IPv4-only),
	// "tcp6" (IPv6-only), "tcp6-then-tcp4" (try IPv6 before IPv4)
	Network string
	// ZeroLinger indicates whether to set SO_LINGER to 0
	ZeroLinger bool
	// Mark sets the SO_MARK socket option (Linux only)
	Mark int
}
