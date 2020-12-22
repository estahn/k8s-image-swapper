# FAQ

### Is pulling from private registries supported?

Not yet, but hopefully soon (see [GitHub #19](https://github.com/estahn/k8s-image-swapper/issues/19))

### Are config changes reloaded gracefully?

Not yet, they require a pod rotation.

### What happens if the image is not found in the target registry?

There is no notion of whether the image exists in the target registry.
`k8s-image-swapper` will blindly replace the images references.

This feature is planned for a future release (see [GitHub #12](https://github.com/estahn/k8s-image-swapper/issues/12)).

### What level of registry outage does this handle?

If the source image registry is not reachable it will replace the reference with the target registry reference.
If the target registry is down it will do the same. It has no notion of the target registry being up or down.

### What happens if `k8s-image-swapper` is unavailable?

Kubernetes will continue to work as if `k8s-image-swapper` was not installed.
The webhook [failure policy](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#failure-policy)
is set to `Ignore`.

!!! tip
    Environments with strict compliance requirements (or air-gapped) may overwrite this with `Fail` to
    avoid falling back to the public images.
