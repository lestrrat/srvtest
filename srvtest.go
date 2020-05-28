// package srvtest provides utilities to run tests against servers

package srvtest

import (
	"context"
	"fmt"
	"math/rand"
	"net"
)

func EmptyPorts(ctx context.Context) <-chan int {
	ch := make(chan int)
	go probeEmptyPorts(ctx, ch)
	return ch
}

func probeEmptyPorts(ctx context.Context, ch chan int) {
	defer close(ch)

	for p := 50000 + rand.Intn(1000); p < 60000; p++ {
		l, e := net.Listen("tcp", fmt.Sprintf(":%d", p))
		if e == nil {
			// yay!
			l.Close()

			select {
			case <-ctx.Done():
				return
			case ch <- p:
			}
		}
	}
}
