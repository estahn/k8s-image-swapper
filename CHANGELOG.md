## [1.3.1](https://github.com/estahn/k8s-image-swapper/compare/v1.3.0...v1.3.1) (2022-10-01)


### :bug: Bug Fixes

* set verbose level & use structured logging ([#346](https://github.com/estahn/k8s-image-swapper/issues/346)) ([9b21320](https://github.com/estahn/k8s-image-swapper/commit/9b21320a52d3f74ae4a6e8233cc3e310d2f5136b))


### :arrow_up: Dependencies

* **deps:** bump github.com/aws/aws-sdk-go from 1.44.92 to 1.44.95 ([#349](https://github.com/estahn/k8s-image-swapper/issues/349)) ([609e915](https://github.com/estahn/k8s-image-swapper/commit/609e91566628b2c89ee0f3a6f582993cb7df8154)), closes [#4553](https://github.com/estahn/k8s-image-swapper/issues/4553) [#4551](https://github.com/estahn/k8s-image-swapper/issues/4551) [#4550](https://github.com/estahn/k8s-image-swapper/issues/4550)
* **deps:** bump github.com/aws/aws-sdk-go from 1.44.95 to 1.44.100 ([#351](https://github.com/estahn/k8s-image-swapper/issues/351)) ([c4aba7d](https://github.com/estahn/k8s-image-swapper/commit/c4aba7dd91b6128c4b6b70b52a3587d81a1b439f)), closes [#4560](https://github.com/estahn/k8s-image-swapper/issues/4560) [#4559](https://github.com/estahn/k8s-image-swapper/issues/4559) [#4558](https://github.com/estahn/k8s-image-swapper/issues/4558) [#4556](https://github.com/estahn/k8s-image-swapper/issues/4556) [#4555](https://github.com/estahn/k8s-image-swapper/issues/4555)
* **deps:** bump github.com/gruntwork-io/terratest from 0.40.21 to 0.40.22 ([#348](https://github.com/estahn/k8s-image-swapper/issues/348)) ([b3fa94d](https://github.com/estahn/k8s-image-swapper/commit/b3fa94df956a05796d8fd396462d0bb6987c8f11)), closes [#1169](https://github.com/estahn/k8s-image-swapper/issues/1169)
* **deps:** bump k8s.io/api from 0.25.0 to 0.25.1 ([#350](https://github.com/estahn/k8s-image-swapper/issues/350)) ([e1b358a](https://github.com/estahn/k8s-image-swapper/commit/e1b358aa28abacbf4e2c125032871d9db6fab401)), closes [#112161](https://github.com/estahn/k8s-image-swapper/issues/112161) [pohly/automated-cherry-pick-of-#112129](https://github.com/pohly/automated-cherry-pick-of-/issues/112129)
* **deps:** bump k8s.io/apimachinery from 0.25.0 to 0.25.1 ([#352](https://github.com/estahn/k8s-image-swapper/issues/352)) ([046ad1e](https://github.com/estahn/k8s-image-swapper/commit/046ad1e07924a4b4e797e5984262ef09872e5e50)), closes [#112330](https://github.com/estahn/k8s-image-swapper/issues/112330) [enj/automated-cherry-pick-of-#112193](https://github.com/enj/automated-cherry-pick-of-/issues/112193) [#112161](https://github.com/estahn/k8s-image-swapper/issues/112161) [pohly/automated-cherry-pick-of-#112129](https://github.com/pohly/automated-cherry-pick-of-/issues/112129)
* **deps:** bump k8s.io/client-go from 0.25.0 to 0.25.1 ([#353](https://github.com/estahn/k8s-image-swapper/issues/353)) ([4525ad4](https://github.com/estahn/k8s-image-swapper/commit/4525ad4a667fda8e86d4c19d3c57f9d2fe9ab7c3)), closes [#112161](https://github.com/estahn/k8s-image-swapper/issues/112161) [pohly/automated-cherry-pick-of-#112129](https://github.com/pohly/automated-cherry-pick-of-/issues/112129) [#112336](https://github.com/estahn/k8s-image-swapper/issues/112336) [enj/automated-cherry-pick-of-#112017](https://github.com/enj/automated-cherry-pick-of-/issues/112017) [#112055](https://github.com/estahn/k8s-image-swapper/issues/112055) [aanm/automated-cherry-pick-of-#111752](https://github.com/aanm/automated-cherry-pick-of-/issues/111752)

## [1.3.0](https://github.com/estahn/k8s-image-swapper/compare/v1.2.3...v1.3.0) (2022-09-07)


### :tada: Features

* cross account caching with role ([#336](https://github.com/estahn/k8s-image-swapper/issues/336)) ([98d138e](https://github.com/estahn/k8s-image-swapper/commit/98d138ece9dc27acf20266994e25bef4d43c3d7b))


### :arrow_up: Dependencies

* **deps:** bump actions/cache from 3.0.6 to 3.0.8 ([#319](https://github.com/estahn/k8s-image-swapper/issues/319)) ([245ab30](https://github.com/estahn/k8s-image-swapper/commit/245ab30bec7155caaad2ee95689ca71574f69252)), closes [#809](https://github.com/estahn/k8s-image-swapper/issues/809) [#833](https://github.com/estahn/k8s-image-swapper/issues/833) [#810](https://github.com/estahn/k8s-image-swapper/issues/810) [#888](https://github.com/estahn/k8s-image-swapper/issues/888) [#891](https://github.com/estahn/k8s-image-swapper/issues/891) [#899](https://github.com/estahn/k8s-image-swapper/issues/899) [#894](https://github.com/estahn/k8s-image-swapper/issues/894)
* **deps:** bump alpine from 3.16.1 to 3.16.2 ([da05fdd](https://github.com/estahn/k8s-image-swapper/commit/da05fdd19e9b2540a1a57b30aadabd00ea260f9e))
* **deps:** bump github.com/alitto/pond from 1.8.0 to 1.8.1 ([#342](https://github.com/estahn/k8s-image-swapper/issues/342)) ([4e50c28](https://github.com/estahn/k8s-image-swapper/commit/4e50c28818fb7db5f2d9b3431a346036109a8f44)), closes [alitto/pond#33](https://github.com/alitto/pond/issues/33) [#34](https://github.com/estahn/k8s-image-swapper/issues/34) [#32](https://github.com/estahn/k8s-image-swapper/issues/32)
* **deps:** bump github.com/aws/aws-sdk-go from 1.44.70 to 1.44.92 ([0f396c5](https://github.com/estahn/k8s-image-swapper/commit/0f396c57a16e97a5ed01dd310cd7fe808cb0c8b1))
* **deps:** bump github.com/aws/aws-sdk-go from 1.44.70 to 1.44.92 ([#338](https://github.com/estahn/k8s-image-swapper/issues/338)) ([fa795ae](https://github.com/estahn/k8s-image-swapper/commit/fa795aef3e847fb0f1526dca9efc6cd44ddd9fd9)), closes [#4548](https://github.com/estahn/k8s-image-swapper/issues/4548) [#4546](https://github.com/estahn/k8s-image-swapper/issues/4546) [#4545](https://github.com/estahn/k8s-image-swapper/issues/4545) [#4544](https://github.com/estahn/k8s-image-swapper/issues/4544) [#4543](https://github.com/estahn/k8s-image-swapper/issues/4543) [#4542](https://github.com/estahn/k8s-image-swapper/issues/4542) [#4539](https://github.com/estahn/k8s-image-swapper/issues/4539) [#4536](https://github.com/estahn/k8s-image-swapper/issues/4536) [#4534](https://github.com/estahn/k8s-image-swapper/issues/4534) [#4533](https://github.com/estahn/k8s-image-swapper/issues/4533)
* **deps:** bump github.com/go-co-op/gocron from 1.16.2 to 1.17.0 ([#340](https://github.com/estahn/k8s-image-swapper/issues/340)) ([645bef3](https://github.com/estahn/k8s-image-swapper/commit/645bef3b6b2ab2c936b0192dd24fd083f64e2034)), closes [go-co-op/gocron#380](https://github.com/go-co-op/gocron/issues/380) [go-co-op/gocron#381](https://github.com/go-co-op/gocron/issues/381) [go-co-op/gocron#375](https://github.com/go-co-op/gocron/issues/375) [#381](https://github.com/estahn/k8s-image-swapper/issues/381) [#380](https://github.com/estahn/k8s-image-swapper/issues/380) [#375](https://github.com/estahn/k8s-image-swapper/issues/375)
* **deps:** bump github.com/gruntwork-io/terratest from 0.40.19 to 0.40.21 ([#334](https://github.com/estahn/k8s-image-swapper/issues/334)) ([d0f6c39](https://github.com/estahn/k8s-image-swapper/commit/d0f6c39c30c6c47c502b036de3687c73912ecec9)), closes [#1166](https://github.com/estahn/k8s-image-swapper/issues/1166) [#1159](https://github.com/estahn/k8s-image-swapper/issues/1159)
* **deps:** bump github.com/rs/zerolog from 1.27.0 to 1.28.0 ([#339](https://github.com/estahn/k8s-image-swapper/issues/339)) ([7fb4ff5](https://github.com/estahn/k8s-image-swapper/commit/7fb4ff588ca7f0d177cc9f5bb36066367f9ca84d)), closes [#457](https://github.com/estahn/k8s-image-swapper/issues/457) [#416](https://github.com/estahn/k8s-image-swapper/issues/416) [#454](https://github.com/estahn/k8s-image-swapper/issues/454) [#453](https://github.com/estahn/k8s-image-swapper/issues/453) [#383](https://github.com/estahn/k8s-image-swapper/issues/383) [#396](https://github.com/estahn/k8s-image-swapper/issues/396) [#414](https://github.com/estahn/k8s-image-swapper/issues/414) [#415](https://github.com/estahn/k8s-image-swapper/issues/415) [#430](https://github.com/estahn/k8s-image-swapper/issues/430) [#432](https://github.com/estahn/k8s-image-swapper/issues/432)
* **deps:** bump github.com/spf13/viper from 1.12.0 to 1.13.0 ([#341](https://github.com/estahn/k8s-image-swapper/issues/341)) ([9b59bd4](https://github.com/estahn/k8s-image-swapper/commit/9b59bd4f308916d207fcfb5c7f3c70eedda1c615)), closes [spf13/viper#1371](https://github.com/spf13/viper/issues/1371) [spf13/viper#1373](https://github.com/spf13/viper/issues/1373) [spf13/viper#1393](https://github.com/spf13/viper/issues/1393) [spf13/viper#1424](https://github.com/spf13/viper/issues/1424) [spf13/viper#1405](https://github.com/spf13/viper/issues/1405) [spf13/viper#1414](https://github.com/spf13/viper/issues/1414) [spf13/viper#1387](https://github.com/spf13/viper/issues/1387) [spf13/viper#1374](https://github.com/spf13/viper/issues/1374) [spf13/viper#1375](https://github.com/spf13/viper/issues/1375) [spf13/viper#1378](https://github.com/spf13/viper/issues/1378) [spf13/viper#1360](https://github.com/spf13/viper/issues/1360) [spf13/viper#1381](https://github.com/spf13/viper/issues/1381) [spf13/viper#1384](https://github.com/spf13/viper/issues/1384) [spf13/viper#1383](https://github.com/spf13/viper/issues/1383) [spf13/viper#1395](https://github.com/spf13/viper/issues/1395) [spf13/viper#1420](https://github.com/spf13/viper/issues/1420) [spf13/viper#1422](https://github.com/spf13/viper/issues/1422) [spf13/viper#1412](https://github.com/spf13/viper/issues/1412) [spf13/viper#1373](https://github.com/spf13/viper/issues/1373) [spf13/viper#1393](https://github.com/spf13/viper/issues/1393) [spf13/viper#1371](https://github.com/spf13/viper/issues/1371) [spf13/viper#1387](https://github.com/spf13/viper/issues/1387) [spf13/viper#1405](https://github.com/spf13/viper/issues/1405) [spf13/viper#1414](https://github.com/spf13/viper/issues/1414)
* **deps:** bump goreleaser/goreleaser-action from 3.0.0 to 3.1.0 ([#328](https://github.com/estahn/k8s-image-swapper/issues/328)) ([a8d2dd1](https://github.com/estahn/k8s-image-swapper/commit/a8d2dd1916be3b7e686cb2e6814710ab73c5f953)), closes [#369](https://github.com/estahn/k8s-image-swapper/issues/369) [#357](https://github.com/estahn/k8s-image-swapper/issues/357) [#356](https://github.com/estahn/k8s-image-swapper/issues/356) [#360](https://github.com/estahn/k8s-image-swapper/issues/360) [#359](https://github.com/estahn/k8s-image-swapper/issues/359) [#358](https://github.com/estahn/k8s-image-swapper/issues/358) [#367](https://github.com/estahn/k8s-image-swapper/issues/367) [#369](https://github.com/estahn/k8s-image-swapper/issues/369) [#367](https://github.com/estahn/k8s-image-swapper/issues/367) [#358](https://github.com/estahn/k8s-image-swapper/issues/358) [#359](https://github.com/estahn/k8s-image-swapper/issues/359) [#360](https://github.com/estahn/k8s-image-swapper/issues/360) [#357](https://github.com/estahn/k8s-image-swapper/issues/357) [#356](https://github.com/estahn/k8s-image-swapper/issues/356)
* **deps:** bump k8s.io/api from 0.24.3 to 0.25.0 ([#325](https://github.com/estahn/k8s-image-swapper/issues/325)) ([ce10907](https://github.com/estahn/k8s-image-swapper/commit/ce10907f31431c641269032b823beaff4932f224)), closes [#111657](https://github.com/estahn/k8s-image-swapper/issues/111657) [#109090](https://github.com/estahn/k8s-image-swapper/issues/109090) [#111258](https://github.com/estahn/k8s-image-swapper/issues/111258) [#111113](https://github.com/estahn/k8s-image-swapper/issues/111113) [#111696](https://github.com/estahn/k8s-image-swapper/issues/111696) [#108692](https://github.com/estahn/k8s-image-swapper/issues/108692)
* **deps:** bump k8s.io/client-go from 0.24.3 to 0.25.0 ([#324](https://github.com/estahn/k8s-image-swapper/issues/324)) ([f7c889f](https://github.com/estahn/k8s-image-swapper/commit/f7c889f4880f0d543c05f70759e8cbfef5c3d7ac))

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
