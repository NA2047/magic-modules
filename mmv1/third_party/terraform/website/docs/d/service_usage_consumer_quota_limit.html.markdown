---
subcategory: "Service Usage"
description: |-
  Get information about a Google Cloud Service Usage Consumer Quota Limit.
---

# google_service_usage_consumer_quota_limit

Get information about a Google Cloud Service Usage Consumer Quota Limit. Consumer quota settings for a quota limit.

To get more information about ConsumerQuotaLimit, see:

* [API documentation](https://cloud.google.com/service-usage/docs/reference/rest/v1beta1/services.consumerQuotaMetrics.limits)
* How-to Guides
    * [Service Usage API](https://cloud.google.com/service-usage/docs)

## Example Usage

```hcl
data "google_service_usage_consumer_quota_limit" "limit" {
  service    = "compute.googleapis.com"
  metric     = "compute.googleapis.com%2Fcpus"
  limit_name = "%2Fproject%2Fregion"
}
```

## Example Usage - With Project

```hcl
data "google_service_usage_consumer_quota_limit" "limit" {
  project    = "my-project-id"
  service    = "compute.googleapis.com"
  metric     = "compute.googleapis.com%2Fcpus"
  limit_name = "%2Fproject%2Fregion"
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The service name for which the quota limit is being retrieved.
  For example: `compute.googleapis.com`

* `metric` - (Required) The metric name for which the quota limit is being retrieved.
  For example: `compute.googleapis.com%2Fcpus`

* `limit_name` - (Required) The limit name for which the quota limit is being retrieved.
  For example: `%2Fproject%2Fregion`

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` - The resource name of the quota limit.
  An example name would be: `projects/123/services/compute.googleapis.com/consumerQuotaMetrics/compute.googleapis.com%2Fcpus/limits/%2Fproject%2Fregion`

* `metric` - The name of the parent metric of this limit.
  An example name would be: `compute.googleapis.com/cpus`

* `unit` - The limit unit.
  An example unit would be `1/{project}/{region}`

* `is_precise` - Whether this limit is precise or imprecise.

* `allows_admin_overrides` - Whether admin overrides are allowed on this limit.

* `quota_buckets` - Summary of the enforced quota buckets, organized by quota dimension, ordered from least specific to most specific.
  Structure is [documented below](#nested_quota_buckets).

* `supported_locations` - List of all supported locations. This field is present only if the limit has a {region} or {zone} dimension.

<a name="nested_quota_buckets"></a>The `quota_buckets` block contains:

* `effective_limit` - The effective limit of this quota bucket. Equal to defaultLimit if there are no overrides.

* `default_limit` - The default limit of this quota bucket, as specified by the service configuration.

* `dimensions` - The dimensions of this quota bucket.

* `producer_override` - Producer override on this quota bucket.
  Structure is [documented below](#nested_override).

* `consumer_override` - Consumer override on this quota bucket.
  Structure is [documented below](#nested_override).

* `admin_override` - Admin override on this quota bucket.
  Structure is [documented below](#nested_override).

* `producer_quota_policy` - Producer policy inherited from the closest ancestor of the current consumer.
  Structure is [documented below](#nested_producer_quota_policy).

* `rollout_info` - Rollout information of this quota bucket.
  Structure is [documented below](#nested_rollout_info).

<a name="nested_override"></a>The `producer_override`, `consumer_override`, and `admin_override` blocks contain:

* `name` - The resource name of the override.

* `override_value` - The overriding quota limit value.

* `dimensions` - The dimensions of the override.

<a name="nested_producer_quota_policy"></a>The `producer_quota_policy` block contains:

* `name` - The resource name of the policy.

* `policy_value` - The quota policy value.

* `dimensions` - The dimensions of the policy.

* `metric` - The name of the metric to which this policy applies.

* `unit` - The limit unit of the limit to which this policy applies.

* `container` - The cloud resource container at which the quota policy is created.

<a name="nested_rollout_info"></a>The `rollout_info` block contains:

* `default_limit_ongoing_rollout` - Whether there is an ongoing rollout for the default limit or not.
