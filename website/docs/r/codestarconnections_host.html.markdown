---
subcategory: "CodeStar Connections"
layout: "aws"
page_title: "AWS: aws_codestarconnections_host"
description: |-
  Provides a CodeStar Host
---

# Resource: aws_codestarconnections_host

Provides a CodeStar Host.

~> **NOTE:** The `aws_codestarconnections_host` resource is created in the state `PENDING`. Authentication with the host provider must be completed in the AWS Console. For more information visit [Set up a pending host](https://docs.aws.amazon.com/dtconsole/latest/userguide/connections-host-setup.html).

## Example Usage

```terraform
resource "aws_codestarconnections_host" "example" {
  name              = "example-host"
  provider_endpoint = "https://example.com"
  provider_type     = "GitHubEnterpriseServer"
}
```

## Argument Reference

This resource supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `name` - (Required) The name of the host to be created. The name must be unique in the calling AWS account.
* `provider_endpoint` - (Required) The endpoint of the infrastructure to be represented by the host after it is created.
* `provider_type` - (Required) The name of the external provider where your third-party code repository is configured.
* `vpc_configuration` - (Optional) The VPC configuration to be provisioned for the host. A VPC must be configured, and the infrastructure to be represented by the host must already be connected to the VPC.

A `vpc_configuration` block supports the following arguments:

* `security_group_ids` - (Required) ID of the security group or security groups associated with the Amazon VPC connected to the infrastructure where your provider type is installed.
* `subnet_ids` - (Required) The ID of the subnet or subnets associated with the Amazon VPC connected to the infrastructure where your provider type is installed.
* `tls_certificate` - (Optional) The value of the Transport Layer Security (TLS) certificate associated with the infrastructure where your provider type is installed.
* `vpc_id` - (Required) The ID of the Amazon VPC connected to the infrastructure where your provider type is installed.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The CodeStar Host ARN.
* `arn` - The CodeStar Host ARN.
* `status` - The CodeStar Host status. Possible values are `PENDING`, `AVAILABLE`, `VPC_CONFIG_DELETING`, `VPC_CONFIG_INITIALIZING`, and `VPC_CONFIG_FAILED_INITIALIZATION`.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import CodeStar Host using the ARN. For example:

```terraform
import {
  to = aws_codestarconnections_host.example-host
  id = "arn:aws:codestar-connections:us-west-1:0123456789:host/79d4d357-a2ee-41e4-b350-2fe39ae59448"
}
```

Using `terraform import`, import CodeStar Host using the ARN. For example:

```console
% terraform import aws_codestarconnections_host.example-host arn:aws:codestar-connections:us-west-1:0123456789:host/79d4d357-a2ee-41e4-b350-2fe39ae59448
```
