// this module test the connection with the database
// it first :-
// 1. Opens the database
// 2. m.Run() runs all the actual unit tests
// 2. Os.Exit(...) exists the test process with correct status code

package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbdriver = "mysql"
	dbsource = "root:password@tcp(127.0.0.1:3307)/simple_bank?parseTime=true"
)

var testqueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {

	// conn will have the pool of connections to MySql database
	var err error
	TestDB, err = sql.Open(dbdriver, dbsource)

	if err != nil {
		log.Fatal("Error in opening the database", err)
	}

	// New function requires the pool of connection with the database
	testqueries = New(TestDB)

	os.Exit(m.Run())

}
