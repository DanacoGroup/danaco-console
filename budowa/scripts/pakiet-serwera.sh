#!/usr/bin/env bash
# Pakiet serwera Danaco Console (.deb) — złożenie wydania dla maszyny SERWERA.
#
# ── Co ten skrypt pakuje i dlaczego akurat to ─────────────────────────────────
# Model wdrożenia (rozstrzygnięcie zamknięte): całość stoi na serwerze, a
# u Operatora zostaje cienkie okno. Ten pakiet jest więc SERCEM platformy, nie
# oknem klienta — i dlatego nie idzie bundlerem Tauri, a dpkg-deb po własnym
# drzewie katalogów.
#
# Niesie:
#   danaco-console      — rdzeń, zbudowany natywnie na Linuksa,
#   danaco-narzedzia    — serwer narzędzi modelu; rdzeń szuka go OBOK SIEBIE
#                         (server/internal/narzedzia/wpiecie.go: os.Executable
#                         → filepath.Dir → nazwa), stąd oba pliki w jednym
#                         katalogu i stąd katalog roboczy jednostki systemd,
#   client/dist         — pakiet interfejsu, który rdzeń serwuje klientom,
#   cztery dokumenty    — README, instalacja i konfiguracja, instrukcja, licencja,
#   jednostkę systemd   — konto usługi, katalog danych, restart, port 17870,
#   scripts/arsenal-serwera.sh — skrypt prowizjonowania arsenału na serwerze.
#
# ── Czego ten skrypt NIE robi ─────────────────────────────────────────────────
# Nie stawia arsenału i nie przepisuje jego wykazu. Wykaz zależności rdzenia stoi
# w jednym rejestrze deklaracji (server/internal/core/zaleznosci_zewnetrzne.go
# wraz z deklaracjami w adapterach) i rdzeń wypisuje go sam
# (danaco-console --wykaz-zaleznosci). Pole Depends w packaging/drzewo/DEBIAN/control
# niesie tę część rejestru, którą ma dystrybucja; resztą zajmuje się
# scripts/arsenal-serwera.sh, wskazany Operatorowi przez postinst.
#
# ── Wywołanie ─────────────────────────────────────────────────────────────────
#   scripts/pakiet-serwera.sh                 # zbuduj wszystko i złóż pakiet
#   scripts/pakiet-serwera.sh --bez-budowy    # złóż z tego, co już zbudowane
#
# Zmienne:
#   DANACO_WERSJA=1.0.0        — wersja pakietu (domyślnie z DEBIAN/control)
#   DANACO_KATALOG_WYDANIA=…   — gdzie odłożyć wynik
set -euo pipefail

SKRYPTY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUDOWA="$(dirname "$SKRYPTY")"
KORZEN="$(dirname "$BUDOWA")"
PAKOWANIE="$BUDOWA/packaging"
SZKIELET="$PAKOWANIE/drzewo"
ROBOCZY="$PAKOWANIE/roboczy"

# Wersji NIE podnosimy — żadne wydanie nie było w użytku.
WERSJA="${DANACO_WERSJA:-1.0.0}"
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

# ── Budowa składników ────────────────────────────────────────────────────────
# Rdzeń idzie natywnie na Linuksa: bez krzyżowej kompilacji, bo pakiet .deb jest
# dla tej samej rodziny maszyn, na której stoi to drzewo.
if [ "$BEZ_BUDOWY" = "nie" ]; then
	zglos "budowa rdzenia i serwera narzędzi ($ARCHITEKTURA)"
	(cd "$BUDOWA" && go build -o "$PAKOWANIE/danaco-console" ./server/cmd/danaco-console) ||
		padnij "rdzeń się nie zbudował"
	(cd "$BUDOWA" && go build -o "$PAKOWANIE/danaco-narzedzia" ./server/cmd/danaco-narzedzia) ||
		padnij "serwer narzędzi się nie zbudował"

	zglos "budowa pakietu interfejsu"
	(cd "$KORZEN" && npm --prefix budowa/client run build) ||
		padnij "pakiet interfejsu się nie zbudował"
else
	zglos "budowa pominięta (--bez-budowy)"
fi

[ -x "$PAKOWANIE/danaco-console" ] || padnij "nie ma binarium rdzenia: $PAKOWANIE/danaco-console"
[ -x "$PAKOWANIE/danaco-narzedzia" ] || padnij "nie ma serwera narzędzi: $PAKOWANIE/danaco-narzedzia"
[ -d "$BUDOWA/client/dist" ] || padnij "nie ma pakietu interfejsu: $BUDOWA/client/dist"

# ── Drzewo pakietu ───────────────────────────────────────────────────────────
# Drzewo roboczne składamy od zera przy każdym przebiegu: pakiet ma nieść to, co
# zbudowane teraz, a nie resztki poprzedniego przebiegu.
zglos "składanie drzewa pakietu"
rm -rf "$ROBOCZY"
mkdir -p "$ROBOCZY"
cp -a "$SZKIELET/." "$ROBOCZY/"

mkdir -p "$ROBOCZY/opt/danaco-console/scripts" \
	"$ROBOCZY/usr/share/doc/danaco-console" \
	"$ROBOCZY/usr/bin"

# Oba pliki wykonywalne w JEDNYM katalogu — inaczej rdzeń nie znajdzie serwera
# narzędzi obok siebie.
install -m 0755 "$PAKOWANIE/danaco-console" "$ROBOCZY/opt/danaco-console/danaco-console"
install -m 0755 "$PAKOWANIE/danaco-narzedzia" "$ROBOCZY/opt/danaco-console/danaco-narzedzia"

# Pakiet interfejsu tam, gdzie jednostka wskazuje rdzeniowi (client/dist).
mkdir -p "$ROBOCZY/opt/danaco-console/client"
cp -a "$BUDOWA/client/dist" "$ROBOCZY/opt/danaco-console/client/dist"
find "$ROBOCZY/opt/danaco-console/client" -type d -exec chmod 0755 {} +
find "$ROBOCZY/opt/danaco-console/client" -type f -exec chmod 0644 {} +

# Rdzeń na ścieżce systemu. Dowiązanie, nie kopia: os.Executable rozwija
# dowiązanie, więc serwer narzędzi zostaje znaleziony obok pliku właściwego.
ln -sf /opt/danaco-console/danaco-console "$ROBOCZY/usr/bin/danaco-console"

# Skrypt arsenału jedzie z pakietem, bo stawia go się na SERWERZE, a nie w drzewie
# budowy. Treści nie powtarzamy — bierzemy plik taki, jaki jest.
if [ -r "$SKRYPTY/arsenal-serwera.sh" ]; then
	install -m 0755 "$SKRYPTY/arsenal-serwera.sh" \
		"$ROBOCZY/opt/danaco-console/scripts/arsenal-serwera.sh"
else
	printf 'UWAGA: brak scripts/arsenal-serwera.sh — pakiet pójdzie bez skryptu arsenału\n' >&2
fi

# Cztery dokumenty produktu z korzenia repozytorium.
BRAK_DOKUMENTU=nie
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
	if [ -r "$KORZEN/$dokument" ]; then
		install -m 0644 "$KORZEN/$dokument" "$ROBOCZY/usr/share/doc/danaco-console/$dokument"
	else
		printf 'UWAGA: brak dokumentu %s w %s\n' "$dokument" "$KORZEN" >&2
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

# ── Złożenie ─────────────────────────────────────────────────────────────────
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
printf '\nodbiór:\n  dpkg-deb -I %s\n  dpkg-deb -c %s\n' "$PAKIET" "$PAKIET"
