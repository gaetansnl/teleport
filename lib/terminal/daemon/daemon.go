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

package daemon

import (
	"context"

	"github.com/gravitational/trace"

	"github.com/gravitational/teleport/api/profile"
	"github.com/gravitational/teleport/lib/client"
)

// Config is the cluster service config
type Config struct {
	// WorkingDir is the directory to store cluster profiles
	WorkingDir string
}

// Service is the cluster service
type Service struct {
	Config
	clusters []*Cluster
}

// Start creates and starts a Teleport Terminal service.
func New(cfg Config) (*Service, error) {
	return &Service{Config: cfg}, nil
}

// GetClusters returns a list of existing clusters
func (s *Service) GetClusters() []*Cluster {
	return s.clusters
}

// CreateCluster creates a new cluster
func (s *Service) CreateCluster(ctx context.Context, clusterName string) error {
	for _, cluster := range s.clusters {
		if cluster.Name == clusterName {
			return trace.BadParameter("Cluster %v already exists", clusterName)
		}
	}

	cluster, err := NewCluster(clusterName, s.WorkingDir)
	if err != nil {
		return trace.Wrap(err)
	}

	s.clusters = append(s.clusters, cluster)
	return nil
}

// GetCluster returns a cluster by its name
func (s *Service) GetCluster(name string) (*Cluster, error) {
	for _, cluster := range s.clusters {
		if cluster.Name == name {
			return cluster, nil
		}
	}

	return nil, trace.NotFound("cluster %v is not found", name)
}

// LoadClusters initializes existing clusters from their profiles
func (s *Service) LoadClusters() error {
	statuses, err := s.getStatuses()
	if err != nil {
		return trace.Wrap(err)
	}

	for _, sts := range statuses {
		cluster := &Cluster{
			Name:   sts.Name,
			Status: sts,
		}

		s.clusters = append(s.clusters, cluster)
	}

	return nil
}

func (s *Service) getStatuses() ([]client.ProfileStatus, error) {
	pNames, err := profile.ListProfileNames(s.WorkingDir)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	statuses := []client.ProfileStatus{}
	for _, name := range pNames {
		status, err := client.StatusFromFile(s.WorkingDir, name)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		statuses = append(statuses, *status)
	}

	return statuses, nil
}
