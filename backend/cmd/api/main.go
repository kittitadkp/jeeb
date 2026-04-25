package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	httpAdapter "github.com/kittitad/jeeb/internal/adapter/in/http"
	"github.com/kittitad/jeeb/internal/adapter/in/http/handler"
	"github.com/kittitad/jeeb/internal/adapter/in/http/middleware"
	mongoAdapter "github.com/kittitad/jeeb/internal/adapter/out/mongo"
	"github.com/kittitad/jeeb/internal/config"
	"github.com/kittitad/jeeb/internal/usecase"
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		// Use default logging for config errors
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logLevel := parseLogLevel(cfg.Log.Level)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	mongoClient, err := mongoAdapter.NewClient(cfg.MongoDB.URI)
	if err != nil {
		slog.Error("failed to connect to mongodb", "error", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB.Database)

	userRepo := mongoAdapter.NewUserRepository(db)
	workoutRepo := mongoAdapter.NewWorkoutRepository(db)
	studyRepo := mongoAdapter.NewStudyRepository(db)
	sleepRepo := mongoAdapter.NewSleepRepository(db)
	financeRepo := mongoAdapter.NewFinanceRepository(db)
	eventRepo := mongoAdapter.NewEventRepository(db)

	userUC := usecase.NewUserUseCase(userRepo)
	workoutUC := usecase.NewWorkoutUseCase(workoutRepo)
	studyUC := usecase.NewStudyUseCase(studyRepo)
	sleepUC := usecase.NewSleepUseCase(sleepRepo)
	financeUC := usecase.NewFinanceUseCase(financeRepo)
	eventUC := usecase.NewEventUseCase(eventRepo, nil)

	provider, err := oidc.NewProvider(context.Background(),
		fmt.Sprintf("%s/realms/%s", cfg.Keycloak.URL, cfg.Keycloak.Realm))
	if err != nil {
		slog.Error("failed to init oidc provider", "error", err)
		os.Exit(1)
	}
	// SkipClientIDCheck: Keycloak access tokens carry aud=["account"] by default,
	// not the client ID. We still verify signature, expiry, and issuer.
	verifier := provider.Verifier(&oidc.Config{
		ClientID:          cfg.Keycloak.ClientID,
		SkipClientIDCheck: true,
	})
	authMiddleware := middleware.NewAuthMiddleware(verifier, userUC)

	handlers := httpAdapter.Handlers{
		User:    handler.NewUserHandler(),
		Workout: handler.NewWorkoutHandler(workoutUC),
		Study:   handler.NewStudyHandler(studyUC),
		Sleep:   handler.NewSleepHandler(sleepUC),
		Finance: handler.NewFinanceHandler(financeUC),
		Event:   handler.NewEventHandler(eventUC),
	}

	router := httpAdapter.NewRouter(handlers, authMiddleware)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.Info("shutting down server")
	srv.Shutdown(ctx)
}
