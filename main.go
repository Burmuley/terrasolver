/*
Terrasolver scans directory structure looking for .hcl files with Terragrunt dependencies definitions.
All dependencies are grouped into a DAG (Directed Acyclic Graph) and then sorted to calculate proper order of execution.
Each module is represented as a directory with `terragrunt.hcl` file.

After dependency graph is calculated the tool simply goes over the list and passes the command to Terragrunt.

Usage:
    terrasolver [flags] [terragrunt command and parameters]

Flags
    -path
        Path to the working directory where to run all activities. Is omitted will use current directory.
    -skip-confirm
        Skip confirmation step after the ordered list modules is displayed,
        will continue with running Terragrunt command against each module.
    -terragrunt
        Path to the Terragrunt binary. The default is /usr/local/bin/terragrunt
    -deepdive
        If set to false will only scan current working directory for dependencies.
        If set to true - will also recursively scan dependencies referenced in files within the working directory
        to build the complete dependency tree if any of modules enlist dependencies out of the working directory.
    -version
        Displays version and build information.

Environment variables

Most of the flags listed above can be also overridden by corresponding environment variable.

Note: flags set with environment variables take precedence over flags in command line!

    TERRASOLVER_PATH - same as -path flag
    TERRASOLVER_SKIP_CONFIRM - same as -skip-confirm flag
    TERRASOLVER_TERRAGRUNT_BIN" - same as -terragrunt flag
    TERRASOLVER_DEEP_DIVE - same as -deepdive flag

Example:
    terrasolver -path=/home/user/infrastructure/dev -deepdive=true apply -auto-approve

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

*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	cfgTerrasolverPath        = "TERRASOLVER_PATH"
	cfgTerrasolverSkipConfirm = "TERRASOLVER_SKIP_CONFIRM"
	cfgTerragruntBinary       = "TERRASOLVER_TERRAGRUNT_BIN"
	cfgTerrasolverDeepDive    = "TERRASOLVER_DEEP_DIVE"
	cfgTerrasolverAutoApprove = "TERRASOLVER_AUTO_APPROVE"
)

var (
	version    string = "no version set"
	commit     string = "no commit set"
	repository        = "github.com/Burmuley/terrasolver"
)

const (
	terragruntBinDefault = "/usr/local/bin/terragrunt"
)

func main() {
	// setup and parse command line args
	cwd, _ := os.Getwd()
	tsPath := flag.String("path", cwd, "Path to Terragrunt working directory")
	tsSkipConfirm := flag.Bool("skip-confirm", false, "Skip confirmation user input request")
	tsTerragruntBin := flag.String("terragrunt", terragruntBinDefault, "Path to Terragrunt binary")
	tsDeepDive := flag.Bool("deepdive", true, "Deep scan for dependencies")
	tsAddAutoApprove := flag.Bool("auto-approve", true, "Automatically add `-auto-approve` flag to the Terragrunt arugs")
	tsVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()
	tgArgs := flag.Args()

	if *tsVersion {
		fmt.Println("Version: ", version)
		fmt.Println("Repository: ", repository)
		fmt.Println("Git commit: ", commit)
		os.Exit(0)
	}

	config := map[string]string{
		cfgTerrasolverPath:        *tsPath,
		cfgTerrasolverSkipConfirm: fmt.Sprintf("%v", *tsSkipConfirm),
		cfgTerragruntBinary:       *tsTerragruntBin,
		cfgTerrasolverDeepDive:    fmt.Sprintf("%v", *tsDeepDive),
		cfgTerrasolverAutoApprove: fmt.Sprintf("%v", *tsAddAutoApprove),
	}
	config = readConfigEnv(config)

	// check if auto-approve is present in tgArgs and add if missing
	autoApprove, _ := strconv.ParseBool(config[cfgTerrasolverAutoApprove])
	if autoApprove {
		hasAutoApprove := false
		applyCommand := false
		for _, arg := range tgArgs {
			if strings.Contains(arg, "-auto-approve") {
				hasAutoApprove = true
				break
			}
		}

		// only add -auto-approve flag if command is `apply`
		// (not supported for other Terraform commands)
		for _, arg := range tgArgs {
			if strings.Contains(arg, "apply") {
				applyCommand = true
				break
			}
		}

		if !hasAutoApprove && applyCommand {
			tgArgs = append(tgArgs, "-auto-approve")
		}
	}

	modulesPath, _ := config[cfgTerrasolverPath]
	modulesPath, _ = filepath.Abs(modulesPath)
	log.Println("Terragrunt modules directory:", modulesPath)
	terragruntBin := config[cfgTerragruntBinary]

	// Find all .hcl files in underlying directory tree
	files, err := FindFilesByExt(modulesPath, ".hcl")
	if err != nil {
		log.Fatal(err)
	}

	// Create a DAG and fill it from the list of detected Terragrunt modules
	dag := NewDAG()
	// check if command line contains `destroy` command and
	// do not reverse topological sort result if `destroy` command present
	for _, arg := range tgArgs {
		if strings.EqualFold(arg, "destroy") {
			dag.SetReverse(false)
			break
		}
	}
	deepDive, _ := strconv.ParseBool(config[cfgTerrasolverDeepDive])
	inds := make(map[string]string)
	if err := dag.FillDAGFromFiles(files, deepDive, inds); err != nil {
		log.Fatal(errConvertIdToPath(err, dag))
	}

	sorted, err := dag.TopologicalSort()
	if err != nil {
		log.Fatal(errConvertIdToPath(err, dag))
	}

	fmt.Printf("Running order for modules in '%s':\n", modulesPath)
	for n, s := range sorted {
		fmt.Printf("#%d: %s\n", n+1, s)
	}

	skipConfirm, _ := strconv.ParseBool(config[cfgTerrasolverSkipConfirm])
	if !skipConfirm {
		fmt.Println("Press ENTER to continue...")
		b := bufio.NewReader(os.Stdin)
		_, _ = b.ReadString('\n')
	}

	q := NewExecQueue(sorted)

	for m := q.Next(); m != nil; m = q.Next() {
		log.Printf("Working on %s ...\n", m.GetPath())
		log.Println(terragruntBin, " ", tgArgs)
		err := m.Exec(terragruntBin, tgArgs...)
		if err != nil {
			log.Fatal(err)
		}
	}
}
