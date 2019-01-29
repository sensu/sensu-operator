// Copyright 2016 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdutil

import "fmt"

// MemberConfig is used to create the template for the members in the StatefulSet
type MemberConfig struct {
	SecurePeer   bool
	SecureClient bool
	Namespace    string
}

func (m *MemberConfig) clientScheme() string {
	if m.SecureClient {
		return "https"
	}
	return "http"
}

// ListenClientURL is the URL to listen on for client connection
func (m *MemberConfig) ListenClientURL() string {
	return fmt.Sprintf("%s://0.0.0.0:2379", m.clientScheme())
}

// ListenPeerURL is the URL to listen on for peer connections
func (m *MemberConfig) ListenPeerURL() string {
	return fmt.Sprintf("%s://0.0.0.0:2380", m.peerScheme())
}

func (m *MemberConfig) peerScheme() string {
	if m.SecurePeer {
		return "https"
	}
	return "http"
}
