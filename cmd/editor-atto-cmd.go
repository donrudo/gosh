package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/Merith-TK/gosh/api"
	"github.com/nsf/termbox-go"
	"gopkg.in/yaml.v2"
)

type attoCmd string

func (t attoCmd) Name() string      { return string(t) }
func (t attoCmd) Usage() string     { return `atto` }
func (t attoCmd) ShortDesc() string { return `opens atto editor` }
func (t attoCmd) LongDesc() string  { return t.ShortDesc() }
func (t attoCmd) Exec(ctx context.Context, args []string) (context.Context, error) {
	cmdArgs := strings.Join(args[1:], " ")
	//attoEdit(cmdArgs)
	attoEditCmd(cmdArgs)
	return ctx, nil
}

// command module
type attoCmds struct{}

func (t *attoCmds) Init(ctx context.Context) error {
	out := ctx.Value("gosh.stdout").(io.Writer)
	fmt.Fprintln(out, "atto module loaded OK")
	return nil
}

func (t *attoCmds) Registry() map[string]api.Command {
	return map[string]api.Command{
		"atto": attoCmd("atto"),
	}
}

// Commands just atto
var Commands attoCmds

// ----------------------------------------------
// Everything below this is an mostly unmodified
// copy of https://github.com/jonpalmisc/atto
// the only modifications are to make the program
// compatable with the gosh shell arguments
// and a few minor modifications to fix the
// program not exiting properly for gosh
// ----------------------------------------------

//##main.go
func attoEditCmd(args string) {
	if len(args) < 1 {
		fmt.Println("Usage: atto <file>")
		//os.Exit(1)
		return
	}

	editor := MakeEditor()
	editor.Open(args)
	editor.Run()
}

//##buffer.go
// BufferLine represents a single line in a buffer.
type BufferLine struct {
	Editor       *Editor
	Text         string
	DisplayText  string
	Highlighting []HighlightType
}

// MakeBufferLine creates a new BufferLine with the given text.
func MakeBufferLine(editor *Editor, text string) (bl BufferLine) {
	bl = BufferLine{
		Editor: editor,
		Text:   text,
	}

	bl.Update()
	return bl
}

// InsertChar inserts a character into the line at the given index.
func (l *BufferLine) InsertChar(i int, c rune) {

	// If a tab is being inserted and the editor is using soft tabs insert a
	// tab's width worth of spaces instead.
	if c == '\t' && l.Editor.Config.SoftTabs {
		l.Text = l.Text[:i] + strings.Repeat(" ", l.Editor.Config.TabSize) + l.Text[i:]
		l.Editor.CursorX += l.Editor.Config.TabSize - 1
	} else {
		l.Text = l.Text[:i] + string(c) + l.Text[i:]
	}

	l.Update()
}

// DeleteChar deletes a character from the line at the given index.
func (l *BufferLine) DeleteChar(i int) {
	if i >= 0 && i < len(l.Text) {
		l.Text = l.Text[:i] + l.Text[i+1:]
		l.Update()
	}
}

// AppendString appends a string to the line.
func (l *BufferLine) AppendString(s string) {
	l.Text += s
	l.Update()
}

// Update refreshes the DisplayText field.
func (l *BufferLine) Update() {
	// Expand tabs to spaces.
	l.DisplayText = strings.ReplaceAll(l.Text, "\t", strings.Repeat(" ", l.Editor.Config.TabSize))

	l.Highlighting = make([]HighlightType, len(l.DisplayText))

	if l.Editor.Config.UseHighlighting {
		switch l.Editor.FileType {
		case FileTypeC, FileTypeCPP:
			HighlightLine(l, &SyntaxC)
		case FileTypeGo:
			HighlightLine(l, &SyntaxGo)
		}
	}
}

// AdjustX corrects the cursor's X position to compensate for rendering effects.
func (l *BufferLine) AdjustX(x int) int {
	tabSize := l.Editor.Config.TabSize
	delta := 0

	for _, c := range l.Text[:x] {
		if c == '\t' {
			delta += (tabSize - 1) - (delta % tabSize)
		}

		delta++
	}

	return delta
}

// IndentLength gets the line's level of indentation in columns.
func (l *BufferLine) IndentLength() (indent int) {
	for j := 0; l.Text[j] == ' ' || l.Text[j] == '\t'; j++ {
		indent++
	}

	return indent
}

//##config.go
// Config holds the editor's configuration and settings.
type Config struct {
	TabSize         int
	SoftTabs        bool
	UseHighlighting bool
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		TabSize:         4,
		SoftTabs:        false,
		UseHighlighting: true,
	}
}

func AttoFolderPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".atto"), nil
}

// ConfigPath returns the path of Atto's config file on the current platform.
func ConfigPath() (string, error) {
	attoFolder, err := AttoFolderPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(attoFolder, "config.yml"), nil
}

// LoadConfig attempts to load the user's config
func LoadConfig() (Config, error) {
	attoFolderPath, err := AttoFolderPath()
	if err != nil {
		panic(err)
	}

	// Create the Atto folder if it doesn't exist already
	if _, err = os.Stat(attoFolderPath); os.IsNotExist(err) {
		err = os.MkdirAll(attoFolderPath, os.ModePerm)
		if err != nil {
			return DefaultConfig(), err
		}
	}

	configPath, err := ConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	// Check if the config file exists. If it does not, create one with the
	// default values.
	if _, err := os.Stat(configPath); err != nil {
		defaultConfig := DefaultConfig()

		// Marshal the default config to YAML.
		yml, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			return DefaultConfig(), err
		}

		// Write the config file.
		err = ioutil.WriteFile(configPath, yml, 0644)
		if err != nil {
			return DefaultConfig(), err
		}

		return DefaultConfig(), nil
	}

	// Read the config file into memory or return the default config if there
	// is an error.
	yml, err := ioutil.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), err
	}

	// Unmarshal the YAML & return the default config if there is an error.
	config := Config{}
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		return DefaultConfig(), err
	}

	return config, nil
}

//##editor.go
// Editor is the editor instance and manages the UI.
type Editor struct {

	// The editor's height and width measured in rows and columns, respectively.
	Width  int
	Height int

	// The cursor's position. The Y value must always be decremented by one when
	// accessing buffer elements since the editor's title bar occupies the first
	// row of the screen. CursorDX is the cursor's X position, with compensation
	// for extra space introduced by rendering tabs.
	CursorX  int
	CursorDX int
	CursorY  int

	// The viewport's
	OffsetX int // The viewport's column offset.
	OffsetY int // The viewport's row offset.

	// The name and type of the file currently being edited.
	FileName string
	FileType FileType

	// The buffer for the current file and whether it has been modifed or not.
	Buffer []BufferLine
	Dirty  bool

	// The current status message and the time it was set.
	StatusMessage     string
	StatusMessageTime time.Time

	Config Config
}

// MakeEditor creates a new Editor instance.
func MakeEditor() Editor {
	editor := Editor{}

	config, err := LoadConfig()
	if err != nil {
		editor.SetStatusMessage("Failed to load config! (%v)", err)
	}

	editor.CursorY = 1
	editor.Config = config

	if err = termbox.Init(); err != nil {
		panic(err)
	}

	return editor
}

// Quit closes the editor and terminates the program.
func (e *Editor) Quit() {
	termbox.Close()
	//os.Exit(0)
	//return
}

// Run starts the editor.
func (e *Editor) Run() {
	e.Draw()

	var infloop string = "true"
	for infloop == "true" {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyArrowUp:
				e.MoveCursor(CursorMoveUp)
			case termbox.KeyArrowDown:
				e.MoveCursor(CursorMoveDown)
			case termbox.KeyArrowLeft:
				e.MoveCursor(CursorMoveLeft)
			case termbox.KeyArrowRight:
				e.MoveCursor(CursorMoveRight)
			case termbox.KeyPgup:
				e.MoveCursor(CursorMovePageUp)
			case termbox.KeyPgdn:
				e.MoveCursor(CursorMovePageDown)
			case termbox.KeyCtrlA:
				e.MoveCursor(CursorMoveLineStart)
			case termbox.KeyCtrlE:
				e.MoveCursor(CursorMoveLineEnd)
			case termbox.KeyCtrlX:
				e.Quit()
				infloop = "false"
			case termbox.KeyCtrlS:
				e.Save()
			case termbox.KeyDelete:
			case termbox.KeyBackspace:
			case termbox.KeyBackspace2:
				e.DeleteChar()
			case termbox.KeyEnter:
				e.BreakLine()
			case termbox.KeyTab:
				e.InsertChar('\t')
			case termbox.KeySpace:
				e.InsertChar(' ')
			default:
				e.InsertChar(event.Ch)
			}
		}
		e.Draw()
	}
}

/* ---------------------------------- I/O ----------------------------------- */

// Open reads a file into a the buffer.
func (e *Editor) Open(path string) {
	e.FileName = path
	e.FileType = GuessFileType(path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		e.InsertLine(0, "")
	}

	f, err := os.Open(path)
	if err != nil {
		e.InsertLine(0, "")
		e.SetStatusMessage("Error: Couldn't open file: %v (%v)", path, err)
		return
	}

	// Read the file line by line and append each line to end of the buffer.
	s := bufio.NewScanner(f)
	for s.Scan() {
		e.InsertLine(len(e.Buffer), s.Text())
	}

	// If the file is completely empty, add an empty line to the buffer.
	if len(e.Buffer) == 0 {
		e.InsertLine(0, "")
	}

	// The file can now be closed since it is loaded into memory.
	f.Close()
}

// Save writes the current buffer back to the file it was read from.
func (e *Editor) Save() {
	var text string

	// Append each line of the buffer (plus a newline) to the string.
	bufferLen := len(e.Buffer)
	for i := 0; i < bufferLen; i++ {
		text += e.Buffer[i].Text + "\n"
	}

	if err := ioutil.WriteFile(e.FileName, []byte(text), 0644); err != nil {
		e.SetStatusMessage("Error: %v.", err)
	} else {
		e.SetStatusMessage("File saved successfully. (%v)", e.FileName)
		e.Dirty = false
	}
}

/* --------------------------------- Buffer --------------------------------- */

// InsertLine inserts a new line to the buffer at the given index.
func (e *Editor) InsertLine(i int, text string) {

	// Ensure the index we are trying to insert at is valid.
	if i >= 0 && i <= len(e.Buffer) {

		// https://github.com/golang/go/wiki/SliceTricks
		e.Buffer = append(e.Buffer, BufferLine{})
		copy(e.Buffer[i+1:], e.Buffer[i:])
		e.Buffer[i] = MakeBufferLine(e, text)
	}
}

// RemoveLine removes the line at the given index from the buffer.
func (e *Editor) RemoveLine(index int) {
	if index >= 0 && index < len(e.Buffer) {
		e.Buffer = append(e.Buffer[:index], e.Buffer[index+1:]...)
		e.Dirty = true
	}
}

// BreakLine inserts a newline character and breaks the line at the cursor.
func (e *Editor) BreakLine() {
	if e.CursorX == 0 {
		e.InsertLine(e.CursorY-1, "")
		e.CursorX = 0
	} else {
		text := e.CurrentRow().Text
		indent := e.CurrentRow().IndentLength()

		e.InsertLine(e.CursorY, text[:indent]+text[e.CursorX:])
		e.CurrentRow().Text = text[:e.CursorX]
		e.CurrentRow().Update()

		e.CursorX = indent
	}

	e.CursorY++
	e.Dirty = true
}

// InsertChar inserts a character at the cursor's position.
func (e *Editor) InsertChar(c rune) {
	if e.CursorY == len(e.Buffer) {
		e.InsertLine(len(e.Buffer), "")
	}

	e.CurrentRow().InsertChar(e.CursorX, c)
	e.CursorX++
	e.Dirty = true
}

// DeleteChar deletes the character to the left of the cursor.
func (e *Editor) DeleteChar() {
	if e.CursorX == 0 && e.CursorY-1 == 0 {
		return
	} else if e.CursorX > 0 {
		e.CurrentRow().DeleteChar(e.CursorX - 1)
		e.CursorX--
	} else {
		e.CursorX = len(e.Buffer[e.CursorY-2].Text)
		e.Buffer[e.CursorY-2].AppendString(e.CurrentRow().Text)
		e.RemoveLine(e.CursorY - 1)
		e.CursorY--
	}

	e.Dirty = true
}

/* ----------------------------- User Interface ----------------------------- */

// DrawTitleBar draws the editor's title bar at the top of the screen.
func (e *Editor) DrawTitleBar() {
	banner := ProgramName + " " + ProgramVersion
	time := time.Now().Local().Format("2006-01-02 15:04")

	name := e.FileName
	if e.Dirty {
		name += " (*)"
	}

	nameLen := len(name)
	timeLen := len(time)

	// Draw the title bar canvas.
	for x := 0; x < e.Width; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorBlack, termbox.ColorWhite)
	}

	// Draw the banner on the left.
	for x := 0; x < len(banner); x++ {
		termbox.SetCell(x, 0, rune(banner[x]),
			termbox.ColorBlack, termbox.ColorWhite)
	}

	// Draw the current file's name in the center.
	namePadding := (e.Width - nameLen) / 2
	for x := 0; x < nameLen; x++ {
		termbox.SetCell(namePadding+x, 0, rune(name[x]),
			termbox.ColorBlack, termbox.ColorWhite)
	}

	// Draw the time on the right.
	for x := 0; x < timeLen; x++ {
		termbox.SetCell(e.Width-timeLen+x, 0, rune(time[x]),
			termbox.ColorBlack, termbox.ColorWhite)
	}
}

// DrawBuffer draws the editor's buffer.
func (e *Editor) DrawBuffer() {
	for y := 0; y < e.Height-2; y++ {
		bufferRow := y + e.OffsetY

		if bufferRow < len(e.Buffer) {
			line := e.Buffer[bufferRow]
			length := len(line.DisplayText) - e.OffsetX

			if length > 0 {
				for x, c := range line.DisplayText[e.OffsetX : e.OffsetX+length] {
					termbox.SetCell(x, y+1, c, line.Highlighting[x].Color(), 0)
				}
			}
		}
	}
}

// DrawStatusBar draws the editor's status bar on the bottom of the screen.
func (e *Editor) DrawStatusBar() {
	right := fmt.Sprintf(" | %v | Line %v, Column %v", e.FileType, e.CursorY, e.CursorDX+1)
	rightLen := len(right)

	// Draw the status bar canvas.
	for x := 0; x < e.Width; x++ {
		termbox.SetCell(x, e.Height-1, ' ',
			termbox.ColorBlack, termbox.ColorWhite)
	}

	// Draw the status message on the left if it hasn't expired.
	if time.Now().Before(e.StatusMessageTime.Add(3 * time.Second)) {
		for x := 0; x < len(e.StatusMessage); x++ {
			termbox.SetCell(x, e.Height-1, rune(e.StatusMessage[x]),
				termbox.ColorBlack, termbox.ColorWhite)
		}
	}

	// Draw the file type and position on the right.
	for x := 0; x < rightLen; x++ {
		termbox.SetCell(e.Width-rightLen+x, e.Height-1, rune(right[x]),
			termbox.ColorBlack, termbox.ColorWhite)
	}
}

// Draw draws the entire editor - UI, buffer, etc. - to the screen & updates the
// cursor's position.
func (e *Editor) Draw() {
	defer termbox.Flush()

	// The screen's height and width should be updated on each render to account
	// for the user resizing the window.
	e.Width, e.Height = termbox.Size()
	e.ScrollView()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	e.DrawTitleBar()
	e.DrawBuffer()
	e.DrawStatusBar()

	termbox.SetCursor(e.CursorDX-e.OffsetX, e.CursorY-e.OffsetY)
}

/* ----------------------------- Input Handling ----------------------------- */

// CursorMove is a type of cursor movement.
type CursorMove int

const (
	// CursorMoveUp moves the cursor up one row.
	CursorMoveUp CursorMove = 0

	// CursorMoveDown moves the cursor down one row.
	CursorMoveDown CursorMove = 1

	// CursorMoveLeft moves the cursor left one column.
	CursorMoveLeft CursorMove = 2

	// CursorMoveRight moves the cursor right one column.
	CursorMoveRight CursorMove = 3

	// CursorMoveLineStart moves the cursor to the first non-whitespace
	// character of the line, or the first character of the line if the cursor
	// is already on the first non-whitespace character.
	CursorMoveLineStart CursorMove = 4

	// CursorMoveLineEnd moves the cursor to the end of the line.
	CursorMoveLineEnd CursorMove = 5

	// CursorMovePageUp moves the cursor up by the  height of the screen.
	CursorMovePageUp CursorMove = 6

	// CursorMovePageDown moves the cursor down by the  height of the screen.
	CursorMovePageDown CursorMove = 7
)

// ScrollView recalculates the offsets for the view window.
func (e *Editor) ScrollView() {
	e.CursorDX = e.CurrentRow().AdjustX(e.CursorX)

	if e.CursorY-1 < e.OffsetY {
		e.OffsetY = e.CursorY - 1
	}

	if e.CursorY+2 >= e.OffsetY+e.Height {
		e.OffsetY = e.CursorY - e.Height + 2
	}

	if e.CursorDX < e.OffsetX {
		e.OffsetX = e.CursorDX
	}

	if e.CursorDX >= e.OffsetX+e.Width {
		e.OffsetX = e.CursorDX - e.Width + 1
	}
}

// MoveCursor moves the cursor according to the operation provided.
func (e *Editor) MoveCursor(move CursorMove) {
	switch move {
	case CursorMoveUp:
		if e.CursorY > 1 {
			e.CursorY--
		}
	case CursorMoveDown:
		if e.CursorY < len(e.Buffer) {
			e.CursorY++
		}
	case CursorMoveLeft:
		if e.CursorX != 0 {
			e.CursorX--
		} else if e.CursorY > 1 {
			e.CursorY--
			e.CursorX = len(e.CurrentRow().Text)
		}
	case CursorMoveRight:
		if e.CursorX < len(e.CurrentRow().Text) {
			e.CursorX++
		} else if e.CursorX == len(e.CurrentRow().Text) && e.CursorY != len(e.Buffer) {
			e.CursorX = 0
			e.CursorY++
		}
	case CursorMoveLineStart:

		// Move the cursor to the end of the indent if the cursor is not there
		// already, otherwise, move it to the start of the line.
		if e.CursorX != e.CurrentRow().IndentLength() {
			e.CursorX = e.CurrentRow().IndentLength()
		} else {
			e.CursorX = 0
		}
	case CursorMoveLineEnd:
		e.CursorX = len(e.CurrentRow().Text)
	case CursorMovePageUp:
		if e.Height > e.CursorY {
			e.CursorY = 1
		} else {
			e.CursorY -= e.Height - 2
		}
	case CursorMovePageDown:
		e.CursorY += e.Height - 2
		e.OffsetY += e.Height

		if e.CursorY > len(e.Buffer) {
			e.CursorY = len(e.Buffer) - 1
		}
	}

	// Prevent the user from moving past the end of the line.
	rowLength := len(e.CurrentRow().Text)
	if e.CursorX > rowLength {
		e.CursorX = rowLength
	}
}

/* -------------------------------- Internal -------------------------------- */

func (e *Editor) CurrentRow() *BufferLine {
	return &e.Buffer[e.CursorY-1]
}

// SetStatusMessage sets the status message and the time it was set at.
func (e *Editor) SetStatusMessage(format string, args ...interface{}) {
	e.StatusMessage = fmt.Sprintf(format, args...)
	e.StatusMessageTime = time.Now()
}

//##highlighter.go
type HighlightType int

const (
	HighlightTypeNormal HighlightType = iota
	HighlightTypePrimaryKeyword
	HighlightTypeSecondaryKeyword
	HighlightTypeNumber
	HighlightTypeString
	HighlightTypeComment
)

func (t HighlightType) Color() termbox.Attribute {
	switch t {
	case HighlightTypePrimaryKeyword:
		return termbox.ColorRed
	case HighlightTypeSecondaryKeyword:
		return termbox.ColorMagenta
	case HighlightTypeNumber:
		return termbox.ColorBlue
	case HighlightTypeString:
		return termbox.ColorGreen
	case HighlightTypeComment:
		return termbox.ColorCyan
	default:
		return termbox.ColorDefault
	}
}

func IsSeparator(c rune) bool {
	switch c {
	case ' ', ',', '.', ';', '(', ')', '[', ']', '+', '-', '/', '*', '=', '%':
		return true
	default:
		return false
	}
}

func HighlightLine(l *BufferLine, s *Syntax) {
	H := &l.Highlighting

	inString := false
	afterSeparator := true

	text := []rune(l.DisplayText)
	for i := 0; i < len(text); i++ {
		c := text[i]

		lastHT := HighlightTypeNormal
		if i > 0 {
			lastHT = (*H)[i-1]
		}

		// If we are already within a string, keep highlighting until we hit
		// another quote character.
		if inString {
			(*H)[i] = HighlightTypeString

			if c == '"' {
				inString = false
			}

			continue
		}

		// If we hit the beginning of a single line comment, highlight the rest
		// of the line and break out of the loop.
		scsPattern := &s.Patterns.SingleLineCommentStart
		scsLength := len(*scsPattern)
		if i+scsLength <= len(text) && string(text[i:i+scsLength]) == *scsPattern {
			for j := i; j < len(text); j++ {
				(*H)[j] = HighlightTypeComment
			}

			break
		}

		// If we hit a quotation mark, set inString to true and highlight it.
		if c == '"' || c == '\'' {
			(*H)[i] = HighlightTypeString
			inString = true
			continue
		}

		// If our character is a digit, is after a separator or trailing another
		// digit, or is a decimal trailing a digit, highlight it as a number.
		if unicode.IsDigit(c) &&
			(afterSeparator || lastHT == HighlightTypeNumber) ||
			(c == '.' && lastHT == HighlightTypeNumber) {
			(*H)[i] = HighlightTypeNumber
			continue
		}

		if afterSeparator {
			for _, k := range s.Keywords {
				kl := len(k)

				tail := ' '
				if i+kl < len(text) {
					tail = text[i+kl]
				}

				if i+kl <= len(text) && string(text[i:i+kl]) == k && IsSeparator(tail) {
					for j := 0; j < kl; j++ {
						(*H)[i+j] = HighlightTypeSecondaryKeyword
					}
				}
			}
		}

		afterSeparator = IsSeparator(c)
	}
}

//##support.go
const (
	ProgramName    string = "Atto"
	ProgramVersion string = "0.2.3"
	ProgramAuthor  string = "Jon Palmisciano <jonpalmisc@gmail.com>"
)

// FileType represents a type of file.
type FileType string

const (
	FileTypeMakefile FileType = "Makefile"
	FileTypeCMake    FileType = "CMake"

	FileTypeGo       FileType = "Go"
	FileTypeGoModule FileType = "Go Module"

	// -- C/C++ --
	FileTypeC   FileType = "C"
	FileTypeCPP FileType = "C++"

	// -- Text Files --
	FileTypeMarkdown  FileType = "Markdown"
	FileTypePlaintext FileType = "Plaintext"

	FileTypeUnknown FileType = "Unknown"
)

// GuessFileType attempts to deduce a file's type from its name and extension.
func GuessFileType(name string) FileType {

	// Handle filetypes which have specific names.
	switch name {
	case "Makefile":
		return FileTypeMakefile
	case "CMakeLists.txt":
		return FileTypeCMake
	}

	parts := strings.Split(name, ".")

	// Return unknown if the file has no extension and wasn't matched earlier.
	if len(parts) < 2 {
		return FileTypeUnknown
	}

	// Attempt to determine the file's type by the extension.
	switch parts[1] {
	case "go":
		return FileTypeGo
	case "mod":
		return FileTypeGoModule
	case "h", "c":
		return FileTypeC
	case "hpp", "cpp", "cc":
		return FileTypeCPP
	case "md":
		return FileTypeMarkdown
	case "txt":
		return FileTypePlaintext
	}

	return FileTypeUnknown
}

//##syntax.go
// SyntaxPatterns is used to define syntax patterns for the highlighter.
type SyntaxPatterns struct {
	SingleLineCommentStart string
	MultiLineCommendStart  string
	MultiLineCommentEnd    string
}

// Syntax represent's a language syntax for highlighting purposes.
type Syntax struct {
	Keywords []string
	Patterns SyntaxPatterns
}

// Syntax definitions are temporarily hardcoded until support for  language
// definition files is added!

// SyntaxC defines the syntax of the C language.
var SyntaxC Syntax = Syntax{
	Keywords: []string{
		"#define", "#include", "NULL", "auto", "break", "case", "char", "const",
		"continue", "default", "do", "double", "else", "enum", "extern", "float",
		"for", "goto", "if", "int", "long", "register", "return", "short",
		"signed", "sizeof", "static", "struct", "switch", "typedef", "union",
		"unsigned", "void", "volatile", "while",
	},
	Patterns: SyntaxPatterns{
		SingleLineCommentStart: "//",
		MultiLineCommendStart:  "/*",
		MultiLineCommentEnd:    "*/",
	},
}

// SyntaxGo defines the syntax of the Go language.
var SyntaxGo Syntax = Syntax{
	Keywords: []string{
		"append", "bool", "break", "byte", "cap", "case", "chan", "close",
		"complex", "complex128", "complex64", "const", "continue", "copy",
		"default", "defer", "delete", "else", "error", "fallthrough", "false",
		"float32", "float64", "for", "func", "go", "goto", "if", "imag",
		"import", "int", "int16", "int32", "int64", "int8", "interface", "len",
		"make", "map", "new", "nil", "package", "panic", "range", "real",
		"recover", "return", "rune", "select", "string", "struct", "switch",
		"true", "type", "uint", "uint16", "uint32", "uint64", "uint8", "uintptr",
		"var",
	},
	Patterns: SyntaxPatterns{
		SingleLineCommentStart: "//",
		MultiLineCommendStart:  "/*",
		MultiLineCommentEnd:    "*/",
	},
}
