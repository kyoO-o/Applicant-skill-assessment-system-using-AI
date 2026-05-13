package mailerman

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

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
	if !s.isConfigured() {
		s.infoLog.Printf("[EMAIL] To: %s | Subject: %s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", s.smtp.Username, s.smtp.Password, s.smtp.Host)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.smtp.From, to, subject, body,
	))

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
  <h2 style="margin-bottom:8px">И-мэйл баталгаажуулалт</h2>
  <p>Сайн байна уу, %s!</p>
  <p>Таны бүртгэлийн баталгаажуулах код:</p>
  <div style="font-size:36px;font-weight:bold;letter-spacing:10px;text-align:center;padding:24px;background:#f4f4f5;border-radius:12px;margin:24px 0;color:#18181b">
    %s
  </div>
  <p style="color:#71717a;font-size:14px">Код <strong>15 минут</strong>ын дотор хүчинтэй.</p>
  <p style="color:#71717a;font-size:14px">Та энэ бүртгэлийг хийгээгүй бол энэ мэйлийг үл тоомсорлоно уу.</p>
</div>`, name, code)

	_ = s.SendEmail(to, "И-мэйл хаягаа баталгаажуулна уу", body)
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
  <h2 style="margin-bottom:8px">И-мэйл хаяг солих</h2>
  <p>Сайн байна уу, %s!</p>
  <p>Шинэ и-мэйл хаягаа баталгаажуулах код:</p>
  <div style="font-size:36px;font-weight:bold;letter-spacing:10px;text-align:center;padding:24px;background:#f4f4f5;border-radius:12px;margin:24px 0;color:#18181b">
    %s
  </div>
  <p style="color:#71717a;font-size:14px">Код <strong>15 минут</strong>ын дотор хүчинтэй.</p>
  <p style="color:#71717a;font-size:14px">Та и-мэйл хаягаа солихыг хүсэлт гаргаагүй бол энэ мэйлийг үл тоомсорлоно уу.</p>
</div>`, name, code)

	_ = s.SendEmail(to, "И-мэйл хаяг солих баталгаажуулах код", body)
}
