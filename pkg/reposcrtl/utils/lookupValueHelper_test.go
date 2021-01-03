package utils_test

import (
	"os"
	"testing"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
)

func TestLookupValueByRef(t *testing.T) {
	type args struct {
		ref     string
		prepare func()
		ending  func()
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "env-exists",
			want: "testToken",
			args: args{
				ref: "ref+env://GITHUB_TOKEN  ",
				ending: func() {
					os.Unsetenv("GITHUB_TOKEN")
				},
				prepare: func() {
					os.Setenv("GITHUB_TOKEN", "testToken")
				},
			},
			wantErr: false,
		},
		{
			name: "cmd-exists",
			want: "testTokenOutput",
			args: args{
				ref: "ref+cmd://$( echo \"testTokenOutput\" )",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.prepare != nil {
				tt.args.prepare()
			}
			if tt.args.ending != nil {
				defer tt.args.ending()
			}
			got, err := utils.LookupValueByRef(tt.args.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupValueByRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LookupValueByRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
