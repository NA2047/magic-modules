---
subcategory: "Kubernetes Engine"
description: |-
  Fetches the details of a node pool in a Google Kubernetes Engine (GKE) cluster.
---

# google_container_node_pool

Fetches the details of a node pool in a Google Kubernetes Engine (GKE) cluster. For more details, see the [API documentation](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters.nodePools).

## Example Usage

```hcl
data "google_container_node_pool" "my_node_pool" {
  name       = "my-node-pool"
  location   = "us-central1-a"
  cluster    = "my-gke-cluster"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The name of the node pool.

* `cluster` -
  (Required)
  The cluster to create the node pool for. Cluster must be present in location provided for zonal clusters.

* `location` -
  (Optional)
  The location (region or zone) of the cluster.

* `project` - 
  (Optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_container_node_pool](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_node_pool) resource for details of all the available attributes.
