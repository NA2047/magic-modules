package firebasedatabase_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFirebaseDatabaseInstanceDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseDatabaseInstanceDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_firebase_database_instance.default", "google_firebase_database_instance.default"),
				),
			},
		},
	})
}

func testAccFirebaseDatabaseInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firebase_database_instance" "default" {
  provider = google-beta
  region   = "us-central1"
  instance_id = "tf-test-db%{random_suffix}"
  type     = "USER_DATABASE"
}

data "google_firebase_database_instance" "default" {
  provider = google-beta
  region   = "us-central1"
  instance_id = google_firebase_database_instance.default.instance_id
}
`, context)
}
