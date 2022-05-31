package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	cfgVar_TERRASOLVER_PATH         = "TERRASOLVER_PATH"
	cfgVar_TERRASOLVER_SKIP_CONFIRM = "TERRASOLVER_SKIP_CONFIRM"
)

const (
	tsPathDefault =
)

func main() {
	// TODO: replace with cmd parameters parsing
	cwd, _ := os.Getwd()
	tsPath := flag.String("path", cwd, "Path to Terragrunt working directory")
	tsSkipConfirm := flag.Bool("skip-confirm", false, "Skip confirmation user input request")
	flag.Parse()
	tgArgs := flag.Args()

	config := map[string]string{
		"TERRASOLVER_PATH": *tsPath,
		"TERRASOLVER_SKIP_CONFIRM": fmt.Sprintf("%b", *tsSkipConfirm),
	}
	config := readConfigEnv()

	modulesPath, _ := os.Getwd()

	if path, ok := config[cfgVar_TERRASOLVER_PATH]; ok {
		modulesPath = path
	}

	//path := "test/env1"
	modulesPath, _ = filepath.Abs(modulesPath)
	log.Println("Terragrunt modules directory:", modulesPath)
	terragrunt_bin := "terragrunt"

	// Find all .hcl files in underlying directory tree
	files, err := FindFilesByExt(modulesPath, ".hcl")

	if err != nil {
		log.Fatal(err)
	}

	// Create a DAG and fill it from the list of detected Terragrunt modules
	d := NewDAG()
	if err := d.FillDAGFromFiles(files); err != nil {
		log.Fatal(errConvertIdToPath(err, d))
	}

	//fmt.Println(d)
	sorted, err := d.TopologicalSort()

	if err != nil {
		log.Fatal(errConvertIdToPath(err, d))
	}

	fmt.Printf("Running order for modules in '%s':\n", modulesPath)
	printList(sorted)
	if _, ok := config[cfgVar_TERRASOLVER_SKIP_CONFIRM]; !ok {
		fmt.Println("Press ENTER to continue...")
		b := bufio.NewReader(os.Stdin)
		b.ReadString('\n')
	}

	q := NewExecQueue(sorted)

	for m := q.Next(); m != nil; m = q.Next() {
		log.Printf("Working on %s ...\n", m.Path)
		log.Println(terragrunt_bin, " ", tgArgs)
		err := m.Exec(terragrunt_bin, tgArgs...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printList(l []string) {
	for n, s := range l {
		fmt.Println(n+1, " - ", s)
	}
}

func readConfigEnv() map[string]string {
	config_vars := []string{
		cfgVar_TERRASOLVER_PATH,
		cfgVar_TERRASOLVER_SKIP_CONFIRM,
	}
	config := make(map[string]string, 0)

	for _, v := range config_vars {
		if e := os.Getenv(v); e != "" {
			config[v] = e
		}
	}

	return config
}
