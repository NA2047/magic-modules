---
name: add-data-source-workflow
description: Workflow for checking if a corresponding Terraform data source exists and adding or updating a data source (resource-backed or resourceless) in Magic Modules. Invoked directly or from add-fields-workflow and new-resource-workflow.
---

# Workflow: Add or Verify a Data Source

This skill guides you through checking whether a corresponding Terraform data source already exists for a resource, and adding or updating the handwritten data source in Magic Modules (`mmv1/third_party/terraform/services/<product>/`), along with its acceptance tests and documentation.

> **When this skill is invoked:**
> - **Directly:** When the user asks to create or update a data source.
> - **From `new-resource-workflow`:** After creating a new resource, to check if a corresponding data source exists and create one if needed.
> - **From `add-fields-workflow`:** When adding fields to an existing resource, to check if a corresponding data source exists and ensure either the existing data source is updated (if needed) or a missing data source is identified/created.

---

## Core Rules

1. **Always Check Existence First**:
   - Before creating any files, check if the corresponding data source already exists in `mmv1/third_party/terraform/services/<product>/` or via MMv1 YAML generation.
2. **`project`, `location`, `region`, and `zone` MUST ALWAYS BE OPTIONAL**:
   - Never mark `project`, `location`, `region`, or `zone` as `Required` on a data source.
   - When a user omits `project`, `location`, `region`, or `zone` in their HCL configuration, the data source **must** fall back to the provider-level default (`tpgresource.GetProject(d, config)`, `tpgresource.GetLocation(d, config)`, `tpgresource.GetRegion(d, config)`, or `tpgresource.GetZone(d, config)`), or to the service's default location (such as `"global"`) when applicable.
3. **Always Populate Data Source Labels (`tpgresource.SetDataSourceLabels(d)`)**:
   - Whenever the underlying resource has `labels` (`KeyValueLabels`), `resource<Product><Resource>Read` only sets `effective_labels` and `terraform_labels`. You **must** call `tpgresource.SetDataSourceLabels(d)` after `resource<Product><Resource>Read(d, meta)`:
     ```go
     if err := tpgresource.SetDataSourceLabels(d); err != nil {
     	return err
     }
     ```
     (Similarly, call `tpgresource.SetDataSourceAnnotations(d)` if the resource has `annotations`.)
4. **Resource-Backed Data Sources — Always Set Unflattened / Defaulted Fields in `Read`**:
   - Delegating to `resource<Product><Resource>Read(d, meta)` is necessary for resource-backed data sources, but **it is rarely sufficient on its own**. You must inspect the underlying resource's YAML/Go definition to determine which additional fields must be explicitly resolved and set via `d.Set(...)` before and after calling `resource<Product><Resource>Read(d, meta)`.
5. **Always Handle `virtual_fields` When Implemented on the Resource**:
   - Check if the resource defines `virtual_fields` in its MMv1 YAML (or client-side lifecycle fields like `deletion_protection`, `deletion_policy`, `force_destroy`). Because client-side virtual fields are never returned by the GCP API—and MMv1's `resource<Name>Read` automatically sets any `virtual_fields` with a `default_value` into state—you must explicitly handle them in `dataSource<Name>Read` and/or the acceptance test (`Step 4`).
6. **Resourceless Data Sources — Implement a Custom `Read` Function**:
   - When a data source queries a read-only or system-managed GCP object that has **no corresponding Terraform resource** (for example, `google_compute_reservation_block` and `google_compute_reservation_sub_block`), there is no `Resource<Name>()` schema or `resource<Name>Read(d, meta)` function to call. You must define the full `map[string]*schema.Schema` and implement the entire `Read` lifecycle yourself using `transport_tpg.SendRequest`, manual response flattening (`d.Set(...)`), and `d.SetId(...)`.

---

## Step 0: Check if the Corresponding Data Source Already Exists

Before writing any new data source code, check whether a data source for `google_<product>_<resource>` already exists:

1. **Search for an existing handwritten data source**:
   - Check `mmv1/third_party/terraform/services/<product>/` for `data_source_*<resource>*.go` or `data_source_*<resource>*.go.tmpl` (using `find_by_name`), or search for `"google_<product>_<resource>"` with `grep_search` in `mmv1/third_party/terraform/services/<product>/`.
   - Also check `mmv1/third_party/terraform/website/docs/d/` for `<product>_<resource>.html.markdown`.
2. **Decision Matrix Based on Existence**:
   - **Case A — The data source DOES NOT exist**:
     - If invoked directly or from `new-resource-workflow` / `add-fields-workflow`, proceed to **Step 1** below to create the data source, acceptance test, and documentation (or ask the user if they want the missing data source created as part of the current change if scope is restricted).
   - **Case B — The data source ALREADY exists**:
     - **If invoked from `add-fields-workflow`**:
       - Check how the existing data source is implemented:
         1. If it uses `tpgresource.DatasourceSchemaFromResourceSchema(Resource<Name>().Schema)`, the new resource fields are automatically included in the data source schema. **However**, inspect the newly added fields against **Step 3A** and **Step 4** below:
            - Did you add `labels` or `annotations` to the resource? Ensure `tpgresource.SetDataSourceLabels(d)` / `tpgresource.SetDataSourceAnnotations(d)` is called in `dataSource<Name>Read`.
            - Does the new field have `ignore_read: true` on the resource with an `effective_*` counterpart? Explicitly copy it with `d.Set(...)` in `dataSource<Name>Read` (Step 3A.5).
            - Is the new field added under `virtual_fields` (e.g., `deletion_protection`, `deletion_policy`, or an `output: true` virtual field)? Follow **Step 4** to ensure `dataSource<Name>Read` and the data source test handle it properly.
         2. If the existing data source defines its own manual `map[string]*schema.Schema` and custom `Read` (Step 3B), add the new field(s) to the data source schema, flatten them via `d.Set(...)` in `Read`, and document them in `website/docs/d/<product>_<resource>.html.markdown`.
       - Also check if `location`, `region`, `zone`, and `project` on the existing data source are properly marked `Optional` and resolved via provider defaults.
       - Run the existing data source acceptance test (`TestAccDataSource...` / `TestAcc...Datasource...`) in addition to the resource tests to confirm `CheckDataSourceStateMatchesResourceState` passes with the new fields.
     - **If invoked directly to add a data source**:
       - Inform the user that the data source already exists at the found path, and check if they want to audit/update its optional `location`/`region`/`zone`/`project` handling, `virtual_fields`, missing fields, or tests.

---

## Step 1: Analyze the Target Resource (or API)

Before writing code, check whether a corresponding Terraform resource exists in `mmv1/products/<product>/<Resource>.yaml` (or `mmv1/third_party/terraform/services/<product>/`).

- **If a corresponding resource exists**, inspect its YAML and generated `resource_<product>_<resource>.go` for:
  1. **URL & ID Structure (`base_url`, `self_link`, `id_format`, `import_format`)**:
     - What template variables (`{{project}}`, `{{location}}`, `{{region}}`, `{{zone}}`, `{{name}}`, `{{instance_id}}`, etc.) are required by `tpgresource.ReplaceVars` in `resource<Name>Read`?
  2. **`url_param_only: true` Fields**:
     - Which parameters have `url_param_only: true`? Generated `Resource<Name>Flatten` functions **do not flatten `url_param_only: true` fields from the API response**.
  3. **Parent References vs. User-Friendly Inputs**:
     - Does the resource require a full parent self-link (e.g., `cluster` in `google_alloydb_instance`) instead of separate `cluster_id`, `location`, and `project` fields?
     - Does the resource name a zonal field `location` (e.g., `google_lustre_instance`) where users of a data source might expect `zone` (or `location`)?
  4. **Labels & Annotations**:
     - Does the resource define `labels` (`KeyValueLabels`) or `annotations` (`KeyValueAnnotations`)?
  5. **`virtual_fields` & Universal `deletion_policy`**:
     - Does the resource YAML define a `virtual_fields:` block (e.g., `deletion_protection`, `deletion_policy`, `force_destroy`, or computed `output: true` virtual fields like `block_names`), or have universal `deletion_policy` enabled? (See **Step 4**.)
  6. **`ignore_read: true` Properties**:
     - Are there properties marked `ignore_read: true` on the resource (e.g., `reserved_ip_range` on `google_redis_instance`) that have a corresponding `effective_*` output field?
  7. **`nested_query` or Custom `decoder` / `pre_read` / `post_read`**:
     - Does the resource's `Read` function or custom decoder read attributes from `d.Get(...)` (such as `d.Get("name")`) to filter a list response or decode the payload?
- **If NO corresponding resource exists (Resourceless Data Source)**:
  - Inspect the GCP REST API documentation / Discovery document for the `GET` (and any child `LIST`) endpoint URL structure, response envelope (e.g., whether fields are nested under a `"resource"` wrapper key), and all camelCase JSON response properties.
- **Beta vs. GA Availability**:
  - If the resource or API is `min_version: beta`, use `.go.tmpl` files with `{{ if ne $.TargetVersionName "ga" -}}` guards so it is only generated into `terraform-provider-google-beta`.

---

## Step 2: Define the Data Source Schema

Create `mmv1/third_party/terraform/services/<product>/data_source_<product>_<resource>.go` (or `data_source_google_<product>_<resource>.go` following the service directory's naming convention).

### Pattern A: Standard Resource-Backed Data Source
When the resource schema already contains the identifier and locational/project fields:

```go
func DataSource<Product><Resource>() *schema.Resource {
	rs := Resource<Product><Resource>().Schema
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(rs)

	// 1. ONLY mark the resource's primary identifier(s) (and non-defaultable parent IDs) as Required:
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "<id_field>")

	// 2. ALWAYS mark project and location/region/zone as Optional (NEVER Required):
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location", "project")

	return &schema.Resource{
		Read:   dataSource<Product><Resource>Read,
		Schema: dsSchema,
	}
}
```

### Pattern B: Merging Custom Schema Fields (`tpgresource.MergeSchemas`)
When the underlying resource schema **lacks** `project`, `location`, `region`, `zone`, or a clean parent ID field, define a custom `map[string]*schema.Schema` and merge it into `dsSchema` using `tpgresource.MergeSchemas`:

- **Example 1 — Resource only has a parent self-link `cluster` (`google_alloydb_instance`)**:
  The resource schema only has `cluster` (`projects/{project}/locations/{location}/clusters/{cluster_id}`) and lacks `project`, `location`, and `cluster_id`. Merge `cluster_id` (`Required: true`), `project` (`Optional: true`), and `location` (`Optional: true`) into `dsSchema`:
  ```go
  dsSchema_custom := map[string]*schema.Schema{
  	"project": {
  		Type:        schema.TypeString,
  		Optional:    true,
  		Description: "Project ID of the project.",
  	},
  	"location": {
  		Type:        schema.TypeString,
  		Optional:    true,
  		Description: "The canonical id of the location.If it is not provided, the provider project is used.",
  	},
  	"cluster_id": {
  		Type:        schema.TypeString,
  		Required:    true,
  		Description: "Identifies the alloydb cluster. Must be in the format 'projects/{project}/locations/{location}/clusters/{cluster_id}'",
  	},
  }
  dsSchema = tpgresource.MergeSchemas(dsSchema_custom, dsSchema)
  ```

- **Example 2 — Resource uses `location` for a zone, data source exposes optional `zone` (`google_lustre_instance`)**:
  ```go
  dsSchema_zone := map[string]*schema.Schema{
  	"zone": {
  		Type:        schema.TypeString,
  		Optional:    true,
  		Description: "The canonical id of the zone. If it is not provided, the provider zone is used.",
  	},
  }
  dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(rs)
  dsSchema = tpgresource.MergeSchemas(dsSchema_zone, dsSchema)
  tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance_id")
  tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
  ```

### Pattern C: Resourceless Data Source (No Corresponding Resource)
When there is no corresponding Terraform resource (e.g., `google_compute_reservation_block`, `google_compute_reservation_sub_block`), define the entire `map[string]*schema.Schema` directly:
- Mark only the object's own name/ID and required parent resource names (e.g., `name`, `reservation`, `reservation_block`) as `Required: true`.
- Always mark `zone` / `region` / `location` as `Optional: true`, and `project` as `Optional: true, Computed: true`.
- Mark all output attributes returned by the API as `Computed: true`. Represent nested objects using `Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{...}}`.
- Avoid reserved Terraform attribute names (`id`, `count`): map API `id` to `resource_id` and API `count` to a descriptive name like `block_count` or `sub_block_count`.

---

## Step 3A: Implement `Read` for a Resource-Backed Data Source

### Standard `Read` Template (Resource-Backed)
Use this baseline structure for `dataSource<Product><Resource>Read` (e.g., `data_source_memorystore_instance.go`), and then apply the scenario checks below and **Step 4 (`virtual_fields`)**:

```go
func dataSource<Product><Resource>Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	// Set location BEFORE ReplaceVars and resourceRead so {{location}} is resolved and saved in state:
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/<collection>/{{<id_field>}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resource<Product><Resource>Read(d, meta)
	if err != nil {
		return err
	}

	// Always set labels when the resource has KeyValueLabels:
	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
```

When wrapping `resource<Product><Resource>Read(d, meta)`, carefully check each of the following scenarios to determine which fields must be resolved and set on `d`:

### 1. Resolving and Setting `location`, `region`, `zone`, and `project` BEFORE `resource<Name>Read`

#### Why this is required:
1. **`tpgresource.ReplaceVars` does NOT auto-resolve `{{location}}` from the provider config**:
   - In `mmv1/third_party/terraform/tpgresource/utils.go`, `BuildReplacementFunc` automatically calls `GetProject(d, config)`, `GetRegion(d, config)`, and `GetZone(d, config)` for `{{project}}`, `{{region}}`, and `{{zone}}`.
   - It **does not** call `GetLocation(d, config)` for `{{location}}`—it only checks `d.GetOkExists("location")`.
   - Therefore, if `location` is optional and omitted by the user in HCL, calling `ReplaceVars` or `resource<Name>Read(d, meta)` without first setting `location` in `d` will produce an empty URL segment (`.../locations//...`) and fail!
2. **`url_param_only: true` fields are NOT set in state by `resource<Name>Read`**:
   - In MMv1's `resource.go.tmpl`, `resource<Name>Read` only sets `region` or `zone` in state if `HasRegion()` / `HasZone()` is true (which requires `ignore_read: true` on the YAML parameter), and **never** automatically sets `location`.
   - If `location`, `region`, or `zone` has `url_param_only: true` (e.g., `location` in `google_alloydb_cluster`, `google_memorystore_instance`, `google_memorystore_acl_policy`, `google_redis_cluster_acl_policy`; `region` in `google_memcache_instance`), `resource<Name>Read` will leave it unset (`""`) in the data source state unless you call `d.Set(...)` explicitly!

#### How to implement:
- **When the resource uses `location` (regional or zonal fallback via provider config)**:
  ```go
  location, err := tpgresource.GetLocation(d, config)
  if err != nil {
  	return err
  }
  if err := d.Set("location", location); err != nil {
  	return fmt.Errorf("Error setting location: %s", err)
  }
  ```
- **When the resource uses `location` with a service-level default like `"global"` (e.g., `google_certificate_manager_dns_authorization`)**:
  ```go
  location, ok := d.GetOk("location")
  if !ok {
  	location = "global"
  }
  if err := d.Set("location", location.(string)); err != nil {
  	return fmt.Errorf("Error setting location: %s", err)
  }
  ```
- **When the resource uses `region`**:
  ```go
  region, err := tpgresource.GetRegion(d, config)
  if err != nil {
  	return err
  }
  if err := d.Set("region", region); err != nil {
  	return fmt.Errorf("Error setting region: %s", err)
  }
  ```
- **When the resource uses `zone`**:
  ```go
  zone, err := tpgresource.GetZone(d, config)
  if err != nil {
  	return err
  }
  if err := d.Set("zone", zone); err != nil {
  	return fmt.Errorf("Error setting zone: %s", err)
  }
  ```
- **Always ensure `project` can fall back to the provider default**:
  If `resource<Name>Read` does not already set `project` (for instance, in resources without `{{project}}` in `base_url`), resolve and set it explicitly:
  ```go
  project, err := tpgresource.GetProject(d, config)
  if err != nil {
  	return err
  }
  if err := d.Set("project", project); err != nil {
  	return fmt.Errorf("Error setting project: %s", err)
  }
  ```

---

### 2. Setting Underlying Resource `url_param_only` Fields from Merged/Custom Inputs

#### Why this is required:
When you use `tpgresource.MergeSchemas` to expose friendlier inputs on the data source (or when a parent reference can be passed as a self-link), `resource<Name>Read` still expects the **resource's** `url_param_only` field names to be populated in `d` so `ReplaceVars` can build the API URL and `CheckDataSourceStateMatchesResourceState` sees matching state attributes.

#### Common Cases:
- **Constructing a parent self-link (`google_alloydb_instance`)**:
  The data source accepts `cluster_id` (Required), `project` (Optional), and `location` (Optional), while `resourceAlloydbInstanceRead` requires `cluster` (`projects/{project}/locations/{location}/clusters/{cluster_id}`):
  ```go
  cluster_id := d.Get("cluster_id").(string)
  location, err := tpgresource.GetLocation(d, config)
  if err != nil {
  	return err
  }
  project, err := tpgresource.GetProject(d, config)
  if err != nil {
  	return err
  }
  cluster := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, cluster_id)
  if err := d.Set("cluster", cluster); err != nil {
  	return fmt.Errorf("Error setting cluster: %s", err)
  }
  ```
- **Mapping `zone` to `location` (`google_lustre_instance`)**:
  The data source accepts optional `zone`, while `resourceLustreInstanceRead` expects `location`:
  ```go
  zone, err := tpgresource.GetZone(d, config)
  if err != nil {
  	return err
  }
  if err := d.Set("location", zone); err != nil {
  	return fmt.Errorf("Error setting location: %s", err)
  }
  ```
- **Shortening a parent self-link (`google_memorystore_token_auth_user`, `google_redis_cluster_token_auth_user`)**:
  When a parent field (`instance` or `cluster`) is used inside a URL template like `projects/{{project}}/locations/{{location}}/instances/{{instance}}/...`, normalize it first with `tpgresource.GetResourceNameFromSelfLink` so passing a full self-link does not break URL construction:
  ```go
  instance := tpgresource.GetResourceNameFromSelfLink(d.Get("instance").(string))
  if err := d.Set("instance", instance); err != nil {
  	return fmt.Errorf("Error setting instance: %s", err)
  }
  ```

---

### 3. Setting `name` (or Composite IDs) Before `resource<Name>Read` for `nested_query` or Custom Decoders

#### Why this is required:
Some MMv1 resources use `nested_query` (fetching a list/parent object and filtering it via `flattenNested<Name>`) or custom decoders (`templates/terraform/decoders/...`) that read `d.Get("name")` during `resource<Name>Read`. If the data source identifies the object via separate ID fields (e.g., `token_auth_user` and `token_id` on `google_memorystore_auth_token` / `google_redis_cluster_auth_token`) rather than `name`, `d.Get("name")` will be empty during `resource<Name>Read` unless you set it beforehand.

#### How to implement:
```go
id := fmt.Sprintf("%s/authTokens/%s", tokenAuthUser, tokenId)
d.SetId(id)
if err := d.Set("name", id); err != nil {
	return fmt.Errorf("Error setting name: %s", err)
}
err := resourceMemorystoreAuthTokenRead(d, meta)
```

---

### 4. Populating `labels` and `annotations` After `resource<Name>Read`

#### Why this is required:
In the Google provider, managed resources separate user-configured `labels` / `annotations` from API-returned `effective_labels` / `effective_annotations`. Consequently, `resource<Name>Read` only populates `effective_labels` (and `terraform_labels`) and leaves `labels` / `annotations` empty unless already in state. On a data source, `labels` and `annotations` are `Computed: true` and must reflect all labels/annotations on the resource.

#### How to implement:
After calling `resource<Name>Read(d, meta)`:
- If the resource has `labels` (`KeyValueLabels`):
  ```go
  if err := tpgresource.SetDataSourceLabels(d); err != nil {
  	return err
  }
  ```
- If the resource has `annotations` (`KeyValueAnnotations`):
  ```go
  if err := tpgresource.SetDataSourceAnnotations(d); err != nil {
  	return err
  }
  ```

---

### 5. Populating Fields Marked `ignore_read: true` on the Resource

#### Why this is required:
When a resource property has `ignore_read: true` in MMv1 YAML (often used to prevent permadiffs when the API populates a server-side default, such as `reserved_ip_range` on `google_redis_instance`), `Resource<Name>Flatten` skips reading that field from the API response. When the data source calls `resource<Name>Read`, that field will remain `null`/`""` in the data source even though the value is available in an `effective_*` attribute!

#### How to implement:
Inspect the resource YAML for `ignore_read: true` properties. If the value is captured by an `effective_*` property on the resource (e.g., `effective_reserved_ip_range`), copy it into the field after `resource<Name>Read(d, meta)`:
```go
if err := d.Set("reserved_ip_range", d.Get("effective_reserved_ip_range")); err != nil {
	return fmt.Errorf("Error setting reserved_ip_range: %s", err)
}
```

---

### 6. Verifying the Resource Exists (`d.Id() == ""`)

#### Why this is required:
When the GCP API returns `404 Not Found`, `transport_tpg.HandleNotFoundError` inside `resource<Name>Read` clears the resource ID (`d.SetId("")`) and returns `nil`. While returning `nil` on 404 is correct for a managed resource refresh, a data source **must** return an error if the target resource does not exist.

#### How to implement:
Always set `d.SetId(id)` before calling `resource<Name>Read(d, meta)` and check `if d.Id() == ""` before returning:
```go
if d.Id() == "" {
	return fmt.Errorf("%s not found", id)
}
return nil
```

---

## Step 3B: Implement a Custom `Read` When No Corresponding Resource Exists (Resourceless Data Source)

When there is no corresponding Terraform resource (e.g., `google_compute_reservation_block` and `google_compute_reservation_sub_block`), the data source must implement its own `Read` function from scratch. Follow this structure:

### 1. Resolve Provider Config, `userAgent`, and Default `project` / `zone` / `region` / `location`
Even without a backing resource, `project` and `zone`/`region`/`location` **must** be resolved using `tpgresource.GetProject` and `tpgresource.GetZone` / `GetRegion` / `GetLocation` so omitting them in HCL falls back to the provider defaults:

```go
func dataSourceGoogleComputeReservationBlockRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
```

### 2. Construct the Request URL and Call `transport_tpg.SendRequest`
Build the REST endpoint URL from the resolved `project`, `zone`/`region`/`location`, and required arguments, then execute a `"GET"` request using `transport_tpg.SendRequest`:

```go
	blockName := d.Get("name").(string)
	reservation := d.Get("reservation").(string)

	// Note: Some deeply nested Compute APIs require URL-encoded slashes (%%2F) in parent path segments,
	// e.g. ".../reservations%%2F%s%%2FreservationBlocks%%2F%s/reservationSubBlocks/%s"
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/reservations/%s/reservationBlocks/%s", project, zone, reservation, blockName)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading ReservationBlock: %s", err)
	}
	if res == nil {
		return fmt.Errorf("ReservationBlock %s not found", blockName)
	}
```

### 3. Unwrap Nested API Response Envelopes (If Applicable)
Some GCP endpoints (such as Compute `reservationBlocks.get` and `reservationSubBlocks.get`) wrap the resource attributes inside a top-level `"resource"` object (`{"resource": { ... }}`). Unwrap it before reading fields:

```go
	if resource, ok := res["resource"]; ok {
		if resourceMap, ok := resource.(map[string]interface{}); ok {
			res = resourceMap
		}
	}
```

### 4. Manually Flatten and `d.Set(...)` Every Schema Field
Because there is no generated `Resource<Name>Flatten` function, you must explicitly call `d.Set(...)` for `project` (and `zone`/`region`/`location` if Computed/Optional) plus every attribute in the schema:

- **Primitive fields & reserved keyword remapping**:
  ```go
  	if err := d.Set("project", project); err != nil {
  		return fmt.Errorf("Error setting project: %s", err)
  	}
  	if err := d.Set("kind", res["kind"]); err != nil {
  		return fmt.Errorf("Error setting kind: %s", err)
  	}
  	// Map API "id" -> "resource_id" and API "count" -> "block_count" to avoid Terraform reserved words:
  	if err := d.Set("resource_id", res["id"]); err != nil {
  		return fmt.Errorf("Error setting resource_id: %s", err)
  	}
  	if err := d.Set("block_count", res["count"]); err != nil {
  		return fmt.Errorf("Error setting count: %s", err)
  	}
  ```
- **Nested objects (`schema.TypeList` with `Elem: &schema.Resource`)**:
  Convert each nested `map[string]interface{}` into a single-element `[]map[string]interface{}` slice with snake_case keys:
  ```go
  	if physicalTopology, ok := res["physicalTopology"].(map[string]interface{}); ok {
  		topologyList := []map[string]interface{}{
  			{
  				"cluster": physicalTopology["cluster"],
  				"block":   physicalTopology["block"],
  			},
  		}
  		if err := d.Set("physical_topology", topologyList); err != nil {
  			return fmt.Errorf("Error setting physical_topology: %s", err)
  		}
  	}
  ```

### 5. Perform Secondary `LIST` Calls for Aggregated Child Fields (If Needed)
If the data source schema includes a convenience list of child resource names not returned by the single-object `GET` call (e.g., `sub_block_names` on `google_compute_reservation_block`), issue a follow-up `GET` request to the child `LIST` endpoint and set the slice in state:

```go
	listUrl := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/reservations%%2F%s%%2FreservationBlocks%%2F%s/reservationSubBlocks?alt=json&maxResults=500", project, zone, reservation, blockName)
	listRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    listUrl,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error listing ReservationSubBlocks: %s", err)
	}

	subBlockNames := []string{}
	if listRes != nil {
		if items, ok := listRes["items"].([]interface{}); ok {
			for _, item := range items {
				if subBlock, ok := item.(map[string]interface{}); ok {
					if subBlockName, ok := subBlock["name"].(string); ok {
						subBlockNames = append(subBlockNames, subBlockName)
					}
				}
			}
		}
	}
	if err := d.Set("sub_block_names", subBlockNames); err != nil {
		return fmt.Errorf("Error setting sub_block_names: %s", err)
	}
```

### 6. Set the Data Source ID (`d.SetId`)
Finally, set a deterministic resource ID matching the GCP resource path:

```go
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s/reservationBlocks/%s", project, zone, reservation, blockName))
	return nil
}
```

---

## Step 4: Handle `virtual_fields` (If Implemented on the Resource)

Check whether the target resource defines `virtual_fields` in `mmv1/products/<product>/<Resource>.yaml` (or client-side fields in a handwritten resource).

### Why `virtual_fields` Require Special Handling in Data Sources
In MMv1's `mmv1/templates/terraform/resource.go.tmpl`, `resource<Product><Resource>Read` executes the following logic for every virtual field with a `default_value`:
```go
// Set virtual fields to default_value if they are not set in state:
if _, ok := d.GetOkExists("<virtual_field>"); !ok {
	if err := d.Set("<virtual_field>", <default_value>); err != nil {
		return fmt.Errorf("Error setting <virtual_field>: %s", err)
	}
}
```
Because a data source starts with an empty state (`d.GetOkExists("<virtual_field>") == false`), calling `resource<Product><Resource>Read(d, meta)` automatically populates `<virtual_field>` with `<default_value>` even though the field does not exist in the GCP API response!

Inspect each entry under `virtual_fields:` in the resource YAML and handle it according to its category:

### Category 1: Client-Side `virtual_fields` WITH a `default_value` (e.g., `deletion_protection: true`, `deletion_policy: "DELETE"`)
- **Problem**:
  - Resources like `google_alloydb_cluster` and `google_compute_storage_pool` define `deletion_protection` in `virtual_fields` with `default_value: true`.
  - When `dataSource<Name>Read` calls `resource<Name>Read(d, meta)`, `resource<Name>Read` sets `deletion_protection = true` on the data source.
  - Meanwhile, in acceptance tests, the managed resource sets `deletion_protection = false` so the test runner can destroy the resource cleanly.
  - As a result, `acctest.CheckDataSourceStateMatchesResourceState` fails with:
    `expected deletion_protection to be false, got true`.
- **How to Handle**:
  1. **In `dataSource<Product><Resource>Read` (after `resource<Product><Resource>Read(d, meta)`)**:
     Override or clear the defaulted client-side virtual field:
     - **Option A — Set to `false` (or the test/zero value)** (as in `mmv1/third_party/terraform/services/compute/data_source_google_compute_storage_pool.go`):
       ```go
       if err := d.Set("deletion_protection", false); err != nil {
       	return fmt.Errorf("Error setting deletion_protection: %s", err)
       }
       ```
     - **Option B — Clear to `nil`** (as in `mmv1/third_party/terraform/services/alloydb/data_source_alloydb_cluster.go` and `mmv1/third_party/terraform/services/sql/data_source_sql_database.go`):
       ```go
       if err := d.Set("deletion_protection", nil); err != nil {
       	return fmt.Errorf("Error setting deletion_protection: %s", err)
       }
       ```
  2. **In `data_source_<product>_<resource>_test.go`**:
     If the virtual field is cleared to `nil` (Option B) or the test resource sets a client-side virtual field that cannot be read from the API, ignore that field in `CheckDataSourceStateMatchesResourceStateWithIgnores` with an adjacent comment justifying why it is ignored:
     ```go
     // deletion_protection is a client-side virtual field not returned by the API
     acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
     	"data.google_<product>_<resource>.default",
     	"google_<product>_<resource>.default",
     	[]string{"deletion_protection"},
     ),
     ```

### Category 2: Client-Side `virtual_fields` WITHOUT a `default_value` (e.g., `deletion_protection` in `memcache/Instance.yaml`, `force_destroy`, `skip_delete`)
- **Behavior**:
  - When a client-side virtual field has no `default_value` in YAML, `resource<Name>Read` leaves it at its zero value (`false` or `""`) on the data source.
- **How to Handle**:
  - If the acceptance test sets a non-zero value on the managed test resource (for example, `force_destroy = true`), the data source will still have `false`. Pass the field name to `acctest.CheckDataSourceStateMatchesResourceStateWithIgnores` (with an adjacent comment explaining that it is a client-side virtual field).
  - If the acceptance test sets `deletion_protection = false` (matching the zero value), `CheckDataSourceStateMatchesResourceState` will pass without extra overrides.

### Category 3: Computed Output `virtual_fields` (`output: true`) Populated in `custom_code.post_read`
- **Behavior**:
  - Some resources define computed `output: true` fields under `virtual_fields` (for example, `block_names` in `mmv1/products/compute/Reservation.yaml`) and populate them inside a `custom_code.post_read` template (such as `mmv1/templates/terraform/post_read/compute_reservetion.go.tmpl`).
- **How to Handle**:
  - Because `post_read` runs at the end of `resource<Name>Read(d, meta)`, inspect the `post_read` template to see which attributes it reads from `d` (such as `d.Get("name")`, `tpgresource.GetZone(d, config)`, or `tpgresource.GetProject(d, config)`).
  - Ensure all required identifier and locational fields are set in `d` **before** calling `resource<Name>Read(d, meta)` so the `post_read` hook succeeds and populates the output virtual field on the data source.

---

## Step 5: Register the Data Source

At the top of your `data_source_<product>_<resource>.go` file, register the data source with the provider's schema registry via `init()`:

```go
func init() {
	registry.Schema{
		Name:        "google_<product>_<resource>",
		ProductName: "<product>",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSource<Product><Resource>(),
	}.Register()
}
```

---

## Step 6: Write Acceptance Tests

Create `mmv1/third_party/terraform/services/<product>/data_source_<product>_<resource>_test.go` (or `.go.tmpl` if beta-only).

1. **For Resource-Backed Data Sources — Use `acctest.CheckDataSourceStateMatchesResourceState`**:
   - Verify that all attributes on the data source match the managed resource:
     ```go
     Check: resource.ComposeTestCheckFunc(
     	acctest.CheckDataSourceStateMatchesResourceState(
     		"data.google_<product>_<resource>.default",
     		"google_<product>_<resource>.default",
     	),
     ),
     ```
   - If the resource has strictly client-side `virtual_fields` (see **Step 4**) or write-only secrets that cannot be read back from the API, use `acctest.CheckDataSourceStateMatchesResourceStateWithIgnores` with an adjacent comment justifying each ignored attribute.
2. **For Resourceless Data Sources — Use `resource.TestCheckResourceAttr` and `resource.TestCheckResourceAttrSet`**:
   - Since there is no corresponding Terraform resource to compare state against, provision the parent resource(s) (or query a bootstrapped test fixture) and assert key computed attributes directly:
     ```go
     Check: resource.ComposeTestCheckFunc(
     	resource.TestCheckResourceAttrSet("data.google_<product>_<resource>.default", "id"),
     	resource.TestCheckResourceAttrSet("data.google_<product>_<resource>.default", "status"),
     ),
     ```
3. **Test Optional `location` / `region` / `zone` / `project` Fallback**:
   - Ensure the data source works both when `location` / `region` / `zone` / `project` are explicitly provided in HCL **and** when they are omitted (falling back to the test provider's default region/zone/project).

---

## Step 7: Write Documentation

Create `mmv1/third_party/terraform/website/docs/d/<product>_<resource>.html.markdown`:

```markdown
---
subcategory: "<Product Display Name>"
description: |-
  Get information about a <Product> <Resource>.
---

# google_<product>_<resource>

Get information about a <Product> <Resource>. For more information see the [official documentation](https://cloud.google.com/...) and [API](https://cloud.google.com/...).

## Example Usage

```hcl
data "google_<product>_<resource>" "default" {
  <id_field> = "my-resource"
  location   = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `<id_field>` - (Required) The ID/name of the <resource>.

* `location` - (Optional) The location/region/zone in which the resource belongs. If it is not provided, the provider location/region/zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_<product>_<resource>](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/<product>_<resource>) resource for details of all the available attributes.
```

- **For Resourceless Data Sources**: Because there is no Terraform resource documentation page to link to, explicitly list all exported attributes under `## Attributes Reference` (including `id`, primitive attributes, and `<a name="nested_..."></a>` sections for nested blocks, as in `mmv1/third_party/terraform/website/docs/d/compute_reservation_block.html.markdown`).

---

## Step 8: Generate, Build, and Verify

1. **Run Pre-Generation Checks**:
   Follow the `run-pre-gen-checks` skill to verify formatting and static checks.
2. **Generate and Build the Provider**:
   Follow the `generate-provider` skill to generate the downstream provider (`google` or `google-beta`) and compile the binary.
3. **Run Acceptance Tests**:
   Follow the `run-acctests` skill to execute `TestAcc...Datasource...` and verify it passes against the live GCP API.
