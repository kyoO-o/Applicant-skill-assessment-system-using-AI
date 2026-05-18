package socket

import (
	"errors"
	"net/http"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/websocket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
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

func DeleteCustomerConnection(customerID int) {
	app.CustomerConnectionMutex.Lock()
	defer app.CustomerConnectionMutex.Unlock()
	if ct := app.CustomerConnections[customerID]; ct != nil {
		delete(app.CustomerConnections, customerID)
	}
}
