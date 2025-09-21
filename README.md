# kdex-crds

This project contains a set of Custom Resource Definitions (CRDs) for managing micro-frontend applications in Kubernetes.

## Description

The `kdex-crds` project provides a declarative way to manage micro-frontends in a Kubernetes environment. It is composed of the following CRDs:

- `MicroFrontEndApp`: Represents a micro-frontend application, including its source code and the custom elements it exposes.
- `MicroFrontEndPageArchetype`: Defines the structure of an App Server page, including optional default header, footer, and navigation.
- `MicroFrontEndPageBinding`: Binds a set of `MicroFrontEndApp`s or raw HTML fragments to a page, and allows overriding the default header, footer, and navigation defined in the archetype.
- `MicroFrontEndPageFooter`: Defines the content of an App Server page footer section.
- `MicroFrontEndPageHeader`: Defines the content of an App Server page header section.
- `MicroFrontEndPageNavigation`: Defines the content of an App Server page navigation section.

These CRDs work together to provide a flexible and declarative way to manage micro-frontends in a Kubernetes environment.

See [CRD_REFERENCE.md](CRD_REFERENCE.md) for reference documentation.

## Design Notes

The architecture of the KDEX App Server's micro-frontend pages follows Semantic HTML as described in [MDN's Structuring documents](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/Structuring_documents).

### TODO

- Add MicroFrontEndHost CRD to form the nexus under which page bindings are collected.
    Investigate if adding a reference to each page binding is the correct approach.
    This CRD should also hold the common metadata about the host such as the name of the Organization, the default stylesheet and so on.
    It's possible that this resource should also result in either a managed Ingress instance or a managed Gateway API HTTPRoute.
- Implement policy for `strict` vs. `non-strict` app compliancy.
    When the strict policy is enabled, an app may not embed JavaScript dependencies. Validation of the application source code will fail if dependencies are not fully externalized.
    This should probably be set on the Host CRD in order to place this policy outside the hands of app developers.
    A Host which defines the `script` app policy must not accept apps which do not comply. While a non-strict Host may accept both strict and non-strict apps.

## Getting Started

### Prerequisites
- go version v1.24.5+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/kdex-crds:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/kdex-crds:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/kdex-crds:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/kdex-crds/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing

We welcome contributions from the community. If you find a bug or have a feature request, please open an issue. If you want to contribute code, please open a pull request.

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
