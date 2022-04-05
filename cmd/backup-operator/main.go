// Copyright 2017 The etcd-operator Authors
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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	_ "github.com/emadolsky/automaxprocs/maxprocs"
	controller "github.com/on2itsecurity/etcd-operator/pkg/controller/backup-operator"
	"github.com/on2itsecurity/etcd-operator/pkg/util/constants"
	"github.com/on2itsecurity/etcd-operator/pkg/util/k8sutil"
	"github.com/on2itsecurity/etcd-operator/version"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var (
	createCRD   bool
	namespace   string
	clusterWide bool
)

func init() {
	flag.BoolVar(&createCRD, "create-crd", true, "The backup operator will not create the EtcdBackup CRD when this flag is set to false.")
	flag.BoolVar(&clusterWide, "cluster-wide", false, "Enable operator to watch clusters in all namespaces")
	flag.Parse()
}

func main() {
	namespace := os.Getenv(constants.EnvOperatorPodNamespace)
	if len(namespace) == 0 {
		logrus.Fatalf("must set env %s", constants.EnvOperatorPodNamespace)
	}
	name := os.Getenv(constants.EnvOperatorPodName)
	if len(name) == 0 {
		logrus.Fatalf("must set env %s", constants.EnvOperatorPodName)
	}
	id, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("failed to get hostname: %v", err)
	}

	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("etcd-backup-operator Version: %v", version.Version)
	logrus.Infof("Git SHA: %s", version.GitSHA)

	kubecli := k8sutil.MustNewKubeClient()

	rl := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "etcd-backup-operator",
			Namespace: namespace,
		},
		Client: kubecli.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", 9091), nil)
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				logrus.Fatalf("leader election lost")
			},
		},
	})
}

func run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := controller.New(newControllerConfig())
	err := c.Start(ctx)
	if err != nil {
		logrus.Fatalf("operator stopped with error: %v", err)
	}
}

func newControllerConfig() controller.Config {
	cfg := controller.Config{
		Namespace:   namespace,
		ClusterWide: clusterWide,
		CreateCRD:   createCRD,
	}

	return cfg
}
