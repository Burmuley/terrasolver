package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/heimdalr/dag"
)

type DAG struct {
	*dag.DAG
}

func NewDAG() *DAG {
	d := &DAG{}
	d.DAG = dag.NewDAG()
	return d
}

func (d *DAG) TopologicalSort() ([]string, error) {
	l, err := topologicalSortRecursive(d.DAG)
	l = d.idsToPaths(l)
	return l, err
}

func (d *DAG) IdToPath(id string) (string, error) {
	path, err := d.GetVertex(id)

	if err != nil {
		return "", err
	}

	return fmt.Sprint(path), nil
}

func (d *DAG) idsToPaths(ids []string) []string {
	p := make([]string, 0, len(ids))

	for _, id := range ids {
		path, err := d.GetVertex(id)

		if err != nil {
			log.Fatal(err)
		}
		p = append(p, fmt.Sprint(path))
	}

	return p
}

func (d *DAG) FillDAGFromFiles(files []string) error {
	inds := make(map[string]string)

	for _, f := range files {
		ff, _ := filepath.Abs(filepath.Dir(f))
		fid := ""

		if _, ok := inds[ff]; !ok {
			fid, _ = d.AddVertex(ff)
			inds[ff] = fid
		} else {
			fid = inds[ff]
		}
		deps, err := ParseDependencies(f)

		if err != nil {
			return err
		}

		for _, p := range deps {
			pp, _ := filepath.Abs(filepath.Join(ff, p))
			pid := ""

			if _, ok := inds[pp]; !ok {
				pid, _ = d.AddVertex(pp)
				inds[pp] = pid
			} else {
				pid = inds[pp]
			}

			if err = d.AddEdge(fid, pid); err != nil {
				return err
			}
		}
	}

	return nil
}

func topologicalSortRecursive(d *dag.DAG) ([]string, error) {
	visited := make(map[string]struct{})
	stack := make([]string, 0, len(d.GetVertices()))

	for id := range d.GetVertices() {
		if _, ok := visited[id]; !ok {
			if err := topologicalSortUtil(d, &stack, visited, id); err != nil {
				return reverseList(stack), err
			}
		}
	}

	return reverseList(stack), nil
}

func topologicalSortUtil(d *dag.DAG, stack *[]string, visited map[string]struct{}, id string) error {
	visited[id] = struct{}{}
	descs, err := d.GetAncestors(id)

	if err != nil {
		return err
	}

	for i := range descs {
		if _, ok := visited[i]; !ok {
			if err = topologicalSortUtil(d, stack, visited, i); err != nil {
				return err
			}
		}
	}

	*stack = append(*stack, id)
	return nil
}

func reverseList(l []string) []string {
	for i := 0; i < len(l)/2; i++ {
		j := len(l) - i - 1
		l[i], l[j] = l[j], l[i]
	}
	return l
}
