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

`chengan/test-zone-mismatch-11257` is test-only. It applies
[kubernetes-sigs/cloud-provider-azure#11257](https://github.com/kubernetes-sigs/cloud-provider-azure/pull/11257)
at commit `ab5167be5dd4023f374f20544b5d38b12753e0e9` while retaining fake VM zones
`eastus2: [1,2,3,4,5]`. The injection has moved from `zone.Repository.ListZones`
to `providerclient.Client.GetVirtualMachineSupportedZones`, so it does not
override the PR's resource-specific Network discovery.

The VM-specific helper is patched identically in `pkg/azclient/providerclient/`
and `vendor/sigs.k8s.io/cloud-provider-azure/pkg/azclient/providerclient/`.
The root CCM build uses the vendored copy; regenerating vendor will overwrite
this test-only patch unless it is reapplied. `TestVMZoneInjection` verifies the
root build's dependency still returns zone 5 without an ARM call. The fixed
Network discovery must not call that helper, so the old fake-zone startup log
is no longer expected. No Network provider response or resource PUT is faked.

Use only a disposable eastus2 cluster where you control the active CCM. The
previous custom 11224 image retained VM zone 5 for internal frontends and
reproduced `SubscriptionDoesNotSupportZone` for a pending internal Service.
Replace that CCM with a newly tagged image from this branch without deleting
the pending Service. Public IPs should use real public-IP zone mappings.
Internal frontends use real load-balancer zone mappings.

The observed ARM response has public-IP zones 1,2,3,4 but
`loadBalancers.zoneMappings: null`. Preserve that response: the PR will omit
zones for a new internal frontend. Verify provisioning, the resulting frontend
zones, and connectivity from inside the VNet separately. Success with unset
zones demonstrates provisioning recovery, not guaranteed zone redundancy.
Do not substitute public-IP or VM mappings for the missing load-balancer data.

For a strict same-hook before/after comparison, use the same VM-helper patches
in the pre-11257 build with its original VM-backed zone repository, then apply
11257 without changing those helper patches. Restore the normal CCM and remove
the test Service after the experiment.

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
