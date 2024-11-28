package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/bmcszk/user-service/api"
	"github.com/bmcszk/user-service/db"
	"github.com/bmcszk/user-service/logic"

	"github.com/joho/godotenv"
)

func main() {
	// envs
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file provided")
	}
	postgresUrl := os.Getenv("POSTGRES_URL")
	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// db
	conn, err := db.InitDB(ctx, postgresUrl)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	defer conn.Close(ctx)
	queries := db.New(conn)
	// logic
	service := logic.NewService(queries)
	// api
	http.ListenAndServe(":8080", api.NewHandler(service))
	// TODO graceful shutdown
}
