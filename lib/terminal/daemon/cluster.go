/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package daemon

import (
	"context"

	"github.com/gravitational/teleport/api/client/webclient"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/client"
)

// Cluster describes user settings and access to various resources.
type Cluster struct {
	// Name is the cluster name
	Name string
	// Dir is the directory where cluster certificates are stored
	Dir string
	// Status is the cluster status
	Status client.ProfileStatus
	// client is the cluster Teleport client
	client *client.TeleportClient
}

// NewCluster creates new cluster
func NewCluster(name, dir string) (*Cluster, error) {
	if name == "" {
		return nil, trace.BadParameter("cluster name is missing")
	}

	if dir == "" {
		return nil, trace.BadParameter("cluster directory is missing")
	}

	cfg := client.MakeDefaultConfig()
	cfg.WebProxyAddr = name
	cfg.HomePath = dir
	cfg.KeysDir = dir

	client, err := client.NewClient(cfg)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &Cluster{
		Dir:    dir,
		Name:   name,
		client: client,
	}, nil

}

// Connected indicates if connection to the cluster can be established
func (c *Cluster) Connected(clock clockwork.Clock) bool {
	return c.Status.IsExpired(clock)
}

// SyncAuthPreference fetches Teleport auth preferences and stores it in the cluster profile
func (c *Cluster) SyncAuthPreference(ctx context.Context) (*webclient.AuthenticationSettings, error) {
	pingResponse, err := c.client.Ping(ctx)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	if err := c.client.SaveProfile(c.Dir, false); err != nil {
		return nil, trace.Wrap(err)
	}

	return &pingResponse.Auth, nil
}

// GetRoles returns currently logged-in user roles
func (c *Cluster) GetRoles(ctx context.Context) ([]*types.Role, error) {
	proxyClient, err := c.client.ConnectToProxy(ctx)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	defer proxyClient.Close()

	roles := []*types.Role{}
	for _, name := range c.Status.Roles {
		role, err := proxyClient.GetRole(ctx, name)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

// GetUser returns currently logged-in user
func (c *Cluster) GetUser(ctx context.Context) (types.User, error) {
	proxyClient, err := c.client.ConnectToProxy(ctx)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	defer proxyClient.Close()

	user, err := proxyClient.GetUser(ctx, c.Status.Username)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return user, nil
}
