# Terragrunt dependencies solver (Terrasolver)

This tool is aimed to solve issue with resolving dependencies among a set of Terragrunt modules and run 
provided Terragrunt command for each module in a sequence automatically.

It works pretty the same way as Terragrunt: you put a Terraform command and command line options and `terrasolver` will
propagate all this to `terragrunt` and it will do its job as usual. :) 

### How to build

Included script `build.sh` builds binary for three of supported platforms: MacOS X (amd64), MacOS X (arm64), Linux (amd64)

```shell
$ ./build.sh
$ ls -1 terrasolver_*                                                                                                           burmuley@RRG-MacPro15
terrasolver_linux_amd64
terrasolver_mac_amd64
terrasolver_mac_arm64
```

### How to use

```shell
terrasolver [flags] [terragrunt command and parameters]
```

Supported flags:
* `-path` - Path to the working directory where to run all activities. Is omitted will use current directory.
* `-skip-confirm` - Skip confirmation step after the ordered list modules is displayed, will continue with running Terragrunt command against each module.
* `-terragrunt` - Path to the Terragrunt binary. The default is /usr/local/bin/terragrunt
* `-deepdive` - If set to `false` will only scan current working directory for dependencies, 
   If set to `true` - will also recursively scan dependencies referenced in files within the working directory
  to build the complete dependency tree if any of modules enlist dependencies out of the working directory.
* `-version` - displays version and build information

Most of the flags listed above can be also overridden by corresponding environment variables.

**Note**: values set with environment variables take precedence over values in command line!

* `TERRASOLVER_PATH` - same as `-path` flag
* `TERRASOLVER_SKIP_CONFIRM` - same as `-skip-confirm` flag
* `TERRASOLVER_TERRAGRUNT_BIN` - same as `-terragrunt` flag
* `TERRASOLVER_DEEP_DIVE` - same as `-deepdive` flag

Example:
```shell
$ terrasolver -path=/home/user/infrastructure/dev -deepdive=true apply -auto-approve
2022/06/08 21:01:18 Terragrunt modules directory: /home/user/infrastructure/dev
Running order for modules in '/home/user/infrastructure/dev':
#1: /home/user/infrastructure/dev/us-west-2/ecs-clusters
#2: /home/user/infrastructure/dev/us-west-2/target-groups
#3: /home/user/infrastructure/dev/us-west-2/load-balancers
#4: /home/user/infrastructure/dev/us-west-2/kms-keys
#5: /home/user/infrastructure/dev/us-east-1/kms-replica-keys
#6: /home/user/infrastructure/dev/us-west-2/s3-buckets
#7: /home/user/infrastructure/dev/us-east-1/s3-buckets
#8: /home/user/infrastructure/dev/global/iam-roles
#9: /home/user/infrastructure/dev/us-west-2/code-deploy
Press ENTER to continue...
...
<then it simply runs `terragrunt [command] for each module in the list above`>
```
