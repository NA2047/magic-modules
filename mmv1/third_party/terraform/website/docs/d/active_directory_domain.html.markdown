---
subcategory: "Active Directory"
description: |-
  Fetches the details of a Managed Microsoft Active Directory domain.
---

# google_active_directory_domain

Use this data source to get information about a Managed Microsoft Active Directory domain. For more details, see the [API documentation](https://cloud.google.com/managed-microsoft-ad/reference/rest/v1/projects.locations.global.domains).

## Example Usage

```hcl
data "google_active_directory_domain" "default" {
  domain_name = "mydomain.org.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` -
  (Required)
  The fully qualified domain name. e.g. mydomain.myorganization.com.

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_active_directory_domain](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/active_directory_domain) resource for details of all the available attributes.
