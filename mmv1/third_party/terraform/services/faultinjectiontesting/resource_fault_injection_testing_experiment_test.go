package faultinjectiontesting_test

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)


func TestAccFaultInjectionTestingExperiment_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFaultInjectionTestingExperiment_basic(context),
			},
			{
				ResourceName:      "google_fault_injection_testing_experiment.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFaultInjectionTestingExperiment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_fault_injection_testing_experiment_template" "template" {
  experiment_template_id = "tf-test-exp-temp-%{random_suffix}"
  location               = "us-central1"
  description            = "basic experiment template"
  duration               = "3600s"

  action {
    gce_fail_compute {
      instance = "projects/${data.google_project.project.project_id}/zones/us-central1-a/instances/tf-test-my-instance-%{random_suffix}"
    }
  }
}

resource "google_fault_injection_testing_experiment" "default" {
  experiment_id = "tf-test-exp-%{random_suffix}"
  location      = "us-central1"
  description   = "basic experiment"
  experiment_template = google_fault_injection_testing_experiment_template.template.name
}
`, context)
}
