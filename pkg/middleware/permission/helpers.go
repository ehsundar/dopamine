package permission

import (
	"context"
	"github.com/ehsundar/dopamine/api/auth/token"
	"github.com/ehsundar/dopamine/pkg/middleware/auth"
	"github.com/spf13/viper"
)

func getSubject(ctx context.Context) *token.Subject {
	s, ok := ctx.Value(auth.UserSubjectContextKey).(*token.Subject)
	if !ok {
		return nil
	}
	return s
}

func getPermissionForTable(table string, apiType string) string {
	tablesConfig := viper.Sub("tables")
	tb := tablesConfig.GetStringMapString(table)
	perm, ok := tb[apiType]

	if ok {
		return perm
	} else {
		if table == "default" {
			return "public"
		} else {
			return getPermissionForTable("default", apiType)
		}
	}
}
