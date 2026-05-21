package mailerman

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Service struct {
	DB       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
	smtp     SMTPConfig
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger, smtpCfg SMTPConfig) *Service {
	return &Service{
		DB:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
		smtp:     smtpCfg,
	}
}

func (s *Service) isConfigured() bool {
	return s.smtp.Host != "" &&
		s.smtp.Username != "" &&
		!strings.HasPrefix(s.smtp.Username, "your-") &&
		s.smtp.Password != "" &&
		s.smtp.Password != "your-app-password"
}

func (s *Service) SendEmail(to, subject, body string) error {
	return s.sendEmailWithReplyTo(to, "", subject, body)
}

func (s *Service) sendEmailWithReplyTo(to, replyTo, subject, body string) error {
	if !s.isConfigured() {
		s.infoLog.Printf("[EMAIL] To: %s | Subject: %s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", s.smtp.Username, s.smtp.Password, s.smtp.Host)
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8", s.smtp.From, to, subject)
	if replyTo != "" {
		headers += "\r\nReply-To: " + replyTo
	}
	msg := []byte(headers + "\r\n\r\n" + body)

	addr := fmt.Sprintf("%s:%d", s.smtp.Host, s.smtp.Port)
	if err := smtp.SendMail(addr, auth, s.smtp.From, []string{to}, msg); err != nil {
		s.errorLog.Printf("failed to send email to %s: %v", to, err)
		return err
	}
	s.infoLog.Printf("email sent to %s (subject: %s)", to, subject)
	return nil
}

func (s *Service) SendVerificationEmail(to, name, code string) {
	// Always log the code so it's visible in the terminal during development
	s.infoLog.Printf("============================================================")
	s.infoLog.Printf("[VERIFY EMAIL] To: %s | Code: %s", to, code)
	s.infoLog.Printf("============================================================")

	if !s.isConfigured() {
		return
	}

	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px">
  <h2 style="margin-bottom:8px">Э-мэйл баталгаажуулалт</h2>
  <p>Сайн байна уу, %s!</p>
  <p>Таны бүртгэлийн баталгаажуулах код:</p>
  <div style="font-size:36px;font-weight:bold;letter-spacing:10px;text-align:center;padding:24px;background:#f4f4f5;border-radius:12px;margin:24px 0;color:#18181b">
    %s
  </div>
  <p style="color:#71717a;font-size:14px">Код <strong>15 минут</strong>ын дотор хүчинтэй.</p>
  <p style="color:#71717a;font-size:14px">Та энэ бүртгэлийг хийгээгүй бол энэ мэйлийг үл тоомсорлоно уу.</p>
</div>`, name, code)

	_ = s.SendEmail(to, "Э-мэйл хаягаа баталгаажуулна уу", body)
}

func (s *Service) SendPasswordResetEmail(to, name, code string) {
	s.infoLog.Printf("============================================================")
	s.infoLog.Printf("[RESET PASSWORD] To: %s | Code: %s", to, code)
	s.infoLog.Printf("============================================================")

	if !s.isConfigured() {
		return
	}

	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px">
  <h2 style="margin-bottom:8px">Нууц үг сэргээх</h2>
  <p>Сайн байна уу, %s!</p>
  <p>Нууц үгээ сэргээх код:</p>
  <div style="font-size:36px;font-weight:bold;letter-spacing:10px;text-align:center;padding:24px;background:#f4f4f5;border-radius:12px;margin:24px 0;color:#18181b">
    %s
  </div>
  <p style="color:#71717a;font-size:14px">Код <strong>15 минут</strong>ын дотор хүчинтэй.</p>
  <p style="color:#71717a;font-size:14px">Та нууц үгээ сэргээхийг хүсэлт гаргаагүй бол энэ мэйлийг үл тоомсорлоно уу.</p>
</div>`, name, code)

	_ = s.SendEmail(to, "Нууц үг сэргээх код", body)
}

func (s *Service) SendEmailChangeEmail(to, name, code string) {
	s.infoLog.Printf("============================================================")
	s.infoLog.Printf("[CHANGE EMAIL] To: %s | Code: %s", to, code)
	s.infoLog.Printf("============================================================")

	if !s.isConfigured() {
		return
	}

	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:480px;margin:0 auto;padding:24px">
  <h2 style="margin-bottom:8px">Э-мэйл хаяг солих</h2>
  <p>Сайн байна уу, %s!</p>
  <p>Шинэ и-мэйл хаягаа баталгаажуулах код:</p>
  <div style="font-size:36px;font-weight:bold;letter-spacing:10px;text-align:center;padding:24px;background:#f4f4f5;border-radius:12px;margin:24px 0;color:#18181b">
    %s
  </div>
  <p style="color:#71717a;font-size:14px">Код <strong>15 минут</strong>ын дотор хүчинтэй.</p>
  <p style="color:#71717a;font-size:14px">Та и-мэйл хаягаа солихыг хүсэлт гаргаагүй бол энэ мэйлийг үл тоомсорлоно уу.</p>
</div>`, name, code)

	_ = s.SendEmail(to, "Э-мэйл хаяг солих баталгаажуулах код", body)
}

func (s *Service) SendInterviewInviteEmail(to, applicantName, recruiterName, recruiterEmail, jobTitle string, interviewAt time.Time, location, note string) {
	// Format for Google Calendar: YYYYMMDDTHHmmSS / YYYYMMDDTHHmmSS (1 hour duration)
	gcalFmt := "20060102T150405"
	startStr := interviewAt.Format(gcalFmt)
	endStr := interviewAt.Add(time.Hour).Format(gcalFmt)
	eventTitle := url.QueryEscape(fmt.Sprintf("Ярилцлага: %s", jobTitle))
	details := url.QueryEscape(fmt.Sprintf("Ажил олгогч: %s\nАнкет илгээгч: %s\n%s", recruiterName, applicantName, note))
	gcalLocation := url.QueryEscape(location)
	gcalLink := fmt.Sprintf(
		"https://calendar.google.com/calendar/render?action=TEMPLATE&text=%s&dates=%s/%s&details=%s&location=%s",
		eventTitle, startStr, endStr, details, gcalLocation,
	)

	dateStr := interviewAt.Format("2006-01-02 15:04")

	s.infoLog.Printf("============================================================")
	s.infoLog.Printf("[INTERVIEW INVITE] To: %s | Job: %s | Date: %s | Location: %s", to, jobTitle, dateStr, location)
	s.infoLog.Printf("============================================================")

	if !s.isConfigured() {
		return
	}

	locationLine := ""
	if location != "" {
		locationLine = fmt.Sprintf(`<p style="margin:4px 0"><strong>Байршил:</strong> %s</p>`, location)
	}
	noteLine := ""
	if note != "" {
		noteLine = fmt.Sprintf(`<p style="margin:4px 0"><strong>Нэмэлт мэдээлэл:</strong> %s</p>`, note)
	}

	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:560px;margin:0 auto;padding:24px">
  <h2 style="margin-bottom:8px;color:#18181b">Ярилцлагын урилга</h2>
  <p>Сайн байна уу, <strong>%s</strong>!</p>
  <p>Та <strong>%s</strong> ажлын байрт анкет илгээсэн бөгөөд ярилцлагад урьж байна.</p>
  <div style="background:#f4f4f5;border-radius:12px;padding:20px;margin:24px 0">
    <p style="margin:4px 0"><strong>Ажлын байр:</strong> %s</p>
    <p style="margin:4px 0"><strong>Огноо, цаг:</strong> %s</p>
    %s
    %s
  </div>
  <a href="%s" target="_blank"
     style="display:inline-block;background:#4285F4;color:#fff;text-decoration:none;padding:12px 24px;border-radius:8px;font-weight:bold;margin-top:8px">
    Google Calendar-д нэмэх
  </a>
  <p style="margin-top:24px;color:#71717a;font-size:13px">
    Энэхүү урилга нь <strong>%s</strong>-аас ирсэн болно. Асуулт байвал хариу мэйл илгээнэ үү.
  </p>
</div>`, applicantName, jobTitle, jobTitle, dateStr, locationLine, noteLine, gcalLink, recruiterName)

	_ = s.sendEmailWithReplyTo(to, recruiterEmail, fmt.Sprintf("Ярилцлагын урилга: %s", jobTitle), body)
}

func (s *Service) SendInterviewConfirmToRecruiter(recruiterEmail, recruiterName, applicantName, applicantEmail, jobTitle string, interviewAt time.Time, location, note string) {
	dateStr := interviewAt.Format("2006-01-02 15:04")

	s.infoLog.Printf("============================================================")
	s.infoLog.Printf("[INTERVIEW CONFIRM] To recruiter: %s | Applicant: %s | Date: %s", recruiterEmail, applicantName, dateStr)
	s.infoLog.Printf("============================================================")

	if !s.isConfigured() {
		return
	}

	locationLine := ""
	if location != "" {
		locationLine = fmt.Sprintf(`<p style="margin:4px 0"><strong>Байршил:</strong> %s</p>`, location)
	}
	noteLine := ""
	if note != "" {
		noteLine = fmt.Sprintf(`<p style="margin:4px 0"><strong>Нэмэлт мэдээлэл:</strong> %s</p>`, note)
	}

	body := fmt.Sprintf(`
<div style="font-family:Arial,sans-serif;max-width:560px;margin:0 auto;padding:24px">
  <h2 style="margin-bottom:8px;color:#18181b">Ярилцлага товлогдлоо</h2>
  <p>Сайн байна уу, <strong>%s</strong>!</p>
  <p><strong>%s</strong> (<a href="mailto:%s">%s</a>)-д ярилцлагын урилга илгээгдлээ.</p>
  <div style="background:#f4f4f5;border-radius:12px;padding:20px;margin:24px 0">
    <p style="margin:4px 0"><strong>Ажлын байр:</strong> %s</p>
    <p style="margin:4px 0"><strong>Горилогч:</strong> %s</p>
    <p style="margin:4px 0"><strong>Огноо, цаг:</strong> %s</p>
    %s
    %s
  </div>
</div>`, recruiterName, applicantName, applicantEmail, applicantEmail, jobTitle, applicantName, dateStr, locationLine, noteLine)

	_ = s.SendEmail(recruiterEmail, fmt.Sprintf("Ярилцлага товлогдлоо: %s — %s", applicantName, jobTitle), body)
}
