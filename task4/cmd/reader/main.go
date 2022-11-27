package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq" // postgres driver
	"lab5/internal"
	"log"

	"os"
)

const dsnTemplate = "user=%s password=%s host=%s dbname=%s sslmode=disable"

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const (
	Host     = "PG_HOST"
	Username = "PG_USER"
	Password = "PG_PASSWORD"
	Database = "PG_DATABASE"
)

func main() {
	ch, q1, q2 := internal.SetupRMQ()
	msgs, err := ch.Consume(
		q1.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	internal.FailOnError(err, "failed to register a consumer")

	nots, err := ch.Consume(
		q2.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	internal.FailOnError(err, "failed to register a consumer")

	desc := fmt.Sprintf(dsnTemplate,
		os.Getenv(Username),
		os.Getenv(Password),
		os.Getenv(Host),
		os.Getenv(Database))
	log.Printf("connecting to Postgres at %s", desc)
	db, err := sql.Open("postgres", desc)
	internal.FailOnError(err, "failed to connect to database")
	err = db.Ping()
	internal.FailOnError(err, "failed to connect to database")

	initDB(db)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			appendToTable(db, "messages", string(d.Body))
		}
	}()
	go func() {
		for d := range nots {
			log.Printf("Received a notification: %s", d.Body)
			appendToTable(db, "notifications", string(d.Body))
		}
	}()

	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func initDB(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS messages(
    id serial PRIMARY KEY,
    data VARCHAR(255)
);
`
	_, err := db.Exec(query)
	internal.FailOnError(err, "failed to create messages table")

	query = `
	CREATE TABLE IF NOT EXISTS notifications(
    id serial PRIMARY KEY,
    data VARCHAR(255)
);
`
	_, err = db.Exec(query)
	internal.FailOnError(err, "failed to create notifications table")
}

func appendToTable(db *sql.DB, table string, data string) {
	query := psql.Insert(table).
		Columns("data").
		Values(data)

	_, err := query.RunWith(db).Exec()
	internal.FailOnError(err, fmt.Sprintf("failed to append to %s table", table))
	log.Printf("appended message '%s' to table %s", data, table)
}
