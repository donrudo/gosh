package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/donrudo/gosh/api"
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

type termCmds struct{}

func (t *termCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)
	fmt.Fprintln(out, "term module loaded OK")
	return nil
}

func (t *termCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"clear": clearCmd("clear"),
		"echo":  echoCmd("echo"),
	}
}

var Commands termCmds
