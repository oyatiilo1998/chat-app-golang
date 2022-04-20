package main

import (
	"Study/websocket/api"
	"Study/websocket/config"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.Load()

	psqlConnString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)
	fmt.Println(psqlConnString)

	db, err := sqlx.Connect("postgres", psqlConnString)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Websocket")
	apiServer := api.New(db, cfg)
	apiServer.Run(":8080")
	for {}
}