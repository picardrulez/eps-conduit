package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

func startup() int {
	database, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening database for startup procedures")
		log.Printf("%s", err)
		return 1
	}
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
	db, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening database for dbCleaner on + " + table)
		log.Printf("%s", err)
	}
	for {
		stmt, err := db.Prepare("DELETE from " + table + " where insertdate <= datetime('now', '-1 hour')")
		if err != nil {
			log.Println("error preparing delete statement for " + table)
			log.Println(err)
			return
		}
		stmtrtn, err := stmt.Exec()
		if err != nil {
			log.Println("error executing delete from " + table)
			log.Println(err)
			return
		}
		rowsAffected, _ := stmtrtn.RowsAffected()
		log.Println("DBCLEANER FOR TABLE " + table + " STATEMENT RETURN ROWS AFFECTED: " + strconv.FormatInt(rowsAffected, 10))
		time.Sleep(1 * time.Minute)
	}
}
