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

package cluster

import (
	"context"
	"testing"

	"github.com/on2itsecurity/etcd-operator/pkg/apis/etcd/v1beta2"
)

func TestCalculateMinAvailable(t *testing.T) {
	c := &Cluster{
		cluster: &v1beta2.EtcdCluster{
			Spec: v1beta2.ClusterSpec{
				Size: 7,
			},
		},
	}
	minAvailable := c.calculateMinAvailable(context.Background())
	expected := 4
	if minAvailable != expected {
		t.Errorf("expect minAvailable=%d, get=%d", expected, minAvailable)
	}

	c = &Cluster{
		cluster: &v1beta2.EtcdCluster{
			Spec: v1beta2.ClusterSpec{
				Size: 3,
			},
		},
	}
	minAvailable = c.calculateMinAvailable(context.Background())
	expected = 2
	if minAvailable != expected {
		t.Errorf("expect minAvailable=%d, get=%d", expected, minAvailable)
	}
	c = &Cluster{
		cluster: &v1beta2.EtcdCluster{
			Spec: v1beta2.ClusterSpec{
				Size: 1,
			},
		},
	}
	minAvailable = c.calculateMinAvailable(context.Background())
	expected = 1
	if minAvailable != expected {
		t.Errorf("expect minAvailable=%d, get=%d", expected, minAvailable)
	}
}
