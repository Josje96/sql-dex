// Command sql-dex is a CLI for learning SQL by querying a rich Pokémon database.
//
// Phase 1: an interactive shell that runs read-only SQL against the pokédex and
// prints the results. Later phases will add a GUI (Phase 2) and an AI helper
// that guides you toward the answer instead of handing it over (Phase 3).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Josje96/sql-dex/internal/cli"
	"github.com/Josje96/sql-dex/internal/gui"
	"github.com/Josje96/sql-dex/internal/pokedb"
)

// defaultDBPath points at the pokédex bundled with the tutorial data.
const defaultDBPath = "PokemonSQLTutorial-master/pokedex.sqlite"

func main() {
	dbPath := flag.String("db", defaultDBPath, "path to the pokédex SQLite file")
	query := flag.String("q", "", "run a single query and exit (non-interactive)")
	guiMode := flag.Bool("gui", false, "launch the web GUI instead of the CLI")
	addr := flag.String("addr", "localhost:8080", "address for the web GUI to listen on")
	flag.Parse()

	if err := run(*dbPath, *query, *guiMode, *addr); err != nil {
		fmt.Fprintf(os.Stderr, "sql-dex: %v\n", err)
		os.Exit(1)
	}
}

func run(dbPath, oneShot string, guiMode bool, addr string) error {
	db, err := pokedb.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()

	// -gui serves the Phase 2 web interface.
	if guiMode {
		srv, err := gui.NewServer(ctx, db)
		if err != nil {
			return err
		}
		return srv.Listen(addr)
	}

	shell := cli.New(db, os.Stdin, os.Stdout)

	// -q runs a single query for scripting/piping, then exits.
	if oneShot != "" {
		return shell.RunOnce(ctx, oneShot)
	}
	return shell.Run(ctx)
}
