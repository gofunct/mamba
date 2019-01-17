package terraform

import (
	"github.com/gofunct/mamba/api/tf"
	"github.com/hashicorp/terraform/config/configschema"
	"github.com/hashicorp/terraform/terraform"
)

type GrpcProvisionerFunc func() tf.ProvisionerServer
type ProvisionerFunc func() terraform.ResourceProvisioner
type SchemaFunc func() (*configschema.Block, error)
type ValidateFunc func(*terraform.ResourceConfig) ([]string, []error)
type ApplyFunc func(terraform.UIOutput, *terraform.InstanceState, *terraform.ResourceConfig) error
type StopFunc func() error
