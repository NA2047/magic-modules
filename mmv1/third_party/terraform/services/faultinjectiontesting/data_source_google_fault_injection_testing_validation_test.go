package faultinjectiontesting_test

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)


func TestAccFaultInjectionTestingValidationDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFaultInjectionTestingValidationDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_fault_injection_testing_validation.default", "google_fault_injection_testing_validation.default"),
				),
			},
		},
	})
}

func testAccFaultInjectionTestingValidationDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_fault_injection_testing_validation" "default" {
  validation_id = "tf-test-val-%{random_suffix}"
  location      = "us-central1"
  description   = "basic validation"

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

data "google_fault_injection_testing_validation" "default" {
  validation_id = google_fault_injection_testing_validation.default.validation_id
  location      = "us-central1"
}
`, context)
}
