package aiman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	claudeModel  = "claude-haiku-4-5-20251001"
	claudeAPIURL = "https://api.anthropic.com/v1/messages"
)

// AssessmentResult is the structured JSON returned by Claude for CV scoring.
type AssessmentResult struct {
	OverallScore    int           `json:"overall_score"`
	Summary         string        `json:"summary"`
	MatchedSkills   []SkillResult `json:"matched_skills"`
	MissingSkills   []SkillResult `json:"missing_skills"`
	Recommendations []string      `json:"recommendations"`
}

type SkillResult struct {
	Skill       string `json:"skill"`
	Explanation string `json:"explanation"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey, httpClient: &http.Client{}}
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	System    string          `json:"system"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type claudeResponse struct {
	Content []claudeContent `json:"content"`
}

func (c *Client) call(system, userMsg string) (string, error) {
	reqBody := claudeRequest{
		Model:     claudeModel,
		MaxTokens: 1024,
		System:    system,
		Messages:  []claudeMessage{{Role: "user", Content: userMsg}},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, claudeAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API error %d: %s", resp.StatusCode, string(respBytes))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBytes, &claudeResp); err != nil {
		return "", err
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return claudeResp.Content[0].Text, nil
}

// AssessCV compares a CV against job requirements and returns a structured assessment.
func (c *Client) AssessCV(cvText, jobTitle string, requirements, skills, duties []string) (*AssessmentResult, error) {
	jobContext := fmt.Sprintf("Ажлын байрны нэр: %s\n\nШаардлагууд:\n%s\n\nУр чадварууд:\n%s\n\nҮүрэг хариуцлага:\n%s",
		jobTitle,
		strings.Join(requirements, "\n"),
		strings.Join(skills, "\n"),
		strings.Join(duties, "\n"),
	)

	system := `Та ажил горилогчийн CV болон ажлын байрны шаардлагыг харьцуулан дүгнэх AI үнэлгээний систем.
Хариултыг заавал дараах JSON форматаар өг, өөр ямар ч текст бүү нэм:
{
  "overall_score": <0-100 хооронд бүхэл тоо>,
  "summary": "<Монгол хэл дээр 2-3 өгүүлбэрийн дүгнэлт>",
  "matched_skills": [{"skill": "<ур чадвар>", "explanation": "<тайлбар>"}, ...],
  "missing_skills": [{"skill": "<дутуу ур чадвар>", "explanation": "<тайлбар>"}, ...],
  "recommendations": ["<зөвлөмж 1>", "<зөвлөмж 2>", ...]
}`

	userMsg := fmt.Sprintf("Ажлын байрны мэдээлэл:\n%s\n\n---\n\nCV текст:\n%s", jobContext, cvText)

	text, err := c.call(system, userMsg)
	if err != nil {
		return nil, err
	}

	// Strip markdown code block if present
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result AssessmentResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %w\nraw: %s", err, text)
	}

	return &result, nil
}

// Chat sends a user message along with available jobs context and returns a response.
func (c *Client) Chat(userMessage string, jobsContext string, history []ChatMessage) (string, error) {
	system := `Та ажил хайгч хэрэглэгчдэд туслах AI чатбот. Монгол хэл дээр хариулна уу.
Доорх ажлын байруудын мэдээлэл дээр тулгуурлан хэрэглэгчийн асуултад хариулна уу.
Хэрэглэгч цалин, байршил, чиглэлийн дагуу тохирох ажлыг хайхад туслаарай.
Хариулт нь тодорхой, хялбар ойлгомжтой байх ёстой.`

	// Build conversation
	var messages []claudeMessage
	for _, h := range history {
		messages = append(messages, claudeMessage{Role: h.Role, Content: h.Content})
	}

	content := userMessage
	if jobsContext != "" {
		content = fmt.Sprintf("Одоогийн ажлын байрны мэдээлэл:\n%s\n\n---\nХэрэглэгчийн асуулт: %s", jobsContext, userMessage)
	}
	messages = append(messages, claudeMessage{Role: "user", Content: content})

	reqBody := claudeRequest{
		Model:     claudeModel,
		MaxTokens: 800,
		System:    system,
		Messages:  messages,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, claudeAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API error %d: %s", resp.StatusCode, string(respBytes))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBytes, &claudeResp); err != nil {
		return "", err
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return claudeResp.Content[0].Text, nil
}

// GenerateTask generates a task description for a job posting.
func (c *Client) GenerateTask(jobTitle string, requirements, skills []string) (string, string, error) {
	system := `Та ажил олгогчид туслах AI систем. Ажлын байрны шаардлагад тулгуурлан практик даалгавар үүсгэнэ.
JSON форматаар хариул:
{"title": "<даалгаврын гарчиг>", "description": "<дэлгэрэнгүй даалгаврын тайлбар, монгол хэл дээр>"}`

	userMsg := fmt.Sprintf("Ажлын байр: %s\nШаардлагууд: %s\nУр чадварууд: %s\n\nЭнэ ажлын байранд тохирсон практик даалгавар үүсгэнэ үү.",
		jobTitle,
		strings.Join(requirements, ", "),
		strings.Join(skills, ", "),
	)

	text, err := c.call(system, userMsg)
	if err != nil {
		return "", "", err
	}

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return "", "", fmt.Errorf("failed to parse Claude response: %w", err)
	}

	return result.Title, result.Description, nil
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
