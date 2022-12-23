package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/heimdalr/dag"
)

type DAG struct {
	*dag.DAG
	reverse bool
}

func NewDAG() *DAG {
	d := &DAG{reverse: true}
	d.DAG = dag.NewDAG()
	return d
}

func (d *DAG) SetReverse(r bool) {
	d.reverse = r
}

func (d *DAG) TopologicalSort() ([]string, error) {
	l, err := topologicalSortRecursive(d.DAG, d.reverse)
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

func (d *DAG) FillDAGFromFiles(files []string, deepDive bool, inds map[string]string, warnings bool) error {
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
				// skip adding edge if error is EdgeDuplicateError
				// return err in all other cases
				if errors.As(err, &dag.EdgeDuplicateError{}) {
					if warnings {
						log.Printf("%s, skipping", errConvertIdToPath(err, d))
					}
				} else {
					return err
				}
			}

			if deepDive {
				files, err := FindFilesByExt(pp, ".hcl")

				if err != nil {
					return err
				}

				if err := d.FillDAGFromFiles(files, deepDive, inds, warnings); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func topologicalSortRecursive(d *dag.DAG, reverse bool) ([]string, error) {
	visited := make(map[string]struct{})
	stack := make([]string, 0, len(d.GetVertices()))

	for id := range d.GetVertices() {
		if _, ok := visited[id]; !ok {
			if err := topologicalSortUtil(d, &stack, visited, id); err != nil {
				if reverse {
					stack = reverseList(stack)
				}
				return stack, err
			}
		}
	}

	if reverse {
		stack = reverseList(stack)
	}
	return stack, nil
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
