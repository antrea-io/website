# Code and Documentation Generation

## CNI

Antrea uses [protoc](https://github.com/protocolbuffers/protobuf), [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)
and [protoc-gen-go-grpc](https://github.com/grpc/grpc-go) to generate CNI gRPC service code.

If you make any change to [cni.proto](https://github.com/antrea-io/antrea/blob/v2.1.1/pkg/apis/cni/v1beta1/cni.proto), you can
re-generate the code by invoking `make codegen`. Make sure that you commit your changes before
running the script, and that you commit the generated code after running it.

## Extension API Resources and Custom Resource Definitions

Antrea extends Kubernetes API with an extension APIServer and Custom Resource Definitions, and uses
[k8s.io/code-generator
(release-1.18)](https://github.com/kubernetes/code-generator/tree/release-1.18) to generate clients,
informers, conversions, protobuf codecs and other helpers. The resource definitions and their
generated codes are located in the conventional paths: `pkg/apis/<resource group>` for internal
types and `pkg/apis/<resource group>/<version>` for versioned types and `pkg/client/clientset` for
clients.

If you make any change to any `types.go`, you can re-generate the code by invoking `make
codegen`. Make sure that you commit your changes before running the script, and that you commit the
generated code after running it.

## Mocks

Antrea uses the [GoMock](https://github.com/uber-go/mock) framework for its unit tests.

If you add or modify interfaces that need to be mocked, please add or update `MOCKGEN_TARGETS` in
[update-codegen-dockerized.sh](https://github.com/antrea-io/antrea/blob/v2.1.1/hack/update-codegen-dockerized.sh) accordingly. All the mocks for a
given package will typically be generated in a sub-package called `testing`. For example, the mock
code for the interface `Baz` defined in the package `pkg/foo/bar` will be generated to
`pkg/foo/bar/testing/mock_bar.go`, and you can import it via `pkg/foo/bar/testing`.

Same as above, you can re-generate the mock source code (with `mockgen`) by invoking `make codegen`.
Make sure that you commit your changes before running the script, and that you commit the
generated code after running it.

## Generated Documentation

[Prometheus integration document](../prometheus-integration.md) contains a list
of supported metrics, which could be affected by third party component
changes. The collection of metrics is done from a running Kind deployment, in
order to reflect the current list of metrics which is exposed by Antrea
Controller and Agents.

To regenerate the metrics list within the document, use [make-metrics-doc.sh](https://github.com/antrea-io/antrea/blob/v2.1.1/hack/make-metrics-doc.sh)
with document location as a parameter.
