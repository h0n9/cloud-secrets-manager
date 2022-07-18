# Cloud Secrets Manager

This simple yet powerful tool named **Cloud Secrets Manager** aims to simplify
the way to inject secrets strored on Cloud-based secrets managers into
Kubernetes Pods, functioning as [HashiCorp Vault's Agent Sidecar
Injector](https://www.vaultproject.io/docs/platform/k8s/injector).

## Supported Cloud Providers
- AWS(Amazon Web Services): Secrets Manager

## Installation

### Prerequisites
- Kubernetes Cluster
- `kubectl`
- `helm`

### Using Helm chart
```bash
kubectl create namespaces cloud-secrets-controller
helm upgrade --install -n cloud-secrets-controller cloud-secrets-controller helm/
```

(The official Helm chart repository should be ready very soon 🙌. Be the first
to get notified of its launch by pressing the `⭐️ Star` button above.)

## Environment Variables

### Common

| **Name**          | **Default**                                                                | **Required** |
|-------------------|----------------------------------------------------------------------------|--------------|
| `PROVIDER_NAME`   | `aws`                                                                      | false        |
| `TEMPLATE_BASE64` | `e3sgcmFuZ2UgJGssICR2IDo9IC4gfX1be3sgJGsgfX1dCnt7ICR2IH19Cgp7eyBlbmQgfX0K` | false        |
| `TEMPLATE_FILE`   |                                                                            | false        |
| `OUTPUT_FILE`     | `output`                                                                   | false        |

### AWS

| **Name**                | **Default** | **Required** |
|-------------------------|-------------|--------------|
| `SECRET_ID`             |             | true         |
| `AWS_ACCESS_KEY_ID`     |             | false        |
| `AWS_SECRET_ACCESS_KEY` |             | false        |

Please don't forget to pass credentials, referring to [Specifying
Credentials](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials)
page.

## Usage

### Controller

### Injector
