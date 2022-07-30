# GCP(Google Cloud Platform)

`cloud-secrets-injector` uses
[`google-cloud-go`](https://github.com/googleapis/google-cloud-go) client to
interact with GCP API server.

## Secret Manager

Create a new secret with [gcloud](https://cloud.google.com/sdk/gcloud) or on
console page.

## Authentication (IAM)

The client requires service account's access key or workload identity to
authenticate to GCP API server.

### Setup

(Recommended) [Allow Pods to authenticate to Google Cloud APIs using Workload
identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
is highly recommended to allow `cloud-secrets-injector` to interact with the API
server. It's simple but secure.

1. Enable Workload Identity
2. Create an IAM service account for your application or use an existing IAM
service account instead
    1. Set permission to get secret value: `roles/secretmanager.secretAccessor`
        ```bash
        gcloud projects add-iam-policy-binding PROJECT_ID \
        --member "serviceAccount:GSA_NAME@GSA_PROJECT.iam.gserviceaccount.com" \
        --role "roles/secretmanager.secretAccessor"
        ```
    2. Allow the Kubernetes service account to impersonate the IAM service
    account by adding an IAM policy binding between the two service accounts.
        ```bash
        gcloud iam service-accounts add-iam-policy-binding GSA_NAME@GSA_PROJECT.iam.gserviceaccount.com \
        --role roles/iam.workloadIdentityUser \
        --member "serviceAccount:PROJECT_ID.svc.id.goog[NAMESPACE/KSA_NAME]"
        ```
3. Associate the IAM service account to a Kubernetes service account
    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      annotations:
        iam.gke.io/gcp-service-account: GSA_NAME@GSA_PROJECT.iam.gserviceaccount.com
    ```

4. Update the Deployment's Pod spec
    ```yaml
    spec:
      serviceAccountName: KSA_NAME
      nodeSelector:
        iam.gke.io/gke-metadata-server-enabled: "true"
    ```

That's all! When you're ready, apply the `Deployment`, `Service Account`
manifests with kubectl.

It's going to work as it should, just like 🧈.

