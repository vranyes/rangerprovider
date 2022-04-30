# ranger_provider
A Terraform Provider for Apache Ranger

## WhY ARe You DoiNG ThIS?!?/

Right now we are managing Ranger policies with Ansible. And it is terrible.
Also I want to learn GO!

## How Can I contribute?

EMail me? Make a PR?
I guess mature projects have like style guide lines and stuff but I think GO obviated the need for those?

## Usage

TBH this is TBD

Thinking it'll be something like

```hcl
provider "ranger" {
    host = data.aws_ssm_parameter.ranger_host.value
    username = data.aws_ssm_parameter.ranger_username.value
    password = data.aws_ssm_parameter.ranger_password.value
}

resource "ranger_policy" {
    name = "My Policy"
    description = "This is a Terraform Managed Policy"
    service_type: "hbase"

    resources = [
        ...
    ]

    policy_items = [
        {
            accesses = [
                "read",
                "create",
                "write"
            ]
            groups = [
                "admins"
            ]
        }
    ]
}

```