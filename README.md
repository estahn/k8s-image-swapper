# k8s-image-swapper

Mirror images into your own registry and swap image references automatically.

A mutating webhook for Kubernetes, pointing the images to a new location.
It is an alternative to a [docker pull-through proxy](docker-mirror).
The feature set was primarily designed with Amazon ECR in mind but may work with other registries.

[docker-mirror]: https://docs.docker.com/registry/recipes/mirror/

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Why?](#why)
- [Architecture](#architecture)
- [Other](#other)
- [POC](#poc)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Why?

* Consolidate all images in a single registry
* Protection against:
  * external registry failure ([quay.io outage](quay-outage))
  * image pull rate limiting ([docker.io rate limits](docker-rate-limiting))
  * accidental image changes
  * removal of images
* Use in air-gaped environments without the need to change manifests
* Reduce NAT ingress traffic/cost

[quay-outage]: https://www.reddit.com/r/devops/comments/f9kiej/quayio_is_experiencing_an_outage/
[docker-rate-limiting]: https://www.docker.com/blog/scaling-docker-to-serve-millions-more-developers-network-egress/

## Architecture

![alt text](architecture.jpg "k8s-image-swapper Architecture")

Components:
* **Image Swapper**, a mutating Webhook adding a prefix to the original image string
* **Image Downloader**
* **Image State Keeper**
* Repository Manager
  * Create repository
  * Delete repository
  * Sync image

## Other

* Manages ECR repository life-cycles (create&delete repository, sync images)

Options:

- registry
- whitelist
- blacklist
- ecr_lifecycle
    - create
    - delete


Components:
* Mutating Webhook
  * Adds a prefix to the original image string
* Repository Manager
  * Create repository
  * Delete repository
  * Sync image

nginx:latest -> docker.io/nginx:latest -> <registry>/docker.io/nginx:latest


## POC

1. Receive admission request
2. Check if image is in ECR
3. If not download image and push to ECR, return false
4. return with new image reference
