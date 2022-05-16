#!/usr/bin/env bash

# Usage:
#   ./hack/release/bump_version.sh 0.8.0 0.8.1

oldv=$1
newv=$2

echo "old version: ${oldv}, new version: ${newv}"

sed -i -e "s/${oldv}+git/${newv}/g" version/version.go
sed -i -e "s/${oldv}.*/${newv}/g" example/deployment.yaml
sed -i -e "s/${oldv}.*/${newv}/g" example/deployment-ha.yaml
sed -i -e "s/${oldv}.*/${newv}/g" example/etcd-backup-operator/deployment.yaml
sed -i -e "s/${oldv}.*/${newv}/g" example/etcd-restore-operator/deployment.yaml

