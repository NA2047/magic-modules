package firebaseextensions_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFirebaseExtensionsInstanceDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseExtensionsInstanceDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_firebase_extensions_instance.default", "google_firebase_extensions_instance.default"),
				),
			},
		},
	})
}

func testAccFirebaseExtensionsInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firebase_extensions_instance" "default" {
  provider = google-beta
  instance_id = "tf-test-ext%{random_suffix}"
  config {
    extension_ref = "firebase/storage-resize-images"
    params = {
      IMG_SIZES = "200x200"
      LOCATION = "us-central1"
      DELETE_ORIGINAL_FILE = "false"
      MAKE_PUBLIC = "false"
    }
  }
}

data "google_firebase_extensions_instance" "default" {
  provider = google-beta
  instance_id = google_firebase_extensions_instance.default.instance_id
}
`, context)
}
