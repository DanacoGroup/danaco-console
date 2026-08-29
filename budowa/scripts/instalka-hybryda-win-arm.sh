#!/usr/bin/env bash
# Instalka Danaco Console dla Windows 11 na ARM64: złożenie instalatora NSIS
# na Linuksie przez cargo xwin i clang, dla produktu hybrydowego, w którym
# rdzeń działa na serwerze wdrożenia, a u operatora staje wyłącznie okno.
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
KLIENT="$KORZEN/budowa/klient"
CEL="aarch64-pc-windows-msvc"
LLVM_MINGW="${LLVM_MINGW:-/opt/llvm-mingw}"
WERSJA="2.0.0"
WYDANIE="$KORZEN/budowa/wydania/$WERSJA-$(date +%Y-%m-%d)"
NAZWA_WYDANIA="Danaco Console_${WERSJA}_hybryda_arm64-setup.exe"

export PATH="$LLVM_MINGW/bin:$PATH"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Sprawdzenie narzędzi"
# Brak któregokolwiek ujawniłby się dopiero po kilkuminutowej budowie, dlatego
# sprawdzenie stoi przed nią, nie po.
for narzedzie in cargo rustup cargo-xwin makensis clang llvm-lib; do
  command -v "$narzedzie" >/dev/null \
    || padnij "brak narzędzia: $narzedzie (apt: nsis, clang; cargo install cargo-xwin; llvm-lib z wydania llvm-mingw rozpakowanego do $LLVM_MINGW)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"
printf '  jest: cel %s\n' "$CEL"

zglos "Sprawdzenie artefaktów wejściowych"
# Powłoka niesie własne okno, więc zbudowany klient jest jej jedynym artefaktem
# wejściowym. Rdzenia nie jest sprawdzany — w tej instalce go nie ma.
[ -e "$KLIENT/dist/index.html" ] \
  || padnij "brak artefaktu: budowa/klient/dist/index.html (zbuduj: npm run budowanie w budowa/klient)"
printf '  jest: budowa/klient/dist/index.html\n'

zglos "Budowa powłoki dla $CEL"
# Budowana jest sama powłoka, bez rdzenia — jedyna binarka tego produktu.
# Cecha tauri/custom-protocol jest obowiązkowa: bez niej gotowy plik zachowuje
# się jak budowa deweloperska i nie szuka rdzenia na serwerze wdrożenia.
( cd "$POWLOKA" \
  && cargo xwin build --release --target "$CEL" \
       --features tauri/custom-protocol --cross-compiler clang )

POWLOKA_EXE="$POWLOKA/target/$CEL/release/danaco-console-powloka.exe"
[ -f "$POWLOKA_EXE" ] || padnij "cargo nie zgłosił błędu, ale binarki powłoki nie ma: $POWLOKA_EXE"
# Cel budowy da się pomylić w wierszu poleceń, a instalka z binarką dla innej
# architektury w środku jest bezużyteczna, dlatego sprawdzian architektury
# stoi przed złożeniem instalatora, nie dopiero po nim.
OPIS_POWLOKI="$(file -b "$POWLOKA_EXE")"
printf '%s' "$OPIS_POWLOKI" | grep -qi 'PE32+ executable' \
  || padnij "powłoka nie jest binarką Windows (PE32+): $OPIS_POWLOKI"
printf '%s' "$OPIS_POWLOKI" | grep -qiE 'ARM64|Aarch64' \
  || padnij "powłoka nie jest binarką PE dla ARM64: $OPIS_POWLOKI"
printf '  architektura powłoki: %s\n' "$OPIS_POWLOKI"

zglos "Zapora: osadzone zasoby interfejsu"
# Zapora sprawdza obecność osadzonych zasobów interfejsu, bo tylko czynna
# cecha custom-protocol je wkompilowuje. Sondą jest nazwa pliku interfejsu
# z sumą treści w nazwie, odczytana z katalogu, a nie wpisana na sztywno.
SONDA="$(find "$KLIENT/dist/assets" -maxdepth 1 -name 'index-*.js' -printf '%f\n' 2>/dev/null | head -1)"
[ -n "$SONDA" ] \
  || padnij "nie ma z czego zrobić sondy: brak budowa/klient/dist/assets/index-*.js — bez niej nie da się sprawdzić, czy zasoby są osadzone"
printf '  sonda: %s\n' "$SONDA"
# Liczenie trafień w podstawieniu polecenia zamiast w potoku z grep -q omija
# pułapkę tego trybu powłoki: grep -q kończy się na pierwszym trafieniu i sygnał
# przerwania zamieniłby sukces sondy w porażkę całego potoku.
LICZBA_SONDY="$(strings -a "$POWLOKA_EXE" | grep -c -- "$SONDA" || true)"
LICZBA_ZASOBOW="$(strings -a "$POWLOKA_EXE" | grep -c '/assets/' || true)"
printf '  trafienia sondy: %s, wystąpienia /assets/: %s\n' "$LICZBA_SONDY" "$LICZBA_ZASOBOW"
{ [ "$LICZBA_SONDY" -gt 0 ] && [ "$LICZBA_ZASOBOW" -gt 0 ]; } \
  || padnij "powłoka NIE ma osadzonych zasobów interfejsu — cecha tauri/custom-protocol nie weszła, więc tauri::is_dev() zwróci prawdę i okno pójdzie do serwera rozwojowego zamiast do rdzenia na serwerze. To produkt zepsuty, choć instalka wyglądałaby poprawnie"
printf '  osadzone zasoby interfejsu: potwierdzone\n'

zglos "Złożenie instalatora NSIS"
# `--bundles nsis` podane jawnie, choć konfiguracja ma ten cel ustawiony:
# skrypt składa instalator NSIS niezależnie od tego, jakie inne cele
# konfiguracja produktu wymienia.
( cd "$POWLOKA" && cargo tauri bundle --target "$CEL" --bundles nsis )

ZLOZONY="$(find "$POWLOKA/target/$CEL/release/bundle/nsis" -maxdepth 1 -name '*-setup.exe' -print -quit)"
[ -n "$ZLOZONY" ] || padnij "makensis nie zgłosił błędu, ale pliku instalatora nie ma"

zglos "Odbiór — czego w środku być nie może"
# Odbiór jest częścią budowy. Gdyby do konfiguracji wróciły zasoby rdzenia,
# instalka złożyłaby się bez błędu i wyszłaby jako produkt pełny pod nazwą
# hybrydowego; dowodem przeciwnym jest wykaz zawartości gotowego pliku.
command -v 7z >/dev/null \
  || padnij "brak 7z (apt: p7zip-full) — bez wykazu zawartości nie ma dowodu, że rdzenia w instalce nie ma"
SPIS="$(mktemp)"
ROZPAK="$(mktemp -d)"
trap 'rm -f "$SPIS"; rm -rf "$ROZPAK"' EXIT
7z l "$ZLOZONY" >"$SPIS"
for zakazany in danaco-console.exe danaco-narzedzia.exe; do
  [ "$(grep -c -- "$zakazany" "$SPIS" || true)" -eq 0 ] \
    || padnij "instalka niesie $zakazany — do produktu wrócił rdzeń; instalka NIE jest cienka"
  printf '  nie ma w środku: %s\n' "$zakazany"
done

# Powyższe zapory badały binarkę z `target/`. Ta bada binarkę WYPAKOWANĄ
# z instalatora, bo tylko ona jest tym, co dostanie Operator: `bundle` bierze
# gotową binarkę z `target/` i sam nie kompiluje, więc podmiany nie byłoby widać
# inaczej.
7z x -o"$ROZPAK" "$ZLOZONY" >/dev/null 2>&1 \
  || padnij "nie udało się wypakować instalatora do sprawdzenia"
WIEZIONA="$ROZPAK/danaco-console-powloka.exe"
[ -f "$WIEZIONA" ] || padnij "w instalatorze nie ma powłoki danaco-console-powloka.exe"
OPIS_WIEZIONEJ="$(file -b "$WIEZIONA")"
printf '%s' "$OPIS_WIEZIONEJ" | grep -qiE 'ARM64|Aarch64' \
  || padnij "powłoka WIEZIONA przez instalator nie jest binarką ARM64: $OPIS_WIEZIONEJ"
TRAFIENIA_WIEZIONEJ="$(strings -a "$WIEZIONA" | grep -c -- "$SONDA" || true)"
[ "$TRAFIENIA_WIEZIONEJ" -gt 0 ] \
  || padnij "powłoka WIEZIONA przez instalator nie ma osadzonych zasobów interfejsu — instalka wiezie budowę rozwojową"
printf '  wieziona powłoka: %s, zasoby osadzone\n' "$OPIS_WIEZIONEJ"

zglos "Odłożenie wyniku do wydania"
mkdir -p "$WYDANIE"
WYNIK="$WYDANIE/$NAZWA_WYDANIA"
cp -f "$ZLOZONY" "$WYNIK"

zglos "Pomiar wyniku"
printf 'ścieżka : %s\n' "$WYNIK"
printf 'rozmiar : %s (%s bajtów)\n' "$(du -h "$WYNIK" | cut -f1)" "$(stat -c%s "$WYNIK")"
printf 'typ     : %s\n' "$(file -b "$WYNIK")"
printf 'suma    : %s\n' "$(sha256sum "$WYNIK" | cut -d' ' -f1)"
# Ostatni wiersz listingu archiwum kończy się słowem oznaczającym liczbę
# plików; sama liczba stoi w wierszu bezpośrednio przed tym słowem.
printf 'wewnątrz: %s plików\n' "$(tail -1 "$SPIS" | awk '{print $(NF-1)}')"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i czy przechodzi do końca,
  - czy zainstalowane okno wstaje (wymaga WebView2 dla ARM64 w Windows 11),
  - czy okno łączy się z rdzeniem na serwerze wdrożenia — wymaga stojącego
    rdzenia i wskazania hosta (DANACO_HOST_RDZENIA albo wskazanie w oknie),
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie wymagają maszyny z Windows na ARM64. Nie zakładaj ich powodzenia.
KONIEC
