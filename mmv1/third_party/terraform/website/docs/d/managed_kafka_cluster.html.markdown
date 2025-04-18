---
subcategory: "Managed Kafka"
description: |-
  Fetches the details of a Managed Service for Apache Kafka cluster.
---

# google_managed_kafka_cluster

Use this data source to get information about a Managed Service for Apache Kafka cluster. For more details, see the [API documentation](https://cloud.google.com/managed-kafka/docs/reference/rest/v1/projects.locations.clusters).

## Example Usage

```hcl
data "google_managed_kafka_cluster" "default" {
  cluster_id = "my-cluster"
  location   = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -
  (Required)
  The ID to use for the cluster, which will become the final component of the cluster's name. The ID must be 1-63 characters long, and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` to comply with RFC 1035.

* `location` -
  (Required)
  ID of the location of the Kafka resource. See https://cloud.google.com/managed-kafka/docs/locations for a list of supported locations.

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_managed_kafka_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/managed_kafka_cluster) resource for details of all the available attributes.
