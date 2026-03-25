package app

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/apputils"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/websocket"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/mailerman"
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
	DB                      *gorm.DB
	// OCRQueue                chan *documentman.Document

	// #region Services
	Users  *userman.Service
	Mailer *mailerman.Service
)

const (
	GB = 1 << 30
	MB = 1 << 20
	KB = 1 << 10
)

//#region Init

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
	Mailer = mailerman.NewService(DB, InfoLog, ErrorLog)

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

}
func Close() {
}

// #region use with caution
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
