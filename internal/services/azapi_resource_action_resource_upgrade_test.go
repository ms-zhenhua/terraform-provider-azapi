package services_test

import (
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAzapiActionResourceUpgrade_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.basic(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_basicWhenDestroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.basicWhenDestroy(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.basicWhenDestroy(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_registerResourceProvider(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.registerResourceProvider(),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.registerResourceProvider(),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_upgradeFromVeryOldVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.registerResourceProvider(),
			Check:  resource.ComposeTestCheckFunc(),
		}, "1.8.0"),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.registerResourceProvider(),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.providerAction(),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_nonstandardLRO(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.nonstandardLRO(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.nonstandardLRO(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.timeouts(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_timeouts_from_v1_13_1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: strings.ReplaceAll(r.timeouts(data), `update = "10m"`, ""),
			Check:  resource.ComposeTestCheckFunc(),
		}, "1.13.1"),
		data.UpgradeTestApplyStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
	})
}
