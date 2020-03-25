package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS groups (
			jid varchar(50) NOT NULL,
			name varchar(50) NULL,
			PRIMARY KEY (jid)
		) DEFAULT CHARSET=utf8`)
	if err != nil {
		return
	}
	_, err = db.Exec(
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
	return
}

func (model *groupModel) add() (err error) {
	_, err = db.Exec(
		`INSERT INTO groups (jid, name) VALUES (?, ?)`,
		model.JID, model.Name)
	return
}

func (model *groupModel) delete() (err error) {
	_, err = db.Exec(`DELETE FROM groups WHERE jid=?`, model.JID)
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

func (model *assignmentModel) query() (result []assignmentModel, err error) {
	if cond, _ := (&groupModel{JID: model.GroupJID}).isExist(); !cond {
		return nil, errors.New("invalid group jid")
	}
	rows, err := db.Query(`SELECT * FROM assignments WHERE group_jid = ?`,
		model.GroupJID)
	if err != nil {
		return
	}
	for rows.Next() {
		row := assignmentModel{}
		err = rows.Scan(&row.ID, &row.Subject, &row.Description,
			&row.Deadline, &row.GroupJID)
		if err != nil {
			return
		}
		result = append(result, row)

	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].deadlineDistance() < 0 {
			return false
		} else if result[j].deadlineDistance() < 0 {
			return true
		}
		return result[i].deadlineDistance() < result[j].deadlineDistance()
	})
	return

}

func (model *assignmentModel) add() (err error) {
	if cond, _ := (&groupModel{JID: model.GroupJID}).isExist(); !cond {
		return errors.New("invalid group jid")
	}
	model.adjustValues()
	_, err = db.Exec(
		`INSERT INTO assignments (subject, description,
		 deadline, group_jid) VALUES (?, ?, ?, ?)`, model.Subject,
		model.Description, model.Deadline, model.GroupJID)
	return
}

func (model *assignmentModel) delete() (err error) {
	_, err = db.Exec(
		`DELETE FROM assignments WHERE id = ? AND group_jid = ?`,
		model.ID, model.GroupJID)
	return
}

func (model *assignmentModel) isExist() (bool, error) {
	stmt := `SELECT id FROM assignments WHERE id = ? AND group_jid = ?`
	err := db.QueryRow(stmt, model.ID, model.GroupJID).Scan(&model.ID)
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
		n, ok := cnf.getDayByName(day)
		if !ok {
			model.Deadline = strings.TrimSpace(strings.Title(model.Deadline))
			return
		}
		deadlineDays = append(deadlineDays, strconv.Itoa(n))
	}
	model.Deadline = strings.Join(deadlineDays, ", ")
}

func (model *assignmentModel) humanReadableValues() {
	deadlines := strings.Split(model.Deadline, ",")
	var deadlineDays []string
	for _, day := range deadlines {
		n, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil || n <= 0 || n > 7 {
			return
		}
		name, ok := cnf.getNameByDay(n)
		if !ok {
			return
		}
		deadlineDays = append(deadlineDays, strings.Title(name))
	}
	model.Deadline = strings.Join(deadlineDays, ", ")
}

func (model *assignmentModel) deadlineDistance() int {
	// dist: 3  deadline: 4(+1) today: (3+1)
	deadlines := strings.Split(model.Deadline, ",")
	lowestDistance := -1
	for _, deadline := range deadlines {
		today := time.Now().Weekday()
		deadlineDay, err := strconv.Atoi(deadline)
		delta := (deadlineDay + 1) - (int(today) + 1)
		if err != nil {
			return -1
		} else if delta >= 0 && (delta < lowestDistance || lowestDistance < 0) {
			lowestDistance = delta
			continue
		}
		distance := 7 + delta
		if distance < lowestDistance || lowestDistance < 0 {
			lowestDistance = distance
		}
	}
	fmt.Println(model.ID, lowestDistance)
	return lowestDistance
}
