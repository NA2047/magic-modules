---
subcategory: "Filestore"
description: |-
  Fetches the details of a Google Cloud Filestore backup.
---

# google_filestore_backup

Fetches the details of a Google Cloud Filestore backup. For more details, see the [API documentation](https://cloud.google.com/filestore/docs/reference/rest/v1/projects.locations.instances.backups).

## Example Usage

```hcl
resource "google_filestore_instance" "instance" {
  name     = "test-instance"
  location = "us-central1-b"
  tier     = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}

resource "google_filestore_backup" "backup" {
  name     = "test-backup"
  location = "us-central1"
  source_instance = google_filestore_instance.instance.id
  source_file_share = "share1"
}

data "google_filestore_backup" "backup" {
  name     = google_filestore_backup.backup.name
  location = google_filestore_backup.backup.location
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The resource name of the backup.

* `location` -
  (Required)
  The name of the location of the backup. This can be a region for ENTERPRISE tier instances.

* `project` - 
  (Optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_filestore_backup](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/filestore_backup) resource for details of all the available attributes.
