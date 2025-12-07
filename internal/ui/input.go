package ui

import (
	"os"
)

type Key int

const (
	KeyNone Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyEsc
	KeyEnter
	KeySpace
	KeyQ
	KeyW
	KeyA
	KeyS
	KeyD
)

type Input struct {
	buf  [3]byte
	keys chan Key
	stop chan struct{}
}

func NewInput() *Input {
	inp := &Input{
		keys: make(chan Key, 16),
		stop: make(chan struct{}),
	}
	go inp.readLoop()
	return inp
}

func (inp *Input) Close() {
	close(inp.stop)
}

func (inp *Input) Poll() Key {
	select {
	case k := <-inp.keys:
		return k
	default:
		return KeyNone
	}
}

func (inp *Input) readLoop() {
	for {
		select {
		case <-inp.stop:
			return
		default:
		}

		n, err := os.Stdin.Read(inp.buf[:])
		if err != nil {
			return
		}

		if n == 0 {
			continue
		}

		if n == 1 {
			switch inp.buf[0] {
			case 27:
				inp.keys <- KeyEsc
			case 13, 10:
				inp.keys <- KeyEnter
			case 32:
				inp.keys <- KeySpace
			case 'q', 'Q':
				inp.keys <- KeyQ
			case 'w', 'W':
				inp.keys <- KeyW
			case 'a', 'A':
				inp.keys <- KeyA
			case 's', 'S':
				inp.keys <- KeyS
			case 'd', 'D':
				inp.keys <- KeyD
			}
			continue
		}

		if n == 3 && inp.buf[0] == 27 && inp.buf[1] == '[' {
			switch inp.buf[2] {
			case 'A':
				inp.keys <- KeyUp
			case 'B':
				inp.keys <- KeyDown
			case 'C':
				inp.keys <- KeyRight
			case 'D':
				inp.keys <- KeyLeft
			}
		}
	}
}
