package client

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	sensu_api_core_v2 "github.com/sensu/sensu-go/api/core/v2"
)

func Test_equal(t *testing.T) {
	type args struct {
		c1 *sensu_api_core_v2.CheckConfig
		c2 *sensu_api_core_v2.CheckConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test round_robin isn't equal",
			args{
				&sensu_api_core_v2.CheckConfig{
					RoundRobin: false,
				},
				&sensu_api_core_v2.CheckConfig{
					RoundRobin: true,
				},
			},
			false,
		},
		{
			"test round_robin is equal",
			args{
				&sensu_api_core_v2.CheckConfig{
					RoundRobin: true,
				},
				&sensu_api_core_v2.CheckConfig{
					RoundRobin: true,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmp.Equal(tt.args.c1, tt.args.c2); got != tt.want {
				t.Errorf("equal() = checkconfigs equality issue, wanted %t, got %t: diff %s", tt.want, got, cmp.Diff(tt.args.c1, tt.args.c2))
			}
		})
	}
}
