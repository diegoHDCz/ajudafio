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
	"github.com/go-chi/cors"

	_ "github.com/diegoHDCz/ajudafio/docs"
	"github.com/diegoHDCz/ajudafio/internal/infra/config"
	"github.com/diegoHDCz/ajudafio/internal/infra/database"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	httpswagger "github.com/swaggo/http-swagger/v2"

	address "github.com/diegoHDCz/ajudafio/internal/address"
	addresshttp "github.com/diegoHDCz/ajudafio/internal/address/adapters/http"
	addresspostgres "github.com/diegoHDCz/ajudafio/internal/address/adapters/postgres"
	appointment "github.com/diegoHDCz/ajudafio/internal/appointment"
	appointmenthttp "github.com/diegoHDCz/ajudafio/internal/appointment/adapters/http"
	appointmentpostgres "github.com/diegoHDCz/ajudafio/internal/appointment/adapters/postgres"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	availability "github.com/diegoHDCz/ajudafio/internal/availability"
	availabilityhttp "github.com/diegoHDCz/ajudafio/internal/availability/adapters/http"
	avalabilityRepo "github.com/diegoHDCz/ajudafio/internal/availability/adapters/postgres"
	bookingrequest "github.com/diegoHDCz/ajudafio/internal/bookingrequest"
	bookingrequesthttp "github.com/diegoHDCz/ajudafio/internal/bookingrequest/adapters/http"
	bookingrequestpostgres "github.com/diegoHDCz/ajudafio/internal/bookingrequest/adapters/postgres"
	contract "github.com/diegoHDCz/ajudafio/internal/contract"
	contracthttp "github.com/diegoHDCz/ajudafio/internal/contract/adapters/http"
	contractpostgres "github.com/diegoHDCz/ajudafio/internal/contract/adapters/postgres"
	professional "github.com/diegoHDCz/ajudafio/internal/professional"
	professionalhttp "github.com/diegoHDCz/ajudafio/internal/professional/adapters/http"
	professionalpostgres "github.com/diegoHDCz/ajudafio/internal/professional/adapters/postgres"
	s3provider "github.com/diegoHDCz/ajudafio/internal/storage/s3"
	user "github.com/diegoHDCz/ajudafio/internal/user"
	userhttp "github.com/diegoHDCz/ajudafio/internal/user/adapters/http"
	userpostgres "github.com/diegoHDCz/ajudafio/internal/user/adapters/postgres"
)

// @title			Ajudafio API
// @version		1.0
// @description	API para gerenciamento de profissionais de cuidado domiciliar
// @host			localhost:8080
// @BasePath		/
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				JWT token no formato: Bearer {token}
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
	storage := s3provider.New(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.AWSRegion, cfg.AWSS3BucketName)
	userSvc := user.NewService(userRepo, storage)

	// ── Wire: shared validator ─────────────────────────────────────────────────
	validator := shared.NewValidator(userSvc)

	userHandler := userhttp.NewHandler(userSvc, validator)

	// ── Wire: middleware request ──────────────────────────────────────────────────────
	authMW, err := authmiddleware.NewAuthMiddleware(context.Background(), cfg.ClerkJWKSURL)
	if err != nil {
		slog.Error("failed to initialize auth middleware", "error", err)
		os.Exit(1)
	}

	// ── Wire: professional slice ────────────────────────────────────────────────────────────────
	professionalRepo := professionalpostgres.NewRepository(db)
	professionalSvc := professional.NewProfessionalService(professionalRepo)
	professionalHandler := professionalhttp.NewProfessionalHandler(professionalSvc, userSvc, validator)

	// ── Wire: avaliabilities slice ────────────────────────────────────────────────────────────────
	avalabilityRepo := avalabilityRepo.NewRepository(db)
	availabilitySvc := availability.NewAvailabilityService(avalabilityRepo)
	availabilityHandler := availabilityhttp.NewAvailabilityHandler(availabilitySvc, validator, professionalSvc)

	// ── Wire: address slice ────────────────────────────────────────────────────────────────
	addressRepo := addresspostgres.NewAddressRepository(db)
	addressSvc := address.NewAddressService(addressRepo)
	addressHandler := addresshttp.NewAddressHandler(addressSvc, validator)

	// ── Wire: contract slice ────────────────────────────────────────────────────────────────
	contractRepo := contractpostgres.NewRepository(db)
	contractSvc := contract.NewContractService(contractRepo)
	contractHandler := contracthttp.NewContractHandler(contractSvc, validator)

	// ── Wire: appointment slice ───────────────────────────────────────────────────────────
	appointmentRepo := appointmentpostgres.NewRepository(db)
	appointmentSvc := appointment.NewAppointmentService(appointmentRepo, avalabilityRepo)
	appointmentHandler := appointmenthttp.NewHandler(appointmentSvc)

	// ── Wire: booking request slice ───────────────────────────────────────────────────────────
	bookingRequestRepo := bookingrequestpostgres.NewRepository(db)
	bookingRequestSvc := bookingrequest.NewBookingRequestService(bookingRequestRepo)
	bookingRequestHandler := bookingrequesthttp.NewHandler(bookingRequestSvc, validator, professionalSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/swagger/*", httpswagger.Handler(
		httpswagger.URL("/swagger/doc.json"),
	))

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

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Group(func(r chi.Router) {
			r.Use(authMW.RequestAuth)
			r.Get("/me", userHandler.Me)
			r.Get("/{id}", userHandler.GetByID)
			r.Patch("/{id}", userHandler.Update)
			r.Delete("/{id}", userHandler.Delete)
			r.Patch("/{id}/avatar", userHandler.UploadAvatar)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(authMW.RequestAuth)
		r.Mount("/availabilities", availabilityhttp.NewAvailabilityRouter(availabilityHandler))
		r.Mount("/addresses", addresshttp.NewRouter(addressHandler))
		r.Mount("/contracts", contracthttp.NewRouter(contractHandler))
		r.Mount("/appointments", appointmenthttp.NewRouter(appointmentHandler))
		r.Mount("/booking-requests", bookingrequesthttp.NewRouter(bookingRequestHandler))
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
