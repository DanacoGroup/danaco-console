#!/usr/bin/env bash
# Instalka od źródeł do gotowego pliku — JEDNO polecenie, w kolejności, która nie
# pozwala wydać produktu niepełnego.
#
# ── Po co ten skrypt, skoro są już dwa ──────────────────────────────────────
# `scripts/wydanie.sh` składa układ ręczny w C:/DanacoConsole_App,
# `scripts/instalka-windows.sh` składa sam instalator NSIS z gotowych artefaktów.
# Między nimi była dziura, w którą wpadło 18.08.2026 wydanie produktu:
#
#   1. instalator nie niósł SERWERA NARZĘDZI (rdzeń szuka go obok siebie), więc
#      zainstalowana aplikacja wstawała bez sterowania platformą,
#   2. instalator nie niósł DOKUMENTACJI, bo skrypt pakowania z zamysłu nie rusza
#      plików *.md — co jest prawdą tylko przy ręcznym składaniu,
#   3. zasób `zasoby/rdzen/dist` był z poprzedniego dnia, a nikt tego nie mierzył:
#      instalka wyglądała świeżo i niosła interfejs sprzed doby.
#
# Pierwsze dwie naprawiono w konfiguracji produktu i w skrypcie NSIS. Trzecia jest
# pułapką KOLEJNOŚCI, a nie konfiguracji, i dlatego istnieje ten plik: buduje
# wszystko od nowa, a potem MIERZY wynik, zamiast założyć, że się udało.
#
# ── Czego pilnuje przed budową ──────────────────────────────────────────────
# Drzewo musi być czyste. Instalka z niezapisanego drzewa jest produktem, którego
# nie da się odtworzyć ze źródeł: suma kontrolna nie odpowiada żadnemu zapisowi
# i nikt po tygodniu nie powie, co w niej jest. Wymuszenie: DANACO_BRUDNE=1.
#
# Użycie:   bash budowa/scripts/instalka-pelna.sh
# Wynik:    budowa/wydania/<wersja>-<zapis>/ wraz z sumami SHA-256
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BUDOWA="$KORZEN/budowa"
POWLOKA="$BUDOWA/desktop/src-tauri"
CEL="x86_64-pc-windows-gnu"

# shellcheck disable=SC1091
[ -f "$HOME/.cargo/env" ] && . "$HOME/.cargo/env"

zglos() { printf '\n═══ %s\n' "$1"; }
padnij() { printf '\nODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Stan drzewa"
cd "$KORZEN"
BRUDNE="$(git status --porcelain | wc -l)"
if [ "$BRUDNE" -ne 0 ] && [ "${DANACO_BRUDNE:-0}" != "1" ]; then
  git status --porcelain | head -10
  padnij "$BRUDNE plików poza zapisem — instalki z niezapisanego drzewa nie da się
        odtworzyć ze źródeł. Zapisz robotę albo wymuś: DANACO_BRUDNE=1"
fi
ZAPIS="$(git rev-parse --short HEAD)"
WERSJA="$(python3 -c "import json;print(json.load(open('$POWLOKA/tauri.conf.json'))['version'])")"
printf '  zapis: %s   wersja produktu: %s\n' "$ZAPIS" "$WERSJA"
[ "$BRUDNE" -eq 0 ] || printf '  UWAGA: drzewo brudne, budowa wymuszona\n'

zglos "Sprawdziany przed wydaniem"
# Produkt wydany z czerwonym drzewem jest produktem, o którym nie wiadomo nic.
# Pełny przebieg rdzenia idzie ponad siedem minut i dlatego stoi TUTAJ, a nie
# w pamięci wykonawcy jako „przecież przechodziło".
( cd "$BUDOWA" && go build ./... && go vet ./server/... ) || padnij "rdzeń się nie składa"
( cd "$BUDOWA" && go test ./server/internal/core/ -count=1 >/dev/null ) || padnij "sprawdziany rdzenia niepomyślne"
( cd "$BUDOWA/client" && npx tsc --noEmit && npx vitest run >/dev/null ) || padnij "klient niepomyślny"
printf '  rdzeń, klient i sprawdziany: pomyślne\n'

zglos "Pakiet klienta"
( cd "$BUDOWA/client" && npm run build >/dev/null )
printf '  złożony: client/dist\n'

zglos "Rdzeń i serwer narzędzi dla $CEL"
# Serwer narzędzi buduje się TU, obok rdzenia, bo do produktu jedzie razem
# z nim — rdzeń szuka go we własnym katalogu (server/internal/narzedzia).
export CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64
mkdir -p "$POWLOKA/zasoby/rdzen"
( cd "$BUDOWA" && go build -o "$POWLOKA/zasoby/rdzen/danaco-console.exe" ./server/cmd/danaco-console )
( cd "$BUDOWA" && go build -o "$POWLOKA/zasoby/rdzen/danaco-narzedzia.exe" ./server/cmd/danaco-narzedzia )
unset CGO_ENABLED CC GOOS GOARCH
for binarium in danaco-console danaco-narzedzia; do
  file -b "$POWLOKA/zasoby/rdzen/$binarium.exe" | grep -q 'PE32+ executable' \
    || padnij "$binarium.exe nie jest binarką Windows — to nie jest produkt dla tego celu"
  printf '  zbudowany: %s.exe (%s)\n' "$binarium" "$(du -h "$POWLOKA/zasoby/rdzen/$binarium.exe" | cut -f1)"
done

zglos "Odświeżenie pakietu klienta w zasobach"
# `--delete`, nie samo kopiowanie: plik usunięty z interfejsu ma zniknąć i tutaj.
# Bez tego zasób rośnie o pliki, których produkt już nie zna.
rsync -a --delete "$BUDOWA/client/dist/" "$POWLOKA/zasoby/rdzen/dist/"
printf '  odświeżony: zasoby/rdzen/dist\n'

zglos "Instalator NSIS wraz z dokumentacją"
bash "$BUDOWA/scripts/instalka-windows.sh" >/tmp/instalka-nsis.log 2>&1 \
  || { tail -20 /tmp/instalka-nsis.log; padnij "budowa instalatora niepomyślna (dziennik: /tmp/instalka-nsis.log)"; }
WYNIK="$(find "$POWLOKA/target/$CEL/release/bundle/nsis" -maxdepth 1 -name '*-setup.exe' -print -quit)"
[ -n "$WYNIK" ] || padnij "pliku instalatora nie ma"

zglos "Pomiar wyniku — czy produkt jest pełny"
# Pomiar, nie założenie: sprawdzamy, że w archiwum SĄ te części, których brak
# przeszedł niezauważony 18.08.2026.
command -v 7z >/dev/null || padnij "brak 7z (apt: p7zip-full) — bez niego nie da się zmierzyć zawartości instalki"
SPIS="$(7z l "$WYNIK" 2>/dev/null)"
for czesc in \
  "rdzen/danaco-console.exe" \
  "rdzen/danaco-narzedzia.exe" \
  "rdzen/dist/index.html" \
  "dokumentacja/README.md" \
  "dokumentacja/LICENSE.md" \
  "dokumentacja/INSTALACJA-I-KONFIGURACJA.md" \
  "dokumentacja/INSTRUKCJA-UZYTKOWANIA.md"; do
  printf '%s' "$SPIS" | grep -qF "$czesc" || padnij "instalka nie niesie: $czesc"
  printf '  jest: %s\n' "$czesc"
done

zglos "Odłożenie wydania"
KATALOG="$BUDOWA/wydania/$WERSJA-$ZAPIS"
mkdir -p "$KATALOG"
cp -f "$WYNIK" "$KATALOG/"
( cd "$KATALOG" && sha256sum ./*.exe > SUMY-SHA256.txt )
printf '  %s\n' "$KATALOG"
printf '  rozmiar : %s\n' "$(du -h "$KATALOG"/*.exe | cut -f1)"
printf '  plików w archiwum: %s\n' "$(printf '%s' "$SPIS" | tail -1 | awk '{print $(NF-1)}')"
printf '  suma    : %s\n' "$(cut -d' ' -f1 "$KATALOG/SUMY-SHA256.txt")"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i przechodzi do końca,
  - czy zainstalowana aplikacja wstaje (wymaga WebView2 w systemie),
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie trzy wymagają maszyny z Windows.
  Paczki Linuksa (.deb, AppImage) to osobny cel — ten skrypt ich nie składa.
KONIEC
