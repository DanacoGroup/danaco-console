#!/usr/bin/env bash
# Instalka Danaco Console dla Windows 11 na ARM64 — złożenie pliku .exe (NSIS)
# na Linuksie.
#
# Produkt ma jedną postać: hybrydę. Ta instalka wiezie SAMĄ POWŁOKĘ. Rdzeń,
# serwer narzędzi i arsenał stoją na serwerze wdrożenia — Operator dostaje okno,
# nie drugą kopię serca platformy. Dlatego jest cienka i dlatego nie sprawdza
# obecności rdzenia: jego brak nie jest tu usterką, jest założeniem.
#
# ── Droga budowy dla ARM64 ────────────────────────────────────────────────
# Cel to `aarch64-pc-windows-msvc`, składany przez `cargo xwin`, które pobiera
# i trzyma zestaw nagłówków oraz bibliotek importu Microsoftu. Kompilatorem
# krzyżowym jest `clang`, a bibliotekarzem `llvm-lib` z llvm-mingw — dlatego
# `$LLVM_MINGW/bin` MUSI być na ścieżce tej budowy. Cel x64 ma wymaganie
# odwrotne: idzie przez systemowe mingw-w64 i llvm-mingw na ścieżce mu
# przeszkadza. Oba skrypty ustawiają ścieżkę same, żeby jedna budowa nie
# psuła drugiej.
#
# Skrypt woła `cargo tauri bundle`, nie `build`: `build` uruchomiłby
# beforeBuildCommand, czyli przebudowę klienta, a ta do budowy powłoki nie
# należy. `bundle` niczego nie kompiluje — bierze gotową binarkę z `target/`.
#
# Nakładki konfiguracyjnej nie ma i nie jest potrzebna: `tauri.conf.json` opisuje
# wprost produkt hybrydowy, bo innego produktu nie ma.
#
# Skrypt nie podpisuje instalatora — podpis Authenticode wymaga certyfikatu
# i hosta Windows — i nie sprawdza, czy instalator się uruchamia; tego nie da
# się sprawdzić bez maszyny z Windows na ARM64.
#
# Użycie: bash budowa/scripts/instalka-hybryda-win-arm.sh
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
KLIENT="$KORZEN/budowa/klient"
CEL="aarch64-pc-windows-msvc"
LLVM_MINGW="${LLVM_MINGW:-/opt/llvm-mingw}"
WERSJA="1.0.0"
WYDANIE="$KORZEN/budowa/wydania/$WERSJA-$(date +%Y-%m-%d)"
NAZWA_WYDANIA="Danaco Console_${WERSJA}_hybryda_arm64-setup.exe"

export PATH="$LLVM_MINGW/bin:$PATH"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Sprawdzenie narzędzi"
# Brak któregokolwiek ujawniłby się dopiero po kilkuminutowej budowie, dlatego
# sprawdzenie stoi przed nią, nie po.
for narzedzie in cargo rustup cargo-xwin makensis clang llvm-lib; do
  command -v "$narzedzie" >/dev/null \
    || padnij "brak narzędzia: $narzedzie (apt: nsis, clang; cargo install cargo-xwin; llvm-lib z wydania llvm-mingw rozpakowanego do $LLVM_MINGW)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"
printf '  jest: cel %s\n' "$CEL"

zglos "Sprawdzenie artefaktów wejściowych"
# Powłoka niesie własne okno, więc zbudowany klient jest jej jedynym artefaktem
# wejściowym. Rdzenia nie sprawdzamy — w tej instalce go nie ma.
[ -e "$KLIENT/dist/index.html" ] \
  || padnij "brak artefaktu: budowa/klient/dist/index.html (zbuduj: npm run budowanie w budowa/klient)"
printf '  jest: budowa/klient/dist/index.html\n'

zglos "Budowa powłoki dla $CEL"
# `--features tauri/custom-protocol` NIE jest ozdobą — bez niej produkt jest
# zepsuty w sposób niewidoczny. W tauri 2 `tauri::is_dev()` to dokładnie
# `!cfg!(feature = "custom-protocol")`, a `Cargo.toml` powłoki tej cechy nie
# włącza i `cargo build` sam jej nie doda (`cargo tauri bundle` też nie — on
# niczego nie kompiluje). Bez cechy `zrodlo_interfejsu::ustal` zatrzymuje się
# na `is_dev()`, NIGDY nie sprawdzając warunku „rdzeń nasłuchujący" — a dla
# hybrydy to jedyna droga do interfejsu. Okno pokazałoby pustą stronę.
( cd "$POWLOKA" \
  && cargo xwin build --release --target "$CEL" \
       --features tauri/custom-protocol --cross-compiler clang )

POWLOKA_EXE="$POWLOKA/target/$CEL/release/danaco-console-powloka.exe"
[ -f "$POWLOKA_EXE" ] || padnij "cargo nie zgłosił błędu, ale binarki powłoki nie ma: $POWLOKA_EXE"
# Cel budowy da się pomylić w wierszu poleceń, a instalka z binarką x64 w środku
# jest dla maszyny ARM64 bezużyteczna. To jedyny sprawdzian, który tu naprawdę
# coś znaczy, więc stoi przed złożeniem, nie tylko po nim.
# `file` nazywa tę architekturę „ARM64", a nie „Aarch64" jak rustup — wzorzec
# dopuszcza obie pisownie, bo to nazwa zależna od wersji `file`, nie od produktu.
OPIS_POWLOKI="$(file -b "$POWLOKA_EXE")"
printf '%s' "$OPIS_POWLOKI" | grep -qi 'PE32+ executable' \
  || padnij "powłoka nie jest binarką Windows (PE32+): $OPIS_POWLOKI"
printf '%s' "$OPIS_POWLOKI" | grep -qiE 'ARM64|Aarch64' \
  || padnij "powłoka nie jest binarką PE dla ARM64: $OPIS_POWLOKI"
printf '  architektura powłoki: %s\n' "$OPIS_POWLOKI"

zglos "Zapora: osadzone zasoby interfejsu"
# Po czym POZNAĆ, że cecha custom-protocol weszła. Nie po napisie
# `localhost:5173`: Tauri wkompilowuje w binarkę całą konfigurację produktu,
# więc obecność albo brak tego napisu mówi o polu `devUrl` w konfiguracji,
# a nie o tym, która droga do interfejsu jest czynna. Zapora na ten napis
# milczy przy pliku zepsutym — jest gorsza niż jej brak.
#
# Rozstrzyga obecność OSADZONYCH ZASOBÓW: wchodzą do binarki wyłącznie przy
# czynnej cesze custom-protocol. Sondą jest nazwa pliku interfejsu z sumą treści
# w nazwie (np. `index-Bx8zeXAO.js`) — taki napis nie ma jak trafić do binarki
# inaczej niż przez osadzenie `klient/dist`, i zmienia się z każdą przebudową
# klienta, dlatego czytamy go z katalogu, a nie wpisujemy tu na sztywno.
#
# Mylące sondy, których tu NIE używamy: `woff2` (tauri ma tablicę typów MIME,
# jedno trafienie jest zawsze) i `index.html` (tauri obsługuje indeks katalogu,
# trafień jest kilka) — obie zapalają się w binarce BEZ osadzonych zasobów.
SONDA="$(find "$KLIENT/dist/assets" -maxdepth 1 -name 'index-*.js' -printf '%f\n' 2>/dev/null | head -1)"
[ -n "$SONDA" ] \
  || padnij "nie ma z czego zrobić sondy: brak budowa/klient/dist/assets/index-*.js — bez niej nie da się sprawdzić, czy zasoby są osadzone"
printf '  sonda: %s\n' "$SONDA"
# `strings … | grep -q` byłoby tu PUŁAPKĄ przy `set -o pipefail`: grep -q kończy
# się na pierwszym trafieniu, `strings` dostaje SIGPIPE, a pipefail zamienia to
# w porażkę całego potoku — zapora zapalałaby się DOKŁADNIE wtedy, gdy sonda się
# znajdzie. Dlatego liczymy trafienia w podstawieniu polecenia.
LICZBA_SONDY="$(strings -a "$POWLOKA_EXE" | grep -c -- "$SONDA" || true)"
LICZBA_ZASOBOW="$(strings -a "$POWLOKA_EXE" | grep -c '/assets/' || true)"
printf '  trafienia sondy: %s, wystąpienia /assets/: %s\n' "$LICZBA_SONDY" "$LICZBA_ZASOBOW"
{ [ "$LICZBA_SONDY" -gt 0 ] && [ "$LICZBA_ZASOBOW" -gt 0 ]; } \
  || padnij "powłoka NIE ma osadzonych zasobów interfejsu — cecha tauri/custom-protocol nie weszła, więc tauri::is_dev() zwróci prawdę i okno pójdzie do serwera rozwojowego zamiast do rdzenia na serwerze. To produkt zepsuty, choć instalka wyglądałaby poprawnie"
printf '  osadzone zasoby interfejsu: potwierdzone\n'

zglos "Złożenie instalatora NSIS"
# `--bundles nsis` podane jawnie, choć konfiguracja ma ten cel ustawiony:
# skrypt składa instalator NSIS niezależnie od tego, jakie inne cele
# konfiguracja produktu wymienia.
( cd "$POWLOKA" && cargo tauri bundle --target "$CEL" --bundles nsis )

ZLOZONY="$(find "$POWLOKA/target/$CEL/release/bundle/nsis" -maxdepth 1 -name '*-setup.exe' -print -quit)"
[ -n "$ZLOZONY" ] || padnij "makensis nie zgłosił błędu, ale pliku instalatora nie ma"

zglos "Odbiór — czego w środku być nie może"
# Odbiór jest częścią budowy, nie osobnym krokiem do zapomnienia. Gdyby do
# konfiguracji wróciły zasoby rdzenia, instalka złożyłaby się bez błędu
# i wyszłaby stąd jako produkt pełny pod nazwą hybrydowego. Jedynym dowodem,
# że tak nie jest, jest wykaz zawartości gotowego pliku.
command -v 7z >/dev/null \
  || padnij "brak 7z (apt: p7zip-full) — bez wykazu zawartości nie ma dowodu, że rdzenia w instalce nie ma"
SPIS="$(mktemp)"
ROZPAK="$(mktemp -d)"
trap 'rm -f "$SPIS"; rm -rf "$ROZPAK"' EXIT
7z l "$ZLOZONY" >"$SPIS"
for zakazany in danaco-console.exe danaco-narzedzia.exe; do
  [ "$(grep -c -- "$zakazany" "$SPIS" || true)" -eq 0 ] \
    || padnij "instalka niesie $zakazany — do produktu wrócił rdzeń; instalka NIE jest cienka"
  printf '  nie ma w środku: %s\n' "$zakazany"
done

# Powyższe zapory badały binarkę z `target/`. Ta bada binarkę WYPAKOWANĄ
# z instalatora, bo tylko ona jest tym, co dostanie Operator: `bundle` bierze
# gotową binarkę z `target/` i sam nie kompiluje, więc podmiany nie byłoby widać
# inaczej.
7z x -o"$ROZPAK" "$ZLOZONY" >/dev/null 2>&1 \
  || padnij "nie udało się wypakować instalatora do sprawdzenia"
WIEZIONA="$ROZPAK/danaco-console-powloka.exe"
[ -f "$WIEZIONA" ] || padnij "w instalatorze nie ma powłoki danaco-console-powloka.exe"
OPIS_WIEZIONEJ="$(file -b "$WIEZIONA")"
printf '%s' "$OPIS_WIEZIONEJ" | grep -qiE 'ARM64|Aarch64' \
  || padnij "powłoka WIEZIONA przez instalator nie jest binarką ARM64: $OPIS_WIEZIONEJ"
TRAFIENIA_WIEZIONEJ="$(strings -a "$WIEZIONA" | grep -c -- "$SONDA" || true)"
[ "$TRAFIENIA_WIEZIONEJ" -gt 0 ] \
  || padnij "powłoka WIEZIONA przez instalator nie ma osadzonych zasobów interfejsu — instalka wiezie budowę rozwojową"
printf '  wieziona powłoka: %s, zasoby osadzone\n' "$OPIS_WIEZIONEJ"

zglos "Odłożenie wyniku do wydania"
mkdir -p "$WYDANIE"
WYNIK="$WYDANIE/$NAZWA_WYDANIA"
cp -f "$ZLOZONY" "$WYNIK"

zglos "Pomiar wyniku"
printf 'ścieżka : %s\n' "$WYNIK"
printf 'rozmiar : %s (%s bajtów)\n' "$(du -h "$WYNIK" | cut -f1)" "$(stat -c%s "$WYNIK")"
printf 'typ     : %s\n' "$(file -b "$WYNIK")"
printf 'suma    : %s\n' "$(sha256sum "$WYNIK" | cut -d' ' -f1)"
# Ostatni wiersz listingu 7z kończy się słowem „files"; liczba stoi przed nim.
printf 'wewnątrz: %s plików\n' "$(tail -1 "$SPIS" | awk '{print $(NF-1)}')"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i czy przechodzi do końca,
  - czy zainstalowane okno wstaje (wymaga WebView2 dla ARM64 w Windows 11),
  - czy okno łączy się z rdzeniem na serwerze wdrożenia — wymaga stojącego
    rdzenia i wskazania hosta (DANACO_HOST_RDZENIA albo wskazanie w oknie),
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie wymagają maszyny z Windows na ARM64. Nie zakładaj ich powodzenia.
KONIEC
