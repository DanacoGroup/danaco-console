# Środowisko maszyny budowy

Opis maszyny, na której powstaje Danaco Console. Dokument prowadzi stan
faktyczny — każda pozycja została sprawdzona uruchomieniem, nie odtworzona
z pamięci.

**Sesja wykonawcza niczego tu nie instaluje; brakujące narzędzie jest zgłoszeniem
do Prowadzącego.** Wyjątek udziela się terenowi wpisem w rejestrze terenów, nazywa
go wprost i wygasa razem z tym terenem. Udzielono go dotąd raz — terenowi
`zaplecze-modeli`, bo Właściciel polecił, żeby aplikacja miała wszystkie narzędzia
czynne, a 16 GB wag w `/opt/danaco-modele` nie miało czym się uruchomić. Wszystko,
co teren postawił, stoi w sekcji 4.

## 1. Zasoby maszyny

Szesnaście rdzeni, 122 GB pamięci operacyjnej, 387 GB przestrzeni dyskowej
(wykorzystane 69 GB). System Linux, powłoka Bash.

## 2. Łańcuchy narzędzi

| Narzędzie | Wersja | Uwaga |
|---|---|---|
| Go | 1.26.5 | `/usr/local/go/bin`; narzędzia w `~/go/bin` |
| Node | 22.23.2 | npm 10.9.8, pnpm 11.22.0, bun 1.4.0 |
| TypeScript | 7.0.2 | |
| Rust | 1.97.1 | `CARGO_HOME=/opt/rust/cargo`, `RUSTUP_HOME=/opt/rust/rustup` |
| tauri-cli | 2.11.4 | `cargo tauri` |
| Playwright | 1.62.1 | przeglądarki w `/opt/ms-playwright` |
| SQLite | 3.46.1 | |
| Docker | 29.1.3 | |
| Git | 2.53.0 | |
| Python | 3.14.4 | |

Narzędzia Go: gopls, golangci-lint, staticcheck, sqlc, oapi-codegen, mockery,
dlv, gotestsum, goimports, dupl.

Warstwa systemowa powłoki desktopowej: webkit2gtk 2.52.3, gtk+ 3.24.52.
Przyspieszenie budowy Rust: sccache 0.17.0 jako otoczka kompilatora oraz mold
2.40.4 jako konsolidator — skonfigurowane w `/opt/rust/cargo/config.toml`.
Kompilator pomocniczy: clang 21.1.8. Budowa krzyżowa pod Windows: cargo-xwin
oraz llvm-mingw w `/opt`.

Narzędzia pomocnicze: ripgrep, jq, Taskfile, ImageMagick, Pandoc, Tesseract,
ffmpeg, LibreOffice, pdftotext, Inkscape, FontForge, vale, semgrep.

## 3. Zmienne środowiska

Ustawione w `~/.bashrc`, wspólne dla wszystkich drzew roboczych:

| Zmienna | Wartość | Powód |
|---|---|---|
| `CARGO_HOME` | `/opt/rust/cargo` | instalacja Rust poza katalogiem domowym |
| `RUSTUP_HOME` | `/opt/rust/rustup` | jak wyżej |
| `PLAYWRIGHT_BROWSERS_PATH` | `/opt/ms-playwright` | przeglądarki pobrane raz, wspólne dla drzew |

Pamięci podręczne wspólne dla wszystkich drzew roboczych: `GOMODCACHE`
(`~/go/pkg/mod`), `GOCACHE` (`~/.cache/go-build`), sccache dla Rust,
`node_modules` przez dowiązanie zakładane skryptem otwarcia terenu.

## 4. Zasoby prowizjonowane

Materiały pobrane raz, wspólne dla całej budowy. Nie należą do repozytorium
i nie podlegają ponownemu pobraniu.

| Katalog | Rozmiar | Zawartość |
|---|---|---|
| `/opt/danaco-modele` | 16 GB | modele lokalne: osadzenia, przeszukiwanie wtórne, rozpoznanie obrazu, synteza mowy, rozpoznanie twarzy, powiększanie obrazu |
| `/opt/danaco-arsenal-windows` | 8,0 GB | zasoby wdrożenia pod Windows |
| `/opt/danaco` | 784 MB | pomocnicze środowiska uruchomieniowe |
| `/opt/danaco-arsenal` | 251 MB | synteza mowy wraz z głosami |
| `/opt/danaco-ikony` | 138 MB | zbiory ikon: heroicons, lucide, phosphor, tabler |
| `/opt/ms-playwright` | — | chromium, firefox, powłoka bezokienkowa, ffmpeg |

## 4a. Narzędzia dołożone po przejęciu rdzenia

Wykaz zależności rdzenia urósł z 30 pozycji przy przejęciu do 54. Poniżej to,
czego nie wymienia sekcja 2, a na czym opierają się komendy dołożone w tej turze.
Każda pozycja sprawdzona uruchomieniem.

| Narzędzie | Skąd | Komenda, która po nie sięga |
|---|---|---|
| `pa11y` | npm globalne | `browser.accessibility.audit` |
| `lighthouse` | npm globalne | `apps.performance.audit` |
| `autocannon` | npm globalne | `developer.api.load.run` |
| `k6` | `/usr/local/bin` | rozważony i odrzucony na rzecz `autocannon` |
| `fastembed` 0.8.0 wraz z czterema pakietami | pip, **~434 MB** | `knowledge.index`, `knowledge.search` |

`fastembed` postawił teren `zaplecze-modeli` na mocy wyjątku opisanego w nagłówku
tego dokumentu. W modelu hybrydowym zaplecze stoi na serwerze wdrożenia, nie
u Operatora — instalacja na maszynie budowlanej jest instalacją tego serwera.

**Srodowisko odtwarzania twarzy — postawione przez teren `odtwarzanie-twarzy`**,
na tym samym wyjatku:

| Rzecz | Waga |
|---|---|
| `/opt/danaco/silniki/twarze` — srodowisko pythonowe (torch, torchvision, facexlib, opencv, numpy, scipy) | **1,6 GB** |
| `/usr/local/bin/danaco-twarze` — opakowanie wolane przez rdzen | 402 B |
| `detection_Resnet50_Final.pth` (RetinaFace) w `/opt/danaco-modele/twarze/` | 104 MB |
| `parsing_parsenet.pth` (ParseNet) tamze | 82 MB |

**Stos modeli PyTorch — postawiony przez teren `zdolnosc-wyszukiwania`**, na tym
samym wyjatku co `fastembed`. Instalacja `pip --user --break-system-packages`,
bo system jest `externally-managed`:

| Pakiet | Waga |
|---|---|
| torch 2.13.0+cpu | 757 MB |
| transformers 5.16.1 | 119 MB |
| sympy, networkx, mpmath, regex, safetensors, typer i drobne | ok. 105 MB |

Razem okolo **980 MB**; caly katalog `~/.local/lib/python3*/site-packages` wazy
1,4 GB. Zadnych wag modeli nie pobrano — przesiew i os obrazu licza na wagach
stojacych w `/opt/danaco-modele`. Wolne miejsce po pracy: 276 GB.

**Brak jedyny:** `typescript-language-server`. Rdzeń nazywa go w sondzie startowej
wraz z drogą naprawy; warstwa językowa TypeScriptu przez to nie działa.

## 5. Usługi w tle

Na maszynie pracują usługi niezwiązane z budową Danaco Console: PostgreSQL,
Redis, Elasticsearch, Neo4j, Ollama, Dovecot (IMAP), Caddy, nginx, Docker.
Zajmują porty i pamięć. Rdzeń produktu, gdy powstanie, nie może kolidować z ich
portami — dobór portu jest rozstrzygnięciem do zapisania w rejestrze decyzji.

## 6. Zasada porządku

Procesy uruchomione przez sesję są przez tę sesję wygaszane przed jej
zamknięciem. Pliki robocze powstają w katalogu tymczasowym sesji i giną wraz
z nią. Drzewo robocze zamkniętego terenu jest usuwane. Sesja, która zostawia po
sobie działający proces albo katalog roboczy, nie zamknęła pracy.

## Postawione 28.08 przy odtwarzaniu drogi budowania instalatora

Toolchain Rusta zniknął z maszyny (`~/.cargo/bin` nie istniał, `~/.rustup` był
pusty), przez co powłoki Tauri ani instalatora Windows nie dało się złożyć.
Postawione na polecenie Właściciela:

| Narzędzie | Wydanie | Waga | Po co |
|---|---|---|---|
| rustup wraz z rustc i cargo | 1.98.0, profil minimalny | ~1,4 GB z celami | budowa powłoki Tauri |
| cel `x86_64-pc-windows-gnu` | — | w powyższym | wariant x64 instalatora |
| cel `aarch64-pc-windows-msvc` | — | w powyższym | wariant ARM64 instalatora |
| `cargo-tauri` | 2.11.4 | ~40 MB | złożenie pakietu NSIS |
| `cargo-xwin` | 0.23.1 | ~25 MB | biblioteki Windows SDK dla celu MSVC na Linuksie |

Nie instalowano niczego poza tym. Vite 8.2.2, TypeScript 7.0.2, `mingw-w64`,
`makensis` oraz llvm-mingw w `/opt/llvm-mingw` stały już na maszynie.

