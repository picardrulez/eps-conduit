package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func insertRoute(destination string) int {
	db, err := sql.Open("sqlite3", "./epsconduit.db?_busy_timeout=15000")
	if err != nil {
		log.Println("error opening db for insert")
		log.Printf("%s", err)
		db.Close()
		return 1
	}
	stmt, err := db.Prepare("INSERT INTO routes(destination, insertdate) values(?, CURRENT_TIMESTAMP)")
	if err != nil {
		log.Println("error preparing insert statemnt")
		log.Printf("%s", err)
		db.Close()
		return 2
	}
	_, err = stmt.Exec(destination)
	if err != nil {
		log.Println("error executing insert statement")
		log.Printf("%s", err)
		db.Close()
		return 3
	}
	db.Close()
	return 0
}

func selectNumLastMinute(table string) int {
	db, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening db for selectNumLastMinute")
		log.Printf("%s", err)
		db.Close()
		return 1
	}
	res, err := db.Query("SELECT count(*) FROM " + table + " where datetime(insertdate) >= datetime('now', '-1 Minute')")
	if err != nil {
		log.Println("error querying last minute")
		log.Printf("%s", err)
		res.Close()
		db.Close()
	}
	count := rowCounter(res)
	db.Close()
	res.Close()
	return count
}

func selectNumLastFiveMinutes(table string) int {
	db, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening db for selectNumLastFiveMinutes")
		log.Printf("%s", err)
		db.Close()
		return 1
	}
	res, err := db.Query("SELECT count(*) FROM " + table + " where datetime(insertdate) >= datetime('now', '-5 Minute')")
	if err != nil {
		db.Close()
		res.Close()
		log.Println("error querying last 5 minutes")
	}
	count := rowCounter(res)
	db.Close()
	res.Close()
	return count
}

func selectNumLastFifteenMinutes(table string) int {
	db, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening db for selectNumLastFifteenMinutes")
		log.Printf("%s", err)
		db.Close()
		return 1
	}
	res, err := db.Query("SELECT count(*) FROM " + table + " where datetime(insertdate) >= datetime('now', '-10 Minute')")
	if err != nil {
		log.Println("error querying last 15 minutes")
		log.Printf("%s", err)
		db.Close()
		res.Close()
	}
	count := rowCounter(res)
	db.Close()
	res.Close()
	return count
}

func selectNumLastHour(table string) int {
	db, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening db for selectNumLastHour")
		log.Printf("%s", err)
		db.Close()
		return 1
	}
	res, err := db.Query("SELECT count(*) FROM " + table + " where datetime(insertdate) >= datetime('now', '-1 Hour')")
	if err != nil {
		log.Println("error querying db")
		log.Printf("%s", err)
		db.Close()
		res.Close()
	}
	count := rowCounter(res)
	db.Close()
	res.Close()
	return count
}

func insertError(destination string, connerr error) int {
	db, err := sql.Open("sqlite3", "./epsconduit.db")
	if err != nil {
		log.Println("error opening db for insertError")
		log.Printf("%s", err)
		db.Close()
		return 1
	}
	stmt, err := db.Prepare("INSERT INTO errors(destination, errormsg, insertdate) values(?,?,CURRENT_TIMESTAMP)")
	if err != nil {
		log.Println("error preparing insertError statement")
		log.Printf("%s", err)
		db.Close()
		return 2
	}
	_, err = stmt.Exec(destination, connerr.Error())
	if err != nil {
		log.Println("error executing insertError statement")
		log.Printf("%s", err)
		db.Close()
		return 3
	}
	db.Close()
	return 0
}

func rowCounter(rows *sql.Rows) (count int) {
	for rows.Next() {
		_ = rows.Scan(&count)
	}
	return count
}
