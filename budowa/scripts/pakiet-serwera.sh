#!/usr/bin/env bash
# Pakiet serwera Danaco Console jest złożeniem wydania dla maszyny serwera:
# całość platformy stoi tam, u operatora zostaje cienkie okno, dlatego pakiet
# idzie przez dpkg-deb, nie przez bundler Tauri.
set -euo pipefail

SKRYPTY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUDOWA="$(dirname "$SKRYPTY")"
KORZEN="$(dirname "$BUDOWA")"
PAKOWANIE="$BUDOWA/packaging"
SZKIELET="$PAKOWANIE/drzewo"
ROBOCZY="$PAKOWANIE/roboczy"

# Wersji nie jest podnoszona, bo żadne wydanie nie było jeszcze w użytku; numer
# wersji zostaje pod kontrolą operatora wydania.
WERSJA="${DANACO_WERSJA:-2.0.0}"
DATA="$(date +%Y-%m-%d)"
KATALOG_WYDANIA="${DANACO_KATALOG_WYDANIA:-$BUDOWA/wydania/$WERSJA-$DATA}"

BEZ_BUDOWY=nie
[ "${1:-}" = "--bez-budowy" ] && BEZ_BUDOWY=tak

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() {
	printf 'ODMOWA: %s\n' "$1" >&2
	exit 1
}

for program in dpkg-deb dpkg go; do
	command -v "$program" >/dev/null || padnij "brak programu $program — bez niego pakietu nie złożę"
done

ARCHITEKTURA="$(dpkg --print-architecture)"

# ── Port nasłuchu ────────────────────────────────────────────────────────────
zglos "port nasłuchu rdzenia"
# Port stoi w pakiecie w dwóch miejscach: w jednostce systemd, która go ustawia,
# i w opisie z DEBIAN/control, który go zapowiada administratorowi. Rozejście
# tych dwóch wartości wysyła administratora pod port, którego rdzeń nie ma
# otwartego, a widać je dopiero po założeniu pakietu na maszynie — dlatego
# zgodność sprawdzana jest tutaj, przed budową.
JEDNOSTKA="$SZKIELET/lib/systemd/system/danaco-console.service"
OPIS="$SZKIELET/DEBIAN/control"
PORT_JEDNOSTKI="$(sed -n 's/^Environment=DANACO_PORT=\([0-9]\+\)$/\1/p' "$JEDNOSTKA")"
PORT_OPISU="$(sed -n 's/.*nasłuchuje na porcie \([0-9]\+\).*/\1/p' "$OPIS")"
[ -n "$PORT_JEDNOSTKI" ] || padnij "jednostka systemd nie niesie wiersza Environment=DANACO_PORT=<port> — pakiet nie wie, na czym ma nasłuchiwać"
[ -n "$PORT_OPISU" ] || padnij "opis w DEBIAN/control nie nazywa portu nasłuchu — administrator nie ma skąd wziąć tej wartości"
[ "$PORT_JEDNOSTKI" = "$PORT_OPISU" ] ||
	padnij "port rozjechał się w pakiecie: jednostka systemd stawia rdzeń na $PORT_JEDNOSTKI, a opis w DEBIAN/control zapowiada $PORT_OPISU"
printf '  jednostka systemd i opis pakietu zgodnie wskazują port %s\n' "$PORT_JEDNOSTKI"

# ── Budowa składników ────────────────────────────────────────────────────────
# Rdzeń idzie natywnie na Linuksa: bez krzyżowej kompilacji, bo pakiet .deb jest
# dla tej samej rodziny maszyn, na której stoi to drzewo.
if [ "$BEZ_BUDOWY" = "nie" ]; then
	zglos "budowa rdzenia i serwera narzędzi ($ARCHITEKTURA)"
	# Hasło skrzynki nadawczej wchodzi przy składaniu, nie ze źródła: repozytorium
	# go nie niesie, więc kopia drzewa nie daje dostępu do poczty platformy.
	# Pakiet złożony bez tej zmiennej działa — rdzeń bierze hasło ze środowiska
	# maszyny wdrożenia, a bez niego wchodzi w drogę bez poczty.
	WPIS_SEKRETU=""
	if [ -n "${DANACO_NADAWCA_SEKRET:-}" ]; then
		WPIS_SEKRETU="-X danacoconsole/server/internal/konfiguracja.NadawcaSekretWbudowany=$DANACO_NADAWCA_SEKRET"
		printf '  hasło skrzynki nadawczej: wpisane w pakiet\n'
	else
		printf '  hasło skrzynki nadawczej: BRAK — pakiet weźmie je ze środowiska maszyny\n'
	fi
	(cd "$BUDOWA" && go build -ldflags "$WPIS_SEKRETU" -o "$PAKOWANIE/danaco-console" ./server/cmd/danaco-console) ||
		padnij "rdzeń się nie zbudował"
	(cd "$BUDOWA" && go build -o "$PAKOWANIE/danaco-narzedzia" ./server/cmd/danaco-narzedzia) ||
		padnij "serwer narzędzi się nie zbudował"

	zglos "budowa pakietu interfejsu"
	(cd "$KORZEN" && npm --prefix budowa/klient run budowanie) ||
		padnij "pakiet interfejsu się nie zbudował"
else
	zglos "budowa pominięta (--bez-budowy)"
fi

[ -x "$PAKOWANIE/danaco-console" ] || padnij "nie ma binarium rdzenia: $PAKOWANIE/danaco-console"
[ -x "$PAKOWANIE/danaco-narzedzia" ] || padnij "nie ma serwera narzędzi: $PAKOWANIE/danaco-narzedzia"
[ -d "$BUDOWA/klient/dist" ] || padnij "nie ma pakietu interfejsu: $BUDOWA/klient/dist"

# ── Drzewo pakietu ───────────────────────────────────────────────────────────
# Drzewo roboczne jest skladane od zera przy każdym przebiegu: pakiet ma nieść to, co
# zbudowane teraz, a nie resztki poprzedniego przebiegu.
zglos "składanie drzewa pakietu"
rm -rf "$ROBOCZY"
mkdir -p "$ROBOCZY"
cp -a "$SZKIELET/." "$ROBOCZY/"

mkdir -p "$ROBOCZY/opt/danaco-console/scripts" \
	"$ROBOCZY/usr/share/doc/danaco-console" \
	"$ROBOCZY/usr/bin"

# Oba pliki wykonywalne stoją w jednym katalogu — inaczej rdzeń nie znajdzie
# serwera narzędzi obok siebie na dysku.
install -m 0755 "$PAKOWANIE/danaco-console" "$ROBOCZY/opt/danaco-console/danaco-console"
install -m 0755 "$PAKOWANIE/danaco-narzedzia" "$ROBOCZY/opt/danaco-console/danaco-narzedzia"

# Pakiet interfejsu trafia dokładnie tam, gdzie jednostka systemd wskazuje
# rdzeniowi jego katalog interfejsu.
mkdir -p "$ROBOCZY/opt/danaco-console/client"
cp -a "$BUDOWA/klient/dist" "$ROBOCZY/opt/danaco-console/client/dist"
find "$ROBOCZY/opt/danaco-console/client" -type d -exec chmod 0755 {} +
find "$ROBOCZY/opt/danaco-console/client" -type f -exec chmod 0644 {} +

# Rdzeń na ścieżce systemu. Dowiązanie, nie kopia: os.Executable rozwija
# dowiązanie, więc serwer narzędzi zostaje znaleziony obok pliku właściwego.
ln -sf /opt/danaco-console/danaco-console "$ROBOCZY/usr/bin/danaco-console"

# Skrypt arsenału jedzie z pakietem, bo stawia go się na SERWERZE, a nie w drzewie
# budowy. Treści się nie powtarza — brany jest plik taki, jaki jest.
if [ -r "$SKRYPTY/arsenal-serwera.sh" ]; then
	install -m 0755 "$SKRYPTY/arsenal-serwera.sh" \
		"$ROBOCZY/opt/danaco-console/scripts/arsenal-serwera.sh"
else
	printf 'UWAGA: brak scripts/arsenal-serwera.sh — pakiet pójdzie bez skryptu arsenału\n' >&2
fi

# Pomocniki pythonowe wołane przez rdzeń; rdzeń szuka ich obok siebie, tak
# samo jak serwera narzędzi, a prowizjonowanie stamtąd bierze plik wymagań
# środowiska rozpoznawania mowy.
if [ -d "$BUDOWA/pomocniki" ]; then
	cp -a "$BUDOWA/pomocniki" "$ROBOCZY/opt/danaco-console/pomocniki"
	find "$ROBOCZY/opt/danaco-console/pomocniki" -type d -exec chmod 0755 {} +
	find "$ROBOCZY/opt/danaco-console/pomocniki" -type f -exec chmod 0644 {} +
else
	padnij "nie ma katalogu pomocników: $BUDOWA/pomocniki — bez niego pakiet nie postawi mowy"
fi

# Cztery dokumenty produktu trafiają do pakietu jako dokumentacja instalowana
# wraz z usługą. Stoją w docs/, bo tam prowadzona jest cała dokumentacja produktu;
# kopia w korzeniu byłaby drugą prawdą o tym samym pliku.
BRAK_DOKUMENTU=nie
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
	if [ -r "$KORZEN/docs/$dokument" ]; then
		install -m 0644 "$KORZEN/docs/$dokument" "$ROBOCZY/usr/share/doc/danaco-console/$dokument"
	else
		printf 'UWAGA: brak dokumentu %s w %s/docs\n' "$dokument" "$KORZEN" >&2
		BRAK_DOKUMENTU=tak
	fi
done
[ "$BRAK_DOKUMENTU" = "nie" ] || padnij "pakiet ma nieść cztery dokumenty produktu — któregoś nie ma"

# ── Kontrola ─────────────────────────────────────────────────────────────────
# Architekturę i rozmiar wpisuje skrypt: zmierzone bije przepisane.
ROZMIAR_KB="$(du -sk --exclude=DEBIAN "$ROBOCZY" | cut -f1)"
sed -i \
	-e "s|@ARCHITEKTURA@|$ARCHITEKTURA|" \
	-e "s|@ROZMIAR_KB@|$ROZMIAR_KB|" \
	-e "s|^Version: .*|Version: $WERSJA|" \
	"$ROBOCZY/DEBIAN/control"

chmod 0755 "$ROBOCZY/DEBIAN/postinst" "$ROBOCZY/DEBIAN/prerm" "$ROBOCZY/DEBIAN/postrm"
chmod 0644 "$ROBOCZY/DEBIAN/control" "$ROBOCZY/DEBIAN/conffiles"
chmod 0644 "$ROBOCZY/lib/systemd/system/danaco-console.service"
chmod 0640 "$ROBOCZY/etc/danaco-console/srodowisko"
find "$ROBOCZY" -path "$ROBOCZY/DEBIAN" -prune -o -type d -exec chmod 0755 {} +

# Sumy kontrolne zawartości — dpkg sprawdza nimi, czy plik na dysku nie odjechał
# od pliku z pakietu (dpkg --verify).
(
	cd "$ROBOCZY" &&
		find . -path ./DEBIAN -prune -o -type f -printf '%P\0' |
		xargs -0 --no-run-if-empty md5sum >DEBIAN/md5sums
)
chmod 0644 "$ROBOCZY/DEBIAN/md5sums"

# Złożenie pakietu odkłada gotowy plik do katalogu wydania i liczy jego sumę
# kontrolną, gotową do odbioru.
zglos "złożenie pakietu"
mkdir -p "$KATALOG_WYDANIA"
PAKIET="$KATALOG_WYDANIA/danaco-console_${WERSJA}_${ARCHITEKTURA}.deb"
# --root-owner-group: pliki w pakiecie należą do root, choć drzewo składał
# użytkownik zwykły. Bez tego dpkg rozłożyłby je na cudze konto.
dpkg-deb --root-owner-group --build "$ROBOCZY" "$PAKIET" >/dev/null ||
	padnij "dpkg-deb nie złożył pakietu"

SUMA="$(sha256sum "$PAKIET" | cut -d' ' -f1)"
printf '%s  %s\n' "$SUMA" "$(basename "$PAKIET")" >"$PAKIET.sha256"

zglos "wynik"
printf 'pakiet:  %s\n' "$PAKIET"
printf 'rozmiar: %s\n' "$(du -h "$PAKIET" | cut -f1)"
printf 'SHA-256: %s\n' "$SUMA"
printf 'nasłuch: port %s na wszystkich interfejsach\n' "$PORT_JEDNOSTKI"
printf '\nodbiór:\n  dpkg-deb -I %s\n  dpkg-deb -c %s\n' "$PAKIET" "$PAKIET"
