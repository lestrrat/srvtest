package srvtest_test

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/lestrrat-go/srvtest"
)

func Example() {
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

	cmd.Wait()
	fmt.Println("done")
	// OUTPUT:
	// done
}
