package main

import (
	"regexp"
	"strings"
)

func errConvertIdToPath(err error, d *DAG) string {
	errStr := err.Error()
	errByte := []byte(errStr)

	reg := regexp.MustCompile(`'[a-z0-9-]+'`)
	newErr := reg.ReplaceAllFunc(errByte, func(b []byte) []byte {
		id := string(b)
		id = strings.ReplaceAll(id, "'", "")

		path, _ := d.IdToPath(id)
		return []byte(path)
	})

	return string(newErr)
}
