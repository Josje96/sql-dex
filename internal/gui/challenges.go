package gui

// Challenge is a task for the learner to solve in their own editor. It states
// what to find and gives a hint about which tables/approach to use, but never
// the query itself — that's the point.
//
// Solution is the reference query used to auto-check an answer. It is kept
// server-side only (json:"-") so it never reaches the browser.
type Challenge struct {
	ID         string `json:"id"`
	Difficulty string `json:"difficulty"` // Easy, Medium, Hard, Bonus
	Title      string `json:"title"`
	Task       string `json:"task"`
	Hint       string `json:"hint"`
	Solution   string `json:"-"`
}

// challenges is the starter set: 3 easy, 3 medium, 3 hard, and 1 bonus. Ordered
// easiest first so the modal reads as a difficulty ladder. Add more over time.
//
// Answers are checked by comparing result sets order-insensitively (both row
// order and column order), so any correct query passes — the tasks name the
// expected columns to keep each answer unambiguous.
var challenges = []Challenge{
	// --- Easy: single table, basic filter/sort ---
	{
		ID:         "kanto-roll-call",
		Difficulty: "Easy",
		Title:      "Kanto roll call",
		Task:       "List the identifier of every Generation 1 species.",
		Hint:       "Everything you need is in pokemon_species. Filter on generation_id.",
		Solution:   `SELECT identifier FROM pokemon_species WHERE generation_id = 1;`,
	},
	{
		ID:         "tallest-ten",
		Difficulty: "Easy",
		Title:      "The tall and the small",
		Task:       "Show the identifier and height of the 10 tallest Pokémon (default forms only).",
		Hint:       "Join pokemon to pokemon_species for the identifier. pokemon has height and an is_default flag. Think ORDER BY ... DESC and LIMIT.",
		Solution: `SELECT ps.identifier, p.height
		           FROM pokemon p
		           JOIN pokemon_species ps ON ps.id = p.species_id
		           WHERE p.is_default = 1
		           ORDER BY p.height DESC, ps.id
		           LIMIT 10;`,
	},
	{
		ID:         "english-names-20",
		Difficulty: "Easy",
		Title:      "Say my name",
		Task:       "Show the English name of the 20 lowest-numbered species (ids 1 through 20).",
		Hint:       "Names live in pokemon_species_names. English is local_language_id = 9. The species id is pokemon_species_id.",
		Solution: `SELECT name FROM pokemon_species_names
		           WHERE local_language_id = 9 AND pokemon_species_id <= 20;`,
	},

	// --- Medium: a join or a group ---
	{
		ID:         "type-census",
		Difficulty: "Medium",
		Title:      "Type census",
		Task:       "For each type, show the type name and how many Pokémon have it (most common first).",
		Hint:       "Join pokemon_types to type_names (language 9), then GROUP BY the type name and COUNT.",
		Solution: `SELECT tn.name, COUNT(*)
		           FROM pokemon_types pt
		           JOIN type_names tn ON tn.type_id = pt.type_id AND tn.local_language_id = 9
		           GROUP BY tn.name
		           ORDER BY COUNT(*) DESC;`,
	},
	{
		ID:         "fire-brigade",
		Difficulty: "Medium",
		Title:      "Fire brigade",
		Task:       "List the identifier of every Fire-type Pokémon (default forms).",
		Hint:       "pokemon_types links a pokemon to its type_id. Filter the type name to 'Fire' and join to pokemon_species for the identifier.",
		Solution: `SELECT DISTINCT ps.identifier
		           FROM pokemon p
		           JOIN pokemon_species ps ON ps.id = p.species_id
		           JOIN pokemon_types pt ON pt.pokemon_id = p.id
		           JOIN type_names tn ON tn.type_id = pt.type_id AND tn.local_language_id = 9
		           WHERE tn.name = 'Fire' AND p.is_default = 1;`,
	},
	{
		ID:         "weight-by-type",
		Difficulty: "Medium",
		Title:      "Weigh station",
		Task:       "For each type, show the type name and the average weight of Pokémon with that type, heaviest average first.",
		Hint:       "pokemon has weight, pokemon_types has the type. AVG(weight) with GROUP BY the type name, then ORDER BY the average.",
		Solution: `SELECT tn.name, AVG(p.weight)
		           FROM pokemon p
		           JOIN pokemon_types pt ON pt.pokemon_id = p.id
		           JOIN type_names tn ON tn.type_id = pt.type_id AND tn.local_language_id = 9
		           GROUP BY tn.name
		           ORDER BY AVG(p.weight) DESC;`,
	},

	// --- Hard: multiple joins, self-joins ---
	{
		ID:         "dual-types",
		Difficulty: "Hard",
		Title:      "Double trouble",
		Task:       "Find Pokémon with two types. Return the pokemon id and both type identifiers.",
		Hint:       "pokemon_types has a slot (1 and 2). Join it to itself on the same pokemon_id where one slot is 1 and the other is 2, then join to types for the identifiers.",
		Solution: `SELECT pt1.pokemon_id, t1.identifier, t2.identifier
		           FROM pokemon_types pt1
		           JOIN pokemon_types pt2 ON pt2.pokemon_id = pt1.pokemon_id AND pt2.slot = 2
		           JOIN types t1 ON t1.id = pt1.type_id
		           JOIN types t2 ON t2.id = pt2.type_id
		           WHERE pt1.slot = 1;`,
	},
	{
		ID:         "evolution-families",
		Difficulty: "Hard",
		Title:      "Family tree",
		Task:       "For every species that evolves from another, show its identifier and the identifier of the species it evolves from.",
		Hint:       "pokemon_species has evolves_from_species_id. Join pokemon_species to itself: one copy for the child, one for the parent.",
		Solution: `SELECT child.identifier, parent.identifier
		           FROM pokemon_species child
		           JOIN pokemon_species parent ON parent.id = child.evolves_from_species_id;`,
	},
	{
		ID:         "move-masters",
		Difficulty: "Hard",
		Title:      "Move masters",
		Task:       "List the 10 Pokémon (by id) that can learn the most distinct moves. Return the pokemon id and the move count.",
		Hint:       "pokemon_moves has one row per learnable move. COUNT(DISTINCT move_id) with GROUP BY pokemon_id, then ORDER BY that count DESC and LIMIT 10.",
		Solution: `SELECT pokemon_id, COUNT(DISTINCT move_id)
		           FROM pokemon_moves
		           GROUP BY pokemon_id
		           ORDER BY COUNT(DISTINCT move_id) DESC, pokemon_id
		           LIMIT 10;`,
	},

	// --- Bonus: the boss fight ---
	{
		ID:         "super-effective-water",
		Difficulty: "Bonus",
		Title:      "Super effective",
		Task:       "List the identifier of every type that is super effective against Water (damage_factor = 200).",
		Hint:       "type_efficacy holds damage_type_id, target_type_id and damage_factor. Join types twice (attacker and target) and filter the target to 'water'.",
		Solution: `SELECT dt.identifier
		           FROM type_efficacy te
		           JOIN types dt ON dt.id = te.damage_type_id
		           JOIN types tt ON tt.id = te.target_type_id
		           WHERE tt.identifier = 'water' AND te.damage_factor = 200;`,
	},
}

// challengeByID returns the challenge with the given id, or false.
func challengeByID(id string) (Challenge, bool) {
	for _, c := range challenges {
		if c.ID == id {
			return c, true
		}
	}
	return Challenge{}, false
}
