# Cloud Secrets Manager 🌤🔐🐳

This simple yet powerful tool **Cloud Secrets Manager** aims to simplify the way
to inject secrets strored on Cloud-based secrets managers into Kubernetes Pods,
functioning as [HashiCorp Vault's Agent Sidecar
Injector](https://www.vaultproject.io/docs/platform/k8s/injector).

Also, it provides a convenient CLI tool with features like `list` and `edit` to
make secret management easier than using the Cloud Console. If you want to jump
into the CLI tool, please refer to the [CLI Tool](#cli-tool) section right away.

## Contents
- [Supported Cloud Providers](#cloud-providers)
  - [Current](#current)
  - [Planned](#planned)
- [Concept](#concept)
  - [Constitution](#constitution)
  - [Step-by-step](#step-by-step)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Using Helm chart](#using-helm-chart)
- [Usage](#usage)
  - [Annotations](#annotations)
  - [Providers](#providers)
  - [CLI Tool](#cli-tool)

## Supported Cloud Providers

### Current
- AWS(Amazon Web Services): [Secrets Manager](https://aws.amazon.com/secrets-manager/)
- GCP(Google Cloud Platform): [Secret Manager](https://cloud.google.com/secret-manager) `(BETA)`

### Planned
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
`cloud-secrets-injector: enabled`, Kubernetes API server sends requests to
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

#### Annotations for Multiple Secrets Injection

From the version `v0.4`, multiple secrets can be injected into pods by defining
the annotations as follows:

```yaml
cloud-secrets-manager.h0n9.postie.chat/provider: aws
cloud-secrets-manager.h0n9.postie.chat/secret-id: secrets-env
cloud-secrets-manager.h0n9.postie.chat/output: /secrets/env
cloud-secrets-manager.h0n9.postie.chat/template: |
  {{ range $k, $v := . }}export {{ $k }}={{ $v }}
  {{ end }}
cloud-secrets-manager.h0n9.postie.chat/provider-config-app: aws
cloud-secrets-manager.h0n9.postie.chat/secret-id-config-app: secrets-config
cloud-secrets-manager.h0n9.postie.chat/output-config-app: /config/application.yaml
cloud-secrets-manager.h0n9.postie.chat/template-config-app: |
  {{ .application-yaml }}
cloud-secrets-manager.h0n9.postie.chat/provider-config-secrets: aws
cloud-secrets-manager.h0n9.postie.chat/secret-id-config-secrets: secrets-config
cloud-secrets-manager.h0n9.postie.chat/output-config-secrets: /config/secrets.yaml
cloud-secrets-manager.h0n9.postie.chat/template-config-secrets: |
  {{ .secrets-yaml }}
```

Just add `<secret-name>` at the end of each annotation key, like
`cloud-secrets-manager.h0n9.postie.chat/provider-<secret-name>`. That's it!

### Providers

Supported providers require the annotations mentioned above in common. However,
the authentication method may differ depending on the provider. Please refer the
following explanation.

- [AWS(Amazon Web Services)](docs/aws.md)
- [GCP(Google Cloud Platform)](docs/gcp.md)

### CLI Tool

#### Installation

Cloud Secrets Manager can be installed via Homebrew:

```bash
$ brew install h0n9/devops/cloud-secrets-manager
```

That's it! You can now use the `cloud-secrets-manager` command.

> For `aws-vault` users, you can use the following command to execute the
`cloud-secrets-manager` command with the specified AWS profile:

```bash
$ aws-vault exec <profile> -- cloud-secrets-manager <command>
```

#### List Secrets

```bash
$ cloud-secrets-manager secrets list --provider aws --limit 3
dev/hello-world
dev/very-precious-secret
dev/another-secret
```
The `--limit` option is available to limit the number of secrets to be listed.

#### Edit Secret

```bash
$ cloud-secrets-manager secrets edit --provider aws --secret-id dev/very-precious-secret
```

A text editor will be opened with the secret value. After editing, save and
close the editor to update the secret value. If you want to cancel the editing,
just close the editor without saving.

If you want to use a specific editor, set the `EDITOR` environment variable.

```bash
$ export EDITOR=nano
$ cloud-secrets-manager secrets edit --provider aws --secret-id dev/very-precious-secret
```
