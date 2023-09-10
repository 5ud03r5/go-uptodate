package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/5ud03r5/uptodate/auth"
	"github.com/5ud03r5/uptodate/db"
	"github.com/5ud03r5/uptodate/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	secret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	auth.TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	auth.RefreshTokenAuth = jwtauth.New("HS256", []byte(refreshSecret), nil)
}

func main() {

	// Environment variables loading
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not found in the env")
	}

	dbUri := os.Getenv("DATABASE_URL")
	if dbUri == "" {
		log.Fatal("DATABASE_URL is not found in the env")
	}
	
	// Defining main router definition
    router := chi.NewRouter()

	// Middlewares
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	  }))
	
	// Modifes Logging output for server logs
	router.Use(middleware.Logger)

	// Redirects all double slash requests to single slash
    router.Use(middleware.CleanPath)
	
	// Rate limitting setup to limit:
	// 100 requests in total from all requests within a second
	router.Use(httprate.LimitAll(100, 1*time.Second))

	// 100 requests performed by 1 IP within a minute
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	// DB Connection
	err := db.DBClient(dbUri)
    if err != nil {
        log.Fatal(err)
    }

	// Indexes creation for defined collections
	// Definitions:
	applicationCollection := db.CollectionIndex{CollectionName: "applications", IndexType: "text", IndexField: "name"}
	
	// Index creation:
	db.CreateIndexes(applicationCollection)

    defer db.MongoDBClient.Disconnect(context.Background())

	// Routes:
	// Debugger, helps to track CPU performance etc
	router.Mount("/debug", middleware.Profiler())
	
	// New sub-router fo v1 handlers
	// This is a core router for the app
	v1Router := chi.NewRouter()
	// Here comes all handlers:

	// Applications route:
	// /v1/applications
	v1Router.Route("/applications", func(r chi.Router) {

		// Public routes
		r.Group(func(routerPublic chi.Router){
			routerPublic.Get("/{applicationName}", handlers.HandlerGetApplicationByName)
		})

		// Private routes
		r.Group(func(routerPrivate chi.Router){
			routerPrivate.Use(jwtauth.Verifier(auth.TokenAuth))
			routerPrivate.Use(authenticatorMiddleware)

			routerPrivate.Post("/", handlers.HandlerUpsertApplication)
			routerPrivate.Post("/subscribe/{applicationName}", handlers.HandlerSubscribeToApplication)
		})
	})

	// Users route:
	// User specific actions to login and register user
	// /v1/users
	v1Router.Route("/users", func(r chi.Router) {
		r.Post("/register", handlers.HandlerRegisterUser)
		r.Post("/login", handlers.HandlerLoginUser)
	})

	// Auth route:
	// Token handling specific actions
	// /v1/auth
	v1Router.Route("/auth", func(r chi.Router) {
		r.Post("/refresh", handlers.HandlerRefreshToken)
	})

	//
	// Mounting a router at the end of the handlers
	router.Mount("/v1", v1Router)

	// Defining server settings
    srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	// Running application
	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}