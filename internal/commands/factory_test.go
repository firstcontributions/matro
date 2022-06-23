package commands

import (
	"reflect"
	"testing"
)

func TestGetCmd(t *testing.T) {

	tests := []struct {
		name string
		args []string
		want Command
	}{
		{
			name: "should return server command if args[0] == server",
			args: []string{"server"},
			want: server,
		},
		{
			name: "should return relay command if args[0] == relay",
			args: []string{"relay"},
			want: relay,
		},
		{
			name: "should return version command if args[0] == version",
			args: []string{"version"},
			want: versionCmd,
		},
		{
			name: "should return help command no args",
			args: nil,
			want: helpCmd,
		},
		{
			name: "should return help command for unknown args",
			args: []string{"asdf"},
			want: helpCmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCmd(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
