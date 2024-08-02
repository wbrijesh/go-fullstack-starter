package internal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-fullstack-starter/schema"
	"go-fullstack-starter/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, email string, password string) (*schema.User, error) {
	var count int
	checkUserQuery := `
    SELECT COUNT(*) FROM users WHERE email = ?;
  `
	err := db.QueryRow(checkUserQuery, email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	if email == "" || password == "" || !utils.IsValidEmail(email) {
		return nil, errors.New("invalid email or password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := schema.User{
		ID:             uuid.New().String(),
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	query := `
    INSERT INTO users (id, email, hashed_password)
    VALUES (?, ?, ?);
  `

	_, err = db.Exec(query, user.ID, user.Email, user.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func VerifyUser(db *sql.DB, email string, password string) bool {
	var hashedPassword string
	query := `
    SELECT hashed_password FROM users WHERE email = ?;
  `
	err := db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

var JWTSecret = []byte("secret")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-fullstack-starter",
			Subject:   email,
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("email", claims.Email)
		next.ServeHTTP(w, r)
	}
}

func UserRegistration(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	utils.HandleMethodNotAllowed(w, r, http.MethodPost)

	var user schema.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleError(w, r, http.StatusBadRequest, "Invalid request body: ", err)
		return
	}

	_, err = CreateUser(db, user.Email, user.Password)
	if err != nil {
		utils.HandleError(w, r, http.StatusInternalServerError, "Failed to create user: ", err)
		return
	}

	token, err := GenerateToken(user.Email)
	if err != nil {
		utils.HandleError(w, r, http.StatusInternalServerError, "Failed to generate token: ", err)
		return
	}

	response := schema.Response{
		Status:    "SUCCESS",
		Message:   "User created successfully",
		RequestID: r.Context().Value(schema.RequestIDKey{}).(string),
		Data:      map[string]string{"token": token},
	}

	utils.SendResponse(w, r, response)
}

func UserLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	utils.HandleMethodNotAllowed(w, r, http.MethodPost)
	var user schema.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.HandleError(w, r, http.StatusBadRequest, "Invalid request body: ", err)
		return
	}
	if !VerifyUser(db, user.Email, user.Password) {
		utils.HandleError(w, r, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}
	token, err := GenerateToken(user.Email)
	if err != nil {
		utils.HandleError(w, r, http.StatusInternalServerError, "Failed to generate token: ", err)
		return
	}
	response := schema.Response{
		Status:    "SUCCESS",
		Message:   "User logged in successfully",
		RequestID: r.Context().Value(schema.RequestIDKey{}).(string),
		Data:      map[string]string{"token": token},
	}
	utils.SendResponse(w, r, response)
}
