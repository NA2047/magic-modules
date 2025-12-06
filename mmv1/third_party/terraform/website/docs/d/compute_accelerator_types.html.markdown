---
subcategory: "Compute Engine"
description: |-
  Provides a list of available Google Compute accelerator types
---

# google_compute_accelerator_types

Provides access to available Google Compute accelerator types in a zone or across all zones for a given project.
See more about [accelerator types](https://cloud.google.com/compute/docs/gpus) in the upstream docs.

To get more information about accelerator types, see:

* [API Documentation](https://cloud.google.com/compute/docs/reference/rest/v1/acceleratorTypes/aggregatedList)
* [GPU Documentation](https://cloud.google.com/compute/docs/gpus)

## Example Usage - List all accelerator types

```hcl
data "google_compute_accelerator_types" "available" {}

output "accelerator_types" {
  value = data.google_compute_accelerator_types.available.types
}
```

## Example Usage - Filter by zone

```
data "google_compute_accelerator_types" "zone_specific" {
  zone = "us-central1-a"
}

resource "google_compute_instance" "gpu_instance" {
  name         = "gpu-instance"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  guest_accelerator {
    type  = data.google_compute_accelerator_types.zone_specific.types[0].name
    count = 1
  }

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}
```

## Example Usage - Filter by accelerator type

```
data "google_compute_accelerator_types" "nvidia_gpus" {
  filter = "name:nvidia*"
  zone   = "us-central1-a"
}

output "nvidia_accelerators" {
  value = [for accelerator in data.google_compute_accelerator_types.nvidia_gpus.types : {
    name = accelerator.name
    max_cards = accelerator.maximum_cards_per_instance
  }]
}
```

## Argument Reference

The following arguments are supported:

* `filter` (Optional) - A filter expression that filters accelerator types listed in the response. The expression must specify the field name, an operator, and the value that you want to use for filtering.

* `zone` (Optional) - The zone to list accelerator types for. If not specified, accelerator types from all zones will be returned.

* `project` (Optional) - The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `types` - The list of accelerator types. Structure is [documented below](#nested_types).

<a name="nested_types"></a>The `types` block supports:

* `name` - The name of the accelerator type.

* `description` - An optional textual description of the resource.

* `zone` - The name of the zone where the accelerator type resides.

* `maximum_cards_per_instance` - Maximum number of accelerator cards allowed per instance.

* `self_link` - The server-defined URL for the resource.

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `deprecated` - The deprecation status associated with this accelerator type. Structure is [documented below](#nested_deprecated).

<a name="nested_deprecated"></a>The `deprecated` block supports:

* `state` - The deprecation state of this resource. This can be `ACTIVE`, `DEPRECATED`, `OBSOLETE`, or `DELETED`.

* `replacement` - The URL of the suggested replacement for a deprecated resource.

* `deprecated` - An optional RFC3339 timestamp on or after which the state of this resource is intended to change to `DEPRECATED`.

* `obsolete` - An optional RFC3339 timestamp on or after which the state of this resource is intended to change to `OBSOLETE`.

* `deleted` - An optional RFC3339 timestamp on or after which the state of this resource is intended to change to `DELETED`.
