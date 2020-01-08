package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Merith-TK/gosh/api"
)

type editCmd string

func (t editCmd) Name() string      { return string(t) }
func (t editCmd) Usage() string     { return `edit` }
func (t editCmd) ShortDesc() string { return `opens nano` }
func (t editCmd) LongDesc() string  { return t.ShortDesc() }
func (t editCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	cmdArgs := strings.Join(args[1:], " ")
	customEdit(cmdArgs)
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
		"edit": editCmd("edit"),
	}
}

// Commands just custom
var Commands customCmds

func customEdit(cmdArgs string) {

	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Println("ERROR: `$EDITOR` variable not set, defaulting to `vi`")
		time.Sleep(time.Duration(2) * time.Second)
		editor = "vi"
	}
	cmd := exec.Command(editor, cmdArgs)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
