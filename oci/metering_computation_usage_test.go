// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	filter = `<<EOF
{
	"operator": "AND",
	"dimentions": [],
	"tags": [],
	"filters": [
		"operator": "OR",
		"dimentions": [
			"key": "compartName"
			"value": "dxterraformtest"
		]
		"filters": []
		"tags": []
	]
}
EOF`

	usageRepresentationWithOptionals = `resource "oci_metering_computation_usage" "test_usage" {
compartment_depth = 1
filter = <<EOF
{
		"operator": "AND",
		"dimensions": [],
		"tags": [],
		"filters": [
			{
				"operator": "OR",
			 	"dimensions": [
					{
						"key": "compartmentName",
						"value": "dxterraformtest"
					}
				],
				"filters": [],
				"tags": []
			}
		]
}
EOF
granularity = "DAILY"
group_by = ["service"]
query_type = "COST"
tenant_id = "${var.tenancy_id}"
time_usage_ended = "${var.time_usage_ended}"
time_usage_started = "${var.time_usage_started}"
time_forecast_ended= "2021-03-21T00:00:00Z"
}
`

	UsageRequiredOnlyResource = UsageResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_metering_computation_usage", "test_usage", Required, Create, usageRepresentation)

	usageRepresentation = map[string]interface{}{
		"granularity":        Representation{RepType: Required, Create: `DAILY`},
		"tenant_id":          Representation{RepType: Required, Create: `${var.tenancy_id}`},
		"time_usage_ended":   Representation{RepType: Required, Create: `2021-03-19T00:00:00Z`},
		"time_usage_started": Representation{RepType: Required, Create: `2021-03-18T00:00:00Z`},
		"compartment_depth":  Representation{RepType: Optional, Create: `1`},
		//"filter":               Representation{RepType: Optional, Create: },
		"filter":   Representation{RepType: Optional, Create: `{\"operator\":\"OR\",\"dimensions\":[{\"key\":\"compartmentName\",\"value\":\"dxterraformtest\"}],\"tags\":[],\"filters\":[]}`, Update: `{\"operator\":\"OR\",\"dimensions\":[{\"key\":\"compartmentName\",\"value\":\"dxterraformtest\"}],\"tags\":[],\"filters\":[]}`},
		"forecast": RepresentationGroup{Optional, usageForecastRepresentation},
		"group_by": Representation{RepType: Optional, Create: []string{`service`}},
		//"group_by_tag":         RepresentationGroup{Optional, usageGroupByTagRepresentation},
		"is_aggregate_by_time": Representation{RepType: Optional, Create: `false`},
		"query_type":           Representation{RepType: Optional, Create: `COST`},
	}
	usageForecastRepresentation = map[string]interface{}{
		"time_forecast_ended":   Representation{RepType: Required, Create: `2021-03-20T00:00:00Z`},
		"forecast_type":         Representation{RepType: Optional, Create: `BASIC`},
		"time_forecast_started": Representation{RepType: Optional, Create: `2021-03-19T00:00:00Z`},
	}
	usageGroupByTagRepresentation = map[string]interface{}{
		"key":       Representation{RepType: Optional, Create: `key`},
		"namespace": Representation{RepType: Optional, Create: `namespace`},
		"value":     Representation{RepType: Optional, Create: `value`},
	}

	UsageResourceDependencies = ""
)

// issue-routing-tag: metering_computation/default
func TestMeteringComputationUsageResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestMeteringComputationUsageResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)
	tenancyId := getEnvSettingWithBlankDefault("tenancy_ocid")
	tenancyIdVariableStr := fmt.Sprintf("variable \"tenancy_id\" { default = \"%s\" }\n", tenancyId)
	usgaeEndTimeStr, usageStartTimeStr := generateUsageRepresentationWithCurrentTimeInputs()
	usgaeEndTimeVariableStr := fmt.Sprintf("variable \"time_usage_ended\" { default = \"%s\" }\n", usgaeEndTimeStr)
	usageStartTimeVariableStr := fmt.Sprintf("variable \"time_usage_started\" { default = \"%s\" }\n", usageStartTimeStr)

	resourceName := "oci_metering_computation_usage.test_usage"

	var resId string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+UsageResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_metering_computation_usage", "test_usage", Optional, Create, usageRepresentation), "usageapi", "usage", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify Create
		{
			PreConfig: func() {
				fmt.Printf("config is : %s", GenerateResourceFromRepresentationMap("oci_metering_computation_usage", "test_usage", Optional, Create, usageRepresentation))
			},
			Config: config + compartmentIdVariableStr + tenancyIdVariableStr + usgaeEndTimeVariableStr + usageStartTimeVariableStr + UsageResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_metering_computation_usage", "test_usage", Required, Create, usageRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "granularity", "DAILY"),
				resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
				resource.TestCheckResourceAttrSet(resourceName, "time_usage_ended"),
				resource.TestCheckResourceAttrSet(resourceName, "time_usage_started"),
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + tenancyIdVariableStr + usgaeEndTimeVariableStr + usageStartTimeVariableStr + UsageResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + tenancyIdVariableStr + usgaeEndTimeVariableStr + usageStartTimeVariableStr + UsageResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_metering_computation_usage", "test_usage", Optional, Create, usageRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_depth", "1"),
				resource.TestCheckResourceAttr(resourceName, "forecast.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "forecast.0.forecast_type", "BASIC"),
				resource.TestCheckResourceAttr(resourceName, "forecast.0.time_forecast_ended", "2021-03-20T00:00:00Z"),
				resource.TestCheckResourceAttr(resourceName, "forecast.0.time_forecast_started", "2021-03-19T00:00:00Z"),
				resource.TestCheckResourceAttrSet(resourceName, "filter"),
				resource.TestCheckResourceAttr(resourceName, "granularity", "DAILY"),
				resource.TestCheckResourceAttr(resourceName, "group_by.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "query_type", "COST"),
				resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
				resource.TestCheckResourceAttrSet(resourceName, "time_usage_ended"),
				resource.TestCheckResourceAttrSet(resourceName, "time_usage_started"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					if isEnableExportCompartment, _ := strconv.ParseBool(getEnvSettingWithDefault("enable_export_compartment", "true")); isEnableExportCompartment {
						if errExport := TestExportCompartmentWithResourceName(&resId, &compartmentId, resourceName); errExport != nil {
							return errExport
						}
					}
					return err
				},
			),
		},
	})
}

func generateUsageRepresentationWithCurrentTimeInputs() (string, string) {
	t := time.Now()
	year, month, day := t.Date()
	endTime := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	startTime := endTime.Add(-24 * time.Hour)
	usgaeEndTimeStr := endTime.Format("2006-01-02T15:04:05Z")
	usageStartTimeStr := startTime.Format("2006-01-02T15:04:05Z")
	return usgaeEndTimeStr, usageStartTimeStr
}
