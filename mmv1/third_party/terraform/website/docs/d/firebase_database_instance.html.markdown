---
subcategory: "Firebase"
description: |-
  Fetches the details of a Firebase Realtime Database instance.
---

# google_firebase_database_instance

Fetches the details of a Firebase Realtime Database instance. For more details, see the [API documentation](https://firebase.google.com/docs/reference/rest/database/database-management/rest).

~> **Note:** This resource requires the [Firebase API](https://console.cloud.google.com/apis/library/firebase.googleapis.com/) to be enabled.

## Example Usage

```hcl
data "google_firebase_database_instance" "default" {
  provider = google-beta
  region   = "us-central1"
  instance_id = "my-database-instance"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` -
  (Required)
  The globally unique identifier of the Firebase Realtime Database instance.

* `region` -
  (Required)
  A reference to the region where the Firebase Realtime database resides.
  Check all [available regions](https://firebase.google.com/docs/projects/locations#rtdb-locations)

* `project` - 
  (Optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_firebase_database_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/firebase_database_instance) resource for details of all the available attributes.
