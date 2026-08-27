# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### zaplecze-modeli

**Szesnaście gigabajtów wag stoi na maszynie i nie ma czym się uruchomić.**
Zmierzone: nie ma `torch`, `transformers`, `sentence-transformers`, `fastembed`,
`kokoro` ani `segment_anything`. Stoją wyłącznie `onnxruntime` i `ctranslate2`.
Skutek najcięższy: **własny silnik wiedzy rdzenia nie działa** — `knowledge.index`
woła `internal/wiedza/pomocnik_osadzen.py`, a ten żąda `fastembed`.

| | |
|---|---|
| **Gałąź** | `teren/zaplecze-modeli` z `main` |
| **Wykaz plików** | `budowa/server/internal/wiedza/`, `budowa/server/internal/core/adapter_modul_mowa*.go`, `adapter_narzedzia_obraz_model_silniki.go`, `zaleznosci_zewnetrzne.go`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/shared/` (kontrakt zmienia wyłącznie teren `pomiar-stron`), `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, adaptery dokumentów, obrazu i przeglądarki |

| Model | Waga | Format | Komenda, która po niego sięgnie |
|---|---|---|---|
| embedder (XLMRoberta, wymiar 1024) | 4,3 GB | `pytorch_model.bin` + ONNX | `knowledge.index`, `knowledge.search` |
| reranker | 2,2 GB | `safetensors` | brak komendy — **zgłoś, nie dokładaj** |
| CLIP | 1,6 GB | `safetensors` | brak komendy — zgłoś |
| twarze | 692 MB | `.pth` | brak komendy — zgłoś |
| Kokoro | 340 MB | 55 × `.bin`/`.pth`, **54 głosy, ani jednego polskiego** | `speech.*` — rozstrzygnij, czy wart deklaracji |
| ESRGAN | 128 MB | `.pth` | `image.upscale` — **program `realesrgan-ncnn-vulkan` już stoi i rdzeń go zna** |
| MobileSAM | 39 MB | `.pth` | `design.photo.select.object` **jest już zrobione w Go**, a zapora Design zabrania wołania procesu w tej rodzinie |

**Instalowanie jest tu dozwolone — wyjątkowo i wyłącznie w tym terenie.**
Właściciel polecił, żeby aplikacja miała wszystkie narzędzia czynne. Zaplecze
modeli stoi w hybrydzie **na serwerze wdrożenia**, nie u Operatora, więc
instalacja na maszynie budowlanej jest instalacją tego serwera. Każdą pozycję,
którą postawisz, wypisujesz w raporcie wraz z wagą na dysku.

**Kontraktu nie zmieniasz.** Model bez komendy w kontrakcie wraca zgłoszeniem
wraz z propozycją obszaru — nie dokładasz komend, bo kontrakt należy w tej turze
do innego terenu.

**Kryteria odbioru.**

1. `knowledge.index` i `knowledge.search` **działają** — wykazane uruchomieniem
   na prawdziwej treści, z przytoczonym żądaniem i odpowiedzią rdzenia.
2. Dla każdego z siedmiu modeli: albo działa i jest wykazany uruchomieniem, albo
   ma podany powód, dla którego dziś nie może, wraz z tym, czego brakuje.
3. Wagi, które stoją w `/opt/danaco-modele`, są **użyte** albo jest wprost
   napisane, dlaczego rdzeń sięga po inne — pobranie drugiej kopii tego samego
   modelu jest uchybieniem, chyba że podasz powód.
4. Brak biblioteki albo wag daje odmowę **nazywającą brak i drogę naprawy**, nie
   błąd wewnętrzny — wykazane sprawdzianem.
5. `gotestsum -- -count=1 ./...` — zero niepowodzeń, wobec stanu zastanego
   2041 zdanych, 17 pominiętych, zero niezdanych.
6. Kontrakt nietknięty — wykazane sumą kontrolną.
7. Wykaz wszystkiego, co postawiłeś na maszynie, wraz z wagą — w raporcie.

### pomiar-stron

Cztery programy mierzące stronę i punkt końcowy. **Ten teren jako jedyny zmienia
kontrakt** — obszary istnieją i już mierzą stronę trzema sondami, ale osi
wydajności, dostępności i obciążenia nie mają.

| | |
|---|---|
| **Gałąź** | `teren/pomiar-stron` z `main` |
| **Wykaz plików** | `budowa/shared/contract.json`, `budowa/server/internal/core/adapter_modul_przegladarka_*.go`, `adapter_modul_apps_*.go`, `adapter_modul_developer_api*.go`, `zaleznosci_zewnetrzne.go`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, adaptery dokumentów i obrazu |

| Narzędzie | Obszar | Czego brakuje w kontrakcie |
|---|---|---|
| **pa11y** | `browser` (47 komend) | audyt WCAG na otwartej karcie, z wykazem naruszeń i wskazaniem węzła DOM. `design.color.accessibility.audit` bada **paletę**, nie stronę |
| **Lighthouse** | `apps` (41 komend) | audyt wydajności zwracający Core Web Vitals. `apps.deployment.health.get` oddaje dostępność, nie pomiar |
| **k6** albo **autocannon** | `developer` (50 komend) | przebieg obciążeniowy wraz z kształtem wyniku — percentyle, przepustowość. `developer.api.request` strzela **jednym** żądaniem |

**Zmiana kontraktu jest tu dozwolona i obwarowana.** Kontrakt nie był tknięty od
przejęcia — suma `2cbb843d33f4531b05cd` stoi od pierwszego dnia. Wolno Ci
**dołożyć** komendy, struktury i wyliczenia. **Nie wolno** zmienić ani usunąć
niczego istniejącego: żadnej komendy, żadnego pola, żadnej wartości wyliczenia.
Generator z `budowa/shared/gen` wytwarza z kontraktu **oba** artefakty —
`contract.go` i `contract.ts` — i musi po Twojej zmianie dawać wynik bajtowo
powtarzalny w dwóch przebiegach.

**Wykaz zależności jest wspólny z dwoma innymi terenami biegnącymi teraz.**
Deklarację narzędzia zakładasz **przy miejscu użycia**, tak jak robi to rdzeń
(Pandoc przy dokumentach, ffmpeg przy nagraniach), a do
`zaleznosci_zewnetrzne.go` dopisujesz wyłącznie odwołanie. Przy scaleniu
rozjazd w tym jednym pliku rozstrzyga Prowadzący — nie jest to Twoja usterka.

### sprawdziany-drogi-wejscia

Trzy sprawdziany w `skutek_wydania_studia_test.go` rozjechały się z rdzeniem
od chwili przejęcia — pomiar wcześniejszy wykazał, że opisują zamiar porzucony,
którego nikt za zmianą nie poprawił. Osobno: droga bez poczty nie ma ani
jednego sprawdzianu własnego — zachowanie rozstrzygnięte pozycją 11 stoi
wyłącznie na komentarzu i pomiarze jednorazowym; pierwsza zmiana w bramce
zniesie je bez niczyjej wiedzy.

| Sprawdzian | Czego żąda | Co rdzeń robi |
|---|---|---|
| `TestRejestracjaBezKontaNadawczegoOdmawiaINieZakladaKonta` | odmowy i zera wierszy przy braku nadajnika | zakłada konto, stawia znacznik, wpuszcza hasłem — pozycja 11 rozstrzyga na rzecz rdzenia |
| `TestNieudaneNadanieListuCofaRejestracje` | cofnięcia rejestracji przy nadajniku nieosiągalnym | do zmierzenia w terenie |
| `TestSkanowanieZUrzadzeniaOdmawiaNazwanie` | odmowy nazywającej brak | do zmierzenia w terenie |

| | |
|---|---|
| **Gałąź** | `teren/sprawdziany-drogi-wejscia` z `main` |
| **Wykaz plików** | `budowa/server/internal/core/skutek_wydania_studia_test.go`; nowy plik sprawdzianów drogi bez poczty w `internal/core/`; wyłącznie przy wykazanej pomiarem usterce rdzenia — plik bramy albo uwierzytelnienia niosący zachowanie |
| **Poza terenem** | kontrakt, `zaleznosci_zewnetrzne.go`, `urzadzenia_skaner.go`, adaptery dokumentów, obrazu, przeglądarki, aplikacji, developer-api, mowy i `wiedza/` (pliki terenów biegnących), `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` |
| **Wykonawca** | sesja wysłana przez Prowadzącego 27.08.2026 |

**Kryteria odbioru.**

1. Dla każdego z trzech sprawdzianów rozstrzygnięte pomiarem, czy zawodzi
   sprawdzian, czy rdzeń — z przytoczonym pomiarem — i poprawiona ta strona,
   która się myli, a nie ta, którą łatwiej. Dla pierwszego rozstrzyga
   pozycja 11: rdzeń ma rację.
2. Droga bez poczty ma sprawdziany własne: rejestracja, wejście hasłem,
   potwierdzenie adresu po ustawieniu nadajnika, zdjęcie znacznika — zgodne
   z pozycją 11 rejestru decyzji.
3. `gotestsum -- -count=1 ./...` — zero niepowodzeń; liczba pominiętych nie
   rośnie, a jeśli któryś z trzech sprawdzianów stoi dziś wśród pominiętych,
   po terenie pominięć ubywa.
4. Kontrakt nietknięty.
5. Rewizje obejmują wyłącznie pliki terenu.

### odmowy-skanera

Brak urządzenia jest już nazwany, ale brak programu `scanimage` na Linuksie
wychodzi odmową arsenału bez wskazania drogi obejścia, podczas gdy
`bladWarstwyWia` (`urzadzenia_skaner.go`) dla tej samej sytuacji na Windowsie
podaje `studio.ingest.queue.add`. Ta sama asymetria w `wykazSkanerow`: gałąź
Windows przekłada odmowę, gałąź Linux oddaje ją surową.

| | |
|---|---|
| **Gałąź** | `teren/odmowy-skanera` z `main` |
| **Wykaz plików** | `budowa/server/internal/core/urzadzenia_skaner.go` oraz sprawdzian skanera (istniejący albo nowy plik sprawdzianu tego zakresu) |
| **Poza terenem** | kontrakt, pozostałe pliki `internal/core/` — w szczególności adaptery czterech terenów biegnących równolegle — `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` |
| **Wykonawca** | sesja wysłana przez Prowadzącego 27.08.2026 |

**Kryteria odbioru.**

1. Brak `scanimage` na Linuksie daje odmowę wskazującą
   `studio.ingest.queue.add` — parytet z `bladWarstwyWia` — wykazane
   sprawdzianem.
2. `wykazSkanerow`: gałąź Linux przekłada odmowę tak samo jak gałąź Windows —
   wykazane sprawdzianem.
3. Komunikaty wzorowane na istniejącej gałęzi Windows — zero nowych nazw,
   kodów i oznaczeń.
4. `gotestsum -- -count=1 ./...` — zero niepowodzeń wobec stanu zastanego
   2041 zdanych, 17 pominiętych, zero niezdanych.
5. Kontrakt nietknięty.

## Zgłoszenia oczekujące na teren

Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

### Trzy narzędzia treści pisanej bez legalnego miejsca wpięcia

Ustalenia terenu `dokumenty-i-tekst`, każde zmierzone, żadne nienaprawialne
w granicach tamtego terenu.

| Rzecz | Zmierzone | Skutek |
|---|---|---|
| LanguageTool przy `translate.quality.check` | tryb `--bitext` (LT 6.6) na parze z jaskrawą rozbieżnością liczb (10:30→9:30, 250→500 EUR) oddaje **zero ustaleń**, a wyjścia maszynowego dla bitext program nie ma wcale | wpięcie dałoby instrument milczący przy realnej usterce; komenda mierzy wierność panelu wobec źródła i dziś nie ma czym |
| unpaper przy `research.source.ocr` | uchwyt `adapter_modul_badania_lektura.go:589` upuszcza pola `preprocess` i `languages`, a `DocumentTextExtractRequest` kontraktu tych pól nie ma | obróbka wstępna nieosiągalna z drogi badań; pole `preprocess` pozostaje martwe |
| Tika przy `library.metadata.get` | obsługa stoi w `adapter_modul_library_technika.go` — pliku terenu `obraz-i-diagramy`, gdzie exiftool wszedł w tę samą komendę | rozłączność terenów; do rozważenia, czy Tika ma tam co dołożyć wobec exiftoola |

### mermaid-cli nie ma legalnego miejsca wpięcia

Ustalenie terenu `obraz-i-diagramy`, zmierzone uruchomieniem. Obie komendy
z tabeli narzędzi są zamknięte: `design.diagram.render`
(`adapter_modul_design_wykresy.go:790`) leży za zaporą fotografii, która
zabrania **każdego** wołania procesu w plikach `adapter_modul_design*.go` —
plik próbny z samą wzmianką `zewnetrzne.Wolaj` wywrócił zaporę, choć zapora
nazw przepuszcza `mmdc`; nadto kontrakt tej komendy przyjmuje węzły i krawędzie,
nie ma pola na źródło Mermaid. `apps.architecture.export`
(`adapter_modul_aplikacje_architektura.go:427`) niesie w nagłówku
rozstrzygnięcie „żaden format nie woła programu z zewnątrz, więc eksport działa
na instalce, która niesie sam rdzeń".

Samo narzędzie działa: z przeglądarką stojącą w `/opt/ms-playwright` renderuje
źródło w kształcie eksportu rdzenia (SVG 12 539 B, PNG 265×278). Bez wskazania
przeglądarki `mmdc` odmawia — puppeteer żąda własnej kopii.

Do rozstrzygnięcia: czy rysowanie diagramu ma powstać poza obszarem Design
(nowa komenda i nowe miejsce), czy zostaje niezrobione. Rozstrzygnięcie dotyka
zapory niosącej rozstrzygnięcie Właściciela, więc nie jest samą robotą.

### Komunikat rdzenia radzi budowę usuniętymi skryptami

`budowa/server/internal/narzedzia/wpiecie.go:52` i `:111` — komentarz i treść
komunikatu błędu wskazują `scripts/wydanie.sh` i `scripts/pakowanie.sh`,
usunięte z repozytorium jako materiał instalek natywnych (pozycja 8 pkt 4).
Odmowa radzi Operatorowi drogę naprawy, której nie ma. Drobny teren rdzenia.

### Skrypty hybrydowe wskazują usunięte profile i nieistniejący katalog

Wszystkie trzy `budowa/scripts/instalka-hybryda-*.sh` wskazują profile
`tauri.hybryda-*.conf.json`, których nie ma (w `budowa/desktop/src-tauri/`
stoi wyłącznie `tauri.conf.json`), oraz `budowa/client/` zamiast
`budowa/klient/`; to samo `pakiet-serwera.sh` i `pokaz.sh`. Hybryda jest
jedyną postacią produktu, więc skrypty są żywe, a odmawiają na starcie.
Naturalny moment terenu: powstanie prawdziwej budowy `budowa/klient/dist`.
Osobno do rozstrzygnięcia: `instalka-hybryda-linux.sh` wobec pozycji 8
(jedyna platforma Windows 11), gdy `droga.rs` wciąż niesie tor AppImage.

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

### Droga SANE bez odpowiednika `bladWarstwyWia` — gotowe do otwarcia

Brak urządzenia jest już nazwany. Brak samego programu `scanimage` dalej wychodzi
odmową arsenału bez wskazania drogi obejścia, podczas gdy `bladWarstwyWia`
(`urzadzenia_skaner.go`) dla tej samej sytuacji na Windowsie podaje
`studio.ingest.queue.add`. Ta sama asymetria dotyczy `wykazSkanerow`: gałąź
Windows przekłada odmowę, gałąź Linux oddaje ją surową. Operator na Linuksie bez
`sane-utils` nie dowie się, że materiał da się wnieść inną drogą.

### Sprawdzian katalogu akcji szuka nieistniejącego katalogu — naprawione

`katalog_akcji_test.go` szukał `../../../client/src/ikony/zrodla`, którego nie
ma w żadnej gałęzi. Rewizja `9ba1f90` wiąże sprawdzian ze źródłami ikon nowego
klienta — `sciezkaZrodelIkon` wskazuje `../../../klient/src/ikony/zrodla`.
Przechodzi w biegu odniesienia 2041 zdanych, zero niezdanych.

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
| `fundament-klienta` | `teren/fundament-klienta` | `d19bfeb` warstwa połączenia i protokołu | weryfikacja Prowadzącego pomiarem: kompilacja bez błędu, 17 sprawdzianów zdanych, rozmowa z żywym rdzeniem, generat bajtowo powtarzalny, zero dotknięć DOM, kontrakt nietknięty |
| `naprawy-rdzenia` | `teren/naprawy-rdzenia` | `060d5b7` naprawy i brama kontraktu | weryfikacja Prowadzącego pomiarem: 2022 zdane wobec 2003 zastanych, te same 4 niezdane, kontrakt nietknięty |
| `proba-prototypow` | `teren/prototypy` | `66e5ee0` przepływ wejścia · `61a5867` moduł Studio · `53bc3d5` odsyłacze | kontrola sesji nadzorującej wykonanie, weryfikacja Prowadzącego pomiarem |
| `obraz-i-diagramy` | `teren/obraz-i-diagramy` | `f7fab06` dogniecenie zapisu i metadane osadzone | kontrola osobnej sesji własnym biegiem: osiem zapór zdanych nietkniętych, 33 sprawdziany zdane w realnych czasach z przytoczonymi rozmiarami przed i po (PNG 3322→2456 B, JPEG 45399→39623 B, model barw paleta dowodzi wejścia pngquanta); sprawdziany maszyny bez programu wytwarzają ją naprawdę (`t.Setenv` na pusty katalog); kontrakt nietknięty; zero śladu mmdc w rdzeniu. Scalone `5b98585` |
| `dokumenty-i-tekst` | `teren/dokumenty-i-tekst` | `ea05603` sześć programów treści pisanej · `aa695e6` poprawka nazwy po zwrocie | kontrola osobnej sesji: teren zwrócony za martwą nazwę `zasiegSyntezy` w komentarzu, po poprawce przyjęty; własny bieg kontrolera 14 zdanych, zero pominiętych; sprawdziany mierzą skutek — PDF czytany drugą komendą i innym programem, korekta aż po treść panelu w bazie, OCR dwoma przebiegami, brak programu wywołany `t.Setenv`; wybory hunspell i typst zweryfikowane uruchomieniem; kontrakt nietknięty. Scalone `43d05f6` |
| `aktualizacja-powloki-i-skrypty` | `teren/aktualizacja-powloki-i-skrypty` | `d1b6207` probne.rs · `54cb287` kanał pobrań · `880b41f` skrypty natywne | kontrola osobnej sesji własnym biegiem: `cargo test` 23 zdane, zero niezdanych (stan zastany: nie kompilował się); sha256 sześciu kopii w materiale zamkniętym tożsame; `ADRES_KANALU` zgodny znak w znak z `kanal.adres` wykazu wydań; scalone `adb2a85`, drzewo scalone tożsame z kontrolowanym. Uwaga trwała: kompilacja powłoki wymaga `budowa/klient/dist` (w `.gitignore`) i zmiennych CARGO/RUSTUP ze środowiska maszyny |

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
