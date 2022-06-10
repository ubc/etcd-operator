#!/usr/bin/env bash

GOBIN="$(go env GOBIN)"
gobin="${GOBIN:-$(go env GOPATH)/bin}"

  "${gobin}/controller-gen" crd paths=./pkg/apis/etcd/v1beta2/... output:crd:artifacts:config=./example/crd
