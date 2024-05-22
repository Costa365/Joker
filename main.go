package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type JokeResponse struct {
	Joke string `json:"joke"`
	Id   string `json:"id"`
}

type JokeData struct {
	Joke string
	Id   string
}

func getJoke() (string, string) {
	url := "https://icanhazdadjoke.com/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "API error: " + err.Error(), ""
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "API error: " + err.Error(), ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "API error: " + resp.Status, ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response body", ""
	}

	var jokeResp JokeResponse
	err = json.Unmarshal(body, &jokeResp)
	if err != nil {
		return "Error decoding JSON from API: " + err.Error(), ""
	}

	return string(jokeResp.Joke), jokeResp.Id
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return IPAddress
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("IP:", ReadUserIP(r))
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
		joke, id := getJoke()
		data := JokeData{
			Joke: joke,
			Id:   id,
		}
		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/fragments/vote.html")

		if err != nil {
			log.Println(err)
			return
		}
		if t == nil {
			panic("template.ParseFiles returned nil")
		}
		id := r.PostFormValue("id")
		strvote := r.PostFormValue("vote")
		vote, err := strconv.Atoi(strvote)
		addJokeVote(id, vote, ReadUserIP(r))

		err = t.Execute(w, vote)
		if err != nil {
			panic(err)
		}
	})

	dbInit()

	log.Println("Starting server on :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
