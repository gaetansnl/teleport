package integration

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	apidefaults "github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/backend/lite"
	"github.com/gravitational/teleport/lib/service"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/teleport/lib/utils"
	"github.com/stretchr/testify/require"
)

func newNodeConfig(t *testing.T, authAddr utils.NetAddr) *service.Config {
	config := service.MakeDefaultConfig()
	config.AWSToken = "test_token"
	config.SSH.Enabled = true
	config.SSH.Addr.Addr = net.JoinHostPort(Host, ports.Pop())
	config.Auth.Enabled = false
	config.Proxy.Enabled = false
	config.DataDir = t.TempDir()
	config.AuthServers = append(config.AuthServers, authAddr)
	config.CachePolicy.Enabled = true
	config.ClientTimeout = time.Second
	return config
}

func newAuthConfig(t *testing.T) *service.Config {
	var err error
	storageConfig := backend.Config{
		Type: lite.GetName(),
		Params: backend.Params{
			"path":               t.TempDir(),
			"poll_stream_period": 50 * time.Millisecond,
		},
	}

	config := service.MakeDefaultConfig()
	config.SSH.Enabled = false
	config.DataDir = t.TempDir()
	config.Auth.SSHAddr.Addr = net.JoinHostPort(Host, ports.Pop())
	config.Auth.PublicAddrs = []utils.NetAddr{
		{
			AddrNetwork: "tcp",
			Addr:        Host,
		},
	}
	config.Auth.ClusterName, err = services.NewClusterNameWithRandomID(types.ClusterNameSpecV2{
		ClusterName: "testcluster",
	})
	config.AuthServers = append(config.AuthServers, config.Auth.SSHAddr)
	config.Auth.StorageConfig = storageConfig
	config.Auth.StaticTokens, err = types.NewStaticTokens(types.StaticTokensSpecV2{
		StaticTokens: []types.ProvisionTokenV1{
			{
				Roles: []types.SystemRole{"Proxy", "Node"},
				Token: "foo",
			},
		},
	})
	require.NoError(t, err)
	config.Proxy.Enabled = false
	/*
		config.Proxy.Enabled = true
		config.Proxy.DisableWebInterface = true
		config.Proxy.DisableWebService = true
		config.Proxy.DisableReverseTunnel = true
		config.Proxy.SSHAddr.Addr = net.JoinHostPort(Host, ports.Pop())
		config.Proxy.WebAddr.Addr = net.JoinHostPort(Host, ports.Pop())
	*/
	config.CachePolicy.Enabled = true
	return config
}

func getIID(t *testing.T) imds.InstanceIdentityDocument {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoError(t, err)
	imdsClient := imds.NewFromConfig(cfg)
	output, err := imdsClient.GetInstanceIdentityDocument(context.TODO(), nil)
	require.NoError(t, err)
	return output.InstanceIdentityDocument
}

func TestSimplifiedNodeJoin(t *testing.T) {
	iid := getIID(t)
	token, err := types.NewProvisionTokenFromSpec(
		"test_token",
		time.Now().Add(time.Hour),
		types.ProvisionTokenSpecV2{
			Roles: []types.SystemRole{types.RoleNode},
			Allow: []*types.TokenRule{
				&types.TokenRule{
					AWSAccount: iid.AccountID,
					AWSRegions: []string{iid.Region},
				},
			},
		})
	require.NoError(t, err)

	authConfig := newAuthConfig(t)
	authSvc, err := service.NewTeleport(authConfig)
	require.NoError(t, err)
	require.NoError(t, authSvc.Start())

	authServer := authSvc.GetAuthServer()

	err = authServer.UpsertToken(context.Background(), token)
	require.NoError(t, err)

	nodes, err := authServer.GetNodes(context.Background(), apidefaults.Namespace)
	require.NoError(t, err)
	require.Empty(t, nodes)

	nodeConfig := newNodeConfig(t, authConfig.Auth.SSHAddr)
	nodeSvc, err := service.NewTeleport(nodeConfig)
	require.NoError(t, err)
	require.NoError(t, nodeSvc.Start())

	require.Eventually(t, func() bool {
		nodes, err := authServer.GetNodes(context.Background(), apidefaults.Namespace)
		require.NoError(t, err)
		return len(nodes) > 0
	}, time.Minute, 5*time.Second, "waiting for node to join cluster")
}
