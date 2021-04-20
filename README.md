# Label Kubernetes Pods

This binary will read kubernetes configuration by checking default locations for `kubectl` and then
apply the `deploy=groupN` label to all pods that feature the `app=sfu` label.

There are four groups currently: 

* `group1` will be applied to ~10% of pods
* `group2` will be applied to ~30% of pods
* `group3` will be applied to ~30% of pods
* `group4` will be applied to the rest of pods (~30%)

## Rationale

We are using a statefulSet deployment in combination with [OnDelete update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#on-delete).

To delete the deployment in batches we are applying labels, splitting pods in groups.
