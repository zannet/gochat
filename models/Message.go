package models

/**
* Message model
 */

import (
	"database/sql"
	"fmt"
	"log"
)

type Message struct {
	Model

	id         int
	From       int
	To         int
	Body       string
	CreateDate int64
	IsRead     bool
}

func NewMessage(conn *sql.DB) *Message {
	result := new(Message)

	result.SetConnection(conn)

	return result
}

func (m *Message) GetUndeliveredMessages(to int) []Message {
	var result []Message

	rows, err := m.GetConnection().Query("SELECT * FROM messages WHERE `to`=? AND is_delivered=0", to)

	if err == nil {
		for rows.Next() {
			var message Message
			if err := rows.Scan(&message.id, &message.From, &message.To, &message.Body, &message.CreateDate, &message.IsRead); err != nil {
				log.Fatal(err)
			} else {
				result = append(result, message)
			}
		}
	} else {
		log.Fatal(err)
	}

	return result
}

func (m *Message) RemoveUndeliveredMessages(to int) {
	m.GetConnection().Query("DELETE FROM messages WHERE `to`=?", to)
}

func (m *Message) Save() bool {
	stmt, err := m.GetConnection().Prepare("INSERT INTO messages (`from`, `to`, `message`, `atime`, `is_delivered`) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		fmt.Println(err)
	} else {
		_, err = stmt.Exec(m.From, m.To, m.Body, m.CreateDate, m.IsRead)
	}

	if err == nil {
		return false
	}

	return true
}
