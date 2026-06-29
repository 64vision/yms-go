package account

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-pg/pg"
)

type Configuration struct {
	System         System         `json:"system"`
	Dbconfig       Dbconfig       `json:"dbconfig"`
	PathPermission PathPermission `json:"path_permission"`
}
type System struct {
	Port         string `json:"port"`
	SalesLogDir  string `json:"sales_log_dir"`
	AccessLogDir string `json:"access_log_dir"`
	ErrorLogDir  string `json:"error_log_dir"`
}
type Dbconfig struct {
	Addr     string `json:"addr"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type PathPermission struct {
	NoAuth []string `json:"no_auth"`
	Mobile []string `json:"mobile"`
}

var DBM *pg.DB
var PRODUCTION = false
var CONFIG *Configuration

var PRODCONFIGPATH = "/home/ubuntu/engine/config.json"

var CONFIGPATH = "../config.json"

func init() {
	fmt.Println("----------------------\n Prod:  Initializing...")
	time.Sleep(3000)
	OpenDB()
}

func CheckAndLoadConfigs() *Configuration {
	if PRODUCTION {
		CONFIGPATH = PRODCONFIGPATH
	}
	// if len(os.Args) < 2 {
	// 	panic("ARGS config path required")
	// }
	// CONFIGPATH = os.Args[1]
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	absPath := filepath.Join(path, CONFIGPATH)

	config := &Configuration{}
	file, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}
	return config
}

func OpenDB() {
	CONFIG = CheckAndLoadConfigs()
	os.Setenv("TZ", "Asia/Manila")
	fmt.Println("----------------------\n Prod:  Initializing account database...")
	os.Setenv("TZ", "Asia/Manila")
	DBM = pg.Connect(&pg.Options{
		Addr:     CONFIG.Dbconfig.Addr,
		User:     CONFIG.Dbconfig.User,
		Password: CONFIG.Dbconfig.Password,
		Database: CONFIG.Dbconfig.Database,
	})
	var n int
	_, err := DBM.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Account database connected!\n--------------------------------------")
	}
}
