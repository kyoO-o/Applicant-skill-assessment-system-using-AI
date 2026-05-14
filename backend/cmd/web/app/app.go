package app

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/apputils"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/websocket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/aiman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/appman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/cvman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/jobman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/mailerman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/taskman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	ErrorLog                *log.Logger
	InfoLog                 *log.Logger
	Session                 *sessions.Session
	Config                  = conf{}
	Mode                    string
	Location                *time.Location
	FrontendWS              *websocket.Websocket
	CustomerConnections     = map[int]*websocket.Connection{}
	CustomerConnectionMutex = new(sync.RWMutex)
	Todu                    *oauth2.Config
	GoogleCalendar          *oauth2.Config
	DB                      *gorm.DB

	// Services
	Users        *userman.Service
	Jobs         *jobman.Service
	Mailer       *mailerman.Service
	Applications *appman.Service
	Tasks        *taskman.Service
	CVProfiles   *cvman.Service
	AI           *aiman.Client
)

const (
	GB = 1 << 30
	MB = 1 << 20
	KB = 1 << 10
)

func Init(path, mode string) {
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	Mode = mode
	loc, err := time.LoadLocation("Asia/Ulaanbaatar")
	if err != nil {
		panic(err)
	}
	Location = loc

	apputils.LoadConfig(&Config, path, mode)

	DB = apputils.OpenDB(Config.DSN)

	Users = userman.NewService(DB, InfoLog, ErrorLog)
	Jobs = jobman.NewService(DB, InfoLog, ErrorLog)
	Mailer = mailerman.NewService(DB, InfoLog, ErrorLog, mailerman.SMTPConfig{
		Host:     Config.SMTP.Host,
		Port:     Config.SMTP.Port,
		Username: Config.SMTP.Username,
		Password: Config.SMTP.Password,
		From:     Config.SMTP.From,
	})
	Applications = appman.NewService(DB, InfoLog, ErrorLog)
	Tasks = taskman.NewService(DB, InfoLog, ErrorLog)
	CVProfiles = cvman.NewService(DB, InfoLog, ErrorLog)
	AI = aiman.NewClient(Config.AnthropicAPIKey)

	FrontendWS = websocket.New()

	Session = sessions.New([]byte(Config.SessionSecret))
	Session.Lifetime = 72 * time.Hour

	Todu = &oauth2.Config{
		ClientID:     Config.Todu.ClientID,
		ClientSecret: Config.Todu.ClientSecret,
		RedirectURL:  Config.Todu.RedirectURL,
		Scopes:       Config.Todu.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  Config.Todu.Endpoint.AuthURL,
			TokenURL: Config.Todu.Endpoint.TokenURL,
		},
	}

	if Config.Google.ClientID != "" && Config.Google.ClientSecret != "" {
		GoogleCalendar = &oauth2.Config{
			ClientID:     Config.Google.ClientID,
			ClientSecret: Config.Google.ClientSecret,
			RedirectURL:  Config.Google.RedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/calendar.events"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		}
	}
}

func Close() {
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintOnError(err error) {
	if err != nil {
		ErrorLog.Println(err)
	}
}
