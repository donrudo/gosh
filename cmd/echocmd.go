package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Merith-TK/gosh/api"
)

type echoCmd string

func (t echoCmd) Name() string      { return string(t) }
func (t echoCmd) Usage() string     { return `echo` }
func (t echoCmd) ShortDesc() string { return `prints args` }
func (t echoCmd) LongDesc() string  { return t.ShortDesc() }
func (t echoCmd) Exec(ctx context.Context, args []string) (context.Context, error) {

	cmdArgs := strings.Join(args[1:], " ")
	fmt.Println(cmdArgs)
	return ctx, nil
}

// command module
type echoCmds struct{}

func (t *echoCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)
	fmt.Fprintln(out, "echo module loaded OK")
	return nil
}

func (t *echoCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"echo": echoCmd("echo"),
	}
}

// Commands just echo
var Commands echoCmds
