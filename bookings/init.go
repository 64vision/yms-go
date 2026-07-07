package bookings

import (
	"gollux/dbconfig"

	"github.com/go-pg/pg"
)

var DBM *pg.DB

func init() {

	DBM = dbconfig.DBM
}
