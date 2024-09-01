package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Joke struct {
	ID     int
	JokeID string
	Vote   int
	IP     string
}

func connectToDB() *sql.DB {
	libsqlURL := fmt.Sprintf("%s?authToken=%s", os.Getenv("DB_URL"), os.Getenv("DB_TOKEN"))
	db, err := sql.Open("libsql", libsqlURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return nil
	}
	return db
}

func dbInit() {
	db := connectToDB()
	if db != nil {
		defer db.Close()

		_, err := db.Exec(`
						CREATE TABLE IF NOT EXISTS jokes (
							id integer primary key,
							timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
							joke_id TEXT,
							vote INTEGER,
							ip TEXT
					);
			`)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func addJokeVote(jokeID string, vote int, ip string) {
	db := connectToDB()
	if db != nil {
		defer db.Close()

		_, err := db.Exec(`
			INSERT INTO jokes (joke_id, vote, ip)
			VALUES (?, ?, ?)
		`, jokeID, vote, ip)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getJokeAverageVote(jokeID string) (int, int) {
	db := connectToDB()
	if db != nil {
		defer db.Close()

		rows, err := db.Query("SELECT vote FROM jokes where joke_id = '" + jokeID + "'")
		if err != nil {
			log.Fatal(err)
			return 0, 0
		}
		defer rows.Close()

		vote_total := 0
		count := 0
		for rows.Next() {
			var vote int

			if err := rows.Scan(&vote); err == nil {
				vote_total += vote
				count++
			}
		}
		if count > 0 {
			return vote_total / count, count
		}
		return 0, 0
	}

	return 0, 0
}
