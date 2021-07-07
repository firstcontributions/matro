package spinner

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errNotStarted = errors.New("spinner.Start routine should be called before Update")

// Spinner implements a command line spinner
type Spinner struct {
	ctx       context.Context
	c         chan string
	msg       string
	msgPrefix string
	started   bool
}

// NewSpinner creates an instance of cmd line spinner
// Args:
// 		ctx: a system context with cancel, the spinner stops
// when context is cancelled
// 		msgPrefix:  a message prefix to be printed on
// begining of every line
func NewSpinner(ctx context.Context, msgPrefix string) *Spinner {
	return &Spinner{
		ctx:       ctx,
		c:         make(chan string),
		msgPrefix: msgPrefix,
	}
}

// Update updates the current message that the spinner is printing
// by default the msg will be an empty string
// 	Args:
// 		msg: the message to be printed
func (s *Spinner) Update(msg string) error {
	if !s.started {
		return errNotStarted
	}
	s.c <- msg
	return nil
}

// Start should be called as a go routine
// this func listen for new messages on the spinner channel
// change the message in the cmd line if a new message is available
// also responsible for an ascii spinner
func (s *Spinner) Start() {
	s.started = true
	// ASCII spinner chars
	spin := `\|/-`
	// Idx used to select a char from ascii spin chars
	// this will be updated in every 150ms
	idx := 0
	// selector for suffix ... spinner
	suffIdx := 0
	suf := []string{"", ".", "..", "...", "...."}

	// events listener routine, listen for new msgs and context changes
	go func() {
		for {
			select {
			case msg := <-s.c:
				suffIdx = 0
				s.msg = msg
			case <-s.ctx.Done():
				close(s.c)
				return
			}
		}
	}()

	for {
		fmt.Printf("\033[2K\r%c %s: %s %s", spin[idx], s.msgPrefix, s.msg, suf[suffIdx])
		time.Sleep(150 * time.Millisecond)
		idx++
		idx %= 4
		if idx%4 == 0 {
			//  this ensures the suffix spinner changes slowly
			suffIdx++
			suffIdx %= len(suf)
		}
	}
}
