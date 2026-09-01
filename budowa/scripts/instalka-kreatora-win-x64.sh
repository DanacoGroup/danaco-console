#!/usr/bin/env bash
# Instalka kreatora Danaco Console dla Windows 11 na x64: złożenie instalatora
# NSIS na Linuksie z okna kreatora sześciu kroków. To jest plik, który pobiera
# Operator ze strony „Pobierz"; powłokę programu kreator ściąga sam z kanału
# wydań, więc w tej instalce jej nie ma.
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
KREATOR="$KORZEN/budowa/instalator/src-tauri"
INTERFEJS="$KORZEN/budowa/instalator/interfejs"
DOSTAWA="$KORZEN/design/zasoby"
WYKAZ="$KORZEN/budowa/witryna/wydania.json"
CEL="x86_64-pc-windows-gnu"
WERSJA="2.0.0"
# Katalog wydania wolno wskazać z zewnątrz, tak samo jak w pakiecie serwera:
# złożenie próbne nie ma odkładać pliku między wydania.
WYDANIE="${DANACO_KATALOG_WYDANIA:-$KORZEN/budowa/wydania/$WERSJA-$(date +%Y-%m-%d)}"
NAZWA_WYDANIA="Danaco Console — Instalator_${WERSJA}_x64-setup.exe"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

# Rozmiar katalogu Security binarium PE, czyli tablicy certyfikatów podpisu
# Authenticode. Czytany wprost z nagłówka, bo osslsigncode nie stoi na tej
# maszynie, a jedynym rozstrzygającym śladem podpisu w pliku jest ten wpis.
rozmiar_podpisu() {
  python3 - "$1" <<'PYTON'
import struct, sys
with open(sys.argv[1], 'rb') as plik:
    dane = plik.read()
poczatek = struct.unpack_from('<I', dane, 0x3C)[0]
if dane[poczatek:poczatek + 4] != b'PE\0\0':
    raise SystemExit('plik nie jest binarium PE')
naglowek = poczatek + 24
magia = struct.unpack_from('<H', dane, naglowek)[0]
# Katalog danych zaczyna sie za polami naglowka opcjonalnego, a te maja inna
# dlugosc w PE32 i w PE32+; wpis czwarty to tablica certyfikatow.
if magia == 0x20B:
    odstep = 144
elif magia == 0x10B:
    odstep = 128
else:
    raise SystemExit('naglowek opcjonalny nie jest ani PE32, ani PE32+')
print(struct.unpack_from('<I', dane, naglowek + odstep + 4)[0])
PYTON
}

zglos "Sprawdzenie narzędzi"
# Brak narzędzia ujawniłby się dopiero po kilkuminutowej budowie, dlatego
# sprawdzenie stoi przed nią. Llvm-mingw nie może być na ścieżce, bo cel
# x86_64-pc-windows-gnu wymaga systemowego mingw-w64 i jego biblioteki libgcc.
for narzedzie in cargo rustup x86_64-w64-mingw32-gcc makensis python3; do
  command -v "$narzedzie" >/dev/null \
    || padnij "brak narzędzia: $narzedzie (apt: mingw-w64, nsis, python3; rustup: https://rustup.rs)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"
printf '  jest: cel %s\n' "$CEL"

zglos "Sprawdzenie artefaktów wejściowych"
# Okno kreatora jest wydawane z katalogu `interfejs/`, a warstwę projektową
# wkłada tam skrypt budowania kreatora, kopiując `design/zasoby/`. Brak dostawy
# wywróciłby budowę paniką build.rs, więc sprawdzenie stoi tutaj.
for artefakt in "$INTERFEJS/index.html" "$INTERFEJS/integracja.js" "$DOSTAWA"; do
  [ -e "$artefakt" ] || padnij "brak artefaktu wejściowego: $artefakt"
  printf '  jest: %s\n' "${artefakt#"$KORZEN/"}"
done

zglos "Sprawdzenie wykazu wydań"
# Krok 5 czyta wykaz najpierw z kanału (`wykaz_biezacy` w
# `instalator/src-tauri/src/pobranie.rs`), a kopia wkompilowana z tego pliku jest
# zapasem na kanał, który nie odpowie. Wykaz bez pozycji dla tej architektury
# dałby kreator, który dochodzi do kroku 5 i odmawia — a odmowę widać dopiero na
# maszynie z Windows, więc wykaz sprawdzany jest tu, przed budową.
python3 - "$WYKAZ" <<'PYTON'
import json, sys
wykaz = json.load(open(sys.argv[1], encoding='utf-8'))
pasujace = [
    p for p in wykaz.get('wydania', [])
    if p.get('system') == 'Windows' and p.get('architektura') == 'x64' and p.get('postac') == 'hybryda'
]
if not pasujace:
    raise SystemExit('wykaz nie niesie powloki Windows/x64 w postaci hybryda')
pozycja = pasujace[0]
for pole in ('plik', 'suma', 'rozmiarBajty'):
    if not pozycja.get(pole):
        raise SystemExit('pozycja powloki Windows/x64 nie niesie pola %s' % pole)
if not pozycja['plik'].startswith(wykaz['kanal']['adres']):
    raise SystemExit('pozycja powloki lezy poza kanalem %s' % wykaz['kanal']['adres'])
print('  powłoka do pobrania: %s' % pozycja['nazwaPliku'])
print('  adres w kanale: %s' % pozycja['plik'])
PYTON

zglos "Sprawdzenie poświadczeń kanału i żądania podpisu"
# Poświadczenia kanału wchodzą do binarium przy KOMPILACJI, przez option_env!
# (`instalator/src-tauri/src/pobranie.rs`). Idą do cargo jawnie albo są jawnie
# zdejmowane ze środowiska: połowa poświadczenia daje nagłówek uwierzytelnienia
# złożony z pustego napisu, a odmowa 401 wygląda wtedy na usterkę kanału.
if [ -n "${DANACO_KANAL_UZYTKOWNIK:-}" ] && [ -n "${DANACO_KANAL_HASLO:-}" ]; then
  POSWIADCZENIA=(DANACO_KANAL_UZYTKOWNIK="$DANACO_KANAL_UZYTKOWNIK" DANACO_KANAL_HASLO="$DANACO_KANAL_HASLO")
  OPIS_POSWIADCZEN="wpisane w instalkę"
elif [ -n "${DANACO_KANAL_UZYTKOWNIK:-}${DANACO_KANAL_HASLO:-}" ]; then
  padnij "podano jedną połowę poświadczeń kanału — nagłówek uwierzytelnienia składa się z nazwy i hasła naraz"
else
  POSWIADCZENIA=(-u DANACO_KANAL_UZYTKOWNIK -u DANACO_KANAL_HASLO)
  OPIS_POSWIADCZEN="BRAK — kreator sięgnie wyłącznie po pozycje spod otwartej ścieżki kanału"
fi
printf '  poświadczenia kanału: %s\n' "$OPIS_POSWIADCZEN"
# Żądanie podpisu sprawdzane jest tutaj, przed kilkuminutową budową, a nie
# dopiero przy odkładaniu wyniku.
case "${DANACO_PODPIS:-pomijany}" in
  wymagany | pomijany) ;;
  *) padnij "DANACO_PODPIS ma wartość ${DANACO_PODPIS}, a przyjmowane są wyłącznie: wymagany, pomijany" ;;
esac
printf '  podpis Authenticode: %s\n' "${DANACO_PODPIS:-pomijany}"

zglos "Budowa kreatora dla $CEL"
# Cecha tauri/custom-protocol jest obowiązkowa: bez niej gotowy plik zachowuje
# się jak budowa deweloperska i szuka okna na serwerze rozwojowym.
( cd "$KREATOR" \
  && env "${POSWIADCZENIA[@]}" \
       cargo build --release --target "$CEL" --features tauri/custom-protocol )

KREATOR_EXE="$KREATOR/target/$CEL/release/danaco-instalator.exe"
[ -f "$KREATOR_EXE" ] || padnij "cargo nie zgłosił błędu, ale binarki kreatora nie ma: $KREATOR_EXE"
OPIS_KREATORA="$(file -b "$KREATOR_EXE")"
printf '%s' "$OPIS_KREATORA" | grep -qi 'PE32+ executable' \
  || padnij "kreator nie jest binarką Windows (PE32+): $OPIS_KREATORA"
printf '%s' "$OPIS_KREATORA" | grep -qi 'x86-64' \
  || padnij "kreator nie jest binarką PE dla x86-64: $OPIS_KREATORA"
printf '  architektura kreatora: %s\n' "$OPIS_KREATORA"

zglos "Zapora: osadzone zasoby okna kreatora"
# Sondą jest ŚCIEŻKA pliku okna wewnątrz binarium, odczytana z katalogu
# interfejsu, a nie wpisana na sztywno. Treści zasobów szukać nie ma sensu:
# tauri wkłada je spakowane, a ścieżki stoją nieskompresowane i to one
# świadczą o osadzeniu.
SONDA="$(find "$INTERFEJS" -maxdepth 1 -name '*.js' -printf '/%f\n' | LC_ALL=C sort | head -1)"
[ -n "$SONDA" ] \
  || padnij "nie ma z czego zrobić sondy: w $INTERFEJS nie ma pliku .js — bez niej nie da się sprawdzić, czy zasoby są osadzone"
printf '  sonda: %s\n' "$SONDA"
# Liczenie trafień w podstawieniu polecenia zamiast w potoku z grep -q omija
# pułapkę tego trybu powłoki: grep -q kończy się na pierwszym trafieniu i sygnał
# przerwania zamieniłby sukces sondy w porażkę całego potoku.
LICZBA_SONDY="$(strings -a "$KREATOR_EXE" | grep -cF -- "$SONDA" || true)"
LICZBA_WARSTWY="$(strings -a "$KREATOR_EXE" | grep -cF -- '/warstwa/zasoby/' || true)"
printf '  trafienia sondy: %s, wystąpienia /warstwa/zasoby/: %s\n' "$LICZBA_SONDY" "$LICZBA_WARSTWY"
{ [ "$LICZBA_SONDY" -gt 0 ] && [ "$LICZBA_WARSTWY" -gt 0 ]; } \
  || padnij "kreator NIE ma osadzonego okna — cecha tauri/custom-protocol nie weszła, więc okno pójdzie do serwera rozwojowego i na maszynie Operatora zostanie puste. To produkt zepsuty, choć instalka wyglądałaby poprawnie"
printf '  osadzone okno kreatora: potwierdzone\n'

# Suma treści okna liczona PO budowie, bo warstwę projektową wkłada do katalogu
# interfejsu dopiero build.rs; jest jedynym sprawdzianem, z jakiego stanu okna
# powstał gotowy plik.
SUMA_INTERFEJSU="$(
  cd "$INTERFEJS" \
    && find . -type f -print0 | LC_ALL=C sort -z | xargs -0 -r sha256sum | sha256sum | cut -d' ' -f1
)"
printf '  suma treści okna: %s\n' "$SUMA_INTERFEJSU"

zglos "Złożenie instalatora NSIS"
# `--bundles nsis` podane jawnie, choć konfiguracja ma ten cel ustawiony:
# skrypt składa instalator NSIS niezależnie od tego, jakie inne cele
# konfiguracja produktu wymienia.
( cd "$KREATOR" && cargo tauri bundle --target "$CEL" --bundles nsis )

ZLOZONY="$KREATOR/target/$CEL/release/bundle/nsis/$NAZWA_WYDANIA"
[ -f "$ZLOZONY" ] || padnij "makensis nie zgłosił błędu, ale pliku instalatora nie ma: $ZLOZONY"

zglos "Odbiór — czego w środku być nie może"
# Kreator POBIERA powłokę i nie wiezie ani jej, ani rdzenia. Gdyby wróciły do
# konfiguracji jako zasoby, instalka złożyłaby się bez błędu i wyszła jako
# produkt pełny pod nazwą kreatora; dowodem przeciwnym jest wykaz zawartości.
command -v 7z >/dev/null \
  || padnij "brak 7z (apt: p7zip-full) — bez wykazu zawartości nie ma dowodu, czego instalka nie wiezie"
SPIS="$(mktemp)"
ROZPAK="$(mktemp -d)"
trap 'rm -f "$SPIS"; rm -rf "$ROZPAK"' EXIT
7z l "$ZLOZONY" >"$SPIS"
for zakazany in danaco-console.exe danaco-narzedzia.exe danaco-console-powloka.exe; do
  [ "$(grep -c -- "$zakazany" "$SPIS" || true)" -eq 0 ] \
    || padnij "instalka kreatora niesie $zakazany — kreator ma to pobrać z kanału, nie wieźć ze sobą"
  printf '  nie ma w środku: %s\n' "$zakazany"
done

# Powyższe zapory badały binarkę z `target/`. Ta bada binarkę WYPAKOWANĄ
# z instalatora, bo tylko ona jest tym, co dostanie Operator: `bundle` bierze
# gotową binarkę z `target/` i sam nie kompiluje, więc podmiany nie byłoby widać
# inaczej.
7z x -o"$ROZPAK" "$ZLOZONY" >/dev/null 2>&1 \
  || padnij "nie udało się wypakować instalatora do sprawdzenia"
WIEZIONY="$ROZPAK/danaco-instalator.exe"
[ -f "$WIEZIONY" ] || padnij "w instalatorze nie ma kreatora danaco-instalator.exe"
OPIS_WIEZIONEGO="$(file -b "$WIEZIONY")"
printf '%s' "$OPIS_WIEZIONEGO" | grep -qi 'x86-64' \
  || padnij "kreator WIEZIONY przez instalator nie jest binarką x86-64: $OPIS_WIEZIONEGO"
TRAFIENIA_WIEZIONEGO="$(strings -a "$WIEZIONY" | grep -cF -- "$SONDA" || true)"
[ "$TRAFIENIA_WIEZIONEGO" -gt 0 ] \
  || padnij "kreator WIEZIONY przez instalator nie ma osadzonego okna — instalka wiezie budowę rozwojową"
printf '  wieziony kreator: %s, okno osadzone\n' "$OPIS_WIEZIONEGO"

zglos "Zapora: podpis Authenticode"
# Katalog Security gotowego pliku niesie tablicę certyfikatów; rozmiar zero
# znaczy plik bez wydawcy. SmartScreen ukrywa wtedy przycisk uruchomienia,
# a zasady firmowe blokują plik całkiem, więc stan podpisu jest mierzony
# i wypisany, a nie zakładany.
PODPIS_BAJTOW="$(rozmiar_podpisu "$ZLOZONY")"
if [ "$PODPIS_BAJTOW" -gt 0 ]; then
  OPIS_PODPISU="katalog Security $PODPIS_BAJTOW bajtów"
elif [ "${DANACO_PODPIS:-pomijany}" = "wymagany" ]; then
  padnij "katalog Security jest pusty, a DANACO_PODPIS=wymagany — wyniku nie odkładam do wydania"
else
  OPIS_PODPISU="BRAK — katalog Security pusty; wydanie wychodzi bez podpisu (DANACO_PODPIS=pomijany, rozstrzygnięcie 27)"
fi
printf '  %s\n' "$OPIS_PODPISU"

zglos "Odłożenie wyniku do wydania"
mkdir -p "$WYDANIE"
WYNIK="$WYDANIE/$NAZWA_WYDANIA"
cp -f "$ZLOZONY" "$WYNIK"

zglos "Pomiar wyniku"
printf 'ścieżka : %s\n' "$WYNIK"
printf 'rozmiar : %s (%s bajtów)\n' "$(du -h "$WYNIK" | cut -f1)" "$(stat -c%s "$WYNIK")"
printf 'typ     : %s\n' "$(file -b "$WYNIK")"
printf 'suma    : %s\n' "$(sha256sum "$WYNIK" | cut -d' ' -f1)"
printf 'okno    : suma treści %s\n' "$SUMA_INTERFEJSU"
printf 'kanał   : poświadczenia %s\n' "$OPIS_POSWIADCZEN"
printf 'podpis  : %s\n' "$OPIS_PODPISU"
# Ostatni wiersz listingu archiwum kończy się słowem oznaczającym liczbę
# plików; sama liczba stoi w wierszu bezpośrednio przed tym słowem.
printf 'wewnątrz: %s plików\n' "$(tail -1 "$SPIS" | awk '{print $(NF-1)}')"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i czy przechodzi wszystkie sześć kroków,
  - czy okno kreatora wstaje (wymaga WebView2 w systemie Windows 11),
  - czy krok 5 pobiera powłokę — wymaga maszyny z Windows i kanału
    odpowiadającego pod adresem z wykazu wydań.
  Wszystkie wymagają maszyny z Windows. Nie zakładaj ich powodzenia.
KONIEC
