#!/usr/bin/env bash
# Instalka Danaco Console dla Windows 11 na x64 — złożenie pliku .exe (NSIS)
# na Linuksie.
#
# Produkt ma jedną postać: hybrydę. U Operatora staje samo okno, a rdzeń, serwer
# narzędzi i arsenał stoją na serwerze wdrożenia. Instalka rdzenia nie niesie
# i nieść nie ma — to jest sens tego produktu, nie oszczędność na rozmiarze:
# rdzeń jest jeden, utrzymuje go administrator w jednym miejscu, a okno
# u Operatora nie ma czego aktualizować poza sobą.
#
# Który to stan wskazania. Powłoka zna trzy (`src/ustawienia.rs`, `src/main.rs`):
# wskazanie niezłożone, rdzeń na tym urządzeniu, rdzeń na serwerze. Ta instalka
# jest dla trzeciego. Po instalacji obowiązuje jeszcze stan pierwszy — powłoka
# NIE stawia niczego i czeka, aż Operator wskaże host rdzenia w oknie
# (`polecenia::wskaz_rdzen`) albo aż wskaże go zmienna `DANACO_HOST_RDZENIA`.
#
# Dlaczego mimo braku rdzenia instalka niesie pakiet interfejsu. Nie niesie go
# jako plik — niesie go w środku pliku wykonywalnego powłoki. `frontendDist`
# w `tauri.conf.json` jest ustawieniem BUDOWY, nie zasobem instalatora: pakiet
# `klient/dist` zostaje wkompilowany w `danaco-console-powloka.exe` i nie da się
# go z instalki wyjąć, nie odbierając powłoce wyjścia awaryjnego.
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
# się sprawdzić bez maszyny z Windows.
#
# Użycie: bash budowa/scripts/instalka-hybryda-win-x64.sh
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
KLIENT="$KORZEN/budowa/klient"
CEL="x86_64-pc-windows-gnu"
WERSJA="1.0.0"
WYDANIE="$KORZEN/budowa/wydania/$WERSJA-$(date +%Y-%m-%d)"
NAZWA_WYDANIA="Danaco Console_${WERSJA}_hybryda_x64-setup.exe"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Sprawdzenie narzędzi"
# Brak któregokolwiek ujawniłby się dopiero po kilkuminutowej budowie, dlatego
# sprawdzenie stoi przed nią, nie po.
#
# `llvm-mingw` NIE ma tu być na ścieżce. Cel `x86_64-pc-windows-gnu` konsoliduje
# się systemowym mingw-w64 i potrzebuje jego `libgcc`; ścieżka z llvm-mingw
# podstawia własny łańcuch i budowa pada na brakującej bibliotece.
for narzedzie in cargo rustup x86_64-w64-mingw32-gcc makensis; do
  command -v "$narzedzie" >/dev/null \
    || padnij "brak narzędzia: $narzedzie (apt: mingw-w64, nsis; rustup: https://rustup.rs)"
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
# Sama powłoka, bez rdzenia. To jedyna binarka, jaką ten produkt wiezie.
#
# `--features tauri/custom-protocol` NIE jest ozdobą i nie wolno go stąd usunąć.
# `tauri::is_dev()` to w tauri 2 dokładnie `!cfg!(feature = "custom-protocol")`,
# a `cargo build` sam tej cechy nie włącza — włącza ją `cargo tauri build`,
# którego tu nie wołamy. Bez cechy gotowy plik Operatora uważa się za budowę
# deweloperską i `zrodlo_interfejsu::ustal` zatrzymuje się na drugim warunku,
# nigdy nie sprawdzając warunku „rdzeń nasłuchujący" — jedynego, na którym ten
# produkt stoi. Okno nie znalazłoby rdzenia na serwerze nigdy.
#
# Cecha podana z wiersza poleceń, nie dopisana do `Cargo.toml`: `Cargo.toml`
# jest wspólny dla obu celów budowy i nie należy do tego skryptu.
( cd "$POWLOKA" && cargo build --release --target "$CEL" --features tauri/custom-protocol )

POWLOKA_EXE="$POWLOKA/target/$CEL/release/danaco-console-powloka.exe"
[ -f "$POWLOKA_EXE" ] || padnij "cargo nie zgłosił błędu, ale binarki powłoki nie ma: $POWLOKA_EXE"
OPIS_POWLOKI="$(file -b "$POWLOKA_EXE")"
printf '%s' "$OPIS_POWLOKI" | grep -qi 'PE32+ executable' \
  || padnij "powłoka nie jest binarką Windows (PE32+): $OPIS_POWLOKI"
printf '%s' "$OPIS_POWLOKI" | grep -qi 'x86-64' \
  || padnij "powłoka nie jest binarką PE dla x86-64: $OPIS_POWLOKI"
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
printf '%s' "$OPIS_WIEZIONEJ" | grep -qi 'x86-64' \
  || padnij "powłoka WIEZIONA przez instalator nie jest binarką x86-64: $OPIS_WIEZIONEJ"
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
  - czy zainstalowane okno wstaje (wymaga WebView2 w systemie Windows 11),
  - czy okno łączy się z rdzeniem na serwerze wdrożenia — wymaga stojącego
    rdzenia i wskazania hosta (DANACO_HOST_RDZENIA albo wskazanie w oknie),
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie wymagają maszyny z Windows. Nie zakładaj ich powodzenia.
KONIEC
