package game

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
	fmt.Println("Prod: Initializing game race database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "quantum-db.c9qcwnngsyds.ap-southeast-1.rds.amazonaws.com:5432",
		User:     "postgres",
		Password: "64vision123!",
		Database: "prod_games_db",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Game Race  Database connected!")
	}
}
func OpenDBLocal() {
	fmt.Println("----------------------\nLocal: Initializing game database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "games",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Game database connected!\n--------------------------------------")
	}
}
