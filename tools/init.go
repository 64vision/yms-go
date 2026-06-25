package tools

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
)

var DBM *pg.DB

func init() {
	OpenDB()
}

func OpenDB() {
	fmt.Println("Prod: Initializing tools database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "quantum-db.c9qcwnngsyds.ap-southeast-1.rds.amazonaws.com:5432",
		User:     "postgres",
		Password: "64vision123!",
		Database: "prod_tools_db",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Tools  Database connected!")
	}
}
func OpenDBLocal() {
	fmt.Println("----------------------\nLocal: Initializing tools database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "tools",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Tools database connected!\n--------------------------------------")
	}
}
