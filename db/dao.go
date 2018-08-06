package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/nettyrnp/go-grpc/util"
	"log"
)

var (
	DB *dbx.DB
)

func init() {
	// TODO: load config file
	var err error
	DB, err = dbx.MustOpen("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// Re-create the database schema. This method is used in tests.
func ResetDB() *dbx.DB {
	if err := runSQLFile(DB, getSQLFile()); err != nil {
		panic(fmt.Errorf("Error while initializing database: %s", err))
	}
	return DB
}

func SaveRecord(p util.Person) (int32, int32, error) {
	p2, err := getRecord(p)
	if err != nil || p2.Id == 0 {
		// didn't found a record with such ID, so insert it
		return 1, 0, insertRecord(p)
	}
	if p == p2 {
		return 0, 0, nil
	} else {
		return 0, 1, updateRecord(p)
	}
}

func getRecord(p util.Person) (util.Person, error) {
	sql := fmt.Sprintf("SELECT * FROM public.people WHERE id=%d", p.Id)
	var p2 util.Person
	err := DB.NewQuery(sql).One(&p2)
	return p2, err
}

func insertRecord(p util.Person) error {
	sql := fmt.Sprintf("INSERT INTO public.people(id, name, email, mobile_number) "+
		"VALUES (%d, '%s', '%s', '%s')", p.Id, p.Name, p.Email, p.MobileNumber)
	log.Printf("sql:", sql)
	if _, err := DB.NewQuery(sql).Execute(); err != nil {
		return err
	}
	return nil
}

func updateRecord(p util.Person) error {
	sql := fmt.Sprintf("UPDATE public.people "+
		"SET name='%s', email='%s', mobile_number='%s' WHERE id=%d",
		p.Name, p.Email, p.MobileNumber, p.Id)
	log.Printf("sql:", sql)
	if _, err := DB.NewQuery(sql).Execute(); err != nil {
		return err
	}
	return nil
}

func getSQLFile() string {
	if _, err := os.Stat("db/migration.sql"); err == nil {
		return "db/migration.sql"
	}
	return "../db/migration.sql"
}

func runSQLFile(db *dbx.DB, file string) error {
	s, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(s), ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if _, err := db.NewQuery(line).Execute(); err != nil {
			fmt.Println(line)
			return err
		}
	}
	return nil
}
