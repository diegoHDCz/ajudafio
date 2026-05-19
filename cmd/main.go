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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/diegoHDCz/ajudafio/internal/infra/config"
	"github.com/diegoHDCz/ajudafio/internal/infra/database"
	"github.com/diegoHDCz/ajudafio/internal/shared"

	address "github.com/diegoHDCz/ajudafio/internal/address"
	addresshttp "github.com/diegoHDCz/ajudafio/internal/address/adapters/http"
	addresspostgres "github.com/diegoHDCz/ajudafio/internal/address/adapters/postgres"
	authhttp "github.com/diegoHDCz/ajudafio/internal/auth/adapters/http"
	keycloak "github.com/diegoHDCz/ajudafio/internal/auth/adapters/keycloak"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	availability "github.com/diegoHDCz/ajudafio/internal/availability"
	availabilityhttp "github.com/diegoHDCz/ajudafio/internal/availability/adapters/http"
	avalabilityRepo "github.com/diegoHDCz/ajudafio/internal/availability/adapters/postgres"
	professional "github.com/diegoHDCz/ajudafio/internal/professional"
	professionalhttp "github.com/diegoHDCz/ajudafio/internal/professional/adapters/http"
	professionalpostgres "github.com/diegoHDCz/ajudafio/internal/professional/adapters/postgres"
	user "github.com/diegoHDCz/ajudafio/internal/user"
	userhttp "github.com/diegoHDCz/ajudafio/internal/user/adapters/http"
	userpostgres "github.com/diegoHDCz/ajudafio/internal/user/adapters/postgres"
)

func main() {

	// ── Config ────────────────────────────────────────────────────────────────
	fmt.Println("Loading configuration...")
	cfg := config.Load()

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// // ── Migrations ────────────────────────────────────────────────────────────
	// if err := database.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
	// 	slog.Error("failed to run migrations", "error", err)
	// 	os.Exit(1)
	// }

	// ── Wire: user slice ──────────────────────────────────────────────────────
	userRepo := userpostgres.NewRepository(db)
	userSvc := user.NewService(userRepo)

	// ── Wire: shared validator ─────────────────────────────────────────────────
	validator := shared.NewValidator(userSvc)

	userHandler := userhttp.NewHandler(userSvc, validator)

	// ── Wire: auth slice ──────────────────────────────────────────────────────
	authRepo := keycloak.NewKeycloakRepository("http://localhost:8080")
	config, _ := authRepo.GetKeycloakConfig()
	authSvc := authhttp.NewHandler(*authRepo, &config, userSvc)

	// ── Wire: middleware request ──────────────────────────────────────────────────────
	authMW, err := authmiddleware.NewAuthMiddleware(context.Background(), "http://localhost:8180/realms/ajudafio", "app-ajudafio")
	if err != nil {
		slog.Error("failed to initialize auth middleware", "error", err)
		os.Exit(1)
	}

	// ── Wire: professional slice ────────────────────────────────────────────────────────────────
	professionalRepo := professionalpostgres.NewRepository(db)
	professionalSvc := professional.NewProfessionalService(professionalRepo)
	professionalHandler := professionalhttp.NewProfessionalHandler(professionalSvc, validator)

	// ── Wire: avaliabilities slice ────────────────────────────────────────────────────────────────
	avalabilityRepo := avalabilityRepo.NewRepository(db)
	availabilitySvc := availability.NewAvailabilityService(avalabilityRepo)
	availabilityHandler := availabilityhttp.NewAvailabilityHandler(availabilitySvc, validator, professionalSvc)

	// ── Wire: address slice ────────────────────────────────────────────────────────────────
	addressRepo := addresspostgres.NewAddressRepository(db)
	addressSvc := address.NewAddressService(addressRepo)
	addressHandler := addresshttp.NewAddressHandler(addressSvc, validator)
	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Mount("/auth", authhttp.NewRouter(authSvc))

	r.Route("/professionals", func(r chi.Router) {
		r.Get("/", professionalHandler.FindWithFilters)
		r.Get("/{id}", professionalHandler.GetByID)
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequestAuth)
			r.Get("/user/{userID}", professionalHandler.GetByUserID)
			r.Post("/", professionalHandler.Create)
			r.Patch("/{id}", professionalHandler.Update)

			r.Delete("/{id}", professionalHandler.Delete)

		})
	})

	r.Group(func(r chi.Router) {
		r.Use(authMW.RequestAuth)
		r.Mount("/users", userhttp.NewRouter(userHandler))
		r.Mount("/availabilities", availabilityhttp.NewAvailabilityRouter(availabilityHandler))
		r.Mount("/addresses", addresshttp.NewRouter(addressHandler))
	})

	// ── Server ────────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
