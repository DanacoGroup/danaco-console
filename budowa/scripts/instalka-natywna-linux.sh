#!/usr/bin/env bash
# Instalka CAŁKOWICIE NATYWNA Danaco Console dla Linuksa — pakiet .deb.
#
# ── Czym się różni od instalki cienkiej ───────────────────────────────────────
# Instalka cienka (scripts/instalka-windows.sh) niesie samo okno, a arsenał i
# rdzeń stoją na serwerze. Ta niesie CAŁE SERCE PLATFORMY na maszynę Operatora,
# bo jest dla tego, kto pracuje bez serwera:
#   powłoka natywna (okno)                       — binarium Tauri
#   rdzeń danaco-console                         — server/cmd/danaco-console
#   serwer narzędzi danaco-narzedzia             — server/cmd/danaco-narzedzia
#   pakiet interfejsu client/dist                — npm run build
#   cztery dokumenty produktu                    — z korzenia repozytorium
#   arsenał: 30 programów zadeklarowanych w
#            server/internal/core/zaleznosci_zewnetrzne.go
#
# ── Jak rozstrzygnięty jest arsenał ──────────────────────────────────────────
# Nie pakujemy cudzych binariów, skoro dystrybucja je niesie. Podział biegnie
# po warstwie, którą liczy sam rdzeń (core.WarstwaZaleznosci) i którą oddaje
# wykazem `danaco-console --wykaz-zaleznosci`:
#
#   obowiazkowa-apt  → Depends pakietu .deb; apt stawia to razem z produktem.
#   snap             → poza .deb (apt nie niesie snapów); arsenał.
#   warsztat-go      → poza .deb (`go install`); arsenał.
#   warsztat-npm     → poza .deb (`npm i -g`); arsenał.
#   model-recznie    → poza .deb (wydania z GitHuba, wagi modeli); arsenał.
#   decyzyjna        → silnik kontenerów; Właściciel go WSTRZYMAŁ, więc nie ma
#                      go ani w Depends, ani w arsenale bez wyraźnej zgody.
#   arsenał mowy     → piper wraz z głosami .onnx oraz wagi faster-whisper;
#                      poza .deb, arsenał (`danaco-console --wykaz-mowy`).
#
# „Arsenał" znaczy tu jeden, już istniejący skrypt: scripts/arsenal-serwera.sh.
# Ten skrypt go NIE POWTARZA — wkłada go do pakietu i wystawia poleceniem
# `danaco-console-arsenal`. Powód jest jeden: arsenał czyta wykaz z binarium
# rdzenia, więc zmiana pola `Pakiet` w deklaracji narzędzia dojeżdża i do sondy
# startowej, i do prowizjonowania jednym ruchem. Druga lista rozjechałaby się
# z pierwszą przy pierwszej zmianie pakietu.
#
# ── Czemu ściąganie NIE dzieje się w skrypcie poinstalacyjnym ────────────────
# Polityka Debiana zabrania skryptom opiekuna sięgać do sieci: `dpkg -i` ma
# skończyć się w czasie skończonym i bez łącza. Wagi modeli i wydania z GitHuba
# to gigabajty — postinst wisiałby, a przerwana instalacja zostawiłaby pakiet
# w stanie połowicznym. Dlatego postinst robi to, co lokalne (prawa, dowiązania,
# sonda), a resztę NAZYWA i podaje jedno polecenie, które ją domyka.
#
# ── Układ w systemie plików i dlaczego właśnie taki ──────────────────────────
# Powłoka szuka rdzenia obok siebie (desktop/src-tauri/src/rdzen/lokalizacja.rs),
# rdzeń szuka serwera narzędzi obok siebie (server/internal/narzedzia/wpiecie.go),
# a pakietu interfejsu w `<katalog rdzenia>/dist` (src/rdzen/pakiet_klienta.rs).
# Ładunek leży więc w JEDNYM katalogu rdzenia, w /usr/lib (117 MB nie należy do
# /usr/bin), a widzialność obok powłoki daje dowiązanie /usr/bin/rdzen.
# Rozstrzygnięcia ścieżek pilnuje odbiór na końcu tego skryptu.
#
# Dowiązania /usr/bin/danaco-console NIE MA celowo: powłoka sprawdza je PRZED
# katalogiem `rdzen`, więc trafiłaby na binarium, obok którego nie leży katalog
# `dist` — i wstałaby bez interfejsu. Rdzeń z wiersza poleceń stoi więc pod
# nazwą danaco-console-rdzen.
#
# ── Dlaczego rozstrzygnięcia opisane są TUTAJ, a nie w profilu ───────────────
# tauri.natywna-linux.conf.json jest czytany z `deny_unknown_fields`, więc klucz
# komentarza (`"//"`) wywraca budowę zdaniem o nieoczekiwanej właściwości. Profil
# niesie same wartości; powody niesie ten skrypt, który go uruchamia.
#
# Profil jest NAKŁADKĄ na tauri.conf.json (json_patch::merge, RFC 7386): wartość
# null USUWA klucz bazowy, tablica podmienia tablicę. Dlatego `bundle.resources`
# ma tam null — baza wskazuje zasoby/rdzen/*.exe dla instalki Windows, a plików
# Windows w pakiecie .deb być nie może. Ładunek natywny wchodzi osobno, przez
# `bundle.linux.deb.files`, bo tylko tam wybieram ŚCIEŻKI DOCELOWE.
#
# Nazwy pakietów w Depends podane wariantem (`a | b`) tam, gdzie różnią się
# między wydaniami dystrybucji — 7-Zip, Chromium, Telnet, łańcuch Go, GTK.
# Wariant nie jest ostrożnością na wszelki wypadek: nazwa nieobecna w
# repozytorium ZATRZYMUJE instalację, a produkt ma wejść na maszynie Operatora,
# nie na tej jednej. Nazwy chwiejniejsze niż wariant je ratuje (shfmt) oraz to,
# czego brak nie zabiera modułu, siedzą w Recommends.
#
# Bundler dokłada do Depends własne trzy pozycje (libayatana-appindicator3-1,
# libwebkit2gtk-4.1-0, libgtk-3-0), więc dwie z nich stoją w polu Depends
# dwukrotnie. To NIE jest przeoczenie i nie usuwamy ich z profilu: bundler
# dokłada `libgtk-3-0` samo, a profil `libgtk-3-0 | libgtk-3-0t64` — czyli
# wariant szerszy, który wchodzi też tam, gdzie dystrybucja przeszła na t64.
# Powtórzenie w Depends jest dla dpkg obojętne (rozstrzyga sumę alternatyw),
# a utrata wariantu nie byłaby.
#
# Użycie:  bash budowa/scripts/instalka-natywna-linux.sh
set -euo pipefail

WERSJA="1.0.0"
DZIEN="2026-08-18"
KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BUDOWA="$KORZEN/budowa"
TAURI="$BUDOWA/desktop/src-tauri"
STAGING="$TAURI/zasoby-natywna-linux"
PROFIL="$TAURI/tauri.natywna-linux.conf.json"
WYDANIE="$BUDOWA/wydania/$WERSJA-$DZIEN"
NAZWA_WYDANIA="Danaco Console_${WERSJA}_natywna_amd64.deb"

powiedz() { printf '\n== %s\n' "$1"; }

# ── Rdzeń i serwer narzędzi ──────────────────────────────────────────────────
# Oba natywnie na Linuksa i oba do JEDNEGO katalogu: rdzeń szuka serwera
# narzędzi obok siebie, więc rozdzielenie ich zabrałoby sterowanie platformą
# (narzedzia.Wpis odmówiłby brakiem binarium).
powiedz "rdzeń i serwer narzędzi"
rm -rf "$STAGING"
mkdir -p "$STAGING/rdzen" "$STAGING/dokumentacja" "$STAGING/arsenal" "$STAGING/maintainer"
cd "$BUDOWA"
go build -o "$STAGING/rdzen/danaco-console" ./server/cmd/danaco-console
go build -o "$STAGING/rdzen/danaco-narzedzia" ./server/cmd/danaco-narzedzia

# ── Pakiet interfejsu ────────────────────────────────────────────────────────
powiedz "pakiet interfejsu"
npm --prefix "$BUDOWA/client" run build
cp -r "$BUDOWA/client/dist" "$STAGING/rdzen/dist"

# ── Dokumenty produktu i arsenał ─────────────────────────────────────────────
powiedz "dokumenty produktu i skrypt arsenału"
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
	install -m 644 "$KORZEN/$dokument" "$STAGING/dokumentacja/$dokument"
done
cp -r "$BUDOWA/pomocniki" "$STAGING/pomocniki"
rm -rf "$STAGING/pomocniki/transkrypcja/__pycache__"
install -m 755 "$BUDOWA/scripts/arsenal-serwera.sh" "$STAGING/arsenal/arsenal-serwera.sh"

# ── Skrypty opiekuna pakietu ─────────────────────────────────────────────────
# Powstają tutaj, a nie leżą w drzewie, bo są wytworem budowy tego jednego
# profilu — tak jak binaria i pakiet interfejsu obok nich.
powiedz "skrypty opiekuna pakietu"
cat >"$STAGING/maintainer/postinst" <<'POSTINST'
#!/bin/sh
# Po rozpakowaniu pakietu natywnego Danaco Console.
# Robi wyłącznie rzeczy LOKALNE — sieci nie tyka (polityka Debiana).
set -e

LIB=/usr/lib/danaco-console-natywna

# Dowiązanie katalogu rdzenia obok powłoki. Powłoka sprawdza kolejno
# /usr/bin/danaco-console i /usr/bin/rdzen/danaco-console; pierwszego NIE
# tworzymy celowo — binarium bez katalogu dist obok nie poda interfejsu.
ln -sfn "$LIB/rdzen" /usr/bin/rdzen
# Rdzeń pod własnym poleceniem — do wykazów i pracy z wiersza poleceń.
# Nazwa inna niż danaco-console, żeby nie wyprzedzić szukania powłoki.
ln -sfn "$LIB/rdzen/danaco-console" /usr/bin/danaco-console-rdzen
ln -sfn "$LIB/arsenal/arsenal-serwera.sh" /usr/bin/danaco-console-arsenal

chmod 755 "$LIB/rdzen/danaco-console" "$LIB/rdzen/danaco-narzedzia" \
	"$LIB/arsenal/arsenal-serwera.sh"

echo "Danaco Console — instalka natywna: rdzeń, serwer narzędzi i interfejs stoją w $LIB"

# Sonda rdzenia mówi wprost, czego brakuje — tym samym wykazem, którym rdzeń
# melduje stan przy starcie. Bez tego Operator dowiadywałby się o braku dopiero
# po naciśnięciu przycisku, osobno przy każdej funkcji.
BRAKI=$("$LIB/rdzen/danaco-console" --wykaz-zaleznosci 2>/dev/null |
	awk -F'\t' '$1 !~ /^#/ && $4 != "tak" { print "  - " $5 " (" $2 ") — " $3 }' || true)
if [ -n "$BRAKI" ]; then
	echo "Arsenał niekompletny — pakiet tego nie niesie i mówi o tym wprost:"
	echo "$BRAKI"
	echo "Domknięcie jednym poleceniem (ściąga wydania z GitHuba, snapy,"
	echo "warsztat Go i npm, modele obrazu oraz arsenał mowy wraz z wagami):"
	echo "    sudo danaco-console-arsenal"
else
	echo "Arsenał kompletny — rdzeń ma czym wykonać każdą czynność."
fi

exit 0
POSTINST

cat >"$STAGING/maintainer/prerm" <<'PRERM'
#!/bin/sh
# Przed usunięciem pakietu — dowiązania z /usr/bin nie należą do listy plików
# pakietu (robi je postinst), więc dpkg ich nie sprzątnie sam.
set -e
rm -f /usr/bin/rdzen /usr/bin/danaco-console-rdzen /usr/bin/danaco-console-arsenal
exit 0
PRERM

chmod 755 "$STAGING/maintainer/postinst" "$STAGING/maintainer/prerm"

# ── Powłoka: budowa z cechą custom-protocol ─────────────────────────────────
# `cargo tauri bundle` NICZEGO NIE KOMPILUJE — bierze gotowe binarium
# z target/release. Kto na to nie uważa, wkłada do produktu powłokę w TRYBIE
# ROZWOJOWYM, bo `tauri::is_dev()` to dokładnie `!cfg!(feature =
# "custom-protocol")`, a Cargo.toml powłoki tej cechy nie włącza (i nie ma jej
# włączać — w pracy przy oknie chodzi devUrl).
#
# Skutek pominięcia jest dla instalki natywnej zgubny: zrodlo_interfejsu::ustal
# idzie wtedy do devUrl (localhost:5173) i NIGDY nie sprawdza warunku „rdzeń
# nasłuchujący" — okno nie wzięłoby interfejsu z rdzenia, którego ten sam pakiet
# niesie, i Operator dostałby pustą stronę.
#
# Dlatego budujemy ZAWSZE i z cechą — nie „jeśli binarium nie istnieje".
# Binarium leżące w target/release mogło zostać zbudowane bez cechy przez
# kogokolwiek innego, a po nazwie pliku tego nie widać.
powiedz "powłoka natywna (cecha custom-protocol)"
# shellcheck source=/dev/null
. "$HOME/.cargo/env"
cd "$TAURI"
cargo build --release --features tauri/custom-protocol

powiedz "pakowanie .deb"
cargo tauri bundle --bundles deb --config "$PROFIL"

ZBUDOWANY=$(find "$TAURI/target/release/bundle/deb" -maxdepth 1 -name '*.deb' \
	-newer "$PROFIL" -print -quit)
if [ -z "$ZBUDOWANY" ]; then
	echo "BŁĄD: nie znalazłem świeżego pakietu .deb" >&2
	exit 1
fi

# ── Nazwa pakietu Debiana ────────────────────────────────────────────────────
# Bundler liczy pole `Package` z productName („Danaco Console" → danaco-console)
# i nie ma na to osobnego pokrętła. Sama ta nazwa KOLIDUJE z pakietem serwera:
# przy instalacji jeden nadpisywałby drugi. Nazwy produktu nie ruszamy — jedna
# platforma ma jedną nazwę w oknie i w menu, niezależnie od postaci dostawy —
# więc poprawiamy samo pole kontroli.
#
# Robimy to wymianą składnika `control.tar.*` w archiwum `ar`, a NIE przez
# `dpkg-deb -R` i `-b`: przepakowanie ścisnęłoby na nowo 117 MB ładunku, którego
# ani jeden bajt się nie zmienia. Składnik kontroli waży kilka kilobajtów.
# Sumy md5sums pozostają prawdziwe, bo `data.tar.*` idzie bit w bit.
powiedz "nazwa pakietu Debiana"
SKLADNIK_KONTROLI=$(ar t "$ZBUDOWANY" | grep '^control\.tar')
ROBOCZY=$(mktemp -d)
trap 'rm -rf "$ROBOCZY"' EXIT
( cd "$ROBOCZY" && ar x "$ZBUDOWANY" )
mkdir -p "$ROBOCZY/kontrola"
tar -xf "$ROBOCZY/$SKLADNIK_KONTROLI" -C "$ROBOCZY/kontrola"
sed -i 's/^Package: .*/Package: danaco-console-natywna/' "$ROBOCZY/kontrola/control"
rm -f "$ROBOCZY/$SKLADNIK_KONTROLI"
# Właściciel numeryczny 0:0 — pliki kontroli należą do roota, a pakuje je
# zwykły użytkownik; bez tego dpkg zobaczyłby ubuntu:ubuntu.
case "$SKLADNIK_KONTROLI" in
*.zst) SCISK=(--zstd) ;;
*.xz) SCISK=(-J) ;;
*.gz) SCISK=(-z) ;;
*) echo "BŁĄD: nieznany ścisk składnika kontroli $SKLADNIK_KONTROLI" >&2; exit 1 ;;
esac
tar "${SCISK[@]}" -cf "$ROBOCZY/$SKLADNIK_KONTROLI" \
	--owner=0 --group=0 --numeric-owner -C "$ROBOCZY/kontrola" .
# Kolejność składników nie jest dowolna: dpkg wymaga `debian-binary` pierwszego.
ar rcD "$ROBOCZY/przepakowany.deb" "$ROBOCZY/debian-binary" \
	"$ROBOCZY/$SKLADNIK_KONTROLI" "$ROBOCZY"/data.tar.*
ZBUDOWANY="$ROBOCZY/przepakowany.deb"
dpkg-deb -I "$ZBUDOWANY" >/dev/null # odmówi, gdy archiwum wyszło niepoprawne

# ── Odbiór: czy pakiet niesie to, co ma nieść ───────────────────────────────
# Sprawdzian jest tu, a nie w oku czytającego wynik `dpkg-deb -c`, bo brak
# jednego z tych czterech znaczy produkt, który nie wstanie: bez rdzenia nie ma
# platformy, bez serwera narzędzi nie ma sterowania, bez dist nie ma okna,
# bez dokumentów nie ma produktu.
powiedz "odbiór pakietu"
# Ścieżki bierzemy z ostatniego pola wiersza i ścinamy wiodące `./`, bo różne
# wydania dpkg-deb wypisują je raz z tym przedrostkiem, raz bez — dopasowanie
# do jednej z tych postaci przepuściłoby pakiet pusty jako pakiet dobry.
ZAWARTOSC=$(dpkg-deb -c "$ZBUDOWANY" | awk '{ sub(/^\.\//, "", $NF); print $NF }')
for wymagane in \
	usr/lib/danaco-console-natywna/rdzen/danaco-console \
	usr/lib/danaco-console-natywna/rdzen/danaco-narzedzia \
	usr/lib/danaco-console-natywna/rdzen/dist/index.html \
	usr/share/doc/danaco-console-natywna/README.md \
	usr/share/doc/danaco-console-natywna/INSTALACJA-I-KONFIGURACJA.md \
	usr/share/doc/danaco-console-natywna/INSTRUKCJA-UZYTKOWANIA.md \
	usr/share/doc/danaco-console-natywna/LICENSE.md \
	usr/lib/danaco-console-natywna/arsenal/arsenal-serwera.sh; do
	if ! printf '%s\n' "$ZAWARTOSC" | grep -qxF "$wymagane"; then
		echo "BŁĄD: pakiet nie niesie $wymagane" >&2
		exit 1
	fi
done

# ── Zapora: powłoka NIE MOŻE być w trybie rozwojowym ────────────────────────
# Mierzymy powłokę WYPAKOWANĄ Z GOTOWEGO PAKIETU, nie tę z target/release:
# między budową a pakowaniem stoi krok, który może wziąć inne binarium, a zapora
# ma pilnować tego, co Operator dostanie do ręki.
#
# Sondą jest nazwa pliku interfejsu z sumą treści w nazwie (index-<suma>.js
# z client/dist/assets). Do binarki nie ma jak trafić inaczej niż przez osadzone
# `client/dist`, a suma zmienia się z każdą budową interfejsu — dlatego sondę
# czytamy z drzewa przy każdym przebiegu, nie wpisujemy na sztywno.
#
# Czego sondą być NIE MOŻE: `woff2`, `index.html` ani `<!doctype html>` zapalają
# się także w powłoce BEZ osadzonych zasobów — tauri niesie własną tablicę typów
# MIME i obsługę indeksu, więc te napisy siedzą w binarce zawsze i taka sonda
# przepuściłaby powłokę rozwojową jako dobrą.
powiedz "zapora: powłoka poza trybem rozwojowym"
SONDY=("$BUDOWA"/client/dist/assets/index-*.js)
if [ ${#SONDY[@]} -ne 1 ] || [ ! -f "${SONDY[0]}" ]; then
	echo "BŁĄD: nie umiem wskazać jednej sondy w client/dist/assets" >&2
	exit 1
fi
SONDA=$(basename "${SONDY[0]}")
# Trafienia liczymy w PODSTAWIENIU POLECENIA, a nie przez `grep -q`: `-q`
# zamyka wejście przy pierwszym trafieniu, więc przy `set -o pipefail` potok
# wywraca się sygnałem SIGPIPE dokładnie wtedy, gdy sonda się ZNAJDZIE —
# zapora odmawiałaby na pakiecie dobrym. `grep -c` czyta wejście do końca.
TRAFIENIA=$(dpkg-deb --fsys-tarfile "$ZBUDOWANY" |
	tar -xO --wildcards '*usr/bin/danaco-console-powloka' |
	strings -a | grep -cF "$SONDA" || true)
printf 'sonda:    %s\n' "$SONDA"
printf 'trafienia w powłoce z pakietu: %s\n' "$TRAFIENIA"
if [ "$TRAFIENIA" -eq 0 ]; then
	echo "BŁĄD: powłoka w pakiecie nie niesie osadzonego client/dist — jest" >&2
	echo "w trybie rozwojowym i pójdzie do devUrl zamiast do rdzenia." >&2
	echo "Powód: binarium zbudowano bez cechy tauri/custom-protocol." >&2
	exit 1
fi

mkdir -p "$WYDANIE"
install -m 644 "$ZBUDOWANY" "$WYDANIE/$NAZWA_WYDANIA"
# Suma obok pakietu powstaje w tym samym kroku, co pakiet. Gdy zostaje z budowy
# poprzedniej, jest gorsza niż jej brak: mówi o pliku, którego już nie ma,
# a wygląda na sprawdzenie.
( cd "$WYDANIE" && sha256sum "$NAZWA_WYDANIA" >"$NAZWA_WYDANIA.sha256" )

powiedz "wynik"
printf 'pakiet:   %s\n' "$WYDANIE/$NAZWA_WYDANIA"
printf 'rozmiar:  %s bajtów\n' "$(stat -c%s "$WYDANIE/$NAZWA_WYDANIA")"
printf 'sha256:   %s\n' "$(sha256sum "$WYDANIE/$NAZWA_WYDANIA" | cut -d' ' -f1)"
dpkg-deb -I "$WYDANIE/$NAZWA_WYDANIA"
