package routes

import (
	"encoding/json"
	"github.com/dhghf/mcauth/internal/bot"
	"github.com/dhghf/mcauth/internal/common"
	"github.com/dhghf/mcauth/internal/common/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	Bot    *bot.Bot
	Store  *db.Store
	Config *common.WebServerConfig
}

func StartAllRoutes(bot *bot.Bot, store *db.Store, config *common.WebServerConfig, router *mux.Router) {
	server := &Server{
		Bot: bot, Store: store, Config: config,
	}

	router.Use(server.authenticator)

	// GET /verifyPlayer/{playerID}
	router.HandleFunc("/verify/{playerID}", server.verifyPlayer).
		Methods("GET")

	/* Alt Account Management */

	// GET /alts
	router.HandleFunc("/alts", server.getAllAlts).
		Methods("GET")

	// GET /alts/{owner name}
	router.HandleFunc("/alts/{owner}", server.getAltsOf).
		Methods("GET")

	// POST /alt/{owner name}
	router.HandleFunc("/alts/{owner}/{alt_name}", server.postAlt).
		Methods("POST")

	// DELETE /alt/{alt name}
	router.HandleFunc("/alts/{alt_name}", server.deleteAlt).
		Methods("DELETE")
}

func Ship(res http.ResponseWriter, response interface{}) {
	res.Header().Set("Content-Type", "application/json")
	serialized, err := json.Marshal(response)

	if err != nil {
		log.Printf("Failed to serialize\n%s\nbecause\n%s\n", response, err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = res.Write(serialized)

	if err != nil {
		log.Println("Failed to respond to a request\n" + err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
