package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Merith-TK/gosh/api"
)

type clearCmd string

func (t clearCmd) Name() string      { return string(t) }
func (t clearCmd) Usage() string     { return `clear` }
func (t clearCmd) ShortDesc() string { return `clear terminal buffer` }
func (t clearCmd) LongDesc() string  { return t.ShortDesc() }
func (t clearCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	cmd := exec.Command("reset")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return ctx, nil
}

type testCmds struct{}

func (t *testCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)
	fmt.Fprintln(out, "test module loaded OK")
	return nil
}

func (t *testCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"clear": clearCmd("clear"),
	}
}

var Commands testCmds
