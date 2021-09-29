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

package terminal

import (
	apidefaults "github.com/gravitational/teleport/api/defaults"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/client"
	"github.com/gravitational/teleport/lib/services"
)

type access struct {
	List   bool `json:"list"`
	Read   bool `json:"read"`
	Edit   bool `json:"edit"`
	Create bool `json:"create"`
	Delete bool `json:"remove"`
}

type ACL struct {
	// Sessions defines access to recorded sessions
	Sessions access
	// AuthConnectors defines access to auth.connectors
	AuthConnectors access
	// Roles defines access to roles
	Roles access
	// Users defines access to users.
	Users access
	// TrustedClusters defines access to trusted clusters
	TrustedClusters access
	// Events defines access to audit logs
	Events access
	// Tokens defines access to tokens.
	Tokens access
	// Nodes defines access to nodes.
	Nodes access
	// AppServers defines access to application servers
	AppServers access
	// DBServers defines access to database servers.
	DBServers access
	// KubeServers defines access to kubernetes servers.
	KubeServers access
	// SSH defines access to servers
	SSHLogins []string
	// AccessRequests defines access to access requests
	AccessRequests access
	// Billing defines access to billing information
	Billing access
}

type Permissions struct {
	// ACL is the access control list of the logged-in user
	ACL ACL
}

// Cluster describes user settings and access to various resources.
type ClusterStatus struct {
	client.ProfileStatus
}

func (cs *ClusterStatus) GetACL() ACL {

	//cs.

	return ACL{}
}

func hasAccess(user types.User, roleSet services.RoleSet, kind string, verbs ...string) bool {
	srvCtx := &services.Context{User: user}
	for _, verb := range verbs {
		// Since this check occurs often and it does not imply the caller is trying
		// to access any resource, silence any logging done on the proxy.
		err := roleSet.CheckAccessToRule(srvCtx, apidefaults.Namespace, kind, verb, true)
		if err != nil {
			return false
		}
	}

	return true
}
