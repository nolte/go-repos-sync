package utils_test

import (
	"testing"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
)

func TestContains(t *testing.T) {
	type args struct {
		s          []string
		searchterm string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "exists",
		args: args{
			s:          []string{"a", "b", "c"},
			searchterm: "a",
		},
		want: true,
	}, {
		name: "exists",
		args: args{
			s:          []string{"a", "b", "c"},
			searchterm: "b",
		},
		want: true,
	}, {
		name: "exists",
		args: args{
			s:          []string{"a", "b", "c"},
			searchterm: "c",
		},
		want: true,
	}, {
		name: "exists",
		args: args{
			s:          []string{"a", "b", "c"},
			searchterm: "xxx",
		},
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Contains(tt.args.s, tt.args.searchterm); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
