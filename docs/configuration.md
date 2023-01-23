# Configuration

The configuration is managed via the config file `.k8s-image-swapper.yaml`.
Some options can be overriden via parameters, e.g. `--dry-run`.

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
* `immediate`: Submits the copy job to a process queue and waits for it to finish (deadline 8s).
* `force`: Attempts to immediately copy the image (deadline 8s).



## Source

This section configures details about the image source.

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

### AWS

The option `target.registry.aws` holds details about the target registry storing the images.
The AWS Account ID and Region is primarily used to construct the ECR domain `[ACCOUNTID].dkr.ecr.[REGION].amazonaws.com`.

!!! example
    ```yaml
    target:
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
