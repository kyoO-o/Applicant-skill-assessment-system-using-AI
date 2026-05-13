package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
)

const codeExpiry = 15 * time.Minute

func generateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// POST /pub/verify-email — confirm 6-digit code after registration
func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Code = strings.TrimSpace(req.Code)

	if req.Email == "" || req.Code == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл болон код шаардлагатай"})
		return
	}

	user, err := app.Users.GetWithEmail(req.Email)
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл олдсонгүй"})
		return
	}

	if user.EmailVerified {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл аль хэдийн баталгаажсан байна"})
		return
	}

	if user.VerifyCode != req.Code {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Код буруу байна"})
		return
	}

	if user.VerifyCodeExpiry == nil || time.Now().After(*user.VerifyCodeExpiry) {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Кодны хугацаа дууссан байна. Шинэ код авна уу."})
		return
	}

	user.EmailVerified = true
	user.VerifyCode = ""
	user.VerifyCodeExpiry = nil
	saved, err := app.Users.Save(user)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	startSession(r, saved, "local")
	oapi.SendResp(w, &authResponse{User: saved, Message: "И-мэйл амжилттай баталгаажлаа"})
}

// POST /pub/resend-verification — resend registration verification code
func resendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := app.Users.GetWithEmail(req.Email)
	if err != nil {
		// Don't reveal whether email exists
		oapi.SendResp(w, map[string]string{"message": "Хэрэв и-мэйл хаяг бүртгэлтэй бол шинэ код илгээгдэнэ."})
		return
	}

	if user.EmailVerified {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл аль хэдийн баталгаажсан байна"})
		return
	}

	code := generateOTP()
	expiry := time.Now().Add(codeExpiry)
	user.VerifyCode = code
	user.VerifyCodeExpiry = &expiry
	if _, err := app.Users.Save(user); err != nil {
		oapi.ServerError(w, err)
		return
	}

	go app.Mailer.SendVerificationEmail(user.Email, user.FullName, code)
	oapi.SendResp(w, map[string]string{"message": "Баталгаажуулах код дахин илгээгдлээ."})
}

// POST /pub/forgot-password — send password reset code to email
func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := app.Users.GetWithEmail(req.Email)
	if err != nil {
		// Don't reveal whether email exists
		oapi.SendResp(w, map[string]string{"message": "Хэрэв и-мэйл хаяг бүртгэлтэй бол нууц үг сэргээх код илгээгдэнэ."})
		return
	}

	code := generateOTP()
	expiry := time.Now().Add(codeExpiry)
	user.ResetCode = code
	user.ResetCodeExpiry = &expiry
	if _, err := app.Users.Save(user); err != nil {
		oapi.ServerError(w, err)
		return
	}

	go app.Mailer.SendPasswordResetEmail(user.Email, user.FullName, code)
	oapi.SendResp(w, map[string]string{"message": "Нууц үг сэргээх код илгээгдлээ."})
}

// POST /pub/reset-password — set new password using reset code
func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Code = strings.TrimSpace(req.Code)

	if req.Email == "" || req.Code == "" || req.Password == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Бүх талбарыг бөглөнө үү"})
		return
	}

	if len(req.Password) < 8 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Нууц үг наад зах нь 8 тэмдэгттэй байна"})
		return
	}

	user, err := app.Users.GetWithEmail(req.Email)
	if err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл олдсонгүй"})
		return
	}

	if user.ResetCode != req.Code || user.ResetCode == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Код буруу байна"})
		return
	}

	if user.ResetCodeExpiry == nil || time.Now().After(*user.ResetCodeExpiry) {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Кодны хугацаа дууссан байна. Шинэ код авна уу."})
		return
	}

	hash, err := common.HashPassword(req.Password)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	user.PasswordHash = hash
	user.ResetCode = ""
	user.ResetCodeExpiry = nil
	if _, err := app.Users.Save(user); err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, map[string]string{"message": "Нууц үг амжилттай шинэчлэгдлээ."})
}

// PUT /api/me/email — initiate email change; sends code to new address
func initiateEmailChangeHandler(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Email == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл хаяг шаардлагатай"})
		return
	}

	if req.Email == user.Email {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Энэ и-мэйл хаяг аль хэдийн ашиглагдаж байна"})
		return
	}

	if _, err := app.Users.GetWithEmail(req.Email); err == nil {
		oapi.CustomError(w, http.StatusConflict, map[string]string{"message": "Энэ и-мэйл хаяг бүртгэлтэй байна"})
		return
	}

	code := generateOTP()
	expiry := time.Now().Add(codeExpiry)
	user.PendingEmail = req.Email
	user.VerifyCode = code
	user.VerifyCodeExpiry = &expiry
	if _, err := app.Users.Save(user); err != nil {
		oapi.ServerError(w, err)
		return
	}

	go app.Mailer.SendEmailChangeEmail(req.Email, user.FullName, code)
	oapi.SendResp(w, map[string]string{"message": "Баталгаажуулах код шинэ и-мэйл хаягт илгээгдлээ."})
}

// POST /api/me/verify-email-change — confirm email change with code
func verifyEmailChangeHandler(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	var req struct {
		Code string `json:"code"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	req.Code = strings.TrimSpace(req.Code)

	if user.PendingEmail == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "И-мэйл солих хүсэлт эхлүүлээгүй байна"})
		return
	}

	if user.VerifyCode != req.Code || user.VerifyCode == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Код буруу байна"})
		return
	}

	if user.VerifyCodeExpiry == nil || time.Now().After(*user.VerifyCodeExpiry) {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Кодны хугацаа дууссан байна"})
		return
	}

	user.Email = user.PendingEmail
	user.PendingEmail = ""
	user.VerifyCode = ""
	user.VerifyCodeExpiry = nil
	saved, err := app.Users.Save(user)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, saved)
}

// PUT /api/me/password — change password (current password required)
func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := authUser(r)
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.CurrentPassword == "" || req.NewPassword == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Бүх талбарыг бөглөнө үү"})
		return
	}

	if len(req.NewPassword) < 8 {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Нууц үг наад зах нь 8 тэмдэгттэй байна"})
		return
	}

	if user.PasswordHash == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Энэ бүртгэлд нууц үг тохируулаагүй байна"})
		return
	}

	if common.CheckPassword(req.CurrentPassword, user.PasswordHash) != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Одоогийн нууц үг буруу байна"})
		return
	}

	hash, err := common.HashPassword(req.NewPassword)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	user.PasswordHash = hash
	if _, err := app.Users.Save(user); err != nil {
		oapi.ServerError(w, err)
		return
	}

	oapi.SendResp(w, map[string]string{"message": "Нууц үг амжилттай солигдлоо"})
}
