package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/objectrocket/sensu-operator/pkg/util/k8sutil"
	"github.com/sensu/sensu-go/cli"
	"github.com/sensu/sensu-go/cli/client"
	"github.com/sensu/sensu-go/cli/client/config/basic"
	"github.com/sensu/sensu-go/types"
)

const defaultTimeout = 10 * time.Second

// SensuClient is the sensu client used to interact with a sensu server within
// a kubernetes cluster, within a certain k8s namespace
//
// NOTE: ** this is tied to a sensu namespace, so a new SensuClient is needed
//   between sensu namespaces **
//
// TODO: add a factory method that returns a sensuclient, and allows easy switching
// of sensu namespaces
type SensuClient struct {
	logger *logrus.Entry

	clusterName string
	namespace   string

	sensuCli *cli.SensuCli

	timeout time.Duration
}

var (
	errSensuClusterObjectNotFound = errors.New("not found")
)

// New will return a new SensuClient tied to a specific cluster within a k8s
// namespace, and tied to a specific sensu namespace.
func New(clusterName, namespace string, sensuNamespace string) *SensuClient {
	sClient := &SensuClient{
		logger:      logrus.WithField("pkg", "sensu_client").WithField("cluster-name", clusterName),
		clusterName: clusterName,
		namespace:   namespace,
		timeout:     time.Duration(defaultTimeout),
	}

	conf := basic.Config{
		Cluster: basic.Cluster{
			APIUrl:  fmt.Sprintf("http://%s:8080", sClient.makeFullyQualifiedSensuClientURL()),
			Edition: "enterprise",
		},
		Profile: basic.Profile{
			Format:    "json",
			Namespace: sensuNamespace,
		},
	}

	sensuCliClient := client.New(&conf)
	logger := logrus.WithFields(logrus.Fields{
		"component": "cli-client",
	})

	sClient.sensuCli = &cli.SensuCli{
		Client: sensuCliClient,
		Config: &conf,
		Logger: logger,
	}

	return sClient
}

// SetTimeout will set the timeout for sensu cluster api operations
func (s *SensuClient) SetTimeout(t time.Duration) {
	s.timeout = t
}

func (s *SensuClient) makeFullyQualifiedSensuClientURL() string {
	return fmt.Sprintf("%s.%s.svc", k8sutil.APIServiceName(s.clusterName), s.namespace)
}

func (s *SensuClient) ensureCredentials() (err error) {
	var (
		tokens *types.Tokens
	)

	currentTokens := s.sensuCli.Config.Tokens()
	if currentTokens == nil || currentTokens.Access == "" {

		c1 := make(chan types.Tokens, 1)
		go func() {
			var tokens *types.Tokens
			if tokens, err = s.sensuCli.Client.CreateAccessToken(fmt.Sprintf("http://%s:8080", s.makeFullyQualifiedSensuClientURL()), "admin", "P@ssw0rd!"); err != nil {
				s.logger.Errorf("create token err: %+v", err)
				return
			}
			c1 <- *tokens
		}()

		select {
		case response := <-c1:
			tokens = &response
		case <-time.After(s.timeout):
			s.logger.Warnf("timeout from sensu server after 10 seconds")
		}

		if tokens == nil {
			return fmt.Errorf("failed to retrieve new access token from sensu server")
		}

		conf := basic.Config{
			Cluster: basic.Cluster{
				APIUrl:  s.sensuCli.Config.APIUrl(),
				Edition: "enterprise",
				Tokens:  tokens,
			},
			Profile: basic.Profile{
				Format:    "json",
				Namespace: s.sensuCli.Config.Namespace(),
			},
		}

		sensuCliClient := client.New(&conf)

		logger := logrus.WithFields(logrus.Fields{
			"component": "cli-client",
		})

		s.sensuCli = &cli.SensuCli{
			Client: sensuCliClient,
			Config: &conf,
			Logger: logger,
		}
	}
	return nil
}
