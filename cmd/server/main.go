package main

import (
	"fmt"
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
	db, closeDB := deploy.OpenDB(config.App.Postgres.DSN)
	defer closeDB()

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	countNoteEditsByAuthorDAO := dao.NewCountNoteEditsByAuthorRepository(db)
	createNoteEditDAO := dao.NewCreateNoteEditRepository(db)
	getLatestNoteEditByAuthorDAO := dao.NewGetLatestNoteEditByAuthorRepository(db)

	canUpdateNoteService := services.NewCanUpdateNoteService(countNoteEditsByAuthorDAO, createNoteEditDAO, getLatestNoteEditByAuthorDAO)

	canUpdateNoteHandler := handlers.NewCanUpdateNoteHandler(canUpdateNoteService)

	listener, server := deploy.StartGRPCServer(fmt.Sprintf(":%d", config.App.Server.Port), "subscription")
	defer deploy.CloseGRPCServer(listener, server)

	subscription_pb.RegisterCanUpdateNoteServer(server, canUpdateNoteHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
