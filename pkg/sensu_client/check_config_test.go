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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equal(tt.args.c1, tt.args.c2); got != tt.want {
				if tt.args.c1.Command != tt.args.c2.Command ||
					tt.args.c1.HighFlapThreshold != tt.args.c2.HighFlapThreshold ||
					tt.args.c1.Interval != tt.args.c2.Interval ||
					tt.args.c1.LowFlapThreshold != tt.args.c2.LowFlapThreshold ||
					tt.args.c1.Publish != tt.args.c2.Publish ||
					tt.args.c1.ProxyEntityName != tt.args.c2.ProxyEntityName ||
					tt.args.c1.Stdin != tt.args.c2.Stdin ||
					tt.args.c1.Cron != tt.args.c2.Cron ||
					tt.args.c1.Ttl != tt.args.c2.Ttl ||
					tt.args.c1.Timeout != tt.args.c2.Timeout ||
					tt.args.c1.RoundRobin != tt.args.c2.RoundRobin ||
					tt.args.c1.OutputMetricFormat != tt.args.c2.OutputMetricFormat ||
					!cmp.Equal(tt.args.c1.ObjectMeta, tt.args.c2.ObjectMeta) {
					t.Errorf("equal() = checkconfigs not equal: diff %s", cmp.Diff(got, tt.want))
				}

				if len(tt.args.c1.Handlers) != len(tt.args.c2.Handlers) ||
					len(tt.args.c1.RuntimeAssets) != len(tt.args.c2.RuntimeAssets) ||
					len(tt.args.c1.Subscriptions) != len(tt.args.c2.Subscriptions) ||
					len(tt.args.c1.OutputMetricHandlers) != len(tt.args.c2.OutputMetricHandlers) ||
					len(tt.args.c1.EnvVars) != len(tt.args.c2.EnvVars) {
					t.Errorf("equal() = checkconfigs not equal: diff %s", cmp.Diff(got, tt.want))
				}
			}
		})
	}
}
