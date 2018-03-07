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
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

func availableImages(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
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
	} else {
		err := "Wrong request method!"
		fmt.Printf("Request failed with %s\n", err)
	}
}

func insertDbEntry(w http.ResponseWriter, r *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	fmt.Println("# Inserting values")
	var lastInsertId int
	//var err error
	err = db.QueryRow("INSERT INTO test(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)
}

/*
 func routeTwo (w http.ResponseWriter, r *http.Request) {
 	fmt.Fprintf(w, "In route Zwei -> gewählte Route: %s", r.URL.Path[1:])
}
*/

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	db.QueryRow("CREATE TABLE test(ID int primary key, spalte1 VARCHAR(30) NOT NULL, spalte2 VARCHAR(30), spalte3 VARCHAR(30));")
	//checkErr(err)

	// fmt.Println("# Updating")
	// stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	// checkErr(err)
	//
	// res, err := stmt.Exec("astaxieupdate", lastInsertId)
	// checkErr(err)
	//
	// affect, err := res.RowsAffected()
	// checkErr(err)
	//
	// fmt.Println(affect, "rows changed")
	//
	// fmt.Println("# Querying")
	// rows, err := db.Query("SELECT * FROM userinfo")
	// checkErr(err)
	//
	// for rows.Next() {
	// 	var uid int
	// 	var username string
	// 	var department string
	// 	var created time.Time
	// 	err = rows.Scan(&uid, &username, &department, &created)
	// 	checkErr(err)
	// 	fmt.Println("uid | username | department | created ")
	// 	fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
	// }
	//
	// fmt.Println("# Deleting")
	// stmt, err = db.Prepare("delete from userinfo where uid=$1")
	// checkErr(err)
	//
	// res, err = stmt.Exec(lastInsertId)
	// checkErr(err)
	//
	// affect, err = res.RowsAffected()
	// checkErr(err)
	//
	// fmt.Println(affect, "rows changed")

	http.HandleFunc("/testserver/getAvailableImages", availableImages)
	http.HandleFunc("/testserver/insertDbEntry", insertDbEntry)
	// http.HandleFunc("/routeTwo", routeTwo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
