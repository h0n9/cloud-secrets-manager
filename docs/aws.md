# AWS(Amazon Web Services)

`cloud-secrets-injector` uses
[`aws-sdk-go-v2`](https://github.com/aws/aws-sdk-go-v2) client to interact with
AWS API server.

## Secrets Manager

Create a new secret with [awscli](https://aws.amazon.com/cli/) or on console
page.

## Credentials (IAM)

The client requires credentials which consist of an access key and secret access
key in general. There are several ways to [specify
credentials](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials)
to `cloud-secrets-injector`.

### Setup

(Recommended) [Using IAM roles for service
accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)
is highly recommended to allow `cloud-secrets-injector` to interact with AWS API
server. It's simple but secure.

1. Create OIDC provider
2. Create an IAM role and policy
    1. Set permission to get secret value
        ```json
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "VisualEditor0",
                    "Effect": "Allow",
                    "Action": [
                        "secretsmanager:GetSecretValue",
                    ],
                    "Resource": "arn:aws:secretsmanager:ap-northeast-2:123456789012:secret:secret-name"
                }
            ]
        }
        ```
    2. Set trusted relationship
        ```json
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Federated": "arn:aws:iam::123456789012:oidc-provider/oidc.eks.ap-northeast-2.amazonaws.com/id/0123456789ABCDEF0123456789ABCDEF"
                    },
                    "Action": "sts:AssumeRoleWithWebIdentity",
                    "Condition": {
                        "StringEquals": {
                            "oidc.eks.ap-northeast-2.amazonaws.com/id/0123456789ABCDEF0123456789ABCDEF:sub": "system:serviceaccount:namespace:service-account-name"
                        }
                    }
                }
            ]
        }
        ```
3. Associate the IAM role to a Kubernetes service account
    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      annotations:
        eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/iam-role-name
    ```

That's all! When you're ready, apply the `Deployment`, `Service Account`
manifests with kubectl.

It's going to work as it should, just like 🧈.

## Example

Please refer the following `sample-deployment.yaml`:
```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: busybox
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/testbed-role
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: busybox
spec:
  selector:
    matchLabels:
      app: busybox
  template:
    metadata:
      labels:
        app: busybox
      annotations:
        cloud-secrets-manager.h0n9.postie.chat/provider: aws
        cloud-secrets-manager.h0n9.postie.chat/secret-id: testbed-secret
        cloud-secrets-manager.h0n9.postie.chat/template: |
          {{ range $k, $v := . }}export {{ $k }}={{ $v }}
          {{ end }}
        cloud-secrets-manager.h0n9.postie.chat/output: /secrets/env
    spec:
      serviceAccountName: busybox
      containers:
      - name: busybox
        image: busybox:1.34.1
        command:
          - /bin/sh
          - -c
          - cat /secrets/env && sleep 3600
        resources:
          limits:
            memory: "64Mi"
            cpu: "100m"
```

Set label `cloud-secrets-injector=enabled` on namespace `testbed`:
```bash
kubectl create namespaces testbed
kubectl label namespaces testbed cloud-secrets-injector=enabled
```

Apply the deployment manifest:
```bash
kubectl apply -f sample-deployment.yaml -n testbed
```
