# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### okno-do-uruchomienia

Klient ma dziś 6382 wiersze TypeScriptu i **nie da się ich ani obejrzeć, ani
spakować**: nie ma dokumentu, nie ma budowania, nie ma `dist/`. Powłoka Tauri
czeka na `budowa/klient/dist`, którego nikt nie wytwarza. Teren zamyka tę lukę
i daje pierwszą rzecz, którą Właściciel może kliknąć.

| | |
|---|---|
| **Gałąź** | `teren/okno-do-uruchomienia` z `main` |
| **Wykaz plików** | `budowa/klient/` — dokument, nastawa budowania, `package.json`, `tsconfig.json` |
| **Do czytania, bez zapisu** | `design/zasoby/` (arkusze i żetony), `design/05-okna/przeplyw/przeplyw-wejscia.html`, `budowa/desktop/src-tauri/tauri.conf.json` |
| **Poza terenem** | `budowa/server/`, `budowa/desktop/`, `budowa/shared/`, `design/`, `prowadzenie/` |

**Przedmiot.** Doprowadzić okno drogi wejścia do postaci, która **uruchamia się
i daje się przeklikać** — a wynik budowania trafia tam, gdzie powłoka go szuka.

Arkusze, których okno potrzebuje, wymienia prototyp w swoim nagłówku:
`zetony/fonty.css`, `zetony/zetony.css`, `css/fundament.css`,
`css/komponenty.css`, `wejscie.css`. **Kolejność wpięcia jest wiążąca** —
`prototyp.css` nie wchodzi, bo należy do warstwy podglądu, nie produktu.

**Arkuszy nie powielasz.** Warstwa projektowa ma jedno źródło w `design/zasoby/`
i pozostaje poza tym terenem. Sposób ich wciągnięcia do pakietu jest pracą
inżynierską — kopia przy budowaniu, dowiązanie albo import — ale skutkiem ma być
**jedno źródło, nie dwa**. Powielony arkusz rozjedzie się przy pierwszej zmianie
żetonu i jest uchybieniem odbioru.

**Klient nie ma dziś ani jednej zależności produkcyjnej.** Jeśli budowanie ich
wymaga, wchodzą wyłącznie jako narzędzie budowania, nie do pakietu — a wybór
uzasadniasz w raporcie. Na maszynie stoją `vite`, `bun`, `pnpm` i `tsc`.

**Kryteria odbioru.**

1. Polecenie budowania wytwarza `budowa/klient/dist` — z przytoczonym wynikiem
   uruchomienia i wykazem wytworzonych plików wraz z rozmiarami.
2. Dwa przebiegi budowania dają ten sam wynik — wykazane sumą kontrolną.
3. Okno **wyświetla się w przeglądarce**: zrzut ekranu odsłony łączenia,
   rejestracji i logowania. Zrzuty odkładasz poza repozytorium.
4. **Pełne przejście wobec żywego rdzenia**: połączenie, założenie konta,
   logowanie, wejście do środowiska — z przytoczonym przebiegiem, wykonane
   w przeglądarce, nie w sprawdzianie warstwy.
5. Zero błędów konsoli i zero nieudanych żądań przy tym przejściu — wykazane
   odczytem konsoli, przy sondzie dodatniej dowodzącej, że odczyt działa.
6. Arkusze pochodzą z `design/zasoby/` i **nie są powielone** — wykazane
   porównaniem sum kontrolnych źródła i tego, co w pakiecie.
7. `tsc --noEmit` bez błędu; sprawdziany zastane (`26 z 26`, `7 z 7`, warstwy
   fundamentu) dalej przechodzą.
8. Zmiany wyłącznie w `budowa/klient/` — wykazane `git show --name-only`.
9. Rzecz wymagająca rozstrzygnięcia wraca zgłoszeniem wraz z przyjętym
   rozstrzygnięciem — nie wstrzymuje reszty.

### warsztat-kodu

Osiem narzędzi stoi na maszynie i jest z rdzenia nieosiągalnych, choć **komendy,
które po nie sięgną, już istnieją w kontrakcie**. Rdzeń w dwóch miejscach nazywa
ten brak wprost we własnych komentarzach. Nic z tego nie jest oknem, więc teren
nie zależy od prototypów.

| | |
|---|---|
| **Gałąź** | `teren/warsztat-kodu` z `main` |
| **Wykaz plików** | `budowa/server/internal/core/zaleznosci_zewnetrzne.go`, adaptery modułów Developer i Terminal w `budowa/server/internal/core/`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, adaptery pozostałych modułów |

**Przedmiot — osiem narzędzi pod istniejące komendy.**

| Narzędzie | Komenda, która po nie sięgnie | Co rdzeń mówi dziś |
|---|---|---|
| `ruff` | `terminal.script.lint`, `developer.lint.get` | mapa analizatorów zna bash i PowerShell, **Pythona nie ma** mimo zadeklarowanego `narzedziePython` |
| `semgrep` | `developer.scan.run` | rdzeń pisze o własnym skanie: „To jest zakres węższy niż `semgrep`" |
| `ast-grep` (`sg`) | `developer.grep.search`, `developer.refactor.apply` | wyszukanie po składni, nie po napisie |
| `jscpd` | `developer.scan.run` | duplikaty w TS/JS |
| `dupl` | `developer.scan.run` | duplikaty w Go |
| `typos` | `developer.lint.get` | literówki w identyfikatorach |
| `stylelint` | `developer.lint.get` | CSS — wzorzec zadeklarowanych Prettiera i ESLinta |
| `typescript-language-server` | `developer.symbol.navigate`, `developer.refactor.apply` | wzorzec `gopls` przeniesiony na TypeScript |

**Wzorzec jest gotowy i masz go powtórzyć, nie wymyślać.** Rdzeń ma jedno miejsce
wołania procesów zewnętrznych (`internal/zewnetrzne/wolanie.go`) i jeden wykaz
zależności zasilający sondę startową. Każde z ośmiu narzędzi wchodzi tak samo jak
trzydzieści już zadeklarowanych.

**Zapory, których nie wolno naruszyć.** Rdzeń niesie sprawdziany zabraniające
powrotu pewnych programów: `zapora_warsztatu_pdf_test.go` (qpdf, Ghostscript)
oraz `zapora_fotografii_test.go` (nazwy silników w obszarze Design). Żadne
z ośmiu narzędzi tego terenu ich nie dotyczy — ale sprawdziany mają dalej
przechodzić.

**Kryteria odbioru.**

1. Każde z ośmiu narzędzi zadeklarowane w wykazie zależności wraz z zakresem
   („co przestaje działać przy braku") — sonda startowa je widzi, wykazane
   przytoczonym wynikiem uruchomienia rdzenia.
2. Każde podłączone do **istniejącej** komendy kontraktu — bez nowych komend
   i bez zmiany kontraktu. Suma `contract.json` nietknięta.
3. Każda komenda wykazana uruchomieniem na prawdziwym pliku: przytoczone
   wywołanie i przytoczona odpowiedź rdzenia.
4. Brak narzędzia daje odmowę **nazywającą brak i drogę naprawy**, nie błąd
   wewnętrzny — wykazane sprawdzianem przy narzędziu niedostępnym.
5. Żaden proces zewnętrzny nie wołany poza `internal/zewnetrzne` — wykazane
   przeszukaniem na `exec.Command`.
6. `gotestsum -- -count=1 ./...` — niepowodzenia nie rosną wobec stanu zastanego;
   obie zapory dalej przechodzą.
7. Narzędzie, którego nie da się podłączyć bez rozstrzygnięcia, wraca
   zgłoszeniem wraz z przyjętym rozstrzygnięciem — nie wstrzymuje reszty.

## Zgłoszenia oczekujące na teren

Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

### Kontrast metadanych w przedsionkach — naprawione

`design/zasoby/przedsionek.css` używał `--dn-tekst-3` (#787C85, szary-500) dla
metadanych 12 px w ośmiu klasach (`pd-kafel-opis`, `pd-sesja-meta`, `pd-kafel-meta`,
`pd-nadtytul`, `pd-motto`, `pd-strefa-opis`, `pd-listwa-meta`, `pd-filtr`) — łamiąc
regułę żetonu „wyłącznie ≥18,66 px" i dając ~28 węzłów poniżej progu WCAG na okno.
Wszystkie osiem przeniesiono na `--dn-tekst-2` (dark szary-400, light szary-600).
Zweryfikowane: 4 przedsionki × 2 motywy = 0 naruszeń axe. Przy okazji zdjęto
zastany `aria-selected="true"` z przycisku `.pd-sesja` (nieprawidłowy ARIA na
`<button>`) na `aria-current="true"` we wszystkich czterech oknach, z hakiem CSS
`.pd-sesja[aria-current='true']`. Decyzja Właściciela: naprawić teraz.

Pozostałe pliki z wpisu „Kontrast warstwy wspólnej" (`rama.css`, `stanowisko.css`,
`css/komponenty.css`, `panel-sesji.css`) nie były tu ruszane — to osobny zakres.

### Żeton mikro-odstępu `--dn-od-05` — wprowadzony przy domknięciu centrum

`zetony.css` dodaje `--dn-od-05: 2px` — półstopień siatki 4 px, między `--dn-od-0`
i `--dn-od-1`. Powód: kryterium maszynowe etapu 1 wymaga, by odstępy pochodziły
z żetonów, a najciaśniejszy realny odstęp (etykieta↔opis w kaflu) wynosił 2 px
bez pokrycia w skali. Zdjęto surowe `gap: 2px` z `okna/centrum-dowodzenia.css`
(2×) i `okna/studio.css` (4×). Decyzja delegowana przez Właściciela: „wariant
zgodny z profesjonalnym standardem" — skala odstępów zawiera używane wartości
jako żetony (wzór: Tailwind `0.5`, Material). Zmiana warstwy wspólnej odnotowana
tutaj zgodnie z ograniczeniem warstwy wspólnej z planu etapów.

### Droga wejścia — braki po stronie rdzenia, ujawnione pomiarem

Wszystkie pochodzą z terenu `droga-wejscia` i są zmierzone na żywym rdzeniu.
Żaden nie wstrzymał pracy; okno podłącza je dziś do odmowy nazywającej brak.

| Rzecz | Zmierzone | Skutek |
|---|---|---|
| brak drogi powtórnego potwierdzenia adresu | `auth.verify` z cudzą drogą → `not_authenticated`; powtórny `auth.register` → `conflict`; żadna z dziewięciu komend `auth.*` nie wydaje drogi drugi raz | Operator, do którego list nie dotarł, **nie ma wyjścia z okna** — dotyczy także instalacji **z pocztą**, nie tylko bez |
| trzy z pięciu etapów przygotowania środowiska bez komendy | prototyp wymienia pięć; droga wejścia obsługuje dwa. „Profil i uprawnienia", „Kanały modeli", „Magistrala kontekstu" nie mają przypisanej komendy | pasek postępu zatrzymuje się na 40% — wartość prawdziwa, nie ozdobna |
| odmowa logowania nie niesie długości zwłoki | osiem prób: 100, 350, 600, 1100, 2100, 4100, 5100, 5100 ms; `details` puste za każdym razem | okno mierzy zwłokę samo, czasem trwania próby poprzedniej |
| magazyn tokenu bramki między uruchomieniami | `connection.hello` przyjmuje token i oddaje `authenticated: true`; gdzie token mieszka, nie rozstrzyga ani prototyp, ani kontrakt, ani rejestr | odsłona „rozpoznano zaufane urządzenie" działa dopiero po wskazaniu magazynu — należy do powłoki |

### Adres serwera wdrożenia nie stoi w żadnym źródle

Teren `powloka-tauri` przeszukał `budowa/`, `design/`, `docs/` i `prowadzenie/`:
pozycja 8 nazwy serwera nie podaje, prototyp instalatora nie ma kroku adresu,
`wydania.json` niesie wyłącznie kanał pobrań. Powłoka bierze adres ze wskazania
Operatora albo ze zmiennej `DANACO_HOST_RDZENIA`. Czy adres ma być wkompilowany
przy składaniu instalki — i skąd wtedy pochodzi — jest rozstrzygnięciem
Właściciela.

### Zastane usterki powłoki, ujawnione przy przejęciu

| Rzecz | Stan |
|---|---|
| `aktualizacja/probne.rs` nie istnieje | `droga.rs:293` importuje `KatalogProbny` w siedmiu sprawdzianach; `cargo test` **nie kompiluje się od chwili przejęcia** |
| `wykonaj_aktualizacje` pobiera spod dowolnego adresu | `pobranie.rs` sprawdza wyłącznie sumę SHA-256, też podaną przez stronę; brak przywiązania do kanału z `wydania.json` |
| cztery z sześciu poleceń IPC bez odbiorcy | `stan_rdzenia`, `wskazanie_rdzenia`, `wybierz_katalog_roboczy`, `wykonaj_aktualizacje` — opisywały odbiorców w usuniętym `client/src/powloka/` |
| `app.security.csp: null` | strona bez polityki treści; ułożenie wymaga interfejsu, który dopiero powstaje |
| sześć skryptów w `budowa/scripts/` | wskazują usunięte profile nastaw, nieistniejący `budowa/client/` i wariant natywny zniesiony pozycją 8; pozycja 8 pkt 4 mówi, że mają leżeć **poza** repozytorium |

### Łańcuchy widoczne poza katalogiem treści w paczce designu

Właściwość, którą pozycja 12 nazywa wiążącą, jest w dwóch miejscach naruszona:
`design/zasoby/okna/wejscie/skladniki/pole-hasla.js:55` niesie zaszyte
`'Pokaż hasło'`, a `design/zasoby/okna/przeplyw-wejscia.js` powiela w kodzie
`SILA_PUSTE`, `SILA_OPISY` oraz dwa napisy o schowku, które stoją już
w katalogu. W kliencie wszystkie te napisy siedzą w katalogu treści.

### Kontrast warstwy wspólnej — gotowe do otwarcia

Żeton `--dn-tekst-3` (`design/zasoby/zetony/zetony.css` w. 222) niesie własną
regułę: wyłącznie metadane i tekst od 18,66 px półgrubego. Warstwa wspólna łamie
ją w około siedemdziesięciu miejscach: `rama.css` 24, `css/komponenty.css` 16,
`okna/centrum-dowodzenia.css` 14, `stanowisko.css` 8, `prototyp.css` 4,
`panel-sesji.css` 4. Żeton jest poprawny — wadliwe jest jego użycie. Dopóki to
stoi, każde okno korzystające z ramy niesie naruszenia wagi `serious` i żaden
teren nie domknie kryterium dostępności bez wyjątku.

| Plik i wiersz | Rzecz | Zmierzony kontrast |
|---|---|---|
| `rama.css:981` | `.dn-stan` — pasek stanu, 6 pozycji | 3,85 ciemny · 3,99 jasny |
| `stanowisko.css:148` | `.sta-kom-pole span`, zaszyte 10 px | 3,85 · 3,99 |
| `stanowisko.css:288` | `.sta-wpis-godzina`, zaszyte 10 px | 4,25 · 4,17 |
| `stanowisko.css:123` | `.sta-okno-znacznik` | 3,85 · 3,99 |
| `css/komponenty.css:226` | `.dn-pole-opis` przy 13 px | 4,17 · 4,25 |
| `panel-sesji.css:59, 139` · `okna/centrum-dowodzenia.css:258, 288` | tekst 12 px | poniżej progu |

Osobno, ta sama warstwa: `zetony.css:69` deklaruje przy `--dn-rama-tekst-3`
kontrast 4,74 : 1 na ramie. Zmierzone: 4,45 : 1, czyli poniżej progu 4,5.
Deklaracja w komentarzu jest nieprawdziwa.

### Mechanizmy warstwy prototypu — gotowe do otwarcia

| Plik i wiersz | Usterka | Waga |
|---|---|---|
| `prototyp.js:66–83` | `przelaczWidok` nadaje `aria-selected` elementom `<button>` i `<a>`, którym atrybut nie przysługuje | krytyczna |
| `prototyp.js:88–102` | wędrujący `tabindex` ustawiany tylko przy starcie, nieodświeżany po przełączeniu | poważna |
| `prototyp.js:300–341` | wstrzykiwany pasek prototypu stoi poza punktami orientacyjnymi | umiarkowana |
| `powloka.js` | montuje pełną powłokę bezwarunkowo; brak trybu „rama dopiero po uwierzytelnieniu" | poważna |
| `css/komponenty.css:1014` | `.dn-postep-wartosc` bez `display: block` — pasek postępu renderuje pusty tor | poważna |
| `css/komponenty.css:1481–1499` | `.dn-alert` bez gniazda ikony, bez części tytułu i treści, bez wariantów błędu i powodzenia | poważna |
| `rama.css:102–114` | `.dn-narzedzia-pas` na barwie gruntu roboczego — przyczyna źródłowa braku rozgraniczenia wstążek | poważna |
| `rama.css` — szyna | strefa środowisk z 38 pozycjami wypycha strefę szybkiego wyboru poza kadr | poważna |

### Wstążka narzędziowa do wyniesienia do warstwy wspólnej

Szkielet wstążki okna roboczego istnieje dziś **wyłącznie w arkuszu modułu
Studio**. Kontrakt wstążki obowiązuje wszystkie moduły — pozycja 9 rejestru
decyzji — więc przy drugim module rozjedzie się bez niczyjej złej woli.

Do wyniesienia: forma paska wraz z trzema strefami, stała wysokość, zwijanie
prawej grupy do menu nadmiaru, trwały stan wybrania narzędzia, nieruchomość przy
przewijaniu treści. Zmienna pozostaje wyłącznie zawartość stref, właściwa
rodzajowi karty.

### Komponent kart okna — wstrzymane

Wstrzymane do rozstrzygnięcia Właściciela w sprawie nazw dwóch pięter kart.
Otwarcie terenu na drugi moduł przed tym rozstrzygnięciem odtworzy pomieszanie
z biblioteki w każdym kolejnym prototypie.

| Plik i wiersz | Rzecz |
|---|---|
| `karty-okna.css` | arkusz nosi nazwę poprawną, a definiuje kontener `.dn-karty-sesji`, w którym stoją elementy `.dn-karta-widoku` — sprzeczność w jednym pliku |
| `karty-okna.js:27` | mechanizm kart zakotwiczony w `.dn-obszar-panel--glowny` z `rama.css`; w oknie roboczym komponent jest martwy |
| `karty-okna.css:140` | przycisk zamknięcia wewnątrz `div[role="tab"]` — zagnieżdżona interaktywność, przenosi się na każde okno używające komponentu |
| `.dn-izolacja-wskaznik` | stoi w ramie aplikacji, a izolacja jest cechą okna roboczego (`izolacja.css`) — etykietuje niewłaściwy poziom |

### Usterki zastane, ujawnione przy próbce — gotowe do otwarcia

Wszystkie sprzed terenu, żadna nie jest jego skutkiem.

| Miejsce | Rzecz |
|---|---|
| `design/INDEKS.html` | martwy odsyłacz `href="kontrakt systemu projektowego"` — jedyne 404 wśród 93 odsyłaczy strony |
| `design/zasoby/okna/studio.js` oraz `studio.html` w. 786 | zdublowana obsługa `[data-srod]`; klik w przycisk trybu wywołuje dwa komunikaty naraz. Zachowana bez zmiany, bo kryterium wymagało zachowania identycznego |
| plansze i indeks | metryki „N linii · M interakcji" oraz „N w. · M kB" rozjechane ze stanem plików także dla kart nietkniętych — `centrum-dowodzenia.html` podane jako 627 wierszy i 42 kB przy faktycznych 1023 wierszach i 305 kB |

Metryki wymagają rozstrzygnięcia Właściciela: albo są normatywne i dostają
definicję sposobu liczenia, albo znikają. Metody liczenia „interakcji" nie da
się odtworzyć z treści plików, więc dziś nikt nie jest w stanie ich utrzymać.

### Sprawdziany drogi wejścia rozjechane z rdzeniem — gotowe do otwarcia

Trzy sprawdziany zawodzą od chwili przejęcia rdzenia. Pomiar wykazał, że nie są
usterką rdzenia — opisują zamiar porzucony i nikt ich za zmianą nie poprawił.

| Sprawdzian | Czego żąda | Co rdzeń robi |
|---|---|---|
| `TestRejestracjaBezKontaNadawczegoOdmawiaINieZakladaKonta` | odmowy i zera wierszy przy braku nadajnika | zakłada konto, stawia znacznik, wpuszcza hasłem — pozycja 11 rejestru decyzji |
| `TestNieudaneNadanieListuCofaRejestracje` | cofnięcia rejestracji przy nadajniku nieosiągalnym | do zmierzenia w terenie |
| `TestSkanowanieZUrzadzeniaOdmawiaNazwanie` | odmowy nazywającej brak | do zmierzenia w terenie |

Osobno: **droga bez poczty nie ma ani jednego sprawdzianu własnego**. Zachowanie
rozstrzygnięte pozycją 11 stoi dziś wyłącznie na komentarzu i na pomiarze
jednorazowym — pierwsza zmiana w bramce zniesie je bez niczyjej wiedzy.

Teren ma dla każdego z trzech sprawdzianów rozstrzygnąć pomiarem, czy zawodzi
sprawdzian, czy rdzeń, i poprawić tę stronę, która się myli — a nie tę, którą
łatwiej. Do tego założyć sprawdziany drogi bez poczty: rejestracja, wejście
hasłem, potwierdzenie adresu po ustawieniu nadajnika, zdjęcie znacznika.

### Droga SANE bez odpowiednika `bladWarstwyWia` — gotowe do otwarcia

Brak urządzenia jest już nazwany. Brak samego programu `scanimage` dalej wychodzi
odmową arsenału bez wskazania drogi obejścia, podczas gdy `bladWarstwyWia`
(`urzadzenia_skaner.go`) dla tej samej sytuacji na Windowsie podaje
`studio.ingest.queue.add`. Ta sama asymetria dotyczy `wykazSkanerow`: gałąź
Windows przekłada odmowę, gałąź Linux oddaje ją surową. Operator na Linuksie bez
`sane-utils` nie dowie się, że materiał da się wnieść inną drogą.

### Sprawdzian katalogu akcji szuka nieistniejącego katalogu — gotowe do otwarcia

`budowa/server/internal/store/katalog_akcji_test.go:177` szuka
`../../../client/src/ikony/zrodla`. Katalog klienta nazywa się `budowa/klient`,
a `budowa/client` nie istnieje w żadnej gałęzi. To jedyne niepowodzenie
pozostałe w całym module. Do rozstrzygnięcia wraz z pierwszym terenem widoku,
bo dotyczy źródeł ikon nowego klienta.

### Reguła odbioru wyprowadzona z pomiarów

`axe.run()` sam wywołuje dwa błędy 404 (`menu.css`, `ruch.css`), bo rozwiązuje
`@import` względem adresu dokumentu, a nie arkusza. Konsola przed wstrzyknięciem
axe jest pusta. Tych dwóch wpisów nie liczy się jako brudnej konsoli.

Pomiar w przeglądarce wymaga jawnego ustawienia `PLAYWRIGHT_BROWSERS_PATH` na
`/opt/ms-playwright` w poleceniu, a nie polegania na środowisku powłoki — powłoka
uruchomiona przed ustawieniem zmiennej jej nie widzi i pobiera przeglądarki
po raz drugi.

## Tereny zamknięte

| Nazwa | Gałąź | Rewizje | Kontrola |
|---|---|---|---|
| `brama-i-droga-wejscia` | `teren/brama-i-droga-wejscia` | `23a6b84` wyjątek powitania · `e41dd78` straże drogi bez poczty · `9a8ec0c` brak skanera | weryfikacja Prowadzącego pomiarem: bieg wymuszony `-count=1` 537 s — 2029 sprawdzianów, 1 niezdany wobec 4 zastanych; powitanie niepełne odpowiada wersją protokołu na żywym rdzeniu, `channel.add` z brakiem pola dalej odmawia; kontrakt nietknięty |
| `prototypy` | `teren/prototypy` | paczki instalatora i drogi wejścia | przyjęte przez Właściciela; weryfikacja Prowadzącego pomiarem: oba okna wczytują się bez błędu konsoli, zero łańcuchów widocznych poza katalogiem treści |
| `brama-i-droga-wejscia` | `teren/brama-i-droga-wejscia` | `23a6b84` · `e41dd78` · `9a8ec0c` | weryfikacja Prowadzącego pomiarem: 2029 sprawdzianów, 1 niepowodzenie zastane spoza terenu wobec 4 zastanych; wyjątek bramy w jednym miejscu; kontrakt nietknięty |
| `fundament-klienta` | `teren/fundament-klienta` | `d19bfeb` warstwa połączenia i protokołu | weryfikacja Prowadzącego pomiarem: kompilacja bez błędu, 17 sprawdzianów zdanych, rozmowa z żywym rdzeniem, generat bajtowo powtarzalny, zero dotknięć DOM, kontrakt nietknięty |
| `naprawy-rdzenia` | `teren/naprawy-rdzenia` | `060d5b7` naprawy i brama kontraktu | weryfikacja Prowadzącego pomiarem: 2022 zdane wobec 2003 zastanych, te same 4 niezdane, kontrakt nietknięty |
| `proba-prototypow` | `teren/prototypy` | `66e5ee0` przepływ wejścia · `61a5867` moduł Studio · `53bc3d5` odsyłacze | kontrola sesji nadzorującej wykonanie, weryfikacja Prowadzącego pomiarem |

Teren `naprawy-rdzenia` scalony do `main`. Piąta usterka — wartość domyślna
`createVersion` — wróciła jako zgłoszenie, bo kontrakt jej nie ustala. Brama
kontraktu ujawniła, że powitanie kanału musi stać poza nią; rozstrzygnięcie
niesie pozycja 10 rejestru decyzji.

Wynik terenu `proba-prototypow`, przeniesionego na gałąź `teren/prototypy`: oba przedmioty wykonane, wszystkie kryteria spełnione. Zakresy trzech
rewizji rozłączne — sprawdzone. Drzewo czyste. Kryterium 7a zwraca zero trafień
w całym repozytorium, nie tylko w `design/`. Gałąź czeka na ocenę kierunku
przez Właściciela; **nie jest scalona** — próbka rozstrzyga kierunek, a nie
wnosi dorobek.

Usterki wykazu popełnione przez Prowadzącego, obie wychwycone przez sesję
nadzorującą: wpisanie katalogu `przeplyw/` zamiast nazw plików oraz wpisanie
nieistniejącego `design/README.md` (pliki README stoją wyłącznie
w podkatalogach). Wniosek na przyszłość: wykaz plików terenu sprawdza się
odczytem drzewa przed otwarciem, nie z pamięci.

## Wzór wpisu

Otwarcie terenu wymaga wypełnienia wszystkich pól. Pole puste blokuje otwarcie.

- **Nazwa** — rzeczownikowa, opisuje przedmiot pracy, bez oznaczeń literowych
  i numerycznych.
- **Gałąź bazowa** — gałąź, z której teren wyrasta i do której wraca.
- **Wykaz plików** — pełne ścieżki albo katalog. Wykaz nie może przecinać się
  z żadnym terenem otwartym.
- **Kryteria odbioru** — zdania sprawdzalne. Kryterium, którego nie da się
  sprawdzić uruchomieniem albo odczytem pliku, nie jest kryterium.
- **Wykonawca** — oznaczenie sesji prowadzącej teren.
