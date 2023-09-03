package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
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
		" and output = $2 limit 1", algorithmName, hash).Scan(&message)

	return message, err
}

func GetHashByMessage(algorithmName string, message string) (string, error) {
	db := DBInstance()

	hash := ""
	err := db.QueryRow("select output from hashes where "+
		" algorithm_id = (select id from algorithms where name = $1) "+
		" and input = $2 limit 1", algorithmName, message).Scan(&hash)

	return hash, err
}

func GetRainbowInputsFromOutputs(algorithmName string, outputs []string) (*[]string, error) {
	db := DBInstance()

	rows, err := db.Query("select input from rainbow where "+
		" algorithm_id = (select id from algorithms where name = $1) "+
		" and output = any($2)", algorithmName, pq.Array(outputs))

	if err != nil {
		return &[]string{}, err
	}

	var inputs []string
	var input string
	for rows.Next() {
		err := rows.Scan(&input)
		CheckErrorPanic(err)
		inputs = append(inputs, input)
	}

	return &inputs, err
}

func InsertMultipleRainbowEntries(hashesToInsert []HashesToInsertType) {
	db := DBInstance()

	sqlStr := "insert into rainbow (algorithm_id, input, output) values "
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

func DeleteTables() {
	log.Println("Dropping tables")
	query := `
	DROP TABLE IF EXISTS public.hashes;
	DROP TABLE IF EXISTS public.rainbow;
	DROP TABLE IF EXISTS public.algorithms;
`
	db := DBInstance()
	_, err := db.Query(query)
	if err == nil {
		log.Println("Tables dropped")
	} else {
		log.Println(err)
	}
}

func CreateTables() {
	log.Println("Creating tables")
	query := `
	CREATE TABLE IF NOT EXISTS  public.algorithms (
		id serial4 NOT NULL,
		"name" varchar(255) NOT NULL,
		CONSTRAINT algorithms_pkey PRIMARY KEY (id),
		CONSTRAINT unique_algorithm_name UNIQUE (name)
	);

	CREATE TABLE IF NOT EXISTS  public.hashes (
		id serial4 NOT NULL,
		algorithm_id int4 NOT NULL,
		"input" text NOT NULL,
		"output" text NOT NULL,
		CONSTRAINT hashes_pkey_1 PRIMARY KEY (id),
		CONSTRAINT unique_algorithm_input_1 UNIQUE (algorithm_id, input)
	);

	CREATE TABLE IF NOT EXISTS  public.rainbow (
		id serial4 NOT NULL,
		algorithm_id int4 NOT NULL,
		"input" text NOT NULL,
		"output" text NOT NULL,
		CONSTRAINT rainbow_pkey_1 PRIMARY KEY (id),
		CONSTRAINT unique_algorithm_input_2 UNIQUE (algorithm_id, input)
	);

	ALTER TABLE public.hashes DROP CONSTRAINT IF EXISTS hashes_algorithm_id_fkey;
	ALTER TABLE public.hashes ADD CONSTRAINT hashes_algorithm_id_fkey FOREIGN KEY (algorithm_id) REFERENCES public.algorithms(id);
	ALTER TABLE public.rainbow ADD CONSTRAINT rainbow_algorithm_id_fkey FOREIGN KEY (algorithm_id) REFERENCES public.algorithms(id);

	INSERT INTO public.algorithms ("name")
	VALUES ('sha256') ON CONFLICT DO NOTHING;
	
	INSERT INTO public.algorithms ("name")
	VALUES ('sha1') ON CONFLICT DO NOTHING;
	
	INSERT INTO public.algorithms ("name")
	VALUES ('md5') ON CONFLICT DO NOTHING;

`
	db := DBInstance()
	_, err := db.Query(query)
	if err == nil {
		log.Println("Tables created")
	} else {
		log.Println(err)
	}
}
