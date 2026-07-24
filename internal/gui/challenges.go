package gui

// Challenge is a task for the learner to solve in their own editor. It states
// what to find and gives a hint about which tables/approach to use, but never
// the query itself — that's the point.
type Challenge struct {
	Difficulty string `json:"difficulty"` // Easy, Medium, Hard, Bonus
	Title      string `json:"title"`
	Task       string `json:"task"`
	Hint       string `json:"hint"`
}

// challenges is the starter set: 3 easy, 3 medium, 3 hard, and 1 bonus. Ordered
// easiest first so the modal reads as a difficulty ladder. Add more over time.
var challenges = []Challenge{
	// --- Easy: single table, basic filter/sort ---
	{
		Difficulty: "Easy",
		Title:      "Kanto roll call",
		Task:       "List the identifiers of every Generation 1 species.",
		Hint:       "Everything you need is in pokemon_species. Filter on generation_id.",
	},
	{
		Difficulty: "Easy",
		Title:      "The tall and the small",
		Task:       "Show the 10 tallest Pokémon with their height.",
		Hint:       "The pokemon table has a height column. Think ORDER BY ... DESC and LIMIT.",
	},
	{
		Difficulty: "Easy",
		Title:      "Say my name",
		Task:       "Get the English display names of the first 20 species (by id).",
		Hint:       "Names live in pokemon_species_names. English is local_language_id = 9.",
	},

	// --- Medium: a join or a group ---
	{
		Difficulty: "Medium",
		Title:      "Type census",
		Task:       "Count how many Pokémon have each type, most common first.",
		Hint:       "Join pokemon_types to type_names (language 9), then GROUP BY the type name.",
	},
	{
		Difficulty: "Medium",
		Title:      "Fire brigade",
		Task:       "List the names of every Fire-type Pokémon.",
		Hint:       "pokemon_types links a pokemon to its type_id. Filter the type name to 'Fire', and join to species names for something readable.",
	},
	{
		Difficulty: "Medium",
		Title:      "Weigh station",
		Task:       "Find the average weight per type, heaviest type first.",
		Hint:       "pokemon has weight, pokemon_types has the type. AVG(weight) with GROUP BY, then ORDER BY the average.",
	},

	// --- Hard: multiple joins, self-joins, HAVING ---
	{
		Difficulty: "Hard",
		Title:      "Double trouble",
		Task:       "Find Pokémon that have two types, and show the name plus both types.",
		Hint:       "pokemon_types has a slot (1 and 2). Join it to itself on the same pokemon_id where one slot is 1 and the other is 2. Or GROUP BY with HAVING COUNT(*) = 2.",
	},
	{
		Difficulty: "Hard",
		Title:      "Family tree",
		Task:       "For every species that evolves from another, show the species name and what it evolves from.",
		Hint:       "pokemon_species has evolves_from_species_id. Join pokemon_species to itself, and bring in names for both sides.",
	},
	{
		Difficulty: "Hard",
		Title:      "Move masters",
		Task:       "List the 10 Pokémon that can learn the most distinct moves.",
		Hint:       "pokemon_moves has one row per learnable move. COUNT(DISTINCT move_id) with GROUP BY pokemon_id, then ORDER BY that count DESC and LIMIT.",
	},

	// --- Bonus: the boss fight ---
	{
		Difficulty: "Bonus",
		Title:      "Super effective",
		Task:       "List every type that is super effective against Water.",
		Hint:       "type_efficacy holds damage_type_id, target_type_id and damage_factor. Super effective is damage_factor = 200. Join type_names twice (once per type) and filter the target to 'Water'.",
	},
}
