package activedirectory_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccActiveDirectoryDomainDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryDomainDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_active_directory_domain.default", "google_active_directory_domain.ad-domain"),
				),
			},
		},
	})
}

func testAccActiveDirectoryDomainDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "ad-network" {
  name                    = "tf-test-ad-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "ad-subnetwork" {
  name          = "tf-test-ad-subnetwork%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.ad-network.id
}

resource "google_active_directory_domain" "ad-domain" {
  domain_name       = "tfgen%{random_suffix}.com"
  locations         = ["us-central1"]
  reserved_ip_range = "192.168.255.0/24"
  authorized_networks = [
    google_compute_network.ad-network.id
  ]
  deletion_protection = false
}

data "google_active_directory_domain" "default" {
  domain_name = google_active_directory_domain.ad-domain.domain_name
}
`, context)
}
