// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package serviceusage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceServiceUsageConsumerQuotaLimit_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceUsageConsumerQuotaLimit_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_service_usage_consumer_quota_limit.limit", "name"),
					resource.TestCheckResourceAttrSet("data.google_service_usage_consumer_quota_limit.limit", "metric"),
					resource.TestCheckResourceAttrSet("data.google_service_usage_consumer_quota_limit.limit", "unit"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "service", "compute.googleapis.com"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "metric", "compute.googleapis.com%2Fcpus"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "limit_name", "%2Fproject%2Fregion"),
				),
			},
		},
	})
}

func TestAccDataSourceServiceUsageConsumerQuotaLimit_withProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_id":    acctest.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServiceUsageConsumerQuotaLimit_withProject(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_service_usage_consumer_quota_limit.limit", "name"),
					resource.TestCheckResourceAttrSet("data.google_service_usage_consumer_quota_limit.limit", "metric"),
					resource.TestCheckResourceAttrSet("data.google_service_usage_consumer_quota_limit.limit", "unit"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "service", "compute.googleapis.com"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "metric", "compute.googleapis.com%2Fcpus"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "limit_name", "%2Fproject%2Fregion"),
					resource.TestCheckResourceAttr("data.google_service_usage_consumer_quota_limit.limit", "project", context["project_id"].(string)),
				),
			},
		},
	})
}

func testAccDataSourceServiceUsageConsumerQuotaLimit_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_service_usage_consumer_quota_limit" "limit" {
  service    = "compute.googleapis.com"
  metric     = "compute.googleapis.com%%2Fcpus"
  limit_name = "%%2Fproject%%2Fregion"
}
`, context)
}

func testAccDataSourceServiceUsageConsumerQuotaLimit_withProject(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_service_usage_consumer_quota_limit" "limit" {
  project    = "%{project_id}"
  service    = "compute.googleapis.com"
  metric     = "compute.googleapis.com%%2Fcpus"
  limit_name = "%%2Fproject%%2Fregion"
}
`, context)
}
