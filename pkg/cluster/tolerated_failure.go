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
	"math"

	"github.com/on2itsecurity/etcd-operator/pkg/util/k8sutil"
)

func (c *Cluster) updatePodDisruptionBudget(ctx context.Context) {
	currentlyMinAvailable := c.calculateMinAvailable(ctx)

	if c.status.MinAvailable != currentlyMinAvailable {
		c.logger.Infof("Maximum failed pods we can tolerate in the cluster has been changed to from %d to %d", c.status.MinAvailable, currentlyMinAvailable)
		c.status.MinAvailable = currentlyMinAvailable
		if c.cluster.Spec.PodDisruptionBudget {
			err := k8sutil.UpdateOrCreatePodDisruptionBudget(ctx, c.config.KubeCli, c.cluster.Namespace, c.cluster.Name, currentlyMinAvailable, c.cluster.AsOwner())
			if err != nil {
				c.logger.Errorf("failed to create/update pod disruption budget: %v", err)
			}
		}
	}
}

// calculateMinAvailable return the amount of pods that cannot be disrupted before we get out of quorum.
func (c *Cluster) calculateMinAvailable(ctx context.Context) int {
	targetSize := float64(c.cluster.Spec.Size)
	halfOfTarget := targetSize / 2
	quorum := math.Floor(halfOfTarget) + 1
	return int(quorum)
}
