// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	privateApplicationPackageSingularDataSourceRepresentation = map[string]interface{}{
		"private_application_package_id": Representation{RepType: Required, Create: `${data.oci_service_catalog_private_application_packages.test_private_application_packages.private_application_package_collection.0.items.0.id}`},
	}

	privateApplicationPackageDataSourceRepresentation = map[string]interface{}{
		"private_application_id":         Representation{RepType: Required, Create: `${oci_service_catalog_private_application.test_private_application.id}`},
		"display_name":                   Representation{RepType: Optional, Create: `displayName`},
		"package_type":                   Representation{RepType: Optional, Create: []string{`packageType`}},
		"private_application_package_id": Representation{RepType: Optional, Create: `${oci_service_catalog_private_application_package.test_private_application_package.id}`},
	}

	PrivateApplicationPackageResourceConfig = GenerateDataSourceFromRepresentationMap("oci_service_catalog_private_application_packages", "test_private_application_packages", Required, Create, privateApplicationPackageDataSourceRepresentation) +
		GenerateResourceFromRepresentationMap("oci_service_catalog_private_application", "test_private_application", Required, Create, privateApplicationRepresentation)
)

// issue-routing-tag: service_catalog/default
func TestServiceCatalogPrivateApplicationPackageResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestServiceCatalogPrivateApplicationPackageResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_service_catalog_private_application_packages.test_private_application_packages"
	singularDatasourceName := "data.oci_service_catalog_private_application_package.test_private_application_package"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config + compartmentIdVariableStr + PrivateApplicationPackageResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "private_application_package_collection.0.items.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "private_application_package_collection.0.items.0.display_name", "version"),
				resource.TestCheckResourceAttrSet(datasourceName, "private_application_package_collection.0.items.0.package_type"),
				resource.TestCheckResourceAttrSet(datasourceName, "private_application_package_collection.0.items.0.private_application_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "private_application_package_collection.0.items.0.id"),

				resource.TestCheckResourceAttrSet(datasourceName, "private_application_package_collection.#"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_service_catalog_private_application_package", "test_private_application_package", Required, Create, privateApplicationPackageSingularDataSourceRepresentation) +
				compartmentIdVariableStr + PrivateApplicationPackageResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "private_application_package_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "content_url"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "display_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "mime_type"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "package_type"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "version"),
			),
		},
	})
}
