package main

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/handlers"
	"ai-support-bot-service/middleware"
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)
	// Cors config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// Load .env
	godotenv.Load()

	// Init Google Client
	config.InitGemini(context.Background())
	// Init Chroma Client
	config.InitChroma()
	// Init MongoDB
	config.InitMongoDB()
	// Ping Ollama
	config.PingOllama()

	// Server Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Working fine!"))
	})

	// API Routes
	r.Route("/api", func(r chi.Router) {
		// Auth Routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.RegisterUserHandler)
			r.Post("/login", handlers.LoginHandler)
			r.Post("/logout", handlers.LogoutHandler)
		})
		// Application Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.UserAuthMiddleware)
			r.Get("/applications", handlers.GetApplicationsHandler)
			r.Route("/application", func(r chi.Router) {
				r.Post("/", handlers.CreateApplicationHandler)
				r.Get("/{id}", handlers.GetApplicationHandler)
				r.Put("/", handlers.UpdateApplicationHandler)

				// Uploading Files
				r.Post("/{id}/upload-file", handlers.UploadFileHandler)

				// Ingest
				r.Get("/{id}/ingest/{fileId}", handlers.IngestDataHandler)

				// Job Status SSE
				r.Get("/ingest-jobs/{job_id}", handlers.IngestJobStatusSSEHandler)
			})
		})
		// AI sdk main protocol API route
		r.Group(func(r chi.Router) {
			// r.Use(middleware.APIKeyAuthMiddleware)
			r.Post("/chat", handlers.ChatHandler)
		})
	})
	log.Println("Server started at port :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("Server Stopped due to the error:", err)
	}
}
