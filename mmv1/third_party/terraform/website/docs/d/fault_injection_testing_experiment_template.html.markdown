---
subcategory: "Fault Injection Testing"
description: |-
  Fetches the details of a Fault Injection Testing Experiment Template.
---

# google_fault_injection_testing_experiment_template

Use this data source to get information about a Fault Injection Testing Experiment Template. For more details, see the [API documentation](https://cloud.google.com/fault-injection/docs/reference/rest/v1alpha/projects.locations.experimentTemplates).

## Example Usage

```hcl
data "google_fault_injection_testing_experiment_template" "default" {
  experiment_template_id = "my-experiment-template"
  location               = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `experiment_template_id` -
  (Required)
  The ID of the Experiment Template.

* `location` -
  (Required)
  The location in which the Experiment Template resides.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_fault_injection_testing_experiment_template](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/fault_injection_testing_experiment_template) resource for details of all the available attributes.
