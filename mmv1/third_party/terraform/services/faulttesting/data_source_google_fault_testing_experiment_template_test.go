package faulttesting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFaultTestingExperimentTemplateDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFaultTestingExperimentTemplateDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_fault_testing_experiment_template.default", "google_fault_testing_experiment_template.default"),
				),
			},
		},
	})
}

func testAccFaultTestingExperimentTemplateDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_fault_testing_experiment_template" "default" {
  experiment_template_id = "tf-test-exp-temp-%{random_suffix}"
  location               = "us-central1"
  description            = "basic experiment template"
  duration               = "3600s"

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

data "google_fault_testing_experiment_template" "default" {
  experiment_template_id = google_fault_testing_experiment_template.default.experiment_template_id
  location               = "us-central1"
}
`, context)
}
