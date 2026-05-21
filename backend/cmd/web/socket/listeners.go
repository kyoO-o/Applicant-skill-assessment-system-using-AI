package socket

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/websocket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

var (
	remindedMu  sync.Mutex
	remindedIDs = map[int]bool{}
)

func NotifyUser(userID int, title, body, notifType string) {
	conn := GetCustomerConnection(userID)
	if conn == nil {
		return
	}
	payload := map[string]string{"title": title, "body": body, "type": notifType}
	if err := conn.Send("Notification", payload); err != nil {
		app.ErrorLog.Printf("notify user %d: %v", userID, err)
	}
}

func OnFrontendWSConnect(r *http.Request, conn *websocket.Connection) error {
	user := r.Context().Value(app.ContextKeyAuthCustomer).(*userman.User)
	if user == nil {
		return errors.New("user is nil")
	}
	app.InfoLog.Printf("OnFrontendWSConnect user=%d", user.ID)

	// Replace any existing connection so browser refresh / reconnects work.
	SetCustomerConnection(user.ID, conn)
	if err := conn.Send("Connected", user); err != nil {
		return err
	}

	conn.OnMessage = func(m websocket.Message) {}
	conn.OnClose = onCustomerWSClose(user.ID, conn)

	return nil
}

// onCustomerWSClose only removes the entry when this specific connection is
// still the active one, preventing a stale close from evicting a newer conn.
func onCustomerWSClose(userID int, conn *websocket.Connection) func() {
	return func() {
		app.CustomerConnectionMutex.Lock()
		defer app.CustomerConnectionMutex.Unlock()
		if app.CustomerConnections[userID] == conn {
			delete(app.CustomerConnections, userID)
		}
	}
}

func GetCustomerConnection(customerID int) *websocket.Connection {
	app.CustomerConnectionMutex.RLock()
	defer app.CustomerConnectionMutex.RUnlock()
	return app.CustomerConnections[customerID]
}

func SetCustomerConnection(customerID int, conn *websocket.Connection) {
	app.CustomerConnectionMutex.Lock()
	defer app.CustomerConnectionMutex.Unlock()
	app.CustomerConnections[customerID] = conn
}

// StartInterviewReminderScheduler sends a 1-hour-before notification to both
// the applicant and recruiter for every scheduled interview. It ticks every
// minute and uses an in-memory set to avoid duplicate reminders.
func StartInterviewReminderScheduler() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			from := now.Add(58 * time.Minute)
			to := now.Add(62 * time.Minute)
			apps, err := app.Applications.FindUpcomingInterviews(from, to)
			if err != nil {
				app.ErrorLog.Printf("interview reminder query: %v", err)
				continue
			}
			for _, a := range apps {
				remindedMu.Lock()
				already := remindedIDs[a.ID]
				if !already {
					remindedIDs[a.ID] = true
				}
				remindedMu.Unlock()
				if already {
					continue
				}

				NotifyUser(int(a.ApplicantID),
					"Ярилцлага удахгүй эхлэх гэж байна",
					"1 цагийн дараа ярилцлага эхлэх тул бэлдэнэ үү",
					"interview_reminder")

				job, err := app.Jobs.Get(int(a.JobPostingID))
				if err == nil {
					loc := a.InterviewLocation
					if loc == "" {
						loc = "байршил тодорхойгүй"
					}
					NotifyUser(int(job.PostedBy),
						"Ярилцлага удахгүй эхлэх гэж байна",
						fmt.Sprintf("1 цагийн дараа ярилцлага байна (%s)", loc),
						"interview_reminder")
				}
			}
		}
	}()
}

func DeleteCustomerConnection(customerID int) {
	app.CustomerConnectionMutex.Lock()
	defer app.CustomerConnectionMutex.Unlock()
	if ct := app.CustomerConnections[customerID]; ct != nil {
		delete(app.CustomerConnections, customerID)
	}
}
