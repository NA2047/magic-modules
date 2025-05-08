package lustre_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLustreInstance_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "default-vpc"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLustreInstance_basic(context),
			},
			{
				ResourceName:      "google_lustre_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLustreInstance_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_lustre_instance.instance",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "gke_support_enabled", "location", "terraform_labels"},
			},
		},
	})
}

func testAccLustreInstance_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parallelstore_instance" "instance" {
  instance_id               = "instance%{random_suffix}"
  location                  = "us-central1-a"
  description               = "test instance"
  capacity_gib              = 12000
  deployment_type           = "SCRATCH"
  network                   = google_compute_network.network.name
  reserved_ip_range         = google_compute_global_address.private_ip_alloc.name
  file_stripe_level         = "FILE_STRIPE_LEVEL_MIN"
  directory_stripe_level    = "DIRECTORY_STRIPE_LEVEL_MIN"
  depends_on                = [google_service_networking_connection.default]
}

resource "google_compute_network" "network" {
  name                      = "network%{random_suffix}"
  auto_create_subnetworks   = true
  mtu                       = 8896
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name                      = "address%{random_suffix}"
  purpose                   = "VPC_PEERING"
  address_type              = "INTERNAL"
  prefix_length             = 24
  network                   = google_compute_network.network.id
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                   = google_compute_network.network.id
  service                   = "servicenetworking.googleapis.com"
  reserved_peering_ranges   = [google_compute_global_address.private_ip_alloc.name]
}

data "google_compute_network" "lustre-network" {
  instance_id                = google_parallelstore_instance.instance.id
}
`, context)
}
