# :rotating_light: ALPHA, may contain :bug:, be warned! :rotating_light:

<p align="center">
  <img alt="Raiders of the Lost Ark" src="docs/img/indiana.gif" height="140" />
  <h3 align="center">k8s-image-swapper</h3>
  <p align="center">Mirror images into your own registry and swap image references automatically.</p>
</p>

---

`k8s-image-swapper` is a mutating webhook for Kubernetes, downloading images into your own registry and pointing the images to that new location.
It is an alternative to a [docker pull-through proxy](docker-mirror).
The feature set was primarily designed with Amazon ECR in mind but may work with other registries.

[docker-mirror]: https://docs.docker.com/registry/recipes/mirror/

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Why?](#why)
- [Getting started](#getting-started)
  - [Helm](#helm)
  - [Kustomize](#kustomize)
- [Stargazers over time](#stargazers-over-time)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Why?

* Consolidate all images in a single registry
* Protect against:
  * external registry failure ([quay.io outage](quay-outage))
  * image pull rate limiting ([docker.io rate limits](docker-rate-limiting))
  * accidental image changes
  * removal of images
* Use in air-gaped environments without the need to change manifests
* Reduce NAT ingress traffic/cost

[quay-outage]: https://www.reddit.com/r/devops/comments/f9kiej/quayio_is_experiencing_an_outage/
[docker-rate-limiting]: https://www.docker.com/blog/scaling-docker-to-serve-millions-more-developers-network-egress/

## Getting started

### Prerequisite

1. Create an IAM user with permissions to create ECR repositories and upload container images.
2. Create a Kubernetes secret (e.g. `k8s-image-swapper-aws`) with the IAM credentials
   
    ```bash
    kubectl create secret generic k8s-image-swapper-aws \
      --from-literal=aws_access_key_id=<...> \
      --from-literal=aws_secret_access_key=<...>
    ```

### Helm

```
helm repo add estahn https://estahn.github.io/charts/
helm install k8s-image-swapper estahn/k8s-image-swapper \
  --set image.tag=1.0.0-alpha.1 \
  --set config.target.registry.aws.accountId=$AWS_ACCOUNT_ID \
  --set config.target.registry.aws.region=$AWS_DEFAULT_REGION \
  --set awsSecretName=k8s-image-swapper-aws
```

## Badges

[![Release](https://img.shields.io/github/release/estahn/k8s-image-swapper.svg?style=for-the-badge)](https://github.com/estahn/k8s-image-swapper/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE.md)
[![Build status](https://img.shields.io/github/workflow/status/estahn/k8s-image-swapper/Test?style=for-the-badge)](https://github.com/estahn/k8s-image-swapper/actions?workflow=build)
[![Codecov branch](https://img.shields.io/codecov/c/github/estahn/k8s-image-swapper/main.svg?style=for-the-badge)](https://codecov.io/gh/estahn/k8s-image-swapper)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/estahn/k8s-image-swapper)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg?style=for-the-badge)](code_of_conduct.md)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/estahn/k8s-image-swapper.svg)](https://starchart.cc/estahn/k8s-image-swapper)
