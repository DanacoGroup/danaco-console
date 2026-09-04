# Przekazanie prowadzenia budowy — stan na 2 września 2026

Dokument opisuje stan zastany, nie przebieg prac. Mówi, co stoi, co jest otwarte
i czego nie wolno ruszyć.

## 1. Czym rzecz stoi

Rdzeń działa. Interfejs prowadzi od wejścia przez przedsionek do Centrum i Studia;
pozostałe moduły stoją na zapowiedzi niegotowości. Praca stoi na `main`: gałąź
`teren/naprawy-audytu` weszła przesunięciem prostym 4 września 2026 po pełnym
przebiegu `go test ./server/...` zamkniętym kodem 0. Punkt cofnięcia sprzed
scalenia to `2881c425`.

Cztery warstwy budowy sprawdzają się czysto: rdzeń Go (`go build`, `go vet`,
`gofmt`, komplet sprawdzianów), klient TypeScript (`tsc --noEmit`, `vite build`),
powłoka Tauri i instalator Tauri — obie także w budowie krzyżowej na
`x86_64-pc-windows-gnu`.

| miara | wartość |
|---|---|
| komendy kontraktu | 1086 |
| komendy wykonalne z klienta | 1086 z 1086: 485 z własnym wołającym, reszta przez katalog operacji |
| tabele w bazie po przejeździe migracji | 421 |
| tabele niosące pracę Operatora | 373, w tym 172 korzenie |
| granica konta w zapytaniach | domknięta; 15 miejsc rozstrzygniętych, nie dziur |
| kroki migracji | do 497 |
| `go test ./server/...` | przechodzi w całości |
| prototypy okien w `design/05-okna/` | 35; okien w produkcie 8 |
| droga wejścia | rejestracja, logowanie, metody, urządzenia, sesje bramki |

## 2. Gałęzie, dokumenty, materiał

| co | gdzie |
|---|---|
| gałąź prac | `teren/naprawy-audytu` |
| stan sprzed napraw | rewizja `2881c425` na `teren/centrum-poprawki` |
| rejestr rozstrzygnięć | `prowadzenie/decyzje.md` — 33 pozycje, wiążą bezwzględnie |
| audyt pierwszy | `~/robocze/przekazanie-2026-09-02/audyt-pierwszy.md` — 189 ustaleń, plan siedmiu etapów |
| audyt powtórny, część | `~/robocze/przekazanie-2026-09-02/audyt2-czesciowy.md` — 7 z 13 kategorii, 79 ustaleń |
| materiał granicy konta | `~/robocze/przekazanie-2026-09-02/` — rozpoznanie korzeni, pozycje otwarte, przebiegi |

Audyt powtórny stanął na 7 z 13 kategorii. **Nie ruszyły**: rdzeń Go, klient
TypeScript, warstwa projektowa, powłoka i instalator, poczta i droga wejścia,
braki produktowe. Skrypt przebiegu leży w materiale przekazania.

## 3. Granica konta — domknięta

**Miara.** Platforma prowadzi wiele kont od migracji 406. Granica konta objęła do
tej sesji cztery tabele: trzy tabele bramki i karty sesji. Rozpoznanie wszystkich
418 tabel wykazało **373 tabele niosące pracę Operatora**; z nich **172 to
korzenie** (wiersz należy do Operatora wprost), 178 to dzieci wiszące na
korzeniach kluczem obcym, 22 to łączniki.

**Mechanizm.** Wskazanie konta stoi w korzeniu; dziecko granicę dziedziczy drogą
po kluczu obcym. Warunek jest jeden na całą bazę:

- `dane.WarunekKonta` — zawężenie wiersza do konta żądania, jeden argument
  zapytania: `dane.KontoOperatora(ctx)`. Porównanie idzie przez `IS`, więc
  instalacja przed rejestracją widzi pracę zastaną.
- `dane.WskazanieKonta` — zapis konta w wierszu zakładanym.
- Wzór dla dziecka: `AND EXISTS (SELECT 1 FROM <korzeń> k WHERE k.id =
  <dziecko>.<klucz_obcy> AND <WarunekKonta z aliasem k>)`.

Konto wchodzi do kontekstu raz, w `core/rdzen.go` przed rozdziałem na
obsługiwacze, więc warstwa danych bierze je z kontekstu, nie argumentem.

**Kroki migracji tej sesji.** 484 dokłada `konto_id` 171 korzeniom (`kanal_modelu`
ma własny). 485 rozstrzyga zderzenie nazw — `kanal_modelu.konto_id` znaczyło
konto dostawcy i nosi teraz nazwę `konto_dostawcy_id` — oraz przebudowuje siedem
wskaźników jednoznacznych na kontowe.

**Co zrobione.** Zapytania 47 korzeni o najwyższej wadze — magazyny poświadczeń
i treści Operatora: `punkt_dostepu`, `terminal_klucz`, `terminal_host`,
`host_zdalny`, `konto`, `kanal_modelu`, `skrzynka_pocztowa`,
`sekret_rozszerzenia`, `poswiadczenie_automatyki`, `dokument_studio`,
`blok_wiadomosci`, `plik_biblioteki` i pozostałe. Zawężonych zapytań: 301
w pierwszym przejściu, 106 dziur domkniętych w drugim.

**Runda trzecia (2 września, wieczór).** Trzynaście terenów domknęło pozycje
otwarte 53 plików i zapisało konto w każdym `INSERT` do korzenia w tych plikach
(rewizja `a9ea0447`); sprawozdania wykonawców co do pliku i stałej:
`~/robocze/przekazanie-2026-09-02/runda-3/wyniki-terenow-01-13.txt`, instrukcja
dokończenia: `runda-3/STAN.md`. Kontrolerzy terenów NIE zdążyli — sesja stanęła na
limicie; rewizja `a9ea0447` przeszła drabinę szybką (budowa, vet, walidator), nie
przeszła sprawdzianów ani kontroli. Teren 14 (nastawy per konto przez
`konfig.Kontekst.KontoOperatora`) przerwany w połowie, stan zapisany rewizją `abc78ea8`
— rozstrzygacz z zerem konta czyta konto najstarsze, więc stan częściowy nie cofa
zachowania. Teren 15 (własność kanału przy użyciu, 11 miejsc w `core`, wykaz
`channel.list` po koncie) nie ruszył.

**Reguła Właściciela dla rund (2 września).** Agent dostaje pracę wartą co najmniej
500K tokenów wyniku i ma dać więcej wyniku, niż przeczyta; rundy dzielić na 3–5
szerokich terenów, nie na kilkanaście wąskich.

**Stan po sesji 4 września 2026.** Granica konta jest domknięta.

- **Zapytania.** Z 173 zapytań sięgających korzenia bez warunku konta zostało 15,
  wszystkie rozstrzygnięte jako praca procesu bez zamawiającego albo odczyt
  kluczem własnym wiersza wewnątrz transakcji, która ten wiersz zapisała
  (rejestr decyzji, poz. 37). Miarę odtwarza skrypt liczący literały SQL sięgające
  korzeni z `korzenie-granicy-konta.json`.
- **Korzenie bez kolumny.** Krok 489 dołożył `konto_id` kolejce, raportowi badania
  i przebiegowi wsadu Studia; okno komunikacji sięga konta drogą przez sesję do
  karty sesji wspólnym `dane.warunekKontaOkna` (poz. 36).
- **Więzy UNIQUE.** Kroki 490–494 przeniosły do konta jednoznaczność każdego
  klucza wpisywanego przez Operatora. Więz na kluczu nadawanym przez rdzeń
  zostaje globalny (poz. 38).
- **Tabele bez drogi do konta.** Krok 495 rozbił przestrzeń badania z jednego
  wiersza na instalację na wiersz na konto i dołożył wskazanie konta pamięci
  tłumaczeń.

**Co zostaje otwarte.**

1. Nastawa pracy Studia (`nastawa_pracy_studio`) własnej kolumny konta nie ma;
   sięga go kodem okna, który nadaje rdzeń. Wystawienia to nie tworzy, ale
   rozstrzygnięcie warto zapisać przy najbliższej pracy nad Studiem.
2. Rodzaj znacznika Studia zostaje z więzem jednoznaczności na całą tabelę:
   jego `nazwa` jest celem klucza obcego z `znakowanie_studio` (poz. 38).

## 4. Naprawy krytyczne — stan

| co | stan |
|---|---|
| sekret nawiązania wykluczał serwer narzędzi | naprawione; poświadczenie zastępuje sekret (decyzja 31) |
| ping zrywał bezczynne gniazdo serwera narzędzi | naprawione; gniazdo czytane biegiem odczytu, ponowienie na gnieździe zastanym |
| sumy migracji uzgadniane po numerze | naprawione; nazwa kroku sprawdzana przy każdym starcie (decyzja 32) |
| trzy produkcyjne maszyny w zaczynie każdej instalacji | naprawione krokiem 483 (decyzja 33) |
| znacznik „bramka bez poczty" jeden na instalację | naprawione; znacznik na konto, wpis zastany przejmowany raz |
| klucze obce wskazywały tabele przejściowe dawnych przebudów | naprawione krokiem 497; zapis wypowiedzi debaty kończył się odmową „no such table: main.debata_tura_nowa" i tura nie zbierała ani jednej wypowiedzi |

Termin gniazda przed bramką, dołożony w trakcie prac, został **zdjęty**: zabijał
gniazdo serwera narzędzi i rejestrację z potwierdzeniem listem. Gniazdo
z poświadczeniem nie liczy się odtąd do granicy gniazd niezwiązanych; samo
ustalenie audytu o granicy 64 zostaje otwarte.

## 5. Ustalenia audytu powtórnego — stan

Pełny wykaz: `audyt2-ustalenia-7-kategorii.json`.

**Zamknięte 4 września 2026.**

| ustalenie | co zrobiono |
|---|---|
| poświadczenie serwera narzędzi w wierszu poleceń procesu | wpis MCP podaje je zmienną `DANACO_POSWIADCZENIE_NARZEDZI`; wiersz poleceń czyta każdy program użytkownika, środowisko — tylko właściciel procesu |
| `progress.changed` idzie do gniazd przed bramką | telemetria idzie `RozglosPoBramce`; gniazdo przed zalogowaniem jej nie dostaje |
| wskaźnik PIN-u nie zna konta | krok 496; para (urządzenie, rodzaj) jest jednoznaczna w koncie |
| `ZamknijKodem` daje 5 s na każde zamknięcie gniazda | własny termin 500 ms, potem gniazdo schodzi twardo |
| `zloz.mjs` ślepy na pozycję za hasłem | odbiór idzie z poświadczeniami kanału, gdy środowisko je niesie; bez nich pozycja melduje się jako NIESPRAWDZONA |
| poświadczenia kanału w `option_env!` bez wznowienia budowy | zastane naprawione: `build.rs` ma `cargo:rerun-if-env-changed` na obu zmiennych |
| wykaz kieruje na plik z 18 sierpnia | zastane naprawione: `wydania.json` niesie wydanie 2.0.0 z 1 września |
| pokrycie klienta 5,34 % | zastane naprawione: interfejs woła 1086/1086 komend kontraktu |
| sekret nawiązania nie jest przez powłokę wytwarzany ani przekazywany | zastane naprawione: `desktop/src-tauri/src/sekret.rs` losuje go na uruchomienie, `rdzen/proces.rs` podaje rdzeniowi zmienną `DANACO_SEKRET_NAWIAZANIA` |

**Otwarte.**

- Pakiet serwera wdrożenia 2.0.0 jest złożony (`scripts/pakiet-serwera.sh`,
  `wydania/2.0.0-2026-09-04/danaco-console_2.0.0_amd64.deb`, 40,4 MB) i wykaz
  wydań kieruje na niego zamiast na 1.0.0 z 18 sierpnia. Pliku nie ma jeszcze
  w kanale: wgranie należy do operatora wydania, a złożenie witryny melduje tę
  pozycję jako NIESPRAWDZONĄ, póki poświadczeń kanału nie ma w środowisku.
  Dwie pozycje Windows z tego samego wykazu są natomiast sprawdzone ruchem:
  kreator uruchomiony na tej maszynie pobrał z kanału
  `Danaco Console_2.0.0_hybryda_x64-setup.exe` — 2 598 033 bajty, instalka NSIS
  PE32 — więc kanał, wykaz i pobranie działają. Braku dotyczy wyłącznie `.deb`.
- Rodzina `notification.*` ma wołającego własnego: dzwonek szyny otwiera centrum
  powiadomień (rozstrzygnięcie 46).
- Pięć opracowań, których spis żądał, a repozytorium nie niosło, wróciło z jego
  własnej historii — rewizja `62c8f0de` wniosła dwanaście opracowań interfejsu,
  `fad6a36d` zdjęła je przy zakładaniu ustroju przebudowy, a na `main` wróciło
  siedem. Pozostałe pięć podjęte stamtąd, nie z materiału zabezpieczonego, więc
  źródłem jest repozytorium. Objętości w spisie poprawione; rubryka objętości
  jest zdjęciem z chwili wpisu i edycja dokumentu jej nie przelicza.
- Okno produktu niesie 371 kształtów ikon, z czego 195 stoi poza zestawem
  (`design/zasoby/ikony/`, 152 pliki, 262 kształty). Zestaw nie wiąże znacznika
  (rozstrzygnięcie Właściciela z 28 sierpnia), więc to miara do zamknięcia
  w warstwie projektowej, nie usterka produktu.

## 5a. Bramka etapu 2 — czym każde kryterium jest wykazane

| kryterium | stan | czym wykazane |
|---|---|---|
| 1. produkt złożony, nie zestaw procesów | wykazane poza instalką | powłoka złożona ze wskazaniem lokalnym i uruchomiona sama jedna stawia rdzeń procesem pobocznym z sekretem nawiązania, otwiera okno i ikonę zasobnika; przejście instalki wymaga maszyny z Windows |
| 2. pełne działanie modułu od kliknięcia do zapisu | wykazane | przejście w przeglądarce: przedsionek, okno Studia, wiadomość, dokument, treść z kanwy leży w `dokument_studio` |
| 3. zerwanie i powrót nie gubią stanu | wykazane | tryb bez sieci i powrót: okno zostaje w Centrum |
| 4. typy z kontraktu, rozjazd wykrywany maszynowo | wykazane | drabina: świeżość wytworów wobec `contract.json`, pokrycie 1086/1086 |
| 5. Studio odpowiada prototypowi | do oceny Właściciela, materiał zebrany | 158 ze 175 klas prototypu Studia stoi w żywym oknie; 17 pozostałych to stany wywoływane, każdy z wiązaniem albo regułą arkusza |
| 6. okno obecne nazywa swoją niegotowość | wykazane | `#cd-mobile` odpowiada zdaniem „To okno nie wchodzi do tego wydania" |
| 7. przejście bez ślepego zaułka | wykazane | uruchomienie → bramka → logowanie → Centrum → przedsionek → Studio, zero błędów konsoli |

Czego przy tym **nie** udało się rozstrzygnąć: okno główne powłoki nie zmapowało
się na ekranie wirtualnym, choć dziennik powłoki notuje „okno otwarte", a klient
nawiązał z rdzeniem (rdzeń odnotował przyłączenie). Okno wstaje ukryte i pokazuje
je sama strona, `aplikacja.ts`, po zmontowaniu ekranu startowego. Dwie przyczyny
tłumaczą to tak samo dobrze i tej maszyny nie da się między nimi rozstrzygnąć:
silnik widoku nie wywołuje `requestAnimationFrame` w oknie nigdy niepokazanym
albo `show()` na oknie bezramkowym nie mapuje się bez menedżera okien, a tej
maszyny żaden nie ma i instalować nie wolno. Zapora czasu w `aplikacja.ts`
obiecywała w komentarzu, że usterka montażu nie zostawi okna niewidocznym na
zawsze, a wybierała tylko chwilę wywołania klatki — klatka została drogą pierwszą,
zegar zapasową. Na Windows sprawę rozstrzyga pierwsze uruchomienie powłoki:
gdyby okno nie wstało, wraca się do niego pozycją „Pokaż okno" w zasobniku.

Kreator instalacji też przeszedł na tej maszynie, pod `xvfb-run`. Pięć kroków
z sześciu wykonał do końca: odczytał parametry urządzenia (Ubuntu 26.4.0,
Intel/AMD x64, 35,0 GB wolnego), wygasił „Dalej" do czasu akceptacji licencji,
zaproponował katalogi dopiero po podaniu `LOCALAPPDATA` — bez tej zmiennej mówi
wprost „Odmowa: zmienna środowiskowa LOCALAPPDATA nie jest ustawiona", zamiast
zmyślać ścieżkę — prawo zapisu potwierdził próbą zapisu, po czym pobrał z kanału
instalkę wydania. Zatrzymał się na kroku 5 z powodem `system-nieobslugiwany`:
„Instalka wydania jest plikiem wykonywalnym Windows — na tym systemie nie ma jej
czym uruchomić". Nieprzebadane zostaje wyłącznie to ostatnie uruchomienie.

Kryterium 1 wykazane na tej maszynie, bez Windows: powłoka złożona
`cargo build --release` ze wskazaniem `DANACO_HOST_WDROZENIA=127.0.0.1`,
`DANACO_PORT_WDROZENIA=17911` i rdzeniem położonym obok pliku wykonywalnego,
uruchomiona pod `xvfb-run`, postawiła rdzeń procesem pobocznym (dziennik powłoki:
„rdzeń: proces poboczny … wystartowany z sekretem nawiązania"), otworzyła port
17911 i okno z ikoną zasobnika. Bez wskazania powłoka nie zgaduje: mówi, że czeka
na podanie serwera wdrożenia.

Do kryterium 5 zebrany materiał, którego ocena wizualna nie zastąpi, ale który
mówi, gdzie patrzeć. Ze 175 klas własnych prototypu Studia 158 stoi w oknie
otwartym na żywo; brakujące siedemnaście — `dn-wersja*`, `dn-diff-*`, `dn-kartka*`,
`dn-zazn`, `dn-plyw`, `dn-stan-miara`, `dn-stan-tor`, `st-obudowa`, `st-miara`,
`st-szyna-grupa` — to stany wywoływane: pojawiają się przy wersjach dokumentu,
różnicy, podglądzie stron i zaznaczeniu. Każda z nich ma albo wiązanie, które
ją stawia (`studio.ts`, `studio-repozytorium.ts`, `studio-sledzenie.ts`), albo
regułę w arkuszu (`studio.css`, `komponenty.css`, `rama.css`). Pomiar prowadzi
się na żywym dokumencie, nie na `dist/index.html`: rama okna powstaje skryptem
warstwy projektowej, więc w pliku wydania jej nie ma.

Kryteria 2, 3, 6 i 7 wykazane klikaniem po naprawie dwóch usterek, które wyszły
dopiero przy tym przejściu: pustego okna „Operacje platformy" i wygaszonego
przycisku nowego dokumentu.

Czynności wiersza sesji w Centrum przeszły tą samą drogą. Menu wystawia dokładnie
pięć pozycji mających pokrycie w komendach — reszta znaczników prototypu schodzi
z okna, zamiast stać martwa. Zmiana nazwy wpisana wprost w wierszu leży w kolumnie
`tytul`, przeniesienie zakłada projekt ze wskazaniem konta operatora i wiąże z nim
sesję, archiwizacja przestawia `stan` na `archiwalna`. Każde bez błędu konsoli.

Kafle komponentów przeszły tak samo: kliknięcie kafla rodzaju zakłada komponent
własny, nazwa wchodzi wprost w pozycji wykazu, a menu wystawia „duplikuj" i „usuń"
— eksportu i importu z prototypu tam nie ma z rozstrzygnięcia 40. Powielenie
odkłada drugi wiersz z przyrostkiem „— kopia" i własnym bytem docelowym.

Wiersz projektu wystawia „nazwa" i „usuń"; oba doszły do bazy. Zmiana nazwy
przedtem ginęła po cichu — powrót ogniska z menu domykał pole wpisu zdarzeniem
`blur`, więc wpis rozstrzygał się wartością pustą i żadna komenda nie wychodziła.
Ścieżka czeka teraz obrót pętli i odsłania panel, tak jak ścieżka zakładania.

Przedsionek TalkIn wystawia sesje, które zwraca `environment.enter`, i zakłada
z siebie projekt ze wskazaniem konta. Szyna sesji zwija się przy wąskim oknie,
więc przejście przeglądarką prowadzi się na szerokości Właściciela (2560 px):
przy domyślnej szerokości karty wiersze są w drzewie, ale poza widokiem, co
łatwo wziąć za pusty wykaz.

Cztery środowiska przeszły klikaniem. TalkIn, WorkSpace i CodeStudio otwierają
przedsionek z wykazem sesji, który zwraca `environment.enter`. MultitaskingAI
nie otwiera go i otwierać nie ma: pusty wykaz modułów jest tam kształtem
zamierzonym (migracja 072 — `navigationKind: orchestration`), a panel orkiestracji
nie wchodzi do pakietu pięciu okien. Zdanie odpowiedzi obwiniało wcześniej rejestr
rdzenia, więc czytało się jak brak danych; teraz nazywa rzecz i wskazuje katalog
operacji.

Dzwonek szyny narzędziowej był afordancją martwą — odpowiadał samą etykietką.
Powstało centrum powiadomień: plakietka na dzwonku liczy zdarzenia nowe w całym
rejestrze, panel wystawia je od najnowszego z klasą, wagą i stanem, a przy każdym
stoją trzy czynności kontraktu — odczytanie, zamknięcie i odłożenie o godzinę.
Sprawdzone klikaniem: zamknięcie zdarzenia przestawia jego stan na `obsluzone`
w tabeli `powiadomienie_centrum`, drugie zdarzenie zostaje nowe. Plakietka słucha
`notification.raised` i `notification.changed`, więc licznik schodzi bez
odświeżania okna.

Dziewięć modułów przedsionka TalkIn otwarto po kolei — Studio, Workspace, Browser,
Research, Library, Translate, Roundtable, Assistant, Agents. Każdy staje z treścią
(od 405 do 5144 znaków, od 49 do 85 pozycji do naciśnięcia) i żaden nie zgłasza
błędu konsoli. Pustego okna, jak przy „Operacjach platformy" przed naprawą,
nie ma ani jednego.

Wstążka okna roboczego w Studiu przeszła to samo przemiatanie. Cztery przyciski
układu — podział pionowy, maksymalizacja okna roboczego, nowe okno pomocnicze
i okno komunikacji — nie miały obsługi nigdzie: ani w wiązaniu, ani w bibliotece
warstwy projektowej; opracowania opisują okna pomocnicze ogólnie, ale tych czynności
nie określają. Zostały w oknie i nazywają swoją niegotowość, bo zdjęcie zmieniłoby
wstążkę wobec prototypu, którą Właściciel ocenia w kryterium 5.

## 6. Rozstrzygnięcia czekające na Właściciela

Czynność, nie rozstrzygnięcie: zakup certyfikatu OV w Certum na dokumenty spółki
i konto SimplySign. Do zakupu wydania idą z `DANACO_PODPIS=pomijany`.

Pozycji audytu bez rozstrzygnięcia już nie ma. Pięć ostatnich zamknęły
rozstrzygnięcia 41–45 z upoważnienia „Rozstrzygaj za mnie": los okien
platformowych, kształt rodziny `project.*`, reguła wykazu narzędzi modelu,
kolejność wpinania podgrup `studio.*` i los kanału `/wydania/` za hasłem.
Trzy z nich okazały się bezprzedmiotowe po pomiarze — rzecz stała już
rozstrzygnięta w źródle, brakowało tylko zapisu.

Naprawa uszkodzeń po dawnych przebudowach (w tym kroku 269) weszła krokiem 497
z upoważnienia „Rozstrzygaj za mnie" — zapis wypowiedzi debaty nie działał
w ogóle, więc zwłoka kosztowałaby moduł, nie tylko porządek.

## 7. Czego nie wolno

- **Nie edytować migracji zastosowanych.** Strażnik sumy kontrolnej wywróci start
  każdej istniejącej bazy. Naprawa schematu idzie wyłącznie nowym krokiem
  o numerze wyższym niż 485.
- **Nie odtwarzać prototypów.** 27 arkuszy i 79 skryptów w `design/zasoby/` oraz
  35 prototypów okien w `design/05-okna/` są biblioteką do wpięcia. Kompozycji
  okna nie przekazuje się zleceniem.
- **Nie budować `klient-poprzedni`** w drzewie repozytorium — `emptyOutDir`
  kasuje pliki terenów roboczych.
- **Nie stawiać nowych testów, bramek i walidatorów.** Sprawdzian pisze się
  wtedy, gdy jest jedynym sposobem wykazania, że naprawa działa.
- **Nie pytać Właściciela o rzeczy rozstrzygalne ze źródeł.** Rozstrzygnięcia
  wiążące stoją w `prowadzenie/decyzje.md`.
- **Nie wystawiać niczego jako gotowego bez pokazania wyniku.** Zrzut z maszyny
  Windows z pliku pobranego z serwera, nie z dysku maszyny budującej.

## 8. Maszyny

**Maszyna budująca** — ta. Podgląd budowy wychodzi pod `http://51.75.62.180/`,
tylko port 80; Caddy kieruje na `127.0.0.1:17896`, gdzie stoi rdzeń podglądowy
z `/tmp/rdzen-podglad`, katalogiem danych `~/robocze/podglad-dane` i pakietem
interfejsu z `budowa/klient/dist`. **To proces, nie usługa** — ginie po restarcie
maszyny.

**danaco-system — 57.128.253.74.** Rdzeń wdrożenia (`danaco-console.service`,
port 17870, dane `/var/lib/danaco-console`), rdzeń wydania Studio
(`danaco-console-studio.service`, 17871), kanał pobrań `pobierz.danaco-group.pl`,
portal, serwer poczty Stalwart. Dostęp: klucz `~/.ssh/danaco_operator`, `root`.
Hasło skrzynki nadawczej stoi w `/etc/danaco-console/srodowisko` na tamtej
maszynie — nie w repozytorium.

**Uwaga wiążąca dla wdrożenia.** Baza `danaco-system` idzie inną linią numeracji:
jej kroki 406–480 noszą nazwy, których repozytorium nie zna. Po naprawie sum
kontrolnych rdzeń **odmówi na niej startu** z czytelnym komunikatem, zamiast paść
kilkadziesiąt kroków dalej na braku kolumny. Przed wdrożeniem rejestr tamtej bazy
trzeba doprowadzić do linii repozytorium; kopie sprzed ręcznych zmian leżą obok
jako `*.bak-2026-09-01-2340`.

**Maszyna Windows do sprawdzania.** Kontener `danaco-win` (dockurr/windows),
noVNC na `http://127.0.0.1:8006/`, przez `sudo docker`. Katalog wspólny
`/srv/win-danaco/wspolny` widziany jako `\\host.lan\Data`. Sterowanie Playwright
przez CDP, `executablePath: '/usr/bin/chromium-browser'`. **Pułapka:** `Alt+F4`
na pulpicie wywołuje zamknięcie systemu i kontener pada.

## 9. Jak sprawdzić, że wszystko stoi

```bash
cd ~/budowa/budowa && go build ./... && go vet ./server/... && gofmt -l server
cd ~/budowa/budowa && go test -count=1 ./server/internal/dane/... ./server/internal/store/... ./server/internal/transport/... ./server/internal/narzedzia/...
cd ~/budowa/budowa/klient && npx tsc --noEmit && npx vite build
cd ~/budowa/budowa/desktop/src-tauri && PATH=$HOME/.local/bin:$PATH cargo check --target x86_64-pc-windows-gnu
python3 ~/budowa/narzedzia/pokrycie-kontraktu.py
```

Pełny `go test -count=1 -timeout 50m ./server/...` przechodzi w całości i tak go
sprawdzać przed scaleniem; cały przebieg idzie w nim niecałe dziewięć minut.
Wcześniej szedł ponad pół godziny: profil wykazał, że dwie trzecie czasu zjadało
parsowanie SQL migracji przy każdym z 433 montaży, więc uprząż bierze teraz bazę
z wzorca złożonego raz na przebieg pakietu. Dziennik startu rdzenia ma mówić
`komend=1086`, `zależności zewnętrzne: 55 z 55 obecnych` i `klient=klient/dist`.

**Czym jest pokrycie 1086/1086.** Miara liczy komendy wykonalne z interfejsu.
485 z nich ma wołającego własnego — panel, przycisk albo wykaz zbudowany pod tę
jedną komendę. Pozostałe wykonuje katalog operacji modułu
(`klient/src/wiazanie/katalog-modulu.ts`): Operator wybiera komendę rodziny,
dokłada parametry JSON, a katalog wysyła żądanie i melduje odmowę rdzenia.
To wołający prawdziwy, nie wpis w wykazie — ale nie jest tym samym co panel
zbudowany pod komendę. Wykazane uruchomieniem: okno „Operacje platformy" (przycisk
konfiguracji w Centrum) składa 67 katalogów, żaden pusty, razem 1086 komend,
a wybrana z katalogu `session.list` wykonuje się bez odmowy. Panele dla modułów spoza łańcucha etapu 2 nie powstają
z rozstrzygnięcia planu: okno zapowiedziane nie staje jako pusta skorupa.

**Okno produktu w przeglądarce.** Osobny rdzeń serwuje pakiet interfejsu, a
Chromium bez okna (puppeteer-core spod `pa11y`) przechodzi drogę: uruchomienie,
bramka logowania, logowanie kontem, menu aplikacji, Centrum dowodzenia ze
środowiskami. Miarą jest brak błędów konsoli i odmów sieci — ostatni przebieg
dał zero i zero. Ekrany stoją w jednym drzewie, więc szukać pól widocznych
(`offsetParent !== null`), nie pierwszych w DOM.

Kryterium 3 wykazane tą samą drogą: przełączenie karty w tryb bez sieci
(`setOfflineMode`) zostawia okno w Centrum, a powrót sieci wraca do Centrum, nie
na bramkę logowania — sesja przeżywa zerwanie.

Kryterium 6 bramki etapu 2 wykazane tą samą drogą: naciśnięcie okna spoza wydania
(`#cd-mobile`) odpowiada zdaniem „Mobile — To okno nie wchodzi do tego wydania",
nie ciszą i nie pustą skorupą.

Kryterium 2 wykazane tą samą drogą, klikaniem: wejście `[data-nowa-sesja-srodowisko]`
otwiera przedsionek TalkIn, kafel Studia stawia okno modułu, pierwsza wiadomość
zakłada stanowisko, przycisk zakłada dokument, a treść wpisana w kanwie i odłożona
skrótem Ctrl+S leży w `dokument_studio` — sprawdzone zapytaniem do pliku bazy.
Przycisk szyny (`button[data-srodowisko]`) sam środowiska nie otwiera: wybiera je.

**Przegląd odczytów na żywym rdzeniu.** Osobny rdzeń na własnym porcie i katalogu
danych, po nim wywołanie każdej komendy wykazu i odczytu: 140 bez pól wymaganych
i 68 zawężonych oknem albo sesją. Odmowa `internal_error` znaczy usterkę rdzenia,
brak danych — nie. Ostatni przebieg: 180 udanych, zero odmów wewnętrznych;
`developer.lint.get` przekracza czas, bo uruchamia analizatory na katalogu
roboczym okna.

## 10. Pułapki tej maszyny

- Cel `x86_64-pc-windows-msvc` buduje się przez `cargo xwin`, a `cc-rs` szuka
  `llvm-lib` bez przyrostka wersji. W `~/.local/bin` stoją dowiązania do plików
  `*-21`. Bez `PATH=$HOME/.local/bin:$PATH` budowa pada.
- Powłoka główna składa się celem `gnu` i wymaga `WebView2Loader.dll` obok
  siebie; niesie ją instalka NSIS. Instalator kreatora składa się celem `msvc`
  i jest samodzielnym plikiem.
- Playwright chodzi wyłącznie z `executablePath: '/usr/bin/chromium-browser'`
  i `--no-sandbox`.
- Katalog tymczasowy sesji znika wraz z sesją. Co ma przetrwać, ląduje
  w `~/robocze/` albo w repozytorium.
