package socket

import (
	"errors"
	"net/http"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/websocket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

func OnFrontendWSConnect(r *http.Request, conn *websocket.Connection) error {
	user := r.Context().Value(app.ContextKeyAuthCustomer).(*userman.User)
	if user == nil {
		return errors.New("user is nil")
	}
	app.InfoLog.Println("OnFrontendWSConnect", user)
	existingConn := app.CustomerConnections[user.ID]
	if existingConn != nil {
		return errors.New("already connected")
	}

	SetCustomerConnection(user.ID, conn)
	if err := conn.Send("Connected", user); err != nil {
		return err
	}

	conn.OnMessage = func(m websocket.Message) {

	}

	conn.OnClose = onCustomerWSClose(user)

	return nil
}

func onCustomerWSClose(user *userman.User) func() {
	return func() {
		DeleteCustomerConnection(user.ID)
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
