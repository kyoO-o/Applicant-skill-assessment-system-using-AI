package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/ocookie"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
	"golang.org/x/oauth2"
)

// #region Mutex locks

var ocrMutex sync.Mutex

// var (
// 	activeContexts = make(map[int]context.CancelFunc) // Map to hold cancel functions for each user
// 	mu             sync.Mutex                         // Mutex to protect access to activeContexts
// )

// #region Authorization

func Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(app.ContextKeyAuthCustomer).(*userman.User)
	oapi.SendResp(w, user)
}

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type authResponse struct {
	User    *userman.User `json:"user"`
	Message string        `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Name, email, password, and role are required"})
		return
	}

	if len(req.Password) < 8 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Password must be at least 8 characters"})
		return
	}

	if req.Role != userman.RoleUser && req.Role != userman.RoleRecruiter {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Role must be either user or recruiter"})
		return
	}

	_, err := app.Users.GetWithEmail(req.Email)
	if err == nil {
		oapi.CustomError(w, http.StatusConflict, map[string]string{"message": "An account with that email already exists"})
		return
	}
	if !errors.Is(err, userman.ErrNotFound) {
		oapi.ServerError(w, err)
		return
	}

	passwordHash, err := common.HashPassword(req.Password)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	user := &userman.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: passwordHash,
		Role:         req.Role,
	}

	user, err = app.Users.Save(user)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	startSession(r, user, "local")
	w.WriteHeader(http.StatusCreated)
	oapi.SendResp(w, &authResponse{
		User:    user,
		Message: "Registration successful",
	})
}

func PasswordLogin(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	if email == "" || password == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Email and password are required"})
		return
	}

	user, err := app.Users.GetWithEmail(email)
	if err != nil {
		if errors.Is(err, userman.ErrNotFound) {
			oapi.CustomError(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
			return
		}
		oapi.ServerError(w, err)
		return
	}

	if user.PasswordHash == "" || common.CheckPassword(password, user.PasswordHash) != nil {
		oapi.CustomError(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
		return
	}

	startSession(r, user, "local")
	oapi.SendResp(w, &authResponse{
		User:    user,
		Message: "Login successful",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()
	url := app.Todu.AuthCodeURL(state, oauth2.AccessTypeOffline)

	app.Session.Put(r, "oauth_state", state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func BAuthenticate(w http.ResponseWriter, r *http.Request) {
	expectedState := app.Session.GetString(r, "oauth_state")
	state := r.URL.Query().Get("state")
	if expectedState == "" || state != expectedState {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	app.Session.Remove(r, "oauth_state")

	code := r.URL.Query().Get("code")
	token, err := app.Todu.Exchange(r.Context(), code)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	client := app.Todu.Client(r.Context(), token)
	resp, err := client.Get(app.Config.Todu.Endpoint.UserInfo)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	defer resp.Body.Close()

	var userInfo userman.ToduUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		oapi.ServerError(w, err)
		return
	}

	app.Session.Put(r, "userID", userInfo.Sub)
	app.InfoLog.Println("userID", userInfo.Sub, userInfo)
	toduID, err := strconv.Atoi(userInfo.Sub)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}
	var user *userman.User
	user, err = app.Users.GetWithToduID(toduID)
	if err != nil {
		if err != userman.ErrNotFound {
			oapi.ServerError(w, err)
			return
		}
		user = &userman.User{
			Email:       userInfo.Email,
			Name:        userInfo.Name,
			ToduID:      toduID,
			PhoneNumber: userInfo.PhoneNumber,
			Role:        userman.RoleUser,
		}
		if _, err := app.Users.Save(user); err != nil {
			oapi.ServerError(w, err)
			return
		}

	}

	app.Session.Put(r, "email", user.Email)
	app.Session.Put(r, "name", user.Name)
	app.Session.Put(r, "todu_id", user.ToduID)
	app.Session.Put(r, "accessToken", token.AccessToken)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// func TokenLogout(w http.ResponseWriter, r *http.Request) {
// 	var requestBody struct {
// 		LogoutToken string `json:"logout_token"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	sessionToken := app.Session.GetString(r, "accessToken")
// 	if sessionToken == requestBody.LogoutToken {
// 		app.Session.Remove(r, "userID")
// 		app.Session.Remove(r, "email")
// 		app.Session.Remove(r, "accessToken")
// 		app.Session.Remove(r, "oauth_state")
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

func Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Remove(r, "email")
	app.Session.Remove(r, "name")
	app.Session.Remove(r, "todu_id")
	app.Session.Remove(r, "userID")
	app.Session.Remove(r, "accessToken")
	ocookie.Remove(w, r, "session")
	oapi.SendResp(w, map[string]string{"message": "Logout successful"})
}

func startSession(r *http.Request, user *userman.User, accessToken string) {
	app.Session.Put(r, "userID", user.ID)
	app.Session.Put(r, "email", user.Email)
	app.Session.Put(r, "name", user.Name)
	app.Session.Put(r, "todu_id", user.ToduID)
	app.Session.Put(r, "accessToken", accessToken)
}
