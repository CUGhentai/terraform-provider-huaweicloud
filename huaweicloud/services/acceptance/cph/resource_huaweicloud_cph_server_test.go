package cph

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCphServerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCphServer: Query the CPH instance
	var (
		getCphServerHttpUrl = "v1/{project_id}/cloud-phone/servers/{id}"
		getCphServerProduct = "cph"
	)
	getCphServerClient, err := cfg.NewServiceClient(getCphServerProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CPH Client: %s", err)
	}

	getCphServerPath := getCphServerClient.Endpoint + getCphServerHttpUrl
	getCphServerPath = strings.ReplaceAll(getCphServerPath, "{project_id}", getCphServerClient.ProjectID)
	getCphServerPath = strings.ReplaceAll(getCphServerPath, "{id}", state.Primary.ID)

	getCphServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCphServerResp, err := getCphServerClient.Request("GET", getCphServerPath, &getCphServerOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CphServer: %s", err)
	}

	getCphServerRespBody, err := utils.FlattenResponse(getCphServerResp)
	if err != nil {
		return nil, err
	}
	statusRaw := utils.PathSearch("status", getCphServerRespBody, nil)

	if fmt.Sprint(statusRaw) == "6" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getCphServerRespBody, nil
}

func TestAccCphServer_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cph_server.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCphServerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCphServer_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "availability_zone"),
					resource.TestCheckResourceAttrSet(rName, "order_id"),
					resource.TestCheckResourceAttrSet(rName, "addresses.#"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testCphServer_basic_update(name, name+"update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_id", "eip_type", "auto_renew", "period", "period_unit"},
			},
		},
	})
}

func testCphServerBase(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_kps_keypair" "test" {
  name        = "%s"
  description = "%s"
}
`, name, name, name, name)
}

func testCphServer_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cph_server" "test" {
  name          = "%s"
  server_flavor = "physical.rx1.xlarge"
  phone_flavor  = "rx1.cp.c15.d46.e1v1"
  image_id      = "22110120221209e101210a3000000b02"
  keypair_name  = huaweicloud_kps_keypair.test.name

  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
  eip_type  = "5_bgp"

  bandwidth {
    share_type  = "0"
    charge_mode = "1"
    size        = 300
  }

  period_unit = "month"
  period      = 1
  auto_renew  = "true"

  lifecycle {
    ignore_changes = [
      image_id, auto_renew, period, period_unit,
    ]
  }
}
`, testCphServerBase(name), name)
}

func testCphServer_basic_update(name, cphServerName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cph_server" "test" {
  name          = "%s"
  server_flavor = "physical.rx1.xlarge"
  phone_flavor  = "rx1.cp.c15.d46.e1v1"
  image_id      = "22110120221209e101210a3000000b02"
  keypair_name  = huaweicloud_kps_keypair.test.name

  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
  eip_type  = "5_bgp"

  bandwidth {
    share_type  = "0"
    charge_mode = "1"
    size        = 300
  }

  period_unit = "month"
  period      = 1
  auto_renew  = "true"

  lifecycle {
    ignore_changes = [
      image_id, auto_renew, period, period_unit,
    ]
  }
}
`, testCphServerBase(name), cphServerName)
}
