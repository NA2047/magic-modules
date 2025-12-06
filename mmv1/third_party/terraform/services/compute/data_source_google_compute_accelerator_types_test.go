// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleComputeAcceleratorTypes_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAcceleratorTypesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeAcceleratorTypes("data.google_compute_accelerator_types.available"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleComputeAcceleratorTypes_withZone(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAcceleratorTypesWithZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeAcceleratorTypes("data.google_compute_accelerator_types.zone_specific"),
					resource.TestCheckResourceAttr("data.google_compute_accelerator_types.zone_specific", "zone", "us-central1-a"),
				),
			},
		},
	})
}

func testAccCheckComputeAcceleratorTypes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find compute accelerator types data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("data source id not set")
		}

		count, ok := rs.Primary.Attributes["types.#"]
		if !ok {
			return errors.New("can't find 'types' attribute")
		}

		cnt, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of types")
		}
		if cnt < 1 {
			return fmt.Errorf("expected at least 1 type, received %d, this is most likely a bug", cnt)
		}

		for i := 0; i < cnt; i++ {
			nameIdx := fmt.Sprintf("types.%d.name", i)
			_, ok := rs.Primary.Attributes[nameIdx]
			if !ok {
				return fmt.Errorf("expected %q, name not found", nameIdx)
			}

			zoneIdx := fmt.Sprintf("types.%d.zone", i)
			_, ok = rs.Primary.Attributes[zoneIdx]
			if !ok {
				return fmt.Errorf("expected %q, zone not found", zoneIdx)
			}

			maxCardsIdx := fmt.Sprintf("types.%d.maximum_cards_per_instance", i)
			_, ok = rs.Primary.Attributes[maxCardsIdx]
			if !ok {
				return fmt.Errorf("expected %q, maximum_cards_per_instance not found", maxCardsIdx)
			}

			selfLinkIdx := fmt.Sprintf("types.%d.self_link", i)
			_, ok = rs.Primary.Attributes[selfLinkIdx]
			if !ok {
				return fmt.Errorf("expected %q, self_link not found", selfLinkIdx)
			}
		}
		return nil
	}
}

const testAccComputeAcceleratorTypesConfig = `
data "google_compute_accelerator_types" "available" {}
`

const testAccComputeAcceleratorTypesWithZoneConfig = `
data "google_compute_accelerator_types" "zone_specific" {
  zone = "us-central1-a"
}
`
