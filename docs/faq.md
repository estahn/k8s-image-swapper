# FAQ

### Is pulling from private registries supported?

Yes, `imagePullSecrets` on `Pod` and `ServiceAccount` level in the hooked pod definition are supported.

It is also possible to provide a list of ECRs to which authentication is handled by `k8s-image-swapper` using the same credentials as for the target registry. Please see [Configuration > Source - AWS](configuration.md#Private-registries).

### Are config changes reloaded gracefully?

Not yet, they require a pod rotation.

### What happens if the image is not found in the target registry?

Please see [Configuration > ImageCopyPolicy](configuration.md#imagecopypolicy).

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

### Why are sidecar images not being replaced?

A Kubernetes cluster can have multiple mutating webhooks.
Mutating webhooks execute sequentiatlly and each can change a submitted object.
Changes may be applied after `k8s-image-swapper` was executed, e.g. Istio injecting a sidecar.

```
... -> k8s-image-swapper -> Istio sidecar injection --> ...
```

Kubernetes 1.15+ allows to re-run webhooks if a mutating webhook modifies an object.
The behaviour is controlled by the [Reinvocation policy](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#reinvocation-policy).

> reinvocationPolicy may be set to `Never` or `IfNeeded`. It defaults to Never.
>
> * `Never`: the webhook must not be called more than once in a single admission evaluation
> * `IfNeeded`: the webhook may be called again as part of the admission evaluation if the object being admitted is modified by other admission plugins after the initial webhook call.

The reinvocation policy can be set in the helm chart as follows:

!!! example "Helm Chart"
    ```yaml
    webhook:
      reinvocationPolicy: IfNeeded
    ```
