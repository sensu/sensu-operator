package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/sensu/sensu-go/cli"
	"github.com/sensu/sensu-go/cli/client"
	"github.com/sensu/sensu-go/cli/client/config/basic"
	"github.com/sensu/sensu-go/types"
	"github.com/sirupsen/logrus"
)

type sensuAPITestClient struct {
	client.APIClient
	namespaceError error
}

func (c *sensuAPITestClient) CreateAccessToken(url, userid, secret string) (*types.Tokens, error) {
	return &types.Tokens{
		Access: "fake",
	}, nil
}

func (c *sensuAPITestClient) FetchNamespace(ns string) (*types.Namespace, error) {
	return &types.Namespace{
		Name: ns,
	}, c.namespaceError
}

func (c *sensuAPITestClient) CreateNamespace(ns *types.Namespace) error {
	return nil
}

func TestNew(t *testing.T) {
	conf := basic.Config{
		Cluster: basic.Cluster{
			APIUrl: "http://testCluster-api.testnamespace.svc:8080",
		},
		Profile: basic.Profile{
			Format:    "json",
			Namespace: "testnamespace",
		},
	}
	sensuCliClient := client.New(&conf)
	logger := logrus.WithFields(logrus.Fields{
		"component": "cli-client",
	})

	type args struct {
		clusterName string
		namespace   string
	}
	tests := []struct {
		name string
		args args
		want *SensuClient
	}{
		{
			"valid",
			args{
				"testCluster",
				"testnamespace",
			},
			&SensuClient{
				logger:      logrus.WithField("pkg", "sensu_client").WithField("cluster-name", "testCluster"),
				clusterName: "testCluster",
				namespace:   "testnamespace",
				sensuCli: &cli.SensuCli{
					Client: sensuCliClient,
					Config: &conf,
					Logger: logger,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.clusterName, tt.args.namespace, tt.args.namespace)
			if got.clusterName != tt.want.clusterName {
				t.Errorf("clustername should match got %s, want %s",
					got.clusterName,
					tt.want.clusterName,
				)
			}
			if got.sensuCli.Config.APIUrl() != tt.want.sensuCli.Config.APIUrl() {
				t.Errorf("sensuCli.Config.Apiurl should match got %s, want %s",
					got.sensuCli.Config.APIUrl(),
					tt.want.sensuCli.Config.APIUrl(),
				)
			}
			if got.sensuCli.Config.Namespace() != tt.want.sensuCli.Config.Namespace() {
				t.Errorf("sensuCli.Profile.Namespace should match got %s, want %s",
					got.sensuCli.Config.Namespace(),
					tt.want.sensuCli.Config.Namespace(),
				)
			}
		})
	}
}

func TestSensuClient_makeFullyQualifiedSensuClientURL(t *testing.T) {
	conf := basic.Config{
		Cluster: basic.Cluster{
			APIUrl: "http://testCluster.testnamespace.svc:8080",
		},
		Profile: basic.Profile{
			Format:    "json",
			Namespace: "testnamespace",
		},
	}
	sensuCliClient := client.New(&conf)
	logger := logrus.WithFields(logrus.Fields{
		"component": "cli-client",
	})

	type fields struct {
		logger      *logrus.Entry
		clusterName string
		namespace   string
		sensuCli    *cli.SensuCli
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"test",
			fields{
				logrus.WithField("pkg", "sensu_client").WithField("cluster-name", "testCluster"),
				"testCluster",
				"testnamespace",
				&cli.SensuCli{
					Client: sensuCliClient,
					Config: &conf,
					Logger: logger,
				},
			},
			"testCluster-api.testnamespace.svc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SensuClient{
				logger:      tt.fields.logger,
				clusterName: tt.fields.clusterName,
				namespace:   tt.fields.namespace,
				sensuCli:    tt.fields.sensuCli,
			}
			if got := s.makeFullyQualifiedSensuClientURL(); got != tt.want {
				t.Errorf("SensuClient.makeFullyQualifiedSensuClientURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSensuClient_ensureCredentials(t *testing.T) {
	conf := basic.Config{
		Cluster: basic.Cluster{
			APIUrl: "http://testCluster.testnamespace.svc:8080",
			Tokens: &types.Tokens{
				Access: "fake",
			},
		},
		Profile: basic.Profile{
			Format:    "json",
			Namespace: "testnamespace",
		},
	}
	confNoToken := basic.Config{
		Cluster: basic.Cluster{
			APIUrl: "http://testCluster.testnamespace.svc:8080",
		},
		Profile: basic.Profile{
			Format:    "json",
			Namespace: "testnamespace",
		},
	}
	sensuCliClient := client.New(&conf)

	logger := logrus.WithFields(logrus.Fields{
		"component": "cli-client",
	})

	type fields struct {
		logger      *logrus.Entry
		clusterName string
		namespace   string
		sensuCli    *cli.SensuCli
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"ensure credentials does nothing with existing tokens",
			fields{
				logrus.WithField("pkg", "sensu_client").WithField("cluster-name", "testCluster"),
				"testCluster",
				"testnamespace",
				&cli.SensuCli{
					Client: sensuCliClient,
					Config: &conf,
					Logger: logger,
				},
			},
			false,
		},
		{
			"ensure credentials attempts create with missing token",
			fields{
				logrus.WithField("pkg", "sensu_client").WithField("cluster-name", "testCluster"),
				"testCluster",
				"testnamespace",
				&cli.SensuCli{
					Client: &sensuAPITestClient{},
					Config: &confNoToken,
					Logger: logger,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SensuClient{
				logger:      tt.fields.logger,
				clusterName: tt.fields.clusterName,
				namespace:   tt.fields.namespace,
				sensuCli:    tt.fields.sensuCli,
				timeout:     2 * time.Second,
			}
			if err := s.ensureCredentials(); (err != nil) != tt.wantErr {
				t.Errorf("SensuClient.ensureCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSensuClient_ensureNamespace(t *testing.T) {
	conf := basic.Config{
		Cluster: basic.Cluster{
			APIUrl: "http://testCluster.testnamespace.svc:8080",
			Tokens: &types.Tokens{
				Access: "fake",
			},
		},
		Profile: basic.Profile{
			Format:    "json",
			Namespace: "testnamespace",
		},
	}
	// sensuCliClient := client.New(&conf)

	logger := logrus.WithFields(logrus.Fields{
		"component": "cli-client",
	})
	type fields struct {
		logger      *logrus.Entry
		clusterName string
		namespace   string
		sensuCli    *cli.SensuCli
		timeout     time.Duration
	}
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"ensure existing namespace does not fail",
			fields{
				logrus.WithField("pkg", "sensu_client").WithField("cluster-name", "testCluster"),
				"testCluster",
				"testnamespace",
				&cli.SensuCli{
					Client: &sensuAPITestClient{},
					Config: &conf,
					Logger: logger,
				},
				2 * time.Second,
			},
			args{
				"testsensunamespace",
			},
			false,
		},
		{
			"ensure non-existing namespace creates",
			fields{
				logrus.WithField("pkg", "sensu_client").WithField("cluster-name", "testCluster"),
				"testCluster",
				"testnamespace",
				&cli.SensuCli{
					Client: &sensuAPITestClient{
						namespaceError: fmt.Errorf("not found"),
					},
					Config: &conf,
					Logger: logger,
				},
				2 * time.Second,
			},
			args{
				"testsensunamespace",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SensuClient{
				logger:      tt.fields.logger,
				clusterName: tt.fields.clusterName,
				namespace:   tt.fields.namespace,
				sensuCli:    tt.fields.sensuCli,
				timeout:     tt.fields.timeout,
			}
			if err := s.ensureNamespace(tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("SensuClient.ensureNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
