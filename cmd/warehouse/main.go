package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/akrovv/warehouse/internal/adapters/postgresql"
	"github.com/akrovv/warehouse/internal/config"
	"github.com/akrovv/warehouse/internal/handlers/jsonrpc"
	"github.com/akrovv/warehouse/internal/services"
	"github.com/akrovv/warehouse/pkg/logger"
	_ "github.com/lib/pq"
)

const (
	configType = "yml"
	path       = "."
	filename   = "config.yml"
	openConns  = 10
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("can't initialize logger, %s", err.Error())
		return
	}

	cfg, err := config.NewConfig(configType, path, filename)
	if err != nil {
		logger.Infof("can't initialize config, %w", err)
		return
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.SslMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatalf("can't open database, %w", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(openConns)

	if err = db.Ping(); err != nil {
		logger.Fatalf("can't connect to database, %w", err)
		return
	}

	var (
		productStorage   = postgresql.NewProductStorage(db)
		warehouseStorage = postgresql.NewWarehouseStorage(db)
	)

	var (
		productService   = services.NewProductService(productStorage)
		warehouseService = services.NewWarehouseService(warehouseStorage)
	)

	server, err := jsonrpc.NewServer(productService, warehouseService, logger)

	if err != nil {
		return
	}

	logger.Infof("starting on :%d", cfg.Server.Port)
	if err = server.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		logger.Fatalf("can't start server on %d, with error: %w", cfg.Server.Port, err)
		return
	}
}
