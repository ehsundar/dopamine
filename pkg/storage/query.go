package storage

import (
	"embed"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

//go:embed sql
var queryFiles embed.FS

func getQuery(name string, unsafeParams ...string) string {
	query, err := queryFiles.ReadFile(fmt.Sprintf("sql/%s.sql", name))
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot read query: %s: %s", name, err))
	}
	q := string(query)

	for _, p := range unsafeParams {
		q = strings.Replace(q, "?", p, 1)
	}

	return q
}
