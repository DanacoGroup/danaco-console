#!/usr/bin/env bash
# Instalka HYBRYDOWA Windows 10/11 x64 Danaco Console — złożenie pliku .exe (NSIS)
# na Linuksie. Wariant CIENKI: u Operatora staje samo okno.
#
# Czym różni się od `instalka-windows.sh`. Tamten skrypt składa instalkę pełną:
# okno, rdzeń, serwer narzędzi i arsenał w jednym pliku, wszystko na maszynie
# Operatora. Ten składa produkt dla modelu wdrożenia, w którym serce platformy
# stoi na serwerze: instalka niesie powłokę natywną i dokumentację, i NIC więcej.
# Rdzenia (`danaco-console.exe`), serwera narzędzi (`danaco-narzedzia.exe`) ani
# pomocników arsenału w środku nie ma i być nie ma — to jest cały sens tego
# produktu, nie oszczędność na rozmiarze. Rdzeń na serwerze jest jeden, utrzymuje
# go administrator w jednym miejscu, a okno u Operatora nie ma czego aktualizować
# poza sobą.
#
# Który to stan wskazania. Powłoka zna trzy (`src/ustawienia.rs`, `src/main.rs`):
# wskazanie niezłożone, rdzeń na tym urządzeniu, rdzeń na serwerze. Ta instalka
# jest dla trzeciego. Po instalacji obowiązuje jeszcze stan pierwszy — powłoka
# NIE stawia niczego i czeka, aż Operator wskaże host rdzenia w oknie
# (`polecenia::wskaz_rdzen`) albo aż wskaże go zmienna `DANACO_HOST_RDZENIA`.
# Milczeniem się to nie kończy: `rdzen::uruchomienie::zawiaz_oczekiwanie`
# i `zawiaz_zdalnie` składają zdanie do opisu stanu, a nie puste pole.
#
# Dlaczego mimo braku rdzenia instalka niesie pakiet interfejsu. Nie niesie go
# jako plik — niesie go w środku pliku wykonywalnego powłoki. `frontendDist`
# w `tauri.conf.json` jest ustawieniem BUDOWY, nie zasobem instalatora: pakiet
# `client/dist` zostaje wkompilowany w `danaco-console-powloka.exe` i nie da się
# go z instalki wyjąć, nie odbierając powłoce wyjścia awaryjnego. Osobnej kopii
# pakietu instalka nie wiezie: zasób `zasoby/rdzen/dist/` należy do rdzenia
# lokalnego, którego tu nie ma, i jest wprost odejmowany niżej.
#
# Skrypt woła `cargo tauri bundle`, nie `build`: `build` uruchomiłby
# beforeBuildCommand, czyli `npm run build` w budowa/client, a przebudowa
# klienta nie należy do budowy powłoki.
#
# Skrypt nie podpisuje instalatora — podpis Authenticode wymaga certyfikatu
# i hosta Windows — i nie sprawdza, czy instalator się uruchamia; tego nie da
# się sprawdzić bez maszyny z Windows.
#
# Użycie: bash budowa/scripts/instalka-hybryda-win-x64.sh
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
CEL="x86_64-pc-windows-gnu"
PROFIL="tauri.hybryda-win-x64.conf.json"
WERSJA="1.0.0"
WYDANIE="$KORZEN/budowa/wydania/$WERSJA-$(date +%Y-%m-%d)"
NAZWA_WYDANIA="Danaco Console_${WERSJA}_hybryda_x64-setup.exe"

# Rustup instaluje się do ~/.cargo, którego nie ma w PATH powłoki
# nieinteraktywnej. Bez tej linii `cargo` jest nieodnajdywalny.
# shellcheck disable=SC1091
[ -f "$HOME/.cargo/env" ] && . "$HOME/.cargo/env"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Sprawdzenie narzędzi"
for narzedzie in cargo rustup x86_64-w64-mingw32-gcc makensis; do
  command -v "$narzedzie" >/dev/null || padnij "brak narzędzia: $narzedzie (apt: mingw-w64, nsis; rustup: https://rustup.rs)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"
printf '  jest: cel %s\n' "$CEL"
[ -f "$POWLOKA/$PROFIL" ] || padnij "brak profilu budowy: $PROFIL"
printf '  jest: profil %s\n' "$PROFIL"

zglos "Sprawdzenie, czego ta instalka NIE bierze"
# Sprawdzenie odwrotne niż w instalce pełnej. Tam brak rdzenia był powodem
# odmowy; tu jego obecność w zasobach nie jest błędem — profil budowy odejmuje
# go z wykazu zasobów. Wypisujemy stan, żeby było jawne, co zostaje na dysku
# maszyny budującej, a co nie wchodzi do produktu.
for odjety in \
  "zasoby/rdzen/danaco-console.exe" \
  "zasoby/rdzen/danaco-narzedzia.exe" \
  "zasoby/rdzen/dist" \
  "zasoby/pomocniki"; do
  if [ -e "$POWLOKA/$odjety" ]; then
    printf '  odjęty z instalki (leży na maszynie budującej): %s\n' "$odjety"
  else
    printf '  odjęty z instalki (nieobecny i tu): %s\n' "$odjety"
  fi
done

zglos "Dokumentacja produktu"
# Mechanizm przeniesiony z `instalka-windows.sh`: instalator ma nieść
# dokumentację, bo to ON jest produktem końcowym — Operator dostaje jeden plik
# i po instalacji nie ma skąd wziąć instrukcji ani licencji. W wariancie
# hybrydowym waży to więcej, nie mniej: rdzenia obok nie ma, więc dokumentacja
# jest jedyną rzeczą w instalce poza samym oknem.
#
# Katalog zasobu jest WYTWOREM budowania, nie źródłem: pliki mają jedno miejsce
# w korzeniu repozytorium i są tam poprawiane. Kopia powstaje tuż przed
# złożeniem, żeby instalka nie wiozła wersji z poprzedniego tygodnia.
#
# Katalog jest własny dla tego wariantu (`dokumentacja-hybryda`, nie
# `dokumentacja`), bo instalka pełna składa się z tej samej kopii drzewa
# i obie budowy mogą iść równolegle. Wspólny katalog znaczyłby, że jedna budowa
# kasuje zasób drugiej w połowie jej pracy. W zainstalowanym produkcie różnicy
# nie ma: profil budowy odwzorowuje ten katalog na `dokumentacja/`, tak samo jak
# instalka pełna.
DOKUMENTACJA="$POWLOKA/zasoby/dokumentacja-hybryda"
rm -rf "$DOKUMENTACJA"
mkdir -p "$DOKUMENTACJA"
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  [ -f "$KORZEN/$dokument" ] || padnij "brak dokumentu produktu: $dokument — instalator bez licencji i instrukcji nie jest produktem końcowym"
  cp -f "$KORZEN/$dokument" "$DOKUMENTACJA/$dokument"
  printf '  spakowany: %s (%s)\n' "$dokument" "$(du -h "$KORZEN/$dokument" | cut -f1)"
done

zglos "Budowa klienta z wkompilowanym adresem serwera"
# Produkt hybrydowy łączy się SAM z rdzeniem na serwerze — bez konfiguracji
# u Operatora. `adresRdzenia()` (client/src/polaczenie/adres-rdzenia.ts) zwraca
# `VITE_ADRES_RDZENIA` PIERWSZĄ, gdy ustawiona przy budowie klienta; wtedy okno
# łączy się wprost tam, niezależnie od stanu powłoki. Zmienną podaje się WYŁĄCZNIE
# tu, w budowie klienta dla hybryd — nie globalnym `.env`: instalka natywna
# (rdzeń lokalny) i interfejs serwowany przez rdzeń (same-origin) muszą zostać
# bez wkompilowanego adresu. Klient budujemy PRZED powłoką, bo custom-protocol
# osadza `client/dist` w binarce podczas `cargo build`.
( cd "$KORZEN/budowa/client" && VITE_ADRES_RDZENIA="wss://console.danaco-group.pl/ws" npm run build )
[ -f "$KORZEN/budowa/client/dist/index.html" ] || padnij "budowa klienta nie zostawiła dist/index.html"
grep -rq "wss://console.danaco-group.pl/ws" "$KORZEN/budowa/client/dist/assets" \
  || padnij "adres serwera nie wszedł do klienta — VITE_ADRES_RDZENIA nie zostało wkompilowane"
printf '  klient zbudowany z adresem wss://console.danaco-group.pl/ws\n'

zglos "Budowa powłoki dla $CEL"
# Sama powłoka, bez rdzenia. To jedyna binarka, jaką ten produkt wiezie.
#
# `--features tauri/custom-protocol` NIE jest ozdobą i nie wolno go stąd usunąć.
# `tauri::is_dev()` to w tauri 2 dokładnie `!cfg!(feature = "custom-protocol")`
# (`tauri-2.11.5/src/lib.rs`), a `cargo build` sam tej cechy nie włącza — włącza
# ją `cargo tauri build`, którego tu nie wołamy, bo uruchomiłby przebudowę
# klienta. Bez tej cechy gotowy plik Operatora uważa się za budowę deweloperską
# i `zrodlo_interfejsu::ustal` zatrzymuje się na drugim warunku: zwraca
# `WebviewUrl::default()`, które w trybie dev prowadzi do `devUrl`, czyli
# `http://localhost:5173`. Na maszynie Operatora nikt tam nie nasłuchuje, a
# warunek trzeci — „rdzeń nasłuchujący", ten jedyny, na którym ten produkt stoi —
# nie zostaje w ogóle sprawdzony. Okno hybrydowe nie znalazłoby rdzenia na
# serwerze nigdy.
#
# Cecha podana z wiersza poleceń, nie dopisana do `Cargo.toml`: `Cargo.toml`
# jest wspólny dla wszystkich wariantów budowy i nie należy do tego skryptu.
# Powłoka dostaje TEN SAM adres co klient, tylko rozłożony na części. Bez tego
# okno łączyło się z serwerem, a powłoka twierdziła, że nic nie wybrano, i w oknie
# „Stan platformy" pokazywała pętlę zwrotną — mówiła Operatorowi nieprawdę o jego
# własnej instalacji. Wartości czyta `ustawienia.rs` przez `option_env!`, więc
# wchodzą w binarkę przy tej kompilacji i wydania niehybrydowe ich nie niosą.
WBUDOWANE=(
  DANACO_HOST_WBUDOWANY=console.danaco-group.pl
  DANACO_PORT_WBUDOWANY=443
  DANACO_SCHEMAT_WBUDOWANY=https
)
( cd "$POWLOKA" && env "${WBUDOWANE[@]}" cargo build --release --target "$CEL" --features tauri/custom-protocol )

# Sprawdzenie, że powyższe zadziałało — na binarce, nie na wierze w przełącznik.
# Adres serwera rozwojowego w gotowym pliku znaczy budowę deweloperską.
if strings -a "$POWLOKA/target/$CEL/release/danaco-console-powloka.exe" 2>/dev/null | grep -q 'localhost:5173'; then
  padnij "powłoka niesie adres serwera rozwojowego (localhost:5173) — zbudowała się jako budowa deweloperska; okno nie szukałoby rdzenia na serwerze"
fi
printf '  powłoka zbudowana jako produkt (bez adresu serwera rozwojowego)\n'

# Wkompilowany adres serwera jest już potwierdzony na kliencie PRZED osadzeniem
# (`grep` na client/dist/assets wyżej). W skompilowanej powłoce osadzone zasoby
# są SPAKOWANE, więc `strings | grep` adresu tam nie znajdzie — to nie znaczy,
# że go nie ma. Dowodem osadzenia jest brak trybu rozwojowego (sprawdzian wyżej)
# oraz obecność świeżego pliku interfejsu; sam adres gwarantuje budowa klienta.
printf '  adres serwera wkompilowany w klienta osadzonego w powłoce
'

zglos "Złożenie instalatora NSIS z profilu $PROFIL"
# `--config` jest w tauri NAKŁADKĄ na `tauri.conf.json` (scalenie w rodzaju
# JSON Merge Patch), nie zamiennikiem. Dlatego profil hybrydowy nie wymienia
# ikon, identyfikatora ani `frontendDist` — bierze je z konfiguracji podstawowej
# — a zasoby rdzenia odejmuje wartością `null` przy każdym kluczu. Klucza nie
# wystarczy pominąć: pominięty klucz zostaje wzięty z konfiguracji podstawowej.
( cd "$POWLOKA" && cargo tauri bundle --target "$CEL" --bundles nsis --config "$PROFIL" )

# Nazwa produktu jest jedna dla wszystkich wydań — „Danaco Console". Wariant
# niesie nazwa pliku wyjściowego, a nie nazwa widziana przez Operatora: instalator
# i pozycja w systemie mają mówić nazwą produktu, nie nazwą jego odmiany budowlanej.
WYNIK="$(find "$POWLOKA/target/$CEL/release/bundle/nsis" -maxdepth 1 -name 'Danaco Console*-setup.exe' -print -quit)"
[ -n "$WYNIK" ] || padnij "makensis nie zgłosił błędu, ale pliku instalatora nie ma"

zglos "Odbiór — czego w środku być nie może"
# Odbiór jest częścią budowy, nie osobnym krokiem do zapomnienia. Gdyby
# scalenie profilu przestało odejmować zasoby rdzenia, instalka złożyłaby się
# bez błędu i wyszłaby stąd jako produkt pełny pod nazwą hybrydowego. Jedyny
# dowód, że tak nie jest, to wykaz zawartości gotowego pliku.
command -v 7z >/dev/null || padnij "brak 7z (apt: p7zip-full) — bez wykazu zawartości nie ma dowodu, że rdzenia w instalce nie ma"
SPIS="$(mktemp)"
trap 'rm -f "$SPIS"' EXIT
7z l "$WYNIK" >"$SPIS"
for zakazany in 'danaco-console\.exe' 'danaco-narzedzia\.exe'; do
  if grep -qE "$zakazany" "$SPIS"; then
    padnij "instalka hybrydowa niesie $zakazany — profil budowy nie odjął zasobów rdzenia; produkt NIE jest cienki"
  fi
  printf '  nie ma w środku: %s\n' "${zakazany//\\/}"
done
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  grep -q "$dokument" "$SPIS" || padnij "instalka nie niesie dokumentu $dokument"
  printf '  jest w środku: %s\n' "$dokument"
done

zglos "Odłożenie wyniku do wydania"
mkdir -p "$WYDANIE"
cp -f "$WYNIK" "$WYDANIE/$NAZWA_WYDANIA"
ODLOZONY="$WYDANIE/$NAZWA_WYDANIA"

zglos "Pomiar wyniku"
printf 'ścieżka : %s\n' "$ODLOZONY"
printf 'rozmiar : %s (%s bajtów)\n' "$(du -h "$ODLOZONY" | cut -f1)" "$(stat -c%s "$ODLOZONY")"
printf 'typ     : %s\n' "$(file -b "$ODLOZONY")"
printf 'suma    : %s\n' "$(sha256sum "$ODLOZONY" | cut -d' ' -f1)"
# Ostatni wiersz listingu 7z kończy się słowem „files"; liczba stoi przed nim.
printf 'wewnątrz: %s plików\n' "$(tail -1 "$SPIS" | awk '{print $(NF-1)}')"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i czy przechodzi do końca,
  - czy zainstalowane okno wstaje (wymaga WebView2 w systemie),
  - czy okno łączy się z rdzeniem na serwerze — wymaga stojącego rdzenia
    i wskazania hosta (DANACO_HOST_RDZENIA albo wskazanie w oknie),
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie wymagają maszyny z Windows. Nie zakładaj ich powodzenia.
KONIEC
