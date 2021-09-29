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

package handler

import (
	"context"

	"github.com/gravitational/teleport/lib/terminal/daemon"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"

	v1 "github.com/gravitational/teleport/lib/terminal/api/protogen/golang/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Config is the terminal service
type Config struct {
	// DaemonService is the instance of daemon service
	DaemonService *daemon.Service
}

// Handler implements teleport.terminal.v1.TerminalService.
type Handler struct {
	// Config is the service config
	Config
}

func New(cfg Config) (*Handler, error) {
	return &Handler{
		cfg,
	}, nil
}

// Lists all existing clusters
func (s *Handler) ListClusters(ctx context.Context, r *v1.ListClustersRequest) (*v1.ListClustersResponse, error) {
	result := []*v1.Cluster{}
	for _, cluster := range s.DaemonService.GetClusters() {
		proto := &v1.Cluster{
			Name:      cluster.Name,
			Connected: cluster.Connected(clockwork.NewRealClock()),
		}
		result = append(result, proto)
	}

	return &v1.ListClustersResponse{
		Clusters: result,
	}, nil
}

// CreateCluster creates a new cluster
func (s *Handler) CreateCluster(ctx context.Context, req *v1.CreateClusterRequest) (*v1.Cluster, error) {
	if err := s.DaemonService.CreateCluster(ctx, req.Name); err != nil {
		return nil, trace.Wrap(err)
	}

	cluster, err := s.DaemonService.GetCluster(req.Name)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	proto := &v1.Cluster{
		Name:      cluster.Name,
		Connected: cluster.Connected(clockwork.NewRealClock()),
	}

	return proto, nil
}

func (s *Handler) CreateClusterLoginChallenge(context.Context, *v1.CreateClusterLoginChallengeRequest) (*v1.ClusterLoginChallenge, error) {
	return nil, nil
}

func (s *Handler) SolveClusterLoginChallenge(context.Context, *v1.SolveClusterLoginChallengeRequest) (*v1.SolveClusterLoginChallengeResponse, error) {
	return nil, nil
}

func (s *Handler) ListDatabases(context.Context, *v1.ListDatabasesRequest) (*v1.ListDatabasesResponse, error) {
	return nil, nil
}

func (s *Handler) CreateGateway(context.Context, *v1.CreateGatewayRequest) (*v1.Gateway, error) {
	return nil, nil
}

func (s *Handler) ListGateways(context.Context, *v1.ListGatewaysRequest) (*v1.ListGatewaysResponse, error) {
	return nil, nil
}

func (s *Handler) DeleteGateway(context.Context, *v1.DeleteGatewayRequest) (*emptypb.Empty, error) {
	return nil, nil
}

// Streams input/output using a gateway.
// Requires the gateway to be created beforehand.
// This has no REST counterpart.
func (s *Handler) StreamGateway(v1.TerminalService_StreamGatewayServer) error {
	return nil
}

// GET /nodes
// Per Teleport nomenclature, a Node is an SSH-capable node.
// Requires login challenge to be solved beforehand.
func (s *Handler) ListNodes(context.Context, *v1.ListNodesRequest) (*v1.ListNodesResponse, error) {
	return nil, nil
}
