# Configuration

The configuration is managed via the config file `.k8s-image-swapper.yaml`.
Some options can be overridden via parameters, e.g. `--dry-run`.

## Dry Run

The option `dryRun` allows to run the webhook without executing the actions, e.g. repository creation,
image download and manifest mutation.

!!! example
    ```yaml
    dryRun: true
    ```

## Log Level & Format

The option `logLevel` & `logFormat` allow to adjust the verbosity and format (e.g. `json`, `console`).

!!! example
    ```yaml
    logLevel: debug
    logFormat: console
    ```

## ImageSwapPolicy

The option `imageSwapPolicy` (default: `exists`) defines the mutation strategy used.

* `always`: Will always swap the image regardless of the image existence in the target registry.
            This can result in pods ending in state ImagePullBack if images fail to be copied to the target registry.
* `exists`: Only swaps the image if it exits in the target registry.
            This can result in pods pulling images from the source registry, e.g. the first pod pulls
            from source registry, subsequent pods pull from target registry.

## ImageCopyPolicy

The option `imageCopyPolicy` (default: `delayed`) defines the image copy strategy used.

* `delayed`: Submits the copy job to a process queue and moves on.
* `immediate`: Submits the copy job to a process queue and waits for it to finish (deadline defined by `imageCopyDeadline`).
* `force`: Attempts to immediately copy the image (deadline defined by `imageCopyDeadline`).
* `none`: Do not copy the image.

## ImageCopyDeadline

The option `imageCopyDeadline` (default: `8s`) defines the duration after which the image copy if aborted.

This option only applies for `immediate` and `force` image copy strategies.


## Source

This section configures details about the image source.

### Registries

The option `source.registries` describes a list of registries to pull images from, using a specific configuration.

#### AWS

By providing configuration on AWS registries you can ask `k8s-image-swapper` to handle the authentication using the same credentials as for the target AWS registry.
This authentication method is the default way to get authorized by a private registry if the targeted Pod does not provide an `imagePullSecret`.

Registries are described with an AWS account ID and region, mostly to construct the ECR domain `[ACCOUNT_ID].dkr.ecr.[REGION].amazonaws.com`.

!!! example
    ```yaml
    source:
      registries:
        - type: aws
          aws:
            accountId: 123456789
            region: ap-southeast-2
        - type: aws
          aws:
            accountId: 234567890
            region: us-east-1
    ```
#### Generic

By providing configuration on Generic registries you can ask `k8s-image-swapper` to handle the authentication using 
username and password.

Registries are described with a repository URL, username and password.

!!! example
    ```yaml
    source:
      registries:
        - type: "generic"
          generic:
            repository: "repo1.azurecr.io"
            username: "username1"
            password: "pass1"
        - type: "generic"
          generic:
            repository: "repo2.azurecr.io"
            username: "username2"
            password: "pass2"
    ```


### Filters

Filters provide control over what pods will be processed.
By default, all pods will be processed.
If a condition matches, the pod will **NOT** be processed.

[JMESPath](https://jmespath.org/) is used as query language and allows flexible rules for most use-cases.

!!! info
    The data structure used for JMESPath is as follows:

    === "Structure"
        ```yaml
        obj:
          <Object Spec>
        container:
          <Container Spec>
        ```

    === "Example"
        ```yaml
        obj:
          metadata:
            name: static-web
            labels:
              role: myrole
          spec:
            containers:
              - name: web
                image: nginx
                ports:
                  - name: web
                    containerPort: 80
                    protocol: TCP
        container:
          name: web
          image: nginx
          ports:
            - name: web
              containerPort: 80
              protocol: TCP
        ```

Below you will find a list of common queries and/or ideas:

!!! tip "List of common queries/ideas"
    * Do not process if namespace equals `kube-system` (_Helm chart default_)
      ```yaml
      source:
        filters:
          - jmespath: "obj.metadata.namespace == 'kube-system'"
      ```
    *  Only process if namespace equals `playground`
       ```yaml
       source:
         filters:
           - jmespath: "obj.metadata.namespace != 'playground'"
       ```
    * Only process if namespace ends with `-dev`
      ```yaml
      source:
        filters:
          - jmespath: "ends_with(obj.metadata.namespace,'-dev')"
      ```
    * Do not process AWS ECR images
      ```yaml
      source:
        filters:
          - jmespath: "contains(container.image, '.dkr.ecr.') && contains(container.image, '.amazonaws.com')"
      ```

`k8s-image-swapper` will log the filter data and result in `debug` mode.
This can be used in conjunction with [JMESPath.org](https://jmespath.org/) which
has a live editor that can be used as a playground to experiment with more complex queries.

## Target

This section configures details about the image target.
The option `target` allows to specify which type of registry you set as your target (AWS, GCP...).
At the moment, `aws` and `gcp` are the only supported values.

### AWS

The option `target.aws` holds details about the target registry storing the images.
The AWS Account ID and Region is primarily used to construct the ECR domain `[ACCOUNTID].dkr.ecr.[REGION].amazonaws.com`.

!!! example
    ```yaml
    target:
      type: aws
      aws:
        accountId: 123456789
        region: ap-southeast-2
    ```


#### ECR Options

##### Tags

This provides a way to add custom tags to newly created repositories. This may be useful while looking at AWS costs.
It's a slice of `Key` and `Value`.

!!! example
    ```yaml
    target:
      type: aws
      aws:
        ecrOptions:
          tags:
            - key: cluster
              value: myCluster
    ```

### GCP

The option `target.gcp` holds details about the target registry storing the images.
The GCP location, projectId, and repositoryId are used to constrct the GCP Artifact Registry domain `[LOCATION]-docker.pkg.dev/[PROJECT_ID]/[REPOSITORY_ID]`.

!!! example
    ```yaml
    target:
      type: gcp
      gcp:
        location: us-central1
        projectId: gcp-project-123
        repositoryId: main
    ```

### Generic

The option `target.generic` holds details about the target registry storing the images.

!!! example
    ```yaml
    target:
      type: generic
      generic:
        repository: "repo2.azurecr.io"
        username: "username2"
        password: "pass2"
    ```
