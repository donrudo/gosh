package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Merith-TK/gosh/api"
)

type pwdCmd string

func (t pwdCmd) Name() string      { return string(t) }
func (t pwdCmd) Usage() string     { return `pwd` }
func (t pwdCmd) ShortDesc() string { return `finds working directory"` }
func (t pwdCmd) LongDesc() string  { return t.ShortDesc() }
func (t pwdCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	customPWD()
	return ctx, nil
}

type dirCmd string

func (t dirCmd) Name() string      { return string(t) }
func (t dirCmd) Usage() string     { return `dir` }
func (t dirCmd) ShortDesc() string { return `lists files in dir` }
func (t dirCmd) LongDesc() string  { return t.ShortDesc() }
func (t dirCmd) Exec(ctx context.Context, args []string) (context.Context, error) {

	//	cmdArgs := strings.Join(args[1:], " ")
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
	return ctx, nil
}

type cdCmd string

func (t cdCmd) Name() string      { return string(t) }
func (t cdCmd) Usage() string     { return `cd` }
func (t cdCmd) ShortDesc() string { return `change dir` }
func (t cdCmd) LongDesc() string  { return t.ShortDesc() }
func (t cdCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	cmdArgs := strings.Join(args[1:], " ")
	os.Chdir(cmdArgs)
	return ctx, nil
}

// command module
type dirCmds struct{}

func (t *dirCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)
	fmt.Fprintln(out, "dir module loaded OK")
	return nil
}

func (t *dirCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"ls":  dirCmd("ls"),
		"dir": dirCmd("dir"),
		"pwd": pwdCmd("pwd"),
		"cd":  cdCmd("cd"),
	}
}

// Commands just dir
var Commands dirCmds

func customPWD() {
	var mydir string
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
}
