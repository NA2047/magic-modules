package composer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"service_account": fmt.Sprintf("tf-test-%d", acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComposerUserWorkloadsConfigMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_basic(context),
			},
			{
				ResourceName:      "google_composer_user_workloads_config_map.config_map",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_composer_user_workloads_config_map.config_map", "data.db_host", "dbhost:5432"),
					resource.TestCheckNoResourceAttr("google_composer_user_workloads_config_map.config_map", "data.api_host"),
				),
			},
			{
				ResourceName:      "google_composer_user_workloads_config_map.config_map",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_delete(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"service_account": fmt.Sprintf("tf-test-%d", acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComposerUserWorkloadsConfigMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_basic(context),
			},
			{
				ResourceName:      "google_composer_user_workloads_config_map.config_map",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_delete(context),
				Check: resource.ComposeTestCheckFunc(
					testAccComposerUserWorkloadsConfigMapDestroyed(t),
				),
			},
		},
	})
}

func testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%{service_account}"
  display_name = "Test Service Account for Composer Environment"
}
resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "environment" {
  name   = "tf-test-test-environment%{random_suffix}"
  region = "us-central1"
  config {
    node_config {
      service_account  = google_service_account.test.name
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
  depends_on = [google_project_iam_member.composer-worker]
}

resource "google_composer_user_workloads_config_map" "config_map" {
  name = "tf-test-test-config-map%{random_suffix}"
  region = "us-central1"
  environment = google_composer_environment.environment.name
  data = {
    api_host: "apihost:443",
  }
}
`, context)
}

func testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%{service_account}"
  display_name = "Test Service Account for Composer Environment"
}
resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "environment" {
  name   = "tf-test-test-environment%{random_suffix}"
  region = "us-central1"
  config {
    node_config {
      service_account  = google_service_account.test.name
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
  depends_on = [google_project_iam_member.composer-worker]
}

resource "google_composer_user_workloads_config_map" "config_map" {
  name = "tf-test-test-config-map%{random_suffix}"
  region = "us-central1"
  environment = google_composer_environment.environment.name
  data = {
    db_host: "dbhost:5432",
  }
}
`, context)
}

func testAccComposerUserWorkloadsConfigMap_composerUserWorkloadsConfigMapBasicExample_delete(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_service_account" "test" {
  account_id   = "%{service_account}"
  display_name = "Test Service Account for Composer Environment"
}
resource "google_project_iam_member" "composer-worker" {
  project = data.google_project.project.project_id
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}

resource "google_composer_environment" "environment" {
  name   = "tf-test-test-environment%{random_suffix}"
  region = "us-central1"
  config {
    node_config {
      service_account  = google_service_account.test.name
    }
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
  depends_on = [google_project_iam_member.composer-worker]
}
`, context)
}

func testAccComposerUserWorkloadsConfigMapDestroyed(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_composer_user_workloads_config_map" {
				continue
			}

			idTokens := strings.Split(rs.Primary.ID, "/")
			if len(idTokens) != 8 {
				return fmt.Errorf("Invalid ID %q, expected format projects/{project}/regions/{region}/environments/{environment}/userWorkloadsConfigMaps/{name}", rs.Primary.ID)
			}
			_, err := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.UserWorkloadsConfigMaps.Get(rs.Primary.ID).Do()
			if err == nil {
				return fmt.Errorf("config map %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}
