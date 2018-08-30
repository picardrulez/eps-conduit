package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

func startup() int {
	database, _ := sql.Open("sqlite3", "./epsconduit.db")
	routesTableCreate, err := database.Prepare("CREATE TABLE IF NOT EXISTS routes (destination TEXT, insertdate DATETIME)")
	if err != nil {
		log.Println("error preparing create routes table statement")
		log.Printf("%s", err)
		return 1
	}
	routesTableCreate.Exec()
	go dbCleaner("routes")

	errorsTableCreate, err := database.Prepare("CREATE TABLE IF NOT EXISTS errors (destination TEXT, errormsg TEXT, insertdate DATETIME)")
	if err != nil {
		log.Println("error preparing create error table statement")
		log.Printf("%s", err)
		return 1
	}
	errorsTableCreate.Exec()
	go dbCleaner("errors")

	return 0
}

func dbCleaner(table string) {
	log.Println("starting db cleaner on " + table)
	db, _ := sql.Open("sqlite3", "./epsconduit.db")
	for {
		stmt, err := db.Prepare("DELETE from " + table + " where insertdate <= date('now', '-1 day')")
		if err != nil {
			log.Println("error preparing delete statement for " + table)
			return
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Println("error executing delete from " + table)
			return
		}
		time.Sleep(1 * time.Minute)
	}
}
