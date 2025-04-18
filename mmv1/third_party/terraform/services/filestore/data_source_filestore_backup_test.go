package filestore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFilestoreBackupDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreBackupDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_filestore_backup.backup", "google_filestore_backup.backup"),
				),
			},
		},
	})
}

func testAccFilestoreBackupDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_filestore_instance" "instance" {
  name     = "tf-test-instance%{random_suffix}"
  location = "us-central1-b"
  tier     = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}

resource "google_filestore_backup" "backup" {
  name     = "tf-test-backup%{random_suffix}"
  location = "us-central1"
  source_instance = google_filestore_instance.instance.id
  source_file_share = "share1"
}

data "google_filestore_backup" "backup" {
  name     = google_filestore_backup.backup.name
  location = google_filestore_backup.backup.location
}
`, context)
}
