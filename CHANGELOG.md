## [1.2.3](https://github.com/estahn/k8s-image-swapper/compare/v1.2.2...v1.2.3) (2022-09-01)

## [1.2.2](https://github.com/estahn/k8s-image-swapper/compare/v1.2.1...v1.2.2) (2022-08-01)

## [1.2.1](https://github.com/estahn/k8s-image-swapper/compare/v1.2.0...v1.2.1) (2022-07-26)

# [1.2.0](https://github.com/estahn/k8s-image-swapper/compare/v1.1.0...v1.2.0) (2022-07-03)


### Bug Fixes

* add missing dash ([228749d](https://github.com/estahn/k8s-image-swapper/commit/228749d98e32a7f90608b37b39d74a108f619f37))
* bump alpine to 3.16 due to security reports ([f7d6564](https://github.com/estahn/k8s-image-swapper/commit/f7d6564e1d607fa53a44e73f8b495a859c31aac1))
* docker references with both tag and digest ([5a17075](https://github.com/estahn/k8s-image-swapper/commit/5a170758a58b0244e6001a3aa5911c3be3d076f8)), closes [#48](https://github.com/estahn/k8s-image-swapper/issues/48)
* failed to solve: executor failed running ([af7df18](https://github.com/estahn/k8s-image-swapper/commit/af7df18a02d6455a4ff8ef1495741ad59cbb4856))
* setup buildx and qemu for image-scan ([c435048](https://github.com/estahn/k8s-image-swapper/commit/c43504873af1c5fd9c2551f8b77f3220f491ab6a))
* standard_init_linux.go:228: exec user process caused: exec format error ([b7d0c89](https://github.com/estahn/k8s-image-swapper/commit/b7d0c89d162ed0d71e01620cb074be68b8612ab2))
* **deps:** update module github.com/aws/aws-sdk-go to v1.40.54 ([7f9dbf5](https://github.com/estahn/k8s-image-swapper/commit/7f9dbf5cf5ddae16e252adc8ce21bb4039cd208d))


### Features

* add arm docker build ([be81815](https://github.com/estahn/k8s-image-swapper/commit/be8181590fb899f1515b78fbc02bf02986d72e9c))
* add full arm support to image copying ([6f14156](https://github.com/estahn/k8s-image-swapper/commit/6f14156acb610541d54d16e85171529de39af6ab))

# [1.1.0](https://github.com/estahn/k8s-image-swapper/compare/v1.0.0...v1.1.0) (2021-10-02)


### Bug Fixes

* provide log record for ImageSwapPolicyExists ([179da70](https://github.com/estahn/k8s-image-swapper/commit/179da706fd43c880d71063b786164f9d2cc862e4))
* timeout for ECR client ([26bdc10](https://github.com/estahn/k8s-image-swapper/commit/26bdc10c3eb21b1dfbea9a659e6b650cb25b335e))
* **deps:** update module github.com/alitto/pond to v1.5.1 ([504e2dd](https://github.com/estahn/k8s-image-swapper/commit/504e2dde58abf1312dab523cb43073a5cc7bc1b1))
* **deps:** update module github.com/aws/aws-sdk-go to v1.38.47 ([#70](https://github.com/estahn/k8s-image-swapper/issues/70)) ([4f30053](https://github.com/estahn/k8s-image-swapper/commit/4f300530ac9a6f8250672b272c24168601f42e62))
* **deps:** update module github.com/aws/aws-sdk-go to v1.40.43 ([266ef01](https://github.com/estahn/k8s-image-swapper/commit/266ef01da6d3caad97dac0f4d0a882dbd75502cc))
* **deps:** update module github.com/containers/image/v5 to v5.11.0 ([#61](https://github.com/estahn/k8s-image-swapper/issues/61)) ([11d6d28](https://github.com/estahn/k8s-image-swapper/commit/11d6d2843dbaa392a418e2a57fdab27fb5249077))
* **deps:** update module github.com/containers/image/v5 to v5.16.0 ([5230b91](https://github.com/estahn/k8s-image-swapper/commit/5230b91a7f37e0f4c6d6370d7c1a9231bf13b983))
* **deps:** update module github.com/dgraph-io/ristretto to v0.1.0 ([#82](https://github.com/estahn/k8s-image-swapper/issues/82)) ([dff1cb1](https://github.com/estahn/k8s-image-swapper/commit/dff1cb186ab1301836f978da1ead02b9ea75bb09))
* **deps:** update module github.com/go-co-op/gocron to v1.9.0 ([c0e9f11](https://github.com/estahn/k8s-image-swapper/commit/c0e9f111eb6b07d54732cc85464bab06dbfdf5e6))
* **deps:** update module github.com/rs/zerolog to v1.22.0 ([#76](https://github.com/estahn/k8s-image-swapper/issues/76)) ([c098326](https://github.com/estahn/k8s-image-swapper/commit/c098326273ab31dbd31869c4749164fde7544b67))
* **deps:** update module github.com/rs/zerolog to v1.23.0 ([#84](https://github.com/estahn/k8s-image-swapper/issues/84)) ([607d5bb](https://github.com/estahn/k8s-image-swapper/commit/607d5bb53a1d7396ae5d504ce49508ceac5e26d6))
* **deps:** update module github.com/rs/zerolog to v1.25.0 ([72822f4](https://github.com/estahn/k8s-image-swapper/commit/72822f42c762455a1a6932631e36418dc3b92d2a))
* **deps:** update module github.com/slok/kubewebhook to v2 ([8bd73d4](https://github.com/estahn/k8s-image-swapper/commit/8bd73d47772c0524c552577805d9f01ae365e77f))
* **deps:** update module github.com/spf13/cobra to v1.2.1 ([ea1e787](https://github.com/estahn/k8s-image-swapper/commit/ea1e7874cdaaa09dea34dd1d4a6f02a7ccb6925c))
* **deps:** update module github.com/spf13/viper to v1.8.1 ([8a055a2](https://github.com/estahn/k8s-image-swapper/commit/8a055a28343d8dbe780f74f99a275a311549576d))
* **deps:** update module k8s.io/api to v0.22.1 ([ab6d898](https://github.com/estahn/k8s-image-swapper/commit/ab6d898a2f9faa49b3c4f61f1443eb55bf79d93b))
* **deps:** update module k8s.io/apimachinery to v0.21.1 ([#79](https://github.com/estahn/k8s-image-swapper/issues/79)) ([aeeeffb](https://github.com/estahn/k8s-image-swapper/commit/aeeeffb4e20c50ecb0e3c0cb46654c3c41f62de0))
* **deps:** update module k8s.io/apimachinery to v0.22.2 ([ef72c66](https://github.com/estahn/k8s-image-swapper/commit/ef72c665f00d6d1fb454cd596c98b3a72cd7614c))


### Features

* Support for imagePullSecrets ([#112](https://github.com/estahn/k8s-image-swapper/issues/112)) ([2d8cf77](https://github.com/estahn/k8s-image-swapper/commit/2d8cf777d32053b8af622cb677d86ac21f526ba8)), closes [#92](https://github.com/estahn/k8s-image-swapper/issues/92) [#19](https://github.com/estahn/k8s-image-swapper/issues/19)
* Support for pod.spec.initContainers ([#118](https://github.com/estahn/k8s-image-swapper/issues/118)) ([725ff2c](https://github.com/estahn/k8s-image-swapper/commit/725ff2cdc45a13d1a31c3694231482ee09ab2cbd)), closes [#73](https://github.com/estahn/k8s-image-swapper/issues/73) [#96](https://github.com/estahn/k8s-image-swapper/issues/96)

# [1.1.0-alpha.1](https://github.com/estahn/k8s-image-swapper/compare/v1.0.0...v1.1.0-alpha.1) (2021-09-30)


### Bug Fixes

* provide log record for ImageSwapPolicyExists ([179da70](https://github.com/estahn/k8s-image-swapper/commit/179da706fd43c880d71063b786164f9d2cc862e4))
* timeout for ECR client ([26bdc10](https://github.com/estahn/k8s-image-swapper/commit/26bdc10c3eb21b1dfbea9a659e6b650cb25b335e))
* **deps:** update module github.com/alitto/pond to v1.5.1 ([504e2dd](https://github.com/estahn/k8s-image-swapper/commit/504e2dde58abf1312dab523cb43073a5cc7bc1b1))
* **deps:** update module github.com/aws/aws-sdk-go to v1.38.47 ([#70](https://github.com/estahn/k8s-image-swapper/issues/70)) ([4f30053](https://github.com/estahn/k8s-image-swapper/commit/4f300530ac9a6f8250672b272c24168601f42e62))
* **deps:** update module github.com/aws/aws-sdk-go to v1.40.43 ([266ef01](https://github.com/estahn/k8s-image-swapper/commit/266ef01da6d3caad97dac0f4d0a882dbd75502cc))
* **deps:** update module github.com/containers/image/v5 to v5.11.0 ([#61](https://github.com/estahn/k8s-image-swapper/issues/61)) ([11d6d28](https://github.com/estahn/k8s-image-swapper/commit/11d6d2843dbaa392a418e2a57fdab27fb5249077))
* **deps:** update module github.com/containers/image/v5 to v5.16.0 ([5230b91](https://github.com/estahn/k8s-image-swapper/commit/5230b91a7f37e0f4c6d6370d7c1a9231bf13b983))
* **deps:** update module github.com/dgraph-io/ristretto to v0.1.0 ([#82](https://github.com/estahn/k8s-image-swapper/issues/82)) ([dff1cb1](https://github.com/estahn/k8s-image-swapper/commit/dff1cb186ab1301836f978da1ead02b9ea75bb09))
* **deps:** update module github.com/go-co-op/gocron to v1.9.0 ([c0e9f11](https://github.com/estahn/k8s-image-swapper/commit/c0e9f111eb6b07d54732cc85464bab06dbfdf5e6))
* **deps:** update module github.com/rs/zerolog to v1.22.0 ([#76](https://github.com/estahn/k8s-image-swapper/issues/76)) ([c098326](https://github.com/estahn/k8s-image-swapper/commit/c098326273ab31dbd31869c4749164fde7544b67))
* **deps:** update module github.com/rs/zerolog to v1.23.0 ([#84](https://github.com/estahn/k8s-image-swapper/issues/84)) ([607d5bb](https://github.com/estahn/k8s-image-swapper/commit/607d5bb53a1d7396ae5d504ce49508ceac5e26d6))
* **deps:** update module github.com/rs/zerolog to v1.25.0 ([72822f4](https://github.com/estahn/k8s-image-swapper/commit/72822f42c762455a1a6932631e36418dc3b92d2a))
* **deps:** update module github.com/slok/kubewebhook to v2 ([8bd73d4](https://github.com/estahn/k8s-image-swapper/commit/8bd73d47772c0524c552577805d9f01ae365e77f))
* **deps:** update module github.com/spf13/cobra to v1.2.1 ([ea1e787](https://github.com/estahn/k8s-image-swapper/commit/ea1e7874cdaaa09dea34dd1d4a6f02a7ccb6925c))
* **deps:** update module github.com/spf13/viper to v1.8.1 ([8a055a2](https://github.com/estahn/k8s-image-swapper/commit/8a055a28343d8dbe780f74f99a275a311549576d))
* **deps:** update module k8s.io/api to v0.22.1 ([ab6d898](https://github.com/estahn/k8s-image-swapper/commit/ab6d898a2f9faa49b3c4f61f1443eb55bf79d93b))
* **deps:** update module k8s.io/apimachinery to v0.21.1 ([#79](https://github.com/estahn/k8s-image-swapper/issues/79)) ([aeeeffb](https://github.com/estahn/k8s-image-swapper/commit/aeeeffb4e20c50ecb0e3c0cb46654c3c41f62de0))
* **deps:** update module k8s.io/apimachinery to v0.22.2 ([ef72c66](https://github.com/estahn/k8s-image-swapper/commit/ef72c665f00d6d1fb454cd596c98b3a72cd7614c))


### Features

* Support for imagePullSecrets ([#112](https://github.com/estahn/k8s-image-swapper/issues/112)) ([2d8cf77](https://github.com/estahn/k8s-image-swapper/commit/2d8cf777d32053b8af622cb677d86ac21f526ba8)), closes [#92](https://github.com/estahn/k8s-image-swapper/issues/92) [#19](https://github.com/estahn/k8s-image-swapper/issues/19)

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
