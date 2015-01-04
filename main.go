package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/iangudger/basicGoAPI/api"
	"github.com/iangudger/basicGoAPI/database"
	"github.com/iangudger/basicGoAPI/website"
)

func main() {
	log.Println("Starting server...")
	rand.Seed(time.Now().Unix())

	dbconn, err := database.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		api.Handler(w, r, dbconn)
	})

	// Catchall
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		website.Handler(w, r, dbconn)
	})

	// Start the server
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
