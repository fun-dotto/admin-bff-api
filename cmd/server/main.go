package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/internal/handler"
	"github.com/fun-dotto/admin-bff-api/internal/infrastructure"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oapi-codegen/gin-middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to get Firebase Auth client: %v", err)
	}

	spec, err := openapi3.NewLoader().LoadFromFile("openapi/openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	spec.Servers = nil

	router := gin.Default()

	router.Use(ginmiddleware.OapiRequestValidatorWithOptions(spec, &ginmiddleware.Options{
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			if authStatusCode, authMessage, ok := middleware.GetAuthenticationError(c); ok {
				c.AbortWithStatusJSON(authStatusCode, gin.H{"error": authMessage})
				return
			}
			c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: middleware.FirebaseAuthenticationFunc(authClient),
		},
	}))

	clients, err := infrastructure.NewExternalClients(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize external clients: %v", err)
	}

	h := handler.NewHandler(clients.Academic, clients.Announcement, clients.Funch, clients.User)
	api.RegisterHandlers(router, h)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
