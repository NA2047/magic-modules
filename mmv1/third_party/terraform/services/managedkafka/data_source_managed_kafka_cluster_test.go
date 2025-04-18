package managedkafka_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccManagedKafkaClusterDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKafkaClusterDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_managed_kafka_cluster.default", "google_managed_kafka_cluster.cluster"),
				),
			},
		},
	})
}

func testAccManagedKafkaClusterDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-subnetwork%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.network.id
}

resource "google_managed_kafka_cluster" "cluster" {
  cluster_id = "tf-test-cluster%{random_suffix}"
  location   = "us-central1"
  
  capacity_config {
    vcpu_count    = "3"
    memory_bytes  = "3221225472" # 3 GiB
  }
  
  gcp_config {
    access_config {
      network_configs {
        subnet = google_compute_subnetwork.subnetwork.id
      }
    }
  }
}

data "google_managed_kafka_cluster" "default" {
  cluster_id = google_managed_kafka_cluster.cluster.cluster_id
  location   = "us-central1"
}
`, context)
}
