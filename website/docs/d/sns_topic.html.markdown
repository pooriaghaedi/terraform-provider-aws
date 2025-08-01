---
subcategory: "SNS (Simple Notification)"
layout: "aws"
page_title: "AWS: aws_sns_topic"
description: |-
  Get information on a Amazon Simple Notification Service (SNS) Topic
---

# Data Source: aws_sns_topic

Use this data source to get the ARN of a topic in AWS Simple Notification
Service (SNS). By using this data source, you can reference SNS topics
without having to hard code the ARNs as input.

## Example Usage

```terraform
data "aws_sns_topic" "example" {
  name = "an_example_topic"
}
```

## Argument Reference

This data source supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `name` - (Required) Friendly name of the topic to match.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arn` - ARN of the found topic, suitable for referencing in other resources that support SNS topics.
* `id` - ARN of the found topic, suitable for referencing in other resources that support SNS topics.
* `tags` - Map of tags for the resource.
