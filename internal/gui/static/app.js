"use strict";

// SQL-Dex frontend: wires the two CodeMirror editors, the Run button, the
// examples modal, and the pokédex sprite strip to the Go JSON API.

const STARTER_QUERY = `-- Welcome to SQL-Dex! Edit this and hit Run (Ctrl+Enter).
SELECT sname.name AS pokemon, tname.name AS type
FROM pokemon_species AS species
JOIN pokemon_species_names AS sname
  ON sname.pokemon_species_id = species.id AND sname.local_language_id = 9
JOIN pokemon AS p ON p.species_id = species.id
JOIN pokemon_types AS pt ON pt.pokemon_id = p.id
JOIN type_names AS tname
  ON tname.type_id = pt.type_id AND tname.local_language_id = 9
WHERE species.generation_id = 1
LIMIT 12;`;

// --- Editors -----------------------------------------------------------------

function makeEditor(id, opts) {
  return CodeMirror.fromTextArea(document.getElementById(id), {
    mode: "text/x-sql",
    theme: "material-darker",
    lineNumbers: true,
    lineWrapping: true,
    ...opts,
  });
}

const userEditor = makeEditor("userEditor", {});
const tutorEditor = makeEditor("tutorEditor", { readOnly: true });

userEditor.setValue(STARTER_QUERY);
tutorEditor.setValue("-- Open Examples to load a query here,\n-- read it, then retype it in your own editor.");

// --- Query execution ---------------------------------------------------------

const runBtn = document.getElementById("runBtn");
const resultMeta = document.getElementById("resultMeta");
const resultArea = document.getElementById("resultArea");
const spriteStrip = document.getElementById("spriteStrip");

async function runQuery() {
  const sql = userEditor.getValue().trim();
  if (!sql) return;

  runBtn.disabled = true;
  setMeta("Running…", "");
  const started = performance.now();

  try {
    const resp = await fetch("/api/query", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ sql }),
    });
    const data = await resp.json();
    const ms = Math.round(performance.now() - started);

    if (data.error) {
      showError(data.error);
      return;
    }
    renderTable(data.columns || [], data.rows || []);
    renderSprites(data.sprites || []);
    const n = (data.rows || []).length;
    setMeta(`${n} row${n === 1 ? "" : "s"} · ${ms} ms`, "ok");
  } catch (err) {
    showError("Could not reach the server: " + err.message);
  } finally {
    runBtn.disabled = false;
  }
}

function setMeta(text, kind) {
  resultMeta.textContent = text;
  resultMeta.className = "result-meta" + (kind ? " " + kind : "");
}

function showError(message) {
  setMeta("Error: " + message, "error");
  resultArea.innerHTML = "";
  spriteStrip.hidden = true;
}

function renderTable(columns, rows) {
  if (columns.length === 0) {
    resultArea.innerHTML = '<p class="placeholder">Query ran, but returned no columns.</p>';
    return;
  }
  const table = document.createElement("table");
  table.className = "grid";

  const thead = document.createElement("thead");
  const hr = document.createElement("tr");
  for (const c of columns) {
    const th = document.createElement("th");
    th.textContent = c;
    hr.appendChild(th);
  }
  thead.appendChild(hr);
  table.appendChild(thead);

  const tbody = document.createElement("tbody");
  for (const row of rows) {
    const tr = document.createElement("tr");
    for (const cell of row) {
      const td = document.createElement("td");
      if (cell === null || cell === undefined) {
        td.textContent = "NULL";
        td.className = "null";
      } else {
        td.textContent = String(cell);
      }
      tr.appendChild(td);
    }
    tbody.appendChild(tr);
  }
  table.appendChild(tbody);

  resultArea.innerHTML = "";
  resultArea.appendChild(table);
}

// renderSprites shows a pokédex-style strip of any Pokémon found in the results.
function renderSprites(sprites) {
  if (!sprites.length) {
    spriteStrip.hidden = true;
    spriteStrip.innerHTML = "";
    return;
  }
  spriteStrip.innerHTML = "";
  for (const s of sprites) {
    const card = document.createElement("div");
    card.className = "sprite-card";

    const img = document.createElement("img");
    img.src = s.url;
    img.alt = s.name;
    img.loading = "lazy";
    // Hide the whole card if the sprite image can't be fetched.
    img.onerror = () => card.remove();

    const dex = document.createElement("div");
    dex.className = "dex";
    dex.textContent = "#" + String(s.dex).padStart(3, "0");

    const nm = document.createElement("div");
    nm.className = "nm";
    nm.textContent = s.name;

    card.append(img, dex, nm);
    spriteStrip.appendChild(card);
  }
  spriteStrip.hidden = false;
}

// --- Examples modal ----------------------------------------------------------

const modal = document.getElementById("examplesModal");
const examplesList = document.getElementById("examplesList");
const examplesBtn = document.getElementById("examplesBtn");
let examplesLoaded = false;

async function openExamples() {
  if (!examplesLoaded) {
    try {
      const resp = await fetch("/api/examples");
      const data = await resp.json();
      buildExamples(data.examples || []);
      examplesLoaded = true;
    } catch (err) {
      examplesList.innerHTML =
        '<li class="placeholder">Could not load examples: ' + err.message + "</li>";
    }
  }
  modal.hidden = false;
}

function closeExamples() {
  modal.hidden = true;
}

function buildExamples(list) {
  examplesList.innerHTML = "";
  for (const ex of list) {
    const li = document.createElement("li");
    li.className = "example";
    li.title = "Load into the Tutor pane";

    const h = document.createElement("h3");
    h.textContent = ex.title;
    const p = document.createElement("p");
    p.textContent = ex.description;
    const pre = document.createElement("pre");
    pre.textContent = ex.sql;

    li.append(h, p, pre);
    // Muscle memory: examples go into the read-only tutor pane, not the user's.
    li.addEventListener("click", () => {
      tutorEditor.setValue(ex.sql);
      closeExamples();
    });
    examplesList.appendChild(li);
  }
}

// --- SQL 101 guide modal -----------------------------------------------------

const guideModal = document.getElementById("guideModal");
const guideBtn = document.getElementById("guideBtn");
const guideSearch = document.getElementById("guideSearch");
const guideResults = document.getElementById("guideResults");
let guideData = null;

async function openGuide() {
  if (guideData === null) {
    try {
      const resp = await fetch("/api/guide");
      const data = await resp.json();
      guideData = data.guide || [];
    } catch (err) {
      guideData = [];
      guideResults.innerHTML =
        '<p class="guide-empty">Could not load the guide: ' + escapeHtml(err.message) + "</p>";
    }
  }
  renderGuide(guideSearch.value);
  guideModal.hidden = false;
  guideSearch.focus();
}

function closeGuide() {
  guideModal.hidden = true;
}

// renderGuide filters topics by the search query and lays them out grouped by
// category, with the matched text highlighted.
function renderGuide(query) {
  if (!guideData) return;
  const q = query.trim().toLowerCase();
  const matches = guideData.filter((t) => {
    if (!q) return true;
    const hay = (t.title + " " + t.body + " " + t.sql + " " + t.keywords + " " + t.category).toLowerCase();
    return hay.includes(q);
  });

  guideResults.innerHTML = "";
  if (matches.length === 0) {
    const p = document.createElement("p");
    p.className = "guide-empty";
    p.textContent = 'No matches for "' + query.trim() + '". Try "sort", "join", or "columns".';
    guideResults.appendChild(p);
    return;
  }

  let lastCat = null;
  for (const t of matches) {
    if (t.category !== lastCat) {
      const h = document.createElement("div");
      h.className = "guide-cat";
      h.textContent = t.category;
      guideResults.appendChild(h);
      lastCat = t.category;
    }

    const card = document.createElement("div");
    card.className = "topic";

    const head = document.createElement("div");
    head.className = "topic-head";
    const h3 = document.createElement("h3");
    h3.innerHTML = highlight(t.title, q);
    const tryBtn = document.createElement("button");
    tryBtn.className = "try-btn";
    tryBtn.textContent = "↳ try";
    tryBtn.title = "Load this into the Tutor pane";
    tryBtn.addEventListener("click", () => {
      tutorEditor.setValue(t.sql);
      closeGuide();
    });
    head.append(h3, tryBtn);

    const p = document.createElement("p");
    p.innerHTML = highlight(t.body, q);
    const pre = document.createElement("pre");
    pre.textContent = t.sql;

    card.append(head, p, pre);
    guideResults.appendChild(card);
  }
}

function escapeHtml(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

// highlight escapes text, then wraps case-insensitive matches of q in <mark>.
function highlight(text, q) {
  const safe = escapeHtml(text);
  if (!q) return safe;
  const escQ = q.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"); // escape regex metachars
  return safe.replace(new RegExp("(" + escQ + ")", "gi"), "<mark>$1</mark>");
}

// --- Tutor (Phase 3 AI helper) -----------------------------------------------

const tutorForm = document.getElementById("tutorForm");
const tutorInput = document.getElementById("tutorInput");
const tutorSend = document.getElementById("tutorSend");
const tutorHints = document.getElementById("tutorHints");

// addBubble appends a chat bubble and returns it (so we can update "thinking").
function addBubble(who, text, extra) {
  const empty = tutorHints.querySelector(".tutor-empty");
  if (empty) empty.remove();
  const b = document.createElement("div");
  b.className = "bubble " + who + (extra ? " " + extra : "");
  b.textContent = text;
  tutorHints.appendChild(b);
  tutorHints.scrollTop = tutorHints.scrollHeight;
  return b;
}

// Render a light bit of markdown: `code` spans and **bold**, escaped safely.
function renderHint(bubble, text) {
  const esc = text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
  bubble.innerHTML = esc
    .replace(/`([^`]+)`/g, "<code>$1</code>")
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
}

async function askTutor(event) {
  event.preventDefault();
  const question = tutorInput.value.trim();
  addBubble("you", question || "(give me a hint on my current query)");
  tutorInput.value = "";
  tutorInput.disabled = true;
  tutorSend.disabled = true;
  const thinking = addBubble("tutor", "thinking…", "thinking");

  try {
    const resp = await fetch("/api/tutor", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ question, sql: userEditor.getValue() }),
    });
    const data = await resp.json();
    if (data.hint) {
      thinking.className = "bubble tutor";
      renderHint(thinking, data.hint);
    } else {
      thinking.className = "bubble tutor error";
      thinking.textContent = data.error || "The tutor didn't respond.";
    }
  } catch (err) {
    thinking.className = "bubble tutor error";
    thinking.textContent = "Could not reach the tutor: " + err.message;
  } finally {
    tutorInput.disabled = false;
    tutorSend.disabled = false;
    tutorInput.focus();
    tutorHints.scrollTop = tutorHints.scrollHeight;
  }
}

// --- Event wiring -------------------------------------------------------------

tutorForm.addEventListener("submit", askTutor);

runBtn.addEventListener("click", runQuery);
examplesBtn.addEventListener("click", openExamples);
guideBtn.addEventListener("click", openGuide);
guideSearch.addEventListener("input", () => renderGuide(guideSearch.value));

// Close the guide via ✕ or backdrop.
guideModal.addEventListener("click", (e) => {
  if (e.target.hasAttribute("data-close")) closeGuide();
});

// "use this" copies the tutor query into the user editor for those who want it.
document.getElementById("copyToUser").addEventListener("click", () => {
  userEditor.setValue(tutorEditor.getValue());
  userEditor.focus();
});

// Close the modal via the ✕, the backdrop, or Escape.
modal.addEventListener("click", (e) => {
  if (e.target.hasAttribute("data-close")) closeExamples();
});
document.addEventListener("keydown", (e) => {
  if (e.key !== "Escape") return;
  if (!modal.hidden) closeExamples();
  if (!guideModal.hidden) closeGuide();
});

// Ctrl/Cmd+Enter runs the query from anywhere in the user editor.
userEditor.setOption("extraKeys", {
  "Ctrl-Enter": runQuery,
  "Cmd-Enter": runQuery,
});

// Run the starter query once on load so the app doesn't open empty.
runQuery();
