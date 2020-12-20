<p align="center">
  <img alt="Raiders of the Lost Ark" src="docs/img/indiana.gif" height="140" />
  <h3 align="center">k8s-image-swapper</h3>
  <p align="center">Mirror images into your own registry and swap image references automatically.</p>
</p>

---

`k8s-image-swapper` is a mutating webhook for Kubernetes, downloading images into your own registry and pointing the images to that new location.
It is an alternative to a [docker pull-through proxy](https://docs.docker.com/registry/recipes/mirror/).
The feature set was primarily designed with Amazon ECR in mind but may work with other registries.

## Benefits

Using `k8s-image-swapper` will improve the overall availability, reliability, durability and resiliency of your
Kubernetes cluster by keeping 3rd-party images mirrored into your own registry.

`k8s-image-swapper` will transparently consolidate all images into a single registry without the need to adjust manifests
therefore reducing the impact of external registry failures, rate limiting, network issues, change or removal of images
while reducing data traffic and therefore cost.

**TL;DR:**

* Protect against:
  * external registry failure ([quay.io outage](https://www.reddit.com/r/devops/comments/f9kiej/quayio_is_experiencing_an_outage/))
  * image pull rate limiting ([docker.io rate limits](https://www.docker.com/blog/scaling-docker-to-serve-millions-more-developers-network-egress/))
  * accidental image changes
  * removal of images
* Use in air-gaped environments without the need to change manifests
* Reduce NAT ingress traffic/cost

## Documentation

The documentation is available at [https://estahn.github.io/k8s-image-swapper/](https://estahn.github.io/k8s-image-swapper/index.html).

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
