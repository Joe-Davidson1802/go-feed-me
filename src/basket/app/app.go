package app

import (
	"basket/hothandler"
	"basket/types"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var inmemorybasket = map[string]types.Restaurant{}

func Run(env string) {
	var c types.Config

	config := "local.toml"

	if env != "" {
		config = env + ".toml"
	}

	fmt.Println("reading config...")

	if _, err := toml.DecodeFile(config, &c); err != nil {
		log.Fatal(err)
	}

	fmt.Println("config loaded.")

	r := mux.NewRouter().
		PathPrefix("/apps").
		Subrouter()

	r.Handle("/basket", hothandler.New(UpdateBasketHandler{Config: c})).Methods("PUT")
	r.Handle("/basket", hothandler.New(GetBasketHandler{Config: c})).Methods("GET")

	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"PUT", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"turbo-frame"})

	err := http.ListenAndServe(c.Listen, handlers.CORS(corsObj, headers, methods)(r))

	log.Fatal(err)
}