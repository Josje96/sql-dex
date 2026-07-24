// Package tutor implements the Phase 3 AI helper: a Socratic SQL coach that
// guides learners toward the answer with hints instead of handing over a
// finished query.
//
// It speaks the OpenAI-compatible "chat completions" API, so it works with many
// providers just by changing the base URL, model, and key — OpenAI, DeepSeek,
// SiliconFlow, OpenRouter, and Google Gemini (via its OpenAI-compatible
// endpoint) all work.
package tutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GoogleOpenAIBaseURL is Gemini's OpenAI-compatible endpoint, used as the
// default when only a Google key is supplied.
const GoogleOpenAIBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-4o-mini"
)

// systemPrompt shapes the model into a hint-giving tutor rather than an answer key.
const systemPrompt = `You are SQL-Dex Tutor, a patient SQL coach inside a learning game.
The learner writes SQLite queries against a Pokémon database (the "veekun" pokédex).

YOUR TEACHING RULES — follow them strictly:
- GUIDE, do not solve. Never output a complete, runnable query that fully answers
  the learner's task. Your goal is that THEY write it.
- Give one small step at a time: name the relevant table(s) or column(s), suggest
  which SQL clause comes next (WHERE, JOIN, GROUP BY, ...), or ask a guiding question.
- Tiny illustrative fragments are fine (e.g. "you'll need a JOIN ... ON a.id = b.a_id"),
  but never the full solution stitched together.
- If the learner's query has an error, point at the cause and nudge them to the fix
  instead of rewriting it for them.
- Be encouraging and concise. Prefer 2-5 short sentences or a couple of bullet points.
- Whenever you name a new SQL clause or concept (WHERE, JOIN, ORDER BY, GROUP BY,
  COUNT, ...), include a SHORT analogous example that uses DIFFERENT tables or a
  different goal than the learner's task — so they see the pattern in action
  without being handed their exact answer. Introduce it with "For example,".
- You have the recent conversation for context; build on earlier hints instead of
  repeating them, and acknowledge progress the learner has made.

USEFUL SCHEMA NOTES:
- pokemon_species(id, identifier, generation_id, ...): one row per species; id is the
  national-dex number. identifier is the lowercase English name.
- pokemon_species_names(pokemon_species_id, local_language_id, name): display names.
  Use local_language_id = 9 for English. Most *_names / type_names tables follow this.
- pokemon(id, species_id, height, weight, base_experience): concrete Pokémon; join to
  a species via pokemon.species_id = pokemon_species.id.
- pokemon_types(pokemon_id, type_id, slot) + types + type_names: a Pokémon's types.
- pokemon_stats(pokemon_id, stat_id, base_stat) + stats + stat_names: base stats.
- pokemon_moves(pokemon_id, version_group_id, move_id, ...) + moves + move_names: moves.
- generations(id): 1 = Kanto/Gen 1, etc.
The database is READ-ONLY, so only SELECT queries work.`

// historyWindow is how many prior messages of the conversation we replay to
// give the tutor a rolling short-term memory.
const historyWindow = 5

// Config selects the provider and model for the tutor.
type Config struct {
	APIKey  string // provider API key
	BaseURL string // OpenAI-compatible base URL, e.g. https://api.deepseek.com/v1
	Model   string // model name, e.g. deepseek-chat
}

// Message is one prior turn in the conversation. Role is "user" or "tutor".
type Message struct {
	Role string
	Text string
}

// Tutor talks to an OpenAI-compatible chat API on the learner's behalf.
type Tutor struct {
	apiKey  string
	baseURL string
	model   string
	http    *http.Client
}

// New builds a Tutor from cfg, filling in sensible defaults. An empty API key
// means the tutor is disabled; callers should check Enabled.
func New(cfg Config) *Tutor {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" {
		base = defaultBaseURL
	}
	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = defaultModel
	}
	return &Tutor{
		apiKey:  strings.TrimSpace(cfg.APIKey),
		baseURL: base,
		model:   model,
		http:    &http.Client{Timeout: 45 * time.Second},
	}
}

// Enabled reports whether an API key is configured.
func (t *Tutor) Enabled() bool { return t.apiKey != "" }

// Describe returns a short "model via host" summary for startup logs (no key).
func (t *Tutor) Describe() string {
	return fmt.Sprintf("%s via %s", t.model, t.baseURL)
}

// Ask sends the learner's question (plus their current query for context) and
// the recent conversation history, returning the tutor's guidance.
func (t *Tutor) Ask(ctx context.Context, question, currentSQL string, history []Message) (string, error) {
	if !t.Enabled() {
		return "", fmt.Errorf("tutor is not configured (no API key)")
	}

	reqBody := chatRequest{
		Model:       t.model,
		Messages:    buildMessages(history, buildUserMessage(question, currentSQL)),
		Temperature: 0.4,
		MaxTokens:   800,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := t.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling the AI provider: %w", err)
	}
	defer resp.Body.Close()

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decoding the AI response (is AI_BASE_URL correct?): %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg := out.Error.Message
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("AI API error (%d): %s", resp.StatusCode, msg)
	}

	text := out.text()
	if text == "" {
		return "", fmt.Errorf("the tutor had no response (the request may have been blocked)")
	}
	return text, nil
}

// buildMessages assembles the system prompt, the rolling history window, and the
// new question into the OpenAI-style messages array.
func buildMessages(history []Message, userText string) []chatMessage {
	if len(history) > historyWindow {
		history = history[len(history)-historyWindow:]
	}
	msgs := make([]chatMessage, 0, len(history)+2)
	msgs = append(msgs, chatMessage{Role: "system", Content: systemPrompt})
	for _, m := range history {
		role := "user"
		if m.Role == "tutor" {
			role = "assistant"
		}
		msgs = append(msgs, chatMessage{Role: role, Content: m.Text})
	}
	msgs = append(msgs, chatMessage{Role: "user", Content: userText})
	return msgs
}

// buildUserMessage combines the learner's question with their current query so
// the tutor can give targeted hints.
func buildUserMessage(question, currentSQL string) string {
	var b strings.Builder
	b.WriteString("Learner's question:\n")
	if strings.TrimSpace(question) == "" {
		b.WriteString("(no question typed — give a hint about their current query)")
	} else {
		b.WriteString(strings.TrimSpace(question))
	}
	if strings.TrimSpace(currentSQL) != "" {
		b.WriteString("\n\nTheir current query so far:\n```sql\n")
		b.WriteString(strings.TrimSpace(currentSQL))
		b.WriteString("\n```")
	}
	return b.String()
}

// --- OpenAI-compatible chat wire types ---------------------------------------

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// text pulls the assistant's reply out of the first choice.
func (r chatResponse) text() string {
	if len(r.Choices) == 0 {
		return ""
	}
	return strings.TrimSpace(r.Choices[0].Message.Content)
}
