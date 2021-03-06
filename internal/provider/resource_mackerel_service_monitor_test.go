package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelServiceMonitor_Basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestServiceMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelServiceMonitorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func TestAccMackerelServiceMonitor_Update(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestServiceMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelServiceMonitorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
			{
				Config: testAccMackerelServiceMonitorConfigUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "85.5"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "95.5"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func TestAccMackerelServiceMonitor_Minimum(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestServiceMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelServiceMonitorConfigMinimum(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func testAccCheckMackerelServiceMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_service_monitor" {
			continue
		}

		monitors, err := client.FindMonitors()
		if err != nil {
			return err
		}
		for _, monitor := range monitors {
			if monitor.MonitorID() == rs.Primary.ID {
				return fmt.Errorf("Monitor still exists")
			}
		}
	}

	return nil
}

func testAccMackerelServiceMonitorConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foobar" {
    name = "%s"
}

resource "mackerel_service_monitor" "foobar" {
  name                  = "%s"
  service               = mackerel_service.foobar.name
  duration              = 10
  metric                = "cpu%%"
  operator              = ">"
  warning               = 80.0
  critical              = 90.0
  notification_interval = 10
}
`, rName, rName)
}

func testAccMackerelServiceMonitorConfigUpdate(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foobar" {
    name = "%s"
}

resource "mackerel_service_monitor" "foobar" {
  name                  = "%s"
  service               = mackerel_service.foobar.name
  duration              = 10
  metric                = "cpu%%"
  operator              = ">"
  warning               = 85.5
  critical              = 95.5
  notification_interval = 10
}
`, rName, rName)
}

func testAccMackerelServiceMonitorConfigMinimum(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foobar" {
    name = "%s"
}

resource "mackerel_service_monitor" "foobar" {
  name                  = "%s"
  service               = mackerel_service.foobar.name
  duration              = 10
  metric                = "cpu%%"
  operator              = ">"
  warning               = 80.0
  critical              = 90.0
  notification_interval = 10
}
`, rName, rName)
}
