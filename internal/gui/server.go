// Package gui implements the Phase 2 web interface for SQL-Dex: a local HTTP
// server that serves a browser frontend (side-by-side editors, an examples
// modal, and pokédex sprite cards) backed by the read-only pokedb engine.
package gui

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/Josje96/sql-dex/internal/pokedb"
	"github.com/Josje96/sql-dex/internal/tutor"
)

//go:embed static
var staticFS embed.FS

// Server holds the dependencies shared across HTTP handlers.
type Server struct {
	db      *pokedb.DB
	sprites *spriteIndex
	tutor   *tutor.Tutor
}

// NewServer builds the server and pre-loads the sprite index from the database.
// coach may be a disabled tutor (no API key); the /api/tutor route reports that.
func NewServer(ctx context.Context, db *pokedb.DB, coach *tutor.Tutor) (*Server, error) {
	idx, err := buildSpriteIndex(ctx, db)
	if err != nil {
		return nil, err
	}
	return &Server{db: db, sprites: idx, tutor: coach}, nil
}

// Handler returns the router with API routes and the embedded static frontend.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/query", s.handleQuery)
	mux.HandleFunc("/api/tables", s.handleTables)
	mux.HandleFunc("/api/schema", s.handleSchema)
	mux.HandleFunc("/api/examples", s.handleExamples)
	mux.HandleFunc("/api/tutor", s.handleTutor)

	// Serve the frontend from the embedded static/ directory at the web root.
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err) // embed path is a compile-time constant; this can't fail
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))
	return mux
}

// Listen starts the HTTP server on addr and blocks until it stops.
func (s *Server) Listen(addr string) error {
	log.Printf("SQL-Dex GUI running at http://%s  (Ctrl-C to stop)", addr)
	return http.ListenAndServe(addr, s.Handler())
}
