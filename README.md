# srvtest

## DESCRIPTION

A long, long, time ago, in a community far, far away, there used to be tools to
choose random local ports, and spin up local servers for testing (i.e. 
[Test::TCP](https://metacpan.org/pod/Test::TCP)

We [ported](https://github.com/lestrrat-go/tcputil) these to Go](https://github.com/lestrrat-go/tcptest) as an exercise,
but this was a) before `context.Context`, and b) before we knew how to write
idiomatic Go code.

`srvtest` is an attempt at re-doing all of this.

## SYNOPSIS

```
func ExampleEmptyPorts() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var cmd *exec.Cmd
	var port int
	for port = range srvtest.EmptyPorts(ctx) {
		// port *should* be empty, at least when EmptyPorts() was
		// probing it. If you failed to listen to it, then use the
		// next one that is returned by EmptyPorts()
		cmd = exec.CommandContext(ctx, "memcached", "-p", fmt.Sprintf("%d", port))
		go cmd.Run() // error checking omitted
		break
	}

	if cmd == nil {
		fmt.Println("cmd is nil (FAIL)")
		return
	}

	srvtest.WaitPort(ctx, port)

	cmd.Wait()
	fmt.Println("done")
	// OUTPUT:
	// done
}
```
