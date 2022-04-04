# Priority Class

Pods can have priority. Priority indicates the importance of a Pod relative to other Pods. If a Pod cannot be scheduled, the scheduler tries to preempt (evict) lower priority Pods to make scheduling of the pending Pod possible.
A PriorityClass is a non-namespaced object that defines a mapping from a priority class name to the integer value of the priority.

Assuming you have a PriorityClass like:
```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: high-priority
value: 1000000
globalDefault: false
description: "This priority class should be used for XYZ service pods only."
```

you are able to reference it in the etcd cluster yaml:
```yaml
spec:
  size: 3
  pod:
    priorityClassName: high-priority
```