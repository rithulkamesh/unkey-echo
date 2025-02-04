package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kordlab/marketplace/config"
	"github.com/kordlab/marketplace/data"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	testConfig *config.Config
	testMongo  *data.MongoDB
	testRedis  *data.RedisDB
)

func setupTestDB(t *testing.T) func() {
	// Setup test config
	testConfig = config.LoadConfig()

	// Connect to test MongoDB
	var err error
	testMongo, err = data.NewMongoDB(testConfig)
	require.NoError(t, err)

	// Create indexes after connecting
	ctx := context.Background()
	_, err = testMongo.Users().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	require.NoError(t, err)

	// Connect to test Redis
	testRedis, err = data.NewRedisDB(testConfig)
	require.NoError(t, err)

	// Return cleanup function
	return func() {
		ctx := context.Background()
		// Drop test database
		err := testMongo.Client().Database(testConfig.DatabaseName).Drop(ctx)
		require.NoError(t, err)

		// Clear Redis database
		err = testRedis.Client().FlushDB(ctx).Err()
		require.NoError(t, err)

		// Close connections
		require.NoError(t, testMongo.Close(ctx))
		require.NoError(t, testRedis.Close())
	}
}

func clearCollections(t *testing.T) {
	_, err := testMongo.Users().DeleteMany(context.Background(), bson.M{})
	require.NoError(t, err)
}

func TestHandleRegister(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := NewAuthHandler(testMongo, testRedis)

	tests := []struct {
		name           string
		registerReq    RegisterRequest
		setupData      func(t *testing.T)
		expectedStatus int
		checkResult    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Registration",
			registerReq: RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				clearCollections(t)
			},
			expectedStatus: http.StatusCreated,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "testuser", response["username"])

				// Verify user in database
				var user data.User
				err = testMongo.Users().FindOne(context.Background(), bson.M{"email": "test@example.com"}).Decode(&user)
				require.NoError(t, err)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name: "Email Already Exists",
			registerReq: RegisterRequest{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				clearCollections(t)
				// Insert existing user
				_, err := testMongo.Users().InsertOne(context.Background(), &data.User{
					Username: "existinguser",
					Email:    "existing@example.com",
				})
				require.NoError(t, err)
			},
			expectedStatus: http.StatusConflict,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], "Email already registered")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupData(t)

			reqBody, _ := json.Marshal(tt.registerReq)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.handleRegister(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResult != nil {
				tt.checkResult(t, rec)
			}
		})
	}
}

func TestHandleLogin(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := NewAuthHandler(testMongo, testRedis)

	// Create test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)
	testUser := &data.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Status:       data.UserStatusActive,
	}

	tests := []struct {
		name           string
		loginReq       LoginRequest
		setupData      func(t *testing.T)
		expectedStatus int
		checkResult    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Login",
			loginReq: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupData: func(t *testing.T) {
				clearCollections(t)
				_, err := testMongo.Users().InsertOne(context.Background(), testUser)
				require.NoError(t, err)
			},
			expectedStatus: http.StatusOK,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			},
		},
		{
			name: "Invalid Credentials",
			loginReq: LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpass",
			},
			setupData: func(t *testing.T) {
				clearCollections(t)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], "Invalid credentials")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupData(t)

			reqBody, _ := json.Marshal(tt.loginReq)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.handleLogin(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResult != nil {
				tt.checkResult(t, rec)
			}
		})
	}
}
