package utils_test

import (
	"os/user"
	"path"
	"testing"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
)

func TestToAbsolutPath(t *testing.T) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple",
			args: args{
				path: "/tmp/test",
			},
			want: "/tmp/test",
		},
		{
			name: "simple-relative",
			args: args{
				path: "tmp/test",
			},
			want: "tmp/test",
		},
		{
			name: "homedir-subpath",
			args: args{
				path: "~/test",
			},
			want: path.Join(user.HomeDir, "/test"),
		},
		{
			name: "homedir-direct",
			args: args{
				path: "~",
			},
			want: user.HomeDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ToAbsolutPath(tt.args.path); got != tt.want {
				t.Errorf("ToAbsolutPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
