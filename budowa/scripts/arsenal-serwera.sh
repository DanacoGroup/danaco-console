#!/usr/bin/env bash
# Arsenał serwera Danaco Console — prowizjonowanie kompletu programów, na których
# stoi rdzeń, na maszynie SERWERA docelowego.
#
# ── Po co ten skrypt ──────────────────────────────────────────────────────────
# Model wdrożenia (rozstrzygnięcie Właściciela, zamknięte): wszystkie programy
# jadą WRAZ Z APLIKACJĄ NA SERWER, a u Operatora stoi tylko cienka instalka
# (okno aplikacji, scripts/instalka-hybryda-win-x64.sh dla x64 oraz
# scripts/instalka-hybryda-win-arm.sh dla ARM64). Skutek: arsenał ma stać na
# serwerze, a wdrożenie ma go stawiać. Instalka Operatora ma zostać cienka i tego
# nie robi — od tego jest ten skrypt, uruchamiany na serwerze podczas wdrożenia.
# Bez arsenału funkcje odmawiają Operatorowi z braku programu (zewnetrzne.Wolaj
# → BrakNarzedzia), osobno przy każdym naciśniętym przycisku.
#
# ── Skąd bierze listę (żeby się nie rozjechała) ───────────────────────────────
# NIE przepisuje nazw pakietów ani nazw warstw. Woła binarium rdzenia w dwóch
# trybach wykazu i konsumuje wynik:
#     danaco-console --wykaz-zaleznosci   → komplet pozycji wykazu, po wierszu
#     danaco-console --wykaz-mowy         → arsenał mowy (piper, głosy, model)
# Źródłem jest jeden rejestr deklaracji narzędzi
# (server/internal/core/zaleznosci_zewnetrzne.go wraz z deklaracjami
# w adapterach), z którego bierze go też sonda startowa rdzenia. Zmiana pola
# `Pakiet` w deklaracji narzędzia dojeżdża więc i do sondy, i tutaj jednym
# ruchem — bez drugiej listy do ręcznego utrzymania.
#
# Wiersz wykazu zależności ma pola rozdzielone tabulacją:
#     warstwa <TAB> program <TAB> pakiet <TAB> stoi <TAB> nazwa <TAB> zakres
# Kolumnę `warstwa` liczy rdzeń (core.WarstwaZaleznosci) — skrypt jej nie
# odgaduje dopasowaniem napisów, bo rozdział warstw jest rozstrzygnięciem, nie
# formatowaniem, i ma jeden sprawdzian w Go.
#
# ── Gdy binarium rdzenia jest nieosiągalne (wykaz awaryjny) ───────────────────
# Wdrożenie zdarza się na maszynie, na której binarium jeszcze nie stoi (świeży
# serwer, wykaz czytany przed rozpakowaniem wydania), a bywa i tak, że stoi, lecz
# nie da się go uruchomić. Odmowa w tym miejscu znaczyłaby, że arsenału NIE MA
# CZYM postawić — a to jest właśnie ta chwila, w której trzeba go postawić.
# Dlatego skrypt niesie wykaz awaryjny: odpis kompletu deklaracji rdzenia w tej
# samej postaci wiersza, użyty TYLKO wtedy, gdy rdzenia nie ma czym zapytać, i za
# każdym razem zapowiedziany na wyjściu błędu — czytelnik ma wiedzieć, że patrzy
# na odpis, nie na wykaz policzony przez rdzeń.
#
# Wykaz awaryjny jest JEDNYM miejscem w tym skrypcie, z którego biorą go wszystkie
# trzy tryby (plan, sprawdz, postaw) — nie ma drugiej listy programów do
# sprawdzania obok listy pakietów do postawienia. Rozjazd z rdzeniem łapie tryb
# `sprawdz` uruchomiony przy dostępnym binarium: liczba pozycji i nazwy pakietów
# muszą wyjść te same. Zmiana deklaracji w rdzeniu ma dojechać tutaj tym samym
# ruchem — pole `Pakiet` przepisane błędnie kieruje Operatora do pakietu, którego
# nie ma.
#
# ── Co pakiet serwera niesie ──────────────────────────────────────────────────
# Warstwy przychodzą z wykazu; poniżej ich znaczenie, nie ich zawartość:
#   obowiazkowa-apt — pakiety dystrybucji, bez których moduły odmawiają:
#                     Tesseract OCR WRAZ z pakietem językowym polskim
#                     (tesseract-ocr-pol — stoi w polu `Pakiet` deklaracji),
#                     7-Zip, eSpeak NG, Pandoc, ffmpeg, ffprobe, LibreOffice,
#                     Chromium, ImageMagick, poppler, OpenSSH, ShellCheck,
#                     shfmt, picocom, telnet, łańcuch Go.
#   warsztat-go     — gopls, goimports, golangci-lint, staticcheck, Delve;
#                     moduł Developer pracuje na serwerze, więc jego warsztat
#                     też należy do serwera.
#   warsztat-npm    — Prettier, ESLint.
#   snap            — kubectl, PowerShell (moduł Terminal); moduły pwsh, np.
#                     PSScriptAnalyzer, to osobny krok Install-Module.
#   model-recznie   — Real-ESRGAN (wydanie z GitHuba) i rembg (środowisko
#                     pythonowe): wydania spoza repozytoriów dystrybucji,
#                     drukowane jako kroki ręczne z treścią pola `Pakiet`.
#   decyzyjna       — silnik kontenerów (docker/podman). Właściciel WSTRZYMAŁ go
#                     świadomie. Skrypt go NIE stawia i mówi o tym wprost;
#                     postawienie wymaga wyraźnego DANACO_SILNIK_KONTENEROW=tak.
#   arsenał mowy    — poza wykazem zależności stoi jeszcze: piper wraz z plikami
#                     głosów `.onnx`, biblioteka pythonowa rozpoznawania
#                     (faster-whisper z pomocniki/transkrypcja/wymagania.txt,
#                     uruchamiana pomocnikiem pomocniki/transkrypcja/transkrypcja.py)
#                     oraz WAGI MODELU pobierane z góry (patrz niżej).
#
# ── Wagi modelu rozpoznawania mowy: pobierane przy stawianiu serwera ──────────
# faster-whisper ściąga wagi przy pierwszym użyciu. Gdyby zostało tak na
# serwerze, pierwsze użycie mikrofonu u Operatora czekałoby na sieć — kilka minut
# ciszy przy pierwszym nagraniu. Dlatego prowizjonowanie pobiera wagi z góry,
# w rozmiarze domyślnym rdzenia (`mowa.ModelDomyslny`, dziś „small"; rozmiar
# przychodzi z wykazu mowy, nie jest tu wpisany).
#
# ROZMIAR POBRANIA — model „small" to około 480 MB na dysku (repozytorium
# Systran/faster-whisper-small; wagi float16, kwantyzacja do int8 dzieje się przy
# ładowaniu, więc pobranie nie jest mniejsze od plików repozytorium). Rozmiary
# pozostałych rozmiarów modelu rosną w tej samej skali — „tiny" i „base" są
# rzędu dziesiątek megabajtów, „medium" i „large-v3" rzędu gigabajtów. Skrypt po
# pobraniu mierzy katalog i wypisuje rozmiar zmierzony, żeby ta liczba nie była
# obietnicą, a pomiarem. Docelowy katalog: DANACO_KATALOG_MODELI (domyślnie
# /opt/danaco-arsenal/modele-mowy) — ten sam, który wskazuje się rdzeniowi
# ustawieniem `mowa_katalog_modeli`.
#
# ── Czego wymaga system operacyjny serwera ────────────────────────────────────
# Warstwa apt zakłada dystrybucję z `apt-get` (Debian/Ubuntu). Warstwa snap
# zakłada `snapd`. Warstwa Go zakłada `go` na ścieżce (pakiet `golang` stawia
# warstwa apt, więc kolejność warstw jest istotna: apt przed go). Warstwa npm
# zakłada `npm`. Warstwa mowy zakłada `python3` wraz z `python3-venv` i `pip`.
# Postawienie (`postaw`) wymaga uprawnień roota dla apt/snap. Tryby `plan`
# i `sprawdz` niczego nie zmieniają i nie wymagają roota.
#
# ── Jak zweryfikować sondą startową rdzenia ───────────────────────────────────
# Po postawieniu arsenału rdzeń przy starcie wypisuje do dziennika wiersz
# zbiorczy „zależności zewnętrzne: N z M obecnych" oraz osobny wiersz dla każdego
# braku wraz z zakresem, który przestaje działać, i podpowiedzią instalacyjną
# (core/zaleznosci_zewnetrzne.go, zglosZaleznosci). Kompletny arsenał to wiersz
# „M z M obecnych" bez wierszy braku. Ten sam stan bez uruchamiania rdzenia
# pokazuje `arsenal-serwera.sh sprawdz` — czyta obecność przez `command -v`
# i niczego nie instaluje. Gotowość samej mowy sprawdza pomocnik:
# `python3 pomocniki/transkrypcja/transkrypcja.py --wersja --model small`.
#
# ── Użycie ────────────────────────────────────────────────────────────────────
#   bash scripts/arsenal-serwera.sh plan       # (domyślnie) wypisz plan, nic nie rusza
#   bash scripts/arsenal-serwera.sh sprawdz    # sprawdź obecność, read-only
#   bash scripts/arsenal-serwera.sh postaw     # POSTAW arsenał na serwerze
# Tryb przyjmujemy w obu zapisach — `sprawdz` i `--sprawdz` — tak samo, jak rdzeń
# przyjmuje swoje znaczniki z jednym i z dwoma minusami (core.zadanoZnacznik).
# Tryb `sprawdz` kończy się kodem 1, gdy brakuje choć jednego programu wykazu:
# wdrożenie ma się na nim zatrzymać, a nie przeczytać braki i jechać dalej.
# Zmienne:
#   DANACO_RDZEN=/ścieżka/danaco-console    — binarium rdzenia wypisujące wykaz
#   DANACO_WYKAZ_PLIK=/ścieżka/wykaz.tsv    — gotowy wykaz zamiast wołania binarium
#   DANACO_WYKAZ_MOWY_PLIK=/ścieżka/mowa.tsv— gotowy wykaz mowy
#   DANACO_KATALOG_MODELI=/ścieżka          — katalog wag modelu mowy
#   DANACO_SRODOWISKO_MOWY=/ścieżka         — środowisko pythonowe rozpoznawania
#   DANACO_BEZ_WAG_MOWY=1                   — pomiń pobranie wag (instalacja bez sieci)
#   DANACO_SILNIK_KONTENEROW=tak            — postaw też WSTRZYMANY silnik kontenerów
#
# Skrypt nie pobiera głosów pipera ani wag Real-ESRGAN/rembg — to wydania spoza
# repozytoriów dystrybucji; wypisuje je jako kroki ręczne wraz z miejscem,
# w które mają trafić, i zmienną, którą można je wskazać.
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

# PAKIETY_POZA_WYKAZEM — pakiety apt, których w wykazie zależności NIE MA, a bez
# których arsenał serwera jest niekompletny. Każdy ma tu powód, bo pakiet bez
# powodu jest pakietem do wyrzucenia przy następnym czytaniu:
#   python3, python3-venv, python3-pip — interpreter i budowa środowiska
#       rozpoznawania mowy oraz środowiska rembg. Rdzeń nie woła interpretera
#       jako narzędzia, tylko pomocnika, więc w wykazie go nie ma.
#   nodejs, npm — nośnik warsztatu npm. Prettier i ESLint stoją w wykazie jako
#       programy (`prettier`, `eslint`), ale `npm i -g` nie ma czym ich postawić,
#       dopóki npm nie stoi; warstwa npm milcząco pomijała się na czystym serwerze.
#   sane-utils — program `scanimage`, warstwa SANE cyfryzacji w module Studio
#       (adapter_modul_studio_cyfryzacja.go deklaruje go osobno, poza wykazem
#       zależności rdzenia). Bez niego wykaz skanerów jest pusty, a Operator
#       dostaje odmowę przy każdym skanowaniu.
PAKIETY_POZA_WYKAZEM="python3 python3-venv python3-pip nodejs npm sane-utils"

# WYKAZ_AWARYJNY — odpis kompletu deklaracji rdzenia w postaci wiersza wykazu:
#     warstwa <TAB> program <TAB> pakiet <TAB> stoi <TAB> nazwa <TAB> zakres
# Używany tylko wtedy, gdy nie ma czym zapytać rdzenia (patrz nagłówek). Kolumna
# `stoi` niesie tu `?`, bo odpis nie jest pomiarem — obecność mierzy `command -v`
# w trybie sprawdz. Zakres skrócony do jednego zdania; pełne zdania stoją
# w deklaracjach (server/internal/core/zaleznosci_zewnetrzne.go).
#
# Pakiety są przepisane z pól `Pakiet` deklaracji i tylko stamtąd. W szczególności
# 7-Zip idzie z pakietu `7zip`, NIE z `p7zip-full`: tego drugiego w dystrybucji
# już nie ma i podpowiedź prowadziłaby donikąd (adapter_narzedzia_archiwum.go
# mówi to wprost).
czytajWykazAwaryjny() {
	# Rozdzielenie pól tabulatorem: $'\t' w literałach poniżej.
	cat <<-WYKAZ
		obowiazkowa-apt	pandoc	pandoc	?	Pandoc	zamiana formatów dokumentu (Studio, Translate, Library)
		obowiazkowa-apt	pdftotext	poppler-utils	?	poppler (pdftotext)	odczyt warstwy tekstowej PDF
		obowiazkowa-apt	pdftoppm	poppler-utils	?	poppler (pdftoppm)	rasteryzacja stron PDF przed OCR
		obowiazkowa-apt	tesseract	tesseract-ocr tesseract-ocr-pol	?	Tesseract OCR	rozpoznanie pisma ze skanu i zdjęcia
		obowiazkowa-apt	libreoffice	libreoffice	?	LibreOffice	formaty biurowe, których nie czyta Pandoc
		obowiazkowa-apt	ffprobe	ffmpeg	?	ffprobe	rozpoznanie zawartości nagrania
		obowiazkowa-apt	ffmpeg	ffmpeg	?	ffmpeg	zamiana formatu nagrania
		obowiazkowa-apt	7z	7zip	?	7-Zip	pakowanie i wydobycie archiwum
		obowiazkowa-apt	magick	imagemagick	?	ImageMagick	zapis obrazu w AVIF i WEBP stratnym, pomiar AVIF
		obowiazkowa-apt	espeak-ng	espeak-ng	?	eSpeak NG	odsłuch debaty syntezą mowy (Roundtable)
		obowiazkowa-apt	gofmt	golang	?	gofmt	formatowanie plików Go (Developer)
		obowiazkowa-apt	ssh	openssh-client	?	OpenSSH	powłoka zdalna i przekierowania portów (Terminal)
		obowiazkowa-apt	shellcheck	shellcheck	?	ShellCheck	analiza statyczna skryptów bash
		obowiazkowa-apt	shfmt	shfmt	?	shfmt	formatowanie skryptów bash
		obowiazkowa-apt	picocom	picocom	?	picocom	konsola portu szeregowego
		obowiazkowa-apt	telnet	telnet	?	Telnet (klient)	sesja Telnet do urządzenia sieciowego
		obowiazkowa-apt	chromium-browser	chromium-browser	?	Chromium	zrzuty stron, DOM, konsola, rejestr sieciowy (Browser)
		warsztat-go	goimports	go install golang.org/x/tools/cmd/goimports@latest	?	goimports	formatowanie Go wraz z porządkowaniem importów
		warsztat-go	gopls	go install golang.org/x/tools/gopls@latest	?	gopls	przejście do definicji i refaktoryzacje semantyczne
		warsztat-go	golangci-lint	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest	?	golangci-lint	analiza statyczna repozytorium Go
		warsztat-go	staticcheck	go install honnef.co/go/tools/cmd/staticcheck@latest	?	staticcheck	pogłębiona analiza statyczna Go
		warsztat-go	dlv	go install github.com/go-delve/delve/cmd/dlv@latest	?	Delve	debugowanie krokowe Go (Run & Debug)
		warsztat-npm	prettier	npm i -g prettier	?	Prettier	formatowanie TS, JS, CSS, Markdown, YAML
		warsztat-npm	eslint	npm i -g eslint	?	ESLint	analiza statyczna TS i JS
		snap	pwsh	powershell (snap) wraz z modułem PSScriptAnalyzer	?	PowerShell z modułem PSScriptAnalyzer	analiza i formatowanie skryptów PowerShell
		snap	kubectl	kubectl (snap)	?	kubectl	karta powłoki wewnątrz poda (Terminal)
		model-recznie	realesrgan-ncnn-vulkan	wydanie realesrgan-ncnn-vulkan z github.com/xinntao/Real-ESRGAN/releases rozpakowane do /usr/local/bin wraz z modelami w /usr/local/share/realesrgan-ncnn-vulkan/models, oraz mesa-vulkan-drivers dla liczenia na procesorze	?	Real-ESRGAN (ncnn)	powiększanie obrazu (Design)
		model-recznie	rembg	rembg[cli] w osobnym środowisku pythonowym, wystawiony opakowaniem /usr/local/bin/rembg ustawiającym U2NET_HOME=/usr/local/share/rembg-modele	?	rembg (U²-Net / ONNX Runtime)	wycinanie tła i rozkład obrazu na warstwy
		decyzyjna	docker	docker.io albo podman	?	Docker	wykaz kontenerów, budowanie obrazu, stos (Containers)
		decyzyjna	docker	docker.io	?	Docker (klient wiersza poleceń)	powłoka wewnątrz kontenera (Terminal)
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

# pobierzWykaz oddaje wiersze danych wykazu zależności (bez komentarzy).
pobierzWykaz() {
	local surowy rdzen
	if [ -n "${DANACO_WYKAZ_PLIK:-}" ]; then
		[ -r "$DANACO_WYKAZ_PLIK" ] || padnij "DANACO_WYKAZ_PLIK nieczytelny: $DANACO_WYKAZ_PLIK"
		surowy="$(cat "$DANACO_WYKAZ_PLIK")"
	elif rdzen="$(odnajdzRdzen)" && surowy="$("$rdzen" --wykaz-zaleznosci 2>/dev/null)" &&
		[ -n "$surowy" ]; then
		: # wykaz policzony przez rdzeń — źródło właściwe
	else
		# Rdzenia nie ma czym zapytać (świeży serwer przed rozpakowaniem wydania,
		# binarium nieuruchamialne, brak sieci przy pobieraniu wydania). Odmowa
		# zabrałaby jedyne narzędzie stawiania arsenału właśnie w tej chwili,
		# w której trzeba go postawić — jedziemy odpisem i mówimy o tym wprost.
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

# wartoscMowy wyjmuje jedną wartość z wykazu mowy po kluczu.
wartoscMowy() {
	printf '%s\n' "$WYKAZ_MOWY" | awk -F'\t' -v k="$1" '$1==k {print $2; exit}'
}

# pakietyApt zbiera tokeny warstwy obowiązkowej w kolejności pierwszego
# wystąpienia, bez powtórzeń — poppler-utils i ffmpeg padają w wykazie dwukrotnie
# (pdftotext i pdftoppm, ffmpeg i ffprobe), a apt ma je dostać raz.
pakietyApt() {
	printf '%s\n' "$1" | awk -F'\t' '$1=="obowiazkowa-apt" {print $3}' |
		tr ' ' '\n' | awk 'NF && !widziane[$0]++'
}

# wypiszWarstwe drukuje pozycje jednej warstwy: nazwa czytelna, program, pakiet.
wypiszWarstwe() {
	printf '%s\n' "$1" | awk -F'\t' -v w="$2" \
		'$1==w {printf "  %-34s [%s]  ← %s\n", $5, $2, $3}'
}

# policzWarstwe oddaje liczbę pozycji warstwy.
policzWarstwe() {
	printf '%s\n' "$1" | awk -F'\t' -v w="$2" '$1==w' | wc -l | tr -d ' '
}

# obecny mówi, czy program jest osiągalny na ścieżce — tak samo, jak sprawdza to
# sonda rdzenia. Sprawdzenie, nie instalacja.
obecny() { command -v "$1" >/dev/null 2>&1; }

# ── Plan ──────────────────────────────────────────────────────────────────────
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

	zglos "WARSZTAT Go — moduł Developer pracuje na serwerze [$(policzWarstwe "$wykaz" warsztat-go)]"
	wypiszWarstwe "$wykaz" warsztat-go
	zglos "WARSZTAT npm [$(policzWarstwe "$wykaz" warsztat-npm)]"
	wypiszWarstwe "$wykaz" warsztat-npm
	zglos "SNAP [$(policzWarstwe "$wykaz" snap)]"
	wypiszWarstwe "$wykaz" snap
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

# ── Sprawdzenie ───────────────────────────────────────────────────────────────
# Nic nie instaluje i nie wymaga roota. Wykaz programów bierze z tego samego
# źródła, z którego biorą go plan i postaw (wykaz rdzenia albo wykaz awaryjny) —
# osobnej listy „co sprawdzić" nie ma, bo rozjechałaby się z listą „co postawić".
#
# Kod wyjścia: 1 przy jakimkolwiek braku programu wykazu, 0 przy komplecie.
# Wdrożenie ma się na tym zatrzymać. Braki arsenału mowy (głosy pipera, wagi
# modelu) są wypisywane, ale kodu nie zmieniają: to pliki, nie programy na
# ścieżce, a synteza ma zejście na eSpeak NG z wykazu.
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

	# Programy poza wykazem zależności rdzenia, a należące do arsenału serwera
	# (powód każdego stoi przy PAKIETY_POZA_WYKAZEM). Liczone osobno, żeby liczba
	# „N/M" pozostała liczbą wykazu rdzenia i dała się zestawić z wierszem sondy
	# startowej. Do kodu wyjścia wchodzą, bo bez nich serwer też nie jest kompletny.
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

# ── Postawienie ───────────────────────────────────────────────────────────────
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
	# Do pakietów wykazu dochodzą pakiety spoza wykazu (interpreter i venv mowy,
	# nośnik warsztatu npm, warstwa SANE) — powód każdego stoi przy
	# PAKIETY_POZA_WYKAZEM, jednym miejscem dla planu, postawienia i sprawdzenia.
	# shellcheck disable=SC2086
	apt-get install -y $apt $PAKIETY_POZA_WYKAZEM

	zglos "Warstwa warsztatu Go"
	if command -v go >/dev/null 2>&1; then
		# GOBIN kierujemy do /usr/local/bin, bo domyślne ~/go/bin należy do roota
		# odpalającego prowizjonowanie i NIE JEST na ścieżce procesu serwera.
		# Program postawiony tam, gdzie go nikt nie widzi, to dla rdzenia brak
		# programu: sonda woła `command -v`, nie zgaduje katalogów.
		printf '%s\n' "$wykaz" | awk -F'\t' '$1=="warsztat-go" {print $3}' |
			while read -r polecenie; do
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
