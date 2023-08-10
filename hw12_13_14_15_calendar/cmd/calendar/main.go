package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	storages "github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/vpsitnik/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	var storage storages.Storager

	if config.Db.Type == "sql" {
		storage = sqlstorage.New(config.Db.Dsn, logg)
	} else {
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
