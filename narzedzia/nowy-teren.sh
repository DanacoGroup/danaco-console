#!/usr/bin/env bash
# Otwarcie terenu tworzy nowe drzewo robocze Git na osobnej gałęzi, dzielące ten sam katalog repozytorium, aby równoległe sesje nie nadpisywały sobie plików.

set -euo pipefail

if [ $# -ne 2 ]; then
  echo "Użycie: bash narzedzia/nowy-teren.sh <nazwa> <galaz-bazowa>" >&2
  echo "Przykład: bash narzedzia/nowy-teren.sh zetony main" >&2
  exit 1
fi

nazwa=$1
galaz_bazowa=$2
katalog_glowny=$(git rev-parse --show-toplevel)
katalog_terenu="$HOME/robocze/$nazwa"

if ! git show-ref --verify --quiet "refs/heads/$galaz_bazowa"; then
  echo "Nie ma gałęzi bazowej $galaz_bazowa." >&2
  exit 1
fi

if [ -e "$katalog_terenu" ]; then
  echo "Katalog $katalog_terenu już istnieje — teren o tej nazwie jest otwarty." >&2
  exit 1
fi

git worktree add "$katalog_terenu" -b "teren/$nazwa" "$galaz_bazowa"

bash "$katalog_glowny/narzedzia/przygotuj-drzewo.sh" "$katalog_terenu"

echo
echo "Teren $nazwa otwarty na gałęzi teren/$nazwa (baza: $galaz_bazowa)."
echo "Drzewo robocze: $katalog_terenu"
echo "Wpis w prowadzenie/rejestr-terenow.md zakłada Prowadzący budowę."
