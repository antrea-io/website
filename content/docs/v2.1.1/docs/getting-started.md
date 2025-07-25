# Getting Started

Antrea is super easy to install. All the Antrea components are
containerized and can be installed using the Kubernetes deployment
manifest.

![antrea-demo](https://user-images.githubusercontent.com/2495809/94325574-e7876500-ff53-11ea-9ecd-6dedef339fac.gif)

## Ensuring requirements are satisfied

### NodeIPAM

Antrea relies on `NodeIPAM` for per-Node CIDR allocation. `NodeIPAM` can run
within the Kubernetes `kube-controller-manager`, or within the Antrea
Controller.

#### NodeIPAM within kube-controller-manager

When using `kubeadm` to create the Kubernetes cluster, passing
`--pod-network-cidr=<CIDR Range for Pods>` to `kubeadm init` will enable
`NodeIpamController`. Clusters created with kubeadm will always have
`CNI` plugins enabled. Refer to
[Creating a cluster with kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm)
for more information about setting up a Kubernetes cluster with `kubeadm`.

When the cluster is deployed by other means then:

* To enable `NodeIpamController`, `kube-controller-manager` should be started
with the following flags:
  - `--cluster-cidr=<CIDR Range for Pods>`
  - `--allocate-node-cidrs=true`

* To enable `CNI` network plugins, `kubelet` should be started with the
`--network-plugin=cni` flag.

* To enable masquerading of traffic for Service cluster IP via iptables,
`kube-proxy` should be started with the `--cluster-cidr=<CIDR Range for Pods>`
flag.

#### NodeIPAM within Antrea Controller

For further info about running NodeIPAM within Antrea Controller, see
[Antrea IPAM Capabilities](antrea-ipam.md)

### Open vSwitch

As for OVS, when using the built-in kernel module, kernel version >= 4.6 is
required. On the other hand, when building it from OVS sources, OVS
version >= 2.6.0 is required.

Red Hat Enterprise Linux and CentOS 7.x use kernel 3.10, but as changes to
OVS kernel modules are regularly backported to these kernel versions, they
should work with Antrea, starting with version 7.4.

In case a node does not have a supported OVS module installed,
you can install it following the instructions at:
[Installing Open vSwitch](https://docs.openvswitch.org/en/latest/intro/install/).
Please be aware that the `vport-stt` module is not in the Linux tree and needs to be
built from source, please build and load it manually before STT tunneling is enabled.

Some experimental features disabled by default may have additional requirements,
please refer to the [Feature Gates documentation](feature-gates.md) to determine
whether it applies to you.

Antrea will work out-of-the-box on most popular Operating Systems. Known issues
encountered when running Antrea on specific OSes are documented
[here](os-issues.md).

There are also a few network prerequisites which need to be satisfied, and they depend
on the tunnel mode you choose, please check [network requirements](./network-requirements.md).

## Installation / Upgrade

To deploy a released version of Antrea, pick a deployment manifest from the
[list of releases](https://github.com/antrea-io/antrea/releases). For any
given release `<TAG>` (e.g. `v0.1.0`), you can deploy Antrea as follows:

```bash
kubectl apply -f https://github.com/antrea-io/antrea/releases/download/<TAG>/antrea.yml
```

To deploy the latest version of Antrea (built from the main branch), use the
checked-in [deployment yaml](https://github.com/antrea-io/antrea/blob/v2.1.1/build/yamls/antrea.yml):

```bash
kubectl apply -f https://raw.githubusercontent.com/antrea-io/antrea/main/build/yamls/antrea.yml
```

You can use the same `kubectl apply` command to upgrade to a more recent version
of Antrea.

Antrea supports some experimental features that can be enabled or disabled,
please refer to the [Feature Gates documentation](feature-gates.md) for more
information.

### Windows support

If you want to add Windows Nodes to your cluster, please refer to these
[installation instructions](windows.md).

### ARM support

Starting with v1.0, Antrea supports arm64 and arm/v7 Nodes. The installation
instructions do not change when some (or all) Linux Nodes in a cluster use an
ARM architecture: the same deployment YAML can be used, as the
`antrea/antrea-agent-ubuntu` and `antrea/antrea-controller-ubuntu` Docker images
are actually manifest lists with support for the amd64, arm64 and arm/v7
architectures.

Note that while we do run a subset of the Kubernetes conformance tests on both
the arm/v7 and arm64 Docker images (using [k3s](https://k3s.io/) as the
Kubernetes distribution), our testing is not as thorough as for the amd64
image. However, we do not anticipate any issue.

### Install with Helm

Starting with v1.8, Antrea can be installed and updated with Helm. Please refer
to these [installation instructions](helm.md).

### Deploying Antrea on a Cluster with Existing CNI

The instructions above only apply when deploying Antrea in a new cluster. If you
need to migrate your existing cluster from another CNI plugin to Antrea, you
will need to do the following:

* Delete previous CNI, including all resources (K8s objects, iptables rules,
interfaces, ...) created by that CNI.
* Deploy Antrea.
* Restart all Pods in the CNI network in order for Antrea to set-up networking
for them. This does not apply to Pods which use the Node's network namespace
(i.e. Pods configured with `hostNetwork: true`). You may use `kubectl drain` to
drain each Node or reboot all your Nodes.

While this is in-progress, networking will be disrupted in your cluster. After
deleting the previous CNI, existing Pods may not be reachable anymore.

For example, when migrating from Flannel to Antrea, you will need to do the
following:

1. Delete Flannel with `kubectl delete -f <path to your Flannel YAML manifest>`.
2. Delete Flannel bridge and tunnel interface with `ip link delete flannel.1 &&
ip link delete flannel cni0` **on each Node**.
3. Ensure [requirements](#ensuring-requirements-are-satisfied) are satisfied.
4. [Deploy Antrea](#installation--upgrade).
5. Drain and uncordon Nodes one-by-one. For each Node, run `kubectl drain
--ignore-daemonsets <node name> && kubectl uncordon <node name>`. The
`--ignore-daemonsets` flag will ignore DaemonSet-managed Pods, including the
Antrea Agent Pods. If you have any other DaemonSet-managed Pods (besides the
Antrea ones and system ones such as kube-proxy), they will be ignored and will
not be drained from the Node. Refer to the [Kubernetes
documentation](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/)
for more information. Alternatively, you can also restart all the Pods yourself,
or simply reboot your Nodes.

To build the image locally, you can follow the instructions in the [Contributor
Guide](../CONTRIBUTING.md#building-and-testing-your-change).

### Deploying Antrea in Kind

To deploy Antrea in a [Kind](https://github.com/kubernetes-sigs/kind) cluster,
please refer to this [guide](kind.md).

### Deploying Antrea in Minikube

To deploy Antrea in a [Minikube](https://github.com/kubernetes/minikube) cluster,
please refer to this [guide](minikube.md).

### Deploying Antrea in Rancher Managed Cluster

To deploy Antrea in a [Rancher](https://github.com/rancher/rancher) managed cluster,
please refer to this [guide](kubernetes-installers.md#rancher).

### Deploying Antrea in AKS, EKS, and GKE

Antrea can work with cloud managed Kubernetes services, and can be deployed to
AKS, EKS, and GKE clusters.

* To deploy Antrea to an AKS or an AKS Engine cluster, please refer to [the AKS installation guide](aks-installation.md).
* To deploy Antrea to an EKS cluster, please refer to [the EKS installation guide](eks-installation.md).
* To deploy Antrea to a GKE cluster, please refer to [the GKE installation guide](gke-installation.md).

### Deploying Antrea with Custom Certificates

By default, Antrea generates the certificates needed for itself to run. To
provide your own certificates, please refer to [Securing Control Plane](securing-control-plane.md).

### Antctl: Installation and Usage

To use antctl, the Antrea command-line tool, please refer to [this guide](antctl.md).

## Features

### Antrea Network Policy

Besides Kubernetes NetworkPolicy, Antrea also implements its own Network Policy
CRDs, which provide advanced features including: policy priority, tiering, deny
action, external entity, and policy statistics. For more information on usage of
Antrea Network Policies, refer to the [Antrea Network Policy document](antrea-network-policy.md).

### Egress

Antrea supports specifying which egress (SNAT) IP the traffic from the selected
Pods to the external network should use and which Node the traffic should leave
the cluster from. For more information, refer to the [Egress document](egress.md).

### Network Flow Visibility

Antrea supports exporting network flow information using IPFIX, and provides a
reference cookbook on how to visualize the exported network flows using Elastic
Stack and Kibana dashboards. For more information, refer to the [network flow
visibility document](network-flow-visibility.md).

### NoEncap and Hybrid Traffic Modes

Besides the default `Encap` mode, in which Pod traffic across Nodes will be
encapsulated and sent over tunnels, Antrea also supports `NoEncap` and `Hybrid`
traffic modes. In `NoEncap` mode, Antrea does not encapsulate Pod traffic, but
relies on the Node network to route the traffic across Nodes. In `Hybrid` mode,
Antrea encapsulates Pod traffic when the source Node and the destination Node
are in different subnets, but does not encapsulate when the source and the
destination Nodes are in the same subnet. Refer to [this guide](noencap-hybrid-modes.md)
to learn how to configure Antrea with `NoEncap` or `Hybrid` mode.

### Antrea Web UI

Antrea comes with a web UI, which can show runtime information of Antrea
components and perform Antrea Traceflow operations. Please refer to the [Antrea
UI repository](https://github.com/antrea-io/antrea-ui) for installation
instructions and more information.

### OVS Hardware Offload

Antrea can offload OVS flow processing to the NICs that support OVS kernel
hardware offload using TC. The hardware offload can improve OVS performance
significantly. For more information on how to configure OVS offload, refer to
the [OVS hardware offload guide](ovs-offload.md).

### Prometheus Metrics

Antrea supports exporting metrics to Prometheus. For more information, refer to
the [Prometheus integration document](prometheus-integration.md).

### Support for Services of type LoadBalancer

By leveraging Antrea's Service external IP management feature or configuring
MetalLB to work with Antrea, Services of type LoadBalancer can be supported
without requiring an external LoadBalancer. To learn more information, please
refer to the [Service LoadBalancer document](service-loadbalancer.md).

### Traceflow

Traceflow is a very useful network diagnosis feature in Antrea. It can trace
and report the forwarding path of a specified packet in the Antrea network.
For usage of this feature, refer to the [Traceflow user guide](traceflow-guide.md).

### Traffic Encryption

Antrea supports encrypting traffic between Linux Nodes using IPsec or WireGuard.
To deploy Antrea with traffic encryption enabled, please refer to [this guide](traffic-encryption.md).

### Antrea Multi-cluster

Antrea Multi-cluster implements Multi-cluster Service API, which allows users to
create multi-cluster Services that can be accessed cross clusters in a
ClusterSet. Antrea Multi-cluster also supports Antrea ClusterNetworkPolicy
replication. Multi-cluster admins can define ClusterNetworkPolicies to be
replicated across the entire ClusterSet, and enforced in all member clusters.
To learn more information about Antrea Multi-cluster, please refer to the
[Antrea Multi-cluster user guide](multicluster/user-guide.md).
