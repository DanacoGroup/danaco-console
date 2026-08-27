# Rejestr terenów

Żywy wykaz terenów. Prowadzi go wyłącznie Prowadzący budowę, na gałęzi `main`.
Teren bez wpisu w tym rejestrze nie jest otwarty, a praca na nim nie zostanie
przyjęta. Zasady podziału opisuje [ustrój budowy](ustroj-budowy.md).

## Tereny otwarte

### warunek-ukonczenia-zadania

Rdzeń zna trzy powody **zatrzymania** biegu — brak postępu, Operator, usterka
(`LoopStopReason`, `stan_obiegu.go:18`, stałe `session.PowodZatrzymania`) —
i żaden nie znaczy „ukończone z wynikiem". Pozycja 7 rejestru decyzji wymaga
rozróżnienia maszynowego, a pozycja otwarta „Warunek ukończenia zadania"
blokuje etap 2. Rozstrzygnięcie przyjęte do czasu rozstrzygnięcia Właściciela
niesie [rejestr decyzji](decyzje.md): czwarta wartość `completed`, **bez bramki
akceptacji** — bramka byłaby sprzeczna z zasadą zero blokad.

Zakotwiczenie sprawdzone przed otwarciem: wartości `complete` nie ma w żadnym
wyliczeniu tur, ale niesie ją `MessageStatus` (`pending`, `streaming`,
`complete`, `stopped`, `error`). Ustalenie, czym dokładnie jest „koordynator
skończył turę", należy do terenu i ma paść pomiarem, nie założeniem.

| | |
|---|---|
| **Gałąź** | `teren/warunek-ukonczenia-zadania` z `main` |
| **Wykaz plików** | `budowa/shared/contract.json` wraz z generatami, `budowa/server/internal/core/stan_obiegu.go`, `budowa/server/internal/session/obieg.go`, `petla.go`, sprawdziany tych pakietów |
| **Poza terenem** | `budowa/server/internal/core/urzadzenia_skaner.go` i `skutek_wydania_studia_test.go` (tereny biegnące), `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/` |

**Zmiana kontraktu jest dozwolona i obwarowana.** Suma zastana
`334705bd88c2efc13779`, 1080 komend. Wolno **dołożyć** wartość wyliczenia.
**Nie wolno** zmienić ani usunąć niczego istniejącego. Generator z
`budowa/shared/gen` musi dawać wynik bajtowo powtarzalny w dwóch przebiegach.
W rejestrze podajesz sumę zastaną i sumę po sobie.

**Kryteria odbioru.**

1. Ustalone **pomiarem**, nie założeniem, przy jakim stanie rdzenia koordynator
   kończy turę i jak sprawdza się, że żadne okno wykonawcze tury nie prowadzi —
   z przytoczonym miejscem w kodzie i wynikiem uruchomienia.
2. `LoopStopReason` niesie czwartą wartość `completed`, a rdzeń ustawia ją
   dokładnie w stanie z punktu 1 — wykazane sprawdzianem, który odróżnia
   ukończenie od trzech zatrzymań.
3. Bramki akceptacji nie ma — bieg po ukończeniu nie czeka na niczyje
   potwierdzenie; wykazane sprawdzianem.
4. Diff kontraktu zawiera wyłącznie dodania; generator bajtowo powtarzalny.
5. `gotestsum -- -count=1 ./...` — zero niepowodzeń wobec stanu zastanego
   podanego w raporcie, zmierzonego przed pierwszą zmianą.
6. Rewizje obejmują wyłącznie pliki terenu.

### odwolania-do-usunietych-skryptow

Sześć skryptów instalek natywnych usunięto z repozytorium pozycją 8 pkt 4
(teren `aktualizacja-powloki-i-skrypty`, rewizja `880b41f`). Rdzeń wciąż je
przywołuje: `budowa/server/internal/narzedzia/wpiecie.go:52` w komentarzu
i `:111` w treści odmowy — komunikat radzi Operatorowi budowę skryptem, którego
nie ma. Odmowa, która wskazuje nieistniejącą drogę naprawy, jest gorsza od
odmowy milczącej, bo wysyła po nic.

| | |
|---|---|
| **Gałąź** | `teren/odwolania-do-usunietych-skryptow` z `main` |
| **Wykaz plików** | `budowa/server/internal/narzedzia/wpiecie.go` wraz ze sprawdzianami tego pakietu; pozostałe pliki **wyłącznie** te, które przeszukanie wskaże jako przywołujące usunięte skrypty, i wyłącznie poza plikami terenów biegnących |
| **Poza terenem** | `budowa/shared/`, `budowa/server/internal/core/` (tereny biegnące), `budowa/klient/`, `budowa/desktop/`, `design/`, `prowadzenie/`, `budowa/scripts/` (skrypty żywe mają osobne zgłoszenie) |

**Kryteria odbioru.**

1. Przeszukanie całego repozytorium za nazwami sześciu usuniętych skryptów
   (`instalka-natywna-win.sh`, `instalka-natywna-linux.sh`,
   `instalka-windows.sh`, `instalka-pelna.sh`, `wydanie.sh`, `pakowanie.sh`)
   z przytoczonym poleceniem i pełnym wykazem trafień — potwierdzone, że wzorzec
   w ogóle łapie.
2. Każde trafienie w rdzeniu albo poprawione, albo wypisane wraz z powodem,
   dla którego zostaje (materiał zamknięty, dokument audytowy).
3. Odmowa `wpiecie.go` wskazuje drogę naprawy, która **istnieje** — wykazane
   sprawdzianem czytającym treść odmowy.
4. `gotestsum -- -count=1 ./...` — zero niepowodzeń wobec stanu zastanego
   podanego w raporcie.
5. Kontrakt nietknięty — suma `334705bd88c2efc13779`.
6. Rewizje obejmują wyłącznie pliki terenu.

## Zgłoszenia oczekujące na teren

Ustalenia z zamkniętych i biegnących terenów, które wykraczają poza ich zakres.
Każde zgłoszenie ma wskazany plik i wiersz. Zgłoszenie staje się terenem, gdy
Prowadzący je otworzy; do tego czasu jest wykazem, nie pracą.

### Granica 15 s w uprzęży kontraktu jest ciasna dla warstwy skanera

Ustalenie zmierzone przez teren `odmowy-skanera` i potwierdzone przez kontrolę.
Generyczna uprząż zgodności kontraktu woła każdą komendę rdzenia z twardym
limitem 15 s: `zgodnosc_kontraktu_test.go:215` i `:254`,
`blokady_skutek_test.go:127`. Tymczasem czynność skanera na tej maszynie
dochodzi do ~14,8 s (`scanimage --format=png` ~7,4 s + wykaz `scanimage -L`
~7,4 s), a własne limity warstwy są znacznie wyższe — `granicaWykazuUrzadzen`
45 s, `granicaSkanowaniaUrzadzenia` 5 min (`urzadzenia_skaner.go:63,68`). Pod
obciążeniem maszyny generyczne 15 s ucina czynność przed jej własnym limitem
i odmowa spada na kod czasu (`internal_error`) zamiast rozpoznać brak
urządzenia (`not_found`). Na maszynie budowlanej bez żywego skanera komendy
odmawiają natychmiast, więc chwiejność ujawnia się dopiero tam, gdzie warstwa
realnie sonduje urządzenie — u Operatora na wolnej maszynie ta sama granica
może uciąć legalną odmowę „brak urządzenia".

Teren obejmuje uprząż zgodności kontraktu (pliki poza terenami zaplecza).
Rozstrzygnięcie progu — podnieść do wartości spójnej z limitami warstwy —
należy rozważyć wraz z tym, jak uprząż ma traktować komendy o własnych,
dłuższych limitach.

### Rejestr rozjechany na gałęzi centrum poprawek — rozstrzygnięcie przy scaleniu

Gałąź `teren/centrum-poprawki` niesie dwie rewizje dotykające tego pliku —
`6361df8` i `31e37de` — założone tam omyłkowo przez poprzedniego Prowadzącego,
bo drzewo główne stoi przełączone na tę gałąź. Ich treść przeniesiono na `main`
jako `a52fa05` i `ea0c673`, ale na gałęzi zostały. Od tamtej pory `main`
przerobił ten plik mocno: pięć terenów zamkniętych, zgłoszenia i reguły odbioru.

**Przy scaleniu gałęzi designu konflikt w tym pliku rozstrzyga się na rzecz
`main`.** Wersja z gałęzi cofnęłaby zamknięcia terenów i zgłoszenia z pomiarów,
a wyglądałoby to na zwykłe scalenie. Praca w `design/` scala się normalnie —
rozstrzygnięcie dotyczy wyłącznie `prowadzenie/rejestr-terenow.md`.

### Kontrakt ruszył pierwszy raz od przejęcia

Do 27.08.2026 kontrakt stał nietknięty pod sumą `2cbb843d33f4531b05cd` — teren
`pomiar-stron` dołożył trzy osie pomiaru strony i suma wynosi dziś
`334705bd88c2efc13779` przy **1080 komendach** wobec 1077 zastanych. Zmiana
przeszła kontrolę porównaniem strukturalnym: nic istniejącego nie ubyło ani się
nie zmieniło. Odtąd sprawdzian nietykalności kontraktu odnosi się do sumy
bieżącej, nie do sumy z pierwszego dnia; kolejny teren, który kontrakt dokłada,
podaje w rejestrze sumę zastaną i sumę po sobie.

### Wdrożenie musi założyć trzy nastawy, inaczej stojące wagi leżą odłogiem

Ustalenie terenu `zaplecze-modeli`, zmierzone na żywym rdzeniu. Silnik wiedzy
działa na wagach z `/opt/danaco-modele/embedder`, ale **dopiero po nastawach** —
wartości domyślne rdzenia wskazują co innego i wdrożenie pobierze drugi model
zamiast użyć stojących 4,3 GB.

| Nastawa | Wartość | Bez niej |
|---|---|---|
| `wiedza_model` | `BAAI/bge-m3` | rdzeń sięga po `mpnet` z migracji 115 |
| `wiedza_katalog_modeli` | `/opt/danaco-modele/embedder` | wagi pobierane na nowo do katalogu danych |
| `mowa_katalog_modeli` | `/opt/danaco-modele/mowa` | wagi mowy stoją w pamięci podręcznej konta, które uruchomiło rdzeń |

Wag mowy dotyczy osobne ustalenie: 464 MB modelu `faster-whisper-small` stoi
dziś w pamięci podręcznej pod katalogiem domowym, bo `mowa/ustawienia.go` przy
pustej nastawie zostawia miejsce bibliotece, podczas gdy `wiedza/pomocnik.go`
przy pustej nastawie **przypina** wagi do katalogu danych rdzenia — droga
wołania procesu nie dziedziczy środowiska, więc pomocnik nie zna nawet `HOME`.
Skutek: czyszczenie pamięci podręcznej kasuje działającą funkcję bez śladu
w produkcie, a wagi są przywiązane do konta uruchamiającego. Rozstrzygnięcie
przyjęte: przenieść wagi do `/opt/danaco-modele/mowa` i wskazać je nastawą.

Do rozstrzygnięcia zostaje jedno: czy zrównać obie rodziny w kodzie, żeby pusta
nastawa mowy znaczyła to samo co pusta nastawa wiedzy. To zmiana zachowania
domyślnego wraz z migracją — `migracja_075_mowa.sql` niesie dziś wartość pustą
zgodną ze stałą co do znaku, więc ruszenie samej stałej stworzyłoby dwie prawdy.

### Trzy modele bez komendy w kontrakcie

Ustalenie terenu `zaplecze-modeli`. Wagi stoją, komend nie ma — rejestr zabronił
ich dokładania, bo kontrakt należał w tej turze do innego terenu.

| Model | Waga | Propozycja obszaru |
|---|---|---|
| reranker (bge-reranker-v2-m3) | 2,2 GB | `knowledge` — przełącznik `rerank` w `knowledge.search` albo osobna `knowledge.rerank` |
| CLIP | 1,6 GB | `library` albo `knowledge` — wyszukiwanie obrazów po znaczeniu |
| twarze (GFPGAN, codeformer) | 692 MB | pole `faces` w `image.upscale` **już istnieje**; brakuje silnika: wydanie ncnn nie niesie sieci twarzowej, a wagi `.pth` żądają stosu torch, którego rdzeń nie woła. Odmowa nazywa dziś brak `gfpgan-ncnn-vulkan` |

Osobno rozstrzygnięte i zamknięte: **Kokoro nie jest wart deklaracji** — 54 głosy,
żadnego polskiego (`pf`/`pm` to portugalski), a piper z `pl_PL-darkman-medium`
działa i rdzeń go zna.

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

**Pełny bieg wymaga podniesionego limitu czasu.** Pakiet `internal/core`
przekroczył domyślną granicę `go test` (10 minut) po dołożeniu przez cztery
tereny około czterdziestu sprawdzianów wołających prawdziwe programy zewnętrzne.
Bieg pada wtedy paniką „test timed out after 10m0s" wskazującą **przypadkowy**
sprawdzian, który akurat biegł — wygląda to na usterkę tego sprawdzianu i nią
nie jest. Objaw uboczny myli podwójnie: `gotestsum` naliczy wtedy około 947
sprawdzianów zamiast ponad dwóch tysięcy, bo pakiet, który spanikował, nie
policzy swoich. Pełny bieg uruchamia się z `-timeout 30m`.

**Pełne biegi szereguje się zamkiem.** Maszyna ma 16 rdzeni; przy czterech
terenach naraz biegło na niej dziewięć procesów sprawdzianów i sprawdziany
z twardymi granicami czasu zaczęły się chwiać — każdy pełny bieg gubił **inny**,
a pojedynczo wszystkie przechodziły. Granice tego rodzaju stoją m.in.
w `zgodnosc_kontraktu_test.go`. Bieg pełny wykonuje się pod zamkiem:

```bash
cd budowa/server && flock /tmp/danaco-bieg-pelny.lock \
  gotestsum -- -count=1 -timeout 30m ./...
```

Sprawdzian, który padł w pełnym biegu, a przechodzi uruchomiony pojedynczo, jest
chwiejnością pod obciążeniem, nie usterką — i tak się go nazywa.

**Pełny bieg regresji jest bramką scalenia u Prowadzącego, nie bramką wyjścia
u wykonawcy.** Wykonawca dowodzi swojego przedmiotu biegiem **celowanym** pakietu,
który dotyka — jest szybki i odporny na obciążenie. Pełny bieg `./...` w drzewie
terenu mierzy drzewo **bez** pozostałych terenów biegnących równolegle, więc
rozstrzyga wyłącznie bieg na scalonym `main` po scaleniu, a ten wykonuje
Prowadzący raz, pod zamkiem, jako warunek bramki scalenia. Zaprzęganie czterech
terenów do czterech pełnych biegów naraz zapycha zamek na godziny i nic nie
rozstrzyga — każdy z nich mierzy inny, niepełny stan.

Czas dostępu do pliku nie dowodzi, że pliku nie czytano. Drzewo stoi na `ext4`
zamontowanym z `relatime`, gdzie jądro odświeża czas dostępu wyłącznie wtedy,
gdy poprzedni jest starszy od czasu zmiany albo starszy niż doba. Odczyt pliku,
którego czas dostępu jest już późniejszy od czasu zmiany, **nie zostawia
śladu** — niezmieniony czas dostępu jest tam brakiem pomiaru, nie dowodem
nietknięcia. Wyszło to przy sporze o sprawstwo pobrania wag mowy: wykonawca
podał czas dostępu jako dowód, a instrument z założenia milczał.

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
| `pomiar-stron` | `teren/pomiar-stron` | `607a3d8` kontrakt · `4191c4b` trzy osie pomiaru · `6b189e7` sprawdziany | kontrola osobnej sesji: kontrakt porównany strukturalnie, nie liniowo — zero sekcji zmienionych, żaden istniejący element nietknięty bajtowo, dodane 3 komendy, 6 struktur, 3 wyliczenia; generator powtórzony dwa razy, sumy zgodne; oś obciążenia zmierzona niezależnie własnym klientem kanału — 102 666 żądań wobec 102 669 na liczniku serwera sprawdzianu; osiem sprawdzianów zdanych, żaden pominięty. Scalone `73e9aad` |
| `zaplecze-modeli` | `teren/zaplecze-modeli` | `090acd6` silnik osadzeń na wagach stojących | kontrola osobnej sesji: teren zwrócony za niepełny wykaz postawionego, po uzupełnieniu przyjęty; kontroler powtórzył pomiar osadzeń niezależnie i przy odciętej sieci (`HF_HUB_OFFLINE=1`) — wymiar 1024, normy 1,000000, cos(kot,kot) 0,8348 wobec 0,2863–0,3353 dla par odległych; wagi `/opt` czytane, nie pobierane po raz drugi; kontrakt nietknięty. Scalone `ba82cb6` |
| `obraz-i-diagramy` | `teren/obraz-i-diagramy` | `f7fab06` dogniecenie zapisu i metadane osadzone | kontrola osobnej sesji własnym biegiem: osiem zapór zdanych nietkniętych, 33 sprawdziany zdane w realnych czasach z przytoczonymi rozmiarami przed i po (PNG 3322→2456 B, JPEG 45399→39623 B, model barw paleta dowodzi wejścia pngquanta); sprawdziany maszyny bez programu wytwarzają ją naprawdę (`t.Setenv` na pusty katalog); kontrakt nietknięty; zero śladu mmdc w rdzeniu. Scalone `5b98585` |
| `dokumenty-i-tekst` | `teren/dokumenty-i-tekst` | `ea05603` sześć programów treści pisanej · `aa695e6` poprawka nazwy po zwrocie | kontrola osobnej sesji: teren zwrócony za martwą nazwę `zasiegSyntezy` w komentarzu, po poprawce przyjęty; własny bieg kontrolera 14 zdanych, zero pominiętych; sprawdziany mierzą skutek — PDF czytany drugą komendą i innym programem, korekta aż po treść panelu w bazie, OCR dwoma przebiegami, brak programu wywołany `t.Setenv`; wybory hunspell i typst zweryfikowane uruchomieniem; kontrakt nietknięty. Scalone `43d05f6` |
| `odmowy-skanera` | `teren/odmowy-skanera` | `1eaa58a` parytet odmów skanera Linux wobec Windows | kontrola osobnej sesji: `bladWarstwySane` wierne lustro `bladWarstwyWia` — ten sam kod `channel_unavailable`, ta sama droga obejścia `studio.ingest.queue.add`, pakiet czytany z braku nie zaszyty; sprawdzian parytetu woła obie warstwy i wymaga jednego kodu; brak wymuszony atrapą PATH z potwierdzeniem `zewnetrzne.Stoi=false`; bieg celowany 3 zdane; kontrakt nietknięty. Scalone `1dca877` |
| `sprawdziany-drogi-wejscia` | `teren/sprawdziany-drogi-wejscia` | `1393e74` sprawdzian zdjęcia znacznika bramki | weryfikacja Prowadzącego mutacją rdzenia: wyłączenie `zdejmijZnacznikBezPoczty` daje sprawdzian niezdany („znacznik przeżył swój powód"), przywrócenie — zdany; trzy sprawdziany rozjazdu z rdzeniem były już przerobione przez `brama-i-droga-wejscia`, wykonawca to zmierzył i nie tknął; bieg celowany osiem zdanych; kontrakt nietknięty, jeden plik terenu. Scalone `4e70ef1` |
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
