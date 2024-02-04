package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)


var testQueries *Queries
var testStore *Store
var testDB *sql.DB


func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	

	// testQueries = New(conn)
	testStore = NewStore(conn) // Use NewStore to create a *Store instance

	

	os.Exit(m.Run())



}