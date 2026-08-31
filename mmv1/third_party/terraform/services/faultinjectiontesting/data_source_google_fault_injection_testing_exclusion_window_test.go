package faultinjectiontesting_test

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)


func TestAccFaultInjectionTestingExclusionWindowDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFaultInjectionTestingExclusionWindowDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_fault_injection_testing_exclusion_window.default", "google_fault_injection_testing_exclusion_window.default"),
				),
			},
		},
	})
}

func testAccFaultInjectionTestingExclusionWindowDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_fault_injection_testing_exclusion_window" "default" {
  exclusion_window_id = "tf-test-ex-win-%{random_suffix}"
  location            = "us-central1"
  description         = "basic exclusion window"
  duration            = "3600s"
}

data "google_fault_injection_testing_exclusion_window" "default" {
  exclusion_window_id = google_fault_injection_testing_exclusion_window.default.exclusion_window_id
  location            = "us-central1"
}
`, context)
}

