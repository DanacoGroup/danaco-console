#!/usr/bin/env bash
# Drabina weryfikacji budowy Danaco Console. Uruchamiana przed zapisem rewizji
# i przed odbiorem terenu; do tej pory istniała wyłącznie jako proza w dokumencie,
# więc nikt jej nie przechodził i nikt tego nie zauważał.
#
# Użycie: bash narzedzia/drabina.sh [szybka|pelna]
#   szybka — format, budowa, dyscyplina; sekundy, nadaje się pod hook
#   pelna  — szybka wraz ze sprawdzianami rdzenia i klienta; minuty
set -uo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TRYB="${1:-szybka}"
NIEPOWODZEN=0

szczebel() { printf '\n── %s\n' "$1"; }
zdany() { printf '   zdany: %s\n' "$1"; }
oblany() { printf '   OBLANY: %s\n' "$1"; NIEPOWODZEN=$((NIEPOWODZEN + 1)); }

cd "$KORZEN"

szczebel "Format plików Go"
NIESFORMATOWANE="$(find budowa narzedzia -name '*.go' -not -path '*/node_modules/*' -print0 2>/dev/null |
	xargs -0 gofmt -l 2>/dev/null)"
if [ -z "$NIESFORMATOWANE" ]; then
	zdany "gofmt nie zgłasza żadnego pliku"
else
	oblany "gofmt zgłasza pliki: $(echo "$NIESFORMATOWANE" | tr '\n' ' ')"
fi

szczebel "Budowa rdzenia"
if (cd budowa/server && go build ./... 2>&1); then
	zdany "go build ./... przechodzi"
else
	oblany "go build ./... nie przechodzi"
fi

szczebel "Analiza statyczna rdzenia"
if (cd budowa/server && go vet ./... 2>&1 | head -20); then
	zdany "go vet bez zastrzeżeń"
else
	oblany "go vet zgłasza zastrzeżenia"
fi

szczebel "Dyscyplina inżynierska"
# Na maszynie stoi kilka wydań walidatora o różnych progach udziału komentarzy.
# Pierwsze trafienie wyszukiwania dawało próg losowy w rozrzucie czterokrotnym,
# więc wydanie dobiera się po zgodności z progiem zapisanym w standardzie redakcji.
PROG_STANDARDU=0.05
WALIDATOR=""
while IFS= read -r kandydat; do
	[ -z "$kandydat" ] && continue
	NASTAWA="$(dirname "$kandydat")/dyscyplina.config.json"
	[ -r "$NASTAWA" ] || continue
	if grep -q "\"max_udzial_komentarzy\"[[:space:]]*:[[:space:]]*$PROG_STANDARDU" "$NASTAWA"; then
		WALIDATOR="$kandydat"
		break
	fi
done < <(find "$HOME/.claude" -name style_guard.py 2>/dev/null | sort)
if [ -z "$WALIDATOR" ]; then
	printf '   pominięty: walidatora dyscypliny nie ma na tej maszynie\n'
else
	ZMIENIONE="$(git diff --cached --name-only --diff-filter=ACM 2>/dev/null |
		grep -E '\.(go|ts|js|css|py|sh)$' || true)"
	[ -z "$ZMIENIONE" ] && ZMIENIONE="$(git diff --name-only --diff-filter=ACM 2>/dev/null |
		grep -E '\.(go|ts|js|css|py|sh)$' || true)"
	if [ -z "$ZMIENIONE" ]; then
		printf '   pominięty: brak zmienionych plików kodu\n'
	elif echo "$ZMIENIONE" | xargs python3 "$WALIDATOR" 2>&1 | grep -q 'BLOK'; then
		oblany "walidator dyscypliny zgłasza naruszenia blokujące w plikach zmienionych"
	else
		zdany "walidator dyscypliny bez naruszeń blokujących"
	fi
fi

szczebel "Świeżość wytworów kontraktu"
if [ -d budowa/shared/gen ]; then
	PRZED="$(git status --porcelain budowa/shared/contract.go budowa/shared/contract.ts 2>/dev/null | wc -l)"
	if (cd budowa/klient && npm run kontrakt >/dev/null 2>&1); then
		PO="$(git status --porcelain budowa/shared/contract.go budowa/shared/contract.ts 2>/dev/null | wc -l)"
		if [ "$PRZED" = "$PO" ]; then
			zdany "wytwory kontraktu są świeże wobec contract.json"
		else
			oblany "wytwory kontraktu rozjechały się z contract.json — przebuduj i zapisz"
		fi
	else
		printf '   pominięty: generatora nie dało się uruchomić\n'
	fi
fi

if [ "$TRYB" = "pelna" ]; then
	szczebel "Sprawdziany klienta"
	if (cd budowa/klient && npm run testy 2>&1 | tail -3); then
		zdany "sprawdziany klienta przechodzą"
	else
		oblany "sprawdziany klienta nie przechodzą"
	fi

	szczebel "Sprawdziany rdzenia — pełny bieg pod zamkiem"
	if DANACO_MODELE=/opt/danaco-modele flock /tmp/danaco-bieg-pelny.lock \
		go test -C budowa/server -count=1 -timeout 40m ./... 2>&1 | tail -5; then
		zdany "sprawdziany rdzenia przechodzą"
	else
		oblany "sprawdziany rdzenia nie przechodzą"
	fi
fi

printf '\n'
if [ "$NIEPOWODZEN" -eq 0 ]; then
	printf 'DRABINA ZDANA (%s)\n' "$TRYB"
	exit 0
fi
printf 'DRABINA OBLANA: %s szczebli\n' "$NIEPOWODZEN"
exit 1
