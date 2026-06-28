package account

import (
	"fmt"
	"os"
	"time"

	"github.com/go-pg/pg"
)

var DBM *pg.DB

// var (
// 	WSCONN   *grpc.ClientConn
// 	WSCLIENT pb.WalletServiceClient
// 	WSONCE   sync.Once
// )

// const (
// 	WALLETSERVER = "localhost:5051"
// )

func init() {
	//InitWalletRPC()
	time.Sleep(3000)
	OpenDB()
}

func OpenDB() {
	fmt.Println("----------------------\n Prod:  Initializing player account database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "",
		User:     "postgres",
		Password: "64vision123!",
		Database: "yms_db",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Account database connected!\n--------------------------------------")
	}
}

func OpenDBLocal() {
	fmt.Println("----------------------\nLOCAL: Initializing player account database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "accounts",
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Account database connected!\n--------------------------------------")
	}
}
