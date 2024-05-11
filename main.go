package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

type JokeResponse struct {
	Joke string `json:"joke"`
}

func getJoke() string {
	url := "https://icanhazdadjoke.com/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "API error: " + err.Error()
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "API error: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "API error: " + resp.Status
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response body"
	}

	var jokeResp JokeResponse
	err = json.Unmarshal(body, &jokeResp)
	if err != nil {
		return "Error decoding JSON from API: " + err.Error()
	}

	return string(jokeResp.Joke)
}

func main() {
	http.Handle("/static",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			log.Println(err)
			return
		}
		if t == nil {
			panic("template.ParseFiles returned nil")
		}
		err = t.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/joke", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/fragments/joke.html")
		if err != nil {
			log.Println(err)
			return
		}
		if t == nil {
			panic("template.ParseFiles returned nil")
		}
		err = t.Execute(w, getJoke())
		if err != nil {
			panic(err)
		}
	})
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
