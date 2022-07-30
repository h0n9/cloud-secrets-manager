# Cloud Secrets Manager 🌤🔐🐳

This simple yet powerful tool **Cloud Secrets Manager** aims to simplify the way
to inject secrets strored on Cloud-based secrets managers into Kubernetes Pods,
functioning as [HashiCorp Vault's Agent Sidecar
Injector](https://www.vaultproject.io/docs/platform/k8s/injector).

## Cloud Providers

### Currently Supported
- AWS(Amazon Web Services): [Secrets Manager](https://aws.amazon.com/secrets-manager/)
- GCP(Google Cloud Platform): [Secret Manager](https://cloud.google.com/secret-manager) `(BETA)`

### TO-BE Supported
- Azure: [Key Vault](https://azure.microsoft.com/services/key-vault/#getting-started)
- Hashicorp: [Vault](https://www.vaultproject.io)

## Concept

### Constitution
- `cloud-secrets-controller`
- `cloud-secrets-injector`

### Step-by-step
1. `cloud-secrets-controller` watches incoming `/mutate`, `/validate` webhooks
from Kubernetes API server.
2. When pods are created or updated in a namespace labeled with
`cloud-secrets-injector: true`, Kubernetes API server sends requests to
`cloud-secrets-controller` webhook server.
3. `cloud-secrets-controller` mutates the pod's manifests by injecting an init
container `cloud-secrets-injector` into the pod and mounting a temporary
directory as a volume on the init and origin containers.
4. When it comes to initializing the pods, the init container
`cloud-secrets-injector` requests secret values, with a secret key id, from
secret providers and stores them in the temporary directory.
5. Once `cloud-secrets-injector` has successfully completed its role, the origin
container starts running as defined on the manifest.

## Installation

### Prerequisites
- Kubernetes Cluster
- `kubectl`
- `helm`

### Using Helm chart
```bash
kubectl create namespaces cloud-secrets-manager
helm repo add h0n9 https://h0n9.github.io/helm-charts
helm upgrade --install -n cloud-secrets-manager cloud-secrets-manager h0n9/cloud-secrets-manager
```

You can check out the official Helm chart repository
[h0n9/helm-charts](https://github.com/h0n9/helm-charts).

By pressing the `⭐️ Star` button above, be the first to get notified of launch
of other new charts.

## Usage

### Annotations

The following annotatins are required to inject `cloud-secrets-injector` into
pods:

| **Key**                                            | **Required** | **Description**           | **Example**                                              |
|----------------------------------------------------|--------------|---------------------------|----------------------------------------------------------|
| `cloud-secrets-manager.h0n9.postie.chat/provider`  | true         | Cloud Provider Name       | `aws`                                                    |
| `cloud-secrets-manager.h0n9.postie.chat/secret-id` | true         | Secret Name               | `very-precious-secret`                                   |
| `cloud-secrets-manager.h0n9.postie.chat/template`  | true         | Template for secret value | ```{{ range $k, $v := . }}{{ $k }}={{ $v }} {{ end }}``` |
| `cloud-secrets-manager.h0n9.postie.chat/output`    | true         | File path for output      | `/secrets/env`                                           |
| `cloud-secrets-manager.h0n9.postie.chat/injected`  | false        | Identifier for injection  | `false`                                                  |

### Providers

Supported providers require the annotations mentioned above in common. However,
the authentication method may differ depending on the provider. Please refer the
following explanation.

- [AWS(Amazon Web Services)](docs/aws.md)

### Example

Please refer the following `sample-deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: busybox
  namespace: testbed
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
        cloud-secrets-manager.h0n9.postie.chat/secret-id: dev/test
        cloud-secrets-manager.h0n9.postie.chat/template: |
          {{ range $k, $v := . }}export {{ $k }}={{ $v }}
          {{ end }}
        cloud-secrets-manager.h0n9.postie.chat/output: /secrets/env
    spec:
      containers:
      - name: busybox
        image: busybox:1.34.1
        command:
          - /bin/sh
          - -c
          - cat /secrets/env && sleep 3600
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

### Environment variables

#### AWS

| **Name**                | **Default** | **Required** |
|-------------------------|-------------|--------------|
| `AWS_ACCESS_KEY_ID`     |             | false        |
| `AWS_SECRET_ACCESS_KEY` |             | false        |

Please don't forget to pass credentials, referring to [Specifying
Credentials](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials)
page.
