---
subcategory: "Firebase"
description: |-
  Fetches the details of a Firebase Extensions instance.
---

# google_firebase_extensions_instance

Fetches the details of a Firebase Extensions instance. For more details, see the [API documentation](https://firebase.google.com/products/extensions).

~> **Note:** This resource requires the [Firebase API](https://console.cloud.google.com/apis/library/firebase.googleapis.com/) to be enabled.

## Example Usage

```hcl
data "google_firebase_extensions_instance" "default" {
  provider = google-beta
  instance_id = "my-extension-instance"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` -
  (Required)
  The ID to use for the Extension Instance, which will become the final component of the instance's name.

* `project` - 
  (Optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_firebase_extensions_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/firebase_extensions_instance) resource for details of all the available attributes.
