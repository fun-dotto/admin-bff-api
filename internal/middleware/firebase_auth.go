package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type firebaseTokenKey struct{}
type rawIDTokenKey struct{}

// FirebaseTokenContextKey は Gin の context および context.Context に
// Firebase ID トークンの検証結果を格納するキーです。
var FirebaseTokenContextKey = firebaseTokenKey{}

// RawIDTokenContextKey は Gin の context および context.Context に
// 生のID トークン文字列を格納するキーです。外部APIへの転送に使用します。
var RawIDTokenContextKey = rawIDTokenKey{}

// FirebaseAuth は Authorization: Bearer <Firebase ID Token> を検証する Gin ミドルウェアです。
// 検証に成功すると、デコードされたトークン（*auth.Token）を context に格納して次のハンドラに渡します。
func FirebaseAuth(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be Bearer <token>",
			})
			return
		}

		idToken := strings.TrimSpace(parts[1])
		if idToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "ID token is required",
			})
			return
		}

		ctx := c.Request.Context()
		token, err := authClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			status, message := authErrorResponse(err)
			c.AbortWithStatusJSON(status, gin.H{"error": message})
			return
		}

		c.Set(FirebaseTokenContextKey, token)
		c.Set(RawIDTokenContextKey, idToken)
		ctx = context.WithValue(ctx, FirebaseTokenContextKey, token)
		ctx = context.WithValue(ctx, RawIDTokenContextKey, idToken)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// authErrorResponse は Firebase Auth の検証エラーから HTTP ステータスとメッセージを返します。
func authErrorResponse(err error) (int, string) {
	if auth.IsIDTokenExpired(err) {
		return http.StatusUnauthorized, "ID token has expired"
	}
	if auth.IsIDTokenInvalid(err) {
		return http.StatusUnauthorized, "Invalid ID token"
	}
	if auth.IsIDTokenRevoked(err) {
		return http.StatusUnauthorized, "ID token has been revoked"
	}
	if auth.IsUserDisabled(err) {
		return http.StatusForbidden, "User has been disabled"
	}
	return http.StatusUnauthorized, "Authentication failed"
}

// GetFirebaseToken は Gin の context から検証済みの Firebase トークンを取得します。
// FirebaseAuth ミドルウェア通過後のハンドラ内でのみ有効です。
func GetFirebaseToken(c *gin.Context) (*auth.Token, bool) {
	val, exists := c.Get(FirebaseTokenContextKey)
	if !exists {
		return nil, false
	}
	token, ok := val.(*auth.Token)
	return token, ok
}

// GetFirebaseUID は Gin の context から認証済みユーザーの UID を取得します。
// トークンが存在しない場合は空文字列を返します。
func GetFirebaseUID(c *gin.Context) string {
	token, ok := GetFirebaseToken(c)
	if !ok || token == nil {
		return ""
	}
	return token.UID
}

// GetFirebaseTokenFromContext は context.Context から Firebase トークンを取得します。
// ミドルウェアで c.Request に設定した context から取得する場合に使用します。
func GetFirebaseTokenFromContext(ctx context.Context) (*auth.Token, bool) {
	val := ctx.Value(FirebaseTokenContextKey)
	if val == nil {
		return nil, false
	}
	token, ok := val.(*auth.Token)
	return token, ok
}

// GetRawIDTokenFromContext は context.Context から生のIDトークン文字列を取得します。
// 外部APIへのリクエスト転送に使用します。
func GetRawIDTokenFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(RawIDTokenContextKey)
	if val == nil {
		return "", false
	}
	token, ok := val.(string)
	return token, ok
}
