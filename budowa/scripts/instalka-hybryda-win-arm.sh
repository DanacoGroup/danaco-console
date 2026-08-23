#!/usr/bin/env bash
# Instalka HYBRYDOWA Windows na ARM64 — złożenie pliku .exe (NSIS) na Linuksie.
#
# Czym ta instalka różni się od `instalka-windows.sh`: tamta wiezie CAŁĄ
# platformę (rdzeń, serwer narzędzi, pomocniki) i instaluje ją na maszynie
# Operatora. Ta wiezie SAMĄ POWŁOKĘ i cztery dokumenty produktu. Rdzeń, serwer
# narzędzi i arsenał stoją na serwerze Danaco — Operator dostaje okno, nie
# drugą kopię serca platformy. Dlatego jest cienka i dlatego nie sprawdza
# obecności `zasoby/rdzen/*`: ich brak nie jest tu usterką, jest założeniem.
#
# ── Droga budowy dla ARM64 ────────────────────────────────────────────────
# Rustup zna dwa cele Windows/ARM64: `aarch64-pc-windows-msvc` i
# `aarch64-pc-windows-gnullvm`. Ten skrypt idzie drogą `gnullvm`, bo:
#   * nie wymaga zestawu nagłówków Microsoftu ani zgody na jego licencję
#     (droga MSVC przez cargo-xwin wymaga jednego i drugiego),
#   * łańcuch narzędzi jest jednym archiwum bez instalatora.
# Apt NIE ma mingw-w64 dla aarch64 — jest tylko i686/x86_64. Łańcuch bierze się
# z gotowego wydania llvm-mingw (clang + lld + nagłówki i CRT dla
# aarch64-w64-mingw32) i rozpakowuje do $LLVM_MINGW. Ustawienie PATH wystarcza:
# domyślnym konsolidatorem celu gnullvm jest `aarch64-w64-mingw32-clang`,
# który z tego archiwum pochodzi.
#
# Skrypt woła `cargo tauri bundle`, nie `build` — z tego samego powodu, co
# `instalka-windows.sh`: `build` uruchomiłby beforeBuildCommand, czyli
# przebudowę klienta, która do budowy powłoki nie należy.
#
# Kształt produktu opisuje `tauri.hybryda-win-arm.conf.json`, dokładany przez
# `--config`. Tauri scala go z `tauri.conf.json` regułą JSON Merge Patch
# (RFC 7386), w której wartość `null` USUWA klucz — tak wypisane są z listy
# zasobów rdzeń, serwer narzędzi i pomocniki. Skrypt nie wierzy, że scalenie
# zadziałało: po złożeniu sprawdza listing i odmawia, jeśli rdzeń w środku jest.
#
# Skrypt nie podpisuje instalatora — podpis Authenticode wymaga certyfikatu
# i hosta Windows — i nie sprawdza, czy instalator się uruchamia; tego nie da
# się sprawdzić bez maszyny z Windows na ARM64.
#
# Użycie: bash budowa/scripts/instalka-hybryda-win-arm.sh
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
CEL="aarch64-pc-windows-gnullvm"
KONFIGURACJA="$POWLOKA/tauri.hybryda-win-arm.conf.json"
LLVM_MINGW="${LLVM_MINGW:-/opt/llvm-mingw}"
WYDANIE="$KORZEN/budowa/wydania/1.0.0-2026-08-18"
NAZWA="Danaco Console_1.0.0_hybryda_arm64-setup.exe"

# Rustup instaluje się do ~/.cargo, którego nie ma w PATH powłoki
# nieinteraktywnej. Bez tej linii `cargo` jest nieodnajdywalny.
# shellcheck disable=SC1091
[ -f "$HOME/.cargo/env" ] && . "$HOME/.cargo/env"
export PATH="$LLVM_MINGW/bin:$PATH"

# Skrypty budowania zależności C (crate `cc`) nie znają celu gnullvm z nazwy
# i sięgnęłyby po linuksowy `cc`. Wskazanie kompilatora per cel jest jawne.
export CC_aarch64_pc_windows_gnullvm=aarch64-w64-mingw32-clang
export CXX_aarch64_pc_windows_gnullvm=aarch64-w64-mingw32-clang++
export AR_aarch64_pc_windows_gnullvm=llvm-ar

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

zglos "Sprawdzenie narzędzi"
# Brak któregokolwiek z nich ujawniłby się dopiero po kilkuminutowej budowie,
# dlatego sprawdzenie stoi przed nią, nie po.
for narzedzie in cargo rustup makensis aarch64-w64-mingw32-clang; do
  command -v "$narzedzie" >/dev/null \
    || padnij "brak narzędzia: $narzedzie (apt: nsis; rustup: https://rustup.rs; łańcuch aarch64: wydanie llvm-mingw rozpakowane do $LLVM_MINGW — apt NIE ma mingw-w64 dla aarch64)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"
printf '  jest: cel %s\n' "$CEL"
[ -f "$KONFIGURACJA" ] || padnij "brak konfiguracji produktu: $KONFIGURACJA"
printf '  jest: %s\n' "${KONFIGURACJA#"$KORZEN"/}"

zglos "Sprawdzenie artefaktów wejściowych"
# Powłoka hybrydowa niesie własne okno, więc zbudowany klient jest jej
# jedynym artefaktem wejściowym. Rdzenia nie sprawdzamy — tej instalki nie ma.
[ -e "$KORZEN/budowa/client/dist/index.html" ] \
  || padnij "brak artefaktu: budowa/client/dist/index.html (zbuduj: npm run build w budowa/client)"
printf '  jest: budowa/client/dist/index.html\n'

zglos "Dokumentacja produktu"
# Instalator ma nieść dokumentację, bo to ON jest produktem końcowym: Operator
# dostaje jeden plik i po instalacji nie ma skąd wziąć instrukcji ani licencji.
# Pliki mają jedno miejsce — korzeń repozytorium — i są tam poprawiane; kopia
# powstaje tuż przed złożeniem, żeby instalka nie wiozła wersji z poprzedniego
# tygodnia.
#
# Katalog jest WŁASNY dla tego wariantu, nie wspólny `zasoby/dokumentacja`.
# Powód jest praktyczny: instalki x64 i ARM64 składa się równolegle, a wspólny
# katalog jest czyszczony przez `rm -rf` na starcie każdej z nich — jedna
# budowa wyrywałaby drugiej pliki spod rąk.
DOKUMENTACJA="$POWLOKA/zasoby/dokumentacja-hybryda-arm64"
rm -rf "$DOKUMENTACJA"
mkdir -p "$DOKUMENTACJA"
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  [ -f "$KORZEN/$dokument" ] || padnij "brak dokumentu produktu: $dokument — instalator bez licencji i instrukcji nie jest produktem końcowym"
  cp -f "$KORZEN/$dokument" "$DOKUMENTACJA/$dokument"
  printf '  spakowany: %s (%s)\n' "$dokument" "$(du -h "$KORZEN/$dokument" | cut -f1)"
done

zglos "Budowa powłoki dla $CEL"
# `--features tauri/custom-protocol` NIE jest ozdobą — bez niej produkt jest
# zepsuty w sposób niewidoczny. W tauri 2 `tauri::is_dev()` to dokładnie
# `!cfg!(feature = "custom-protocol")`, a `Cargo.toml` powłoki tej cechy nie
# włącza i `cargo build` sam jej nie doda (`cargo tauri bundle` też nie —
# on niczego nie kompiluje, bierze gotową binarkę z target/). Bez cechy
# `zrodlo_interfejsu::ustal` zatrzymuje się na `is_dev()` i idzie do devUrl,
# NIGDY nie sprawdzając warunku „rdzeń nasłuchujący" — a dla hybrydy to jedyna
# droga do interfejsu, bo rdzeń stoi na serwerze. Okno pokazałoby pustą stronę.
# Cechę podajemy z wiersza poleceń; Cargo.toml zostaje nietknięty, bo to samo
# źródło buduje też wariant pełny.
#
# TAURI_CONFIG dokłada tę samą nakładkę, którą niżej dostaje `bundle`. Bez tego
# kompilacja widziałaby `tauri.conf.json` (z devUrl i z rdzeniem na liście
# zasobów), a złożenie — nakładkę: dwa różne opisy jednego produktu w jednej
# budowie. Przy okazji usunięty `devUrl` nie zostaje wkompilowany w binarkę.
( cd "$POWLOKA" \
  && TAURI_CONFIG="$(cat "$KONFIGURACJA")" \
     cargo build --release --target "$CEL" --features tauri/custom-protocol )

POWLOKA_EXE="$POWLOKA/target/$CEL/release/danaco-console-powloka.exe"
[ -f "$POWLOKA_EXE" ] || padnij "cargo nie zgłosił błędu, ale binarki powłoki nie ma: $POWLOKA_EXE"
# Cel budowy da się pomylić w wierszu poleceń, a instalka z binarką x64 w środku
# jest dla maszyny ARM64 bezużyteczna. To jedyny sprawdzian, który tu naprawdę
# coś znaczy, więc stoi przed złożeniem, nie tylko po nim.
# `file` nazywa tę architekturę „ARM64", a nie „Aarch64" jak rustup — wzorzec
# dopuszcza obie pisownie, bo to nazwa zależna od wersji `file`, nie od produktu.
OPIS_POWLOKI="$(file -b "$POWLOKA_EXE")"
printf '%s' "$OPIS_POWLOKI" | grep -qi 'PE32+ executable' \
  || padnij "powłoka nie jest binarką Windows (PE32+): $OPIS_POWLOKI"
printf '%s' "$OPIS_POWLOKI" | grep -qiE 'ARM64|Aarch64' \
  || padnij "powłoka nie jest binarką PE dla ARM64: $OPIS_POWLOKI"
printf '  architektura powłoki: %s\n' "$OPIS_POWLOKI"

zglos "Zapora: osadzone zasoby interfejsu"
# Po czym POZNAĆ, że cecha custom-protocol weszła. Nie po napisie
# `localhost:5173`: ten siedzi w binarce ZAWSZE, bo Tauri wkompilowuje w nią
# całą konfigurację produktu wraz z polem devUrl, więc jego obecność nie mówi
# nic o tym, która droga jest czynna. Zapora na ten napis zapala się przy pliku
# poprawnym i milczy przy zepsutym — jest gorsza niż jej brak.
#
# Rozstrzyga obecność OSADZONYCH ZASOBÓW: wchodzą do binarki wyłącznie przy
# czynnej cesze custom-protocol, czyli tylko gdy budowa nie jest rozwojowa.
# Sondą jest nazwa pliku interfejsu z sumą treści w nazwie (np.
# `index-C8csBLsr.js`) — taki napis nie ma jak trafić do binarki inaczej niż
# przez osadzenie `client/dist`, i zmienia się z każdą przebudową klienta,
# dlatego czytamy go z katalogu, a nie wpisujemy tu na sztywno.
#
# Mylące sondy, których tu NIE używamy: `woff2` (tauri ma tablicę typów MIME,
# jedno trafienie jest zawsze) i `index.html` (tauri obsługuje indeks katalogu,
# trafień jest kilka) — obie zapalają się w binarce BEZ osadzonych zasobów.
SONDA="$(find "$KORZEN/budowa/client/dist/assets" -maxdepth 1 -name 'index-*.js' -printf '%f\n' 2>/dev/null | head -1)"
[ -n "$SONDA" ] \
  || padnij "nie ma z czego zrobić sondy: brak budowa/client/dist/assets/index-*.js — bez niej nie da się sprawdzić, czy zasoby są osadzone"
printf '  sonda: %s\n' "$SONDA"
LICZBA_SONDY="$(strings -a "$POWLOKA_EXE" | grep -c -- "$SONDA" || true)"
LICZBA_ASSETS="$(strings -a "$POWLOKA_EXE" | grep -c '/assets/' || true)"
printf '  trafienia sondy: %s, wystąpienia /assets/: %s\n' "$LICZBA_SONDY" "$LICZBA_ASSETS"
[ "$LICZBA_SONDY" -gt 0 ] && [ "$LICZBA_ASSETS" -gt 0 ] \
  || padnij "powłoka NIE ma osadzonych zasobów interfejsu — cecha tauri/custom-protocol nie weszła, więc tauri::is_dev() zwróci prawdę i okno pójdzie do serwera rozwojowego zamiast do rdzenia na serwerze. To produkt zepsuty, choć instalka wyglądałaby poprawnie"
printf '  osadzone zasoby interfejsu: potwierdzone\n'

zglos "Złożenie instalatora NSIS"
# `--bundles nsis` podane jawnie, choć konfiguracja ma ten cel ustawiony:
# skrypt składa instalator NSIS niezależnie od tego, jakie inne cele
# konfiguracja produktu wymienia.
( cd "$POWLOKA" && cargo tauri bundle --target "$CEL" --bundles nsis --config "$KONFIGURACJA" )

ZLOZONY="$(find "$POWLOKA/target/$CEL/release/bundle/nsis" -maxdepth 1 -name '*-setup.exe' -print -quit)"
[ -n "$ZLOZONY" ] || padnij "makensis nie zgłosił błędu, ale pliku instalatora nie ma"

zglos "Sprawdzenie, że instalka jest CIENKA"
# Scalenie konfiguracji przez `null` jest cichym mechanizmem: gdyby Tauri
# przestało honorować RFC 7386, rdzeń wróciłby do środka bez jednego ostrzeżenia,
# a produkt przestałby być tym, co obiecuje nazwa. Stąd sprawdzenie listingu.
if command -v 7z >/dev/null; then
  LISTING="$(7z l "$ZLOZONY" 2>/dev/null || true)"
  for zakazane in danaco-console.exe danaco-narzedzia.exe; do
    [ "$(printf '%s' "$LISTING" | grep -c "$zakazane" || true)" -gt 0 ] \
      && padnij "instalka wiezie $zakazane — to nie jest wariant hybrydowy, scalenie konfiguracji nie usunęło rdzenia"
  done
  printf '  brak rdzenia: potwierdzony\n'
  for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
    [ "$(printf '%s' "$LISTING" | grep -c "$dokument" || true)" -gt 0 ] \
      || padnij "instalka nie wiezie dokumentu $dokument"
  done
  printf '  cztery dokumenty: obecne\n'
  # Powyższe zapory badały binarkę z target/. Ta bada binarkę WYPAKOWANĄ
  # z instalatora, bo tylko ona jest tym, co dostanie Operator: gdyby `bundle`
  # wziął binarkę z poprzedniej budowy (a bierze gotową z target/, sam nie
  # kompiluje), różnicy nie byłoby widać inaczej.
  ROZPAK="$(mktemp -d)"
  trap 'rm -rf "$ROZPAK"' EXIT
  7z x -o"$ROZPAK" "$ZLOZONY" >/dev/null 2>&1 \
    || padnij "nie udało się wypakować instalatora do sprawdzenia"
  WIEZIONA="$ROZPAK/danaco-console-powloka.exe"
  [ -f "$WIEZIONA" ] || padnij "w instalatorze nie ma powłoki danaco-console-powloka.exe"
  OPIS_WIEZIONEJ="$(file -b "$WIEZIONA")"
  printf '%s' "$OPIS_WIEZIONEJ" | grep -qiE 'ARM64|Aarch64' \
    || padnij "powłoka WIEZIONA przez instalator nie jest binarką ARM64: $OPIS_WIEZIONEJ"
  # `strings … | grep -q` byłoby tu PUŁAPKĄ przy `set -o pipefail`: grep -q
  # kończy się na pierwszym trafieniu, `strings` dostaje SIGPIPE (141),
  # a pipefail zamienia to w porażkę całego potoku — zapora zapalałaby się
  # DOKŁADNIE wtedy, gdy sonda się znajdzie. Dlatego liczymy trafienia
  # w podstawieniu polecenia, gdzie `strings` dobiega do końca.
  TRAFIENIA_WIEZIONEJ="$(strings -a "$WIEZIONA" | grep -c -- "$SONDA" || true)"
  [ "$TRAFIENIA_WIEZIONEJ" -gt 0 ] \
    || padnij "powłoka WIEZIONA przez instalator nie ma osadzonych zasobów interfejsu — instalka wiezie budowę rozwojową"
  printf '  wieziona powłoka: %s, zasoby osadzone\n' "$OPIS_WIEZIONEJ"
else
  printf '  nie sprawdzono (brak 7z — apt: p7zip-full)\n'
fi

zglos "Odłożenie wyniku do wydania"
mkdir -p "$WYDANIE"
WYNIK="$WYDANIE/$NAZWA"
cp -f "$ZLOZONY" "$WYNIK"

zglos "Pomiar wyniku"
printf 'ścieżka : %s\n' "$WYNIK"
printf 'rozmiar : %s\n' "$(du -h "$WYNIK" | cut -f1)"
printf 'typ     : %s\n' "$(file -b "$WYNIK")"
printf 'suma    : %s\n' "$(sha256sum "$WYNIK" | cut -d' ' -f1)"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia i czy przechodzi do końca,
  - czy zainstalowana aplikacja wstaje (wymaga WebView2 dla ARM64 w systemie),
  - czy powłoka dogaduje się z serwerem Danaco, na którym stoi rdzeń,
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie wymagają maszyny z Windows na ARM64. Nie zakładaj ich powodzenia.
KONIEC
