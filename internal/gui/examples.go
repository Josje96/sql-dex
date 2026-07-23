package gui

// Example is a labelled sample query shown in the Examples modal. Learners load
// it into the tutor pane to read and re-type, rather than copying it wholesale.
type Example struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SQL         string `json:"sql"`
}

// examples is a small, curated ladder of queries that ramps up in difficulty:
// simple SELECTs first, then filtering, then joins across the pokédex.
var examples = []Example{
	{
		Title:       "All Pokémon names",
		Description: "A plain SELECT — every species identifier in the pokédex.",
		SQL:         "SELECT identifier\nFROM pokemon_species\nLIMIT 20;",
	},
	{
		Title:       "First 151 (Gen 1)",
		Description: "Filter with WHERE to get only the original Kanto species.",
		SQL:         "SELECT identifier\nFROM pokemon_species\nWHERE generation_id = 1\nORDER BY id;",
	},
	{
		Title:       "Readable English names",
		Description: "Join pokemon_species_names to get display names (language 9 = English).",
		SQL: "SELECT sname.name\n" +
			"FROM pokemon_species AS species\n" +
			"JOIN pokemon_species_names AS sname\n" +
			"  ON sname.pokemon_species_id = species.id\n" +
			"WHERE sname.local_language_id = 9\n" +
			"LIMIT 20;",
	},
	{
		Title:       "Each Pokémon's type(s)",
		Description: "A multi-table join linking species → pokemon → types → type names.",
		SQL: "SELECT sname.name AS pokemon, tname.name AS type\n" +
			"FROM pokemon_species AS species\n" +
			"JOIN pokemon_species_names AS sname\n" +
			"  ON sname.pokemon_species_id = species.id AND sname.local_language_id = 9\n" +
			"JOIN pokemon AS p ON p.species_id = species.id\n" +
			"JOIN pokemon_types AS pt ON pt.pokemon_id = p.id\n" +
			"JOIN type_names AS tname\n" +
			"  ON tname.type_id = pt.type_id AND tname.local_language_id = 9\n" +
			"WHERE species.generation_id = 1\n" +
			"LIMIT 20;",
	},
	{
		Title:       "Heaviest Pokémon",
		Description: "ORDER BY with DESC — weight is stored in hectograms.",
		SQL: "SELECT sname.name, p.weight\n" +
			"FROM pokemon AS p\n" +
			"JOIN pokemon_species_names AS sname\n" +
			"  ON sname.pokemon_species_id = p.species_id AND sname.local_language_id = 9\n" +
			"ORDER BY p.weight DESC\n" +
			"LIMIT 10;",
	},
	{
		Title:       "Count Pokémon per type",
		Description: "GROUP BY with COUNT — how many species share each type.",
		SQL: "SELECT tname.name AS type, COUNT(*) AS how_many\n" +
			"FROM pokemon_types AS pt\n" +
			"JOIN type_names AS tname\n" +
			"  ON tname.type_id = pt.type_id AND tname.local_language_id = 9\n" +
			"GROUP BY tname.name\n" +
			"ORDER BY how_many DESC;",
	},
}
