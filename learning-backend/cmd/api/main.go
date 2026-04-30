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
	httpAdapter "github.com/kittitadkp/jeeb-learning/internal/adapter/in/http"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/handler"
	"github.com/kittitadkp/jeeb-learning/internal/adapter/in/http/middleware"
	mongoAdapter "github.com/kittitadkp/jeeb-learning/internal/adapter/out/mongo"
	"github.com/kittitadkp/jeeb-learning/internal/config"
	"github.com/kittitadkp/jeeb-learning/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.Log.Level),
	})))

	mongoClient, err := mongoAdapter.NewClient(cfg.MongoDB.URI)
	if err != nil {
		slog.Error("failed to connect to mongodb", "error", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB.Database)

	userRepo := mongoAdapter.NewUserRepository(db)
	topicRepo := mongoAdapter.NewTopicRepository(db)
	itemRepo := mongoAdapter.NewItemRepository(db)
	progressRepo := mongoAdapter.NewProgressRepository(db)

	userUC := usecase.NewUserUseCase(userRepo)
	topicUC := usecase.NewTopicUseCase(topicRepo)
	itemUC := usecase.NewItemUseCase(itemRepo)
	progressUC := usecase.NewProgressUseCase(progressRepo, topicRepo, itemRepo)

	var verifier *oidc.IDTokenVerifier
	if !cfg.UpstreamAuth {
		provider, err := oidc.NewProvider(context.Background(),
			fmt.Sprintf("%s/realms/%s", cfg.Keycloak.URL, cfg.Keycloak.Realm))
		if err != nil {
			slog.Error("failed to init oidc provider", "error", err)
			os.Exit(1)
		}
		verifier = provider.Verifier(&oidc.Config{
			ClientID:          cfg.Keycloak.ClientID,
			SkipClientIDCheck: true,
		})
	}
	authMiddleware := middleware.NewAuthMiddleware(verifier, userUC, cfg.UpstreamAuth)

	handlers := httpAdapter.Handlers{
		User:     handler.NewUserHandler(),
		Topic:    handler.NewTopicHandler(topicUC),
		Item:     handler.NewItemHandler(itemUC),
		Progress: handler.NewProgressHandler(progressUC, itemUC),
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

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
