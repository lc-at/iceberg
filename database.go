package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initiateDatabase() (err error) {
	db, err = sql.Open("mysql", cnf.DbConnectionString)
	if err != nil {
		return
	}
	err = createTable()
	return
}

func createTable() (err error) {
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS groups (
			jid varchar(50) NOT NULL,
			name varchar(50) NULL,
			PRIMARY KEY (jid)
		) DEFAULT CHARSET=utf8`)
	if err != nil {
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		return
	}
	stmt, err = db.Prepare(
		`CREATE TABLE IF NOT EXISTS assignments (
			id int NOT NULL AUTO_INCREMENT,
			subject varchar(10) NOT NULL,
			description text NOT NULL,
			group_jid varchar(50) NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (group_jid)
				REFERENCES groups(jid)
				ON DELETE CASCADE
		) DEFAULT CHARSET=utf8`)
	if err != nil {
		return
	}
	_, err = stmt.Exec()
	return
}

func addGroup(jid, name string) (err error) {
	stmt, err := db.Prepare(
		`INSERT INTO groups (jid, name) VALUES (?, ?)`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(jid, name)
	return
}

func deleteGroup(jid string) (err error) {
	stmt, err := db.Prepare(
		`DELETE FROM groups WHERE jid=?`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(jid)
	return
}

func groupExists(jid string) (bool, error) {
	stmt := `SELECT jid FROM groups WHERE jid= ?`
	err := db.QueryRow(stmt, jid).Scan(&jid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	checkError(err)
	return true, nil
}

func addAssignment(subject, desc, groupJid string) (err error) {
	if cond, _ := groupExists(groupJid); !cond {
		return errors.New("invalid group jid")
	}
	stmt, err := db.Prepare(
		`INSERT INTO sassignments (subject, description,
		 group_jid) VALUES (?, ?, ?)`)
	_, err = stmt.Exec(subject, desc, groupJid)
	return
}
