---
subcategory: "Fault Injection Testing"
description: |-
  Fetches the details of a Fault Injection Testing Experiment.
---

# google_fault_injection_testing_experiment

Use this data source to get information about a Fault Injection Testing Experiment. For more details, see the [API documentation](https://cloud.google.com/fault-injection/docs/reference/rest/v1alpha/projects.locations.experiments).

## Example Usage

```hcl
data "google_fault_injection_testing_experiment" "default" {
  experiment_id = "my-experiment"
  location      = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `experiment_id` -
  (Required)
  The ID of the Experiment.

* `location` -
  (Required)
  The location in which the Experiment resides.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_fault_injection_testing_experiment](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/fault_injection_testing_experiment) resource for details of all the available attributes.
