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

#### Using ECR registries cross-account

Although ECR allows creating registry policy that allows reposistories creation from different account, there's no way to push anything to these repositories.
ECR resource-level policy can not be applied during creation, and to apply it afterwards we need ecr:SetRepositoryPolicy permission, which foreign account doesn't have.

One way out of this conundrum is to assume the role in target account

```yaml
target:
  type: aws
  aws:
    accountId: 123456789
    region: ap-southeast-2
    role: arn:aws:iam::123456789012:role/roleName
```
!!! note
    Make sure that target role has proper trust permissions that allow to assume it cross-account

!!! note
    In order te be able to pull images from outside accounts, you will have to apply proper access policy


#### Access policy

You can specify the access policy that will be applied to the created repos in config. Policy should be raw json string.
For example:
```yaml
target:
  aws:
    accountId: 123456789
    region: ap-southeast-2
    role: arn:aws:iam::123456789012:role/roleName
    accessPolicy: '{
  "Statement": [
    {
      "Sid": "AllowCrossAccountPull",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": [
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:BatchCheckLayerAvailability"
      ],
      "Condition": {
        "StringEquals": {
          "aws:PrincipalOrgID": "o-xxxxxxxxxx"
        }
      }
    }
  ],
  "Version": "2008-10-17"
}'
```

#### Lifecycle policy

Similarly to access policy, lifecycle policy can be specified, for example:

```yaml
target:
  aws:
    accountId: 123456789
    region: ap-southeast-2
    role: arn:aws:iam::123456789012:role/roleName
    lifecyclePolicy: '{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Rule 1",
      "selection": {
        "tagStatus": "any",
        "countType": "imageCountMoreThan",
        "countNumber": 1000
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
'
```

#### Service Account

1. Create an Webidentity IAM role (e.g. `k8s-image-swapper`) with the following trust policy, e.g
```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${your_aws_account_id}:oidc-provider/${oidc_image_swapper_role_arn}"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "${oidc_image_swapper_role_arn}:sub": "system:serviceaccount:${k8s_image_swapper_namespace}:${k8s_image_swapper_serviceaccount_name}"
        }
      }
    }
  ]
}
```

2. Create and attach permission policy[^2] to the role from Step 1..

Note: You can see a complete example below in [Terraform](Terraform)

## Helm

1. Add the Helm chart repository:
   ```bash
   helm repo add estahn https://estahn.github.io/charts/
   ```
2. Update the local chart information:
   ```bash
   helm repo update
   ```
3. Install `k8s-image-swapper`
   ```
   helm install k8s-image-swapper estahn/k8s-image-swapper \
     --set config.target.aws.accountId=$AWS_ACCOUNT_ID \
     --set config.target.aws.region=$AWS_DEFAULT_REGION \
     --set awsSecretName=k8s-image-swapper-aws
   ```

!!! note
    `awsSecretName` is not required for the Service Account method and instead the service account is annotated:
    ```yaml
    serviceAccount:
      create: true
      annotations:
        eks.amazonaws.com/role-arn: ${oidc_image_swapper_role_arn}
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

### Terraform

- Full example of helm chart deployment with AWS service account setup.


```
data "aws_caller_identity" "current" {
}

variable "cluster_oidc_provider" {
  default = "oidc.eks.ap-southeast-1.amazonaws.com/id/ABCDEFGHIJKLMNOPQRSTUVWXYZ012345"
  description = "example oidc endpoint that is created during eks deployment"
}

variable  "cluster_name" {
  default = "test"
  description = "name of the eks cluster being deployed to"
}


variable  "region" {
  default = "ap-southeast-1"
  description = "name of the eks cluster being deployed to"
}

variable "k8s_image_swapper_namespace" {
  default     = "kube-system"
  description = "namespace to install k8s-image-swapper"
}

variable "k8s_image_swapper_name" {
  default     = "k8s-image-swapper"
  description = "name for k8s-image-swapper release and service account"
}

#k8s-image-swapper helm chart
resource "helm_release" "k8s_image_swapper" {
  name       = var.k8s_image_swapper_name
  namespace  = "kube-system"
  repository = "https://estahn.github.io/charts/"
  chart   = "k8s-image-swapper"
  keyring = ""
  version = "1.0.1"
  values = [
    <<YAML
config:
  dryRun: true
  logLevel: debug
  logFormat: console

  source:
    # Filters provide control over what pods will be processed.
    # By default all pods will be processed. If a condition matches, the pod will NOT be processed.
    # For query language details see https://jmespath.org/
    filters:
      - jmespath: "obj.metadata.namespace != 'default'"
      - jmespath: "contains(container.image, '.dkr.ecr.') && contains(container.image, '.amazonaws.com')"
  target:
    aws:
      accountId: "${data.aws_caller_identity.current.account_id}"
      region: ${var.region}

secretReader:
  enabled: true

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # Specifies annotations for this service account
  annotations:
    eks.amazonaws.com/role-arn: "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${aws_iam_role.k8s_image_swapper.name}"
YAML
    ,
  ]
}

#iam policy for k8s-image-swapper service account
resource "aws_iam_role_policy" "k8s_image_swapper" {
  name = "${var.cluster_name}-${var.k8s_image_swapper_name}"
  role = aws_iam_role.k8s_image_swapper.id

  policy = <<-EOF
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
            "Resource": [
              "arn:aws:ecr:*:${data.aws_caller_identity.current.account_id}:repository/docker.io/*",
              "arn:aws:ecr:*:${data.aws_caller_identity.current.account_id}:repository/quay.io/*"
	    ]
        }
    ]
}
EOF
}

#role for k8s-image-swapper service account
resource "aws_iam_role" "k8s_image_swapper" {
  name               = "${var.cluster_name}-${var.k8s_image_swapper_name}"
  assume_role_policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${replace(var.cluster_oidc_provider, "/https:///", "")}"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "${replace(var.cluster_oidc_provider, "/https:///", "")}:sub": "system:serviceaccount:${var.k8s_image_swapper_namespace}:${var.k8s_image_swapper_name}"
        }
      }
    }
  ]
}
EOF
}

```
