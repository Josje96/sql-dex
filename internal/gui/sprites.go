package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/Josje96/sql-dex/internal/pokedb"
)

// spriteBaseURL is the PokeAPI sprite set on GitHub. A species' national-dex
// number is its filename, e.g. .../pokemon/1.png for Bulbasaur.
const spriteBaseURL = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/"

// Sprite is a Pokémon to render as a pokédex-style card in the results view.
type Sprite struct {
	Dex  int    `json:"dex"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// spriteIndex maps lowercased species identifiers and English names to their
// national-dex number, so we can spot Pokémon anywhere in a result set.
type spriteIndex struct {
	byKey map[string]int // lookup key -> dex number
	names map[int]string // dex number -> display name
}

// buildSpriteIndex loads every species' identifier and English name once, up
// front, so sprite lookups during queries are just map hits.
func buildSpriteIndex(ctx context.Context, db *pokedb.DB) (*spriteIndex, error) {
	const q = `SELECT ps.id, ps.identifier, sname.name
	           FROM pokemon_species AS ps
	           LEFT JOIN pokemon_species_names AS sname
	             ON sname.pokemon_species_id = ps.id AND sname.local_language_id = 9`
	res, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	idx := &spriteIndex{byKey: make(map[string]int), names: make(map[int]string)}
	for _, row := range res.Rows {
		dex, ok := toInt(row[0])
		if !ok {
			continue
		}
		display := ""
		if name, ok := row[2].(string); ok && name != "" {
			display = name
			idx.byKey[normalize(name)] = dex
		}
		if ident, ok := row[1].(string); ok && ident != "" {
			idx.byKey[normalize(ident)] = dex
			if display == "" {
				display = ident
			}
		}
		idx.names[dex] = display
	}
	return idx, nil
}

// find scans a query result for any cell that names a known Pokémon and returns
// the matching sprites, de-duplicated and capped so a huge result set can't
// flood the UI with images.
func (idx *spriteIndex) find(res *pokedb.Result) []Sprite {
	const maxSprites = 60
	var sprites []Sprite
	seen := make(map[int]bool)
	for _, row := range res.Rows {
		for _, v := range row {
			s, ok := v.(string)
			if !ok {
				continue
			}
			dex, ok := idx.byKey[normalize(s)]
			if !ok || seen[dex] {
				continue
			}
			seen[dex] = true
			sprites = append(sprites, Sprite{
				Dex:  dex,
				Name: idx.names[dex],
				URL:  fmt.Sprintf("%s%d.png", spriteBaseURL, dex),
			})
			if len(sprites) >= maxSprites {
				return sprites
			}
		}
	}
	return sprites
}

// normalize lowercases and trims a value so "Bulbasaur", "bulbasaur", and
// " bulbasaur " all match the same species.
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// toInt coerces the driver's numeric types (typically int64) to an int.
func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int64:
		return int(n), true
	case int:
		return n, true
	case float64:
		return int(n), true
	default:
		return 0, false
	}
}
