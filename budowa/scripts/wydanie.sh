#!/usr/bin/env bash
# Wydanie Danaco Console — jedno polecenie od źródeł do uruchamialnego produktu.
#
# Buduje cztery części produktu i przekazuje je skryptowi pakowania:
#   1. pakiet klienta   — npm ci (przy braku node_modules) + tsc + vite build,
#   2. rdzeń Go         — go build ./server/cmd/danaco-console,
#   3. serwer narzędzi  — go build ./server/cmd/danaco-narzedzia,
#   4. powłokę Tauri    — cargo tauri build --no-bundle (pomijalna, patrz niżej).
#
# Serwer narzędzi nie jest pomijalny tak jak powłoka. Rdzeń dokłada wpis
# `danaco` do konfiguracji MCP każdej tury i wskazuje nim plik danaco-narzedzia
# leżący obok siebie (server/internal/narzedzia/wpiecie.go); bez tego pliku
# narzędzia sterowania platformą nie działają.
#
# Struktura wyjścia (składana przez scripts/pakowanie.sh):
#   C:/DanacoConsole_App/
#   ├── danaco-console.exe    ← rdzeń (rola all, nasłuch 17870); serwuje interfejs
#   ├── danaco-narzedzia.exe  ← serwer narzędzi modelu (MCP po stdio); rdzeń szuka go obok siebie
#   ├── Danaco Console.exe    ← powłoka natywna (obecna, gdy budowa powłoki się powiodła)
#   ├── client/dist/          ← pakiet interfejsu; rdzeń znajduje go obok siebie
#   └── *.md                  ← dokumentacja produktu — wydanie jej nie rusza
#
# Układ odpowiada wyszukiwaniu w powłoce (desktop/src-tauri/src/rdzen/lokalizacja.rs
# szuka danaco-console.exe obok pliku powłoki; rdzen/pakiet_klienta.rs szuka
# client/dist obok binarki rdzenia) oraz README.md (C:\DanacoConsole_App).
#
# Użycie:   bash budowa/scripts/wydanie.sh [katalog-wyjścia]
#           domyślnie C:/DanacoConsole_App
# Zmienne:  DANACO_BEZ_POWLOKI=1  — pomija budowę powłoki (wariant rdzeń+klient);
#           powłoka jest też pomijana automatycznie, gdy brak cargo albo tauri-cli.
# Instalator NSIS to osobny krok (wymaga NSIS): cargo tauri build w desktop/src-tauri.

set -euo pipefail

SKRYPTY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUDOWA="$(dirname "$SKRYPTY")"
WYJSCIE="${1:-C:/DanacoConsole_App}"
STAGING="${TEMP:-/tmp}/danaco-wydanie"
mkdir -p "$STAGING"

echo "══ 1. Zależności klienta"
cd "$BUDOWA/client"
if [ ! -d node_modules ]; then
  npm ci --no-audit --no-fund
else
  echo "  node_modules istnieje — instalacja pominięta (wymuszenie: usuń node_modules)"
fi

echo "══ 2. Pakiet klienta (tsc + vite build → client/dist)"
npm run build

echo "══ 3. Rdzeń (go build ./server/cmd/danaco-console)"
cd "$BUDOWA"
go build -o "$STAGING/danaco-console.exe" ./server/cmd/danaco-console
echo "  zbudowany: $STAGING/danaco-console.exe"

echo "══ 4. Serwer narzędzi modelu (go build ./server/cmd/danaco-narzedzia)"
go build -o "$STAGING/danaco-narzedzia.exe" ./server/cmd/danaco-narzedzia
echo "  zbudowany: $STAGING/danaco-narzedzia.exe"

echo "══ 5. Powłoka natywna (Tauri)"
POWLOKA=""
if [ "${DANACO_BEZ_POWLOKI:-0}" = "1" ]; then
  echo "  pominięta na żądanie (DANACO_BEZ_POWLOKI=1) — wariant rdzeń+klient"
elif ! command -v cargo >/dev/null 2>&1 || ! cargo tauri --version >/dev/null 2>&1; then
  echo "  pominięta: brak cargo albo tauri-cli — wariant rdzeń+klient"
  echo "  (instalacja narzędzia: cargo install tauri-cli)"
else
  # --no-bundle: sam plik wykonywalny; instalator NSIS wymaga narzędzia NSIS
  # i pozostaje osobnym krokiem wydania. beforeBuildCommand z tauri.conf.json
  # przebudowuje pakiet klienta — powtórka kroku 2 jest zamierzona i tania.
  (cd "$BUDOWA/desktop/src-tauri" && cargo tauri build --no-bundle)
  POWLOKA="$BUDOWA/desktop/src-tauri/target/release/danaco-console-powloka.exe"
  if [ ! -f "$POWLOKA" ]; then
    echo "BŁĄD: budowa powłoki zgłosiła powodzenie, ale brak pliku: $POWLOKA" >&2
    exit 1
  fi
  echo "  zbudowana: $POWLOKA"
fi

echo "══ 6. Pakowanie → $WYJSCIE"
bash "$SKRYPTY/pakowanie.sh" "$WYJSCIE" "$STAGING/danaco-console.exe" \
  "$STAGING/danaco-narzedzia.exe" "$BUDOWA/client/dist" ${POWLOKA:+"$POWLOKA"}

echo
echo "══ Wydanie zamknięte: $WYJSCIE"
if [ -z "$POWLOKA" ]; then
  echo "  wariant rdzeń+klient — uruchomienie: $WYJSCIE/danaco-console.exe,"
  echo "  interfejs w przeglądarce pod http://127.0.0.1:17870/"
else
  echo "  uruchomienie: \"$WYJSCIE/Danaco Console.exe\" (powłoka stawia rdzeń w tle)"
fi
