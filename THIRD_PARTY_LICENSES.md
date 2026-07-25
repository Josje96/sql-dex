# Third-Party Licenses

SQL-Dex bundles the following third-party software and data. Each is the
property of its respective authors and is used under the terms below.

---

## CodeMirror

Bundled at `internal/gui/static/vendor/codemirror.js`, `codemirror.css`, and
`material-darker.css`. Used to power the SQL editor panes.

- Project: https://codemirror.net
- License: MIT

```
MIT License

Copyright (C) by Marijn Haverbeke <marijnh@gmail.com> and others

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

---

## sql.js

Bundled at `internal/gui/static/vendor/sql.js`. SQLite compiled to
WebAssembly / JavaScript, used to run queries in the browser.

- Project: https://github.com/sql-js/sql.js
- License: MIT

```
MIT License

Copyright (c) 2017 sql.js authors (see AUTHORS file in the sql.js project)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

---

## veekun / pokedex (Pokémon database)

The bundled SQLite database (`PokemonSQLTutorial-master/pokedex.sqlite`) and the
example queries are derived from the veekun pokedex project.

- Project: https://github.com/veekun/pokedex
- The database schema and tooling are provided by veekun. See the veekun
  repository for its licensing terms.

**Note on Pokémon data:** Pokémon and Pokémon character names are trademarks of
Nintendo, Game Freak, and The Pokémon Company. This project is an unofficial,
non-commercial educational tool and is not affiliated with or endorsed by them.
The data is included here solely to make learning SQL more engaging.
