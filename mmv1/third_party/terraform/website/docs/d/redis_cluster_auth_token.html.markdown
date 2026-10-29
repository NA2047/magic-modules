---
subcategory: "Redis"
description: |-
  Get information about an Auth Token for a Google Cloud Redis Cluster Token Auth User.
---

# google_redis_cluster_auth_token

Get information about an Auth Token associated with a Google Cloud Redis Cluster Token Auth User. For more details refer to the [API documentation](https://cloud.google.com/memorystore/docs/cluster/reference/rest/v1/projects.locations.clusters).

## Example Usage

```hcl
data "google_redis_cluster_auth_token" "qa" {
  token_auth_user = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
  token_id        = "version-1"
}
```

## Argument Reference

The following arguments are supported:

* `token_auth_user` -
  (Required)
  The full resource name of the parent Token Auth User (`projects/{project}/locations/{region}/clusters/{cluster}/tokenAuthUsers/{user_id}`) or short path (`{cluster}/tokenAuthUsers/{user_id}`).

* `token_id` -
  (Optional)
  The ID of the auth token (e.g. `version-1`). If not provided, the latest active token version for the user is returned.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `region` -
  (Optional)
  The name of the region of the Redis cluster. If not provided, the provider region is used.

## Attributes Reference

See [`google_redis_cluster_auth_token`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_cluster_auth_token) resource for details of all available attributes.
