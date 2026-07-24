package gui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Josje96/sql-dex/internal/pokedb"
)

// rowKeySep separates cell values in a canonical row key. It's a control char
// so it won't collide with real data.
const rowKeySep = "\x1f"

// resultsMatch reports whether two result sets contain the same rows, ignoring
// both row order and column order. This lets any correct query pass regardless
// of how the learner ordered their SELECT list or their rows.
func resultsMatch(a, b *pokedb.Result) bool {
	ca, cb := canonRows(a), canonRows(b)
	if len(ca) != len(cb) {
		return false
	}
	for k, n := range ca {
		if cb[k] != n {
			return false
		}
	}
	return true
}

// canonRows builds a multiset of canonical row keys. Each row's cell values are
// stringified and sorted, so column order doesn't matter.
func canonRows(res *pokedb.Result) map[string]int {
	m := make(map[string]int, len(res.Rows))
	for _, row := range res.Rows {
		cells := make([]string, len(row))
		for i, v := range row {
			cells[i] = canonCell(v)
		}
		sort.Strings(cells)
		m[strings.Join(cells, rowKeySep)]++
	}
	return m
}

// canonCell renders a scanned SQL value to a stable string for comparison.
func canonCell(v any) string {
	switch x := v.(type) {
	case nil:
		return "\x00NULL"
	case []byte:
		return string(x)
	default:
		return fmt.Sprintf("%v", x)
	}
}
