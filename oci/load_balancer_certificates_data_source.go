// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	oci_load_balancer "github.com/oracle/oci-go-sdk/loadbalancer"
)

func init() {
	RegisterDatasource("oci_load_balancer_certificates", LoadBalancerCertificatesDataSource())
}

func LoadBalancerCertificatesDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readLoadBalancerCertificates,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     LoadBalancerCertificateResource(),
			},
		},
	}
}

func readLoadBalancerCertificates(d *schema.ResourceData, m interface{}) error {
	sync := &LoadBalancerCertificatesDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).loadBalancerClient()

	return ReadResource(sync)
}

type LoadBalancerCertificatesDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_load_balancer.LoadBalancerClient
	Res    *oci_load_balancer.ListCertificatesResponse
}

func (s *LoadBalancerCertificatesDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *LoadBalancerCertificatesDataSourceCrud) Get() error {
	request := oci_load_balancer.ListCertificatesRequest{}

	if loadBalancerId, ok := s.D.GetOkExists("load_balancer_id"); ok {
		tmp := loadBalancerId.(string)
		request.LoadBalancerId = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "load_balancer")

	response, err := s.Client.ListCertificates(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response
	return nil
}

func (s *LoadBalancerCertificatesDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(GenerateDataSourceID())
	resources := []map[string]interface{}{}

	for _, r := range s.Res.Items {
		certificate := map[string]interface{}{}

		if r.CaCertificate != nil {
			certificate["ca_certificate"] = *r.CaCertificate
		}

		if r.CertificateName != nil {
			certificate["certificate_name"] = *r.CertificateName
		}

		if r.PublicCertificate != nil {
			certificate["public_certificate"] = *r.PublicCertificate
		}

		resources = append(resources, certificate)
	}

	if f, fOk := s.D.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, LoadBalancerCertificatesDataSource().Schema["certificates"].Elem.(*schema.Resource).Schema)
	}

	if err := s.D.Set("certificates", resources); err != nil {
		return err
	}

	return nil
}
