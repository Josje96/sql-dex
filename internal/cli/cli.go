// Package cli implements the interactive Phase 1 shell for SQL-Dex: it reads
// SQL (or dot-commands) from the user, runs it against the pokédex, and prints
// the results as an aligned table.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Josje96/sql-dex/internal/pokedb"
)

// Shell drives a read-eval-print loop over a pokédex database.
type Shell struct {
	db  *pokedb.DB
	in  *bufio.Scanner
	out io.Writer
}

// New creates a Shell reading commands from in and writing output to out.
func New(db *pokedb.DB, in io.Reader, out io.Writer) *Shell {
	sc := bufio.NewScanner(in)
	// Pokédex CREATE statements and long rows can exceed the default 64K token.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return &Shell{db: db, in: sc, out: out}
}

// Run starts the loop and returns when the user quits or input ends (EOF).
func (s *Shell) Run(ctx context.Context) error {
	s.printBanner()
	// buf accumulates lines until a statement ends with ';' so users can spread
	// a query across several lines, just like the sqlite3 shell.
	var buf strings.Builder
	s.prompt(false)
	for s.in.Scan() {
		line := s.in.Text()
		trimmed := strings.TrimSpace(line)

		// Dot-commands are only recognised at the start of a fresh statement.
		if buf.Len() == 0 && strings.HasPrefix(trimmed, ".") {
			if quit := s.runDotCommand(ctx, trimmed); quit {
				return nil
			}
			s.prompt(false)
			continue
		}

		buf.WriteString(line)
		buf.WriteByte('\n')

		if strings.HasSuffix(trimmed, ";") {
			stmt := strings.TrimSpace(buf.String())
			buf.Reset()
			if stmt != ";" {
				s.runQuery(ctx, strings.TrimSuffix(stmt, ";"))
			}
			s.prompt(false)
			continue
		}
		s.prompt(true) // continuation prompt
	}
	if err := s.in.Err(); err != nil {
		return err
	}
	fmt.Fprintln(s.out) // clean newline after EOF (Ctrl-D / Ctrl-Z)
	return nil
}

// RunOnce executes a single query and prints its results, without the
// interactive banner or prompt. Used by the -q flag for scripting.
func (s *Shell) RunOnce(ctx context.Context, query string) error {
	s.runQuery(ctx, strings.TrimSuffix(strings.TrimSpace(query), ";"))
	return nil
}

func (s *Shell) prompt(continuation bool) {
	if continuation {
		fmt.Fprint(s.out, "  ...> ")
	} else {
		fmt.Fprint(s.out, "sql-dex> ")
	}
}

func (s *Shell) printBanner() {
	fmt.Fprintln(s.out, "SQL-Dex — learn SQL by exploring the Pokémon world.")
	fmt.Fprintln(s.out, "Type SQL ending in ';', or .help for commands. The database is read-only.")
	fmt.Fprintln(s.out)
}

// runDotCommand handles the meta-commands. It returns true when the user asked
// to quit.
func (s *Shell) runDotCommand(ctx context.Context, line string) (quit bool) {
	fields := strings.Fields(line)
	cmd := fields[0]
	args := fields[1:]

	switch cmd {
	case ".quit", ".exit":
		return true
	case ".help":
		s.printHelp()
	case ".tables":
		s.printTables(ctx)
	case ".schema":
		if len(args) == 0 {
			fmt.Fprintln(s.out, "usage: .schema <table>")
			return false
		}
		schema, err := s.db.Schema(ctx, args[0])
		if err != nil {
			fmt.Fprintf(s.out, "Error: %v\n", err)
			return false
		}
		fmt.Fprintln(s.out, schema)
	default:
		fmt.Fprintf(s.out, "Unknown command %q. Type .help for the list.\n", cmd)
	}
	return false
}

func (s *Shell) printHelp() {
	fmt.Fprintln(s.out, `Commands:
  .tables          list all tables in the pokédex
  .schema <table>  show the CREATE statement for a table
  .help            show this help
  .quit / .exit    leave SQL-Dex

Everything else is treated as SQL. End a statement with ';' to run it.
Try:  SELECT identifier FROM pokemon_species LIMIT 5;`)
}

func (s *Shell) printTables(ctx context.Context) {
	tables, err := s.db.Tables(ctx)
	if err != nil {
		fmt.Fprintf(s.out, "Error: %v\n", err)
		return
	}
	for _, t := range tables {
		fmt.Fprintln(s.out, t)
	}
	fmt.Fprintf(s.out, "\n%d tables\n", len(tables))
}

func (s *Shell) runQuery(ctx context.Context, query string) {
	start := time.Now()
	res, err := s.db.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(s.out, "Error: %v\n", err)
		return
	}
	s.printTable(res)
	fmt.Fprintf(s.out, "%s in %s\n\n", rowsLabel(len(res.Rows)), took(start))
}

// printTable renders a Result as an aligned, bordered table.
func (s *Shell) printTable(res *pokedb.Result) {
	if len(res.Columns) == 0 {
		return
	}
	// Compute the display width of each column from its header and every cell.
	widths := make([]int, len(res.Columns))
	for i, c := range res.Columns {
		widths[i] = utf8.RuneCountInString(c)
	}
	cells := make([][]string, len(res.Rows))
	for r, row := range res.Rows {
		cells[r] = make([]string, len(row))
		for i, v := range row {
			text := format(v)
			cells[r][i] = text
			if w := utf8.RuneCountInString(text); w > widths[i] {
				widths[i] = w
			}
		}
	}

	s.printSeparator(widths)
	s.printRow(res.Columns, widths)
	s.printSeparator(widths)
	for _, row := range cells {
		s.printRow(row, widths)
	}
	s.printSeparator(widths)
}

func (s *Shell) printSeparator(widths []int) {
	var b strings.Builder
	b.WriteByte('+')
	for _, w := range widths {
		b.WriteString(strings.Repeat("-", w+2))
		b.WriteByte('+')
	}
	fmt.Fprintln(s.out, b.String())
}

func (s *Shell) printRow(cells []string, widths []int) {
	var b strings.Builder
	b.WriteByte('|')
	for i, w := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		pad := w - utf8.RuneCountInString(cell)
		b.WriteByte(' ')
		b.WriteString(cell)
		b.WriteString(strings.Repeat(" ", pad))
		b.WriteString(" |")
	}
	fmt.Fprintln(s.out, b.String())
}

// format turns a scanned SQL value into its display string.
func format(v any) string {
	switch val := v.(type) {
	case nil:
		return "NULL"
	case []byte:
		return string(val)
	case string:
		return val
	case bool:
		if val {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func rowsLabel(n int) string {
	if n == 1 {
		return "1 row"
	}
	return fmt.Sprintf("%d rows", n)
}

func took(start time.Time) time.Duration {
	return time.Since(start).Round(time.Millisecond)
}
