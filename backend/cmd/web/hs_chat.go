package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/aiman"
)

type chatRequest struct {
	Message string             `json:"message"`
	History []aiman.ChatMessage `json:"history"`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

// POST /api/chat
func chat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Хүсэлт буруу байна"})
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Мессеж хоосон байна"})
		return
	}

	// Fetch active jobs to give chatbot context
	jobs, err := app.Jobs.ListActive(nil)
	if err != nil {
		oapi.ServerError(w, err)
		return
	}

	var sb strings.Builder
	for _, job := range jobs {
		skills := make([]string, 0, len(job.Skills))
		for _, s := range job.Skills {
			skills = append(skills, s.Name)
		}
		sb.WriteString(fmt.Sprintf("- %s | Байршил: %s | Цалин: %.0f-%.0f ₮ | Ур чадвар: %s\n",
			job.Title, job.Location, job.MinSalary, job.MaxSalary, strings.Join(skills, ", ")))
	}

	reply, err := app.AI.Chat(req.Message, sb.String(), req.History)
	if err != nil {
		app.ErrorLog.Printf("chat error: %v", err)
		oapi.CustomError(w, http.StatusServiceUnavailable, map[string]string{"message": "AI үйлчилгээ түр ажиллахгүй байна"})
		return
	}

	oapi.SendResp(w, chatResponse{Reply: reply})
}
