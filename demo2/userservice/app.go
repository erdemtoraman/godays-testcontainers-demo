package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"godays-testcontainers-demo/userservice/api"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// run migrations
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS users 
			(id SERIAL PRIMARY KEY,  name TEXT NOT NULL UNIQUE)
			`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("running service")
	r := mux.NewRouter().StrictSlash(true)
	r.NewRoute().Path("/health").Methods(http.MethodGet).HandlerFunc(api.Health)

	r.NewRoute().Path("/users").Methods(http.MethodPost).HandlerFunc(api.PostUser(db))
	r.NewRoute().Path("/users/{id:[0-9]+}").Methods(http.MethodGet).HandlerFunc(api.GetUser(db))

	http.ListenAndServe(":"+os.Getenv("PORT"), r)

}
