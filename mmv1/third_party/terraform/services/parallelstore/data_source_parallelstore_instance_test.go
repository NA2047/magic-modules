package parallelstore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccParallelstoreInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParallelstoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParallelstoreInstanceDatasourceConfig(context),
			},
		},
	})
}

func testAccParallelstoreInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parallelstore_instance" "instance" {
  instance_id = "test-instance-%{random_suffix}"
  location = "us-central1-a"
  capacity_gib = "1024"
  network = "default"
}

data "google_parallelstore_instance" "default" {
  name = google_parallelstore_instance.instance.instance_id
  location = "us-central1-a"
}
`, context)
}
