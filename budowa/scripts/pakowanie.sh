#!/usr/bin/env bash
# Pakowanie Danaco Console — złożenie katalogu produktu z gotowych artefaktów.
#
# Jedyna odpowiedzialność: przenieść zbudowane części do katalogu wyjściowego.
# Niczego nie buduje — budową zajmuje się scripts/wydanie.sh, który wywołuje ten
# skrypt jako ostatni.
#
# Katalog wyjściowy zawiera równolegle dokumentację produktu (*.md), więc skrypt
# nie usuwa ani nie nadpisuje plików *.md; czyszczenie obejmuje wyłącznie
# podkatalog client/ i to z wyłączeniem *.md.
#
# Serwer narzędzi jest częścią obowiązkową: rdzeń wpisuje binarium
# danaco-narzedzia do konfiguracji MCP każdej tury i szuka go obok siebie
# (server/internal/narzedzia/wpiecie.go, sciezkaProgramu). Bez tego pliku model
# nie dostaje ani jednego narzędzia sterowania platformą.
#
# Użycie: bash budowa/scripts/pakowanie.sh WYJŚCIE RDZEŃ_EXE NARZĘDZIA_EXE KLIENT_DIST [POWŁOKA_EXE]

set -euo pipefail

if [ $# -lt 4 ]; then
  echo "Użycie: pakowanie.sh WYJŚCIE RDZEŃ_EXE NARZĘDZIA_EXE KLIENT_DIST [POWŁOKA_EXE]" >&2
  exit 2
fi

WYJSCIE="$1"
RDZEN="$2"
NARZEDZIA="$3"
KLIENT="$4"
POWLOKA="${5:-}"

[ -f "$RDZEN" ] || { echo "BŁĄD: brak binarki rdzenia: $RDZEN" >&2; exit 1; }
[ -f "$NARZEDZIA" ] || { echo "BŁĄD: brak binarki serwera narzędzi: $NARZEDZIA" >&2; exit 1; }
[ -f "$KLIENT/index.html" ] || { echo "BŁĄD: brak pakietu klienta (index.html): $KLIENT" >&2; exit 1; }

mkdir -p "$WYJSCIE"

# Rdzeń — obok niego powłoka i pakiet klienta (rdzen/lokalizacja.rs, pakiet_klienta.rs).
cp -f "$RDZEN" "$WYJSCIE/danaco-console.exe"
echo "  rdzeń → $WYJSCIE/danaco-console.exe"

# Serwer narzędzi modelu — obowiązkowo obok rdzenia (wpiecie.go: sciezkaProgramu
# składa ścieżkę z katalogu bieżącego procesu). Bez niego model nie ma ani
# jednego narzędzia sterowania platformą.
cp -f "$NARZEDZIA" "$WYJSCIE/danaco-narzedzia.exe"
echo "  serwer narzędzi → $WYJSCIE/danaco-narzedzia.exe"

# Pakiet klienta: stare pliki pakietu ustępują nowym (nazwy z sumami się zmieniają),
# ale *.md pozostają nietknięte — mogłaby tam leżeć dokumentacja.
if [ -d "$WYJSCIE/client" ]; then
  find "$WYJSCIE/client" -type f ! -name '*.md' -delete
  find "$WYJSCIE/client" -mindepth 1 -type d -empty -delete
fi
mkdir -p "$WYJSCIE/client/dist"
cp -R "$KLIENT/." "$WYJSCIE/client/dist/"
echo "  klient → $WYJSCIE/client/dist ($(find "$WYJSCIE/client/dist" -type f | wc -l) plików)"

# Powłoka — opcjonalna: wydanie.sh pomija ją przy braku narzędzi (wariant rdzeń+klient).
if [ -n "$POWLOKA" ]; then
  [ -f "$POWLOKA" ] || { echo "BŁĄD: wskazana powłoka nie istnieje: $POWLOKA" >&2; exit 1; }
  cp -f "$POWLOKA" "$WYJSCIE/Danaco Console.exe"
  echo "  powłoka → $WYJSCIE/Danaco Console.exe"
fi
