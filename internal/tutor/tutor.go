// Package tutor implements the Phase 3 AI helper: a Socratic SQL coach backed
// by Google's Gemini API. It is deliberately built to *guide* learners toward
// the answer with hints, not to hand them a finished query.
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

// defaultModel is Gemini 2.5 Flash — fast and inexpensive, a good fit for
// short, interactive tutoring turns.
const defaultModel = "gemini-2.5-flash"

const apiBase = "https://generativelanguage.googleapis.com/v1beta/models/"

// systemPrompt shapes Gemini into a hint-giving tutor rather than an answer key.
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

// Tutor talks to the Gemini API on the learner's behalf.
type Tutor struct {
	apiKey string
	model  string
	http   *http.Client
}

// New returns a Tutor using the given API key. An empty key means the tutor is
// disabled; callers should check Enabled.
func New(apiKey string) *Tutor {
	return &Tutor{
		apiKey: strings.TrimSpace(apiKey),
		model:  defaultModel,
		http:   &http.Client{Timeout: 30 * time.Second},
	}
}

// Enabled reports whether an API key is configured.
func (t *Tutor) Enabled() bool { return t.apiKey != "" }

// Ask sends the learner's question along with their current query for context
// and returns the tutor's guidance.
func (t *Tutor) Ask(ctx context.Context, question, currentSQL string) (string, error) {
	if !t.Enabled() {
		return "", fmt.Errorf("tutor is not configured (no API key)")
	}

	userText := buildUserMessage(question, currentSQL)
	reqBody := geminiRequest{
		SystemInstruction: &content{Parts: []part{{Text: systemPrompt}}},
		Contents:          []content{{Role: "user", Parts: []part{{Text: userText}}}},
		GenerationConfig: &generationConfig{
			Temperature:     0.4,
			MaxOutputTokens: 800,
		},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := apiBase + t.model + ":generateContent"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", t.apiKey)

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling Gemini: %w", err)
	}
	defer resp.Body.Close()

	var out geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decoding Gemini response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg := out.Error.Message
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("Gemini API error: %s", msg)
	}

	text := out.text()
	if text == "" {
		return "", fmt.Errorf("the tutor had no response (the request may have been blocked)")
	}
	return text, nil
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

// --- Gemini JSON wire types --------------------------------------------------

type geminiRequest struct {
	SystemInstruction *content          `json:"systemInstruction,omitempty"`
	Contents          []content         `json:"contents"`
	GenerationConfig  *generationConfig `json:"generationConfig,omitempty"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// text pulls the concatenated text out of the first candidate.
func (r geminiResponse) text() string {
	if len(r.Candidates) == 0 {
		return ""
	}
	var b strings.Builder
	for _, p := range r.Candidates[0].Content.Parts {
		b.WriteString(p.Text)
	}
	return strings.TrimSpace(b.String())
}
