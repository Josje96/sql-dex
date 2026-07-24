package gui

// GuideTopic is one entry in the searchable SQL 101 reference. Each topic pairs
// a plain-English explanation with a small, pokédex-flavored example the learner
// can drop into the tutor pane and try.
type GuideTopic struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	SQL      string `json:"sql"`
	// Keywords are extra search terms that aren't in the visible text, so a
	// search for "fields" still finds the "columns" topic.
	Keywords string `json:"keywords"`
}

// guide is a stripped-down SQL primer: just enough to find tables, inspect
// columns, filter, sort, join, and summarize. Ordered from "explore" to
// "combine" so reading top-to-bottom is itself a mini course.
var guide = []GuideTopic{
	// --- Explore the database ---
	{
		Category: "Explore",
		Title:    "List all tables",
		Body:     "Every table lives in the built-in sqlite_master catalog. This is how you discover what data exists.",
		SQL:      "SELECT name\nFROM sqlite_master\nWHERE type = 'table'\nORDER BY name;",
		Keywords: "tables show list names discover catalog what exists",
	},
	{
		Category: "Explore",
		Title:    "See a table's columns",
		Body:     "PRAGMA table_info lists each column's name and type, so you know what you can select.",
		SQL:      "PRAGMA table_info(pokemon_species);",
		Keywords: "columns fields structure describe schema pragma types",
	},
	{
		Category: "Explore",
		Title:    "Peek at some rows",
		Body:     "The quickest way to see what a table holds: select everything (*) and cap it with LIMIT.",
		SQL:      "SELECT *\nFROM pokemon\nLIMIT 5;",
		Keywords: "sample preview star all rows peek look",
	},

	// --- Choosing what to show ---
	{
		Category: "Select",
		Title:    "Pick specific columns",
		Body:     "List the columns you want after SELECT instead of *. Cleaner and faster.",
		SQL:      "SELECT identifier, generation_id\nFROM pokemon_species\nLIMIT 10;",
		Keywords: "choose columns projection select specific",
	},
	{
		Category: "Select",
		Title:    "Limit how many rows",
		Body:     "LIMIT caps the number of rows returned — handy while you're experimenting.",
		SQL:      "SELECT identifier\nFROM pokemon_species\nLIMIT 20;",
		Keywords: "limit top first cap rows how many",
	},
	{
		Category: "Select",
		Title:    "Rename with AS (aliases)",
		Body:     "AS gives a column or table a friendlier name in your results.",
		SQL:      "SELECT identifier AS name\nFROM pokemon_species\nLIMIT 5;",
		Keywords: "alias as rename label heading",
	},

	// --- Filtering ---
	{
		Category: "Filter",
		Title:    "Keep matching rows (WHERE)",
		Body:     "WHERE keeps only the rows that satisfy a condition.",
		SQL:      "SELECT identifier\nFROM pokemon_species\nWHERE generation_id = 1;",
		Keywords: "where filter condition equals only matching",
	},
	{
		Category: "Filter",
		Title:    "Combine conditions (AND / OR)",
		Body:     "Use AND when every condition must be true, OR when any of them can be.",
		SQL:      "SELECT identifier\nFROM pokemon_species\nWHERE generation_id = 1\n  AND capture_rate > 100;",
		Keywords: "and or combine multiple conditions boolean both either",
	},
	{
		Category: "Filter",
		Title:    "Text patterns (LIKE)",
		Body:     "LIKE matches text patterns. % stands for \"any characters\" — 'char%' means starts with 'char'.",
		SQL:      "SELECT identifier\nFROM pokemon_species\nWHERE identifier LIKE 'char%';",
		Keywords: "like pattern wildcard search text contains starts ends %",
	},
	{
		Category: "Filter",
		Title:    "Ranges and lists (BETWEEN, IN)",
		Body:     "BETWEEN matches a range; IN matches any value from a list.",
		SQL:      "SELECT identifier\nFROM pokemon_species\nWHERE id IN (1, 4, 7, 25);",
		Keywords: "in between range list multiple values set",
	},

	// --- Sorting & shaping ---
	{
		Category: "Sort",
		Title:    "Sort results (ORDER BY)",
		Body:     "ORDER BY sorts the rows. Add DESC to go high-to-low instead of low-to-high.",
		SQL:      "SELECT identifier, capture_rate\nFROM pokemon_species\nORDER BY capture_rate DESC\nLIMIT 10;",
		Keywords: "order by sort ascending descending asc desc rank highest lowest largest",
	},
	{
		Category: "Sort",
		Title:    "Unique values (DISTINCT)",
		Body:     "DISTINCT collapses repeated values so each one appears once.",
		SQL:      "SELECT DISTINCT generation_id\nFROM pokemon_species;",
		Keywords: "distinct unique duplicates dedupe different",
	},

	// --- Combining & summarizing ---
	{
		Category: "Combine",
		Title:    "Join two tables (JOIN)",
		Body:     "JOIN links rows across tables using a column they share (a key). ON says how they match.",
		SQL: "SELECT ps.identifier, sname.name\n" +
			"FROM pokemon_species AS ps\n" +
			"JOIN pokemon_species_names AS sname\n" +
			"  ON sname.pokemon_species_id = ps.id\n" +
			"WHERE sname.local_language_id = 9;",
		Keywords: "join combine link tables relationship on key foreign inner",
	},
	{
		Category: "Combine",
		Title:    "Count rows (COUNT)",
		Body:     "COUNT(*) counts how many rows match — often paired with WHERE.",
		SQL:      "SELECT COUNT(*)\nFROM pokemon_species\nWHERE generation_id = 1;",
		Keywords: "count aggregate how many total number rows",
	},
	{
		Category: "Combine",
		Title:    "Group and count (GROUP BY)",
		Body:     "GROUP BY buckets rows by a column so you can count or sum within each group.",
		SQL: "SELECT generation_id, COUNT(*) AS species\n" +
			"FROM pokemon_species\n" +
			"GROUP BY generation_id\n" +
			"ORDER BY generation_id;",
		Keywords: "group by aggregate count per bucket summarize sum average totals",
	},

	// --- Pokédex-specific tip ---
	{
		Category: "Pokédex tips",
		Title:    "Get English names",
		Body:     "Name tables hold many languages. Filter with local_language_id = 9 to get English.",
		SQL:      "SELECT name\nFROM pokemon_species_names\nWHERE local_language_id = 9\nLIMIT 10;",
		Keywords: "english names language local_language_id 9 display readable",
	},
}
