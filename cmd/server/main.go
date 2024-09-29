package main

import (
	"fmt"
	"github.com/in-rich/lib-go/deploy"
	"github.com/in-rich/lib-go/monitor"
	subscription_pb "github.com/in-rich/proto/proto-go/subscription"
	"github.com/in-rich/uservice-subscription/config"
	"github.com/in-rich/uservice-subscription/migrations"
	"github.com/in-rich/uservice-subscription/pkg/dao"
	"github.com/in-rich/uservice-subscription/pkg/handlers"
	"github.com/in-rich/uservice-subscription/pkg/services"
	"github.com/rs/zerolog"
	"os"
)

func getLogger() monitor.GRPCLogger {
	if deploy.IsReleaseEnv() {
		return monitor.NewGCPGRPCLogger(zerolog.New(os.Stdout), "uservice-subscription")
	}

	return monitor.NewConsoleGRPCLogger()
}

func main() {
	logger := getLogger()

	logger.Info("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Fatal(err, "failed to connect to database")
	}
	defer closeDB()

	logger.Info("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		logger.Fatal(err, "failed to migrate")
	}

	depCheck := deploy.DepsCheck{
		Dependencies: func() map[string]error {
			return map[string]error{
				"Postgres": db.Ping(),
			}
		},
		Services: deploy.DepCheckServices{
			"CanUpdateNote": {"Postgres"},
		},
	}

	countNoteEditsByAuthorDAO := dao.NewCountNoteEditsByAuthorRepository(db)
	createNoteEditDAO := dao.NewCreateNoteEditRepository(db)
	getLatestNoteEditByAuthorDAO := dao.NewGetLatestNoteEditByAuthorRepository(db)

	canUpdateNoteService := services.NewCanUpdateNoteService(countNoteEditsByAuthorDAO, createNoteEditDAO, getLatestNoteEditByAuthorDAO)

	canUpdateNoteHandler := handlers.NewCanUpdateNoteHandler(canUpdateNoteService, logger)

	logger.Info(fmt.Sprintf("Starting to listen on port %v", config.App.Server.Port))
	listener, server, health := deploy.StartGRPCServer(logger, config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	subscription_pb.RegisterCanUpdateNoteServer(server, canUpdateNoteHandler)

	logger.Info("Server started")
	if err := server.Serve(listener); err != nil {
		logger.Fatal(err, "failed to serve")
	}
}
