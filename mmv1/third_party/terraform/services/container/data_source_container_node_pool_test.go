package container_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContainerNodePoolDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerNodePoolDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_container_node_pool.default", "google_container_node_pool.default"),
				),
			},
		},
	})
}

func testAccContainerNodePoolDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "cluster" {
  name               = "tf-test-cluster%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  deletion_protection = false
}

resource "google_container_node_pool" "default" {
  name       = "tf-test-nodepool%{random_suffix}"
  location   = "us-central1-a"
  cluster    = google_container_cluster.cluster.name
  node_count = 1

  node_config {
    machine_type = "e2-medium"
  }
}

data "google_container_node_pool" "default" {
  name     = google_container_node_pool.default.name
  location = "us-central1-a"
  cluster  = google_container_cluster.cluster.name
}
`, context)
}
