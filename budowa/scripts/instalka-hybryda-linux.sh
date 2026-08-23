#!/usr/bin/env bash
# Instalka HYBRYDOWA dla Linuksa — okno klienta Danaco Console.
#
# Co składa: pakiety .deb i .AppImage z powłoką natywną i czterema dokumentami
# produktu z korzenia repozytorium. Bez rdzenia, bez serwera narzędzi, bez
# arsenału — rdzeń pracuje na serwerze, okno łączy się z nim po sieci.
#
# Dwie rzeczy, których pominięcie daje pakiet wyglądający na dobry i martwy
# u użytkownika:
#
#  1. Cecha `tauri/custom-protocol`. `tauri::is_dev()` to dokładnie
#     `!cfg!(feature = "custom-protocol")`; Cargo.toml powłoki tej cechy nie
#     włącza, a `cargo tauri bundle` niczego nie kompiluje — bierze gotową
#     binarkę z `target/`. Powłoka zbudowana bez cechy idzie do `devUrl`
#     (localhost:5173), nigdy nie sprawdza warunku „rdzeń nasłuchujący" i nie
#     niesie osadzonych zasobów interfejsu. Dlatego skrypt kompiluje sam,
#     z cechą podaną z wiersza poleceń, i tylko potem woła bundler.
#
#  2. Zapora odbioru. Sondą jest nazwa pliku interfejsu z sumą treści w nazwie
#     (`client/dist/assets/index-*.js`) — ta nie ma jak trafić do binarki inaczej
#     niż przez osadzenie pakietu. Sondy `index.html`, `<!doctype html>` ani
#     `woff2` do niczego nie służą: zapalają się w powłoce bez osadzonych
#     zasobów (tablica typów MIME, obsługa indeksu katalogu). Zapora liczy
#     trafienia w powłoce WYPAKOWANEJ Z GOTOWEGO PAKIETU, nie w `target/`.
#
# Profil bundlera: budowa/desktop/src-tauri/tauri.hybryda-linux.conf.json.
# Wyniki: budowa/wydania/<wersja>-<data>/.
#
# Wywołanie:
#   budowa/scripts/instalka-hybryda-linux.sh            # deb + appimage
#   budowa/scripts/instalka-hybryda-linux.sh deb        # tylko deb
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
PROFIL="tauri.hybryda-linux.conf.json"
CELE="${1:-deb,appimage}"

# Nazwa pakietu Debiana. Bundler składa ją z `productName`, co dawało
# `danaco-console` — nazwę zajętą przez pakiet serwera; dwa pakiety o jednej
# nazwie nadpisują się na maszynie użytkownika. Twarz produktu zostaje
# „Danaco Console"; rozdzielenie siedzi wyłącznie w polu Package.
PAKIET_DEB="danaco-console-okno"

WERSJA="$(sed -n 's/.*"version": "\([^"]*\)".*/\1/p' "$POWLOKA/tauri.conf.json" | head -1)"
DATA="$(date +%Y-%m-%d)"
WYDANIE="$KORZEN/budowa/wydania/$WERSJA-$DATA"
BUNDLE="$POWLOKA/target/release/bundle"

powiedz() { printf '\n== %s\n' "$1"; }
odmow() { printf '\nODMOWA WYDANIA: %s\n' "$1" >&2; exit 1; }

powiedz "Powłoka hybrydowa $WERSJA ($DATA), cele: $CELE"

# Rust bywa poza ścieżką powłoki nieinteraktywnej.
if [ -f "$HOME/.cargo/env" ]; then
  # shellcheck disable=SC1091
  . "$HOME/.cargo/env"
fi

# Cztery dokumenty produktu muszą leżeć w korzeniu — profil wskazuje je wprost.
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  [ -f "$KORZEN/$dokument" ] || odmow "brak dokumentu produktu: $dokument"
done

# Sondę czyta się przy każdej budowie — nazwa zmienia się z treścią pakietu.
SONDA="$(cd "$KORZEN/budowa/client/dist/assets" 2>/dev/null && ls -1 index-*.js 2>/dev/null | grep -v '\.map$' | head -1 || true)"
[ -n "$SONDA" ] || odmow "brak pakietu interfejsu (budowa/client/dist/assets/index-*.js) — nie ma czego osadzić ani czym sprawdzać"
echo "sonda osadzenia: $SONDA"

# Liczy trafienia sondy we wskazanej binarce. Trafienia liczy się w podstawieniu
# polecenia, nie przez `grep -q`: przy `set -o pipefail` zamknięcie potoku po
# pierwszym trafieniu wywraca `strings` sygnałem SIGPIPE, czyli krok padałby
# dokładnie wtedy, gdy sonda się ZNAJDUJE.
trafienia() {
  strings -a "$1" 2>/dev/null | grep -c -- "$SONDA" || true
}

# Druga zapora, na treści. Sonda z nazwy pliku łapie powłokę bez zasobów i
# powłokę z zasobami przeterminowanymi, ale przepuściłaby powłokę zbudowaną bez
# cechy `custom-protocol` z pakietem interfejsu ŚWIEŻYM: wykaz kluczy zasobów
# powstaje tak samo w obu trybach, a różni je to, czy do binarki wchodzi TREŚĆ.
# Dlatego drugi krok szuka w powłoce dosłownych bajtów największego zasobu
# skompresowanego przez codegen (`target/release/build/*/out/tauri-codegen-assets`).
# Pliku `.map` się nie liczy — codegen zostawia jego klucz bez treści.
tresc_osadzona() {
  python3 - "$1" <<'PY'
import glob, os, sys
powloka = open(sys.argv[1], 'rb').read()
katalogi = glob.glob('target/release/build/danaco-console-powloka-*/out/tauri-codegen-assets')
if not katalogi:
    print('BRAK-KATALOGU-ZASOBOW'); sys.exit(0)
katalog = max(katalogi, key=os.path.getmtime)
zasoby = sorted(
    ((os.path.getsize(os.path.join(katalog, n)), n) for n in os.listdir(katalog)),
    reverse=True,
)
duze = [(r, n) for r, n in zasoby if r >= 500_000]
if not duze:
    print('BRAK-DUZYCH-ZASOBOW'); sys.exit(0)
for rozmiar, nazwa in duze:
    dane = open(os.path.join(katalog, nazwa), 'rb').read()
    if dane[:4096] in powloka:
        print(f'OSADZONA {rozmiar}')
        sys.exit(0)
print('BRAK-TRESCI')
PY
}

powiedz "Kompilacja powłoki (cecha tauri/custom-protocol)"
cd "$POWLOKA"
cargo build --release --features tauri/custom-protocol

powiedz "Bundler"
rm -rf "$BUNDLE/deb" "$BUNDLE/appimage"
cargo tauri bundle --bundles "$CELE" --config "$PROFIL"

powiedz "Nazwa pakietu Debiana: $PAKIET_DEB"
ROBOCZY="$(mktemp -d)"
trap 'rm -rf "$ROBOCZY"' EXIT
while IFS= read -r -d '' pakiet; do
  rozpakowany="$ROBOCZY/przepakowanie"
  rm -rf "$rozpakowany"
  dpkg-deb -R "$pakiet" "$rozpakowany"
  sed -i "s/^Package: .*/Package: $PAKIET_DEB/" "$rozpakowany/DEBIAN/control"
  dpkg-deb --build "$rozpakowany" "$pakiet" >/dev/null
done < <(find "$BUNDLE" -maxdepth 2 -name '*.deb' -print0 2>/dev/null)

powiedz "Zapora odbioru: osadzenie interfejsu w gotowych pakietach"
mkdir -p "$WYDANIE"

# Nazwy odróżniają wariant dostawy: hybryda (okno klienta) od pakietu serwera.
znalezione=0
while IFS= read -r -d '' pakiet; do
  znalezione=$((znalezione + 1))
  nazwa="$(basename "$pakiet")"
  cel="$WYDANIE/${nazwa/_${WERSJA}_/_${WERSJA}_hybryda_}"

  case "$nazwa" in
    *.deb)
      rozpakowany="$ROBOCZY/odbior-deb"
      rm -rf "$rozpakowany"
      dpkg-deb -x "$pakiet" "$rozpakowany"
      powloka="$rozpakowany/usr/bin/danaco-console-powloka"
      ;;
    *.AppImage)
      # AppImage sprawdza się sobą: rozpakowanie nie wymaga montowania.
      rozpakowany="$ROBOCZY/odbior-appimage"
      rm -rf "$rozpakowany"
      mkdir -p "$rozpakowany"
      ( cd "$rozpakowany" && "$pakiet" --appimage-extract >/dev/null 2>&1 ) || true
      powloka="$(find "$rozpakowany" -type f -name 'danaco-console-powloka' | head -1)"
      # Bez rozpakowania sondę liczy się w samym pliku obrazu (squashfs bywa
      # nieskompresowany dla dużych zasobów) — słabsze, ale nie zerowe.
      [ -n "$powloka" ] || powloka="$pakiet"
      ;;
  esac

  [ -f "$powloka" ] || odmow "$nazwa — nie znalazłem powłoki w pakiecie"
  ile="$(trafienia "$powloka")"
  if [ "$ile" -eq 0 ]; then
    odmow "$nazwa — sonda $SONDA nie występuje w powłoce; pakiet nie niesie osadzonego interfejsu (powłoka zbudowana bez cechy tauri/custom-protocol)"
  fi
  echo "$nazwa — trafienia sondy: $ile"

  wynik_tresci="$(cd "$POWLOKA" && tresc_osadzona "$powloka")"
  case "$wynik_tresci" in
    OSADZONA*)
      echo "$nazwa — treść interfejsu w powłoce: ${wynik_tresci#OSADZONA }" ;;
    BRAK-TRESCI)
      odmow "$nazwa — powłoka niesie klucze zasobów, ale nie ich treść; zbudowana bez cechy tauri/custom-protocol" ;;
    *)
      echo "$nazwa — zapory treści nie dało się rozstrzygnąć ($wynik_tresci)" ;;
  esac

  cp -f "$pakiet" "$cel"
  sha256sum "$(basename "$cel")" >/dev/null 2>&1 || true
done < <(find "$BUNDLE" -maxdepth 2 \( -name '*.deb' -o -name '*.AppImage' \) -print0 2>/dev/null)

[ "$znalezione" -gt 0 ] || odmow "bundler nie zostawił żadnego pakietu"

powiedz "Odbiór"
cd "$WYDANIE"
shopt -s nullglob
for artefakt in *hybryda*.deb *hybryda*.AppImage; do
  sha256sum "$artefakt" >"$artefakt.sha256"
  printf '%s  %s  %s\n' "$(du -h "$artefakt" | cut -f1)" "$(file -b "$artefakt" | cut -c1-40)" "$artefakt"
done
for artefakt in *hybryda*.deb; do
  dpkg-deb -f "$artefakt" Package Version Installed-Size Depends
done

powiedz "Wynik w $WYDANIE"
