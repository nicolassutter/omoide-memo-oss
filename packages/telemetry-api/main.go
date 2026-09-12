package main

import (
	"log"
	"net/http"
	"os"

	"omoidememo.com/telemetry-api/internal/config"
	"omoidememo.com/telemetry-api/internal/httpapi"
	"omoidememo.com/telemetry-api/internal/store/postgres"
	"omoidememo.com/telemetry-api/internal/telemetry"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "openapi" || os.Args[1] == "--openapi") {
		_, telemetryAPI := httpapi.NewRouter(config.Config{}, nil)
		yamlBytes, serializationError := telemetryAPI.OpenAPI().YAML()
		if serializationError != nil {
			log.Fatalf("failed to generate openapi yaml: %v", serializationError)
		}
		os.Stdout.Write(yamlBytes)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	db, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer sqlDB.Close()

	store := telemetry.NewStore(db)

	handler, _ := httpapi.NewRouter(cfg, store)
	log.Printf("Listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}
