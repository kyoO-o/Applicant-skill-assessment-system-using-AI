package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
	"golang.org/x/oauth2"
)

func gcalConnect(w http.ResponseWriter, r *http.Request) {
	if app.GoogleCalendar == nil {
		oapi.CustomError(w, http.StatusServiceUnavailable, map[string]string{"message": "Google Calendar интеграц тохиргоогүй байна"})
		return
	}
	state := generateRandomState()
	app.Session.Put(r, "gcal_oauth_state", state)
	url := app.GoogleCalendar.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func gcalCallback(w http.ResponseWriter, r *http.Request) {
	expected := app.Session.GetString(r, "gcal_oauth_state")
	state := r.URL.Query().Get("state")
	if expected == "" || state != expected {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	app.Session.Remove(r, "gcal_oauth_state")

	user := authUser(r)
	code := r.URL.Query().Get("code")
	token, err := app.GoogleCalendar.Exchange(r.Context(), code)
	if err != nil {
		app.ErrorLog.Printf("gcal token exchange failed: %v", err)
		frontendURL := app.Config.FrontendURL
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		http.Redirect(w, r, frontendURL+"/settings?gcal=error", http.StatusTemporaryRedirect)
		return
	}

	if token.RefreshToken != "" {
		user.GoogleRefreshToken = token.RefreshToken
		if _, err := app.Users.Save(user); err != nil {
			oapi.ServerError(w, err)
			return
		}
	}

	frontendURL := app.Config.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	http.Redirect(w, r, frontendURL+"/settings?gcal=connected", http.StatusTemporaryRedirect)
}

func gcalDisconnect(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	user.GoogleRefreshToken = ""
	if _, err := app.Users.Save(user); err != nil {
		oapi.ServerError(w, err)
		return
	}
	oapi.SendResp(w, map[string]string{"message": "Google Calendar холболт цуцлагдлаа"})
}

func gcalStatus(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	oapi.SendResp(w, map[string]bool{"connected": user.GoogleRefreshToken != ""})
}

func createGCalEvent(ctx context.Context, recruiter *userman.User, title, description, location string, start time.Time) error {
	if app.GoogleCalendar == nil || recruiter.GoogleRefreshToken == "" {
		return nil
	}

	ts := app.GoogleCalendar.TokenSource(ctx, &oauth2.Token{RefreshToken: recruiter.GoogleRefreshToken})
	token, err := ts.Token()
	if err != nil {
		return fmt.Errorf("gcal token refresh: %w", err)
	}

	end := start.Add(time.Hour)
	tz := "Asia/Ulaanbaatar"
	event := map[string]any{
		"summary":     title,
		"description": description,
		"location":    location,
		"start":       map[string]string{"dateTime": start.Format(time.RFC3339), "timeZone": tz},
		"end":         map[string]string{"dateTime": end.Format(time.RFC3339), "timeZone": tz},
	}

	body, _ := json.Marshal(event)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://www.googleapis.com/calendar/v3/calendars/primary/events",
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("gcal API returned %d", resp.StatusCode)
	}
	return nil
}
