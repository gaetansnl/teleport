// Copyright 2021 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package terminal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gravitational/teleport/api/profile"
	"github.com/gravitational/teleport/lib/terminal"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"

	"github.com/gravitational/teleport/lib/client"

	v1 "github.com/gravitational/teleport/lib/terminal/api/protogen/golang/v1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	profileDir = "/home/alexey/go/src/github.com/gravitational/_terminal"
)

func TestStart(t *testing.T) {
	tests := []struct {
		name string
		cfg  terminal.Config
	}{
		{
			name: "Unix socket",
			cfg: terminal.Config{
				Addr:    fmt.Sprintf("unix://%v/terminal.sock", t.TempDir()),
				HomeDir: fmt.Sprintf("%v/", t.TempDir()),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			wait := make(chan error)
			go func() {
				err := terminal.Start(ctx, test.cfg)
				wait <- err
			}()

			cc, err := grpc.Dial(test.cfg.Addr, grpc.WithInsecure())
			require.NoError(t, err)

			term := v1.NewTerminalServiceClient(cc)
			_, err = term.ListClusters(ctx, &v1.ListClustersRequest{})
			require.NoError(t, err)

			cluster, err := term.CreateCluster(ctx, &v1.CreateClusterRequest{Name: "platform.teleport.sh"})
			fmt.Print("ZZZZZZZZZZZZZZZZZZZ", cluster)

			require.NoError(t, err)

			defer func() {
				cancel() // Stop the server.
				require.NoError(t, <-wait)
			}()

		})
	}
}

func FTestStart(t *testing.T) {
	cc, err := grpc.Dial("unix://tshdffdsf", grpc.WithInsecure())
	fmt.Print("AAAAAAAAAAAAAAAAAAAAAAA", cc.GetState())
	require.Error(t, err)

	defer cc.Close()
}

func CreateClusterProfile(clusterName string) error {
	c := client.MakeDefaultConfig()
	c.WebProxyAddr = clusterName
	c.HomePath = profileDir
	c.KeysDir = profileDir

	tc, err := client.NewClient(c)
	if err != nil {
		return trace.Wrap(err)
	}

	if _, err := tc.Ping(context.TODO()); err != nil {
		return trace.Wrap(err)
	}

	if err := tc.SaveProfile(profileDir, false); err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func GetClusterStatuses() ([]terminal.ClusterStatus, error) {
	pNames, err := profile.ListProfileNames(profileDir)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	statuses := []terminal.ClusterStatus{}
	for _, name := range pNames {
		status, err := client.StatusFromFile(profileDir, name)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		clusterStatus := terminal.ClusterStatus{
			ProfileStatus: *status,
		}

		statuses = append(statuses, clusterStatus)
	}

	return statuses, nil
}

func ListClusters() ([]*v1.Cluster, error) {
	statuses, err := GetClusterStatuses()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	clusters := []*v1.Cluster{}
	for _, sts := range statuses {
		cluster := &v1.Cluster{
			Name:      sts.Name,
			Connected: sts.IsExpired(clockwork.NewRealClock()),
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func NewClusterProfile(clusterName string) (*profile.Profile, error) {
	c := client.MakeDefaultConfig()
	c.WebProxyAddr = clusterName
	c.HomePath = profileDir
	c.KeysDir = profileDir

	tc, err := client.NewClient(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	pingResponse, err := tc.Ping(context.TODO())
	if err != nil {
		return nil, trace.Wrap(err)
	}

	profile := profile.Profile{
		WebProxyAddr:           clusterName,
		SiteName:               clusterName,
		SSHProxyAddr:           pingResponse.Proxy.SSH.PublicAddr,
		KubeProxyAddr:          pingResponse.Proxy.Kube.PublicAddr,
		PostgresProxyAddr:      pingResponse.Proxy.DB.PostgresPublicAddr,
		MySQLProxyAddr:         pingResponse.Proxy.DB.MySQLPublicAddr,
		ALPNSNIListenerEnabled: pingResponse.Proxy.ALPNSNIListenerEnabled,
	}

	return &profile, nil
}

func FTestMama(t *testing.T) {
	c := client.MakeDefaultConfig()
	//c.SiteName = "https://platform.teleport.sh"
	c.WebProxyAddr = "platform.teleport.sh"
	c.HomePath = "/home/alexey/go/src/github.com/gravitational/_terminal"
	c.AuthConnector = "okta"
	c.KeysDir = "/home/alexey/go/src/github.com/gravitational/_terminal"

	//	profile, profiles, err := client.Status("/home/alexey/go/src/github.com/gravitational/_terminal", "https://platform.teleport.sh")
	//	require.NoError(t, err)

	tc, err := client.NewClient(c)
	require.NoError(t, err)

	webclient, err := tc.Ping(context.TODO())
	require.NoError(t, err)

	//webclient.Auth.

	fmt.Print("AAAAAAAAAAAAAAAAAAAAA", webclient)
	//	require.Error(t, err)

	key, err := tc.Login(context.TODO())
	require.NoError(t, err)
	tc.AddKey(key)

	err = tc.ActivateKey(context.TODO(), key)
	require.NoError(t, err)

	tc.SaveProfile(c.HomePath, true)

	servers, err := tc.ListNodes(context.TODO())
	require.NoError(t, err)
	fmt.Print("SERVERS:", servers)
	require.Error(t, err)

}

func FTestBrother(t *testing.T) {
	c := client.MakeDefaultConfig()
	//c.SiteName = "https://platform.teleport.sh"
	//c.AuthMethods
	c.WebProxyAddr = "platform.teleport.sh"
	c.HomePath = profileDir
	//c.AuthConnector = "okta"
	c.KeysDir = profileDir
	c.Username = "alexey@goteleport.com"

	//active, _, err := client.Status("/home/alexey/go/src/github.com/gravitational/_terminal", "platform.teleport.sh")
	err := c.LoadProfile("/home/alexey/go/src/github.com/gravitational/_terminal", "platform.teleport.sh")

	require.NoError(t, err)

	tc, err := client.NewClient(c)
	require.NoError(t, err)
	err = tc.LoadKeyForCluster("platform.teleport.sh")

	servers, err := tc.ListNodes(context.TODO())
	require.NoError(t, err)
	fmt.Print("SERVERS:", servers)
}

func FTestAddCluster(t *testing.T) {

	//	AddCluster("platform.teleport.sh")

	profileNames, err := profile.ListProfileNames(profileDir)
	require.NoError(t, err)

	for _, name := range profileNames {
		status, err := client.StatusFromFile(profileDir, name)
		require.NoError(t, err)
		fmt.Print("STATUS:", status)

		//pr, err := profile.FromDir(profileDir, v)
		//require.NoError(t, pr.)

	}

	_, profs, err := client.Status("/home/alexey/go/src/github.com/gravitational/_terminal", "")
	//fmt.Print("AAAAAAAAAAAAAAAAAAAAAAAAA:", profs)
	fmt.Print("ACTIVE:", profs[0])
	require.Error(t, err)

}

func FTestPapa(t *testing.T) {
	profileNames, err := profile.ListProfileNames(profileDir)
	require.NoError(t, err)

	for _, name := range profileNames {
		status, err := client.StatusFromFile(profileDir, name)
		require.NoError(t, err)
		fmt.Print("STATUS:", status)

		//pr, err := profile.FromDir(profileDir, v)
		//require.NoError(t, pr.)

	}

	_, profs, err := client.Status("/home/alexey/go/src/github.com/gravitational/_terminal", "")
	//fmt.Print("AAAAAAAAAAAAAAAAAAAAAAAAA:", profs)
	fmt.Print("ACTIVE:", profs[0])
	require.Error(t, err)

}

/*
https://platform.teleport.sh/webapi/find

// Get the status of the active profile as well as the status
	// of any other proxies the user is logged into.
	profile, profiles, err := client.Status(cf.HomePath, cf.Proxy)
	if err != nil {
		if !trace.IsNotFound(err) {
			return trace.Wrap(err)
		}
	}

*/
