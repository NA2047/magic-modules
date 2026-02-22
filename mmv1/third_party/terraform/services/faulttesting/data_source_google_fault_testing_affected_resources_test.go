package faulttesting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFaultTestingAffectedResourcesDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFaultTestingAffectedResourcesDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_fault_testing_affected_resources.default", "google_fault_testing_affected_resources.default"),
				),
			},
		},
	})
}

func testAccFaultTestingAffectedResourcesDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_fault_testing_affected_resources" "default" {
  affected_resources_id = "tf-test-affected-res-%{random_suffix}"
  location              = "us-central1"
  description           = "basic affected resources"

  action {
    l7_lb_http_fault {
      forwarding_rule = "projects/${local.project}/regions/us-central1/forwardingRules/my-forwarding-rule"
      abort {
        percentage  = 50
        status_code = 503
      }
    }
  }
}

data "google_fault_testing_affected_resources" "default" {
  affected_resources_id = google_fault_testing_affected_resources.default.affected_resources_id
  location              = "us-central1"
}
`, context)
}
