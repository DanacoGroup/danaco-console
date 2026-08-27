# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### centrum-poprawki

**Prototyp centrum dowodzenia po serii poprawek Właściciela.** Teren prowadzi
wygląd i mechanikę okna `centrum-dowodzenia.html` — separatory pasa kart,
kotwiczenie dymków, wstążki okien, kafle środowisk i modułów, animacje kafli
oraz obrys w barwie znaku. Praca toczy się poprawka po poprawce, na wskazanie
Właściciela ze zrzutu ekranu.

| | |
|---|---|
| **Gałąź** | `teren/centrum-poprawki` z `main` |
| **Wykaz plików** | `design/05-okna/przeplyw/centrum-dowodzenia.html`, `design/zasoby/okna/centrum-dowodzenia.css`, `design/zasoby/okna/centrum-dowodzenia.js`, `design/zasoby/okna/danaco-anim-3d.css`, `design/01-dokumentacja-md/11-uzasadnienia-okien.md` |
| **Poza terenem** | `design/zasoby/zetony/` — bez zgody Właściciela nie rusza się palety ani skali; `design/zasoby/rama.css`, `karty-okna.css`, `panel-sesji.css` i pozostała warstwa wspólna; wszystkie okna poza centrum dowodzenia; `budowa/server/`, `budowa/klient/`, `budowa/desktop/`, `shared/` |

**Warstwa wspólna jest poza terenem, a mimo to została tknięta.** Rewizje
`9e73f37` i `aa8820e` zmieniły `zasoby/rama.css`: uniesienie przycisku wstążki
przeszło z `transform` na `top`, a glif przycisku okna z 12 na 16 px. Pierwsza
zmiana była konieczna — `transform` czyni z przodka blok odniesienia dla
potomków umocowanych do okna widoku, przez co dymek etykiety tracił kotwicę.
Obie obejmują wszystkie okna aplikacji i **wymagają kontroli poza tym terenem**.

**Kryteria odbioru**

- Każda poprawka zamknięta pomiarem w przeglądarce, nie deklaracją; wartość
  przed i po podana w raporcie.
- Zero błędów konsoli i zero odpowiedzi 4xx przy załadowaniu okna.
- Wysokości kafli i kart równe w obrębie strefy; opis kafla bez przypadkowego
  łamania wiersza.
- Dymek każdego wyzwalacza mieści się w widoku i stoi pod swoim przyciskiem.
- Wyłącznie żetony `--dn-*`; barwy, odstępy i rozmiary wpisane wprost wyłącznie
  tam, gdzie Właściciel rozstrzygnął inaczej (barwy animacji kafli środowisk).
- Kontrolę przeprowadza sesja inna niż wykonawcza.

**Nierozstrzygnięte, przeniesione poza teren**

- Punkty łamania siatki środowisk stoją na 1180 i 860 px, poza skalą
  `--dn-bp-*` (640 / 960 / 1280 / 1600). Pochodzą sprzed tego terenu.
- Rozstrzygnięcie o obrysie w barwie znaku żyje jako `--cd-obrys-marki`
  w arkuszu jednego okna; przy rozszerzeniu na pozostałe okna należy do palety.
- Wyściółka `.cd-tresc` przełącza się progiem ekranu (640 px), bo element nie
  może odpytywać własnego pojemnika. Pozostałe punkty łamania idą za płótnem.
- Trzy wiersze dokumentu okna mają 23 252, 20 198 i 10 108 znaków — znacznik bez
  łamania jest nieczytelny w przeglądzie i w różnicy rewizji.
- Stopka lewego okna sięga kanału wydań `pobierz.danaco-group.pl`, zgodnego
  z `wydania.json`, lecz nieopisanego w [rejestrze decyzji](decyzje.md).


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

### dokumenty-i-tekst

Siedem programów do treści pisanej. Wszystkie stoją na maszynie, komendy, które
po nie sięgną, **już są w kontrakcie**. Rdzeń w jednym miejscu nazywa brak
wprost: `adapter_narzedzia_dokument_formaty.go:62` — „brak silnika składu".

| | |
|---|---|
| **Gałąź** | `teren/dokumenty-i-tekst` z `main` |
| **Wykaz plików** | `budowa/server/internal/core/adapter_modul_tlumaczenie_*.go`, `adapter_narzedzia_dokument_*.go`, `adapter_modul_studio_*ingest*.go`, `zaleznosci_zewnetrzne.go`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, adaptery obrazu i przeglądarki |

| Narzędzie | Komenda | Uwaga |
|---|---|---|
| **Apache Tika** `/opt/tika` | `document.text.extract`, `library.metadata.get`, `studio.document.import.file` | Java 25 stoi |
| **LanguageTool** `/opt/languagetool` | `translate.proofread.run`, `translate.quality.check` | dziś korekta idzie wyrażeniami regularnymi |
| **hunspell / enchant-2** | to samo | pisownia; uzasadnij wybór jednego |
| **xelatex** albo **typst** | `document.convert` → PDF | **wybierz jedno i uzasadnij**; TeX to ~1 GB, typst to jeden plik |
| **vale** | `translate.qa.profile.set` | styl prozy wedle profilu |
| **unpaper** | `studio.ingest.recognize`, `research.source.ocr` | czyszczenie skanu **przed** Tesseractem |

**Wykaz zależności jest wspólny z dwoma innymi terenami biegnącymi teraz.**
Deklarację narzędzia zakładasz **przy miejscu użycia**, tak jak robi to rdzeń
(Pandoc przy dokumentach, ffmpeg przy nagraniach), a do
`zaleznosci_zewnetrzne.go` dopisujesz wyłącznie odwołanie. Przy scaleniu
rozjazd w tym jednym pliku rozstrzyga Prowadzący — nie jest to Twoja usterka.

### obraz-i-diagramy

Sześć programów do obrazu. Komendy istnieją. Rdzeń **wypuszcza źródło Mermaid**
(`AppExportFormatMermaid`), ale nie umie go narysować.

| | |
|---|---|
| **Gałąź** | `teren/obraz-i-diagramy` z `main` |
| **Wykaz plików** | `budowa/server/internal/core/adapter_narzedzia_obraz_*.go`, `adapter_modul_biblioteka_*.go`, `zaleznosci_zewnetrzne.go`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/shared/`, `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, **cały obszar `design.*` rdzenia**, adaptery dokumentów i przeglądarki |

| Narzędzie | Komenda |
|---|---|
| **mermaid-cli** (`mmdc`) | `design.diagram.render`, `apps.architecture.export` — **sprawdź zaporę fotografii, zanim tkniesz cokolwiek w `design.*`** |
| **optipng, jpegoptim, pngquant, cwebp** | `image.convert` — kompresja, której biblioteka wkompilowana nie robi |
| **exiftool** | `library.metadata.get`, `media.inspect` — **nie** `design.photo.metadata.get` ani `studio.security.metadata.strip`, te leżą za zaporami |

**Zapora fotografii jest napisana grubo — na wystąpienie nazwy silnika w treści
pliku.** `zapora_fotografii_test.go` zabrania w obszarze Design nazw `vips`,
`imagemagick`, `potrace`, `inkscape`, `fontforge`, `fonttools`, `graphicsmagick`.
Sprawdź uruchomieniem, czy przepuści `mmdc` i kompresory, **zanim** na tym
oprzesz pracę. Jeśli nie przepuści — to zgłoszenie, nie powód do jej zmiany.

**Wykaz zależności jest wspólny z dwoma innymi terenami biegnącymi teraz.**
Deklarację narzędzia zakładasz **przy miejscu użycia**, tak jak robi to rdzeń
(Pandoc przy dokumentach, ffmpeg przy nagraniach), a do
`zaleznosci_zewnetrzne.go` dopisujesz wyłącznie odwołanie. Przy scaleniu
rozjazd w tym jednym pliku rozstrzyga Prowadzący — nie jest to Twoja usterka.

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
