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

package apiserver

import (
	"net"
	"strings"

	"github.com/gravitational/teleport/lib/terminal/apiserver/handler"
	"github.com/gravitational/teleport/lib/terminal/daemon"
	"github.com/gravitational/trace"
	"google.golang.org/grpc"

	terminalpb "github.com/gravitational/teleport/lib/terminal/api/protogen/golang/v1"
)

// Config is the APIServer configuration
type Config struct {
	// HostAddr is the APIServer host address
	HostAddr string
	// Daemon is the terminal daemon service
	Daemon *daemon.Service
}

// CheckAndSetDefaults checks and sets default config values.
func (c *Config) CheckAndSetDefaults() error {
	if c.HostAddr == "" {
		return trace.BadParameter("missing HostAddr")
	}

	if c.Daemon == nil {
		return trace.BadParameter("missing daemon service")
	}

	return nil
}

// Server is a combination of the underlying grpc.Server and its RuntimeOpts.
type APIServer struct {
	Config
	// ls is the server listener
	ls net.Listener
	// grpc is an instance of grpc server
	grpcServer *grpc.Server
}

func New(cfg Config) (apiServer *APIServer, err error) {
	if err := cfg.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	ls, err := newListener(cfg.HostAddr)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	defer func() {
		if err != nil {
			ls.Close()
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.Creds(nil),
	)

	serviceHandler, err := handler.New(
		handler.Config{
			DaemonService: cfg.Daemon,
		},
	)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	terminalpb.RegisterTerminalServiceServer(grpcServer, serviceHandler)

	return &APIServer{cfg, ls, grpcServer}, nil
}

// ServeAndWait starts the server goroutine and waits until it exits.
func (s *APIServer) Serve() error {
	return s.grpcServer.Serve(s.ls)
}

// Close terminates the server and closes all open connections
func (s *APIServer) Stop() {
	s.grpcServer.GracefulStop()
}

func newListener(host string) (net.Listener, error) {
	var network, addr string

	parts := strings.SplitN(host, "://", 2)
	network = parts[0]
	switch network {
	case "unix":
		addr = parts[1]
	default:
		return nil, trace.BadParameter("Unsupported network protocol: %s", network)
	}

	lis, err := net.Listen(network, addr)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return lis, nil
}
