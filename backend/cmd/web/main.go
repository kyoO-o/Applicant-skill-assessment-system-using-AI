package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/socket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/appman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/companyman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/cvman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/jobman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/taskman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

func main() {

	mode := flag.String("mode", "debug", "Choose mode. debug, test or production")
	configPath := flag.String("conf", "../confs/web.yaml", "Configuration file path")
	flag.Parse()

	app.Init(*configPath, *mode)
	defer app.Close()

	if err := app.DB.AutoMigrate(
		new(jobman.JobPosting),
		new(jobman.Duty),
		new(jobman.Requirement),
		new(jobman.Skill),
		new(jobman.Bonus),
		new(userman.User),
		new(companyman.Company),
		new(companyman.CompanyBenefit),
		new(appman.Application),
		new(taskman.Task),
		new(taskman.TaskSubmission),
		new(cvman.CVProfile),
	); err != nil {
		app.ErrorLog.Panic(err)
	}

	// Mark users created before email verification was introduced as already verified
	app.DB.Exec("UPDATE users SET email_verified = true WHERE email_verified = false AND (verify_code = '' OR verify_code IS NULL)")

	// Ensure CV upload directory exists
	if app.Config.StoragePath != "" {
		if err := os.MkdirAll(app.Config.StoragePath+"/cvs", 0755); err != nil {
			app.ErrorLog.Printf("warning: could not create storage dir: %v", err)
		}
	}

	if _, err := app.Users.Get(1); err != nil {
		if !errors.Is(err, userman.ErrNotFound) {
			app.ErrorLog.Panic(err) // If error is not plan not found, panic
		}
		adminUsers := &userman.User{
			Email:       "anujinnn.ts@gmail.com",
			FullName:    "Anujin",
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
