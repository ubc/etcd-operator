# Roadmap

This document defines a high level roadmap for the etcd cluster operator development.



#### Stability/Reliability
- make release container images smaller using multistage
- make examples deployment more secure by default
- make examples deployment more reliable by default (two replicas with topology aware scheduling)
- add priorityClassName support to etcd crd
- add crd validation schema