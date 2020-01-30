package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
	"v3Osm/endpoints/rest"
	"v3Osm/handlers/sqlDb"
	"v3Osm/pkg/graceful"
	"v3Osm/pkg/logger"
	"v3Osm/pkg/middlewares"
	"v3Osm/usecases/change"
)

func main() {
	cfg := loadConfig()

	lg := logger.New(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	db, err := sqlDb.Connect(cfg.DbConn, cfg.DbDriver)
	if err != nil {
		lg.Fatalf("DB connection failed %s", err.Error())
	}

	cSt := sqlDb.NewChangeStore(lg, db)

	chr := change.NewChanger(lg, cSt)

	restHandler := rest.New(lg, chr)

	srv := setupServer(cfg, lg, restHandler)
	lg.Infof("listening for requests on :8080...")
	if err := srv.ListenAndServe(); err != nil {
		lg.Fatalf("http server exited: %s", err)
	}
}

func setupServer(cfg config, lg logger.Logger, rest http.Handler) *graceful.Server {
	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(http.StripPrefix("/api", rest))

	handler := middlewares.WithRequestLogging(lg, router)
	handler = middlewares.WithRecovery(lg, handler)

	srv := graceful.NewServer(handler, cfg.GracefulTimeout, os.Interrupt)
	srv.Log = lg.Errorf
	srv.Addr = cfg.Addr
	return srv
}

type config struct {
	LogLevel        string
	LogFormat       string
	DbConn          string
	DbDriver        string
	GracefulTimeout time.Duration
	Addr            string
}

func loadConfig() config {
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("DB_CONN_LINE", "postgres://zetkolink:remmuh23@localhost/postgres")
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("GRACEFUL_TIMEOUT", 20*time.Second)
	viper.SetDefault("ADDR", ":8080")

	return config{
		LogLevel:        viper.GetString("LOG_LEVEL"),
		LogFormat:       viper.GetString("LOG_FORMAT"),
		DbConn:          viper.GetString("DB_CONN_LINE"),
		DbDriver:        viper.GetString("DB_DRIVER"),
		GracefulTimeout: viper.GetDuration("GRACEFUL_TIMEOUT"),
		Addr:            viper.GetString("ADDR"),
	}
}
