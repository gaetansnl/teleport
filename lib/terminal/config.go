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

package terminal

import (
	"os"

	"github.com/gravitational/trace"
)

// ServerOpts contains configuration options for a Teleport Terminal service.
type Config struct {
	// Addr is the bind address for the server, either in "scheme://host:port"
	// format (allowing for "tcp", "unix", etc) or in "host:port" format (defaults
	// to "tcp").
	Addr string `json:"addr"`
	// ShutdownSignals is the set of captured signals that cause server shutdown.
	ShutdownSignals []os.Signal `json:"-"`
	// HomeDir is the directory to store cluster profiles
	HomeDir string `json:"homeDir"`
}

// CheckAndSetDefaults checks and sets default config values.
func (c *Config) CheckAndSetDefaults() error {
	if c.Addr == "" {
		return trace.BadParameter("missing addr")
	}

	if c.HomeDir == "" {
		return trace.BadParameter("missing home directory")
	}

	return nil
}
