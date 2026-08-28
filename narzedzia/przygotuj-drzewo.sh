#!/usr/bin/env bash
# Przygotowanie drzewa roboczego wiąże zależności symbolicznie z drzewa głównego zamiast kopiować je osobno, a wywołuje je narzedzia/nowy-teren.sh.

set -euo pipefail

katalog_drzewa=${1:-$(pwd)}
katalog_glowny=$(git -C "$katalog_drzewa" rev-parse --path-format=absolute --git-common-dir)
katalog_glowny=$(dirname "$katalog_glowny")

if [ "$katalog_drzewa" = "$katalog_glowny" ]; then
  echo "To jest drzewo główne — nie wymaga przygotowania."
  exit 0
fi

# Katalogi z zależnościami Node dowiązywane są tam, gdzie stoi manifest package.json, oddzielnie dla każdego pakietu wewnątrz drzewa roboczego.
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
