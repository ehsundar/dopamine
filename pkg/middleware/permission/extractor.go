package permission

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Extractor func(r *http.Request) string

var TablePermissionExtractorFactory = func(apiType string) Extractor {
	return func(r *http.Request) string {
		vars := mux.Vars(r)
		table := vars["table"]
		permissionForTable := getPermissionForTable(table, apiType)
		return permissionForTable
	}
}

var StaticExtractorFactory = func(staticPermission string) Extractor {
	return func(_ *http.Request) string {
		return staticPermission
	}
}
