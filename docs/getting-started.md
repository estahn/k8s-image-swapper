# Getting started

This document will provide guidance for installing `k8s-image-swapper`.

## Prerequisites

`k8s-image-swapper` will automatically create image repositories and mirror images into them.
This requires certain permissions for your target registry (_only AWS ECR supported atm_).

Before you get started choose a namespace to install `k8s-image-swapper` in, e.g. `operations` or `k8s-image-swapper`.
Ensure the namespace exists and is configured as your current context[^1].
All examples below will omit the namespace.

### AWS ECR as target registry

AWS supports a variety of authentication strategies.
`k8s-image-swapper` uses the official Amazon AWS SDK and therefore supports [all available authentication strategies](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).
Choose from one of the strategies below or an alternative if needed.

#### IAM credentials

1. Create an IAM user (e.g. `k8s-image-swapper`) with permissions[^2] to create ECR repositories and upload container images.
   An IAM policy example can be found in the footnotes[^2].
2. Create a Kubernetes secret (e.g. `k8s-image-swapper-aws`) containing the IAM credentials you just obtained, e.g.

    ```bash
    kubectl create secret generic k8s-image-swapper-aws \
      --from-literal=aws_access_key_id=<...> \
      --from-literal=aws_secret_access_key=<...>
    ```

#### Service Account

TBD

## Helm

1. Add the Helm chart repository:
   ```bash
   helm repo add estahn https://estahn.github.io/charts/
   ```
2. Update the local chart information:
   ```bash
   helm repo update
   ```
3. Install `k8x-image-swapper`
   ```
   helm install k8s-image-swapper estahn/k8s-image-swapper \
     --set config.target.registry.aws.accountId=$AWS_ACCOUNT_ID \
     --set config.target.registry.aws.region=$AWS_DEFAULT_REGION \
     --set awsSecretName=k8s-image-swapper-aws
   ```

!!! note
    `awsSecretName` is not required for the Service Account method and instead the service account is annotated:
    ```yaml
    serviceAccount:
    create: true
    annotations:
      eks.amazonaws.com/role-arn: ${oidc_image_sawpper_role_arn}
    ```


[^1]: Use a tool like [kubectx & kubens](https://github.com/ahmetb/kubectx) for convienience.
[^2]:
    ??? tldr "IAM Policy"
        ```json
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "",
                    "Effect": "Allow",
                    "Action": [
                        "ecr:GetAuthorizationToken",
                        "ecr:DescribeRepositories",
                        "ecr:DescribeRegistry"
                    ],
                    "Resource": "*"
                },
                {
                    "Sid": "",
                    "Effect": "Allow",
                    "Action": [
                        "ecr:UploadLayerPart",
                        "ecr:PutImage",
                        "ecr:ListImages",
                        "ecr:InitiateLayerUpload",
                        "ecr:GetDownloadUrlForLayer",
                        "ecr:CreateRepository",
                        "ecr:CompleteLayerUpload",
                        "ecr:BatchGetImage",
                        "ecr:BatchCheckLayerAvailability"
                    ],
                    "Resource": "arn:aws:ecr:*:123456789:repository/*"
                }
            ]
        }
        ```

        !!! tip "Further restricting access"
            The resource configuration allows access to all AWS ECR repositories within the account 123456789.
            Restrict this further by repository name or tag.
            `k8s-image-swapper` will create repositories with the source registry as prefix, e.g. `nginx` --> `docker.io/library/nginx:latest`.
