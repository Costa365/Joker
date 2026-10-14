package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type rankedJoke struct {
	ID      string
	Text    string
	Average float64
	Votes   int
	Rank    int
}

func topJokesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	db := connectToDB()
	defer db.Close()
	rows, err := db.QueryContext(ctx, `
		SELECT joke_id, AVG(vote), COUNT(*) FROM jokes
		WHERE vote BETWEEN 1 AND 5 AND joke_id IS NOT NULL AND joke_id != ''
		GROUP BY joke_id ORDER BY AVG(vote) DESC, COUNT(*) DESC, joke_id ASC LIMIT 10`)
	if err != nil {
		log.Println("Leaderboard:", err)
		http.Error(w, "Could not load ratings. Please try Refresh.", http.StatusServiceUnavailable)
		return
	}
	var jokes []rankedJoke
	for rows.Next() {
		var joke rankedJoke
		if err = rows.Scan(&joke.ID, &joke.Average, &joke.Votes); err != nil {
			break
		}
		joke.Rank = len(jokes) + 1
		jokes = append(jokes, joke)
	}
	rowsErr := rows.Err()
	rows.Close()
	if err != nil || rowsErr != nil {
		log.Println("Leaderboard rows:", err, rowsErr)
		http.Error(w, "Could not load ratings. Please try Refresh.", http.StatusServiceUnavailable)
		return
	}

	// Fetch only the ten winning jokes, concurrently, within the request deadline.
	var wg sync.WaitGroup
	for i := range jokes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://icanhazdadjoke.com/j/"+url.PathEscape(jokes[i].ID), nil)
			if err != nil {
				return
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("User-Agent", "Joker (https://github.com/Costa365/Joker)")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var result JokeResponse
			if resp.StatusCode == http.StatusOK && json.NewDecoder(resp.Body).Decode(&result) == nil {
				jokes[i].Text = result.Joke
			}
		}(i)
	}
	wg.Wait()

	t, err := template.New("top-jokes.html").Funcs(template.FuncMap{
		"score":   func(value float64) string { return fmt.Sprintf("%.2f", value) },
		"percent": func(value float64) string { return fmt.Sprintf("%.2f", value*20) },
	}).ParseFiles("templates/fragments/top-jokes.html")
	if err != nil {
		log.Println("Leaderboard template:", err)
		http.Error(w, "Could not display leaderboard.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err := t.Execute(w, jokes); err != nil {
		log.Println("Leaderboard render:", err)
	}
}
