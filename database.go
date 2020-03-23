package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database, _ = sql.Open("sqlite3", cnf.DbFilename)

func createTable() {
	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS [group] (
			[jid] NVARCHAR(50) NOT NULL PRIMARY KEY,
			[name] NVARCHAR(50) NULL
		);
		CREATE TABLE IF NOT EXISTS [assignment] (
			[id] INTEGER NOT NULL PRIMARY_KEY AUTOINCREMENT,
			[subject] NVARCHAR(10) NOT NULL,
			[description] TEXT NOT NULL,
			[group_jid] NVARCHAR(50) NOT NULL,
			FOREIGN KEY(group_jid) REFERENCES group(jid)
		);`)
	if err != nil {
		log.Fatalln(err)
	}
	statement.Exec()
}

func addGroup(jid, name string) (err error) {
	statement, _ := database.Prepare(
		`INSERT INTO group VALUES (?, ?)`)
	_, err = statement.Exec(jid, name)
	return
}

func delGroup(jid string) (err error) {
	statement, _ := database.Prepare(
		`DELETE FROM assignment WHERE group_jid=?;
		 DELETE FROM group WHERE jid=?`)
	_, err = statement.Exec(jid, jid)
	return 
}

func groupExists(jid string) bool {
	_, err := database.Query(`SELECT jid FROM group WHERE jid=?`, jid)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func addAssignment(subject, desc, groupJid string) (err error) {
	if !groupExists(groupJid) {
		return errors.New("invalid group jid")
	}
	statement, err := database.Prepare(
		`INSERT INTO assignment (subject, description,
		 group_jid) VALUES (?, ?, ?)`)
	statement.Exec(subject, desc, groupJid)
	return
}
