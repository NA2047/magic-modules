package workflows_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccWorkflowsWorkflowDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowsWorkflowDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_workflows_workflow.default", "google_workflows_workflow.workflow"),
				),
			},
		},
	})
}

func testAccWorkflowsWorkflowDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workflows_workflow" "workflow" {
  name            = "workflow%{random_suffix}"
  region          = "us-central1"
  description     = "A sample workflow"
  source_contents = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - Logs a message
  # - Waits for 10 seconds
  # - Logs another message
  
  main:
    params: [input]
    steps:
      - log_info:
          call: sys.log
          args:
            severity: INFO
            message: Hello from a workflow!
      - wait:
          call: sys.sleep
          args:
            seconds: 10
      - log_info_again:
          call: sys.log
          args:
            severity: INFO
            message: Another message!
  EOF
}

data "google_workflows_workflow" "default" {
  name   = google_workflows_workflow.workflow.name
  region = "us-central1"
}
`, context)
}
