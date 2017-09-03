<h1>Ghee</h1>

Simple Kubernetes Multi-Cluster Controller - a commandline tool to create, update, and delete Kubernetes resources across multiple clusters.

<h2>Clusters</h2>
The user needs to add cluster to Ghee in order to control them. This includes where Ghee can reach it and authorization information.

Example:
```
name: first-cluster
address: 192.168.99.100
certAuthority: <ca>
clientCert: <cert>
clientKey: <key>
```

<h2>Resources</h2>
Each manipulation of a Kubernetes resource takes in a `Gheefile` which defines the Kubernetes manifests for the resource as well as the names of the clusters it should be created on and how many replicas should be available in each cluster. This allows a user to tailor the manifests to each cluster without having duplicate manifests with slight differences.

Example `Gheefile`:
```
- manifest:
  - apiVersion: v1
    kind: Namespace
    metadata:
      name: ghee
  - apiVersion: v1/beta1
    kind: Deployment
    metadata:
      name: hello-world
      namespace: ghee
    spec:
      containers:
        image: hello-world
  clusters:
  - first-cluster
  - second-cluster
  replicas:
    first-cluster: 1
    second-cluster: 3
- manifest:
  <another-manifest>
  clusters:
  - second-cluster
  - third-cluster
  replicas:
    second-cluster: 5
    third-cluster: 19
```

<h4>Supported resources:</h4>

- deployment (apiVerison apps/v1beta1)
- statefulSet (apiVersion apps/v1beta1)
- clusterRole (apiVersion rbac.authorization.k8s.io/v1beta1)
- clusterRolebinding (apiVersion rbac.authorization.k8s.io/v1beta1)
- configMap (apiVersion v1)
- namespace (apiVersion v1)
- role (apiVersion rbac.authorization.k8s.io/v1beta1)
- roleBinding (apiVersion rbac.authorization.k8s.io/v1beta1)
- secret (apiVersion v1)
- serviceAccount (apiVersion v1)

<h4>Notes:</h4>

- If the cluster is not listed under "clusters" but is under "replicas" the resources will not be created on the cluster.
- If the manifest for a deployment or statefulSet specifies "replicas" it will be overwritted with the number listed in the `Gheefile` "replicas".
