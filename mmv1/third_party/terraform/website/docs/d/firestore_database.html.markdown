---
subcategory: "Firestore"
description: |-
  Fetches the details of a Firestore Database.
---

# google_firestore_database

Use this data source to get information about a Firestore Database. For more details, see the [API documentation](https://cloud.google.com/firestore/docs/reference/rest/v1/projects.databases).

## Example Usage

```hcl
data "google_firestore_database" "default" {
  name = "(default)"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The ID to use for the database, which will become the final component of the database's resource name. This value should be 4-63 characters. Valid characters are /[a-z][0-9]-/ with first character a letter and the last a letter or a number. Must not be UUID-like /[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}/. "(default)" database id is also valid.

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_firestore_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/firestore_database) resource for details of all the available attributes.
