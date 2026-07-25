# SQL-Dex

A little SQL game. You write SQL against a full Pokemon database and learn SQL along the way.

## Credits

Built on the shoulders of these projects:

- [veekun/pokedex](https://github.com/veekun/pokedex) — the Pokemon SQLite database and example queries. Thanks to Veekun for creating it.
- [CodeMirror](https://codemirror.net) — the in-browser SQL editor (MIT).
- [sql.js](https://github.com/sql-js/sql.js) — SQLite compiled to WebAssembly, runs the queries in your browser (MIT).

Full license texts are in [THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md).

Pokemon and Pokemon character names are trademarks of Nintendo, Game Freak, and The Pokemon Company. This is an unofficial, non-commercial educational tool, not affiliated with or endorsed by them.

## What it is

SQL-Dex hands you a 170-table Pokemon database (the veekun pokedex) and lets you query it. It comes in two flavors:

- a CLI where you type SQL and get results back
- a web app with two editors (yours and a tutor pane), a searchable SQL 101 guide, an examples list, a notepad for looking things up, and pokedex sprites that show up next to your results

There is also an optional AI tutor. It does not hand you the answer. It nudges you toward it with hints and small examples, so you learn instead of copy pasting.

## Why it's here

Learning SQL from toy tables is boring. This uses data you actually (might) know about and gives you a safe place to poke around. The database is opened read only, so you cannot break anything.

## Running it on a fresh machine

You only need Go (1.25 or newer). Get it from https://go.dev/dl. The database is already in the repo, so there is nothing else to download.

Clone it:

```bash
git clone https://github.com/Josje96/sql-dex.git
cd sql-dex
```

Run the CLI:

```bash
go run .
```

Or run the web app and open http://localhost:8080:

```bash
go run . -gui
```

That is it. Everything works without any keys.

## Turning on the AI tutor (optional)

The tutor talks to any OpenAI compatible provider (OpenAI, DeepSeek, SiliconFlow, OpenRouter, or Google Gemini). To switch it on, copy the example env file and drop in a key:

```bash
cp .env.example .env
```

Open `.env`, pick a provider, and fill in the three values (`AI_API_KEY`, `AI_BASE_URL`, `AI_MODEL`). Every provider's base URL, model name, and where to get a key is listed right there in the file. Restart the app and the tutor is live. Your `.env` is git ignored, so your key stays local.

## Branches

- `main` is the full app (CLI plus web app plus tutor)
- `CLI` is a snapshot if you only want the terminal experience the cli doesnt include the ai tutor
