package health

import (
	"database/sql"
	"github.com/google/wire"
	"gocloud.dev/health"
	"gocloud.dev/health/sqlhealth"
)

var Set = wire.NewSet(
	New,
)

func New(db *sql.DB) ([]health.Checker, func()) {
	dbCheck := sqlhealth.New(db)
	list := []health.Checker{dbCheck}
	return list, func() {
		dbCheck.Stop()
	}
}
