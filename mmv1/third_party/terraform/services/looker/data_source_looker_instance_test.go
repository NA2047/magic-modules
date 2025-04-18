package looker_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLookerInstanceDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLookerInstanceDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_looker_instance.default", "google_looker_instance.looker-instance"),
				),
			},
		},
	})
}

func testAccLookerInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_looker_instance" "looker-instance" {
  name     = "tf-test-looker%{random_suffix}"
  region   = "us-central1"
  platform_edition = "LOOKER_CORE_STANDARD_ANNUAL"
  
  oauth_config {
    client_id     = "client-id"
    client_secret = "client-secret"
  }
}

data "google_looker_instance" "default" {
  name   = google_looker_instance.looker-instance.name
  region = "us-central1"
}
`, context)
}
