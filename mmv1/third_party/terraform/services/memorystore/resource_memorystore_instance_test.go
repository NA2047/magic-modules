package memorystore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Validate that replica count is updated for the instance
func TestAccMemorystoreInstance_updateReplicaCount(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create instance with replica count 1
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 1, shardCount: 3, preventDestroy: true, zoneDistributionMode: "MULTI_ZONE", deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update replica count to 2
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 2, shardCount: 3, preventDestroy: true, zoneDistributionMode: "MULTI_ZONE", deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// clean up the resource
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 2, shardCount: 3, preventDestroy: false, zoneDistributionMode: "MULTI_ZONE", deletionProtectionEnabled: false}),
			},
		},
	})
}

// Validate that shard count is updated for the cluster
func TestAccMemorystoreInstance_updateShardCount(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with shard count 3
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 1, shardCount: 3, preventDestroy: true, zoneDistributionMode: "MULTI_ZONE", deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update shard count to 5
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 1, shardCount: 5, preventDestroy: true, zoneDistributionMode: "MULTI_ZONE", deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// clean up the resource
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 1, shardCount: 5, preventDestroy: false, zoneDistributionMode: "MULTI_ZONE", deletionProtectionEnabled: false}),
			},
		},
	})
}

// Validate that engineConfigs is updated for the cluster
func TestAccMemorystoreInstance_updateRedisConfigs(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					engineConfigs: map[string]string{
						"maxmemory-policy": "volatile-ttl",
					},
					deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// add a new redis config key-value pair and update existing redis config
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					engineConfigs: map[string]string{
						"maxmemory-policy":  "allkeys-lru",
						"maxmemory-clients": "90%",
					},
					deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// remove all redis configs
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					engineConfigs: map[string]string{
						"maxmemory-policy":  "allkeys-lru",
						"maxmemory-clients": "90%",
					},
					deletionProtectionEnabled: false}),
			},
		},
	})
}

// Validate that deletion protection is updated for the cluster
func TestAccMemorystoreInstance_updateDeletionProtection(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with deletion protection true
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                      name,
					shardCount:                3,
					zoneDistributionMode:      "MULTI_ZONE",
					deletionProtectionEnabled: true,
				}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update cluster with deletion protection false
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                      name,
					shardCount:                3,
					zoneDistributionMode:      "MULTI_ZONE",
					deletionProtectionEnabled: false,
				}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Validate that engine version is updated for the cluster
func TestAccMemorystoreInstance_updateEngineVersion(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with engine version 7.2
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					engineVersion:        "VALKEY_7_2",
				}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update cluster with engine version 8.0
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					engineVersion:        "VALKEY_8_0",
				}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Validate that deletion protection enabled/disabled cluster is created updated
func TestAccMemorystoreInstance_switchoverAndDetachSecondary(t *testing.T) {
	t.Parallel()

	context_create := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"primary_role":     "PRIMARY",
		"secondary_role":   "SECONDARY",
		"primary_region":   "us-central1",
		"secondary_region": "europe-west1",
		"zone_distubution": "MULTI_ZONE",
	}

	context_switch := map[string]interface{}{
		"random_suffix":    acctest.RandString(t, 10),
		"role":             "PRIMARY",
		"primary_region":   "us-central1",
		"secondary_region": "europe-west1",
		"zone_distubution": "MULTI_ZONE",
	}

	piName := fmt.Sprintf("tf-test-prim-%d", acctest.RandInt(t))
	siName := fmt.Sprintf("tf-test-sec-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{

				Config: testAccMemorystoreInstance_switchoverAndDetachSecondaryCreate(context_create),
				// create primary and secondary clusters cluster
				// Config: testAccMemorystoreInstance_switchoverAndDetachSecondary(&InstanceParams{name: piName, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", shouldCreateSecondary: true, secondaryInstanceName: siName, icrRole: "SECONDARY"}),
			},
			{
				ResourceName:            "google_memorystore_instance.test_secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				Config: testAccMemorystoreInstance_switchoverAndDetachSecondarySwitch(context_switch),

				// // Switchover to secondary cluster
				// Config: testAccMemorystoreInstance_switchoverAndDetachSecondary(&InstanceParams{name: piName, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", shouldCreateSecondary: true, secondaryInstanceName: siName, icrRole: "PRIMARY"}),
			},
			{
				ResourceName:            "google_memorystore_instance.test_secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
		},
	})
}

func testAccMemorystoreInstance_switchoverAndDetachSecondaryCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
// Primary instance
resource "google_memorystore_instance" "primary_instance" {
  instance_id                    = "tf-prime-instance-%{random_suffix}"
  shard_count                    = 1
  location                       = "%{primary_region}"
  replica_count                  = 2
  node_type                      = "SHARED_CORE_NANO"
  transit_encryption_mode        = "TRANSIT_ENCRYPTION_DISABLED"
  authorization_mode             = "AUTH_DISABLED"
  engine_configs = {
    maxmemory-policy             = "volatile-ttl"
  }
  zone_distribution_config {
    mode                         = "%{zone_distubution}"
    zone                         = "%{primary_region}-b"
  }

// replication config
   cross_instance_replication_config { 
      instance_role              = "%{primary_role}"
  }

  deletion_protection_enabled    = false
  mode                           = "CLUSTER_DISABLED"
  labels = {
    "abc" : "xyz"
  }
  depends_on                     = [google_network_connectivity_service_connection_policy.primary_policy]
}


// Primary policy
resource "google_network_connectivity_service_connection_policy" "primary_policy" {
	name                          = "tf-test-prime-policy-%{random_suffix}"
	location                      = "%{primary_region}"
	service_class                 = "gcp-memorystore"
	network                       = google_compute_network.primary_producer_net.id
	psc_config {                 
	  subnetworks                 = [google_compute_subnetwork.primary_producer_subnet.id]
	}
  }

  // Primary subnet
  resource "google_compute_subnetwork" "primary_producer_subnet" {
	name                           = "tf-test-prime-subnet-%{random_suffix}"
	ip_cidr_range                  = "10.0.1.0/29"
	region                         = "%{primary_region}"
	network                        = google_compute_network.primary_producer_net.id
  }

  // Primary vpc
  resource "google_compute_network" "primary_producer_net" {
	name                           = "tf-test-prime-net-%{random_suffix}"
	auto_create_subnetworks        = false
  }

// Secondary instance
resource "google_memorystore_instance" "secondary_instance" {
  instance_id                     = "tf-test-sec-instance-%{random_suffix}"
  shard_count                     = 1
  location                        = "%{secondary_region}"
  replica_count                   = 2
  node_type                       = "SHARED_CORE_NANO"
  transit_encryption_mode         = "TRANSIT_ENCRYPTION_DISABLED"
  authorization_mode              = "AUTH_DISABLED"
  engine_configs = {
    maxmemory-policy              = "volatile-ttl"
  }
  zone_distribution_config {
    mode                          = "%{zone_distubution}"
    zone                          = "%{secondary_region}-b"
  }
  deletion_protection_enabled     = false
  mode                            = "CLUSTER_DISABLED"
  // Cross instance replication config
  cross_instance_replication_config {
    instance_role                 = "%{secondary_role}"
  }
     primary_instance {
       instance                   = google_memorystore_instance.primary_instance.id
     }
     secondary_instances {
        instance                  = google_memorystore_instance.primary_instance.id
        }

  }
  labels = {
    "abc" : "xyz"
  }
  depends_on                      = [google_network_connectivity_service_connection_policy.secondary_policy]

}

// Secondary policy
resource "google_network_connectivity_service_connection_policy" "secondary_policy" {
  name                            = "tf-test-sec-ploicy-%{random_suffix}"
  location                        = "%{secondary_region}"
  service_class                   = "gcp-memorystore-autopush"
  description                     = "my basic service connection policy"
  network                         = google_compute_network.secondary_producer_net.id
  psc_config {                 
    subnetworks                   = [google_compute_subnetwork.secondary_producer_subnet.id]
  }
}

// Secondary subnet
resource "google_compute_subnetwork" "secondary_producer_subnet" {
  name                            = "tf-test-sec-subnet-%{random_suffix}"
  ip_cidr_range                   = "10.0.2.0/29"
  region                          = "%{secondary_region}"
  network                         = google_compute_network.secondary_producer_net.id
}

// Secondary vpc
resource "google_compute_network" "secondary_producer_net" {
  name                            = "tf-test-sec-net-%{random_suffix}"
  auto_create_subnetworks         = false
}

data "google_project" "project" {
}

`, context)
}

func testAccMemorystoreInstance_switchoverAndDetachSecondarySwitch(context map[string]interface{}) string {
	return acctest.Nprintf(`
// Primary instance
resource "google_memorystore_instance" "primary_instance" {
  instance_id                    = "tf-prime-instance-%{random_suffix}"
  shard_count                    = 1
  location                       = "%{primary_region}"
  replica_count                  = 2
  node_type                      = "SHARED_CORE_NANO"
  transit_encryption_mode        = "TRANSIT_ENCRYPTION_DISABLED"
  authorization_mode             = "AUTH_DISABLED"
  engine_configs = {
    maxmemory-policy             = "volatile-ttl"
  }
  zone_distribution_config {
    mode                         = "%{zone_distubution}"
    zone                         = "%{primary_region}-b"
  }

// replication config
   cross_instance_replication_config { 
      instance_role              = "%{primary_role}"
  }

  deletion_protection_enabled    = false
  mode                           = "CLUSTER_DISABLED"
  labels = {
    "abc" : "xyz"
  }
  depends_on                     = [google_network_connectivity_service_connection_policy.primary_policy]
}


// Primary policy
resource "google_network_connectivity_service_connection_policy" "primary_policy" {
	name                          = "tf-test-prime-policy-%{random_suffix}"
	location                      = "%{primary_region}"
	service_class                 = "gcp-memorystore"
	network                       = google_compute_network.primary_producer_net.id
	psc_config {                 
	  subnetworks                 = [google_compute_subnetwork.primary_producer_subnet.id]
	}
  }

  // Primary subnet
  resource "google_compute_subnetwork" "primary_producer_subnet" {
	name                           = "tf-test-prime-subnet-%{random_suffix}"
	ip_cidr_range                  = "10.0.1.0/29"
	region                         = "%{primary_region}"
	network                        = google_compute_network.primary_producer_net.id
  }

  // Primary vpc
  resource "google_compute_network" "primary_producer_net" {
	name                           = "tf-test-prime-net-%{random_suffix}"
	auto_create_subnetworks        = false
  }

// Secondary instance
resource "google_memorystore_instance" "secondary_instance" {
  instance_id                     = "tf-test-sec-instance-%{random_suffix}"
  shard_count                     = 1
  location                        = "%{secondary_region}"
  replica_count                   = 2
  node_type                       = "SHARED_CORE_NANO"
  transit_encryption_mode         = "TRANSIT_ENCRYPTION_DISABLED"
  authorization_mode              = "AUTH_DISABLED"
  engine_configs = {
    maxmemory-policy              = "volatile-ttl"
  }
  zone_distribution_config {
    mode                          = "%{zone_distubution}"
    zone                          = "%{secondary_region}-b"
  }
  deletion_protection_enabled     = false
  mode                            = "CLUSTER_DISABLED"
  // Cross instance replication config
  cross_instance_replication_config {
    instance_role                 = "%{secondary_role}"
  }
     primary_instance {
       instance                   = google_memorystore_instance.primary_instance.id
     }
     secondary_instances {
        instance                  = google_memorystore_instance.primary_instance.id
        }

  }
  labels = {
    "abc" : "xyz"
  }
  depends_on                      = [google_network_connectivity_service_connection_policy.secondary_policy]

}

// Secondary policy
resource "google_network_connectivity_service_connection_policy" "secondary_policy" {
  name                            = "tf-test-sec-ploicy-%{random_suffix}"
  location                        = "%{secondary_region}"
  service_class                   = "gcp-memorystore-autopush"
  description                     = "my basic service connection policy"
  network                         = google_compute_network.secondary_producer_net.id
  psc_config {                 
    subnetworks                   = [google_compute_subnetwork.secondary_producer_subnet.id]
  }
}

// Secondary subnet
resource "google_compute_subnetwork" "secondary_producer_subnet" {
  name                            = "tf-test-sec-subnet-%{random_suffix}"
  ip_cidr_range                   = "10.0.2.0/29"
  region                          = "%{secondary_region}"
  network                         = google_compute_network.secondary_producer_net.id
}

// Secondary vpc
resource "google_compute_network" "secondary_producer_net" {
  name                            = "tf-test-sec-net-%{random_suffix}"
  auto_create_subnetworks         = false
}

data "google_project" "project" {
}

`, context)
}

// Validate that persistence config is updated for the cluster
func TestAccMemorystoreInstance_updatePersistence(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create instance with AOF enabled
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, preventDestroy: true, zoneDistributionMode: "MULTI_ZONE", persistenceMode: "AOF", deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update persitence to RDB
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, preventDestroy: true, zoneDistributionMode: "MULTI_ZONE", persistenceMode: "RDB", deletionProtectionEnabled: false}),
			},
			{
				ResourceName:      "google_memorystore_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// clean up the resource
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, preventDestroy: false, zoneDistributionMode: "MULTI_ZONE", persistenceMode: "RDB", deletionProtectionEnabled: false}),
			},
		},
	})
}

// Validate that instance endpoints are updated for the cluster
func TestAccMemorystoreInstance_updateInstanceEndpoints(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with no user created connections
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", userEndpointCount: 0}),
			},
			{
				ResourceName:            "google_memorystore_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// create cluster with one user created connection
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", userEndpointCount: 1}),
			},
			{
				ResourceName:            "google_memorystore_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 2 endpoints
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", userEndpointCount: 2}),
			},
			{
				ResourceName:            "google_memorystore_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 1 endpoint
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", userEndpointCount: 1}),
			},
			{
				ResourceName:            "google_memorystore_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 0 endpoints
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", userEndpointCount: 0}),
			},
			{
				ResourceName:            "google_memorystore_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateMemorystoreInstance(&InstanceParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", userEndpointCount: 0}),
			},
		},
	})
}

type InstanceParams struct {
	name                      string
	replicaCount              int
	shardCount                int
	preventDestroy            bool
	nodeType                  string
	engineConfigs             map[string]string
	zoneDistributionMode      string
	zone                      string
	deletionProtectionEnabled bool
	persistenceMode           string
	engineVersion             string
	userEndpointCount         int
	shouldCreateSecondary     bool
	secondaryInstanceName     string
	icrRole                   string
}

func createMemorystoreInstanceEndpointsWithOneUserCreatedConnections(params *InstanceParams) string {
	return fmt.Sprintf(`
		resource "google_memorystore_instance_desired_user_created_endpoints" "default" {

		name                           = "%s"
		region                         = "europe-west1"
		desired_user_created_endpoints {
			connections {
				psc_connection {
					psc_connection_id  = google_compute_forwarding_rule.forwarding_rule1_network1.psc_connection_id
					ip_address         = google_compute_address.ip1_network1.address
					forwarding_rule    = google_compute_forwarding_rule.forwarding_rule1_network1.id
					network            = google_compute_network.network1.id
					project_id         = data.google_project.project.project_id
					service_attachment = google_memorystore_instance.test.psc_attachment_details[0].service_attachment
				}
			}
		desired_user_created_endpoints {
				psc_connection {
					psc_connection_id  = google_compute_forwarding_rule.forwarding_rule2_network1.psc_connection_id
					ip_address         = google_compute_address.ip2_network1.address
					forwarding_rule    = google_compute_forwarding_rule.forwarding_rule2_network1.id
					network            = google_compute_network.network1.id
					service_attachment = google_memorystore_instance.test.psc_attachment_details[1].service_attachment
				}
			}
		}
		}
		%s
		`,
		params.name,
		createMemorystoreUserCreatedConnection1(params),
	)

}

func createMemorystoreInstanceEndpointsWithTwoUserCreatedConnections(params *InstanceParams) string {
	return fmt.Sprintf(`
		resource "google_memorystore_instance_desired_user_created_endpoints" "default" {
		name                           = "%s"
		region                         = "europe-west1"
		desired_user_created_endpoints {
			connections {
				psc_connection {
					psc_connection_id  = google_compute_forwarding_rule.forwarding_rule1_network1.psc_connection_id
					ip_address         = google_compute_address.ip1_network1.address
					forwarding_rule    = google_compute_forwarding_rule.forwarding_rule1_network1.id
					network            = google_compute_network.network1.id
					project_id         = data.google_project.project.project_id
					service_attachment = google_memorystore_instance.test.psc_attachment_details[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id  = google_compute_forwarding_rule.forwarding_rule2_network1.psc_connection_id
					ip_address         = google_compute_address.ip2_network1.address
					forwarding_rule    = google_compute_forwarding_rule.forwarding_rule2_network1.id
					network            = google_compute_network.network1.id
					service_attachment = google_memorystore_instance.test.psc_attachment_details[1].service_attachment
				}
			}
		}
		desired_user_created_endpoints {
			connections {
				psc_connection {
					psc_connection_id  = google_compute_forwarding_rule.forwarding_rule1_network2.psc_connection_id
					ip_address         = google_compute_address.ip1_network2.address
					forwarding_rule    = google_compute_forwarding_rule.forwarding_rule1_network2.id
					network            = google_compute_network.network2.id
					service_attachment = google_memorystore_instance.test.psc_attachment_details[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id  = google_compute_forwarding_rule.forwarding_rule2_network2.psc_connection_id
					ip_address         = google_compute_address.ip2_network2.address
					forwarding_rule    = google_compute_forwarding_rule.forwarding_rule2_network2.id
					network            = google_compute_network.network2.id
					service_attachment = google_memorystore_instance.test.psc_attachment_details[1].service_attachment
				}
			}
		}
		}
		%s
		%s
		`,
		params.name,
		createMemorystoreUserCreatedConnection1(params),
		createMemorystoreUserCreatedConnection2(params),
	)
}
func createMemorystoreUserCreatedConnection1(params *InstanceParams) string {
	return fmt.Sprintf(`
		resource "google_compute_forwarding_rule" "forwarding_rule1_network1" {
		name                          = "%s"
		region                        = "europe-west1"
		ip_address                    = google_compute_address.ip1_network1.id
		load_balancing_scheme         = ""
		network                       = google_compute_network.network1.id
		target                        = google_memorystore_instance.test.psc_attachment_details[0].service_attachment
		} 

		resource "google_compute_forwarding_rule" "forwarding_rule2_network1" {
		name                          = "%s"
		region                        = "europe-west1"
		ip_address                    = google_compute_address.ip2_network1.id
		load_balancing_scheme         = ""
		network                       = google_compute_network.network1.id
		target                        = google_memorystore_instance.test.psc_attachment_details[1].service_attachment
		}

		resource "google_compute_address" "ip1_network1" {
		name                          = "%s"
		region                        = "europe-west1"
		subnetwork                    = google_compute_subnetwork.subnet_network1.id
		address_type                  = "INTERNAL"
		purpose                       = "GCE_ENDPOINT"
		}

		resource "google_compute_address" "ip2_network1" {
		name                         = "%s"
		region                       = "europe-west1"
		subnetwork                   = google_compute_subnetwork.subnet_network1.id
		address_type                 = "INTERNAL"
		purpose                      = "GCE_ENDPOINT"
		}

		resource "google_compute_subnetwork" "subnet_network1" {
		name                         = "%s"
		ip_cidr_range                = "10.0.0.248/29"
		region                       = "europe-west1"
		network                      = google_compute_network.network1.id
		}

		resource "google_compute_network" "network1" {
		name                         = "%s"
		auto_create_subnetworks      = false
		}
		`,
		params.name+"-11", // fwd-rule1-net1
		params.name+"-12", // fwd-rule2-net1
		params.name+"-11", // ip1-net1
		params.name+"-12", // ip2-net1
		params.name+"-1",  // subnet-net1
		params.name+"-1",  // net1
	)
}

func createMemorystoreUserCreatedConnection2(params *InstanceParams) string {
	return fmt.Sprintf(`
		resource "google_compute_forwarding_rule" "forwarding_rule1_network2" {
		name                         = "%s"
		region                       = "europe-west1"
		ip_address                   = google_compute_address.ip1_network2.id
		load_balancing_scheme        = ""
		network                      = google_compute_network.network2.id
		target                       = google_memorystore_instance.test.psc_attachment_details[0].service_attachment
		}

		resource "google_compute_forwarding_rule" "forwarding_rule2_network2" {
		name                         = "%s"
		region                       = "europe-west1"
		ip_address                   = google_compute_address.ip2_network2.id
		load_balancing_scheme        = ""
		network                      = google_compute_network.network2.id
		target                       = google_memorystore_instance.test.psc_attachment_details[1].service_attachment
		}

		resource "google_compute_address" "ip1_network2" {
		name                         = "%s"
		region                       = "europe-west1"
		subnetwork                   = google_compute_subnetwork.subnet_network2.id
		address_type                 = "INTERNAL"     
		purpose                      = "GCE_ENDPOINT"
		}

		resource "google_compute_address" "ip2_network2" {
		name                         = "%s"
		region                       = "europe-west1"
		subnetwork                   = google_compute_subnetwork.subnet_network2.id
		address_type                 = "INTERNAL"
		purpose                      = "GCE_ENDPOINT"
		}

		resource "google_compute_subnetwork" "subnet_network2" {
		name                         = "%s"
		ip_cidr_range                = "10.0.0.248/29"
		region                       = "europe-west1"
		network                      = google_compute_network.network2.id
		}

		resource "google_compute_network" "network2" {
		name                         = "%s"
		auto_create_subnetworks      = false
		}
		`,
		params.name+"-21", // fwd-rule1-net2
		params.name+"-22", // fwd-rule2-net2
		params.name+"-21", // ip1-net2
		params.name+"-22", // ip2-net2
		params.name+"-2",  // subnet-net2
		params.name+"-2",  // net2
	)
}

func createOrUpdateMemorystoreInstance(params *InstanceParams) string {
	lifecycleBlock := ""
	if params.preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	var strBuilder strings.Builder
	for key, value := range params.engineConfigs {
		strBuilder.WriteString(fmt.Sprintf("%s =  \"%s\"\n", key, value))
	}

	zoneDistributionConfigBlock := ``
	if params.zoneDistributionMode != "" {
		zoneDistributionConfigBlock = fmt.Sprintf(`
		zone_distribution_config {
			mode = "%s"
			zone = "%s"
		}
		`, params.zoneDistributionMode, params.zone)
	}
	persistenceBlock := ``
	if params.persistenceMode != "" {
		persistenceBlock = fmt.Sprintf(`
		persistence_config {
			mode = "%s"
		}
		`, params.persistenceMode)
	}

	if params.userEndpointCount == 2 {
		createMemorystoreInstanceEndpointsWithTwoUserCreatedConnections(params)
	} else if params.userEndpointCount == 1 {
		createMemorystoreInstanceEndpointsWithOneUserCreatedConnections(params)
	}

	return fmt.Sprintf(`
resource "google_memorystore_instance" "test" {
    instance_id  = "%s"
	replica_count = %d
	shard_count = %d
	node_type = "%s"
	location         = "europe-west1"
	desired_psc_auto_connections  {
			network = google_compute_network.producer_net.id
            project_id = data.google_project.project.project_id
	}
    deletion_protection_enabled = %t
	engine_version = "%s"
	engine_configs = {
		%s
	}
  %s
  %s
	depends_on = [
			google_network_connectivity_service_connection_policy.default
		]
	%s
}

resource "google_network_connectivity_service_connection_policy" "default" {
	name = "%s"
	location = "europe-west1"
	service_class = "gcp-memorystore"
	description   = "my basic service connection policy"
	network = google_compute_network.producer_net.id
	psc_config {
	subnetworks = [google_compute_subnetwork.producer_subnet.id]
	}
}

resource "google_compute_subnetwork" "producer_subnet" {
	name          = "%s"
	ip_cidr_range = "10.0.0.248/29"
	region        = "europe-west1"
	network       = google_compute_network.producer_net.id
}

resource "google_compute_network" "producer_net" {
	name                    = "%s"
	auto_create_subnetworks = false
}

data "google_project" "project" {
}
`, params.name, params.replicaCount, params.shardCount, params.nodeType, params.deletionProtectionEnabled, params.engineVersion, strBuilder.String(), zoneDistributionConfigBlock, persistenceBlock, lifecycleBlock, params.name, params.name, params.name)
}
