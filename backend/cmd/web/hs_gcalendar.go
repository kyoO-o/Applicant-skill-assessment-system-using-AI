package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
	"golang.org/x/oauth2"
)

var errGCalEventNotFound = errors.New("gcal event not found")

func gcalConnect(w http.ResponseWriter, r *http.Request) {
	if app.GoogleCalendar == nil {
		oapi.CustomError(w, http.StatusServiceUnavailable, map[string]string{"message": "Google Calendar интеграц тохиргоогүй байна"})
		return
	}
	user := authUser(r)
	state := generateRandomState()
	app.Session.Put(r, "gcal_oauth_state", state)
	url := app.GoogleCalendar.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("login_hint", user.Email),
	)
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

type gcalEventResponse struct {
	ID             string `json:"id"`
	ConferenceData struct {
		ConferenceID  string `json:"conferenceId"`
		CreateRequest struct {
			Status struct {
				StatusCode string `json:"statusCode"`
			} `json:"status"`
		} `json:"createRequest"`
		EntryPoints []struct {
			EntryPointType string `json:"entryPointType"`
			URI            string `json:"uri"`
		} `json:"entryPoints"`
	} `json:"conferenceData"`
}

// createGCalEvent creates a calendar event in the recruiter's Google Calendar,
// optionally attaching a Google Meet conference. attendeeEmail (applicant) is
// added as a guest so they receive a calendar invitation.
// Returns the Meet link (empty when generateMeet is false) and the event ID.
func createGCalEvent(ctx context.Context, recruiter *userman.User, title, description, location, attendeeEmail string, start time.Time, generateMeet bool) (string, string, error) {
	if app.GoogleCalendar == nil || recruiter.GoogleRefreshToken == "" {
		return "", "", nil
	}

	ts := app.GoogleCalendar.TokenSource(ctx, &oauth2.Token{RefreshToken: recruiter.GoogleRefreshToken})
	token, err := ts.Token()
	if err != nil {
		return "", "", fmt.Errorf("gcal token refresh: %w", err)
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

	if attendeeEmail != "" {
		event["attendees"] = []map[string]string{{"email": attendeeEmail}}
	}

	if generateMeet {
		event["conferenceData"] = map[string]any{
			"createRequest": map[string]any{
				"requestId": uuid.New().String(),
				"conferenceSolutionKey": map[string]string{
					"type": "hangoutsMeet",
				},
			},
		}
	}

	body, _ := json.Marshal(event)

	// conferenceDataVersion=1 is required for Meet link generation;
	// sendUpdates=all delivers the calendar invite to attendees.
	apiURL := "https://www.googleapis.com/calendar/v3/calendars/primary/events?conferenceDataVersion=1&sendUpdates=all"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Always read the full body so we can log errors and decode the response.
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf("gcal API %d: %s", resp.StatusCode, string(respBody))
	}

	app.InfoLog.Printf("gcal event created (status=%d)", resp.StatusCode)

	var evResp gcalEventResponse
	if err := json.Unmarshal(respBody, &evResp); err != nil {
		return "", "", fmt.Errorf("gcal response decode: %w", err)
	}

	if !generateMeet {
		return "", evResp.ID, nil
	}

	app.InfoLog.Printf("gcal conferenceData: id=%q status=%q entryPoints=%d",
		evResp.ConferenceData.ConferenceID,
		evResp.ConferenceData.CreateRequest.Status.StatusCode,
		len(evResp.ConferenceData.EntryPoints))

	// If conference is still being provisioned, poll once for the final data.
	if evResp.ConferenceData.CreateRequest.Status.StatusCode == "pending" && evResp.ID != "" {
		time.Sleep(2 * time.Second)
		getURL := "https://www.googleapis.com/calendar/v3/calendars/primary/events/" + evResp.ID
		if getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil); err == nil {
			getReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
			if getResp, err := http.DefaultClient.Do(getReq); err == nil {
				defer getResp.Body.Close()
				if body, _ := io.ReadAll(getResp.Body); len(body) > 0 {
					json.Unmarshal(body, &evResp)
				}
			}
		}
	}

	// Primary: extract from entryPoints (conference created synchronously).
	for _, ep := range evResp.ConferenceData.EntryPoints {
		if ep.EntryPointType == "video" {
			return ep.URI, evResp.ID, nil
		}
	}

	// Fallback: Google sometimes returns status="pending" with no entryPoints yet
	// but always returns the conferenceId; the Meet URL is predictable from it.
	if evResp.ConferenceData.ConferenceID != "" {
		return "https://meet.google.com/" + evResp.ConferenceData.ConferenceID, evResp.ID, nil
	}

	return "", evResp.ID, fmt.Errorf("gcal Meet conference not returned (status=%q)", evResp.ConferenceData.CreateRequest.Status.StatusCode)
}

// patchGCalEventTime updates the start/end time of an existing calendar event.
func patchGCalEventTime(ctx context.Context, recruiter *userman.User, eventID string, newStart time.Time) error {
	if app.GoogleCalendar == nil || recruiter.GoogleRefreshToken == "" || eventID == "" {
		return nil
	}

	ts := app.GoogleCalendar.TokenSource(ctx, &oauth2.Token{RefreshToken: recruiter.GoogleRefreshToken})
	token, err := ts.Token()
	if err != nil {
		return fmt.Errorf("gcal token refresh: %w", err)
	}

	tz := "Asia/Ulaanbaatar"
	patch := map[string]any{
		"start": map[string]string{"dateTime": newStart.Format(time.RFC3339), "timeZone": tz},
		"end":   map[string]string{"dateTime": newStart.Add(time.Hour).Format(time.RFC3339), "timeZone": tz},
	}
	body, _ := json.Marshal(patch)

	patchURL := "https://www.googleapis.com/calendar/v3/calendars/primary/events/" + eventID + "?sendUpdates=all"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, patchURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return errGCalEventNotFound
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("gcal patch API %d: %s", resp.StatusCode, string(respBody))
	}

	app.InfoLog.Printf("gcal event %s rescheduled to %s", eventID, newStart.Format(time.RFC3339))
	return nil
}
