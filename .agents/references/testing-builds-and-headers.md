# Tests, Builds, and File Headers

## Unit Tests

```sh
make test-unit                        # Run all unit tests
go test -v ./pkg/provider/...         # Run tests for a specific folder
```

If changes are restricted to one folder, prefer `go test -v ./<folder>/...`
over `make test-unit` for faster feedback.

## Testing Expectations

- Test public functions.
- Test edge cases and error conditions.
- Mock external dependencies.
- Match the existing test style in the file being edited.

### azclient Recording Cleanup

Private DNS zones may still report deleted virtual-network links briefly after
the link deletion poller completes. Both affected recording suites use
`recording.DeletePrivateDNSZone`, which retries only HTTP 409
`CannotDeleteResource`, including deletion-poller errors, for at most five
minutes or the caller's deadline. Other errors fail immediately.

Offline cleanup and public-IP ETag regression tests do not require Azure
credentials. Live suites need fresh recordings after their request sequences
change. Preserve successful fixtures and re-record only affected package/cloud
pairs; a fixture saved by a failed test is not a successful recording.

### Service Gateway in Azure China

The generated Service Gateway `TestClient` explicitly skips when
`AZURE_CLOUD=AZURECHINACLOUD`, before running Ginkgo setup or cleanup. Its
required `Microsoft.Network/AllowServiceGateways` feature cannot currently be
self-service registered in the available China test subscriptions. This is a
test-environment limitation, not a claim that the service is unavailable
throughout Azure China. Public-Azure recording tests and offline Service Gateway
unit tests remain enabled. The exception lives in the client generator so
regeneration preserves it; revisit it when an enabled China subscription is
available.

## Building Images

### Zone-mismatch reproduction branch

`chengan/test-zone-mismatch` is test-only. Its `zone.Repository.ListZones`
implementation bypasses ARM VM zone discovery and returns only
`eastus2: [1,2,3,4,5]`, including on periodic refresh. Other regions have no
discovered VM zones in this build. The public-IP zone-selection fix from
[kubernetes-sigs/cloud-provider-azure#11224](https://github.com/kubernetes-sigs/cloud-provider-azure/pull/11224)
at commit `4d8b45db8d7105c85c143987ff913f5d8e32100b` is now applied:
new Standard public IPs use real `Microsoft.Network/publicIPAddresses` zone
mappings instead of the fake VM zones. Public-IP creation still calls real ARM.

Use only a disposable eastus2 cluster where you control the active CCM and have
confirmed that `Microsoft.Network/publicIPAddresses` supports zones 1,2,3,4 but
rejects 5. The injection-only image produces `InvalidAvailabilityZone` and a
pending Service. Build this fixed revision with a new image tag and replace that
CCM, leaving the same pending external `LoadBalancer` Service in place. Expect
the public IP to use zones 1,2,3,4 (order may vary) and the Service to acquire an
external IP despite the injected VM zone 5. A new external Service without an
existing or pre-created public IP also exercises the fixed path.

Existing public IPs do not exercise this path. Internal load-balancer frontends
still consume the fake VM zone list, so this build must not be used in production.
The PR's cache and discovery behavior is preserved for this comparison, including
its lack of public-IP zone refresh and Azure Stack bypass. Restore the normal CCM
and remove the test Service after the experiment.

### Image build commands

Before building images, ask the user for `IMAGE_REGISTRY` and `IMAGE_TAG`. For
the reusable agent workflow, use `.agents/skills/build-images/`.

```sh
IMAGE_REGISTRY=<registry> IMAGE_TAG=<tag> make build-ccm-image
IMAGE_REGISTRY=<registry> IMAGE_TAG=<tag> make build-node-image-linux
```

The CCM and CNM Dockerfiles use Microsoft Go. Starting with version 1.27,
`systemcrypto` is no longer configured through `GOEXPERIMENT`. For non-FIPS
development images, use `MS_GO_NOSYSTEMCRYPTO=1 make build-ccm-image`. This
produces a non-FIPS-compliant binary and must not be used for production or
FIPS validation.

## File Headers

Newly created tracked source files must carry the standard Kubernetes Apache 2.0
boilerplate header using the appropriate comment style for that file type.

Run `hack/verify-boilerplate.sh` or `make test-check` after adding new tracked
files when there is any doubt about header compliance.

Markdown files do not need the boilerplate header. For shared-skill docs under
`.agents/skills/`, also follow [shared skill authoring](../skills/authoring.md).
