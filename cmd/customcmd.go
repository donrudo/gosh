package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Merith-TK/gosh/api"
)

type pwdCmd string

func (t pwdCmd) Name() string      { return string(t) }
func (t pwdCmd) Usage() string     { return `pwd` }
func (t pwdCmd) ShortDesc() string { return `finds working directory"` }
func (t pwdCmd) LongDesc() string  { return t.ShortDesc() }
func (t pwdCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	//out := ctx.Value("gosh.stdout").(io.Writer)
	//fmt.Fprintln(out, customPWD())
	customPWD()
	return ctx, nil
}

type editCmd string

func (t editCmd) Name() string      { return string(t) }
func (t editCmd) Usage() string     { return `edit` }
func (t editCmd) ShortDesc() string { return `opens nano` }
func (t editCmd) LongDesc() string  { return t.ShortDesc() }
func (t editCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	//out := ctx.Value("gosh.stdout").(io.Writer)
	//fmt.Fprintln(out, customPWD())
	customEdit()
	return ctx, nil
}

// command module
type customCmds struct{}

func (t *customCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)
	fmt.Fprintln(out, "custom module loaded OK")
	return nil
}

func (t *customCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"pwd":  pwdCmd("pwd"),
		"edit": editCmd("edit"),
	}
}

// Commands just custom
var Commands customCmds

func customPWD() {
	var mydir string
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
}

func customEdit() {
	cmd := exec.Command("nano")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
