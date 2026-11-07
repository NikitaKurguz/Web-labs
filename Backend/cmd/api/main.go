package main

import (
	"log"
	"primarte/internal/app/apiserver"
	"primarte/internal/app/store"
)

func main() {
	servCfg := apiserver.NewConfig()
	if err := servCfg.Load("configs/config.toml"); err != nil {
		log.Fatalf("server %v", err)
	}

	DbCfg := store.NewConfig()
	if err := DbCfg.Load("configs/config.toml"); err != nil {
		log.Fatalf("DataBase %v", err)
	}

	db, err := store.New(DbCfg.DataBase.ConnectionString())

	if err != nil {
		log.Fatalf("store %v", err)
	}
	defer db.Close()

	serv := apiserver.New(servCfg, db)
	log.Printf("starting server on %s", servCfg.Server.BindAddr)
	if err := serv.Start(); err != nil {
		log.Fatalf("server: %v", err)
	}
}
