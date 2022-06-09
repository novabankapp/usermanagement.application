package utilities

import (
	"fmt"
	noSql "github.com/novabankapp/usermanagement.data/repositories/base/cassandra"
)

func MakeQueries(queries []map[string]string, field, compare, value string) []map[string]string {

	m := make(map[string]string)
	m[noSql.COLUMN] = field
	m[noSql.COMPARE] = compare
	m[noSql.VALUE] = value
	queries = append(queries, m)
	return queries
}
func FormatPhonePasswordResetMessage(pin string, expiryDate string) string {
	return fmt.Sprintf("Your Password reset pin is %s and will expire after %s", pin, expiryDate)
}
func FormatEmailPasswordResetMessage(hash string, expiryDate string) string {
	return fmt.Sprintf("Your Password reset pin is %s and will expire after %s", hash, expiryDate)
}
