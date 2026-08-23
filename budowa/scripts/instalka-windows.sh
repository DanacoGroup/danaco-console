#!/usr/bin/env bash
# Instalka Windows Danaco Console — złożenie pliku .exe (NSIS) na Linuksie.
#
# Cel budowy: `x86_64-pc-windows-gnu` z konsolidatorem mingw-w64 z apt. Wariant
# MSVC przez cargo-xwin jest drogą zapasową na wypadek, gdyby któraś zależność
# przestała się składać pod GNU; wine nie jest potrzebne. Zależności powłoki
# (`webview2-com`, `tao`, `windows-rs`, `ureq` z rustls) składają się pod GNU
# bez poprawek, a cel GNU nie wymaga zestawu nagłówków Microsoftu ani zgody
# na jego licencję.
#
# Skrypt woła `cargo tauri bundle`, nie `build`: `build` uruchomiłby
# beforeBuildCommand, czyli `npm run build` w budowa/client, a przebudowa
# klienta nie należy do budowy powłoki. `bundle` bierze gotowe artefakty,
# dlatego skrypt sam sprawdza ich obecność i odmawia zamiast wydać pusty produkt.
#
# Cel budowy podawany jest z wiersza poleceń, bo tauri.conf.json opisuje produkt
# (nazwa, identyfikator, ikony, zasoby), a cel jest sprawą maszyny budującej.
#
# Skrypt nie podpisuje instalatora — podpis Authenticode wymaga certyfikatu
# i hosta Windows — i nie sprawdza, czy instalator się uruchamia; tego nie da
# się sprawdzić bez maszyny z Windows.
#
# Użycie: bash budowa/scripts/instalka-windows.sh
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
CEL="x86_64-pc-windows-gnu"

# Rustup instaluje się do ~/.cargo, którego nie ma w PATH powłoki
# nieinteraktywnej. Bez tej linii `cargo` jest nieodnajdywalny.
# shellcheck disable=SC1091
[ -f "$HOME/.cargo/env" ] && . "$HOME/.cargo/env"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Sprawdzenie narzędzi"
# Konsolidator mingw daje binarkę PE; makensis składa z niej instalator NSIS.
# Oba są w apt (mingw-w64, nsis) i oba muszą być, zanim cokolwiek ruszy —
# brak makensis ujawniłby się dopiero po kilkuminutowej budowie.
for narzedzie in cargo rustup x86_64-w64-mingw32-gcc makensis; do
  command -v "$narzedzie" >/dev/null || padnij "brak narzędzia: $narzedzie (apt: mingw-w64, nsis; rustup: https://rustup.rs)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"
printf '  jest: cel %s\n' "$CEL"

zglos "Sprawdzenie artefaktów wejściowych"
# Instalator pakuje te pliki (bundle.resources w tauri.conf.json). Bez nich NSIS
# złożyłby instalator pozbawiony rdzenia. Buduje je budowa/scripts/wydanie.sh.
for artefakt in \
  "$POWLOKA/zasoby/rdzen/danaco-console.exe" \
  "$POWLOKA/zasoby/rdzen/dist/index.html" \
  "$KORZEN/budowa/client/dist/index.html"; do
  [ -e "$artefakt" ] || padnij "brak artefaktu: $artefakt (zbuduj: bash budowa/scripts/wydanie.sh)"
  printf '  jest: %s\n' "${artefakt#"$KORZEN"/}"
done
# Rdzeń ma być binarką Windows, nie linuksową o mylącej nazwie.
file -b "$POWLOKA/zasoby/rdzen/danaco-console.exe" | grep -q 'PE32+ executable' \
  || padnij "zasoby/rdzen/danaco-console.exe nie jest binarką Windows (PE32+) — to nie jest rdzeń dla tego celu"

zglos "Dokumentacja produktu"
# Instalator ma nieść dokumentację, bo to ON jest produktem końcowym: Operator
# dostaje jeden plik i po instalacji nie ma skąd wziąć instrukcji ani licencji.
# Dotąd nie niósł jej wcale — `scripts/pakowanie.sh` zakłada, że pliki *.md już
# leżą w katalogu wyjściowym (i dlatego ich nie rusza), a to jest prawdą tylko
# przy ręcznym składaniu do C:/DanacoConsole_App.
#
# Katalog zasobu jest WYTWOREM budowania, nie źródłem: pliki mają jedno miejsce
# w korzeniu repozytorium i są tam poprawiane. Kopia powstaje tuż przed
# złożeniem, żeby instalka nie wiozła wersji z poprzedniego tygodnia.
DOKUMENTACJA="$POWLOKA/zasoby/dokumentacja"
rm -rf "$DOKUMENTACJA"
mkdir -p "$DOKUMENTACJA"
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  [ -f "$KORZEN/$dokument" ] || padnij "brak dokumentu produktu: $dokument — instalator bez licencji i instrukcji nie jest produktem końcowym"
  cp -f "$KORZEN/$dokument" "$DOKUMENTACJA/$dokument"
  printf '  spakowany: %s (%s)\n' "$dokument" "$(du -h "$KORZEN/$dokument" | cut -f1)"
done

zglos "Budowa powłoki dla $CEL"
( cd "$POWLOKA" && cargo build --release --target "$CEL" )

zglos "Złożenie instalatora NSIS"
# `--bundles nsis` podane jawnie, choć tauri.conf.json ma ten cel ustawiony:
# skrypt składa instalator NSIS niezależnie od tego, jakie inne cele
# konfiguracja produktu wymienia.
( cd "$POWLOKA" && cargo tauri bundle --target "$CEL" --bundles nsis )

WYNIK="$(find "$POWLOKA/target/$CEL/release/bundle/nsis" -maxdepth 1 -name '*-setup.exe' -print -quit)"
[ -n "$WYNIK" ] || padnij "makensis nie zgłosił błędu, ale pliku instalatora nie ma"

zglos "Pomiar wyniku"
printf 'ścieżka : %s\n' "$WYNIK"
printf 'rozmiar : %s\n' "$(du -h "$WYNIK" | cut -f1)"
printf 'typ     : %s\n' "$(file -b "$WYNIK")"
printf 'suma    : %s\n' "$(sha256sum "$WYNIK" | cut -d' ' -f1)"
if command -v 7z >/dev/null; then
  # Ostatni wiersz listingu 7z kończy się słowem „files"; liczba stoi przed nim.
  printf 'wewnątrz: %s plików\n' "$(7z l "$WYNIK" 2>/dev/null | tail -1 | awk '{print $(NF-1)}')"
else
  printf 'wewnątrz: nie sprawdzono (brak 7z — apt: p7zip-full)\n'
fi

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i czy przechodzi do końca,
  - czy zainstalowana aplikacja wstaje (wymaga WebView2 w systemie),
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie trzy wymagają maszyny z Windows. Nie zakładaj ich powodzenia.
KONIEC
