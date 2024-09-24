package main

import (
	"github.com/in-rich/lib-go/deploy"
	subscription_pb "github.com/in-rich/proto/proto-go/subscription"
	"github.com/in-rich/uservice-subscription/config"
	"github.com/in-rich/uservice-subscription/migrations"
	"github.com/in-rich/uservice-subscription/pkg/dao"
	"github.com/in-rich/uservice-subscription/pkg/handlers"
	"github.com/in-rich/uservice-subscription/pkg/services"
	"log"
)

func main() {
	log.Println("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer closeDB()

	log.Println("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	depCheck := func() map[string]bool {
		errDB := db.Ping()

		return map[string]bool{
			"CanUpdateNote": errDB == nil,
			"":              errDB == nil,
		}
	}

	countNoteEditsByAuthorDAO := dao.NewCountNoteEditsByAuthorRepository(db)
	createNoteEditDAO := dao.NewCreateNoteEditRepository(db)
	getLatestNoteEditByAuthorDAO := dao.NewGetLatestNoteEditByAuthorRepository(db)

	canUpdateNoteService := services.NewCanUpdateNoteService(countNoteEditsByAuthorDAO, createNoteEditDAO, getLatestNoteEditByAuthorDAO)

	canUpdateNoteHandler := handlers.NewCanUpdateNoteHandler(canUpdateNoteService)

	log.Println("Starting to listen on port", config.App.Server.Port)
	listener, server, health := deploy.StartGRPCServer(config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	subscription_pb.RegisterCanUpdateNoteServer(server, canUpdateNoteHandler)

	log.Println("Server started")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
