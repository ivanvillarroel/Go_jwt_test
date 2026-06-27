package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ivanvillarroel/go_jwt_gin/internal/auth"
	"github.com/ivanvillarroel/go_jwt_gin/internal/config"
)

func TestListUsersRequiresJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server, err := NewServer(testConfig())
	if err != nil {
		t.Fatalf("expected server to start: %v", err)
	}
	defer server.DB.Close()

	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestLoginReturnsToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server, err := NewServer(testConfig())
	if err != nil {
		t.Fatalf("expected server to start: %v", err)
	}
	defer server.DB.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/auth/token", strings.NewReader(`{"username":"api-user","password":"password"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if body["token_type"] != "Bearer" {
		t.Fatalf("expected bearer token type, got %q", body["token_type"])
	}

	if body["token"] == "" {
		t.Fatal("expected token to be present")
	}
}

func TestListUsersReturnsSeededUsersWithJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server, err := NewServer(testConfig())
	if err != nil {
		t.Fatalf("expected server to start: %v", err)
	}
	defer server.DB.Close()

	token := loginAndGetToken(t, server)
	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var users []map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &users); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 seeded users, got %d", len(users))
	}

	if users[0]["email"] != "ada@example.com" {
		t.Fatalf("expected first user to be Ada, got %#v", users[0])
	}
}

func TestListUsersRejectsJWTWithoutRequiredPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server, err := NewServer(testConfig())
	if err != nil {
		t.Fatalf("expected server to start: %v", err)
	}
	defer server.DB.Close()

	tokenService := auth.NewJWTService(testConfig().JWTSecret)
	token, err := tokenService.Generate("api-user", []string{"profile:read"})
	if err != nil {
		t.Fatalf("expected token generation to succeed: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d: %s", http.StatusForbidden, response.Code, response.Body.String())
	}
}

func TestValidTokenRouteReturnsTokenStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server, err := NewServer(testConfig())
	if err != nil {
		t.Fatalf("expected server to start: %v", err)
	}
	defer server.DB.Close()

	token := loginAndGetToken(t, server)
	request := httptest.NewRequest(http.MethodGet, "/api/auth/valid", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if body["valid"] != true {
		t.Fatalf("expected token to be valid, got %#v", body)
	}

	if body["subject"] != "api-user" {
		t.Fatalf("expected subject %q, got %#v", "api-user", body["subject"])
	}
}

func TestReadTokenRouteReturnsSubject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server, err := NewServer(testConfig())
	if err != nil {
		t.Fatalf("expected server to start: %v", err)
	}
	defer server.DB.Close()

	token := loginAndGetToken(t, server)
	request := httptest.NewRequest(http.MethodGet, "/api/auth/read", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if body["subject"] != "api-user" {
		t.Fatalf("expected subject %q, got %#v", "api-user", body["subject"])
	}
}

func loginAndGetToken(t *testing.T, server *Server) string {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, "/api/auth/token", strings.NewReader(`{"username":"api-user","password":"password"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	server.Engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected login status %d, got %d", http.StatusOK, response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	return body["token"]
}

func testConfig() config.Config {
	return config.Config{
		Port:         "0",
		JWTSecret:    "test-secret",
		AuthUser:     "api-user",
		AuthPassword: "password",
	}
}
