package main

import (
	"errors"
	"flag"
	"net/http"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/socket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

func main() {

	mode := flag.String("mode", "debug", "Choose mode. debug, test or production")
	configPath := flag.String("conf", "../confs/web.yaml", "Configuration file path")
	flag.Parse()

	app.Init(*configPath, *mode)
	defer app.Close()

	if err := app.DB.AutoMigrate(
		new(userman.User),
	); err != nil {
		app.ErrorLog.Panic(err)
	}

	if _, err := app.Users.Get(1); err != nil {
		if !errors.Is(err, userman.ErrNotFound) {
			app.ErrorLog.Panic(err) // If error is not plan not found, panic
		}
		adminUsers := &userman.User{
			Email:       "anujinnn.ts@gmail.com",
			Name:        "Anujin",
			PhoneNumber: "99560628",
			ToduID:      228757,
			Role:        "superadmin",
		}

		if _, err := app.Users.Save(adminUsers); err != nil {
			app.ErrorLog.Fatalf("failed to insert free plan: %v", err)
		}
	}

	app.FrontendWS.OnConnect = socket.OnFrontendWSConnect

	// ele, ok := queue.Dequeue()
	// if doc, isDoc := ele.(documentman.Document); ok && isDoc {
	// 	fmt.Println("Queue After removing", doc, ":", queue.Data)
	// }
	srv := &http.Server{
		Addr:         app.Config.Port,
		ErrorLog:     app.ErrorLog,
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		// MaxHeaderBytes: 1 << 20,
		MaxHeaderBytes: 300 << 20,
	}

	app.InfoLog.Printf("Starting server on %s", app.Config.Port)
	app.ErrorLog.Fatal(srv.ListenAndServe())
}
