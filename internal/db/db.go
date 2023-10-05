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
	SigningKey  string
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

func SelectSingleUser(address string, db *sql.DB) (User, error) {

	var user User

	row := db.QueryRow(selectSingleUser, address)

	err := row.Scan(&user.Chat_id, &user.Tg_name, &user.Eth_address, &user.SigningKey)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CheckAddr(addr string, db *sql.DB) bool {

	var amount int

	row := db.QueryRow(checkAddr, addr)

	err := row.Scan(amount)

	if err != nil {
		return false
	}

	if amount == 0 {
		return false
	} else {
		return true
	}
}

func CheckNumberAddr(tgNick string, number int, db *sql.DB) bool {

	var amount int

	row := db.QueryRow(checkNumberAddr, tgNick)

	err := row.Scan(&amount)
	if err != nil {
		return false
	}

	if amount >= number {
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

		err = rows.Scan(&user.Chat_id, &user.Tg_name, &user.Eth_address, &user.SigningKey)
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

func InsertIntoTable(chat_id int64, tgNick string, eth_address string, signingKey string, db *sql.DB) error {
	stmt, err := db.Prepare(insertIntoTable)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(chat_id, tgNick, eth_address, signingKey)

	if err != nil {
		return err
	}

	return nil
}

const create string = `
  CREATE TABLE IF NOT EXISTS users (
  CHAT_ID INTEGER NOT NULL,
  TG_NAME TEXT NOT NULL,
  ADDRESS TEXT UNIQUE,
  SIGNING_KEY TEXT NOT NULL UNIQUE
  );`

const selectSingleUser string = `
Select *
From (
  Select CHAT_ID, TG_NAME, lower(ADDRESS) as addr, SIGNING_KEY
  From users
)
Where addr = ?
  `
const insertIntoTable string = `
INSERT INTO users VALUES ( ?, ?, ?, ?)
`

const selectUsers string = `
  Select *
  From (
	Select CHAT_ID, TG_NAME, lower(ADDRESS) as addr, SIGNING_KEY
	From users
  )
  Where addr = ?
  `

const checkAddr string = `
  Select COUNT(*)
  From (
	Select TG_NAME, lower(ADDRESS) as addr
	From users)
  Where addr = ?
  `

const checkNumberAddr string = `
  Select COUNT(*)
  From users
  Where TG_NAME = ?
  `
