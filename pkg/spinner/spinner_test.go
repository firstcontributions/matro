package spinner

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestNewSpinner(t *testing.T) {
	type args struct {
		ctx       context.Context
		msgPrefix string
	}
	ctx := context.TODO()
	tests := []struct {
		name string
		args args
		want *Spinner
	}{
		{
			name: "should return an instance of NewSpinner",
			args: args{ctx: ctx, msgPrefix: "pref"},
			want: &Spinner{
				ctx:       ctx,
				msg:       "",
				msgPrefix: "pref",
				c:         make(chan string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpinner(tt.args.ctx, tt.args.msgPrefix); !reflect.DeepEqual(got.ctx, tt.want.ctx) ||
				!reflect.DeepEqual(got.msgPrefix, tt.want.msgPrefix) ||
				!reflect.DeepEqual(got.msg, tt.want.msg) {
				t.Errorf("NewSpinner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpinner_Update(t *testing.T) {
	type fields struct {
		ctx       context.Context
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
				ctx: context.TODO(),
				c:   make(chan string),
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
				ctx: context.TODO(),
				c:   make(chan string),
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
				ctx:       tt.fields.ctx,
				c:         tt.fields.c,
				msg:       tt.fields.msg,
				msgPrefix: tt.fields.msgPrefix,
			}
			if tt.shouldStart {
				go s.Start()
				fmt.Println(s)
			}
			if got := s.Update(tt.args.msg); got != tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
			if tt.shouldStart {
				if tt.args.msg != s.msg {
					t.Errorf("Update() did not update the msg, current= %v, want = %v", s.msg, tt.args.msg)
				}
			}
		})
	}
}

func TestSpinner_Start(t *testing.T) {
	t.Run("should close the channel once the ctx is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		s := NewSpinner(ctx, "a")
		go s.Start()
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
