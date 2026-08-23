#!/usr/bin/env bash
# Przygotowanie drzewa roboczego do pracy: dowiązanie zależności zamiast kopii.
#
# Instalacja zależności osobno w każdym drzewie roboczym zajmowałaby setki
# megabajtów i minuty na teren. Drzewa dzielą jeden katalog node_modules przez
# dowiązanie symboliczne do drzewa głównego; pamięci podręczne Go i Rust są
# wspólne z ustawienia środowiska, więc nie wymagają tu niczego.
#
# Wywoływany przez narzedzia/nowy-teren.sh; wolno uruchomić samodzielnie.
# Użycie: bash narzedzia/przygotuj-drzewo.sh [katalog-drzewa]

set -euo pipefail

katalog_drzewa=${1:-$(pwd)}
katalog_glowny=$(git -C "$katalog_drzewa" rev-parse --path-format=absolute --git-common-dir)
katalog_glowny=$(dirname "$katalog_glowny")

if [ "$katalog_drzewa" = "$katalog_glowny" ]; then
  echo "To jest drzewo główne — nie wymaga przygotowania."
  exit 0
fi

# Katalogi z zależnościami Node dowiązywane są tam, gdzie stoi package.json.
podpiete=0
while IFS= read -r manifest; do
  katalog_pakietu=$(dirname "$manifest")
  wzgledny=${katalog_pakietu#"$katalog_drzewa"/}
  zrodlo="$katalog_glowny/$wzgledny/node_modules"
  cel="$katalog_pakietu/node_modules"
  if [ -d "$zrodlo" ] && [ ! -e "$cel" ]; then
    ln -s "$zrodlo" "$cel"
    echo "Dowiązano node_modules: $wzgledny"
    podpiete=$((podpiete + 1))
  fi
done < <(find "$katalog_drzewa" -name package.json -not -path "*/node_modules/*")

if [ "$podpiete" -eq 0 ]; then
  echo "Nie ma czego dowiązywać — drzewo nie zawiera zależności Node."
fi
