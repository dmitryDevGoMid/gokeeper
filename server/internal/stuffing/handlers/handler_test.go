package handlers

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/asimencrypt"
	mock_asimencrypt "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/asimencrypt/mocks"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/jwttoken"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/security"
	mock_security "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/security/mocks"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/files"
	mock_files "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/files/mocks"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"
	userRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"
	mock_user "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var PublicClientKey string = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsDnFluV73ef/gz8jIUZO
/P3bZ6gyBn8KLFXaPnNh5ajuju4HnXZTJd42JWMWRyYFsA16And3Hf4w6AcHwS9E
TGxClV69LE/OPPbug0bITUWEMQgOEafW29BKBgmAtBEFL5M4OxbomCK7DQTnaBOD
VZqFmS87DZe5l8l31HZooowVIF+fhHjRsRqS3L7t5PLG1WWo606rtj2cuJtdz6df
z8/gSl/IQwMbeu6q0Gii8vna9Yxw1ENOhYWdF4bhTXjposefj4ICV8eLO5iqudV8
Unb3DTQNTJ6aRFfEaPbuKb3xS7wS63i1qFxh9cdyZGpxRzHoaRmVD4abH1kmtTVQ
CnOX95syTgJ4IwHBcgIlt1lpZoVWqQdSJQsCjM36Ax4I0H17vZyKZ7EXiXnpnk/X
HK4+sCQ9BqqfUgSxz7DpMJVGwgUCIYFO7Xgh0iDEjb0M1nBWZy6yc+SOSrOBmQuM
5dgytRhPaUnQcwPqERJClSwavM2C6lsBINGute+HPG2aS4mROedCu/iREGj+zGVB
BWNtim0/4iiqJtyGXyLHfS4vTeThHNjZnb+4FLq5UWdQ7gviL3m/YpcbVTuIYPjO
h6V98oFVtOsJhVoD04nu7X5X+nYu472qUeJ+cBVxPeC+xuPaOqj0CKdVntEh0KeR
X5l46qclblQnvpY7q9cpsyMCAwEAAQ==
-----END RSA PUBLIC KEY-----`

var PrivateClientKey string = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

// Возвращаем публичный rsa ключ
func GetRSAPublicKey(data []byte) (*rsa.PublicKey, error) { // Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse RSA")
	}

	// Parse the public key in PKIX format
	publick, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publick.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

// Возвращаем приватный rsa ключ
func GetRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse RSA")
	}

	// Parse the private key in PKCS#1 format
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil

}

func TestRefreshToken(t *testing.T) {
	// Define a type for the mock behavior function
	userData := &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"}

	// Generate a JWT token for the user
	token, err := jwttoken.SetToken(userData)
	if err != nil {
		t.Fatalf("Failed to SetToken by user data: %v", err)
	}

	// Generate a refresh token for the user
	tokenRefresh, err := jwttoken.CreateRefreshToken(userData)
	if err != nil {
		t.Fatalf("Failed to CreateRefreshToken by user data: %v", err)
	}

	// Define the test cases
	tests := []struct {
		user         *user.User
		name         string
		username     string
		body         string
		tokenRefresh string
		token        string
		statusCode   int
	}{
		{
			user:         &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			username:     "asd",
			name:         "RefreshToken_200",
			tokenRefresh: tokenRefresh,
			token:        token,
			body:         "",
			statusCode:   200,
		},
		{
			user:         &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			username:     "asd",
			name:         "RefreshToken_200",
			tokenRefresh: tokenRefresh + "_test",
			token:        token,
			body:         "",
			statusCode:   401,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the test controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Create a new mock user repository
			s := mock_user.NewMockUserRepository(c)

			// Create a new encryption instance
			encrypt := asimencrypt.NewAsimEncrypt()
			security := security.NewSecurity()

			// Create a new user handler
			userHandler := NewHandler(s, nil, encrypt, security)

			// Initialize the router
			r := gin.New()
			r.POST("/api/user/refresh/token", func(context *gin.Context) {}, userHandler.RefreshToken)

			// Create a new request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/refresh/token", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Token-Refresh", tt.tokenRefresh)
			req.Header.Set("Token", tt.token)

			// Serve the HTTP request
			r.ServeHTTP(w, req)

			// Assert the expected status code
			assert.Equal(t, tt.statusCode, w.Code)

		})
	}
}

func TestLogin(t *testing.T) {
	// Define a type for the mock behavior function
	type mockBehaviorGetUserByName func(ctx context.Context, mocks *mock_user.MockUserRepository, username string, user *user.User)

	// Define the test cases
	tests := []struct {
		user         *user.User
		name         string
		username     string
		body         string
		mockBehavior mockBehaviorGetUserByName
		statusCode   int
	}{
		{
			user:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			username: "asd",
			name:     "Login_200",
			body:     fmt.Sprintf(`{"username":"%s","password":"%s"}`, "asd", "asd"),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, username string, user *user.User) {
				// Mock the GetByUsername method to return the user
				mocks.EXPECT().GetByUsername(ctx, username).Return(user, nil).AnyTimes()
			},
			statusCode: 200,
		},
		{
			user:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asds", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			username: "asds",
			name:     "Login 401 user not found",
			body:     fmt.Sprintf(`{"username":"%s","password":"%s"}`, "asds", "asd"),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, username string, user *user.User) {
				// Mock the GetByUsername method to return an error indicating the user is not found
				mocks.EXPECT().GetByUsername(ctx, username).Return(nil, userRepository.ErrDataNotFound).AnyTimes()
			},
			statusCode: 401,
		},
		{
			user:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			username: "asd",
			name:     "Login 401 password not verified",
			body:     fmt.Sprintf(`{"username":"%s","password":"%s"}`, "asd", "asd123"),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, username string, user *user.User) {
				// Mock the GetByUsername method to return an error indicating the user is not found
				mocks.EXPECT().GetByUsername(ctx, username).Return(nil, userRepository.ErrDataNotFound).AnyTimes()
			},
			statusCode: 401,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the test controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Create a new mock user repository
			s := mock_user.NewMockUserRepository(c)

			// Create a new encryption instance
			encrypt := asimencrypt.NewAsimEncrypt()
			security := security.NewSecurity()

			// Create a new user handler
			userHandler := NewHandler(s, nil, encrypt, security)

			// Initialize the router
			r := gin.New()
			r.POST("/api/user/login", func(context *gin.Context) {
				// Set the mock behavior for the current test case
				tt.mockBehavior(context, s, tt.username, tt.user)
			}, userHandler.Login)

			// Create a new request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			// Serve the HTTP request
			r.ServeHTTP(w, req)

			// Assert the expected status code
			assert.Equal(t, tt.statusCode, w.Code)

			// Additional assertions for the successful login case
			if tt.name == "Login_200" {
				assert.NotEmpty(t, w.Header().Get("Token"))
				assert.NotEmpty(t, w.Header().Get("TokenRefresh"))

				// Check the token returned from the server
				user, err := jwttoken.Welcom(w.Header().Get("Token"))
				assert.Empty(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, user.Username, tt.username)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	// Define a type for the mock behavior function
	type mockBehaviorCreate func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User)
	type mockBehaviorSecurity func(ctx context.Context, mocks *mock_security.MockISecurity, username string)

	// Define the test cases
	tests := []struct {
		user                 *user.User
		name                 string
		username             string
		body                 string
		mockBehavior         mockBehaviorCreate
		mockBehaviorSecutiry mockBehaviorSecurity
		statusCode           int
	}{
		{
			user:     &user.User{Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			name:     "Register_201",
			username: "asd",
			body:     fmt.Sprintf(`{"username":"%s","password":"%s"}`, "asd", "asd"),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, userData *user.User) {
				// Mock the GetByUsername method to return the user
				mocks.EXPECT().Create(ctx, userData).Return(nil).AnyTimes()
			},
			mockBehaviorSecutiry: func(ctx context.Context, mocks *mock_security.MockISecurity, username string) {
				mocks.EXPECT().EncryptPassword(ctx, username).Return("$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", nil).AnyTimes()
			},
			statusCode: 201,
		},
		{
			user:     &user.User{Username: "asd", Password: ""},
			name:     "Register_400_passwoprd_empty",
			username: "asd",
			body:     fmt.Sprintf(`{"username":"%s","password":"%s"}`, "asd", ""),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, userData *user.User) {
				// Mock the GetByUsername method to return the user
				mocks.EXPECT().Create(ctx, userData).Return(nil).AnyTimes()
			},
			mockBehaviorSecutiry: func(ctx context.Context, mocks *mock_security.MockISecurity, username string) {
				mocks.EXPECT().EncryptPassword(ctx, username).Return("$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", nil).AnyTimes()
			},
			statusCode: 400,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the test controller
			c := gomock.NewController(t)
			defer c.Finish()

			// Create a new mock user repository
			s := mock_user.NewMockUserRepository(c)
			sec := mock_security.NewMockISecurity(c)

			// Create a new encryption instance
			encrypt := asimencrypt.NewAsimEncrypt()

			// Create a new user handler
			userHandler := NewHandler(s, nil, encrypt, sec)

			// Initialize the router
			r := gin.New()
			r.POST("/api/user/register", func(context *gin.Context) {
				// Set the mock behavior for the current test case
				tt.mockBehavior(context, s, tt.user)
				tt.mockBehaviorSecutiry(context, sec, tt.username)
			}, userHandler.Register)

			// Create a new request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			// Serve the HTTP request
			r.ServeHTTP(w, req)

			// Assert the expected status code
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestPasswordList(t *testing.T) {
	type mockBehaviorGetPasswordByUser func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User)
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"}

	requestBody := &model.RequestBody{Body: []byte("get_list_password"), User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	tests := []struct {
		name                 string
		userData             *user.User
		body                 string
		mockBehavior         mockBehaviorGetPasswordByUser
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:     "PasswordList_200",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     string(requestBodyMarshal),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User) {
				mocks.EXPECT().GetPasswordByUser(ctx, user).Return(&[]userRepository.ResponseSaveData{
					{ID: "1", Data: "password1"},
					{ID: "2", Data: "password2"},
				}, nil).AnyTimes()
			},
			statusCode:           200,
			expectedResponseBody: `[{"id":"1","data":"password1"},{"id":"2","data":"password2"}]`,
		},
		{
			name:     "PasswordList_400_invalid_json",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     `invalid json`,
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:     "PasswordList_500_internal_error",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     string(requestBodyMarshal),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User) {
				mocks.EXPECT().GetPasswordByUser(ctx, user).Return(nil, fmt.Errorf("internal error")).AnyTimes()
			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserRepo := mock_user.NewMockUserRepository(c)
			encrypt := asimencrypt.NewAsimEncrypt()
			userHandler := NewHandler(mockUserRepo, nil, encrypt, nil)

			r := gin.Default()

			r.POST("/api/user/list/passwords",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockUserRepo, tt.userData)
				},
				userHandler.PasswordList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/list/passwords", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestCardsList(t *testing.T) {
	type mockBehaviorGetCardsByUser func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User)
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"}

	requestBody := &model.RequestBody{Body: []byte("get_list_files"), User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	tests := []struct {
		name                 string
		userData             *user.User
		body                 string
		mockBehavior         mockBehaviorGetCardsByUser
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:     "CardsList_200",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     string(requestBodyMarshal),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User) {
				mocks.EXPECT().GetCardsByUser(ctx, user).Return(&[]userRepository.ResponseSaveData{
					{ID: "1", Data: "card1"},
					{ID: "2", Data: "card2"},
				}, nil).AnyTimes()
			},
			statusCode:           200,
			expectedResponseBody: `[{"id":"1","data":"card1"},{"id":"2","data":"card2"}]`,
		},
		{
			name:     "CardsList_400_invalid_json",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     `invalid json`,
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:     "CardsList_500_internal_error",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     string(requestBodyMarshal),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, user *user.User) {
				mocks.EXPECT().GetCardsByUser(ctx, user).Return(nil, fmt.Errorf("internal error")).AnyTimes()
			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserRepo := mock_user.NewMockUserRepository(c)
			encrypt := asimencrypt.NewAsimEncrypt()
			userHandler := NewHandler(mockUserRepo, nil, encrypt, nil)

			r := gin.Default()

			r.POST("/api/user/list/cards",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockUserRepo, tt.userData)
				},
				userHandler.CardsList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/list/cards", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestFilesList(t *testing.T) {
	type ResponseLists struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}

	encrypt := asimencrypt.NewAsimEncrypt()

	encryptPrivateClientKey, err := GetRSAPrivateKey([]byte(PrivateClientKey))
	if err != nil {
		t.Fatalf("Failed to encrypt rsa key: %v", err)
	}
	encrypt.PrivateClientKey = encryptPrivateClientKey

	type mockBehaviorGetFilesByUser func(ctx context.Context, mocks *mock_files.MockFilesRepository, user *user.User)
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey}

	requestBody := &model.RequestBody{Body: []byte("get_list_files"), User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	timeStr := "2024-09-24T19:07:34.871519+03:00"
	// Определяем формат строки даты
	layout := time.RFC3339

	// Преобразуем строку в time.Time
	times, err := time.Parse(layout, timeStr)
	if err != nil {
		t.Fatalf("Failed to time.Time parsing: %v", err)
	}

	type decryptResponseBody func(body []byte, expectedResponseBody string)

	tests := []struct {
		name                 string
		userData             *user.User
		body                 string
		mockBehavior         mockBehaviorGetFilesByUser
		checkResponseBody    decryptResponseBody
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:     "FilesList_200",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			body:     string(requestBodyMarshal),
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, user *user.User) {

				respList := []ResponseLists{}

				filesMetadata := files.Metadata{ClientID: "66f1792b0ce65bda608ef22b", CountPart: 100, UID: "10050066f1792b0ce65bda608ef22b"}

				files := &[]files.Files{
					{ID_File: "66f1792b0ce65bda608ef22b", Filename: "files1", ChunkSize: 5000, Metadata: filesMetadata, UploadDate: times},
				}
				for _, val := range *files {
					valJson, err := json.Marshal(val)
					if err != nil {
						t.Fatalf("Failed to marshal request body: %v", err)
					}
					data, err := encrypt.EncryptByClientKeyParts(string(valJson), PublicClientKey)
					if err != nil {
						t.Fatalf("asimencrypt failed to encrypt: %v", err)
					}

					_ = append(respList, ResponseLists{ID: val.ID.Hex(), Data: base64.StdEncoding.EncodeToString(data)})
				}
				mocks.EXPECT().GetByUserIdListFiles(ctx, user).Return(files, nil).AnyTimes()

			},
			checkResponseBody: func(body []byte, expectedResponseBody string) {
				respList := []ResponseLists{}

				err = json.Unmarshal(body /*[]byte(w.Body.String())*/, &respList)
				if err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}

				var valDecript []byte

				for _, val := range respList {

					valDecode, err := base64.StdEncoding.DecodeString(val.Data)
					if err != nil {
						t.Fatalf("Failed to DecodeString request body: %v", err)
					}

					valDecript, err = encrypt.DecryptOAEPClient(valDecode)
					if err != nil {
						t.Fatalf("Failed to DecryptOAEPClient request body: %v", err)
					}

				}

				assert.JSONEq(t, expectedResponseBody, string(valDecript))

			},
			statusCode:           200,
			expectedResponseBody: `{"_id":"000000000000000000000000","id_file":"66f1792b0ce65bda608ef22b","filename":"files1","chunk_size":5000,"metadata":{"ClientID":"66f1792b0ce65bda608ef22b","CountPart":100,"UID":"10050066f1792b0ce65bda608ef22b"},"upload_date":"2024-09-24T19:07:34.871519+03:00","length":0}`,
		},
		{
			name:     "FilesList_400_invalid_json",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			body:     `invalid json`,
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, user *user.User) {
				// No mock behavior needed for invalid JSON
			},
			checkResponseBody: func(body []byte, expectedResponseBody string) {
				assert.JSONEq(t, expectedResponseBody, string(body))
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:     "FilesList_500_internal_error",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			body:     string(requestBodyMarshal),
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, user *user.User) {
				mocks.EXPECT().GetByUserIdListFiles(ctx, user).Return(nil, fmt.Errorf("internal error")).AnyTimes()
			},
			checkResponseBody: func(body []byte, expectedResponseBody string) {
				assert.JSONEq(t, expectedResponseBody, string(body))
			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockFileRepo := mock_files.NewMockFilesRepository(c)
			userHandler := NewHandler(nil, mockFileRepo, encrypt, nil)

			r := gin.Default()

			r.POST("/api/user/list/files",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockFileRepo, tt.userData)
				},
				userHandler.FilesList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/list/files", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			tt.checkResponseBody(w.Body.Bytes(), tt.expectedResponseBody)

		})
	}
}

func GetRequestForPasswordDeleteByID(t *testing.T, id string) []byte {
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"}
	bodyData := PasswordRequest{ID: id, Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	requestBody := &model.RequestBody{Body: requestBodyData, User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyMarshal
}

func TestPasswordDelete(t *testing.T) {
	type mockBehaviorPasswordDeleteByID func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string)

	tests := []struct {
		objectIDStr          string
		name                 string
		userData             *user.User
		body                 string
		mockBehavior         mockBehaviorPasswordDeleteByID
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:        "DeletePassword_200",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForPasswordDeleteByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {

				objectID, err := primitive.ObjectIDFromHex(objectIDStr)
				if err != nil {
					t.Fatalf("Failed to ObjectID strint to object: %v", err)
				}
				mocks.EXPECT().DelerePasswordById(ctx, objectID).Return(nil).AnyTimes()
			},
			statusCode:           200,
			expectedResponseBody: `{"result":"delete is successful"}`,
		},
		{
			name:        "DeletePassword_400_invalid_json",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        `invalid json`,
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:        "DeletePassword_500_internal_error",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForPasswordDeleteByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {
				objectID, err := primitive.ObjectIDFromHex(objectIDStr)
				if err != nil {
					t.Fatalf("Failed to ObjectID strint to object: %v", err)
				}
				mocks.EXPECT().DelerePasswordById(ctx, objectID).Return(errors.New("internal error")).AnyTimes()
			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:        "ObjectIdISEmpty_400_internal_error",
			objectIDStr: "",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForPasswordDeleteByID(t, "")),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"id is empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserRepo := mock_user.NewMockUserRepository(c)
			encrypt := asimencrypt.NewAsimEncrypt()
			userHandler := NewHandler(mockUserRepo, nil, encrypt, nil)

			r := gin.Default()

			r.POST("/api/user/delete/password",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockUserRepo, tt.objectIDStr)
				},
				userHandler.PasswordDelete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/delete/password", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func GetRequestForCardDeleteByID(t *testing.T, id string) []byte {
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"}
	bodyData := CardRequest{ID: id, Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	requestBody := &model.RequestBody{Body: requestBodyData, User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyMarshal
}

func TestCardDelete(t *testing.T) {
	type mockBehaviorCardDeleteByID func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string)

	tests := []struct {
		objectIDStr          string
		name                 string
		userData             *user.User
		body                 string
		mockBehavior         mockBehaviorCardDeleteByID
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:        "DeleteCard_200",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForCardDeleteByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {

				objectID, err := primitive.ObjectIDFromHex(objectIDStr)
				if err != nil {
					t.Fatalf("Failed to ObjectID strint to object: %v", err)
				}
				mocks.EXPECT().DelereCardById(ctx, objectID).Return(nil).AnyTimes()
			},
			statusCode:           200,
			expectedResponseBody: `{"result":"delete is successful"}`,
		},
		{
			name:        "DeleteCard_400_invalid_json",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        `invalid json`,
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:        "DeleteCard_500_internal_error",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForCardDeleteByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {
				objectID, err := primitive.ObjectIDFromHex(objectIDStr)
				if err != nil {
					t.Fatalf("Failed to ObjectID strint to object: %v", err)
				}
				mocks.EXPECT().DelereCardById(ctx, objectID).Return(errors.New("internal error")).AnyTimes()
			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:        "ObjectIdISEmpty_400_internal_error",
			objectIDStr: "",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForCardDeleteByID(t, "")),
			mockBehavior: func(ctx context.Context, mocks *mock_user.MockUserRepository, objectIDStr string) {
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"id is empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserRepo := mock_user.NewMockUserRepository(c)
			encrypt := asimencrypt.NewAsimEncrypt()
			userHandler := NewHandler(mockUserRepo, nil, encrypt, nil)

			r := gin.Default()

			r.POST("/api/user/delete/card",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockUserRepo, tt.objectIDStr)
				},
				userHandler.CardDelete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/delete/card", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func GetRequestForFileDeleteByID(t *testing.T, id string) []byte {
	timeStr := "2024-09-24T19:07:34.871519+03:00"
	// Определяем формат строки даты
	layout := time.RFC3339

	// Преобразуем строку в time.Time
	times, err := time.Parse(layout, timeStr)
	if err != nil {
		t.Fatalf("Failed to time.Time parsing: %v", err)
	}

	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"}
	bodyData := files.Files{ID_File: id, Filename: "File move TarakanTV", Length: 3215, ChunkSize: 567678, UploadDate: times}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	requestBody := &model.RequestBody{Body: requestBodyData, User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyMarshal
}

func TestFileDelete(t *testing.T) {
	type mockBehaviorFileDeleteByID func(ctx context.Context, mocks *mock_files.MockFilesRepository, objectIDStr string)
	//`{"_id": "2","filename": "fotodog.jpg", "length": 16808746,"chunk_size": 50000,"upload_date": "2024-09-21T17:25:36.433Z"}`
	tests := []struct {
		objectIDStr          string
		name                 string
		userData             *user.User
		body                 string
		mockBehavior         mockBehaviorFileDeleteByID
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:        "DeleteFile_200",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForFileDeleteByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, objectIDStr string) {
				mocks.EXPECT().DeleteFilesByID(ctx, objectIDStr).Return(nil).AnyTimes()
			},
			statusCode:           200,
			expectedResponseBody: `{"result":"delete is successful"}`,
		},
		{
			name:        "DeleteFile_400_invalid_json",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        `invalid json`,
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, objectIDStr string) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:        "DeleteCard_500_internal_error",
			objectIDStr: "66f1792b0ce65bda608ef22b",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForFileDeleteByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, objectIDStr string) {
				mocks.EXPECT().DeleteFilesByID(ctx, objectIDStr).Return(errors.New("internal error")).AnyTimes()
			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:        "ObjectIdISEmpty_400_internal_error",
			objectIDStr: "",
			userData:    &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:        string(GetRequestForFileDeleteByID(t, "")),
			mockBehavior: func(ctx context.Context, mocks *mock_files.MockFilesRepository, objectIDStr string) {
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"id is empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockFileRepo := mock_files.NewMockFilesRepository(c)
			encrypt := asimencrypt.NewAsimEncrypt()
			userHandler := NewHandler(nil, mockFileRepo, encrypt, nil)

			r := gin.Default()

			r.POST("/api/user/delete/file",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockFileRepo, tt.objectIDStr)
				},
				userHandler.DeleteFile)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/delete/file", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func GetRequestForSaveCardByID(t *testing.T, id string) []byte {
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey}
	bodyData := CardRequest{ID: id, Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	requestBody := &model.RequestBody{Body: requestBodyData, User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyMarshal
}

func GetForSaveData(t *testing.T, id string) []byte {
	bodyData := CardRequest{ID: id, Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyData
}

func TestCardSave(t *testing.T) {
	encrypt := asimencrypt.NewAsimEncrypt()

	type mockBehaviorCardSaveByUser func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, user *user.User, card *CardRequest)

	tests := []struct {
		name                 string
		userData             *user.User
		cardData             *CardRequest
		body                 string
		mockBehavior         mockBehaviorCardSaveByUser
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:     "CardsSaveUpdate_200",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			cardData: &CardRequest{ID: "66f1792b0ce65bda608ef22b", Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"},
			//body:     string(requestBodyMarshal),
			body: string(GetRequestForSaveCardByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, card *CardRequest) {
				saveCardByData := &user.SaveData{}

				saveCardByData.TypeData = "card"

				if card.ID != "" {
					docID, err := primitive.ObjectIDFromHex(card.ID)
					if err != nil {
						log.Fatal(err)
					}
					saveCardByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSaveData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				saveCardByData.Data = base64.StdEncoding.EncodeToString(data)
				saveCardByData.User_ID = userData.ID_User

				//mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()
				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSaveData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()

			},
			statusCode:           200,
			expectedResponseBody: `{"result": "update is successful"}`,
		},
		{
			name:     "CardsSaveCreate_201",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			cardData: &CardRequest{Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"},
			//body:     string(requestBodyMarshal),
			body: string(GetRequestForSaveCardByID(t, "")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, card *CardRequest) {
				saveCardByData := &user.SaveData{}

				saveCardByData.TypeData = "card"

				if card.ID != "" {
					docID, err := primitive.ObjectIDFromHex(card.ID)
					if err != nil {
						log.Fatal(err)
					}
					saveCardByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSaveData(t, "")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				saveCardByData.Data = base64.StdEncoding.EncodeToString(data)
				saveCardByData.User_ID = userData.ID_User

				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSaveData(t, "")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().CreateCardByUser(ctx, saveCardByData).Return(nil).AnyTimes()

			},
			statusCode:           201,
			expectedResponseBody: `{"result": "create is successful"}`,
		},
		{
			name:     "CardsList_400_invalid_json",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     `invalid json`,
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, card *CardRequest) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:     "CardsSaveCreate_500",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			cardData: &CardRequest{Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"},
			//body:     string(requestBodyMarshal),
			body: string(GetRequestForSaveCardByID(t, "")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, card *CardRequest) {
				saveCardByData := &user.SaveData{}

				saveCardByData.TypeData = "card"

				if card.ID != "" {
					docID, err := primitive.ObjectIDFromHex(card.ID)
					if err != nil {
						log.Fatal(err)
					}
					saveCardByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSaveData(t, "")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				saveCardByData.Data = base64.StdEncoding.EncodeToString(data)
				saveCardByData.User_ID = userData.ID_User

				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSaveData(t, "")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().CreateCardByUser(ctx, saveCardByData).Return(errors.New("internal error")).AnyTimes()

			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:     "CardsSaveUpdate_500",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			cardData: &CardRequest{ID: "66f1792b0ce65bda608ef22b", Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"},
			body:     string(GetRequestForSaveCardByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, card *CardRequest) {
				saveCardByData := &user.SaveData{}

				saveCardByData.TypeData = "card"

				if card.ID != "" {
					docID, err := primitive.ObjectIDFromHex(card.ID)
					if err != nil {
						log.Fatal(err)
					}
					saveCardByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSaveData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				saveCardByData.Data = base64.StdEncoding.EncodeToString(data)
				saveCardByData.User_ID = userData.ID_User

				//mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()
				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSaveData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(errors.New("internal error")).AnyTimes()

			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new controller for the test
			c := gomock.NewController(t)
			defer c.Finish()

			// Create mock repositories and services
			mockUserRepo := mock_user.NewMockUserRepository(c)
			mockEncrypt := mock_asimencrypt.NewMockAsimEncrypt(c)
			userHandler := NewHandler(mockUserRepo, nil, mockEncrypt, nil)

			// Create a new Gin router
			r := gin.Default()

			// Define the route for the test
			r.POST("/api/user/save/card",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockEncrypt, mockUserRepo, tt.userData, tt.cardData)
				},
				userHandler.CardSave)

			// Create a new HTTP response recorder
			w := httptest.NewRecorder()
			// Create a new HTTP request
			req := httptest.NewRequest("POST", "/api/user/save/card", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			// Serve the HTTP request
			r.ServeHTTP(w, req)

			// Assert the expected status code
			assert.Equal(t, tt.statusCode, w.Code)
			// Assert the expected response body
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}

}

func GetRequestForSavePasswordByID(t *testing.T, id string) []byte {
	userData := user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey}
	//bodyData := CardRequest{ID: id, Description: "Card 1 Zoom Bank", Number: "3215", Exp: "12/24", Cvc: "456"}
	bodyData := PasswordRequest{ID: id, Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	requestBody := &model.RequestBody{Body: requestBodyData, User: userData}

	requestBodyMarshal, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyMarshal
}

func GetForSavePasswordData(t *testing.T, id string) []byte {
	bodyData := PasswordRequest{ID: id, Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"}

	requestBodyData, err := json.Marshal(bodyData)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	return requestBodyData
}

func TestPasswordSave(t *testing.T) {
	encrypt := asimencrypt.NewAsimEncrypt()

	type mockBehaviorPasswordSaveByUser func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, user *user.User, password *PasswordRequest)

	tests := []struct {
		name                 string
		userData             *user.User
		passwordData         *PasswordRequest
		body                 string
		mockBehavior         mockBehaviorPasswordSaveByUser
		statusCode           int
		expectedResponseBody string
	}{
		{
			name:         "PasswordSaveUpdate_200",
			userData:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			passwordData: &PasswordRequest{ID: "66f1792b0ce65bda608ef22b", Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"},
			body:         string(GetRequestForSavePasswordByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, password *PasswordRequest) {
				savePasswordByData := &user.SaveData{}

				savePasswordByData.TypeData = "password"

				if password.ID != "" {
					docID, err := primitive.ObjectIDFromHex(password.ID)
					if err != nil {
						log.Fatal(err)
					}
					savePasswordByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSavePasswordData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				savePasswordByData.Data = base64.StdEncoding.EncodeToString(data)
				savePasswordByData.User_ID = userData.ID_User

				//mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()
				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSavePasswordData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().UpdatePasswordByKey(ctx, savePasswordByData).Return(nil).AnyTimes()

			},
			statusCode:           200,
			expectedResponseBody: `{"result": "update is successfull"}`,
		},
		{
			name:         "PasswordSaveCretae_201",
			userData:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			passwordData: &PasswordRequest{ID: "", Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"},
			body:         string(GetRequestForSavePasswordByID(t, "")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, password *PasswordRequest) {
				savePasswordByData := &user.SaveData{}

				savePasswordByData.TypeData = "password"

				if password.ID != "" {
					docID, err := primitive.ObjectIDFromHex(password.ID)
					if err != nil {
						log.Fatal(err)
					}
					savePasswordByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSavePasswordData(t, "")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				savePasswordByData.Data = base64.StdEncoding.EncodeToString(data)
				savePasswordByData.User_ID = userData.ID_User

				//mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()
				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSavePasswordData(t, "")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().CreatePasswordByUser(ctx, savePasswordByData).Return(nil).AnyTimes()

			},
			statusCode:           201,
			expectedResponseBody: `{"result": "create is successfull"}`,
		},

		{
			name:     "PasswordSave_400_invalid_json",
			userData: &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK"},
			body:     `invalid json`,
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, password *PasswordRequest) {
				// No mock behavior needed for invalid JSON
			},
			statusCode:           400,
			expectedResponseBody: `{"error":"invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:         "PasswordSaveUpdate_500",
			userData:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			passwordData: &PasswordRequest{ID: "66f1792b0ce65bda608ef22b", Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"},
			body:         string(GetRequestForSavePasswordByID(t, "66f1792b0ce65bda608ef22b")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, password *PasswordRequest) {
				savePasswordByData := &user.SaveData{}
				savePasswordByData.TypeData = "password"

				if password.ID != "" {
					docID, err := primitive.ObjectIDFromHex(password.ID)
					if err != nil {
						log.Fatal(err)
					}
					savePasswordByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSavePasswordData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				savePasswordByData.Data = base64.StdEncoding.EncodeToString(data)
				savePasswordByData.User_ID = userData.ID_User

				//mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()
				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSavePasswordData(t, "66f1792b0ce65bda608ef22b")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().UpdatePasswordByKey(ctx, savePasswordByData).Return(errors.New("internal error")).AnyTimes()

			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
		{
			name:         "PasswordSaveCretae_201",
			userData:     &user.User{ID_User: "66f1792b0ce65bda608ef22b", Username: "asd", Password: "$2a$10$S8xVbBSSXVSTHYVUPsmG9uqnnbs1iaS1D9a0Cso67KEo9Wun.lcJK", PublicKey: PublicClientKey},
			passwordData: &PasswordRequest{ID: "", Description: "new site movie login http://movie.dovie", Username: "iLikeMovie", Password: "$2a$10$S8xVbBSSXVST"},
			body:         string(GetRequestForSavePasswordByID(t, "")),
			mockBehavior: func(ctx context.Context, mocksEncrypt *mock_asimencrypt.MockAsimEncrypt, mocks *mock_user.MockUserRepository, userData *user.User, password *PasswordRequest) {
				savePasswordByData := &user.SaveData{}

				savePasswordByData.TypeData = "password"

				if password.ID != "" {
					docID, err := primitive.ObjectIDFromHex(password.ID)
					if err != nil {
						log.Fatal(err)
					}
					savePasswordByData.ID = docID
				}

				//Берем ключ клиента публичны и шифруем данные
				data, err := encrypt.EncryptByClientKeyParts(string(GetForSavePasswordData(t, "")), PublicClientKey)
				if err != nil {
					log.Println("asimencrypt failed to encrypt", err)
				}

				//Зашифрованные данные преобразуем в base64 чтобы сохранить в базе
				savePasswordByData.Data = base64.StdEncoding.EncodeToString(data)
				savePasswordByData.User_ID = userData.ID_User

				//mocks.EXPECT().UpdateCardByKey(ctx, saveCardByData).Return(nil).AnyTimes()
				mocksEncrypt.EXPECT().EncryptByClientKeyParts(string(GetForSavePasswordData(t, "")), PublicClientKey).Return(data, nil).AnyTimes()
				mocks.EXPECT().CreatePasswordByUser(ctx, savePasswordByData).Return(errors.New("internal error")).AnyTimes()

			},
			statusCode:           500,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new controller for the test
			c := gomock.NewController(t)
			defer c.Finish()

			// Create mock repositories and services
			mockUserRepo := mock_user.NewMockUserRepository(c)
			mockEncrypt := mock_asimencrypt.NewMockAsimEncrypt(c)
			userHandler := NewHandler(mockUserRepo, nil, mockEncrypt, nil)

			// Create a new Gin router
			r := gin.Default()

			// Define the route for the test
			r.POST("/api/user/save/password",
				func(context *gin.Context) {
					// Set the mock behavior for the current test case
					tt.mockBehavior(context, mockEncrypt, mockUserRepo, tt.userData, tt.passwordData)
				},
				userHandler.PasswordSave)

			// Create a new HTTP response recorder
			w := httptest.NewRecorder()
			// Create a new HTTP request
			req := httptest.NewRequest("POST", "/api/user/save/password", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			// Serve the HTTP request
			r.ServeHTTP(w, req)

			// Assert the expected status code
			assert.Equal(t, tt.statusCode, w.Code)
			// Assert the expected response body
			assert.JSONEq(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
