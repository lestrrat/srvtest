// package srvtest provides utilities to run tests against servers

package srvtest

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	DefaultMinPortNumber = 50000
	DefaultMaxPortNumber = 60000
)

func tcpAddr(p int) string {
	return ":" + strconv.Itoa(p)
}

func tcp6Addr(p int) string {
	return "[::]:" + strconv.Itoa(p)
}

func DialPort(ctx context.Context, port int) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "tcp", tcpAddr(port))
}

func DialPortV4(ctx context.Context, port int) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "tcp4", tcpAddr(port))
}

func DialPortV6(ctx context.Context, port int) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "tcp6", tcp6Addr(port))
}

// ListenPort is a utility function to an integer port number
func ListenPort(p int) (net.Listener, error) {
	return net.Listen("tcp", tcpAddr(p))
}

// ListenPort is a utility function to an integer port number from
// a IPv4 network
func ListenPortV4(p int) (net.Listener, error) {
	return net.Listen("tcp4", tcpAddr(p))
}

// ListenPort is a utility function to an integer port number from
// a IPv6 network
func ListenPortV6(p int) (net.Listener, error) {
	return net.Listen("tcp6", tcp6Addr(p))
}

// Returns "empty", i.e. available ports that you can listen to.
// EmptyPort returns only one port number. This suits most casual
// needs, but if you *reall* must grab a port, you might want to
// keep trying until you really get one. In that case use the
// EmptyPort"s" tool
func EmptyPort(octx context.Context, options ...Option) int {
	ctx, cancel := context.WithCancel(octx)
	defer cancel()

	ch := EmptyPorts(ctx, options...)
	return <-ch
}

// Returns "empty", i.e. available ports that you can listen to.
// it returns a channel that (eventually) produces all possible
// ports within the given boundary, to make sure you have maxium
// probability of grabbing an actual empty port
//
// Note that this operation is NOT atomic, and therefore it is
// very possible to that the port has been taken after `EmportPorts`
// but before the user could bind to it.
//
// Port numbers in the range of (min + random number up to 10% o max - min) to max will be used.
func EmptyPorts(ctx context.Context, options ...Option) <-chan int {
	var lo = DefaultMinPortNumber
	var hi = DefaultMaxPortNumber
	for _, option := range options {
		if option == nil {
			continue
		}

		switch option.Name() {
		case optkeyMaxPortNumber:
			hi = option.Value().(int)
		case optkeyMinPortNumber:
			lo = option.Value().(int)
		}
	}

	v := hi - lo
	if v < 0 {
		v = 0
	}

	ch := make(chan int, 1)
	go probeEmptyPorts(ctx, ch, lo, hi, v)
	return ch
}

func probeEmptyPorts(ctx context.Context, ch chan int, lo, hi, v int) {
	defer close(ch)

	for p := lo + rand.Intn(v); p < hi; p++ {
		l, e := ListenPort(p)
		if e == nil {
			// yay!
			l.Close()

			select {
			case <-ctx.Done():
				return
			case ch <- p:
			}
			continue
		}
	}

	select {
	case <-ctx.Done():
		return
	default:
	}
}

// Waits for the specified port to be available.
func WaitPort(ctx context.Context, port int) error {
	for {
		l, e := ListenPort(port)
		if e == nil {
			l.Close()
			return nil
		}

		tm := time.NewTimer(time.Duration(rand.Intn(1000)) * time.Millisecond)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tm.C:
		}
	}
}
