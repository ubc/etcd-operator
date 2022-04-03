# Roadmap

This document defines a high level roadmap for the etcd cluster operator development.



#### Stability/Reliability
- make examples deployment more secure by default
- make examples deployment more reliable by default (two replicas with topology aware scheduling)
- add priorityClassName support to etcd crd
- add crd validation schema
- add PodDisruptionBudget creation support