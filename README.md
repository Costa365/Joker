# Joker

This is a web application that gets a random joke from [icanhazdadjoke](https://icanhazdadjoke.com) and displays it. The user can then vote how funny they find the joke.

It was created to practice/learn several technologies including HTMX, Go and Turso. It can be accessed at [joker.costa365.site](https://joker.costa365.site/).

## Prerequisites

- Create an account on [turso.tech](https://turso.tech/).
  - Create a database - you'll need the database URL and the token.
- Install Go or Docker

## Running The App

If you want to run it with Go directly, create the following environment variables:
- DB_URL=libsql://site.turso.io
- DB_TOKEN=token

```go run main.go```

If you prefer to use Docker, update the docker-compose file with the url of your Turso database and your Turso token.

```docker compose up```


