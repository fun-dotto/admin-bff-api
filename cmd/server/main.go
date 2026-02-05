package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/generated/external/announcement_api"
	"github.com/fun-dotto/api-template/internal/handler"
	"github.com/fun-dotto/api-template/internal/middleware"
	"github.com/fun-dotto/api-template/internal/repository"
	"github.com/fun-dotto/api-template/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oapimiddleware "github.com/oapi-codegen/gin-middleware"
	"google.golang.org/api/idtoken"
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

	router.Use(oapimiddleware.OapiRequestValidator(spec))
	router.Use(middleware.FirebaseAuth(authClient))

	announcementAPIURL := os.Getenv("ANNOUNCEMENT_API_URL")
	if announcementAPIURL == "" {
		log.Fatal("ANNOUNCEMENT_API_URL is required")
	}
	// 認証付きHTTPクライアントを作成
	announcementAPIAuthClient, err := idtoken.NewClient(ctx, announcementAPIURL)
	if err != nil {
		log.Fatal("Failed to create auth client:", err)
	}

	// 生成されたクライアントに認証付きHTTPクライアントを注入
	apiClient, err := announcement_api.NewClientWithResponses(
		announcementAPIURL,
		announcement_api.WithHTTPClient(announcementAPIAuthClient),
	)
	if err != nil {
		log.Fatal("Failed to create API client:", err)
	}

	// Initialize layers
	announcementRepo := repository.NewAnnouncementRepository(apiClient)
	announcementService := service.NewAnnouncementService(announcementRepo)

	// Register handlers
	h := handler.NewHandler(announcementService)
	api.RegisterHandlers(router, h)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
