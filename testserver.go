package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/lib/pq"
)

const (
	host		= "127.0.0.1"
	port		= "5432"
	user		= "postgres"
	password	= "postgres"
	dbname		= "postgres"
)

func createTable(w http.ResponseWriter, r *http.Request) {
	if checkRequestMethod(r.Method, "GET") {
		db, err := openDB()
		if noErr(err) {
			fmt.Printf("Creating Table..\n")
			_, err = db.Exec("CREATE TABLE IF NOT EXISTS test(ID SERIAL UNIQUE PRIMARY KEY, " +
				"spalte1 VARCHAR (30) NOT NULL, spalte2 VARCHAR(30) NOT NULL, spalte3 VARCHAR(30) NOT NULL);")
			if noErr(err) {
				fmt.Printf("Table created\n")
			}
		}
		closeDB(db)
	}
}

func availableImages(w http.ResponseWriter, r *http.Request) {
	if checkRequestMethod(r.Method, "GET") {
		cmd := exec.Command("docker", "images")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	}
}

func insertDbEntry(w http.ResponseWriter, r *http.Request) {
	if checkRequestMethod(r.Method, "GET") {
		db, err := openDB()
		//checkDbConnection(db)
		if noErr(err) {
			fmt.Println("# Inserting values..")
			var lastInsertId int
			err = db.QueryRow("INSERT INTO test(spalte1 , spalte2, spalte3) VALUES($1,$2,$3) returning ID;",
				"astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
			if noErr(err) {
				fmt.Println("# Inserted value =", lastInsertId)
			}
			closeDB(db)
		}
	}
}


func main() {
	http.HandleFunc("/testserver/getAvailableImages", availableImages)
	http.HandleFunc("/testserver/insertDbEntry", insertDbEntry)
	http.HandleFunc("/testserver/createTable", createTable)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func openDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}

func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func noErr(err error) bool {
	if err != nil {
		fmt.Printf("Failed with %s \n", err)
		return false
	} else {
		return true
	}
}

func checkRequestMethod(usedMethod string, method string) bool{
	if usedMethod == method{
		return true
	} else {
		err := "Wrong request method!"
		fmt.Printf("Request failed with %s\n", err)
		return false
	}
}
