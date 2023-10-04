package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Chat_id     int64
	Tg_name     string
	Eth_address string
}

func Connect(dbconn string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbconn)

	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec(create)

	if err != nil {
		return nil, err
	}

	return sqlDB, nil

}

func SelectSingleUser(tgNick string, db *sql.DB) (User, error) {

	var user User

	row := db.QueryRow(selectSingleUser, tgNick)

	err := row.Scan(&user.Chat_id, &user.Tg_name, &user.Eth_address)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CheckAddr(addr string, db *sql.DB) bool {

	var amount int

	row := db.QueryRow(checkAddr, addr)

	err := row.Scan(&amount)

	log.Println(addr)
	log.Println(amount)
	log.Println(err)

	if err != nil {
		return false
	}

	if amount == 0 {
		return false
	} else {
		return true
	}
}

func SelectUsers(addressFrom string, addressTo string, db *sql.DB) ([]User, error) {

	var users []User
	var rows *sql.Rows
	var addr string

	if CheckAddr(addressFrom, db) {
		addr = addressFrom
	} else {
		addr = addressTo
	}

	rows, err := db.Query(selectUsers, addr)

	for rows.Next() {

		var user User

		err = rows.Scan(&user.Chat_id, &user.Tg_name, &user.Eth_address)
		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func InsertIntoTable(chat_id int64, tgNick string, eth_address string, db *sql.DB) error {
	stmt, err := db.Prepare(insertIntoTable)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(chat_id, tgNick, eth_address)

	if err != nil {
		return err
	}

	return nil
}

const create string = `
  CREATE TABLE IF NOT EXISTS users (
  CHAT_ID INTEGER NOT NULL PRIMARY KEY,
  TG_NAME TEXT NOT NULL,
  ADDRESS TEXT
  );`

const selectSingleUser string = `
  Select *
  From users
  Where TG_NAME = ?
  `
const insertIntoTable string = `
INSERT INTO users VALUES ( ?, ?, ?)
`

const selectUsers string = `
  Select *
  From users
  Where ADDRESS = ?
  `

const checkAddr string = `
  Select COUNT(*)
  From users
  Where ADDRESS = ?
  `
