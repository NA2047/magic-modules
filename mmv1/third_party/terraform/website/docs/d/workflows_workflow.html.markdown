---
subcategory: "Workflows"
description: |-
  Fetches the details of a Workflows Workflow.
---

# google_workflows_workflow

Use this data source to get information about a Workflows Workflow. For more details, see the [API documentation](https://cloud.google.com/workflows/docs/reference/rest/v1/projects.locations.workflows).

## Example Usage

```hcl
data "google_workflows_workflow" "default" {
  name   = "my-workflow"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Name of the Workflow.

* `region` -
  (Required)
  The region of the workflow.

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_workflows_workflow](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/workflows_workflow) resource for details of all the available attributes.
