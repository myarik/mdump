package dump

import "fmt"

type PGCredentials struct {
	Host, Username, DBname, Password string
	Port                             int
}

func (pg PGCredentials) String() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s", pg.Username, pg.Password, pg.Host, pg.Port, pg.DBname)
}
