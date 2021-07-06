package spinner

import (
	"context"
	"fmt"
	"time"
)

type Spinner struct {
	ctx       context.Context
	c         chan string
	msg       string
	msgPrefix string
}

func NewSpinner(ctx context.Context, msgPrefix string) *Spinner {
	return &Spinner{
		ctx:       ctx,
		c:         make(chan string),
		msgPrefix: msgPrefix,
	}
}
func (s *Spinner) Update(msg string) {
	s.c <- msg
}

// Start should be called as a go routine
func (s *Spinner) Start() {
	spin := `\|/-`
	idx := 0
	idx1 := 0
	suf := []string{"", ".", "..", "...", "...."}

	go func() {
		for {
			select {
			case msg := <-s.c:
				idx1 = 0
				s.msg = msg
			case <-s.ctx.Done():
				return
			}
		}
	}()

	for {
		fmt.Printf("\033[2K\r%c %s: %s %s", spin[idx], s.msgPrefix, s.msg, suf[idx1])
		time.Sleep(150 * time.Millisecond)
		idx++
		idx %= 4
		if idx%4 == 0 {
			idx1++
			idx1 %= len(suf)
		}

	}
}
