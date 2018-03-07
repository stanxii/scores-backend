package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/raphi011/scores/sqlite"
)

type credentials struct {
	Cid     string `json:"client_id"`
	Csecret string `json:"client_secret"`
}

type app struct {
	db         *sql.DB
	conf       *oauth2.Config
	production bool
}

func main() {
	dbPath := flag.String("db", "./scores.db", "Path to sqlite db")
	gSecret := flag.String("goauth", "./client_secret.json", "Path to google oauth secret")
	flag.Parse()

	var redirectURL string
	var cred credentials
	production := os.Getenv("APP_ENV") == "production"
	file, err := ioutil.ReadFile(*gSecret)

	if err != nil {
		log.Printf("Client secret error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	if production {
		redirectURL = "https://scores.raphi011.com/api/auth"
	} else {
		redirectURL = "http://localhost:3000/api/auth"
	}

	db, err := sqlite.Open(*dbPath)

	if err != nil {
		log.Printf("DB error: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	app := app{
		production: production,
		conf: &oauth2.Config{
			ClientID:     cred.Cid,
			ClientSecret: cred.Csecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
		db: db,
	}

	router := initRouter(app)
	router.Run()
}