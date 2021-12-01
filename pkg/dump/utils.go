package dump

import "strings"

func getDatabaseNameFromURI(dbURI string) string {
	return dbURI[strings.LastIndex(dbURI, "/")+1:]
}
