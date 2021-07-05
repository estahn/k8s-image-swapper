## [1.0.1](https://github.com/DolevAlgam/k8s-image-swapper/compare/v1.0.0...v1.0.1) (2021-07-05)


### Bug Fixes

* **deps:** update module github.com/aws/aws-sdk-go to v1.38.47 ([#70](https://github.com/DolevAlgam/k8s-image-swapper/issues/70)) ([4f30053](https://github.com/DolevAlgam/k8s-image-swapper/commit/4f300530ac9a6f8250672b272c24168601f42e62))
* **deps:** update module github.com/containers/image/v5 to v5.11.0 ([#61](https://github.com/DolevAlgam/k8s-image-swapper/issues/61)) ([11d6d28](https://github.com/DolevAlgam/k8s-image-swapper/commit/11d6d2843dbaa392a418e2a57fdab27fb5249077))
* **deps:** update module github.com/dgraph-io/ristretto to v0.1.0 ([#82](https://github.com/DolevAlgam/k8s-image-swapper/issues/82)) ([dff1cb1](https://github.com/DolevAlgam/k8s-image-swapper/commit/dff1cb186ab1301836f978da1ead02b9ea75bb09))
* **deps:** update module github.com/rs/zerolog to v1.22.0 ([#76](https://github.com/DolevAlgam/k8s-image-swapper/issues/76)) ([c098326](https://github.com/DolevAlgam/k8s-image-swapper/commit/c098326273ab31dbd31869c4749164fde7544b67))
* **deps:** update module github.com/rs/zerolog to v1.23.0 ([#84](https://github.com/DolevAlgam/k8s-image-swapper/issues/84)) ([607d5bb](https://github.com/DolevAlgam/k8s-image-swapper/commit/607d5bb53a1d7396ae5d504ce49508ceac5e26d6))
* **deps:** update module k8s.io/apimachinery to v0.21.1 ([#79](https://github.com/DolevAlgam/k8s-image-swapper/issues/79)) ([aeeeffb](https://github.com/DolevAlgam/k8s-image-swapper/commit/aeeeffb4e20c50ecb0e3c0cb46654c3c41f62de0))

# 1.0.0 (2020-12-25)


### Bug Fixes

* bump skopeo from 0.2.0 to 1.2.0 ([84025aa](https://github.com/estahn/k8s-image-swapper/commit/84025aaf06d287a306fba98f848e272a19ff8aa0))
* hardcoded AWS region ([3cc0d49](https://github.com/estahn/k8s-image-swapper/commit/3cc0d492bc17a6ad022cb2794786079759f7bc41)), closes [#20](https://github.com/estahn/k8s-image-swapper/issues/20) [#17](https://github.com/estahn/k8s-image-swapper/issues/17)
* **chart:** serviceaccount missing annotation tag ([#21](https://github.com/estahn/k8s-image-swapper/issues/21)) ([7164626](https://github.com/estahn/k8s-image-swapper/commit/71646266e54d043f3bba2ee59975e7f9d11f8f13))
* trace for verbose logs and improve context ([58e05dc](https://github.com/estahn/k8s-image-swapper/commit/58e05dc66644de22183e39dcdc85cf8ce139d8db)), closes [#15](https://github.com/estahn/k8s-image-swapper/issues/15)


### Features

* allow filters for container context ([37d0a4d](https://github.com/estahn/k8s-image-swapper/commit/37d0a4d9ac3bd37128c92ede0bff3f4071483b1d)), closes [#32](https://github.com/estahn/k8s-image-swapper/issues/32)
* automatic token renewal before expiry ([a7c45b8](https://github.com/estahn/k8s-image-swapper/commit/a7c45b8b093efa00e7a04f89a57d5909b4ce068a)), closes [#31](https://github.com/estahn/k8s-image-swapper/issues/31)
* helm chart ([00f6b74](https://github.com/estahn/k8s-image-swapper/commit/00f6b7409c1f0ab59ea227f5d3b995d532beb623))
* ImageSwapPolicy defines the mutation strategy used by the webhook. ([9d61659](https://github.com/estahn/k8s-image-swapper/commit/9d616596013d7b1cbb121b0cf137273867bdb19f))
* POC ([fedcb22](https://github.com/estahn/k8s-image-swapper/commit/fedcb22c2fef26a76bd0fd9dacff70d0d952c077))

# [1.0.0-beta.4](https://github.com/estahn/k8s-image-swapper/compare/v1.0.0-beta.3...v1.0.0-beta.4) (2020-12-23)


### Bug Fixes

* bump skopeo from 0.2.0 to 1.2.0 ([09fdb6e](https://github.com/estahn/k8s-image-swapper/commit/09fdb6eb2383c30a45d1a5a7fb3d10a4c6b891e0))

# [1.0.0-beta.3](https://github.com/estahn/k8s-image-swapper/compare/v1.0.0-beta.2...v1.0.0-beta.3) (2020-12-23)


### Features

* ImageSwapPolicy defines the mutation strategy used by the webhook. ([e64bc6d](https://github.com/estahn/k8s-image-swapper/commit/e64bc6d120bea925a06cf06f3b22c8184a24fb35))

# [1.0.0-beta.2](https://github.com/estahn/k8s-image-swapper/compare/v1.0.0-beta.1...v1.0.0-beta.2) (2020-12-22)


### Features

* allow filters for container context ([c7e4c51](https://github.com/estahn/k8s-image-swapper/commit/c7e4c51a5a04ef9ae8689ffe73ff7d1411f43450)), closes [#32](https://github.com/estahn/k8s-image-swapper/issues/32)
* automatic token renewal before expiry ([d557c23](https://github.com/estahn/k8s-image-swapper/commit/d557c23e798f4cae61cd412d99f482ec4d310b9f)), closes [#31](https://github.com/estahn/k8s-image-swapper/issues/31)

# 1.0.0-beta.1 (2020-12-21)


### Bug Fixes

* hardcoded AWS region ([3cc0d49](https://github.com/estahn/k8s-image-swapper/commit/3cc0d492bc17a6ad022cb2794786079759f7bc41)), closes [#20](https://github.com/estahn/k8s-image-swapper/issues/20) [#17](https://github.com/estahn/k8s-image-swapper/issues/17)
* **chart:** serviceaccount missing annotation tag ([#21](https://github.com/estahn/k8s-image-swapper/issues/21)) ([7164626](https://github.com/estahn/k8s-image-swapper/commit/71646266e54d043f3bba2ee59975e7f9d11f8f13))
* trace for verbose logs and improve context ([58e05dc](https://github.com/estahn/k8s-image-swapper/commit/58e05dc66644de22183e39dcdc85cf8ce139d8db)), closes [#15](https://github.com/estahn/k8s-image-swapper/issues/15)


### Features

* helm chart ([00f6b74](https://github.com/estahn/k8s-image-swapper/commit/00f6b7409c1f0ab59ea227f5d3b995d532beb623))
* POC ([fedcb22](https://github.com/estahn/k8s-image-swapper/commit/fedcb22c2fef26a76bd0fd9dacff70d0d952c077))

# 1.0.0-alpha.1 (2020-12-18)


### Bug Fixes

* **chart:** serviceaccount missing annotation tag ([#21](https://github.com/estahn/k8s-image-swapper/issues/21)) ([7164626](https://github.com/estahn/k8s-image-swapper/commit/71646266e54d043f3bba2ee59975e7f9d11f8f13))
* trace for verbose logs and improve context ([58e05dc](https://github.com/estahn/k8s-image-swapper/commit/58e05dc66644de22183e39dcdc85cf8ce139d8db)), closes [#15](https://github.com/estahn/k8s-image-swapper/issues/15)


### Features

* helm chart ([00f6b74](https://github.com/estahn/k8s-image-swapper/commit/00f6b7409c1f0ab59ea227f5d3b995d532beb623))
* POC ([fedcb22](https://github.com/estahn/k8s-image-swapper/commit/fedcb22c2fef26a76bd0fd9dacff70d0d952c077))
