package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"github.com/donrudo/gosh/api"
	"strings"
)

type resolveCmd string

/**
 * Basic network operations: resolve (convert DNS Record to IP addresses)
 */

func (t resolveCmd) Name() string      { return string(t) }                //OP
func (t resolveCmd) Usage() string     { return `resolve [HOST]` }                //OP
func (t resolveCmd) ShortDesc() string { return `resolves a hostname addresses` } //OP
func (t resolveCmd) LongDesc() string  { return t.ShortDesc() }            //OP
func (t resolveCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	out := ctx.Value("gosh.stdout").(io.Writer)

	if len(args) != 2 {

		fmt.Fprintln(out, t.Usage())
		return ctx, nil
	}

	addressList, err := net.LookupHost(args[1])
	if err != nil {
		return ctx, err
	}

	fmt.Fprintln(out, strings.Join(addressList,"\n") )

	return ctx, nil
}

type networkCmds struct{}

func (t *networkCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)

	fmt.Fprintln(out, "Network: resolve module loaded OK")
	return nil
}

func (t *networkCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"resolve": resolveCmd("resolve"),
	}
}

var Commands networkCmds

// If your plugin needs extra functions, declare
// them down here to call upon, or import their
// library.
