# DEPRECATED

This code was used to create https://github.com/akamai/cli-terraform

# Akamai Property Manager to Terraform Provider CLI plugin

An [Akamai CLI](https://developer.akamai.com/cli) plugin for generating an [Akamai Terraform Provider](https://github.com/terraform-providers/terraform-provider-akamai) configuration from existing Property Manager configuration.

Note, these instructions assume you already have a working CLI and Terraform installation.

## Getting Started

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install https://github.com/IanCassTwo/akamai-property-to-terraform.git
```
 ## Usage
 ```
 Usage:
  akamai terraform [global flags] command [command flags] [arguments...]

Description:
   Generate Terraform configs from your existing Property Manager configs

Built-In Commands:
  create
  list
  help

Global Flags:
   --edgerc value  Location of the credentials file (default: "/home/icass/.edgerc") [$AKAMAI_EDGERC]
   --section value     Section of the credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --accountkey value  Account switch key
```

1. Create a new empty subdirectory to hold your new Terraform configuration and change into it.
2. ```akamai terraform create <name of property manager config>```
3. ```terraform init```
4. ```terraform apply```

