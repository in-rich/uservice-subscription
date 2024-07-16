package main

import (
	"database/sql"
	"fmt"
	subscription_pb "github.com/in-rich/proto/proto-go/subscription"
	"github.com/in-rich/uservice-subscription/config"
	"github.com/in-rich/uservice-subscription/migrations"
	"github.com/in-rich/uservice-subscription/pkg/dao"
	"github.com/in-rich/uservice-subscription/pkg/handlers"
	"github.com/in-rich/uservice-subscription/pkg/services"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.App.Postgres.DSN)))
	db := bun.NewDB(sqldb, pgdialect.New())

	defer func() {
		_ = db.Close()
		_ = sqldb.Close()
	}()

	err := db.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	countNoteEditsByAuthorDAO := dao.NewCountNoteEditsByAuthorRepository(db)
	createNoteEditDAO := dao.NewCreateNoteEditRepository(db)
	getLatestNoteEditByAuthorDAO := dao.NewGetLatestNoteEditByAuthorRepository(db)

	canUpdateNoteService := services.NewCanUpdateNoteService(countNoteEditsByAuthorDAO, createNoteEditDAO, getLatestNoteEditByAuthorDAO)

	canUpdateNoteHandler := handlers.NewCanUpdateNoteHandler(canUpdateNoteService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	defer func() {
		server.GracefulStop()
		_ = listener.Close()
	}()

	subscription_pb.RegisterCanUpdateNoteServer(server, canUpdateNoteHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
