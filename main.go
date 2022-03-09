package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	cfgVar_TERRASOLVER_PATH         = "TERRASOLVER_PATH"
	cfgVar_TERRASOLVER_SKIP_CONFIRM = "TERRASOLVER_SKIP_CONFIRM"
)

func main() {
	// TODO: replace with cmd parameters parsing
	config := readConfigEnv()

	modules_path, _ := os.Getwd()

	if path, ok := config[cfgVar_TERRASOLVER_PATH]; ok {
		modules_path = path
	}

	//path := "test/env1"
	modules_path, _ = filepath.Abs(modules_path)
	fmt.Println("Terragrunt modules directory:", modules_path)
	terragrunt_bin := "terragrunt"

	// Find all .hcl files in underlying directory tree
	files, err := FindFilesByExt(modules_path, ".hcl")

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

	fmt.Printf("Running order for modules in '%s':\n", modules_path)
	printList(sorted)
	if _, ok := config[cfgVar_TERRASOLVER_SKIP_CONFIRM]; !ok {
		fmt.Println("Press ENTER to continue...")
		b := bufio.NewReader(os.Stdin)
		b.ReadString('\n')
	}

	q := NewExecQueue(sorted)

	for m := q.Next(); m != nil; m = q.Next() {
		fmt.Println(terragrunt_bin, " ", os.Args[1:])
		err := m.Exec(terragrunt_bin, os.Args[1:]...)
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
