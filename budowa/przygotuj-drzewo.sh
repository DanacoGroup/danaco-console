#!/usr/bin/env bash
# Dowiązanie zależności klienta w równoległym drzewie roboczym gita.
#
# `node_modules` nie jest w repozytorium, więc `git worktree add` tworzy drzewo
# bez zależności i kompilator TypeScriptu nie ma czym ruszyć. Skrypt wskazuje
# katalog zależności z drzewa głównego, zamiast pobierać je drugi raz.
#
# W drzewie głównym nie ma nic do zrobienia — skrypt to rozpoznaje i kończy.
# Wolno go uruchamiać wielokrotnie.

set -eu

MOJE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GLOWNE="$(git -C "$MOJE" worktree list --porcelain | awk '/^worktree /{print $2; exit}')/budowa"

if [ "$MOJE" = "$GLOWNE" ]; then
  echo "To jest drzewo główne — zależności są na miejscu."
  exit 0
fi

ZRODLO="$GLOWNE/client/node_modules"
CEL="$MOJE/client/node_modules"

if [ -e "$CEL" ]; then
  echo "Zależności klienta już są."
elif [ -d "$ZRODLO" ]; then
  ln -sfn "$ZRODLO" "$CEL"
  echo "Zależności klienta dowiązane z drzewa głównego."
else
  echo "Drzewo główne nie ma zależności klienta." >&2
  echo "Pobierz je tam poleceniem: cd $GLOWNE/client && npm ci" >&2
  exit 1
fi
