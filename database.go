package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type groupModel struct {
	JID  string
	Name string
}

type assignmentModel struct {
	ID          int
	GroupJID    string
	Subject     string
	Description string
	Deadline    string
}

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
			deadline varchar(30) NOT NULL,
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

func (model *groupModel) add() (err error) {
	stmt, err := db.Prepare(
		`INSERT INTO groups (jid, name) VALUES (?, ?)`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(model.JID, model.Name)
	return
}

func (model *groupModel) delete() (err error) {
	stmt, err := db.Prepare(
		`DELETE FROM groups WHERE jid=?`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(model.JID)
	return
}

func (model *groupModel) isExist() (bool, error) {
	stmt := `SELECT jid FROM groups WHERE jid= ?`
	err := db.QueryRow(stmt, model.JID).Scan(&model.JID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	checkError(err)
	return true, nil
}

func (model *assignmentModel) adjustValues() {
	deadlines := strings.Split(model.Deadline, ",")
	var deadlineDays []string
	for _, day := range deadlines {
		day = strings.ToLower(strings.TrimSpace(day))
		n, ok := cnf.Days[day]
		if !ok {
			model.Deadline = strings.TrimSpace(strings.Title(model.Deadline))
			return
		}
		deadlineDays = append(deadlineDays, strconv.Itoa(n))
	}
	model.Deadline = strings.Join(deadlineDays, ",")
}

func (model *assignmentModel) add() (err error) {
	if cond, _ := (&groupModel{JID: model.GroupJID}).isExist(); !cond {
		return errors.New("invalid group jid")
	}
	stmt, err := db.Prepare(
		`INSERT INTO assignments (subject, description,
		 deadline, group_jid) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return
	}
	model.adjustValues()
	_, err = stmt.Exec(model.Subject, model.Description,
		model.Deadline, model.GroupJID)
	return
}

func (model *assignmentModel) delete() (err error) {
	stmt, err := db.Prepare(
		`DELETE FROM assignments WHERE id=?`)
	if err != nil {
		return
	}
	_, err = stmt.Exec(model.ID)
	return
}

func (model *assignmentModel) isExist() (bool, error) {
	stmt := `SELECT id FROM assignments WHERE id= ?`
	err := db.QueryRow(stmt, model.ID).Scan(&model.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	checkError(err)
	return true, nil
}

// TODO: LIST ASSIGNMENT
