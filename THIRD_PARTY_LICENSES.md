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

## teamdandelion / PokemonSQLTutorial (bundled database & example queries)

The `PokemonSQLTutorial-master/` directory — including `pokedex.sqlite` and the
example `.sql` files — comes from this SQL tutorial project.

- Project: https://github.com/teamdandelion/PokemonSQLTutorial
- See that repository for its licensing terms.

## veekun / pokedex (source of the Pokémon data)

The bundled SQLite database is generated from the veekun pokedex project.

- Project: https://github.com/veekun/pokedex
- The database schema and tooling are provided by veekun. See the veekun
  repository for its licensing terms.

## PokéAPI / sprites (Pokémon sprite images)

The web app fetches Pokémon sprite images at runtime from the PokéAPI sprite
set (`internal/gui/sprites.go`). The images are not redistributed in this repo;
they are loaded on demand from:

- Project: https://github.com/PokeAPI/sprites
- License: BSD 3-Clause

```
Copyright (c) © 2013–2023 Paul Hallett and PokéAPI contributors
(https://github.com/PokeAPI/pokeapi#contributing). Pokémon and Pokémon
character names are trademarks of Nintendo.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.
* Neither the name of PokéAPI nor the names of its contributors may be used to
  endorse or promote products derived from this software without specific prior
  written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```

**Note on Pokémon data:** Pokémon and Pokémon character names are trademarks of
Nintendo, Game Freak, and The Pokémon Company. This project is an unofficial,
non-commercial educational tool and is not affiliated with or endorsed by them.
The data is included here solely to make learning SQL more engaging.
