# Terragrunt dependencies solver (Terrasolver)

This tool is aimed to solve issue with resolving dependencies among a set of Terragrunt modules and run 
provided Terragrunt command for each module in a sequence automatically.

It works pretty the same way as Terragrunt: you put a Terraform command and command line options and `terrasolver` will
propagate all this to `terragrunt` and it will do its job as usual. :) 

### How to build
```shell
go build -o terrasolver
```

### How to use

Configuration for the tool is done via environment variables to lower any probability of overlapping options among all 
three tools involved: Terrasolver, Terragrunt, Terraform.

Supported configuration ~options~variables:
* `TERRASOLVER_PATH` - path to start scanning for the Terragrunt modules; if not defined `terrasolver` will use current working difectory
* `TERRASOLVER_SKIP_CONFIRM` - if set to any value, then `terrasolver` will not make a pause for user input after printing the running sequence

Usage is simple:
```shell
terrasolver <Terrafor command> ...<Terraform/Terragrunt options>
```

Examples:
```shell
terrasolver apply -auto-approve
```

```shell
terrasolver plan
```
