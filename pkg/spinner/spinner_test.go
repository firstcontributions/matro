package spinner

import (
	"context"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestNewSpinner(t *testing.T) {
	type args struct {
		msgPrefix string
	}
	tests := []struct {
		name string
		args args
		want *Spinner
	}{
		{
			name: "should return an instance of NewSpinner",
			args: args{msgPrefix: "pref"},
			want: &Spinner{
				msg:       "",
				msgPrefix: "pref",
				c:         make(chan string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpinner(tt.args.msgPrefix); !reflect.DeepEqual(got.msgPrefix, tt.want.msgPrefix) ||
				!reflect.DeepEqual(got.msg, tt.want.msg) {
				t.Errorf("NewSpinner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpinner_Update(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		c         chan string
		msg       string
		msgPrefix string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		shouldStart bool
		want        error
	}{
		{
			name: "should return an error if start routine is not called",
			fields: fields{
				c: make(chan string),
			},
			shouldStart: false,
			args: args{
				msg: "hello world",
			},
			want: errNotStarted,
		},
		{
			name: "should pass the message to the channel",
			fields: fields{
				c: make(chan string),
			},
			args: args{
				msg: "hello world",
			},
			shouldStart: true,
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Spinner{
				c:         tt.fields.c,
				msg:       tt.fields.msg,
				msgPrefix: tt.fields.msgPrefix,
			}
			if tt.shouldStart {
				go s.Start(ctx)
				time.Sleep(1 * time.Second)
			}
			if got := s.Update(tt.args.msg); got != tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpinner_Start(t *testing.T) {
	t.Run("should close the channel once the ctx is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		s := NewSpinner("a")
		go s.Start(ctx)
		cancel()
		time.Sleep(1 * time.Second)
		if !isChanClosed(s.c) {
			t.Errorf("channel should have closed once the context is cancelled")
		}
	})
}

func isChanClosed(ch interface{}) bool {
	if reflect.TypeOf(ch).Kind() != reflect.Chan {
		panic("only channels!")
	}

	// get interface value pointer, from cgo_export
	// typedef struct { void *t; void *v; } GoInterface;
	// then get channel real pointer
	cptr := *(*uintptr)(unsafe.Pointer(
		unsafe.Pointer(uintptr(unsafe.Pointer(&ch)) + unsafe.Sizeof(uint(0))),
	))

	// this function will return true if chan.closed > 0
	// see hchan on https://github.com/golang/go/blob/master/src/runtime/chan.go
	// type hchan struct {
	// qcount   uint           // total data in the queue
	// dataqsiz uint           // size of the circular queue
	// buf      unsafe.Pointer // points to an array of dataqsiz elements
	// elemsize uint16
	// closed   uint32
	// **

	cptr += unsafe.Sizeof(uint(0)) * 2
	cptr += unsafe.Sizeof(unsafe.Pointer(uintptr(0)))
	cptr += unsafe.Sizeof(uint16(0))
	return *(*uint32)(unsafe.Pointer(cptr)) > 0
}
