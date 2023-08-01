package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var dbInstance *sql.DB = nil

func createDBInstance() {
	psqlconn := os.Getenv("DB_CONN")
	if psqlconn == "" {
		panic("DB connection string is empty")
	}

	log.Println("Connecting to DB...")

	var err error
	dbInstance, err = sql.Open("postgres", psqlconn)
	CheckErrorPanic(err)

	err = dbInstance.Ping()
	CheckErrorPanic(err)
	log.Println("Connected to DB")
}

func DBInstance() *sql.DB {
	if dbInstance == nil {
		createDBInstance()
	}

	return dbInstance
}

func GetAlgorithmIdFromName(algorithmName string) (int, error) {
	db := DBInstance()

	foundId := -1
	err := db.QueryRow("select id from algorithms where name = $1", algorithmName).Scan(&foundId)

	return foundId, err
}

func InsertAlgorithmHash(algorithmId int, input string, output string) bool {
	db := DBInstance()

	insertedId := -1
	err := db.QueryRow(
		"insert into hashes (algorithm_id, input, output) values ($1, $2, $3) returning id",
		algorithmId, input, output).Scan(&insertedId)

	return err == nil && insertedId != -1
}

func InsertMultipleAlgorithmHashes(hashesToInsert []HashesToInsertType) {
	db := DBInstance()

	sqlStr := "insert into hashes (algorithm_id, input, output) values "
	vals := []interface{}{}

	valId := 1
	for _, row := range hashesToInsert {
		sqlStr += fmt.Sprintf("($%d, $%d, $%d),", valId, valId+1, valId+2)
		vals = append(vals, row.Id, row.Input, row.Output)
		valId += 3
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += " on conflict (algorithm_id, input) do nothing"

	stmt, _ := db.Prepare(sqlStr)
	stmt.Exec(vals...)

	stmt.Close()
}

func GetMessageByHash(algorithmName string, hash string) (string, error) {
	db := DBInstance()

	message := ""
	err := db.QueryRow("select input from hashes where "+
		" algorithm_id = (select id from algorithms where name = $1) "+
		" and output = $2", algorithmName, hash).Scan(&message)

	return message, err
}

func GetHashByMessage(algorithmName string, message string) (string, error) {
	db := DBInstance()

	hash := ""
	err := db.QueryRow("select output from hashes where "+
		" algorithm_id = (select id from algorithms where name = $1) "+
		" and input = $2", algorithmName, message).Scan(&hash)

	return hash, err
}
