# Danaco Console — przewodnik wykonawcy

Aplikacja desktopowa (Tauri 2): rdzeń w Go, interfejs TypeScript/Vite, poczta IMAP, WebSocket.

**Przed pierwszą zmianą przeczytaj [ZRODLA-PRAWDY.md](ZRODLA-PRAWDY.md).** W skrócie:
- `docs/` — jedyne źródło prawdy funkcji (51 opracowań, zbiór zamknięty). TYLKO DO CZYTANIA.
- `design/` — prototypy HTML, żetony `--dn-*`, komponenty `.dn-*`. TYLKO DO CZYTANIA.
- `budowa/` — cały kod. Tu się pisze.
- Wartości normatywne: `budowa/shared/contract.json` (komendy, zdarzenia, kody błędów),
  `design/zasoby/zetony/zetony.css` (żetony), `design/05-okna/` (brzmienia etykiet).
- Niczego nie dopowiadasz z głowy — brak nazwy/wartości w opracowaniu to usterka do zgłoszenia,
  nie swoboda wykonawcy. Miejsca otwarte są oznaczone **[DO DECYZJI OPERATORA]**.

## Układ `budowa/`

- `server/` — rdzeń Go (`cmd/danaco-console` — wejście, `internal/` — pakiety).
- `shared/` — kontrakt (`contract.json`, `contract.go`); moduł Go `danacoconsole` obejmuje `server/` i `shared/`.
- `client/` — interfejs Vite + TypeScript (bez frameworka), testy Vitest.
- `desktop/src-tauri/` — powłoka Tauri 2 (Rust).

## Polecenia

```bash
# Go (z katalogu budowa/)
go build ./...                 # kompilacja całości
gotestsum ./...                # testy z czytelnym raportem
golangci-lint run ./...        # lint
# Klient (z budowa/client/)
npm run typy                   # tsc --noEmit
npm run testy                  # vitest run
npm run build                  # pełna budowa
# Desktop (z budowa/desktop/src-tauri/) — sccache+mold już skonfigurowane w ~/.cargo/config.toml
cargo build
```

## Praca równoległa (worktree)

- `main` — stan scalony; sesja pracuje na gałęzi `teren/<nazwa>` we własnym worktree:
  `git worktree add ~/robocze/<nazwa> -b teren/<nazwa>`.
- Podziały terenów (listy plików per teren) leżą w `~/robocze/tereny/*.txt`.
- Tereny się nie nakładają. Commit obejmuje wyłącznie pliki własnego terenu — nigdy `git add -A`.
- Po `worktree add` uruchom `budowa/przygotuj-drzewo.sh` (dowiązuje `node_modules` z drzewa głównego).
- Cache'e wspólne dla wszystkich worktree: GOMODCACHE/GOCACHE (Go), sccache (Rust), `node_modules` przez dowiązanie.
- Wpis w dzienniku: tryb oznajmujący, jedno zdanie, do 70 znaków.

## Narzędzia w środowisku (wszystkie zainstalowane — nie instaluj ponownie)

Go 1.26 (gopls, golangci-lint, staticcheck, sqlc, oapi-codegen, mockery, dlv, gotestsum, goimports) ·
Node 22, npm, pnpm, bun, Vite, tsc, Vitest · Playwright 1.62 z pobranymi przeglądarkami (chromium, firefox) ·
Rust 1.97 + tauri-cli 2.11 + webkit2gtk/gtk3 (pełne zależności systemowe) · sccache + mold + clang ·
Docker, SQLite 3.46, ripgrep, jq, Taskfile (`task`), ImageMagick. Maszyna: 16 rdzeni, 122 GB RAM.

Weryfikacja zmiany przed zamknięciem terenu: `go build ./...` + `gotestsum` (serwer),
`npm run typy` + `npm run testy` (klient). Nie zamykaj terenu z czerwonym wynikiem.
