package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nikvas0/dc-homework/auth/objects"
	"github.com/nikvas0/dc-homework/auth/queues"
	"github.com/nikvas0/dc-homework/auth/storage"
	"github.com/nikvas0/dc-homework/auth/utils"
	"golang.org/x/crypto/bcrypt"
)

const salt = "salt"
const refreshExpirationTime = 7 * 24 * time.Hour

func SignUp(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("SignUp request error: Error while reading request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userData := objects.UserData{}
	err = json.Unmarshal(reqBody, &userData)
	if err != nil {
		log.Println("SignUp request error: Got broken JSON.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := objects.User{}
	user.Email = userData.Email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userData.Password+salt), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)
	user.Confirmed = false
	err = storage.CreateUser(&user)
	if storage.IsErrorAlreadyExists(err) {
		log.Println("SignUp request error: already exists.")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		log.Println("SignUp request error: Database error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	confirmToken, err := utils.GenerateConfirmToken()
	if err != nil {
		log.Println("SignUp request error: failed to generate token.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = storage.CreateConfirmationToken(&user, confirmToken)
	if err != nil {
		log.Println("SignUp error: Failed to create confirmation token.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = queues.SheduleConfirmation(&user, confirmToken)
	if err != nil {
		log.Println("SignUp error: Failed to put request to queue.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusCreated)
	log.Printf("SignUp request: success (id=%d).", user.ID)
}

func Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.String()[12:]
	var userID uint32
	err := storage.GetUserIDByConfirmationToken(token, &userID)
	if err != nil {
		log.Println("Confirm error: Failed to find userID by token.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = storage.SetUserConfirmedByID(userID)
	if err != nil {
		log.Println("Confirm error: Failed to find user by id.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Confirmed user %d", userID)
	w.WriteHeader(http.StatusOK)
	return
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("SignIn request error: Error while reading request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userData := objects.UserData{}
	err = json.Unmarshal(reqBody, &userData)
	if err != nil {
		log.Println("SignIn request error: Got broken JSON.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := objects.User{}
	err = storage.GetUserByEmail(userData.Email, &user)
	if storage.IsNotFoundError(err) {
		log.Println("SignIn request error: Not found.")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		log.Println("SignIn request error: Database error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !user.Confirmed || user.Email != userData.Email || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userData.Password+salt)) != nil {
		log.Println("SignIn request error: Wrong email or password")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(&user)
	if err != nil {
		log.Println("Error creating tokens")
		return
	}

	session := objects.Session{}
	session.UserID = user.ID
	session.RefreshToken = refreshToken
	session.ExpireAt = time.Now().Add(refreshExpirationTime)
	err = storage.CreateSession(&session)
	if err != nil {
		log.Println("Error creating session")
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access":  accessToken,
		"refresh": refreshToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Get SignIn request error: Encoded broken JSON.")
		return
	}

	log.Printf("SignIn request: success (id=%d).", session.UserID)
}

type TokenRequest struct {
	Token string
}

func Validate(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Validate request error: Error while reading request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validateRequest := TokenRequest{}
	err = json.Unmarshal(reqBody, &validateRequest)
	if err != nil {
		log.Println("Validate request error: Got broken JSON.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := objects.Token{}
	tokenInfo, err := jwt.ParseWithClaims(validateRequest.Token, &token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_KEY")), nil
	})
	if err != nil || !tokenInfo.Valid {
		log.Println("Validate request error: Got bad token.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    token.UserID,
		"email": token.Email,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Get Validate request error: Encoded broken JSON.")
		return
	}

	log.Printf("Validate request: success (id=%d).", token.UserID)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Refresh request error: Error while reading request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshRequest := TokenRequest{}
	err = json.Unmarshal(reqBody, &refreshRequest)
	if err != nil {
		log.Println("Refresh request error: Got broken JSON.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session := objects.Session{}
	err = storage.GetSessionByToken(refreshRequest.Token, &session)
	if storage.IsNotFoundError(err) {
		log.Println("Refresh request error: Bad token.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("Refresh request error: Failed to update session.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := objects.User{}
	err = storage.GetUserByID(session.UserID, &user)
	if storage.IsNotFoundError(err) {
		log.Println("Refresh request error: User not found.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("Refresh request error: Database error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(&user)
	if err != nil {
		log.Println("Error creating tokens")
		return
	}

	session.RefreshToken = refreshToken
	session.ExpireAt = time.Now().Add(refreshExpirationTime)
	err = storage.UpdateSession(&session)
	if err != nil {
		log.Println("Error creating session")
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access":  accessToken,
		"refresh": refreshToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Get Refresh request error: Encoded broken JSON.")
		return
	}

	log.Printf("Refresh request: success (id=%d).", session.UserID)
	return
}
