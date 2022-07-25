package utils

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/romanyx/polluter"
	"gopkg.in/yaml.v2"
)

func GetTestPGDB() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr:     fmt.Sprintf("%s:5432", os.Getenv("DB_HOST")),
	})
}

var db *sql.DB
var Teardown func() error

func InsertSeedData(seedFile string) {
	preparePostgresDB()
	p := polluter.New(polluter.PostgresEngine(db))
	seed, err := os.Open(seedFile)
	if err != nil {
		fmt.Printf("failed to open seed file: %s", err.Error())
	}
	defer seed.Close()
	if err := p.Pollute(seed); err != nil {
		fmt.Printf("failed to pollute: %s", err.Error())
		os.Exit(1)
	}
}

func CleanUpSeedData(seedFile string) {
	fmt.Println("============== CLEANUP =================")
	m := make(map[interface{}]interface{})
	seed, err := os.Open(seedFile)
	if err != nil {
		fmt.Printf("failed to open seed file: %s", err.Error())
	}
	defer seed.Close()
	s, _ := seed.Stat()
	data := make([]byte, s.Size())
	seed.Read(data)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("invalid seed file: %s", err.Error())
	}
	for k := range m {
		r, err := db.Exec(fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE;", k))
		if err != nil {
			fmt.Printf("error truncating: %s %v", err.Error(), r)
		}
	}
}

func preparePostgresDB() {
	var (
		host     = os.Getenv("DB_HOST")
		port     = 5432
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		pem      = os.Getenv("RDS_ROOT_CERT")
		sslmode  = "disable"
	)

	if pem != "" {
		rootCertPool := x509.NewCertPool()

		if ok := rootCertPool.AppendCertsFromPEM([]byte(pem)); !ok {
			log.Fatal("Failed to append PEM.")
		}

		mysql.RegisterTLSConfig("rds", &tls.Config{
			RootCAs: rootCertPool,
		})

		sslmode = "require"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	openDB, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("failed to open postgres connection: %s", err.Error())
	}

	db = openDB
	Teardown = db.Close
}

func RunTestWithSeed(testSeedMap map[string][]func(*testing.T), t *testing.T, tx *pg.Tx) {
	i := 0
	for seed, tests := range testSeedMap {
		InsertSeedData(seed)
		defer CleanUpSeedData(seed)
		for _, test := range tests {
			i++
			t.Run(GetTestFunctionName(test), func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						if tx != nil {
							tx.Rollback()
						}
						CleanUpSeedData(seed)
						panic(r)
					}
				}()
				test(t)
			})
		}
	}
}

func GetTestFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func getTestAllValuesEqual(v reflect.Value, value string) bool {
	allEqual := true
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			allEqual = allEqual && getTestAllValuesEqual(v.Index(i), value)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			allEqual = allEqual && getTestAllValuesEqual(v.MapIndex(k), value)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			allEqual = allEqual && getTestAllValuesEqual(v.Field(i), value)
		}
	default:
		if v.Kind() == reflect.String && fmt.Sprintf("%v", v.Interface()) != value {
			allEqual = false
		}
	}
	return allEqual
}
