package ui

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBoldRed = "\033[1;31m"
	ColorDim     = "\033[2m"
)

type Terminal struct {
	oldState *term.State
	width    int
	height   int
	buf      []byte
	mu       sync.Mutex
}

func NewTerminal() (*Terminal, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("make raw: %w", err)
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		return nil, fmt.Errorf("get size: %w", err)
	}

	t := &Terminal{
		oldState: oldState,
		width:    w,
		height:   h,
		buf:      make([]byte, 0, w*h*4),
	}

	t.HideCursor()
	t.Clear()

	return t, nil
}

func (t *Terminal) Close() {
	t.ShowCursor()
	t.Clear()
	_ = term.Restore(int(os.Stdin.Fd()), t.oldState)
}

func (t *Terminal) Width() int  { return t.width }
func (t *Terminal) Height() int { return t.height }

func (t *Terminal) RefreshSize() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		t.width = w
		t.height = h
	}
}

func (t *Terminal) Clear() {
	fmt.Fprint(os.Stdout, "\033[2J\033[H")
}

func (t *Terminal) HideCursor() {
	fmt.Fprint(os.Stdout, "\033[?25l")
}

func (t *Terminal) ShowCursor() {
	fmt.Fprint(os.Stdout, "\033[?25h")
}

func (t *Terminal) BeginFrame() {
	t.mu.Lock()
	t.buf = t.buf[:0]
	// Move cursor home + clear screen in buffer (not direct stdout)
	t.buf = append(t.buf, "\033[H\033[2J"...)
}

func (t *Terminal) EndFrame() {
	_, _ = os.Stdout.Write(t.buf)
	t.mu.Unlock()
}

func (t *Terminal) SetCell(x, y int, ch rune, color string) {
	if x < 0 || x >= t.width || y < 0 || y >= t.height {
		return
	}
	t.buf = append(t.buf, fmt.Sprintf("\033[%d;%dH%s%c%s", y+1, x+1, color, ch, ColorReset)...)
}

func (t *Terminal) WriteStr(x, y int, s string, color string) {
	if y < 0 || y >= t.height {
		return
	}
	t.buf = append(t.buf, fmt.Sprintf("\033[%d;%dH%s%s%s", y+1, x+1, color, s, ColorReset)...)
}

func (t *Terminal) FillRow(y int, ch rune) {
	if y < 0 || y >= t.height {
		return
	}
	t.buf = append(t.buf, fmt.Sprintf("\033[%d;1H", y+1)...)
	for i := 0; i < t.width; i++ {
		t.buf = append(t.buf, string(ch)...)
	}
}
