# Cloud Run httpbin (OpenTofu)

Provisions a public Cloud Run service running `ghcr.io/mccutchen/go-httpbin:latest` behind the custom domain `httpbin-test.goat.build.kyma-project.io`.

## Layout

- `bootstrap/` — one-shot module that creates the GCS bucket used as remote state for this module. Run once by a human with `roles/storage.admin`. See `bootstrap/README.md`.
- `*.tf` (this directory) — the actual deployment. State lives in the bucket created by `bootstrap/`.

## Prerequisites

Granted **out-of-band** on project `sap-se-cx-kyma-goat` to the GitHub Actions service account `github-actions@sap-se-cx-kyma-goat.iam.gserviceaccount.com` (and to any human applying locally):

| Role                                            | Why                                                                                          |
|-------------------------------------------------|----------------------------------------------------------------------------------------------|
| `roles/run.admin`                               | Manage the Cloud Run service and set its IAM policy (covers the `allUsers` invoker binding). |
| `roles/iam.serviceAccountUser`                  | Allow the SA to act as the Cloud Run runtime SA when deploying revisions.                    |
| `roles/storage.objectAdmin` on the state bucket | Read/write tfstate.                                                                          |
| `roles/dns.admin`                               | Manage record sets in the `goat-build-kyma-project-io` zone.                                 |

Other preconditions:

1. The bootstrap module has been applied (state bucket exists).
2. The parent domain (`kyma-project.io` or `goat.build.kyma-project.io`) is **verified** for the project in Google Search Console. Without verification, `google_cloud_run_domain_mapping` fails with `domain not verified`.

## Local apply

```sh
cd terraform
gcloud auth application-default login
tofu init
tofu plan
tofu apply
```

## First-apply DNS quirk

`google_cloud_run_domain_mapping.status[0].resource_records` is **populated asynchronously** by GCP. On the very first apply the list may be empty, in which case `google_dns_record_set` creates no records. Run `tofu apply` a second time once GCP has issued the records. This is documented standard behavior, not a bug.

After DNS records exist, Cloud Run issues a TLS certificate automatically. Expect a few minutes before `https://httpbin-test.goat.build.kyma-project.io/get` returns 200.

## Smoke test

```sh
curl -s -o /dev/null -w "%{http_code}\n" https://httpbin-test.goat.build.kyma-project.io/get
```

Expected: `200`.

## CI

`.github/workflows/cloud-run-deploy.yaml` runs `tofu plan` on PRs touching `terraform/**` and `tofu apply` on push to `main`. Authenticates via Workload Identity Federation — see PR #26 for how the WIF pool was set up.
