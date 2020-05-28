package srvtest_test

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/lestrrat-go/srvtest"
	"github.com/stretchr/testify/assert"
)

func ExampleEmptyPorts() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var cmd *exec.Cmd
	for port := range srvtest.EmptyPorts(ctx) {
		// port *should* be empty, at least when EmptyPorts() was
		// probing it. If you failed to listen to it, then use the
		// next one that is returned by EmptyPorts()
		cmd = exec.CommandContext(ctx, "memcached", "-p", fmt.Sprintf("%d", port))
		go cmd.Run()
		break
	}

	if cmd == nil {
		fmt.Println("cmd is nil (FAIL)")
		return
	}

	cmd.Wait()
	fmt.Println("done")
	// OUTPUT:
	// done
}

func TestEmptyPorts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var ports []int
	for port := range srvtest.EmptyPorts(ctx, srvtest.WithMinPortNumber(55000), srvtest.WithMaxPortNumber(55010)) {
		ports = append(ports, port)
	}

	// We can't be sure if all the ports will be returned, but let's just be
	// optimistic that we will have at least one, and no more than 10
	if !assert.True(t, len(ports) < 10, "should have less than 10 ports") {
		return
	}
	if !assert.True(t, len(ports) > 0, "should have more than 1 ports") {
		return
	}
}

func TestWaitPort(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	port := srvtest.EmptyPort(ctx)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if !assert.NoError(t, err, `should succeed listening to port`) {
		return
	}

	// listening again to this port should fail
	var oldL = l
	l, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if !assert.Error(t, err, `should fail listening to same port`) {
		return
	}

	// release this port in 50 millisecs
	time.AfterFunc(50*time.Millisecond, func() {
		oldL.Close()
	})

	// should *eventually* be able to listen to this address again
	timeout := time.NewTimer(time.Second)
	tick := time.NewTicker(10 * time.Millisecond)
	for l == nil {
		select {
		case <-tick.C:
			l, err = net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				t.Logf("failed to listen")
			}
		case <-timeout.C:
			t.Errorf(`timeout reached`)
		case <-ctx.Done():
			t.Errorf(ctx.Err().Error())
		}
	}
}
