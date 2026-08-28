#!/usr/bin/env bash
# Skrypt prowizjonuje na maszynie serwera komplet programów, na których stoi rdzeń,
# czytając wykaz zależności z binarium rdzenia i stawiając warstwy w kolejności.
set -euo pipefail

SKRYPTY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUDOWA="$(dirname "$SKRYPTY")"
# Minusy przed nazwą trybu obcinamy, żeby `sprawdz` i `--sprawdz` znaczyły to
# samo. Wołający pisze flagę odruchowo i nie ma powodu, żeby odbić go za to
# odmową.
TRYB="${1:-plan}"
TRYB="${TRYB#--}"
TRYB="${TRYB#-}"

# GOBIN_ARSENALU — katalog, do którego trafiają programy warsztatu Go. Domyślne
# ~/go/bin należy do roota prowizjonującego i nie jest na ścieżce procesu serwera,
# a program niewidoczny na ścieżce jest dla sondy rdzenia brakiem programu.
GOBIN_ARSENALU="${DANACO_GOBIN:-/usr/local/bin}"
KATALOG_MODELI="${DANACO_KATALOG_MODELI:-/opt/danaco-arsenal/modele-mowy}"
SRODOWISKO_MOWY="${DANACO_SRODOWISKO_MOWY:-/opt/danaco-arsenal/mowa}"

# Zaplecze pythonowe stoi poza wykazem zależności, ponieważ rdzeń woła oba
# środowiska interpreterem, a nie nazwą programu; ich brak widać dopiero odmową.
SRODOWISKO_WIEDZY="${DANACO_SRODOWISKO_WIEDZY:-/opt/danaco/silniki/wiedza}"
SRODOWISKO_TWARZY="${DANACO_SRODOWISKO_TWARZY:-/opt/danaco/silniki/twarze}"
KATALOG_WAG_TWARZY="${DANACO_KATALOG_WAG_TWARZY:-/opt/danaco-modele/twarze}"
OPAKOWANIE_TWARZY="${DANACO_OPAKOWANIE_TWARZY:-/usr/local/bin/danaco-twarze}"

# Nazwy plików wag przebiegu twarzowego przepisane ze stałych rdzenia, który
# sprawdza obecność wszystkich trzech przed startem pomocnika.
WAGI_ODTWARZANIA="GFPGANv1.4.pth"
WAGI_WYKRYWANIA="detection_Resnet50_Final.pth"
WAGI_PODZIALU="parsing_parsenet.pth"

# INDEKS_TORCH wskazuje składnicę kół PyTorcha liczących na procesorze, ponieważ
# wdrożenie jest procesorowe, a indeks domyślny ciąga wielogigabajtową warstwę CUDA.
INDEKS_TORCH="${DANACO_INDEKS_TORCH:-https://download.pytorch.org/whl/cpu}"

# PAKIETY_POZA_WYKAZEM zbiera pakiety dystrybucji nieobecne w wykazie zależności,
# bez których arsenał serwera pozostaje niekompletny; powód każdego stoi w dokumentacji.
PAKIETY_POZA_WYKAZEM="python3 python3-venv python3-pip nodejs npm sane-utils"

# WYKAZ_AWARYJNY niesie odpis kompletu deklaracji rdzenia w postaci wiersza wykazu,
# używany wyłącznie wtedy, gdy nie ma czym zapytać binarium rdzenia.
czytajWykazAwaryjny() {
	# Pola wiersza wykazu rozdziela tabulator, zapisywany poniżej literałem.
	cat <<-'WYKAZ'
		obowiazkowa-apt	7z	7zip	?	7-Zip	pakowanie i wydobycie zawartości archiwum
		obowiazkowa-recznie	java	środowisko uruchomieniowe Javy (default-jre) wraz z wydaniem Apache Tika w /opt/tika albo w katalogu wskazanym zmienną DANACO_TIKA	?	Apache Tika (uruchamiana środowiskiem Javy)	odczyt treści pliku w formacie spoza słownika rdzenia (document.text.extract) — arkusz, prezentacja, wiadomość poczty
		obowiazkowa-apt	chromium-browser	chromium-browser	?	Chromium	zrzuty stron, drzewo DOM, konsola, rejestr sieciowy, emulacja urządzenia i przewijanie w module Browser, a także strona otwierana przez audyt dostępności i audyt wydajności — te dwa dostają tę samą przeglądarkę, zamiast pobierać własną
		warsztat-go	dlv	go install github.com/go-delve/delve/cmd/dlv@latest	?	Delve	debugowanie krokowe programów Go w oknie Run & Debug
		decyzyjna	docker	docker.io albo podman	?	Docker	wykaz kontenerów i obrazów, budowanie obrazu i uruchomienie stosu w zakładce Containers
		decyzyjna	docker	docker.io	?	Docker (klient wiersza poleceń)	karta powłoki wewnątrz kontenera w module Terminal
		warsztat-npm	eslint	npm i -g eslint	?	ESLint	analiza statyczna kodu TypeScript i JavaScript
		obowiazkowa-apt	exiftool	libimage-exiftool-perl	?	ExifTool	odczyt metadanych IPTC, XMP i ID3 osadzonych w zasobie (library.metadata.get z includeTechnical) — EXIF i GPS czyta czytnik wkompilowany i te pola zostają także bez tego programu
		model-recznie	danaco-twarze	torch, torchvision i facexlib w osobnym środowisku pythonowym wraz z architekturą GFPGAN (gfpganv1_clean_arch, stylegan2_clean_arch z wydania github.com/TencentARC/GFPGAN), wystawione opakowaniem /usr/local/bin/danaco-twarze; wagi w /opt/danaco-modele/twarze	?	GFPGAN (pomocnik pythonowy)	osobny przebieg poprawiania twarzy przy powiększaniu obrazu (image.upscale z faces: true) — bez niego powiększanie pracuje dalej, a żądanie z tym polem odmawia zamiast oddać obraz bez poprawki twarzy
		obowiazkowa-apt	magick	imagemagick	?	ImageMagick	zapis obrazu w AVIF oraz w WEBP stratnym, a także pomiar pliku AVIF (image.convert, image.inspect) — pozostałe czynności rodziny image.* liczy biblioteka wkompilowana i przy braku tego programu pracują dalej
		obowiazkowa-recznie	java	środowisko uruchomieniowe Javy (default-jre) wraz z wydaniem LanguageToola w /opt/languagetool albo w katalogu wskazanym zmienną DANACO_LANGUAGETOOL	?	LanguageTool (uruchamiany środowiskiem Javy)	gramatyka, ortografia, interpunkcja, typografia i styl w korekcie językowej modułu Translate (translate.proofread.run)
		obowiazkowa-apt	libreoffice	libreoffice	?	LibreOffice	zamiana formatów biurowych, których nie czyta Pandoc
		warsztat-npm	lighthouse	npm i -g lighthouse	?	Lighthouse	audyt wydajności strony produktu wraz z Core Web Vitals (apps.performance.audit) — miary powstają w przeglądarce po wykonaniu skryptów, więc rdzeń nie policzy ich własnym pobraniem
		obowiazkowa-apt	node	nodejs	?	Node.js	orzeczenie o składni skryptu karty node w module Terminal (terminal.script.lint) — innego analizatora ta karta nie ma
		obowiazkowa-apt	ssh	openssh-client	?	OpenSSH	powłoka zdalna, odczyt pliku na maszynie karty i przekierowania portów w module Terminal
		obowiazkowa-apt	optipng	optipng	?	OptiPNG	dogniecenie zapisu PNG przy zamianie formatu (image.convert) — bez tego programu obraz zapisuje się tak samo, tylko dłuższym strumieniem
		warsztat-npm	pa11y	npm i -g pa11y	?	Pa11y	audyt dostępności bieżącej strony okna wobec normy WCAG (browser.accessibility.audit) — reguły normy są cudzą wiedzą i rdzeń ich nie przepisuje
		obowiazkowa-apt	pandoc	pandoc	?	Pandoc	zamiana formatów dokumentu w modułach Studio, Translate i Library
		snap	pwsh	powershell (snap) wraz z modułem PSScriptAnalyzer	?	PowerShell z modułem PSScriptAnalyzer	analiza statyczna i formatowanie skryptów PowerShell w module Terminal
		warsztat-npm	prettier	npm i -g prettier	?	Prettier	formatowanie plików TypeScript, JavaScript, CSS, Markdown i YAML
		obowiazkowa-apt	python3	python3	?	Python	orzeczenie o składni skryptu karty python w module Terminal (terminal.script.lint) na maszynie bez Ruffa — Ruff ma pierwszeństwo i obejmuje składnię wraz z regułami, interpreter zostaje drogą zapasową
		model-recznie	realesrgan-ncnn-vulkan	wydanie realesrgan-ncnn-vulkan z github.com/xinntao/Real-ESRGAN/releases rozpakowane do /usr/local/bin wraz z modelami w /usr/local/share/realesrgan-ncnn-vulkan/models, oraz mesa-vulkan-drivers dla liczenia na procesorze	?	Real-ESRGAN (ncnn)	powiększanie obrazu w module Design
		obowiazkowa-apt	ruff	pip install ruff	?	Ruff	analiza statyczna plików Pythona w module Developer oraz analiza i formatowanie treści skryptu karty python w module Terminal
		obowiazkowa-apt	semgrep	pip install semgrep	?	Semgrep	poszerzenie skanu kodu w module Developer o reguły semantyczne — wyłącznie w repozytorium niosącym własny zestaw reguł, bo zestaw z rejestru wymagałby sieci
		obowiazkowa-apt	shellcheck	shellcheck	?	ShellCheck	analiza statyczna skryptów powłoki bash w module Terminal
		warsztat-npm	stylelint	npm i -g stylelint	?	Stylelint	analiza statyczna arkuszy CSS w module Developer — wyłącznie w repozytorium niosącym własną konfigurację Stylelinta, bo program nie ma wbudowanego zestawu reguł
		obowiazkowa-apt	telnet	telnet	?	Telnet (klient)	karta sesji Telnet do urządzenia sieciowego w module Terminal
		obowiazkowa-apt	tesseract	tesseract-ocr tesseract-ocr-pol	?	Tesseract OCR	rozpoznanie pisma ze skanu i zdjęcia kartki
		warsztat-npm	ast-grep	npm i -g @ast-grep/cli	?	ast-grep	wyszukanie i zamiana po składni w module Developer — wzorzec z metazmienną (`$NAZWA`) idzie tą drogą zamiast po napisie
		warsztat-npm	autocannon	npm i -g autocannon	?	autocannon	przebieg obciążeniowy punktu końcowego wraz z percentylami czasu odpowiedzi i przepustowością (developer.api.load.run) — developer.api.request strzela jednym żądaniem i rozkładu nie ma z czego policzyć
		obowiazkowa-apt	cwebp	webp	?	cwebp	dogniecenie zapisu WEBP bezstratnego przy zamianie formatu (image.convert) — bez tego programu obraz zapisuje się koderem wkompilowanym; zapis stratny idzie inną drogą i nie jest dogniatany
		warsztat-go	dupl	go install github.com/mibk/dupl@latest	?	dupl	wykrywanie powtórzonych fragmentów w plikach Go podczas skanu kodu w module Developer
		obowiazkowa-apt	espeak-ng	espeak-ng	?	eSpeak NG	odsłuch przebiegu debaty syntezą mowy w module Roundtable
		obowiazkowa-apt	ffmpeg	ffmpeg	?	ffmpeg	zamiana formatu nagrania
		obowiazkowa-apt	ffprobe	ffmpeg	?	ffprobe	rozpoznanie zawartości nagrania dźwiękowego i filmowego
		obowiazkowa-apt	gofmt	golang	?	gofmt	formatowanie plików Go w module Developer
		warsztat-go	goimports	go install golang.org/x/tools/cmd/goimports@latest	?	goimports	formatowanie plików Go wraz z porządkowaniem importów
		warsztat-go	golangci-lint	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest	?	golangci-lint	analiza statyczna repozytorium Go w module Developer
		warsztat-go	gopls	go install golang.org/x/tools/gopls@latest	?	gopls	przejście do definicji, wystąpienia symbolu i refaktoryzacje semantyczne w module Developer
		obowiazkowa-recznie	hunspell	hunspell wraz ze słownikiem języka (hunspell-pl, hunspell-en-us)	?	hunspell	ortografia w korekcie językowej modułu Translate na maszynie bez LanguageToola — LanguageTool ma pierwszeństwo i obejmuje pisownię wraz z gramatyką, słownik zostaje drogą zapasową
		obowiazkowa-apt	jpegoptim	jpegoptim	?	jpegoptim	dogniecenie zapisu JPEG przy zamianie formatu (image.convert) — bez tego programu obraz zapisuje się tak samo, tylko dłuższym strumieniem
		warsztat-npm	jscpd	npm i -g jscpd	?	jscpd	wykrywanie powtórzonych fragmentów w plikach TypeScriptu i JavaScriptu podczas skanu kodu w module Developer
		snap	kubectl	kubectl (snap)	?	kubectl	karta powłoki wewnątrz poda w module Terminal
		obowiazkowa-apt	picocom	picocom	?	picocom	karta konsoli portu szeregowego w module Terminal
		obowiazkowa-apt	pngquant	pngquant	?	pngquant	sprowadzenie PNG do palety przy zamianie formatu na zapis stratny (image.convert z lossless: false) — przy zapisie bezstratnym nie jest wołany
		obowiazkowa-apt	pdftoppm	poppler-utils	?	poppler (pdftoppm)	rasteryzacja stron PDF przed rozpoznaniem pisma
		obowiazkowa-apt	pdftotext	poppler-utils	?	poppler (pdftotext)	odczyt warstwy tekstowej dokumentu PDF
		model-recznie	rembg	rembg[cli] w osobnym środowisku pythonowym, wystawiony opakowaniem /usr/local/bin/rembg ustawiającym U2NET_HOME=/usr/local/share/rembg-modele	?	rembg (U²-Net / ONNX Runtime)	wycinanie tła obrazu oraz rozkład obrazu na warstwy (image.background.remove, image.layers.split) — obie czynności stoją na tej samej sieci segmentującej i znikają razem z nią
		warsztat-npm	typescript-language-server	npm i -g typescript-language-server typescript	?	serwer języka TypeScript	przejście do definicji, wystąpienia symbolu i refaktoryzacje semantyczne plików TypeScriptu w module Developer
		obowiazkowa-apt	shfmt	shfmt	?	shfmt	formatowanie skryptów powłoki bash w module Terminal
		warsztat-go	staticcheck	go install honnef.co/go/tools/cmd/staticcheck@latest	?	staticcheck	pogłębiona analiza statyczna kodu Go
		obowiazkowa-recznie	typos	cargo install typos-cli	?	typos	wykrywanie literówek w identyfikatorach i treści plików repozytorium w module Developer
		obowiazkowa-recznie	typst	typst (jeden plik wykonywalny z wydania projektu)	?	typst	skład dokumentu do PDF-u w komendzie document.convert dla materiału, którego LibreOffice nie otwiera wprost (markdown, epub) — bez niego ta droga wraca do wersji zapasowej przez HTML i LibreOffice
		obowiazkowa-apt	unpaper	unpaper	?	unpaper	prostowanie skosu, odszumianie, progowanie i przycinanie marginesów skanu przed rozpoznaniem pisma (studio.ingest.recognize)
		obowiazkowa-recznie	vale	vale (jeden plik wykonywalny z wydania projektu)	?	vale	styl prozy w korekcie językowej modułu Translate — powtórzenia i terminy, zestawem reguł wbudowanym w program
	WYKAZ
}

zglos() { printf '\n=== %s ===\n' "$1"; }
padnij() {
	printf 'ODMOWA: %s\n' "$1" >&2
	exit 1
}

# odnajdzRdzen wskazuje binarium rdzenia, które wypisze wykaz. Kolejność:
# wskazanie wprost, binarium obok skryptu albo w wydaniu, ścieżka systemu.
odnajdzRdzen() {
	if [ -n "${DANACO_RDZEN:-}" ]; then
		[ -x "$DANACO_RDZEN" ] || padnij "DANACO_RDZEN wskazuje plik nieuruchamialny: $DANACO_RDZEN"
		printf '%s\n' "$DANACO_RDZEN"
		return
	fi
	for kandydat in \
		"$SKRYPTY/danaco-console" \
		"$BUDOWA/danaco-console" \
		"$BUDOWA/wydania/1.0.0/danaco-console"; do
		if [ -x "$kandydat" ]; then
			printf '%s\n' "$kandydat"
			return
		fi
	done
	if sciezka="$(command -v danaco-console 2>/dev/null)"; then
		printf '%s\n' "$sciezka"
		return
	fi
	return 1
}

# pobierzWykaz oddaje wiersze danych wykazu zależności bez wierszy komentarza,
# biorąc je z pliku wskazanego zmienną albo z binarium rdzenia.
pobierzWykaz() {
	local surowy rdzen
	if [ -n "${DANACO_WYKAZ_PLIK:-}" ]; then
		[ -r "$DANACO_WYKAZ_PLIK" ] || padnij "DANACO_WYKAZ_PLIK nieczytelny: $DANACO_WYKAZ_PLIK"
		surowy="$(cat "$DANACO_WYKAZ_PLIK")"
	elif rdzen="$(odnajdzRdzen)" && surowy="$("$rdzen" --wykaz-zaleznosci 2>/dev/null)" &&
		[ -n "$surowy" ]; then
		: # wykaz policzony przez rdzeń — źródło właściwe
	else
		# Rdzenia nie ma czym zapytać, więc wykaz idzie z odpisu awaryjnego.
		printf 'UWAGA: nie odczytano wykazu z rdzenia (danaco-console --wykaz-zaleznosci).\n' >&2
		printf '       Jadę wykazem awaryjnym wpisanym w ten skrypt. Wykaz policzony przez\n' >&2
		printf '       rdzeń wskażesz zmienną DANACO_RDZEN albo DANACO_WYKAZ_PLIK.\n' >&2
		surowy="$(czytajWykazAwaryjny)"
	fi
	printf '%s\n' "$surowy" | awk -F'\t' 'NF>=3 && $1 !~ /^#/ && $1 != ""'
}

# pobierzWykazMowy oddaje wiersze arsenału mowy (klucz <TAB> wartość). Brak
# wykazu mowy nie zatrzymuje prowizjonowania — warstwa mowy zostaje wtedy
# pominięta z komunikatem, bo reszta arsenału jest od niej niezależna.
pobierzWykazMowy() {
	if [ -n "${DANACO_WYKAZ_MOWY_PLIK:-}" ]; then
		[ -r "$DANACO_WYKAZ_MOWY_PLIK" ] || padnij "DANACO_WYKAZ_MOWY_PLIK nieczytelny: $DANACO_WYKAZ_MOWY_PLIK"
		cat "$DANACO_WYKAZ_MOWY_PLIK" | awk -F'\t' 'NF>=2 && $1 !~ /^#/'
		return 0
	fi
	local rdzen
	rdzen="$(odnajdzRdzen)" || return 1
	"$rdzen" --wykaz-mowy 2>/dev/null | awk -F'\t' 'NF>=2 && $1 !~ /^#/' || return 1
}

# wartoscMowy wyjmuje z wykazu mowy jedną wartość wskazaną kluczem, oddając
# pusty łańcuch, gdy klucza w wykazie nie ma.
wartoscMowy() {
	printf '%s\n' "$WYKAZ_MOWY" | awk -F'\t' -v k="$1" '$1==k {print $2; exit}'
}

# Pole Pakiet warstwy obowiązkowej rozbiera się z rozpoznaniem postaci.

# nazwyPakietow oddaje 0, gdy całe pole składa się z nazw pakietów dystrybucji.
# Kształt nazwy: mała litera albo cyfra, dalej litery, cyfry, kropka, plus, minus
# (polityka nazw Debiana).
polePakietowe() {
	printf '%s\n' "$1" | grep -Eq '^[a-z0-9][a-z0-9.+-]*( [a-z0-9][a-z0-9.+-]*)*$'
}

# poleAptDoInstalacji oddaje te pola warstwy obowiązkowej, które nadają się do
# podania programowi apt, pomijając pola opisowe i polecenia innych narzędzi.
polaWarstwyObowiazkowej() {
	printf '%s\n' "$1" | awk -F'\t' '$1=="obowiazkowa-apt" {print $3}'
}

# pakietyApt zbiera tokeny warstwy obowiązkowej w kolejności pierwszego
# wystąpienia, bez powtórzeń — poppler-utils i ffmpeg padają w wykazie dwukrotnie
# (pdftotext i pdftoppm, ffmpeg i ffprobe), a apt ma je dostać raz.
pakietyApt() {
	local pole
	while IFS= read -r pole; do
		# Pytanie o pip install idzie przed pytaniem o kształt nazwy pakietu.
		case "$pole" in "pip install "*) continue ;; esac
		polePakietowe "$pole" && printf '%s\n' "$pole"
	done < <(polaWarstwyObowiazkowej "$1") |
		tr ' ' '\n' | awk 'NF && !widziane[$0]++'
}

# pakietyPip oddaje nazwy pakietów wyjęte z pól rozpoczynających się poleceniem
# pip install, w kolejności pierwszego wystąpienia i bez powtórzeń.
pakietyPip() {
	local pole
	while IFS= read -r pole; do
		case "$pole" in
		"pip install "*) printf '%s\n' "${pole#pip install }" ;;
		esac
	done < <(polaWarstwyObowiazkowej "$1") |
		tr ' ' '\n' | awk 'NF && !widziane[$0]++'
}

# podpowiedziOpisowe drukuje pozycje warstwy obowiązkowej, których pola skrypt
# nie rozpoznał jako polecenia — jako kroki ręczne, treścią pola.
podpowiedziOpisowe() {
	printf '%s\n' "$1" | awk -F'\t' '$1=="obowiazkowa-apt" {print $2"\t"$5"\t"$3}' |
		while IFS=$'\t' read -r program nazwa pole; do
			polePakietowe "$pole" && continue
			case "$pole" in "pip install "*) continue ;; esac
			printf '  %-30s [%s]\n      %s\n' "$nazwa" "$program" "$pole"
		done
}

# wypiszWarstwe drukuje pozycje jednej warstwy wykazu, podając dla każdej nazwę
# czytelną, nazwę programu oraz pakiet, z którego program pochodzi.
wypiszWarstwe() {
	printf '%s\n' "$1" | awk -F'\t' -v w="$2" \
		'$1==w {printf "  %-34s [%s]  ← %s\n", $5, $2, $3}'
}

# policzWarstwe oddaje liczbę pozycji należących do wskazanej warstwy wykazu,
# służąc wierszom zbiorczym planu oraz sprawdzenia.
policzWarstwe() {
	printf '%s\n' "$1" | awk -F'\t' -v w="$2" '$1==w' | wc -l | tr -d ' '
}

# obecny mówi, czy program jest osiągalny na ścieżce — tak samo, jak sprawdza to
# sonda rdzenia. Sprawdzenie, nie instalacja.
obecny() { command -v "$1" >/dev/null 2>&1; }

# Plan wypisuje zamierzone czynności prowizjonowania warstwa po warstwie, niczego
# nie instalując i nie wymagając uprawnień roota.
trybPlan() {
	local wykaz="$1"
	printf 'Wykaz niesie %s pozycji. Plan prowizjonowania serwera:\n' \
		"$(printf '%s\n' "$wykaz" | wc -l | tr -d ' ')"

	zglos "OBOWIĄZKOWA (apt) — bez tego moduły odmawiają [$(policzWarstwe "$wykaz" obowiazkowa-apt) pozycji]"
	wypiszWarstwe "$wykaz" obowiazkowa-apt
	printf '\n  apt-get install -y %s\n' "$(pakietyApt "$wykaz" | tr '\n' ' ')"
	printf '  oraz poza wykazem: %s\n' "$PAKIETY_POZA_WYKAZEM"
	printf '  (interpreter i venv mowy, nośnik warsztatu npm dla Prettiera i ESLinta,\n'
	printf '  warstwa SANE (program scanimage) cyfryzacji Studia — powody stoją przy\n'
	printf '  PAKIETY_POZA_WYKAZEM w tym skrypcie)\n'
	local pip
	pip="$(pakietyPip "$wykaz" | tr '\n' ' ')"
	[ -z "${pip// /}" ] || printf '\n  pip install %s   (tak brzmi podpowiedź rdzenia)\n' "$pip"
	printf '\n  KROKI RĘCZNE tej warstwy — pola, których skrypt nie wykona za Operatora:\n'
	podpowiedziOpisowe "$wykaz"

	zglos "WIEDZA — biblioteki pythonowe silnika wiedzy, przesiewu i osi obrazu"
	planWiedzy
	zglos "TWARZE — pomocnik odtwarzania twarzy (danaco-twarze)"
	planTwarzy

	zglos "WARSZTAT Go — moduł Developer pracuje na serwerze [$(policzWarstwe "$wykaz" warsztat-go)]"
	wypiszWarstwe "$wykaz" warsztat-go
	zglos "WARSZTAT npm [$(policzWarstwe "$wykaz" warsztat-npm)]"
	wypiszWarstwe "$wykaz" warsztat-npm
	zglos "SNAP [$(policzWarstwe "$wykaz" snap)]"
	wypiszWarstwe "$wykaz" snap
	zglos "OBOWIĄZKOWE — kroki ręczne, podpowiedź zdaniem albo polecenie spoza apt [$(policzWarstwe "$wykaz" obowiazkowa-recznie)]"
	wypiszWarstwe "$wykaz" obowiazkowa-recznie
	zglos "MODELE/SILNIKI — kroki ręczne, wydania spoza repozytoriów [$(policzWarstwe "$wykaz" model-recznie)]"
	wypiszWarstwe "$wykaz" model-recznie

	zglos "ARSENAŁ MOWY — poza wykazem zależności"
	planMowy

	zglos "DECYZYJNA — silnik kontenerów WSTRZYMANY decyzją Właściciela [$(policzWarstwe "$wykaz" decyzyjna)]"
	wypiszWarstwe "$wykaz" decyzyjna
	printf '  NIE stawiany. Włączenie wymaga wyraźnego żądania: DANACO_SILNIK_KONTENEROW=tak\n'

	printf '\nPlan wypisany. Nic nie zmieniono. Postawienie: bash %s postaw\n' "${BASH_SOURCE[0]}"
}

# planMowy wypisuje warstwę mowy: syntezę (piper wraz z głosami) i rozpoznanie
# (biblioteka pythonowa wraz z wagami modelu).
planMowy() {
	if [ -z "$WYKAZ_MOWY" ]; then
		printf '  pominięto: rdzeń nie wypisał arsenału mowy (--wykaz-mowy)\n'
		return
	fi
	local model wymagania pomocnik piper glosy
	model="$(wartoscMowy rozpoznanie.model-domyslny)"
	wymagania="$(wartoscMowy rozpoznanie.wymagania)"
	pomocnik="$(wartoscMowy rozpoznanie.pomocnik-skrypt)"
	piper="$(wartoscMowy synteza.piper-arsenal)"
	glosy="$(wartoscMowy synteza.piper-glosy)"

	printf '  synteza — piper (głos dobry): binarium do %s\n' "$piper"
	printf '            wskazanie zmienną %s; głosy .onnx do %s (zmienna %s)\n' \
		"$(wartoscMowy synteza.piper-zmienna)" "$glosy" "$(wartoscMowy synteza.piper-glosy-zmienna)"
	printf '            KROK RĘCZNY: piper-tts oraz plik głosu pl (pl_PL-*.onnx wraz z .onnx.json)\n'
	printf '  synteza — %s (głos zapasowy): stoi w warstwie obowiązkowej wykazu\n' \
		"$(wartoscMowy synteza.espeak-program)"
	printf '  rozpoznanie — środowisko pythonowe: %s\n' "$SRODOWISKO_MOWY"
	if [ -n "$wymagania" ]; then
		printf '            pip install -r %s   (nazwa i wersja biblioteki stoją TAM, nie tutaj)\n' "$wymagania"
	else
		printf '            UWAGA: rdzeń nie wskazał pliku wymagań (pomocnika nie widać obok binarium)\n'
	fi
	printf '            pomocnik: %s\n' "${pomocnik:-nie znaleziony obok binarium rdzenia}"
	printf '  rozpoznanie — WAGI MODELU pobierane z góry, rozmiar „%s"\n' "${model:-?}"
	printf '            katalog: %s (rdzeniowi wskazać ustawieniem %s)\n' \
		"$KATALOG_MODELI" "$(wartoscMowy rozpoznanie.ustawienie-katalogu-modeli)"
	printf '            pobranie „small" to około 480 MB; rozmiar zmierzony po pobraniu wypisze tryb postaw\n'
	printf '            pomiń pobranie: DANACO_BEZ_WAG_MOWY=1 (pierwsze użycie mikrofonu czeka wtedy na sieć)\n'
}

# Zaplecze wiedzy stawia środowisko pythonowe silnika wiedzy, przesiewu
# wyszukiwania oraz osi obrazu, którego wykaz zależności rdzenia nie wypisuje.
planWiedzy() {
	printf '  środowisko pythonowe: %s\n' "$SRODOWISKO_WIEDZY"
	printf '    pip install --index-url %s torch torchvision\n' "$INDEKS_TORCH"
	printf '    pip install fastembed transformers pillow\n'
	printf '  rdzeniowi wskazać interpreter ustawieniem wiedza_program = %s/bin/python\n' \
		"$SRODOWISKO_WIEDZY"
	printf '  WAGI MODELI nie idą tędy — rdzeń szuka ich w katalogach ustawień\n'
	printf '    wiedza_katalog_modeli, wiedza_model_przesiewu i wiedza_model_obrazu\n'
	printf '    (domyślnie pod /opt/danaco-modele); pobranie jest krokiem osobnym\n'
}

postawWiedze() {
	command -v python3 >/dev/null 2>&1 || {
		printf '  UWAGA: brak python3 — zaplecze wiedzy pominięte\n'
		return
	}
	zbudujSrodowisko "$SRODOWISKO_WIEDZY" || return
	# Dwa wywołania, bo indeks procesorowy niesie tylko torch i torchvision.
	"$SRODOWISKO_WIEDZY/bin/pip" install --index-url "$INDEKS_TORCH" torch torchvision ||
		printf '  UWAGA: instalacja torch/torchvision nie powiodła się\n'
	"$SRODOWISKO_WIEDZY/bin/pip" install fastembed transformers pillow ||
		printf '  UWAGA: instalacja fastembed/transformers/pillow nie powiodła się\n'
	printf '  rdzeniowi wskazać interpreter ustawieniem wiedza_program = %s/bin/python\n' \
		"$SRODOWISKO_WIEDZY"
}

# Pomocnik odtwarzania twarzy dostaje własne środowisko pythonowe, opakowanie na
# ścieżce systemu oraz trzy zestawy wag wymagane przed startem przebiegu.
planTwarzy() {
	printf '  środowisko pythonowe: %s\n' "$SRODOWISKO_TWARZY"
	printf '    pip install --index-url %s torch torchvision\n' "$INDEKS_TORCH"
	printf '    pip install facexlib\n'
	printf '  opakowanie: %s (rdzeń woła je, nie plik z wnętrza środowiska)\n' "$OPAKOWANIE_TWARZY"
	printf '  wagi w %s:\n' "$KATALOG_WAG_TWARZY"
	printf '    %-30s pobiera facexlib własnym pobieraczem\n' "$WAGI_WYKRYWANIA"
	printf '    %-30s pobiera facexlib własnym pobieraczem\n' "$WAGI_PODZIALU"
	printf '    %-30s KROK RĘCZNY — wydanie github.com/TencentARC/GFPGAN\n' "$WAGI_ODTWARZANIA"
	printf '  KROK RĘCZNY: architektura GFPGAN (gfpganv1_clean_arch, stylegan2_clean_arch\n'
	printf '    z tego samego wydania) ma stanąć w site-packages środowiska jako pakiet\n'
	printf '    gfpgan_clean — pomocnik importuje tę nazwę. Adresu pobrania rdzeń nie\n'
	printf '    podaje, więc skrypt go nie zgaduje.\n'
}

postawTwarze() {
	command -v python3 >/dev/null 2>&1 || {
		printf '  UWAGA: brak python3 — pomocnik twarzy pominięty\n'
		return
	}
	zbudujSrodowisko "$SRODOWISKO_TWARZY" || return
	"$SRODOWISKO_TWARZY/bin/pip" install --index-url "$INDEKS_TORCH" torch torchvision ||
		printf '  UWAGA: instalacja torch/torchvision nie powiodła się\n'
	"$SRODOWISKO_TWARZY/bin/pip" install facexlib ||
		printf '  UWAGA: instalacja facexlib nie powiodła się\n'

	# Zapis idzie przez plik tymczasowy i przemianowanie, zawsze tą samą treścią.
	mkdir -p "$(dirname "$OPAKOWANIE_TWARZY")"
	cat >"$OPAKOWANIE_TWARZY.czesciowy" <<-OPAKOWANIE
		#!/bin/sh
		# Środowisko nazwane tutaj, bo rdzeń woła pomocnika bez dziedziczenia.
		export HOME=/tmp
		export XDG_CACHE_HOME=/tmp
		export OMP_NUM_THREADS=4
		exec $SRODOWISKO_TWARZY/bin/python "\$@"
	OPAKOWANIE
	chmod 0755 "$OPAKOWANIE_TWARZY.czesciowy"
	mv -f "$OPAKOWANIE_TWARZY.czesciowy" "$OPAKOWANIE_TWARZY"
	printf '  opakowanie: %s\n' "$OPAKOWANIE_TWARZY"

	pobierzWagiTwarzy
}

# pobierzWagiTwarzy ściąga dwa zestawy wag pobieraczem samej biblioteki, ponieważ
# adres wydania zapisany tutaj rozjechałby się z biblioteką przy jej nowym wydaniu.
pobierzWagiTwarzy() {
	mkdir -p "$KATALOG_WAG_TWARZY"
	if [ -f "$KATALOG_WAG_TWARZY/$WAGI_WYKRYWANIA" ] &&
		[ -f "$KATALOG_WAG_TWARZY/$WAGI_PODZIALU" ]; then
		printf '  wagi wykrywacza i podziału twarzy już leżą w %s — nie pobieram\n' \
			"$KATALOG_WAG_TWARZY"
	elif "$SRODOWISKO_TWARZY/bin/python" - "$KATALOG_WAG_TWARZY" <<-'PYTHON'; then
		import sys
		from facexlib.utils.face_restoration_helper import FaceRestoreHelper

		# Te same trzy parametry, którymi pomocnik składa przebieg twarzowy.
		FaceRestoreHelper(1, det_model="retinaface_resnet50", device="cpu",
		                  use_parse=True, model_rootpath=sys.argv[1])
	PYTHON
		printf '  wagi wykrywacza i podziału pobrane do %s\n' "$KATALOG_WAG_TWARZY"
	else
		printf '  UWAGA: pobranie wag facexlib nie powiodło się — przebieg twarzowy odmówi\n'
	fi
	if [ -f "$KATALOG_WAG_TWARZY/$WAGI_ODTWARZANIA" ]; then
		printf '  %s stoi\n' "$WAGI_ODTWARZANIA"
	else
		printf '  KROK RĘCZNY: %s (333 MB) z wydania github.com/TencentARC/GFPGAN do %s\n' \
			"$WAGI_ODTWARZANIA" "$KATALOG_WAG_TWARZY"
	fi
	printf '  KROK RĘCZNY: architektura gfpgan_clean w site-packages %s\n' "$SRODOWISKO_TWARZY"
}

# zbudujSrodowisko stawia środowisko pythonowe albo zostawia stojące nietknięte,
# nazywając wprost, który przebieg co zrobił.
zbudujSrodowisko() {
	local katalog="$1"
	if [ -x "$katalog/bin/pip" ]; then
		printf '  środowisko stoi: %s\n' "$katalog"
		return 0
	fi
	printf '  stawiam środowisko: %s\n' "$katalog"
	python3 -m venv "$katalog" 2>/dev/null || {
		printf '  UWAGA: venv nie powstał w %s (brak python3-venv?)\n' "$katalog"
		return 1
	}
	"$katalog/bin/pip" install --upgrade pip >/dev/null 2>&1 || true
	[ -x "$katalog/bin/pip" ] || {
		printf '  UWAGA: w %s nie ma pip\n' "$katalog"
		return 1
	}
	return 0
}

# Sprawdzenie nie instaluje niczego i nie wymaga roota, a wykaz programów bierze
# z tego samego źródła, z którego biorą go plan oraz postawienie.
trybSprawdz() {
	local wykaz="$1" brakow=0 wszystkich=0
	zglos "Obecność programów wykazu (command -v) — read-only, nic nie instaluję"
	while IFS=$'\t' read -r warstwa program pakiet _stoi nazwa _zakres; do
		wszystkich=$((wszystkich + 1))
		if obecny "$program"; then
			printf '  %-24s → jest    %s\n' "$program" "$nazwa"
		else
			brakow=$((brakow + 1))
			printf '  %-24s → BRAK    %s  ← %s  {%s}\n' \
				"$program" "$nazwa" "$pakiet" "$warstwa"
		fi
	done <<<"$wykaz"
	printf '\nzależności zewnętrzne: %d/%d obecnych na tej maszynie\n' \
		"$((wszystkich - brakow))" "$wszystkich"
	if [ "$brakow" -gt 0 ]; then
		printf 'BRAKUJE %d — postawienie: sudo bash %s postaw\n' \
			"$brakow" "${BASH_SOURCE[0]}"
	fi

	# Programy spoza wykazu rdzenia liczone osobno, lecz wchodzące do kodu wyjścia.
	zglos "Poza wykazem zależności — nośniki i warstwa SANE"
	local program
	for program in python3 node npm scanimage; do
		if obecny "$program"; then
			printf '  %-24s → jest\n' "$program"
		else
			brakow=$((brakow + 1))
			printf '  %-24s → BRAK    ← %s\n' "$program" "$PAKIETY_POZA_WYKAZEM"
		fi
	done

	# Miarą jest import biblioteki, nie command -v: rdzeń woła je interpreterem.
	zglos "Zaplecze wiedzy — biblioteki pythonowe"
	local biblioteka
	for biblioteka in fastembed torch transformers PIL; do
		if [ -x "$SRODOWISKO_WIEDZY/bin/python" ] &&
			"$SRODOWISKO_WIEDZY/bin/python" -c "import $biblioteka" >/dev/null 2>&1; then
			printf '  %-24s → jest    w %s\n' "$biblioteka" "$SRODOWISKO_WIEDZY"
		else
			brakow=$((brakow + 1))
			printf '  %-24s → BRAK    w %s\n' "$biblioteka" "$SRODOWISKO_WIEDZY"
		fi
	done

	zglos "Pomocnik odtwarzania twarzy"
	if obecny danaco-twarze; then
		printf '  %-24s → jest\n' "danaco-twarze"
	else
		brakow=$((brakow + 1))
		printf '  %-24s → BRAK    ← opakowanie %s\n' "danaco-twarze" "$OPAKOWANIE_TWARZY"
	fi
	for biblioteka in torch torchvision facexlib gfpgan_clean; do
		if [ -x "$SRODOWISKO_TWARZY/bin/python" ] &&
			"$SRODOWISKO_TWARZY/bin/python" -c "import $biblioteka" >/dev/null 2>&1; then
			printf '  %-24s → jest    w %s\n' "$biblioteka" "$SRODOWISKO_TWARZY"
		else
			brakow=$((brakow + 1))
			printf '  %-24s → BRAK    w %s\n' "$biblioteka" "$SRODOWISKO_TWARZY"
		fi
	done
	local waga
	for waga in "$WAGI_ODTWARZANIA" "$WAGI_WYKRYWANIA" "$WAGI_PODZIALU"; do
		if [ -f "$KATALOG_WAG_TWARZY/$waga" ]; then
			printf '  %-30s → jest    %s\n' "$waga" \
				"$(du -h "$KATALOG_WAG_TWARZY/$waga" 2>/dev/null | cut -f1)"
		else
			brakow=$((brakow + 1))
			printf '  %-30s → BRAK    w %s\n' "$waga" "$KATALOG_WAG_TWARZY"
		fi
	done

	zglos "Arsenał mowy"
	if [ -z "$WYKAZ_MOWY" ]; then
		printf '  pominięto: rdzeń nie wypisał arsenału mowy\n'
	else
		local piper glosy model
		piper="$(wartoscMowy synteza.piper-program)"
		glosy="$(wartoscMowy synteza.piper-glosy)"
		model="$(wartoscMowy rozpoznanie.model-domyslny)"
		if obecny "$piper" || [ -x "$(wartoscMowy synteza.piper-arsenal)" ]; then
			printf '  STOI   piper\n'
		else
			printf '  BRAK   piper (głos dobry; zejście na espeak-ng działa)\n'
		fi
		if compgen -G "$glosy/*.onnx" >/dev/null 2>&1; then
			printf '  STOI   głosy pipera w %s\n' "$glosy"
		else
			printf '  BRAK   głosów .onnx w %s\n' "$glosy"
		fi
		if [ -d "$KATALOG_MODELI" ] && [ -n "$(ls -A "$KATALOG_MODELI" 2>/dev/null)" ]; then
			printf '  STOI   wagi modelu „%s" w %s (%s)\n' "$model" "$KATALOG_MODELI" \
				"$(du -sh "$KATALOG_MODELI" 2>/dev/null | cut -f1)"
		else
			printf '  BRAK   wag modelu „%s" w %s — pierwsze nagranie czekałoby na sieć\n' \
				"$model" "$KATALOG_MODELI"
		fi
	fi
	printf '\nPełną diagnozę wraz z zakresem, który nie zadziała, daje sonda startowa rdzenia.\n'
	[ "$brakow" -eq 0 ] || return 1
}

# Postawienie instaluje warstwy w kolejności wymuszonej ich zależnościami
# i wymaga uprawnień roota dla warstw apt oraz snap.
trybPostaw() {
	local wykaz="$1"
	[ "$(id -u)" -eq 0 ] || padnij "postaw wymaga roota (apt/snap). Uruchom przez sudo."

	zglos "Warstwa obowiązkowa (apt)"
	local apt
	apt="$(pakietyApt "$wykaz" | tr '\n' ' ')"
	[ -n "${apt// /}" ] || padnij "warstwa obowiązkowa wyszła pusta — wykaz nie został poprawnie odczytany"
	command -v apt-get >/dev/null 2>&1 ||
		padnij "brak apt-get; warstwy obowiązkowej nie da się postawić automatycznie na tej dystrybucji. Pakiety: $apt"
	apt-get update
	# Do pakietów wykazu dochodzą pakiety spoza wykazu, z jednego miejsca.
	# shellcheck disable=SC2086
	apt-get install -y $apt $PAKIETY_POZA_WYKAZEM

	local pip
	pip="$(pakietyPip "$wykaz" | tr '\n' ' ')"
	if [ -n "${pip// /}" ]; then
		printf '  pip install %s\n' "$pip"
		# Znacznik konieczny, bo dystrybucja zarządza swoim Pythonem zewnętrznie.
		# shellcheck disable=SC2086
		python3 -m pip install --break-system-packages $pip ||
			printf '  UWAGA: instalacja pipem nie powiodła się: %s\n' "$pip"
	fi

	printf '\n  KROKI RĘCZNE warstwy obowiązkowej — skrypt ich nie wykonuje:\n'
	podpowiedziOpisowe "$wykaz"

	zglos "Warstwa warsztatu Go"
	if command -v go >/dev/null 2>&1; then
		# GOBIN idzie do /usr/local/bin, bo ~/go/bin nie stoi na ścieżce serwera.
		printf '%s\n' "$wykaz" | awk -F'\t' '$1=="warsztat-go" {print $3}' |
			while read -r polecenie; do
				# Polecenie nie zaczynające się od go install idzie jako krok ręczny.
				case "$polecenie" in
				"go install "*) ;;
				*)
					printf '  KROK RĘCZNY (nie jest poleceniem Go): %s\n' "$polecenie"
					continue
					;;
				esac
				printf '  GOBIN=%s %s\n' "$GOBIN_ARSENALU" "$polecenie"
				(
					export GOBIN="$GOBIN_ARSENALU"
					eval "$polecenie"
				) || printf '  UWAGA: nie powiodło się: %s\n' "$polecenie"
			done
	else
		printf '  pominięto: brak go na ścieżce (pakiet golang stawia warstwa apt)\n'
	fi

	zglos "Warstwa warsztatu npm"
	if command -v npm >/dev/null 2>&1; then
		printf '%s\n' "$wykaz" | awk -F'\t' '$1=="warsztat-npm" {print $3}' |
			while read -r polecenie; do
				printf '  %s\n' "$polecenie"
				eval "$polecenie" || printf '  UWAGA: nie powiodło się: %s\n' "$polecenie"
			done
	else
		printf '  pominięto: brak npm na ścieżce\n'
	fi

	zglos "Warstwa snap"
	if command -v snap >/dev/null 2>&1; then
		printf '%s\n' "$wykaz" | awk -F'\t' '$1=="snap" {print $3"\t"$5}' |
			while IFS=$'\t' read -r pakiet nazwa; do
				local nazwaSnapa="${pakiet%% *}"
				printf '  snap install %s (dla %s)\n' "$nazwaSnapa" "$nazwa"
				snap install "$nazwaSnapa" ||
					printf '  UWAGA: snap install %s nie powiódł się\n' "$nazwaSnapa"
			done
		printf '  PAMIĘTAJ: człony poza snapem (np. moduł PSScriptAnalyzer) to osobny krok Install-Module\n'
	else
		printf '  pominięto: brak snap na ścieżce\n'
	fi

	zglos "Arsenał mowy"
	postawMowe

	zglos "Zaplecze wiedzy (fastembed, torch, transformers, pillow)"
	postawWiedze

	zglos "Pomocnik odtwarzania twarzy (danaco-twarze)"
	postawTwarze

	zglos "OBOWIĄZKOWE — kroki ręczne (nieautomatyzowane)"
	wypiszWarstwe "$wykaz" obowiazkowa-recznie

	zglos "MODELE/SILNIKI — kroki ręczne (nieautomatyzowane)"
	wypiszWarstwe "$wykaz" model-recznie

	zglos "Warstwa decyzyjna — silnik kontenerów"
	if [ "${DANACO_SILNIK_KONTENEROW:-}" = "tak" ]; then
		printf '  stawiam na wyraźne żądanie (DANACO_SILNIK_KONTENEROW=tak)\n'
		printf '%s\n' "$wykaz" | awk -F'\t' '$1=="decyzyjna" {print $3}' | head -1 |
			while read -r podpowiedz; do
				printf '  podpowiedź wykazu: %s → stawiam docker.io\n' "$podpowiedz"
			done
		apt-get install -y docker.io || printf '  UWAGA: instalacja docker.io nie powiodła się\n'
	else
		printf '  WSTRZYMANY decyzją Właściciela — NIE stawiam.\n'
		printf '  Włączenie wprost: DANACO_SILNIK_KONTENEROW=tak (zakładka Containers odmawia bez niego)\n'
	fi

	printf '\nArsenał postawiony. Weryfikacja: sonda startowa rdzenia albo bash %s sprawdz\n' "${BASH_SOURCE[0]}"
}

# postawMowe buduje środowisko pythonowe rozpoznawania z pliku wymagań pomocnika
# i pobiera wagi modelu z góry, żeby pierwsze użycie mikrofonu nie czekało na sieć.
postawMowe() {
	if [ -z "$WYKAZ_MOWY" ]; then
		printf '  pominięto: rdzeń nie wypisał arsenału mowy (--wykaz-mowy)\n'
		return
	fi
	local wymagania model
	wymagania="$(wartoscMowy rozpoznanie.wymagania)"
	model="$(wartoscMowy rozpoznanie.model-domyslny)"

	if [ -z "$wymagania" ] || [ ! -r "$wymagania" ]; then
		printf '  UWAGA: pliku wymagań pomocnika nie widać (%s) — rozpoznawanie mowy zostaje bez biblioteki\n' \
			"${wymagania:-nie wskazany}"
		printf '  Naprawa: uruchom skrypt na maszynie, gdzie obok binarium rdzenia stoi katalog pomocniki/\n'
		return
	fi

	command -v python3 >/dev/null 2>&1 || {
		printf '  UWAGA: brak python3 — rozpoznawanie mowy pominięte\n'
		return
	}
	printf '  środowisko pythonowe: %s\n' "$SRODOWISKO_MOWY"
	python3 -m venv "$SRODOWISKO_MOWY" 2>/dev/null || printf '  UWAGA: venv nie powstał (brak python3-venv?)\n'
	if [ -x "$SRODOWISKO_MOWY/bin/pip" ]; then
		printf '  pip install -r %s\n' "$wymagania"
		"$SRODOWISKO_MOWY/bin/pip" install --upgrade pip >/dev/null 2>&1 || true
		"$SRODOWISKO_MOWY/bin/pip" install -r "$wymagania" ||
			printf '  UWAGA: instalacja biblioteki rozpoznawania nie powiodła się\n'
	else
		printf '  UWAGA: w %s nie ma pip — pomijam instalację biblioteki\n' "$SRODOWISKO_MOWY"
		return
	fi
	printf '  rdzeniowi wskazać interpreter ustawieniem %s = %s/bin/python\n' \
		"$(wartoscMowy rozpoznanie.ustawienie-interpretera)" "$SRODOWISKO_MOWY"

	if [ "${DANACO_BEZ_WAG_MOWY:-}" = "1" ]; then
		printf '  wagi modelu POMINIĘTE na żądanie (DANACO_BEZ_WAG_MOWY=1) — pierwsze nagranie czeka na sieć\n'
		return
	fi
	pobierzWagiMowy "$model"
}

# pobierzWagiMowy ściąga wagi modelu z góry. Parametry ładowania (procesor, int8)
# są te same, którymi pracuje pomocnik (pomocniki/transkrypcja/silnik.py) — wagi
# mają zostać pobrane w tej postaci, w której zostaną potem użyte.
pobierzWagiMowy() {
	local model="$1"
	[ -n "$model" ] || {
		printf '  UWAGA: rdzeń nie podał rozmiaru modelu — pomijam pobranie wag\n'
		return
	}
	mkdir -p "$KATALOG_MODELI"
	printf '  pobieram wagi modelu „%s" do %s (około 480 MB dla „small"; jednorazowo)\n' \
		"$model" "$KATALOG_MODELI"
	if "$SRODOWISKO_MOWY/bin/python" - "$model" "$KATALOG_MODELI" <<'PYTHON'; then
import sys
from faster_whisper import WhisperModel

# device/compute_type są te same, które pomocnik ma wpisane na sztywno:
# instalacja jest CPU-only i taka ma pozostać.
WhisperModel(sys.argv[1], device="cpu", compute_type="int8", download_root=sys.argv[2])
PYTHON
		printf '  wagi pobrane; rozmiar zmierzony: %s\n' \
			"$(du -sh "$KATALOG_MODELI" 2>/dev/null | cut -f1)"
		printf '  rdzeniowi wskazać katalog ustawieniem %s = %s\n' \
			"$(wartoscMowy rozpoznanie.ustawienie-katalogu-modeli)" "$KATALOG_MODELI"
	else
		printf '  UWAGA: pobranie wag nie powiodło się — pierwsze użycie mikrofonu spróbuje ściągnąć je samo\n'
	fi
}

WYKAZ="$(pobierzWykaz)"
[ -n "$WYKAZ" ] || padnij "wykaz zależności jest pusty — sprawdź źródło (binarium rdzenia albo DANACO_WYKAZ_PLIK)"
WYKAZ_MOWY="$(pobierzWykazMowy || true)"

case "$TRYB" in
plan) trybPlan "$WYKAZ" ;;
sprawdz) trybSprawdz "$WYKAZ" ;;
postaw) trybPostaw "$WYKAZ" ;;
*) padnij "nieznany tryb: $TRYB (użyj: plan | sprawdz | postaw, także z minusami: --sprawdz)" ;;
esac
