package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// fetchChatbotModels proxies a GET /v1/models request to the AI provider
// so the admin UI can show the real model list.
func fetchChatbotModels(c *fiber.Ctx) error {
	var req struct {
		ApiUrl string `json:"api_url"`
		ApiKey string `json:"api_key"`
	}
	if err := c.BodyParser(&req); err != nil {
		log.Println("BodyParser error in fetchChatbotModels:", err, "Body:", string(c.Body()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}
	if req.ApiKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "API key is required"})
	}

	// Derive the /v1/models URL from the chat completions URL
	baseURL := req.ApiUrl
	if baseURL == "" {
		baseURL = "https://ai.sumopod.com/v1/chat/completions"
	}
	// Strip /chat/completions to get /v1/models
	modelsURL := strings.TrimSuffix(baseURL, "/chat/completions")
	modelsURL = strings.TrimSuffix(modelsURL, "/") + "/models"

	httpReq, err := http.NewRequest("GET", modelsURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to build request"})
	}
	httpReq.Header.Set("Authorization", "Bearer "+req.ApiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println("fetchChatbotModels error:", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Failed to connect to AI provider"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": fmt.Sprintf("Provider returned status %d", resp.StatusCode),
		})
	}

	// Forward the JSON response as-is
	c.Set("Content-Type", "application/json")
	return c.Send(body)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

type ChatbotSettings struct {
	Enabled        bool              `json:"enabled"`
	ApiKey         string            `json:"api_key"`
	ApiUrl         string            `json:"api_url"`
	Model          string            `json:"model"`
	Provider       string            `json:"provider"`
	SystemPrompt   string            `json:"system_prompt"`
	KnowledgeFiles []KnowledgeFile   `json:"knowledge_files"`
	MaxTokens      int               `json:"max_tokens"`
}

// buildFullPrompt combines the system prompt with knowledge file contents
func buildFullPrompt(settings ChatbotSettings) string {
	prompt := settings.SystemPrompt
	if len(settings.KnowledgeFiles) == 0 {
		return prompt
	}

	var sb strings.Builder
	sb.WriteString(prompt)
	sb.WriteString("\n\n--- KNOWLEDGE BASE (gunakan informasi berikut untuk menjawab pertanyaan) ---\n")
	for _, kf := range settings.KnowledgeFiles {
		if kf.Content == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("\n### Dokumen: %s\n", kf.Name))
		sb.WriteString(kf.Content)
		sb.WriteString("\n")
	}
	sb.WriteString("\n--- AKHIR KNOWLEDGE BASE ---")
	return sb.String()
}

func handleChat(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(req.Messages) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Messages array cannot be empty"})
	}

	// Fetch chatbot settings from the database
	var settingJSON []byte
	err := db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'chatbot'").Scan(&settingJSON)
	if err != nil {
		log.Println("Chatbot settings not found or error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Chatbot is not configured"})
	}

	var settings ChatbotSettings
	if err := json.Unmarshal(settingJSON, &settings); err != nil {
		log.Println("Failed to parse chatbot settings:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid chatbot settings"})
	}

	if !settings.Enabled {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Chatbot is currently disabled"})
	}

	if settings.ApiKey == "" {
		// Mock response if API key is not set
		time.Sleep(1 * time.Second)
		return c.JSON(fiber.Map{
			"status": "success",
			"reply":  "Maaf, fitur AI sedang dalam tahap simulasi karena API Key belum dikonfigurasi di dashboard admin.",
		})
	}

	switch settings.Provider {
	case "openai":
		settings.SystemPrompt = buildFullPrompt(settings)
		reply, err := callOpenAICompatible(settings, req.Messages)
		if err != nil {
			log.Println("OpenAI-compatible API Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to communicate with AI provider"})
		}
		return c.JSON(fiber.Map{"status": "success", "reply": reply})

	case "gemini":
		fullPrompt := buildFullPrompt(settings)
		reply, err := callGeminiAPI(settings.ApiKey, fullPrompt, settings.MaxTokens, req.Messages)
		if err != nil {
			log.Println("Gemini API Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to communicate with AI provider"})
		}
		return c.JSON(fiber.Map{"status": "success", "reply": reply})

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported AI provider"})
	}
}

// ─── OpenAI-Compatible API (Sumopod / LiteLLM / OpenAI / etc.) ───

func callOpenAICompatible(settings ChatbotSettings, messages []ChatMessage) (string, error) {
	apiURL := settings.ApiUrl
	if apiURL == "" {
		apiURL = "https://ai.sumopod.com/v1/chat/completions"
	}
	// Normalize: jika hanya base URL (misal /v1), tambahkan /chat/completions
	if !strings.HasSuffix(apiURL, "/chat/completions") {
		apiURL = strings.TrimSuffix(apiURL, "/") + "/chat/completions"
	}
	model := settings.Model
	if model == "" {
		model = "gpt-5-nano"
	}

	// Build messages array with system prompt
	var chatMessages []map[string]string
	if settings.SystemPrompt != "" {
		chatMessages = append(chatMessages, map[string]string{
			"role":    "system",
			"content": settings.SystemPrompt,
		})
	}
	for _, msg := range messages {
		if msg.Content == "" {
			continue
		}
		chatMessages = append(chatMessages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	maxTokens := settings.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1000
	}

	requestBody := map[string]interface{}{
		"model":    model,
		"messages": chatMessages,
		"temperature":      0.7,
		"max_tokens":       maxTokens,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.ApiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai-compatible api returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "Maaf, saya tidak bisa memberikan jawaban saat ini.", nil
}

// ─── Google Gemini Native API ───

type GeminiMessage struct {
	Role  string `json:"role"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

func callGeminiAPI(apiKey, systemPrompt string, maxTokens int, messages []ChatMessage) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey)

	var history []GeminiMessage
	for _, msg := range messages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		if msg.Content == "" {
			continue
		}
		history = append(history, GeminiMessage{
			Role: role,
			Parts: []struct {
				Text string `json:"text"`
			}{{Text: msg.Content}},
		})
	}

	if maxTokens <= 0 {
		maxTokens = 1000
	}

	requestBody := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": systemPrompt},
			},
		},
		"contents": history,
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,
			"maxOutputTokens": maxTokens,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "Maaf, saya tidak bisa memberikan jawaban saat ini.", nil
}

type AdminAIGenerateRequest struct {
	Prompt string `json:"prompt"`
	Style  string `json:"style"`
}

func handleAdminAIGenerate(c *fiber.Ctx) error {
	var req AdminAIGenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if strings.TrimSpace(req.Prompt) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Prompt/ide tulisan tidak boleh kosong"})
	}

	// Fetch chatbot settings from the database for API Key and Provider info
	var settingJSON []byte
	err := db.QueryRowContext(context.Background(), "SELECT value FROM site_settings WHERE key = 'chatbot'").Scan(&settingJSON)
	if err != nil {
		log.Println("Chatbot settings not found or error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI tidak terkonfigurasi di server"})
	}

	var settings ChatbotSettings
	if err := json.Unmarshal(settingJSON, &settings); err != nil {
		log.Println("Failed to parse chatbot settings:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Pengaturan AI tidak valid"})
	}

	// Determine style description
	styleDescription := "Artikel Umum (menarik dan informatif)"
	switch req.Style {
	case "news", "berita":
		styleDescription = "Berita Formal (5W+1H, bahasa jurnalisme baku, format berita)"
	case "article", "artikel":
		styleDescription = "Artikel Populer/Edukasi (sub-heading terstruktur, mengalir, informatif)"
	case "formal":
		styleDescription = "Pengumuman Resmi (sangat formal, sopan, padat, dan jelas)"
	case "casual", "santai":
		styleDescription = "Kasual/Santai (bahasa santai/populer namun sopan, cocok untuk mahasiswa/cerita kampus)"
	}

	if settings.ApiKey == "" {
		// Mock response if API key is not set
		time.Sleep(1 * time.Second)
		mockContent := fmt.Sprintf(`## Hasil Generate AI (Simulasi)

Ini adalah simulasi hasil generate tulisan berdasarkan input Anda:
*"%s"*

Dengan gaya penulisan: **%s**.

### Poin Utama
- Poin penting pertama terkait topik yang Anda berikan.
- Poin kedua yang mendetailkan informasi lebih lanjut.
- Kesimpulan dan ajakan bertindak (call to action).

*Catatan: Konfigurasikan API Key AI di menu Pengaturan Chatbot untuk mengaktifkan AI asli.*`, req.Prompt, styleDescription)
		return c.JSON(fiber.Map{
			"status": "success",
			"reply":  mockContent,
		})
	}

	// Custom copywriting system prompt for admin generation
	settings.SystemPrompt = "Kamu adalah seorang jurnalis, editor, dan copywriter profesional. Tugasmu adalah menulis artikel, berita, atau postingan lengkap berdasarkan poin-poin yang diberikan oleh user. Tulislah dalam Bahasa Indonesia yang baik dan benar. Output harus menggunakan format Markdown lengkap (dengan sub-heading, bold, dll jika sesuai). Jangan menuliskan intro/outro basa-basi seperti 'Tentu, berikut adalah...', langsung keluarkan isi tulisan/artikelnya saja dalam format Markdown."

	messages := []ChatMessage{
		{
			Role: "user",
			Content: fmt.Sprintf("Tolong buatkan artikel atau berita lengkap dalam format Markdown berdasarkan ide/poin-poin berikut:\n\"%s\"\n\nGaya penulisan yang diinginkan: %s.", req.Prompt, styleDescription),
		},
	}

	switch settings.Provider {
	case "openai":
		reply, err := callOpenAICompatible(settings, messages)
		if err != nil {
			log.Println("OpenAI AI Generate Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal berkomunikasi dengan AI provider"})
		}
		reply = cleanMarkdownCodeBlocks(reply)
		return c.JSON(fiber.Map{"status": "success", "reply": reply})

	case "gemini":
		reply, err := callGeminiAPI(settings.ApiKey, settings.SystemPrompt, settings.MaxTokens, messages)
		if err != nil {
			log.Println("Gemini AI Generate Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal berkomunikasi dengan AI provider"})
		}
		reply = cleanMarkdownCodeBlocks(reply)
		return c.JSON(fiber.Map{"status": "success", "reply": reply})

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Provider AI tidak didukung"})
	}
}

// cleanMarkdownCodeBlocks strips leading/trailing markdown block comments from the text.
func cleanMarkdownCodeBlocks(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```markdown") {
		s = strings.TrimPrefix(s, "```markdown")
		if strings.HasSuffix(s, "```") {
			s = strings.TrimSuffix(s, "```")
		}
	} else if strings.HasPrefix(s, "```html") {
		s = strings.TrimPrefix(s, "```html")
		if strings.HasSuffix(s, "```") {
			s = strings.TrimSuffix(s, "```")
		}
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		if strings.HasSuffix(s, "```") {
			s = strings.TrimSuffix(s, "```")
		}
	}
	return strings.TrimSpace(s)
}

