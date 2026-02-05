# Code Style & Conventions

- Go standard formatting
- Comments in Japanese
- Interface definitions in handler.go alongside the Handler struct
- Converter functions separated into converter.go
- Handler methods named to match oapi-codegen generated interface (e.g., `AnnouncementsV1List`)
- Middleware uses Gin context (`c.Set`/`c.Get`) and `context.WithValue` for passing auth data
- Private context key types (unexported structs) for context keys
- Error responses as `gin.H{"error": ...}`
