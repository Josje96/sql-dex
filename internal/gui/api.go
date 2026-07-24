package gui

import (
	"encoding/json"
	"net/http"
)

// queryRequest is the JSON body posted to /api/query.
type queryRequest struct {
	SQL string `json:"sql"`
}

// queryResponse is what /api/query returns: either an error message, or the
// columns/rows plus any Pokémon sprites detected in the results.
type queryResponse struct {
	Columns []string `json:"columns,omitempty"`
	Rows    [][]any  `json:"rows,omitempty"`
	Sprites []Sprite `json:"sprites,omitempty"`
	Error   string   `json:"error,omitempty"`
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req queryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, queryResponse{Error: "invalid request body"})
		return
	}

	res, err := s.db.Query(r.Context(), req.SQL)
	if err != nil {
		// SQL errors (syntax, read-only writes) are expected while learning, so
		// return them as data with 200 rather than an HTTP error.
		writeJSON(w, http.StatusOK, queryResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, queryResponse{
		Columns: res.Columns,
		Rows:    res.Rows,
		Sprites: s.sprites.find(res),
	})
}

func (s *Server) handleTables(w http.ResponseWriter, r *http.Request) {
	tables, err := s.db.Tables(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tables": tables})
}

func (s *Server) handleSchema(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	if table == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing ?table="})
		return
	}
	schema, err := s.db.Schema(r.Context(), table)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"schema": schema})
}

func (s *Server) handleExamples(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"examples": examples})
}

func (s *Server) handleGuide(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"guide": guide})
}

// tutorRequest is the JSON body posted to /api/tutor.
type tutorRequest struct {
	Question string `json:"question"`
	SQL      string `json:"sql"`
}

// tutorResponse carries the tutor's guidance, or an error / disabled flag.
type tutorResponse struct {
	Hint     string `json:"hint,omitempty"`
	Error    string `json:"error,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

func (s *Server) handleTutor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.tutor == nil || !s.tutor.Enabled() {
		writeJSON(w, http.StatusOK, tutorResponse{
			Disabled: true,
			Error:    "The tutor is off. Add GOOGLE=<your Gemini API key> to a .env file and restart.",
		})
		return
	}
	var req tutorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, tutorResponse{Error: "invalid request body"})
		return
	}

	hint, err := s.tutor.Ask(r.Context(), req.Question, req.SQL)
	if err != nil {
		writeJSON(w, http.StatusOK, tutorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, tutorResponse{Hint: hint})
}

// writeJSON encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
