package faultinjectiontesting_test

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)


func TestAccFaultInjectionTestingExperimentTemplate_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFaultInjectionTestingExperimentTemplate_basic(context),
			},
			{
				ResourceName:      "google_fault_injection_testing_experiment_template.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFaultInjectionTestingExperimentTemplate_update(context),
			},
			{
				ResourceName:      "google_fault_injection_testing_experiment_template.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFaultInjectionTestingExperimentTemplate_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_fault_injection_testing_experiment_template" "default" {
  experiment_template_id = "tf-test-exp-temp-%{random_suffix}"
  location               = "us-central1"
  description            = "basic experiment template"
  duration               = "3600s"

  action {
    cloud_sql_failover {
      instance = "projects/${data.google_project.project.project_id}/instances/tf-test-exp-instance-%{random_suffix}"
    }
  }
}
`, context)
}

func testAccFaultInjectionTestingExperimentTemplate_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_fault_injection_testing_experiment_template" "default" {
  experiment_template_id = "tf-test-exp-temp-%{random_suffix}"
  location               = "us-central1"
  description            = "updated experiment template"
  duration               = "7200s"

  action {
    cloud_sql_failover {
      instance = "projects/${data.google_project.project.project_id}/instances/tf-test-exp-instance-%{random_suffix}"
    }
  }
}
`, context)
}

