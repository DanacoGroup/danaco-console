#!/usr/bin/env bash
# Instalka CAŁKOWICIE NATYWNA dla Windows — złożenie pliku .exe (NSIS) na Linuksie.
#
# ── Czym ta instalka różni się od `instalka-windows.sh` ───────────────────────
# Tamta wiezie powłokę, rdzeń, interfejs i dokumentację, czyli sam produkt.
# Ta wiezie ponadto ARSENAŁ: programy spoza instalki, których rdzeń używa do
# obrazu, filmu, dźwięku, dokumentu i archiwum, w wydaniach dla Windows. Bez
# nich Operator dostaje produkt, który na kilkanaście czynności odpowiada
# uprzejmą odmową „brak programu…" — i nie ma skąd tych programów wziąć, bo na
# Windowsie nie ma `apt`.
#
# ── Dlaczego arsenał leży w `rdzen/pomocniki`, a nie gdzie indziej ────────────
# `server/internal/zewnetrzne/odnajdywanie.go` szuka binarium arsenału najpierw
# w katalogu `pomocniki` OBOK pracującego pliku wykonywalnego, dopiero potem na
# ścieżce systemu. Pracującym plikiem jest rdzeń, a rdzeń jedzie do
# `rdzen/danaco-console.exe`, więc arsenał musi wylądować w `rdzen/pomocniki/`.
# Ta sama zasada każe trzymać serwer narzędzi obok rdzenia
# (`server/internal/narzedzia/wpiecie.go`).
#
# ── Dwa układy arsenału: płaski i własny katalog ──────────────────────────────
# Programy jednoplikowe (i te niosące kilka własnych bibliotek) leżą PŁASKO
# w `pomocniki/<program>.exe` — Windows szuka DLL w katalogu pliku wykonywalnego,
# więc program znajduje swoje biblioteki obok siebie. Programy niosące CAŁE
# WŁASNE ŚRODOWISKO (pwsh, LibreOffice, Chromium, rembg z Pythonem) dostają
# WŁASNY KATALOG `pomocniki/<program>/…`, bo zsypanie ich bibliotek do wspólnego
# `pomocniki` znaczyłoby kolizję nazw. Oba układy obsługuje `zewnetrzne.Odnajdz`
# (`server/internal/zewnetrzne/odnajdywanie.go`, funkcja `miejscaPakietu`).
#
# ── Czego ten skrypt NIE udaje ───────────────────────────────────────────────
# Jedyna deklarowana zależność bez drogi przenośnej to SILNIK KONTENERÓW: na
# Windowsie wymaga WSL2 i uprawnień administratora, więc nie da się go położyć
# obok rdzenia. Skrypt go nie pobiera i nie podszywa się pod niego — rdzeń sam
# nazywa ten brak zdaniem `zewnetrzne.BrakNarzedzia`. Reszta arsenału (w tym
# LibreOffice, Chromium, rembg, pwsh) JEST w pakiecie.
#
# Użycie: bash budowa/scripts/instalka-natywna-win.sh
set -euo pipefail

KORZEN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
POWLOKA="$KORZEN/budowa/desktop/src-tauri"
ARSENAL="$POWLOKA/zasoby/arsenal-windows"
POBRANE="$ARSENAL/_pobrane"
ROZPAKOWANE="$ARSENAL/_rozpakowane"
POMOCNIKI="$ARSENAL/pomocniki"
NASTAWA="$POWLOKA/tauri.natywna-win.conf.json"
CEL="x86_64-pc-windows-gnu"
WERSJA="1.0.0"
WYDANIE="$KORZEN/budowa/wydania/$WERSJA-$(date +%Y-%m-%d)"

# shellcheck disable=SC1091
[ -f "$HOME/.cargo/env" ] && . "$HOME/.cargo/env"
command -v go >/dev/null || PATH="$PATH:/usr/local/go/bin"

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() { printf 'ODMOWA: %s\n' "$1" >&2; exit 1; }

# wersje arsenału są przypięte w jednym miejscu, bo raport odbioru ma mówić,
# co dokładnie pojechało do Operatora — „najnowsze" nie jest odpowiedzią.
IM_WERSJA="7.1.2-29"
PANDOC_WERSJA="3.10.2"
POPPLER_WERSJA="26.02.0-0"
TESSERACT_WERSJA="5.4.0.20240606"
SIEDEMZIP_WERSJA="2501"
PIPER_WERSJA="2023.11.14-2"
ESRGAN_WERSJA="v0.2.5.0"
ESPEAK_WERSJA="1.52.0"
PWSH_WERSJA="7.6.5"
LIBREOFFICE_WERSJA="26.2.5"
CHROMIUM_WERSJA="151.0.7922.137-1.1"
PYTHON_WERSJA="3.11.9"
ONNXRUNTIME_WERSJA="1.29.0"

zglos "Sprawdzenie narzędzi"
for narzedzie in cargo rustup x86_64-w64-mingw32-gcc makensis go npm curl unzip 7z; do
  command -v "$narzedzie" >/dev/null \
    || padnij "brak narzędzia: $narzedzie (apt: mingw-w64, nsis, p7zip-full, unzip; rustup: https://rustup.rs)"
  printf '  jest: %s\n' "$narzedzie"
done
rustup target list --installed | grep -qx "$CEL" \
  || padnij "brak celu Rusta $CEL — dodaj: rustup target add $CEL"

zglos "Rdzeń i serwer narzędzi dla Windows"
# CGO jest włączone, bo rdzeń wiąże biblioteki C; bez `CC` wskazującego mingw
# kompilator próbowałby budować pod Linuksa i przewrócił się na nagłówkach.
# Oba pliki lądują w jednym katalogu, bo rdzeń szuka serwera narzędzi obok siebie.
mkdir -p "$POWLOKA/zasoby/rdzen"
for program in danaco-console danaco-narzedzia; do
  CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
    go build -C "$KORZEN/budowa" -o "$POWLOKA/zasoby/rdzen/$program.exe" "./server/cmd/$program"
  file -b "$POWLOKA/zasoby/rdzen/$program.exe" | grep -q 'PE32+ executable' \
    || padnij "$program.exe nie jest binarką Windows (PE32+)"
  printf '  zbudowany: %s (%s)\n' "$program.exe" "$(du -h "$POWLOKA/zasoby/rdzen/$program.exe" | cut -f1)"
done

zglos "Interfejs"
npm --prefix "$KORZEN/budowa/client" run build
# Rdzeń podaje interfejs z katalogu obok siebie, więc kopia jedzie do zasobu rdzenia.
rm -rf "$POWLOKA/zasoby/rdzen/dist"
cp -r "$KORZEN/budowa/client/dist" "$POWLOKA/zasoby/rdzen/dist"

zglos "Dokumentacja produktu"
# Katalog zasobu jest WYTWOREM budowania: pliki mają jedno miejsce w korzeniu
# repozytorium, a kopia powstaje tuż przed złożeniem, żeby instalka nie wiozła
# wersji z poprzedniego tygodnia.
DOKUMENTACJA="$POWLOKA/zasoby/dokumentacja"
rm -rf "$DOKUMENTACJA"; mkdir -p "$DOKUMENTACJA"
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  [ -f "$KORZEN/$dokument" ] \
    || padnij "brak dokumentu produktu: $dokument — instalator bez licencji i instrukcji nie jest produktem końcowym"
  cp -f "$KORZEN/$dokument" "$DOKUMENTACJA/$dokument"
  printf '  spakowany: %s\n' "$dokument"
done

zglos "Pobranie arsenału w wydaniach Windows"
mkdir -p "$POBRANE"
# Pobranie jest pomijane, gdy plik już leży — powtórne złożenie instalki nie ma
# ciągnąć z sieci gigabajta, który się nie zmienił.
pobierz() {
  local plik="$1" adres="$2"
  if [ -s "$POBRANE/$plik" ]; then
    printf '  jest już: %-28s %s\n' "$plik" "$(du -h "$POBRANE/$plik" | cut -f1)"
    return 0
  fi
  if curl -sSL --retry 2 --max-time 600 -o "$POBRANE/$plik" "$adres" && [ -s "$POBRANE/$plik" ]; then
    printf '  pobrany : %-28s %s\n' "$plik" "$(du -h "$POBRANE/$plik" | cut -f1)"
  else
    rm -f "$POBRANE/$plik"
    printf '  NIE POBRANY: %s (%s) — produkt powie o braku sam\n' "$plik" "$adres"
    return 1
  fi
}

pobierz imagemagick.7z "https://github.com/ImageMagick/ImageMagick/releases/download/$IM_WERSJA/ImageMagick-$IM_WERSJA-portable-Q16-HDRI-x64.7z" || true
pobierz ffmpeg.zip "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip" || true
pobierz tesseract-setup.exe "https://github.com/UB-Mannheim/tesseract/releases/download/v$TESSERACT_WERSJA/tesseract-ocr-w64-setup-$TESSERACT_WERSJA.exe" || true
for jezyk in pol eng osd; do
  pobierz "$jezyk.traineddata" "https://github.com/tesseract-ocr/tessdata_fast/raw/main/$jezyk.traineddata" || true
done
pobierz 7z-inst.exe "https://www.7-zip.org/a/7z$SIEDEMZIP_WERSJA-x64.exe" || true
pobierz pandoc.zip "https://github.com/jgm/pandoc/releases/download/$PANDOC_WERSJA/pandoc-$PANDOC_WERSJA-windows-x86_64.zip" || true
pobierz poppler.zip "https://github.com/oschwartz10612/poppler-windows/releases/download/v$POPPLER_WERSJA/Release-$POPPLER_WERSJA.zip" || true
pobierz piper.zip "https://github.com/rhasspy/piper/releases/download/$PIPER_WERSJA/piper_windows_amd64.zip" || true
pobierz pl_PL-gosia-medium.onnx "https://huggingface.co/rhasspy/piper-voices/resolve/main/pl/pl_PL/gosia/medium/pl_PL-gosia-medium.onnx" || true
pobierz pl_PL-gosia-medium.onnx.json "https://huggingface.co/rhasspy/piper-voices/resolve/main/pl/pl_PL/gosia/medium/pl_PL-gosia-medium.onnx.json" || true
pobierz esrgan.zip "https://github.com/xinntao/Real-ESRGAN/releases/download/$ESRGAN_WERSJA/realesrgan-ncnn-vulkan-20220424-windows.zip" || true
pobierz espeak.msi "https://github.com/espeak-ng/espeak-ng/releases/download/$ESPEAK_WERSJA/espeak-ng.msi" || true
pobierz pwsh.zip "https://github.com/PowerShell/PowerShell/releases/download/v$PWSH_WERSJA/PowerShell-$PWSH_WERSJA-win-x64.zip" || true
pobierz libreoffice.msi "https://download.documentfoundation.org/libreoffice/stable/$LIBREOFFICE_WERSJA/win/x86_64/LibreOffice_${LIBREOFFICE_WERSJA}_Win_x86-64.msi" || true
pobierz chromium.zip "https://github.com/ungoogled-software/ungoogled-chromium-windows/releases/download/$CHROMIUM_WERSJA/ungoogled-chromium_${CHROMIUM_WERSJA}_windows_x64.zip" || true
pobierz python-embed.zip "https://www.python.org/ftp/python/$PYTHON_WERSJA/python-$PYTHON_WERSJA-embed-amd64.zip" || true
pobierz u2net.onnx "https://github.com/danielgatis/rembg/releases/download/v0.0.0/u2net.onnx" || true

zglos "Rozpakowanie arsenału"
rm -rf "$ROZPAKOWANE"; mkdir -p "$ROZPAKOWANE"
# Instalator Tesseractu i instalator 7-Zipa są archiwami NSIS, więc `7z x`
# wyjmuje z nich pliki bez uruchamiania czegokolwiek pod Windows.
[ -s "$POBRANE/imagemagick.7z" ]     && 7z x -y -o"$ROZPAKOWANE/im" "$POBRANE/imagemagick.7z" >/dev/null
[ -s "$POBRANE/tesseract-setup.exe" ] && 7z x -y -o"$ROZPAKOWANE/tess" "$POBRANE/tesseract-setup.exe" >/dev/null
[ -s "$POBRANE/7z-inst.exe" ]        && 7z x -y -o"$ROZPAKOWANE/7zip" "$POBRANE/7z-inst.exe" >/dev/null
[ -s "$POBRANE/espeak.msi" ]         && 7z x -y -o"$ROZPAKOWANE/espeak" "$POBRANE/espeak.msi" >/dev/null
for para in ffmpeg pandoc poppler piper esrgan; do
  [ -s "$POBRANE/$para.zip" ] && unzip -q -o "$POBRANE/$para.zip" -d "$ROZPAKOWANE/$para"
done

zglos "Ułożenie arsenału obok rdzenia"
rm -rf "$POMOCNIKI"; mkdir -p "$POMOCNIKI/tessdata" "$POMOCNIKI/realesrgan-ncnn-vulkan/models" "$POMOCNIKI/glosy"
# ImageMagick w wydaniu przenośnym niesie samodzielny `magick.exe` wraz z plikami
# nastaw XML — bez nich program wstaje, ale nie zna delegatów i profili.
[ -d "$ROZPAKOWANE/im" ] && cp -f "$ROZPAKOWANE/im/magick.exe" "$ROZPAKOWANE"/im/*.xml "$ROZPAKOWANE/im/sRGB.icc" "$POMOCNIKI/"
# ffmpeg w wydaniu statycznym to dwa pliki bez żadnej biblioteki obok.
[ -d "$ROZPAKOWANE/ffmpeg" ] && cp -f "$ROZPAKOWANE"/ffmpeg/*/bin/ffmpeg.exe "$ROZPAKOWANE"/ffmpeg/*/bin/ffprobe.exe "$POMOCNIKI/"
[ -d "$ROZPAKOWANE/pandoc" ] && cp -f "$ROZPAKOWANE"/pandoc/*/pandoc.exe "$POMOCNIKI/"
[ -d "$ROZPAKOWANE/7zip" ]   && cp -f "$ROZPAKOWANE/7zip/7z.exe" "$ROZPAKOWANE/7zip/7z.dll" "$POMOCNIKI/"
# poppler idzie przed Tesseractem, bo przy dwóch zbieżnych bibliotekach
# (`libcrypto-3-x64.dll`, `libzstd.dll`) pierwszeństwo ma wydanie Tesseractu —
# to ono jest nowsze i to jego zestaw sprawdzono razem z danymi języków.
if [ -d "$ROZPAKOWANE/poppler" ]; then
  cp -f "$ROZPAKOWANE"/poppler/*/Library/bin/pdftotext.exe \
        "$ROZPAKOWANE"/poppler/*/Library/bin/pdftoppm.exe \
        "$ROZPAKOWANE"/poppler/*/Library/bin/pdfinfo.exe \
        "$ROZPAKOWANE"/poppler/*/Library/bin/pdfimages.exe "$POMOCNIKI/"
  cp -f "$ROZPAKOWANE"/poppler/*/Library/bin/*.dll "$POMOCNIKI/"
fi
if [ -d "$ROZPAKOWANE/tess" ]; then
  cp -f "$ROZPAKOWANE/tess/tesseract.exe" "$POMOCNIKI/"
  cp -f "$ROZPAKOWANE"/tess/*.dll "$POMOCNIKI/"
  for jezyk in pol eng osd; do
    [ -s "$POBRANE/$jezyk.traineddata" ] && cp -f "$POBRANE/$jezyk.traineddata" "$POMOCNIKI/tessdata/"
  done
fi
if [ -d "$ROZPAKOWANE/piper" ]; then
  cp -f "$ROZPAKOWANE"/piper/piper/piper.exe "$ROZPAKOWANE"/piper/piper/*.dll "$ROZPAKOWANE"/piper/piper/*.ort "$POMOCNIKI/"
  cp -r "$ROZPAKOWANE/piper/piper/espeak-ng-data" "$POMOCNIKI/"
  cp -f "$POBRANE"/pl_PL-gosia-medium.onnx* "$POMOCNIKI/glosy/" 2>/dev/null || true
fi
# eSpeak NG jest zapasową syntezą mowy, gdy piper nie wstanie. Instalator MSI
# trzyma pliki pod nazwami z podkreśleniem, a rdzeń woła `espeak-ng`, zaś sam
# program ładuje `libespeak-ng.dll` — obie nazwy wracają tu do postaci z myślnikiem.
# Danych głosowych nie wyjmujemy z MSI: `espeak-ng-data` przyszło już z piperem
# i jest tym samym zbiorem, więc druga kopia byłaby tylko cięższa.
if [ -d "$ROZPAKOWANE/espeak" ]; then
  cp -f "$ROZPAKOWANE/espeak/espeak_ng.exe" "$POMOCNIKI/espeak-ng.exe"
  cp -f "$ROZPAKOWANE/espeak/libespeak_ng.dll" "$POMOCNIKI/libespeak-ng.dll"
fi
# Real-ESRGAN: plik wykonywalny i jego biblioteka leżą PŁASKO (rdzeń woła
# `realesrgan-ncnn-vulkan` i znajduje go płaską drogą), ale WAGI jadą do
# `realesrgan-ncnn-vulkan/models` — dokładnie tam, gdzie po poprawce rdzenia
# (`adapter_narzedzia_obraz_model.go`, `katalogModeliPowiekszenia` na Windowsie)
# szuka ich argument `-m`. Rdzeń sprawdza tę ścieżkę PRZED wołaniem silnika,
# więc wagi w innym miejscu znaczą odmowę mimo obecnej paczki.
if [ -d "$ROZPAKOWANE/esrgan" ]; then
  cp -f "$ROZPAKOWANE/esrgan/realesrgan-ncnn-vulkan.exe" "$ROZPAKOWANE/esrgan/vcomp140.dll" "$POMOCNIKI/"
  cp -r "$ROZPAKOWANE"/esrgan/models/* "$POMOCNIKI/realesrgan-ncnn-vulkan/models/"
fi
# ── Programy niosące własne środowisko: własny katalog, nie płasko ────────────
# Układ `pomocniki/<program>/<program>.exe` obsługuje `zewnetrzne.Odnajdz`
# (patrz `miejscaPakietu` w `server/internal/zewnetrzne/odnajdywanie.go`).
# PowerShell jest tu pozycją PIERWSZĄ, bo warstwa skanera przez WIA woła `pwsh`,
# a WIA jest na Windowsie jedyną drogą skanera — SANE tam nie istnieje.
if [ -s "$POBRANE/pwsh.zip" ]; then
  mkdir -p "$POMOCNIKI/pwsh"
  unzip -q -o "$POBRANE/pwsh.zip" -d "$POMOCNIKI/pwsh"
  printf '  własny katalog: pwsh (%s)\n' "$(du -sh "$POMOCNIKI/pwsh" | cut -f1)"
fi

# LibreOffice wyjmujemy `msiextract`, a NIE `7z x`: `7z` czyta z pakietu MSI same
# strumienie plików i wysypuje 19 494 pliki do jednego katalogu bez struktury,
# a LibreOffice bez drzewa `program/`, `share/` i `presets/` nie wstanie.
# `msiextract` (pakiet msitools) odtwarza katalogi z tabel pakietu.
if [ -s "$POBRANE/libreoffice.msi" ]; then
  command -v msiextract >/dev/null || padnij "brak msiextract (apt: msitools) — bez niego LibreOffice wyjdzie bez struktury katalogów"
  mkdir -p "$ROZPAKOWANE/lo" "$POMOCNIKI/libreoffice"
  msiextract -C "$ROZPAKOWANE/lo" "$POBRANE/libreoffice.msi" >/dev/null 2>&1
  KORZEN_LO="$(dirname "$(find "$ROZPAKOWANE/lo" -maxdepth 3 -name soffice.exe -print -quit)")"
  cp -r "$(dirname "$KORZEN_LO")"/* "$POMOCNIKI/libreoffice/"
  # Zamiana formatów woła LibreOffice BEZ OKNA (`--convert-to`), więc motywy
  # ikon, tłumaczenia interfejsu, pomoc i skryptowy Python nie biorą w niej
  # udziału. Bez tego odchudzenia pakiet nie zmieściłby się pod granicą NSIS.
  find "$POMOCNIKI/libreoffice/program/resource" -mindepth 1 -maxdepth 1 -type d \
    ! -name 'en-US' ! -name 'pl' -exec rm -rf {} + 2>/dev/null || true
  rm -f "$POMOCNIKI"/libreoffice/share/config/images_*.zip
  rm -rf "$POMOCNIKI/libreoffice/program/python-core-"* \
         "$POMOCNIKI/libreoffice/share/extensions" \
         "$POMOCNIKI/libreoffice/share/gallery" \
         "$POMOCNIKI/libreoffice/help" "$POMOCNIKI/libreoffice/readmes"
  printf '  własny katalog: libreoffice (%s)\n' "$(du -sh "$POMOCNIKI/libreoffice" | cut -f1)"
fi

# Opakowania startowe składamy mingw-em, bo rdzeń woła `libreoffice`, a wydanie
# Windows nazywa ten plik `soffice.exe`. Powód, dla którego to opakowanie jest
# plikiem PE, a nie skryptem `.cmd`, stoi w nagłówku `uruchamiacz/uruchamiacz.c`.
if [ -d "$POMOCNIKI/libreoffice" ]; then
  x86_64-w64-mingw32-gcc -municode -O2 -s -o "$POMOCNIKI/libreoffice/libreoffice.exe" \
    "$ARSENAL/uruchamiacz/uruchamiacz.c" -DCEL_PROGRAMU='L"program\\soffice.exe"' \
    || padnij "nie złożyło się opakowanie libreoffice.exe"
  printf '  opakowanie: libreoffice.exe -> program\\soffice.exe\n'
fi

# Chromium: wydanie przenośne (ungoogled) rozpakowuje się do własnego katalogu,
# a rdzeń woła `chromium`, gdy tymczasem plik wydania nazywa się `chrome.exe` —
# stąd opakowanie o nazwie, której szuka rdzeń, wskazujące plik wydania obok.
if [ -s "$POBRANE/chromium.zip" ]; then
  mkdir -p "$ROZPAKOWANE/chromium" "$POMOCNIKI/chromium"
  unzip -q -o "$POBRANE/chromium.zip" -d "$ROZPAKOWANE/chromium"
  KORZEN_CHR="$(dirname "$(find "$ROZPAKOWANE/chromium" -maxdepth 3 -name chrome.exe -print -quit)")"
  cp -r "$KORZEN_CHR"/* "$POMOCNIKI/chromium/"
  x86_64-w64-mingw32-gcc -municode -O2 -s -o "$POMOCNIKI/chromium/chromium.exe" \
    "$ARSENAL/uruchamiacz/uruchamiacz.c" -DCEL_PROGRAMU='L"chrome.exe"' \
    || padnij "nie złożyło się opakowanie chromium.exe"
  printf '  własny katalog: chromium (%s), opakowanie chromium.exe -> chrome.exe\n' "$(du -sh "$POMOCNIKI/chromium" | cut -f1)"
fi

# rembg: silnik wycinania tła jest modułem Pythona, więc jego środowisko składa
# się z osadzalnego Pythona dla Windows, kół `win_amd64` i wag U²-Net. `pip` przy
# `rembg[cli]` NIE zaciąga `onnxruntime` (runtime wnioskowania) — dociągamy go
# osobno, bo bez niego rembg nie ruszy wcale.
if [ -s "$POBRANE/python-embed.zip" ] && [ -s "$POBRANE/u2net.onnx" ]; then
  mkdir -p "$POMOCNIKI/rembg/python" "$POMOCNIKI/rembg/modele"
  unzip -q -o "$POBRANE/python-embed.zip" -d "$POMOCNIKI/rembg/python"
  KOLA="$ROZPAKOWANE/rembg-kola"; mkdir -p "$KOLA"
  # Koła bierzemy pod cel Windows/cp311, nie pod host — instalka jedzie na Windows.
  python3 -m pip download --quiet --only-binary=:all: --platform win_amd64 \
    --python-version 3.11 --implementation cp --dest "$KOLA" "rembg[cli]" onnxruntime \
    || padnij "nie pobrały się koła rembg dla Windows"
  SITE="$POMOCNIKI/rembg/python/Lib/site-packages"
  python3 -m pip install --quiet --no-deps --no-index --target "$SITE" \
    --platform win_amd64 --python-version 3.11 --implementation cp --only-binary=:all: \
    "$KOLA"/*.whl || padnij "nie zainstalowały się koła rembg"
  # Python osadzony domyślnie nie widzi site-packages — bez tego wiersza żaden
  # zainstalowany pakiet nie jest importowalny.
  printf 'python311.zip\r\n.\r\nLib\\site-packages\r\n\r\n#import site\r\n' > "$POMOCNIKI/rembg/python/python311._pth"
  # Web UI (`rembg s`, gradio+pandas) nie jest wołane przez rdzeń, a waży ~285 MB
  # — usuwamy je razem z importem, który by się o nie potknął.
  find "$SITE" -name '__pycache__' -type d -prune -exec rm -rf {} + 2>/dev/null || true
  rm -rf "$SITE"/gradio "$SITE"/gradio_client "$SITE"/pandas "$SITE"/gradio-*.dist-info "$SITE"/pandas-*.dist-info
  rm -f "$SITE/rembg/commands/s_command.py"
  cp -f "$ARSENAL/sfx/rembg-commands-init.py" "$SITE/rembg/commands/__init__.py"
  cp -f "$ARSENAL/sfx/rembg-start.py" "$POMOCNIKI/rembg/start.py"
  cp -f "$POBRANE/u2net.onnx" "$POMOCNIKI/rembg/modele/"
  # Opakowanie ustawia `U2NET_HOME` na `pomocniki/rembg/modele` — ten sam katalog,
  # który zwraca `katalogWagWycinania` rdzenia na Windowsie.
  x86_64-w64-mingw32-gcc -municode -O2 -s -o "$POMOCNIKI/rembg/rembg.exe" \
    "$ARSENAL/uruchamiacz/uruchamiacz.c" \
    -DCEL_PROGRAMU='L"python\\python.exe"' -DCEL_ARGUMENTY='L"\"start.py\""' \
    -DZMIENNA_NAZWA='L"U2NET_HOME"' -DZMIENNA_PODKATALOG='L"modele"' \
    || padnij "nie złożyło się opakowanie rembg.exe"
  printf '  własny katalog: rembg (%s), opakowanie rembg.exe ustawia U2NET_HOME=modele\n' "$(du -sh "$POMOCNIKI/rembg" | cut -f1)"
fi

printf '  arsenał: %s razem, %s programów płaskich + pwsh/libreoffice/chromium/rembg\n' "$(du -sh "$POMOCNIKI" | cut -f1)" "$(find "$POMOCNIKI" -maxdepth 1 -name '*.exe' | wc -l)"

# ── Dlaczego 7-Zip SFX, a nie NSIS ───────────────────────────────────────────
# NSIS adresuje dane pakietu 32 bitami i przewraca się na ~2 GB zawartości
# NIESPAKOWANEJ (zdaniem „error mmapping file (2093971843, …) is out of range”).
# Pełny arsenał natywny waży ~2,9 GB, więc w NSIS się NIE mieści. 7-Zip SFX nie
# ma tej granicy — dlatego wersja z pełnym arsenałem jest samorozpakowującym się
# archiwum 7-Zip, a nie instalatorem NSIS. Składanie stoi niżej, po powłoce.
find "$POMOCNIKI" -maxdepth 1 -name '*.exe' -printf '    %f\n' | sort

zglos "Budowa powłoki dla $CEL"
# Powłoka WKOMPILOWUJE `client/dist`, więc przebudowany interfejs musi wymusić
# przebudowę powłoki. Cargo patrzy na pliki źródłowe skrzyni, a nie na katalog
# `frontendDist`, i przy niezmienionym kodzie Rusta zostawiłby binarkę
# z interfejsem z poprzedniej budowy — okno pokazywałoby wczorajszy interfejs
# przy dzisiejszym rdzeniu. Dotknięcie skryptu budowy jest tanie i pewne.
touch "$POWLOKA/build.rs"
# `--features tauri/custom-protocol` NIE jest ozdobą: `tauri::is_dev()` to
# dokładnie `!cfg!(feature = "custom-protocol")`, a `cargo build --release` sam
# tej cechy nie włącza. Bez niej `zrodlo_interfejsu::ustal` wchodzi w gałąź
# deweloperską, `WebviewUrl::default()` rozwiązuje się do `devUrl`, czyli
# `localhost:5173`, i okno u Operatora szuka serwera rozwojowego, którego na
# jego maszynie nie ma i nigdy nie będzie. Cecha idzie z wiersza poleceń, bo
# `Cargo.toml` opisuje powłokę, a nie to, czy właśnie składamy wydanie.
( cd "$POWLOKA" && cargo build --release --target "$CEL" --features tauri/custom-protocol )

zglos "Złożenie drzewa instalacji"
# Układ pod plikiem wykonywalnym powłoki musi być dokładnie taki, jakiego szuka
# `src/rdzen/lokalizacja.rs`: rdzeń w `rdzen/danaco-console.exe`, a obok niego
# serwer narzędzi, interfejs i arsenał. Powłoka + WebView2Loader.dll na wierzchu.
DRZEWO="$ARSENAL/_sfx/tree"
rm -rf "$ARSENAL/_sfx"; mkdir -p "$DRZEWO/rdzen"
# Arsenał (~2,9 GB) wnosimy TWARDYMI DOWIĄZANIAMI, żeby nie dublować go na dysku.
cp -al "$POMOCNIKI" "$DRZEWO/rdzen/pomocniki"
cp -f "$POWLOKA/zasoby/rdzen/danaco-console.exe" "$POWLOKA/zasoby/rdzen/danaco-narzedzia.exe" "$DRZEWO/rdzen/"
cp -al "$POWLOKA/zasoby/rdzen/dist" "$DRZEWO/rdzen/dist"
cp -f "$POWLOKA/target/$CEL/release/danaco-console-powloka.exe" "$DRZEWO/Danaco Console.exe"
cp -f "$POWLOKA/target/$CEL/release/WebView2Loader.dll" "$DRZEWO/"
mkdir -p "$DRZEWO/dokumentacja"
for dokument in README.md INSTALACJA-I-KONFIGURACJA.md INSTRUKCJA-UZYTKOWANIA.md LICENSE.md; do
  cp -f "$KORZEN/$dokument" "$DRZEWO/dokumentacja/"
done
cp -f "$ARSENAL/sfx/zainstaluj.cmd" "$DRZEWO/"

zglos "Zapora odbioru 1: powłoka niesie interfejs, nie adres serwera rozwojowego"
# Sprawdzamy obecność OSADZONEGO PAKIETU, nie brak napisu `localhost:5173` — ten
# napis siedzi w binarce ZAWSZE, bo Tauri wkompilowuje całą nastawę razem z polem
# `devUrl`. Rozstrzyga to, czego bez cechy `custom-protocol` w binarce NIE MA:
# obsługa `tauri://localhost` i nazwy plików z `client/dist`.
ZASOB="$(find "$KORZEN/budowa/client/dist/assets" -maxdepth 1 -type f -printf '%f\n' | head -1)"
strings -a "$DRZEWO/Danaco Console.exe" | grep -q 'tauri://localhost' \
  || padnij "powłoka nie zna protokołu tauri://localhost — zbudowano ją BEZ cechy tauri/custom-protocol"
strings -a "$DRZEWO/Danaco Console.exe" | grep -qF "$ZASOB" \
  || padnij "powłoka nie niesie osadzonego pakietu client/dist (brak $ZASOB)"
printf '  powłoka niesie osadzony pakiet interfejsu (znacznik: %s)\n' "$ZASOB"

zglos "Złożenie samorozpakowującego się instalatora (7-Zip SFX)"
# SFX to konkatenacja: moduł GUI 7-Zip + konfiguracja + archiwum. `7z l` czyta
# taki plik jak zwykłe archiwum, co wykorzystuje zapora 2.
MODUL_SFX="$ROZPAKOWANE/7zip/7z.sfx"
[ -s "$MODUL_SFX" ] || padnij "brak modułu 7z.sfx — rozpakuj instalator 7-Zip (7z x 7z-inst.exe do $ROZPAKOWANE/7zip)"
printf ';!@Install@!UTF-8!\r\nTitle="Danaco Console %s"\r\nBeginPrompt="Zainstalowac Danaco Console %s (wersja natywna) do wybranego katalogu?"\r\nProgress="yes"\r\nRunProgram="zainstaluj.cmd"\r\n;!@InstallEnd@!\r\n' "$WERSJA" "$WERSJA" > "$ARSENAL/_sfx/config.txt"
# `-mx=3 -mmt`: kompromis czasu i wagi; kompletność jest celem, nie stopień ściśnięcia.
( cd "$DRZEWO" && 7z a -mx=3 -mmt=on -bso0 -bsp0 "$ARSENAL/_sfx/payload.7z" . >/dev/null ) \
  || padnij "nie złożyło się archiwum 7-Zip"
mkdir -p "$WYDANIE"
WYNIK="$WYDANIE/Danaco Console_${WERSJA}_natywna_x64-setup.exe"
cat "$MODUL_SFX" "$ARSENAL/_sfx/config.txt" "$ARSENAL/_sfx/payload.7z" > "$WYNIK"
7z l "$WYNIK" >/dev/null 2>&1 || padnij "złożony SFX nie daje się odczytać przez 7z — konkatenacja jest wadliwa"

zglos "Zapora odbioru 2: pełny arsenał jest w instalatorze"
# „Wszystko ma działać" znaczy, że sprawdzamy, że wszystko JEST. Lista `7z l`
# wystarcza — nie trzeba rozpakowywać 2,9 GB, żeby policzyć obecność ścieżek.
LST="$(mktemp)"; 7z l -slt "$WYNIK" > "$LST" 2>/dev/null
brak=0
wymagaj() {
  grep -qiE "Path = ${1}" "$LST" || { printf '  BRAK: %s\n' "$2"; brak=1; }
}
wymagaj 'rdzen[\\/]pomocniki[\\/]rembg[\\/]modele[\\/]u2net\.onnx' 'rembg/modele/u2net.onnx'
wymagaj 'rdzen[\\/]pomocniki[\\/]realesrgan-ncnn-vulkan[\\/]models[\\/].*\.bin' 'realesrgan-ncnn-vulkan/models/*.bin'
wymagaj 'rdzen[\\/]pomocniki[\\/]chromium[\\/]chrome\.exe' 'chromium/'
wymagaj 'rdzen[\\/]pomocniki[\\/]libreoffice[\\/]program[\\/]soffice\.exe' 'libreoffice/'
wymagaj 'rdzen[\\/]pomocniki[\\/]pwsh[\\/]pwsh\.exe' 'pwsh/pwsh.exe'
wymagaj 'rdzen[\\/]danaco-console\.exe' 'rdzeń'
wymagaj 'rdzen[\\/]danaco-narzedzia\.exe' 'serwer narzędzi'
rm -f "$LST"
[ "$brak" -eq 0 ] || padnij "instalator nie niesie kompletnego arsenału — patrz braki wyżej"
printf '  komplet: rembg, Real-ESRGAN (wagi), Chromium, LibreOffice, pwsh, rdzeń, serwer narzędzi\n'

zglos "Wydanie"
( cd "$WYDANIE" && sha256sum "$(basename "$WYNIK")" > "$(basename "$WYNIK").sha256" )
printf 'postać  : 7-Zip SFX (samorozpakowujący PE32 GUI), rozszerzenie .exe\n'
printf 'ścieżka : %s\n' "$WYNIK"
printf 'rozmiar : %s\n' "$(du -h "$WYNIK" | cut -f1)"
printf 'suma    : %s\n' "$(sha256sum "$WYNIK" | cut -d' ' -f1)"
printf 'wewnątrz: %s\n' "$(7z l "$WYNIK" 2>/dev/null | tail -1 | sed 's/^[0-9 :-]*//')"

zglos "Czego ten skrypt NIE sprawdził"
cat <<'KONIEC'
  - czy instalator się uruchamia, wypakowuje i czy krok zainstaluj.cmd tworzy skróty,
  - czy zainstalowana aplikacja wstaje (wymaga WebView2 w systemie — SFX go nie niesie),
  - czy programy arsenału startują pod Windows — na Linuksie nie da się ich uruchomić,
  - czy opakowania (libreoffice.exe, chromium.exe, rembg.exe) wołają swój cel,
  - czy podpis Authenticode jest — NIE MA, plik jest niepodpisany.
  Wszystkie wymagają maszyny z Windows. Nie zakładaj ich powodzenia.
KONIEC
