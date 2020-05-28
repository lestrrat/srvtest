package srvtest

import "github.com/lestrrat-go/srvtest/internal/option"

type Option = option.Interface

const (
	optkeyMinPortNumber       = "min-port-number"
	optkeyMaxPortNumber       = "max-port-number"
)

// WithMinPortNumber indicates the absolute minimum value used for
// the port numbers used in `EmptyPorts`. Note that this is not
// necessarily the lowest port number that the tool starts looking
// for availability, as the actual starting value is derived by
// the minimum value specified by this option plus a random number
// up to 10% of the difference between the maximum and minimum.
// So for example, if min = 50000 and max = 60000, the starting
// value is a random number between 50000 and 51000:
func WithMinPortNumber(p int) Option {
	return option.New(optkeyMinPortNumber, p)
}

func WithMaxPortNumber(p int) Option {
	return option.New(optkeyMaxPortNumber, p)

}
